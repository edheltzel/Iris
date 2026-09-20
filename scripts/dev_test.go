package scripts

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInstallDevRemovesOnlyOwnedLegacySymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX development installer")
	}
	directory := t.TempDir()
	installer := filepath.Join(directory, "install-dev.sh")
	source, err := filepath.Abs("install-dev.sh")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(source, installer); err != nil {
		t.Fatal(err)
	}
	built := filepath.Join(directory, "built-iris")
	if err := os.WriteFile(built, []byte("#!/bin/sh\nexit 0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "dev.sh"), []byte("#!/bin/sh\nprintf '%s\\n' \"$DEV_BUILD\"\n"), 0o700); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name        string
		setup       func(string) error
		wantRemoved bool
		wantWarning bool
	}{
		{
			name: "owned symlink",
			setup: func(path string) error {
				return os.Symlink(built, path)
			},
			wantRemoved: true,
		},
		{
			name: "unowned file",
			setup: func(path string) error {
				return os.WriteFile(path, []byte("unrelated\n"), 0o700)
			},
			wantWarning: true,
		},
		{
			name: "unowned symlink",
			setup: func(path string) error {
				other := filepath.Join(filepath.Dir(path), "other")
				if err := os.WriteFile(other, []byte("unrelated\n"), 0o700); err != nil {
					return err
				}
				return os.Symlink(other, path)
			},
			wantWarning: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			bin := t.TempDir()
			legacy := filepath.Join(bin, "spynel")
			if err := test.setup(legacy); err != nil {
				t.Fatal(err)
			}
			command := exec.Command("sh", installer, "--bin-dir", bin)
			command.Env = append(os.Environ(), "DEV_BUILD="+built)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("install-dev failed: %v output=%q", err, output)
			}
			_, statErr := os.Lstat(legacy)
			if test.wantRemoved != os.IsNotExist(statErr) {
				t.Fatalf("legacy target removal = %t, want %t: %v", os.IsNotExist(statErr), test.wantRemoved, statErr)
			}
			if strings.Contains(string(output), "leaving unowned legacy command") != test.wantWarning {
				t.Fatalf("warning output = %q, want warning %t", output, test.wantWarning)
			}
		})
	}
}

func TestColdCacheHelperCleansCancelledRunOutsideWorkspace(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux procfs cancellation contract")
	}
	tempBase := t.TempDir()
	cachePathFile := filepath.Join(tempBase, "cache-path")
	childPIDFile := filepath.Join(tempBase, "child-pid")
	descendantPIDFile := filepath.Join(tempBase, "descendant-pid")
	command := exec.Command("sh", "cold-cache.sh", "sh", "-c", `printf '%s' "$GOCACHE" >"$1"; printf '%s' "$$" >"$2"; trap '' TERM; (trap '' TERM; sleep 2; mkdir -p "$GOCACHE/recreated"; while :; do sleep 1; done) & printf '%s' "$!" >"$3"; while :; do sleep 1; done`, "sh", cachePathFile, childPIDFile, descendantPIDFile)
	command.Env = append(os.Environ(), "TMPDIR="+tempBase)
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if command.Process != nil {
			_ = command.Process.Kill()
		}
	})

	var cacheDir string
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(cachePathFile)
		if err == nil && len(data) != 0 {
			cacheDir = string(data)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if cacheDir == "" {
		t.Fatal("cold-cache child did not publish its cache path")
	}
	if err := exec.Command("kill", "-TERM", strconv.Itoa(command.Process.Pid)).Run(); err != nil {
		t.Fatalf("send TERM to cold-cache helper: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancelled cold-cache helper unexpectedly succeeded")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cold-cache helper did not exit promptly after TERM")
	}
	assertColdCacheRemoved(t, cacheDir, childPIDFile, descendantPIDFile)
	time.Sleep(2500 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(cacheDir, "recreated")); !os.IsNotExist(err) {
		t.Fatalf("cold-cache descendant recreated the cache after cancellation: %v", err)
	}
}

