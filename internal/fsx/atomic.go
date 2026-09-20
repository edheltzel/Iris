package fsx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var renameDirectory = os.Rename

// AtomicWriteFile durably writes a temporary sibling and replaces path with it.
func AtomicWriteFile(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".spynel-write-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(mode); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return replace(tempName, path)
}

// AtomicCreateFile publishes a complete sibling temporary file only when path
// does not already exist. The hard-link publication cannot replace a manual or
// concurrent writer's file.
func AtomicCreateFile(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".spynel-create-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(mode); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Link(tempName, path)
}

// MigrateDir moves the first existing source directory into dest without
// requiring the source and destination to share a filesystem.
func MigrateDir(dest string, sources ...string) error {
	parent := filepath.Dir(dest)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return err
	}
	if _, err := os.Lstat(dest); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	for _, source := range sources {
		info, err := os.Lstat(source)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if !info.IsDir() {
			continue
		}
		if err := renameDirectory(source, dest); err == nil {
			return nil
		}
		if _, err := os.Lstat(dest); err == nil {
			return nil
		}
		if _, err := os.Lstat(source); os.IsNotExist(err) {
			continue
		}
		temp, err := os.MkdirTemp(parent, ".iris-migrate-")
		if err != nil {
			return err
		}
		removeTemp := true
		defer func() {
			if removeTemp {
				_ = os.RemoveAll(temp)
			}
		}()
		if err := copyDirectory(source, temp); err != nil {
			return err
		}
		if err := renameDirectory(temp, dest); err != nil {
			if _, destErr := os.Lstat(dest); destErr == nil {
				return nil
			}
			return err
		}
		removeTemp = false
		return os.RemoveAll(source)
	}
	return nil
}

func copyDirectory(source, dest string) error {
	root, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if err := os.Chmod(dest, root.Mode().Perm()); err != nil {
		return err
	}
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil || relative == "." {
			return err
		}
		target := filepath.Join(dest, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		switch {
		case info.IsDir():
			if err := os.Mkdir(target, info.Mode().Perm()); err != nil {
				return err
			}
			return os.Chmod(target, info.Mode().Perm())
		case info.Mode().IsRegular():
			return copyRegularFile(path, target, info)
		case info.Mode()&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		default:
			return fmt.Errorf("cannot migrate unsupported file %q", path)
		}
	})
}

func copyRegularFile(source, dest string, expected os.FileInfo) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	actual, err := input.Stat()
	if err != nil {
		return err
	}
	if !actual.Mode().IsRegular() || !os.SameFile(expected, actual) {
		return fmt.Errorf("source changed during migration: %q", source)
	}
	output, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, expected.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err = io.Copy(output, input); err == nil {
		err = output.Sync()
	}
	if closeErr := output.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Chmod(dest, expected.Mode().Perm())
}
