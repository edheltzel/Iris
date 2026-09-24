package quota_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func run(t *testing.T, workspace string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell companion")
	}
	command := exec.Command("sh", "wake.sh")
	command.Dir = "."
	command.Env = []string{"PATH=/usr/bin:/bin", "SPYNEL_WORKSPACE=" + workspace}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("wake.sh: %v: %s", err, output)
	}
	return string(output)
}

func TestEmptyWorkspaceAbsorbs(t *testing.T) {
	workspace := t.TempDir()
	output := run(t, workspace)
	if !strings.HasPrefix(output, "absorb\n") {
		t.Fatalf("output = %q", output)
	}
}

func TestWaitingTaskDispatchesWithoutQuotaCLI(t *testing.T) {
	workspace := t.TempDir()
	dir := filepath.Join(workspace, ".spynel", "tasks", "waiting")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "wait.md"), []byte("# Wait\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	output := run(t, workspace)
	if !strings.HasPrefix(output, "dispatch\n") || !strings.Contains(output, "quota-axi is not installed") {
		t.Fatalf("output = %q", output)
	}
}