func TestColdCacheHelperCleansFailedRunOutsideWorkspace(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux procfs ownership contract")
	}
	tempBase := t.TempDir()
	command := exec.Command("sh", "cold-cache.sh", "sh", "-c", `printf '%s' "$GOCACHE"; exit 7`)
	command.Env = append(os.Environ(), "TMPDIR="+tempBase)
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatal("cold-cache helper unexpectedly succeeded")
	}
	cacheDir := string(output)
	if !strings.HasPrefix(cacheDir, tempBase+string(filepath.Separator)+"spynel-gocache.") {
		t.Fatalf("cold cache path %q is not a unique child of %q", cacheDir, tempBase)
	}
	if _, statErr := os.Stat(cacheDir); !os.IsNotExist(statErr) {
		t.Fatalf("cold cache was not removed after failure: %v", statErr)
	}

	projectRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	rejected := exec.Command("sh", "cold-cache.sh", "true")
	rejected.Env = append(os.Environ(), "TMPDIR="+projectRoot)
	output, err = rejected.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "must be outside the project workspace") {
		t.Fatalf("cold-cache helper accepted workspace TMPDIR: err=%v output=%q", err, output)
	}
}

func TestColdCacheHelperCleansCompletedRunWithDescendant(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux procfs process-tree contract")
	}
	tempBase := t.TempDir()
	command := exec.Command("sh", "cold-cache.sh", "sh", "-c", `printf '%s' "$GOCACHE"; (trap '' TERM; sleep 2; mkdir -p "$GOCACHE/recreated"; while :; do sleep 1; done) &`)
	command.Env = append(os.Environ(), "TMPDIR="+tempBase)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("completed cold-cache diagnostic failed: %v output=%q", err, output)
	}
	cacheDir := string(output)
	if !strings.HasPrefix(cacheDir, tempBase+string(filepath.Separator)+"spynel-gocache.") {
		t.Fatalf("cold cache path %q is not a unique child of %q", cacheDir, tempBase)
	}
	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		t.Fatalf("cold cache was not removed after success: %v", err)
	}
	time.Sleep(2500 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(cacheDir, "recreated")); !os.IsNotExist(err) {
		t.Fatalf("cold-cache descendant recreated the cache after successful command exit: %v", err)
	}
}

func TestColdCacheHelperCleansCompletedRunWithEscapedSession(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux procfs ownership contract")
	}
	tempBase := t.TempDir()
	cachePathFile := filepath.Join(tempBase, "cache-path")
	descendantPIDFile := filepath.Join(tempBase, "descendant-pid")
	command := exec.Command("sh", "cold-cache.sh", "sh", "-c", `printf '%s' "$GOCACHE" >"$1"; setsid sh -c 'trap "" TERM; printf "%s" "$$" >"$1"; sleep 2; mkdir -p "$GOCACHE/escaped-recreated"; while :; do sleep 1; done' sh "$2" >/dev/null 2>&1 & while [ ! -s "$2" ]; do sleep 0.01; done`, "sh", cachePathFile, descendantPIDFile)
	command.Env = append(os.Environ(), "TMPDIR="+tempBase)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("completed cold-cache diagnostic failed: %v output=%q", err, output)
	}
	cacheDir, err := os.ReadFile(cachePathFile)
	if err != nil {
		t.Fatal(err)
	}
	assertColdCacheRemoved(t, string(cacheDir), descendantPIDFile)
	time.Sleep(2500 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(string(cacheDir), "escaped-recreated")); !os.IsNotExist(err) {
		t.Fatalf("escaped cold-cache descendant recreated the cache: %v", err)
	}
}

func assertColdCacheRemoved(t *testing.T, cacheDir string, pidFiles ...string) {
	t.Helper()
	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		t.Fatalf("cold cache was not removed: %v", err)
	}
	for _, pidFile := range pidFiles {
		pid, err := os.ReadFile(pidFile)
		if err != nil {
			t.Fatal(err)
		}
		state, err := exec.Command("ps", "-o", "stat=", "-p", string(pid)).Output()
		if err == nil && !strings.HasPrefix(strings.TrimSpace(string(state)), "Z") {
			t.Fatalf("cold-cache process %s survived with state %q", pid, state)
		}
	}
}
