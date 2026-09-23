package gitbutler_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// stub records every `but` invocation and answers the three subcommands the
// hook depends on: status, branch list, and commit.
const stub = `#!/bin/sh
printf '%s\n' "$*" >> "$BUT_CALLS"
case "$1 $2" in
"branch list") cat "$BUT_BRANCHES" 2>/dev/null; exit 0 ;;
"branch new") printf '  "name": "%s",\n' "$3" >> "$BUT_BRANCHES"; exit 0 ;;
esac
case "$1" in
status) exit 0 ;;
commit) exit "${BUT_COMMIT_STATUS:-0}" ;;
esac
exit 0
`

type harness struct {
	workspace string
	calls     string
	branches  string
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
	if err := os.WriteFile(filepath.Join(bin, "but"), []byte(stub), 0o700); err != nil {
		t.Fatal(err)
	}
	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		t.Fatal(err)
	}
	return harness{
		workspace: workspace,
		calls:     filepath.Join(root, "calls"),
		branches:  filepath.Join(root, "branches"),
		path:      bin + ":/usr/bin:/bin",
	}
}

func (h harness) run(t *testing.T, hook, payload string, extra ...string) {
	t.Helper()
	command := exec.Command("./hook.sh", hook)
	command.Stdin = strings.NewReader(payload)
	command.Env = append(os.Environ(),
		"PATH="+h.path,
		"SPYNEL_HOOK="+hook,
		"SPYNEL_EXTENSION=gitbutler",
		"SPYNEL_WORKSPACE="+h.workspace,
		"BUT_CALLS="+h.calls,
		"BUT_BRANCHES="+h.branches,
	)
	command.Env = append(command.Env, extra...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("%s hook exited nonzero: %v: %s", hook, err, output)
	}
}

func (h harness) butCalls(t *testing.T) string {
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

const claimed = `{"hook":"task.claimed","payload":{"route":"developer","phase":"work","file":"tasks/add-cache.md","id":"t1"}}`
const completed = `{"hook":"task.completed","payload":{"route":"developer","file":"tasks/add-cache.md","thread_id":"th1","outcome":"done","event_id":"lease-1:done"}}`

func TestClaimedTaskCreatesBranchOnce(t *testing.T) {
	h := setup(t)
	h.run(t, "task.claimed", claimed)
	if calls := h.butCalls(t); !strings.Contains(calls, "branch new iris/add-cache") {
		t.Fatalf("but calls = %q", calls)
	}
	if err := os.Truncate(h.calls, 0); err != nil {
		t.Fatal(err)
	}
	h.run(t, "task.claimed", claimed)
	if calls := h.butCalls(t); strings.Contains(calls, "branch new") {
		t.Fatalf("existing branch was recreated: %q", calls)
	}
}

func TestClaimedTaskDoesNotMatchLongerBranchName(t *testing.T) {
	h := setup(t)
	if err := os.WriteFile(h.branches, []byte("  \"name\": \"iris/add-cache-later\",\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	h.run(t, "task.claimed", claimed)
	if calls := h.butCalls(t); !strings.Contains(calls, "branch new iris/add-cache") {
		t.Fatalf("branch with a longer namesake was treated as existing: %q", calls)
	}
}

func TestCompletedTaskCommitsOncePerEvent(t *testing.T) {
	h := setup(t)
	h.run(t, "task.claimed", claimed)
	h.run(t, "task.completed", completed)
	calls := h.butCalls(t)
	if !strings.Contains(calls, "commit --branch iris/add-cache -m iris(developer): add-cache done") {
		t.Fatalf("but calls = %q", calls)
	}
	if err := os.Truncate(h.calls, 0); err != nil {
		t.Fatal(err)
	}
	h.run(t, "task.completed", completed)
	if calls := h.butCalls(t); strings.Contains(calls, "commit") {
		t.Fatalf("redelivered event committed again: %q", calls)
	}
}

func TestFailedCommitLeavesEventRetryable(t *testing.T) {
	h := setup(t)
	h.run(t, "task.claimed", claimed)
	h.run(t, "task.completed", completed, "BUT_COMMIT_STATUS=1")
	if err := os.Truncate(h.calls, 0); err != nil {
		t.Fatal(err)
	}
	h.run(t, "task.completed", completed)
	if calls := h.butCalls(t); !strings.Contains(calls, "commit --branch iris/add-cache") {
		t.Fatalf("retry did not commit: %q", calls)
	}
}

func TestHookSucceedsWithoutGitButler(t *testing.T) {
	h := setup(t)
	h.path = "/usr/bin:/bin"
	h.run(t, "task.claimed", claimed)
	if calls := h.butCalls(t); calls != "" {
		t.Fatalf("but was invoked without the CLI installed: %q", calls)
	}
}

func TestHookSucceedsWithoutWorkspace(t *testing.T) {
	h := setup(t)
	command := exec.Command("./hook.sh", "task.claimed")
	command.Stdin = strings.NewReader(claimed)
	command.Env = append(os.Environ(), "PATH="+h.path, "SPYNEL_HOOK=task.claimed", "SPYNEL_WORKSPACE=", "BUT_CALLS="+h.calls)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("hook exited nonzero without a workspace: %v: %s", err, output)
	}
}
