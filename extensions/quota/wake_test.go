package quota_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func workspaceWith(t *testing.T, rel, body string) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if body != "" {
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func runWatch(t *testing.T, workspace string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell companion")
	}
	command := exec.Command("sh", "watch.sh")
	command.Env = []string{"PATH=/usr/bin:/bin", "SPYNEL_WORKSPACE=" + workspace}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("watch.sh: %v: %s", err, output)
	}
	return string(output)
}

func TestUnchangedWorkflowAbsorbs(t *testing.T) {
	workspace := workspaceWith(t, ".spynel/tasks/working/job.md", "# Job\n")
	first := runWatch(t, workspace)
	if !strings.Contains(first, "wake: workflow files changed") {
		t.Fatalf("first poll = %q", first)
	}
	second := runWatch(t, workspace)
	if second != "" {
		t.Fatalf("unchanged poll = %q", second)
	}
}

func TestWorkingFileChangeWakes(t *testing.T) {
	workspace := workspaceWith(t, ".spynel/tasks/working/job.md", "# Job\n")
	runWatch(t, workspace)
	path := filepath.Join(workspace, ".spynel", "tasks", "working", "job.md")
	if err := os.WriteFile(path, []byte("# Job\nchanged\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if output := runWatch(t, workspace); !strings.Contains(output, "wake: workflow files changed") {
		t.Fatalf("changed poll = %q", output)
	}
}

func TestClaimHookRecordsEvidenceAndExitsZero(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell companion")
	}
	workspace := workspaceWith(t, ".spynel/tasks/working/job.md", "---\nid: job\n---\n# Job\n")
	file := filepath.Join(workspace, ".spynel", "tasks", "working", "job.md")
	command := exec.Command("sh", "hook.sh")
	command.Env = []string{
		"PATH=/usr/bin:/bin",
		"SPYNEL_HOOK=task.claimed",
		"SPYNEL_WORKSPACE=" + workspace,
	}
	command.Stdin = strings.NewReader(`{"hook":"task.claimed","payload":{"file":"` + file + `"}}`)
	var stdout strings.Builder
	command.Stdout = &stdout
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q", stdout.String())
	}
	body, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := filepath.Join(workspace, ".spynel", "extensions-state", "quota", "latest.txt")
	if !strings.Contains(string(body), "Quota evidence: "+snapshot) {
		t.Fatalf("task body = %q", body)
	}
	recorded, err := os.ReadFile(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(recorded), "quota-axi is not installed") {
		t.Fatalf("snapshot = %q", recorded)
	}
	again := exec.Command("sh", "hook.sh")
	again.Env = command.Env
	again.Stdin = strings.NewReader(`{"hook":"task.claimed","payload":{"file":"` + file + `"}}`)
	if err := again.Run(); err != nil {
		t.Fatal(err)
	}
	reread, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(reread), "Quota evidence:") != 1 {
		t.Fatalf("evidence repeated: %s", reread)
	}
}
