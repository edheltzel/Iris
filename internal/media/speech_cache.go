package media

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/edheltzel/iris/internal/fsx"
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
	if err := fsx.MigrateDir(root, sources...); err != nil {
		return "", fmt.Errorf("migrate speech cache to %q: %w", root, err)
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		return "", fmt.Errorf("create shared speech cache %q: %w; check directory permissions or configure speech.model_dir explicitly", root, err)
	}
	return root, nil
}
