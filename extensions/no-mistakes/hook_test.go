package nomistakes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const stub = `#!/bin/sh
printf '%s\n' "$*" >> "$NM_CALLS"
exit 0
`

type harness struct {
	workspace string
	calls     string
	journal   string
	path      string
}

func setup(t *testing.T) harness {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell hook")
	}
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, "no-mistakes"), []byte(stub), 0o700); err != nil {
		t.Fatal(err)
	}
	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		t.Fatal(err)
	}
	return harness{
		workspace: workspace,
		calls:     filepath.Join(root, "calls"),
		journal:   filepath.Join(workspace, ".spynel", "extensions-state", "no-mistakes", "journal"),
		path:      bin + ":/usr/bin:/bin",
	}
}

func (h harness) run(t *testing.T, hook, payload string) (int, string, string) {
	t.Helper()
	command := exec.Command("./hook.sh", hook)
	command.Stdin = strings.NewReader(payload)
	command.Env = []string{
		"PATH=" + h.path,
		"SPYNEL_HOOK=" + hook,
		"SPYNEL_EXTENSION=no-mistakes",
		"SPYNEL_WORKSPACE=" + h.workspace,
		"NM_CALLS=" + h.calls,
	}
	var stdout, stderr strings.Builder
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	code := 0
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatal(err)
		}
		code = exit.ExitCode()
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	return code, stderr.String(), h.read(t, h.journal)
}

func (h harness) read(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func (h harness) mode(t *testing.T, mode string) {
	t.Helper()
	dir := filepath.Join(h.workspace, ".spynel")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "no-mistakes-mode"), []byte(mode+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

const done = `{"hook":"task.completed","payload":{"route":"developer","file":"tasks/add-cache.md","outcome":"done","event_id":"lease-1:done"}}`

func TestDoneTaskRefusesWhenCLIMissing(t *testing.T) {
	h := setup(t)
	h.path = "/usr/bin:/bin"
	code, stderr, journal := h.run(t, "task.completed", done)
	if code != 1 || !strings.Contains(stderr, "not installed") || !strings.Contains(journal, "action=fail") {
		t.Fatalf("code=%d stderr=%q journal=%q", code, stderr, journal)
	}
}

func TestHarnessAfterRefusesWhenCLIMissing(t *testing.T) {
	h := setup(t)
	h.path = "/usr/bin:/bin"
	code, stderr, _ := h.run(t, "harness.after", `{"hook":"harness.after","payload":{"kind":"final"}}`)
	if code != 1 || !strings.Contains(stderr, "not installed") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestCLIPresentDoesNotRunThePipeline(t *testing.T) {
	h := setup(t)
	code, stderr, journal := h.run(t, "task.completed", done)
	if code != 0 || !strings.Contains(stderr, "/no-mistakes") || !strings.Contains(journal, "action=pass") {
		t.Fatalf("code=%d stderr=%q journal=%q", code, stderr, journal)
	}
	if calls := h.read(t, h.calls); calls != "" {
		t.Fatalf("hook invoked the CLI: %q", calls)
	}
}

func TestDirectPRAndLocalOnlySkip(t *testing.T) {
	for _, mode := range []string{"direct-PR", "local-only"} {
		t.Run(mode, func(t *testing.T) {
			h := setup(t)
			h.path = "/usr/bin:/bin"
			h.mode(t, mode)
			code, _, journal := h.run(t, "task.completed", done)
			if code != 0 || !strings.Contains(journal, "action=skip") {
				t.Fatalf("code=%d journal=%q", code, journal)
			}
			if calls := h.read(t, h.calls); calls != "" {
				t.Fatalf("CLI called on skip: %q", calls)
			}
		})
	}
}

func TestRegistryPolicyIsNotATaskMode(t *testing.T) {
	h := setup(t)
	h.mode(t, "no-mistakes-prod-only")
	code, stderr, _ := h.run(t, "task.completed", done)
	if code != 1 || !strings.Contains(stderr, "registry policy") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	if calls := h.read(t, h.calls); calls != "" {
		t.Fatalf("policy mode invoked the CLI: %q", calls)
	}
}

func TestFailedOutcomeSkips(t *testing.T) {
	h := setup(t)
	h.path = "/usr/bin:/bin"
	payload := strings.Replace(done, `"done"`, `"failed"`, 1)
	code, _, journal := h.run(t, "task.completed", payload)
	if code != 0 || !strings.Contains(journal, "not a ship") {
		t.Fatalf("code=%d journal=%q", code, journal)
	}
}

func TestUnknownModeRefuses(t *testing.T) {
	h := setup(t)
	h.mode(t, "yolo")
	code, stderr, _ := h.run(t, "task.completed", done)
	if code != 1 || !strings.Contains(stderr, "unknown mode") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}
