#!/usr/bin/env python3
"""Native isolated Spynel-to-Iris upgrade check. Pass prior and current archives."""
import hashlib
import http.server
import json
import os
from pathlib import Path
import re
import select
import shutil
import subprocess
import sys
import tempfile
import threading


def archive_details(archive, product):
    pattern = rf"{product}_([0-9.]+)_(linux|darwin)_(amd64|arm64)\.tar\.gz"
    match = re.fullmatch(pattern, archive.name)
    assert match is not None, archive.name
    return match.groups()


def main():
    older, newer = [Path(path).resolve() for path in sys.argv[1:]]
    old_version, target_os, target_arch = archive_details(older, "spynel")
    new_version, new_os, new_arch = archive_details(newer, "iris")
    assert (target_os, target_arch) == (new_os, new_arch)
    repo = Path(__file__).resolve().parent.parent
    script = (repo / "install.sh").read_bytes()
    with tempfile.TemporaryDirectory(prefix=".tmp-standalone-", dir=repo) as temporary:
        temp = Path(temporary).resolve()
        install = temp / "installed runtime Ω"
        user_bin = temp / "user bin"
        user_bin.mkdir()
        workspace = temp / "workspace Ω"
        workspace.mkdir()
        user_home = temp / "home"
        user_home.mkdir()
        tools = temp / "system tools"
        tools.mkdir()
        for name in ("sh", "curl", "tar", "awk", "mktemp", "uname", "sha256sum", "shasum", "wc", "mkdir", "rm", "rmdir", "chmod", "readlink", "ln", "gzip", "id", "sed", "grep", "cat", "cmp"):
            source = shutil.which(name)
            if source:
                (tools / name).symlink_to(source)
        env = {key: value for key, value in os.environ.items() if not key.startswith("SPYNEL_")}
        env.update(
            HOME=str(user_home),
            SHELL="/bin/bash",
            PATH=str(tools),
            XDG_RUNTIME_DIR=str(temp / "isolated runtime"),
            XDG_CONFIG_HOME=str(user_home / ".config"),
            SPYNEL_INSTALL_DIR=str(install),
            SPYNEL_BIN_DIR=str(user_bin),
            SPYNEL_VERSION=new_version,
        )
        files = {path.name: path.read_bytes() for path in (older, newer)}
        sums = "".join(f"{hashlib.sha256(data).hexdigest()}  {name}\n" for name, data in files.items()).encode()
        checksums = temp / "checksums.txt"
        checksums.write_bytes(sums)
        fixture = {"failure": ""}
        download_started = threading.Event()
        release_download = threading.Event()

        class Handler(http.server.BaseHTTPRequestHandler):
            def log_message(self, *args):
                pass

            def do_GET(self):
                if self.path == "/install.sh":
                    data = script
                elif self.path == "/uninstall.sh":
                    data = uninstall_script
                elif self.path == "/latest":
                    data = json.dumps({"tag_name": "v" + new_version, "prerelease": False, "draft": False}).encode()
                elif self.path == "/checksums.txt":
                    data = sums if fixture["failure"] != "checksum" else re.sub(rb"^[0-9a-f]{64}", b"0" * 64, sums, flags=re.M)
                elif self.path[1:] in files:
                    data = files[self.path[1:]]
                    if fixture["failure"] == "incomplete":
                        self.send_response(200)
                        self.send_header("Content-Length", str(len(data)))
                        self.end_headers()
                        self.wfile.write(data[:1024])
                        self.close_connection = True
                        return
                else:
                    self.send_error(404)
                    return
                self.send_response(200)
                self.send_header("Content-Length", str(len(data)))
                self.end_headers()
                if fixture["failure"] == "pause" and self.path[1:] in files:
                    download_started.set()
                    release_download.wait(15)
                self.wfile.write(data)

        server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), Handler)
        worker = threading.Thread(target=server.serve_forever, daemon=True)
        worker.start()
        base = f"http://127.0.0.1:{server.server_port}"
        env.update(
            SPYNEL_DOWNLOAD_BASE=base,
            SPYNEL_GITHUB_API_URL=base + "/latest",
            SPYNEL_NPM_REGISTRY_URL=base + "/forbidden-npm",
        )
        uninstall_script = (repo / "uninstall.sh").read_bytes().replace(
            b"https://spynel.agent-zero.ai/install.sh", (base + "/install.sh").encode()
        ).replace(b"=https", b"=http,https")
        try:
            immediate = temp / "immediate installation"
            immediate_env = {**env, "SPYNEL_INSTALL_DIR": str(immediate)}
            immediate_env.pop("SPYNEL_BIN_DIR")
            fixture["failure"] = "pause"
            child = subprocess.Popen(
                ["sh", "-c", 'curl -LsSf "$1/install.sh" | sh && iris --version', "sh", base],
                cwd=workspace,
                env=immediate_env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            try:
                assert download_started.wait(10), "installer did not start downloading"
                assert select.select([child.stderr], [], [], 5)[0], "installer is silent during the download"
                first_line = child.stderr.readline()
                assert b"Downloading Iris" in first_line, first_line
                assert child.poll() is None, "fixture did not hold the download open"
            finally:
                release_download.set()
                output, progress = child.communicate(timeout=120)
            assert child.returncode == 0, progress.decode()
            assert ("iris " + new_version).encode() in output, output.decode()
            assert b"Run: iris" in output and b"open a new terminal" not in output
            assert b"%" in progress and b"Verifying" in progress and b"Installing Iris" in progress, progress.decode()
            assert not (user_home / ".bashrc").exists(), "on-PATH install edited shell profiles"
            fixture["failure"] = ""
            removed = subprocess.run(
                ["sh", "-c", 'curl -LsSf "$1/uninstall.sh" | sh && ! command -v iris', "sh", base],
                cwd=workspace,
                env=immediate_env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=120,
            )
            assert removed.returncode == 0, removed.stderr.decode()
            assert not immediate.exists() and not (tools / "iris").is_symlink()

            if target_os in ("linux", "darwin"):
                with tempfile.TemporaryDirectory(prefix="iris-user-test-") as regular_temporary:
                    nonroot = Path(regular_temporary).resolve()
                    nonroot_home = nonroot / "home"
                    nonroot_home.mkdir()
                    readonly_tools = nonroot / "tools"
                    readonly_tools.mkdir()
                    for tool in tools.iterdir():
                        (readonly_tools / tool.name).symlink_to(tool.resolve())
                    sudo = readonly_tools / "sudo"
                    sudo.write_text('#!/bin/sh\necho "$*" >> "$HOME/sudo.calls"\n[ "$1" != -v ] || exit 0\nchmod 755 "$PATH"\nexec "$@"\n')
                    sudo.chmod(0o755)
                    identity = {}
                    if os.geteuid() == 0:
                        identity = {"user": 65534, "group": 65534}
                        for directory in (nonroot, nonroot_home, readonly_tools):
                            os.chown(directory, 65534, 65534)
                    temp.chmod(0o755)
                    readonly_tools.chmod(0o555)
                    regular_env = {
                        **env,
                        "HOME": str(nonroot_home),
                        "PATH": str(readonly_tools),
                        "SPYNEL_INSTALL_DIR": str(nonroot / "installation"),
                    }
                    regular_env.pop("SPYNEL_BIN_DIR")
                    command = 'before=$PATH; curl -LsSf "$1/install.sh" | sh && [ "$PATH" = "$before" ] && iris --version && curl -LsSf "$1/uninstall.sh" | sh && hash -r && ! command -v iris'
                    result = subprocess.run(
                        ["sh", "-c", command, "sh", base],
                        cwd=nonroot,
                        env=regular_env,
                        **identity,
                        stdout=subprocess.PIPE,
                        stderr=subprocess.PIPE,
                        timeout=180,
                    )
                    assert result.returncode == 0, result.stderr.decode()
                    assert ("iris " + new_version).encode() in result.stdout
                    assert "-v" in (nonroot_home / "sudo.calls").read_text()
                    assert not (nonroot / "installation").exists()

            old_runtime = temp / "prior runtime"
            old_runtime.mkdir()
            subprocess.run(["tar", "-xzf", str(older), "-C", str(old_runtime)], env=env, check=True, timeout=120)
            old_binary = old_runtime / "spynel"
            seeded = subprocess.run(
                [
                    str(old_binary),
                    "install-bundle",
                    "--root",
                    str(install),
                    "--archive",
                    str(older),
                    "--checksums",
                    str(checksums),
                    "--version",
                    old_version,
                ],
                cwd=workspace,
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=120,
            )
            assert seeded.returncode == 0, seeded.stderr.decode()
            legacy_launcher = install / "spynel"
            legacy_executable = legacy_launcher.resolve()
            (user_bin / "spynel").symlink_to(legacy_launcher)
            (install / ".bin-dir").write_text(str(user_bin) + "\n")
            old_output = subprocess.run(
                [str(legacy_launcher), "--version"],
                cwd=workspace,
                env=env,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=30,
            )
            assert old_output.returncode == 0 and old_output.stdout.strip() == "spynel " + old_version, old_output.stderr

            state = workspace / ".spynel" / "tasks" / "todo"
            state.mkdir(parents=True)
            sentinel = state / "preserved.md"
            sentinel.write_text("workspace state stays here\n")
            for failure in ("checksum", "incomplete"):
                fixture["failure"] = failure
                rejected = subprocess.run(
                    ["sh"],
                    input=script,
                    cwd=workspace,
                    env=env,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    timeout=120,
                )
                assert rejected.returncode != 0, "upgrade accepted a corrupt Iris archive"
                assert legacy_launcher.resolve() == legacy_executable
                assert not (install / "iris").exists()
            fixture["failure"] = ""

            upgraded = subprocess.run(
                ["sh"],
                input=script,
                cwd=workspace,
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=120,
            )
            assert upgraded.returncode == 0, upgraded.stderr.decode()
            iris_launcher = user_bin / "iris"
            version = subprocess.run(
                [str(iris_launcher), "--version"],
                cwd=workspace,
                env=env,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=30,
            )
            assert version.returncode == 0 and version.stdout.strip() == "iris " + new_version, version.stderr
            assert not legacy_launcher.exists() and not legacy_launcher.is_symlink()
            assert not (user_bin / "spynel").exists() and not (user_bin / "spynel").is_symlink()
            assert legacy_executable.exists(), "running prior-generation libraries were removed"
            members = {
                name.removeprefix("./").rstrip("/")
                for name in subprocess.check_output(["tar", "-tzf", str(newer)], text=True).splitlines()
            }
            assert "iris" in members and "spynel" not in members
            assert sentinel.read_text() == "workspace state stays here\n"
            profile_check = subprocess.run(
                ["sh", "-c", '. "$HOME/.bashrc"; iris --version'],
                cwd=workspace,
                env=env,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=30,
            )
            assert profile_check.returncode == 0 and profile_check.stdout.strip() == "iris " + new_version, profile_check.stderr

            removed = subprocess.run(
                ["sh"],
                input=uninstall_script,
                cwd=workspace,
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=120,
            )
            assert removed.returncode == 0, removed.stderr.decode()
            assert not install.exists() and not iris_launcher.exists()
            assert sentinel.read_text() == "workspace state stays here\n"
            print(
                json.dumps(
                    {
                        "result": "passed",
                        "classification": "observed-native",
                        "target": f"{target_os}/{target_arch}",
                        "checks": [
                            "current Iris bootstrap and uninstall",
                            "live download progress",
                            "regular-user immediate command with modeled sudo",
                            "corrupt Iris upgrade rejected",
                            "owned Spynel installation upgraded without alias",
                            "Iris archive contains only the Iris command",
                            "prior bundle and workspace state preserved",
                        ],
                    }
                )
            )
        finally:
            server.shutdown()
            server.server_close()
            worker.join()


if __name__ == "__main__":
    main()
