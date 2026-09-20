package fsx

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestAtomicCreateFileNeverReplacesExistingContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "instruction.md")
	if err := AtomicCreateFile(path, []byte("owner edit"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := AtomicCreateFile(path, []byte("default"), 0o600); !errors.Is(err, os.ErrExist) {
		t.Fatalf("second create error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "owner edit" {
		t.Fatalf("existing content = %q, %v", data, err)
	}
}

func TestMigrateDirFallsBackWhenRenameFails(t *testing.T) {
	source := filepath.Join(t.TempDir(), "legacy")
	dest := filepath.Join(t.TempDir(), "current")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "value"), []byte("kept"), 0o600); err != nil {
		t.Fatal(err)
	}
	originalRename := renameDirectory
	calls := 0
	renameDirectory = func(oldPath, newPath string) error {
		calls++
		if calls == 1 {
			return syscall.EXDEV
		}
		return originalRename(oldPath, newPath)
	}
	t.Cleanup(func() { renameDirectory = originalRename })
	if err := MigrateDir(dest, source); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "nested", "value"))
	if err != nil || string(data) != "kept" || calls < 2 {
		t.Fatalf("migrated data = %q, calls = %d, err = %v", data, calls, err)
	}
	if _, err := os.Stat(source); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("legacy directory remains: %v", err)
	}
	if info, err := os.Stat(filepath.Join(dest, "nested", "value")); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("migrated mode = %v, err = %v", info, err)
	}
}

func TestMigrateDirAcceptsConcurrentCopyWinner(t *testing.T) {
	source := filepath.Join(t.TempDir(), "legacy")
	dest := filepath.Join(t.TempDir(), "current")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	originalRename := renameDirectory
	originalCopy := copyMigrationDirectory
	renameDirectory = func(string, string) error { return syscall.EXDEV }
	copyMigrationDirectory = func(string, string) error {
		if err := os.Mkdir(dest, 0o700); err != nil {
			return err
		}
		if err := os.RemoveAll(source); err != nil {
			return err
		}
		return &os.PathError{Op: "open", Path: source, Err: syscall.ENOENT}
	}
	t.Cleanup(func() {
		renameDirectory = originalRename
		copyMigrationDirectory = originalCopy
	})
	if err := MigrateDir(dest, source); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(dest); err != nil || !info.IsDir() {
		t.Fatalf("published destination = %v, err = %v", info, err)
	}
}
