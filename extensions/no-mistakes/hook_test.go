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
pwd >> "$NM_CALLS"
printf '%s\n' "$*" >> "$NM_CALLS"
case "$1" in
status)
	cat "$NM_STATUS" 2>/dev/null
	exit "${NM_STATUS_CODE:-0}"
	;;
esac
exit 0
`

type harness struct {
	workspace string
	calls     string
	status    string
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
		status:    filepath.Join(root, "status"),
		journal:   filepath.Join(workspace, ".spynel", "extensions-state", "no-mistakes", "journal"),
		path:      bin + ":/usr/bin:/bin",
	}
}

func (h harness) run(t *testing.T, hook, payload string, extra ...string) (int, string, string) {
	t.Helper()
	command := exec.Command("./hook.sh", hook)
	command.Stdin = strings.NewReader(payload)
	command.Env = append([]string{
		"PATH=" + h.path,
		"SPYNEL_HOOK=" + hook,
		"SPYNEL_EXTENSION=no-mistakes",
		"SPYNEL_WORKSPACE=" + h.workspace,
		"NM_CALLS=" + h.calls,
		"NM_STATUS=" + h.status,
	}, extra...)
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
	return code, stderr.String(), h.journalText(t)
}

func (h harness) journalText(t *testing.T) string {
	t.Helper()
	body, err := os.ReadFile(h.journal)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func (h harness) callsText(t *testing.T) string {
	t.Helper()
	body, err := os.ReadFile(h.calls)
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

func TestDoneTaskRefusesWhenGateMissing(t *testing.T) {
	h := setup(t)
	if err := os.WriteFile(h.status, []byte("daemon: running\ngate: missing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	code, stderr, _ := h.run(t, "task.completed", done)
	if code != 1 || !strings.Contains(stderr, "no no-mistakes gate") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestDoneTaskPassesWhenGatePresent(t *testing.T) {
	h := setup(t)
	if err := os.WriteFile(h.status, []byte("gate: /tmp/gate.git\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	code, _, journal := h.run(t, "task.completed", done)
	if code != 0 || !strings.Contains(journal, "action=pass") {
		t.Fatalf("code=%d journal=%q", code, journal)
	}
	calls := h.callsText(t)
	if !strings.Contains(calls, h.workspace) || !strings.Contains(calls, "status") {
		t.Fatalf("calls = %q", calls)
	}
}

func TestStatusFailureRefuses(t *testing.T) {
	h := setup(t)
	code, stderr, _ := h.run(t, "task.completed", done, "NM_STATUS_CODE=1")
	if code != 1 || !strings.Contains(stderr, "status failed") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestDirectPRSkipsWithoutCallingCLI(t *testing.T) {
	h := setup(t)
	h.mode(t, "direct-PR")
	code, _, journal := h.run(t, "task.completed", done)
	if code != 0 || !strings.Contains(journal, "action=skip") {
		t.Fatalf("code=%d journal=%q", code, journal)
	}
	if calls := h.callsText(t); calls != "" {
		t.Fatalf("CLI called on skip: %q", calls)
	}
}

func TestFailedOutcomeSkips(t *testing.T) {
	h := setup(t)
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

func TestProdOnlyRequiresGate(t *testing.T) {
	h := setup(t)
	h.mode(t, "no-mistakes-prod-only")
	h.path = "/usr/bin:/bin"
	code, stderr, _ := h.run(t, "task.completed", done)
	if code != 1 || !strings.Contains(stderr, "not installed") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestHarnessAfterJournalsMissingCLIWithoutFailingChat(t *testing.T) {
	h := setup(t)
	h.path = "/usr/bin:/bin"
	code, stderr, journal := h.run(t, "harness.after", `{"hook":"harness.after","payload":{"kind":"final","text":"done"}}`)
	if code != 0 || !strings.Contains(stderr, "not installed") || !strings.Contains(journal, "action=fail") {
		t.Fatalf("code=%d stderr=%q journal=%q", code, stderr, journal)
	}
}
