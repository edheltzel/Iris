package media

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const speechCacheVersion = "v1"

// SpeechCacheDir resolves and creates the stable per-user namespace for
// automatically managed speech assets. Composition injects this path once.
func SpeechCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		if err == nil {
			err = errors.New("platform returned an empty path")
		}
		return "", fmt.Errorf("determine operating-system user home directory for automatic speech assets: %w; configure speech.model_dir explicitly to avoid automatic model provisioning", err)
	}
	root := filepath.Join(home, ".agents", "Iris", "speech", speechCacheVersion, "parakeet")
	var sources []string
	if cache, err := os.UserCacheDir(); err == nil && strings.TrimSpace(cache) != "" {
		sources = []string{
			filepath.Join(cache, "iris", "speech", speechCacheVersion, "parakeet"),
			filepath.Join(cache, "spynel", "speech", speechCacheVersion, "parakeet"),
		}
	}
	if err := migrateInto(root, sources...); err != nil {
		return "", fmt.Errorf("migrate speech cache to %q: %w", root, err)
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		return "", fmt.Errorf("create shared speech cache %q: %w; check directory permissions or configure speech.model_dir explicitly", root, err)
	}
	return root, nil
}

func migrateInto(dest string, sources ...string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return err
	}
	if _, err := os.Lstat(dest); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	for _, src := range sources {
		info, err := os.Lstat(src)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if !info.IsDir() {
			continue
		}
		err = os.Rename(src, dest)
		if err == nil {
			return nil
		}
		if _, destErr := os.Lstat(dest); destErr == nil {
			return nil
		}
		if os.IsNotExist(err) {
			if _, destErr := os.Lstat(dest); destErr == nil {
				return nil
			}
			continue
		}
		return err
	}
	return nil
}
