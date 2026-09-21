package startup

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/edheltzel/iris/internal/fsx"
)

// RemoveInstallation stops and removes only registrations for this executable
// or npm launcher. Ordinary preference changes retain Sync's non-stopping behavior.
func (m *Manager) RemoveInstallation(ctx context.Context, userID int) error {
	return m.stopInstallation(ctx, userID, true)
}

// StopInstallation unloads services without changing future startup settings.
func (m *Manager) StopInstallation(ctx context.Context, userID int) error {
	return m.stopInstallation(ctx, userID, false)
}

func (m *Manager) MigrateInstallation(ctx context.Context, from, to string) error {
	from, to = filepath.Clean(from), filepath.Clean(to)
	fromName, toName := filepath.Base(from), filepath.Base(to)
	if !filepath.IsAbs(from) || !filepath.IsAbs(to) || filepath.Dir(from) != filepath.Dir(to) ||
		!((fromName == "spynel" && toName == "iris") || (fromName == "iris" && toName == "spynel")) {
		return errors.New("startup migration requires exact sibling iris and spynel launchers")
	}
	var changed []changedRegistration
	scopes := []bool{false}
	if m.SystemWide {
		scopes = append(scopes, true)
	}
	for _, system := range scopes {
		directory := filepath.Join(m.Home, ".config", "systemd", "user")
		pattern := "spynel-????????.service"
		if m.GOOS == "darwin" {
			directory = filepath.Join(m.Home, "Library", "LaunchAgents")
			pattern = "dev.spynel.workspace.????????.plist"
		}
		if system {
			directory = m.SystemUnitDirectory
			if m.GOOS == "darwin" {
				directory = m.SystemLaunchDirectory
			}
		}
		paths, err := filepath.Glob(filepath.Join(directory, pattern))
		if err != nil {
			return m.rollbackMigration(ctx, changed, err)
		}
		if len(paths) > 4096 {
			return m.rollbackMigration(ctx, changed, errors.New("too many startup registrations to inspect safely"))
		}
		for _, path := range paths {
			data, err := readRegistration(path)
			if err != nil {
				return m.rollbackMigration(ctx, changed, err)
			}
			updated, matches, err := migratedRegistration(data, m.GOOS, from, to)
			if err != nil {
				return m.rollbackMigration(ctx, changed, err)
			}
			if !matches {
				continue
			}
			registration := changedRegistration{path: path, data: data, system: system}
			if m.GOOS == "linux" {
				registration.loaded, err = m.systemdRegistrationActive(ctx, path, system)
				if err != nil {
					return m.rollbackMigration(ctx, changed, err)
				}
			}
			if err := m.writeMigratedRegistration(ctx, path, updated, system); err != nil {
				return m.rollbackMigration(ctx, changed, err)
			}
			changed = append(changed, registration)
			if m.GOOS == "darwin" {
				loaded, err := m.reloadLaunchdRegistration(ctx, path, system, false)
				changed[len(changed)-1].loaded = loaded
				if err != nil {
					return m.rollbackMigration(ctx, changed, err)
				}
			}
		}
	}
	if m.GOOS == "linux" {
		for _, system := range scopes {
			var active []string
			found := false
			for _, registration := range changed {
				if registration.system != system {
					continue
				}
				found = true
				if registration.loaded {
					active = append(active, filepath.Base(registration.path))
				}
			}
			if !found {
				continue
			}
			if err := m.reloadMigratedRegistrations(ctx, system, active); err != nil {
				return m.rollbackMigration(ctx, changed, err)
			}
		}
	}
	return nil
}

type changedRegistration struct {
	path   string
	data   []byte
	system bool
	loaded bool
}

func (m *Manager) rollbackMigration(ctx context.Context, changed []changedRegistration, cause error) error {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	errorsFound := []error{cause}
	reload := make(map[bool][]string)
	for index := len(changed) - 1; index >= 0; index-- {
		registration := changed[index]
		if err := fsx.AtomicWriteFile(registration.path, registration.data, 0o600); err != nil {
			errorsFound = append(errorsFound, err)
			continue
		}
		if m.GOOS == "darwin" && registration.loaded {
			if _, err := m.reloadLaunchdRegistration(cleanupCtx, registration.path, registration.system, true); err != nil {
				errorsFound = append(errorsFound, err)
			}
		}
		if m.GOOS == "linux" {
			if _, exists := reload[registration.system]; !exists {
				reload[registration.system] = nil
			}
			if registration.loaded {
				reload[registration.system] = append(reload[registration.system], filepath.Base(registration.path))
			}
		}
	}
	if m.GOOS == "linux" {
		for _, system := range []bool{false, true} {
			active, exists := reload[system]
			if !exists {
				continue
			}
			if err := m.reloadMigratedRegistrations(cleanupCtx, system, active); err != nil {
				errorsFound = append(errorsFound, err)
			}
		}
	}
	return errors.Join(errorsFound...)
}

func (m *Manager) writeMigratedRegistration(ctx context.Context, path string, data []byte, system bool) error {
	if m.GOOS == "darwin" {
		return m.writeValidated(ctx, path, data, "plutil", "-lint")
	}
	arguments := []string{"verify", "--man=no"}
	if !system {
		arguments = append(arguments, "--user")
	}
	return m.writeValidated(ctx, path, data, "systemd-analyze", arguments...)
}

func (m *Manager) reloadLaunchdRegistration(ctx context.Context, path string, system, wantLoaded bool) (bool, error) {
	domain := "gui/" + strconv.Itoa(os.Getuid())
	if system {
		domain = "system"
	}
	name := strings.TrimSuffix(filepath.Base(path), ".plist")
	service := domain + "/" + name
	if _, err := m.run(ctx, "launchctl", "bootout", service); err == nil {
		wantLoaded = true
	} else {
		output, queryErr := m.run(ctx, "launchctl", "print", service)
		if queryErr == nil || ctx.Err() != nil || !launchdServiceMissing(output, queryErr) {
			return true, fmt.Errorf("stop startup registration %s: %w", name, err)
		}
	}
	if !wantLoaded {
		return false, nil
	}
	if _, err := m.run(ctx, "launchctl", "bootstrap", domain, path); err != nil {
		return true, fmt.Errorf("reload startup registration %s: %w", name, err)
	}
	return true, nil
}

func launchdServiceMissing(output string, err error) bool {
	detail := output + fmt.Sprint(err)
	return strings.Contains(detail, "Could not find service") || strings.Contains(detail, "Could not find domain")
}

func (m *Manager) systemdCommand(system bool) ([]string, bool, error) {
	runtimePath := filepath.Join("/run/user", strconv.Itoa(os.Getuid()))
	managerPath := filepath.Join(runtimePath, "systemd", "private")
	arguments := []string{"--no-ask-password", "--user"}
	if value := os.Getenv("XDG_RUNTIME_DIR"); !system && value != "" {
		managerPath = filepath.Join(value, "systemd", "private")
	}
	if system {
		managerPath = "/run/systemd/system"
		arguments = []string{"--no-ask-password"}
	}
	if _, err := os.Stat(managerPath); errors.Is(err, os.ErrNotExist) {
		return arguments, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return arguments, true, nil
}

func (m *Manager) systemdRegistrationActive(ctx context.Context, path string, system bool) (bool, error) {
	arguments, available, err := m.systemdCommand(system)
	if err != nil || !available {
		return false, err
	}
	name := filepath.Base(path)
	output, err := m.run(ctx, "systemctl", append(arguments, "show", "--property=ActiveState", "--value", name)...)
	if err != nil {
		return false, fmt.Errorf("inspect startup registration %s: %w", name, err)
	}
	switch strings.TrimSpace(output) {
	case "active", "activating", "reloading":
		return true, nil
	case "inactive", "failed", "deactivating":
		return false, nil
	default:
		return false, fmt.Errorf("inspect startup registration %s: unrecognized systemd state %q", name, strings.TrimSpace(output))
	}
}

func (m *Manager) reloadMigratedRegistrations(ctx context.Context, system bool, active []string) error {
	arguments, available, err := m.systemdCommand(system)
	if err != nil || !available {
		return err
	}
	if _, err := m.run(ctx, "systemctl", append(arguments, "daemon-reload")...); err != nil {
		return err
	}
	for _, name := range active {
		if _, err := m.run(ctx, "systemctl", append(arguments, "restart", name)...); err != nil {
			return fmt.Errorf("restart startup registration %s: %w", name, err)
		}
	}
	return nil
}

func migratedRegistration(data []byte, goos, from, to string) ([]byte, bool, error) {
	if goos == "linux" {
		old := []byte("ExecStart=:" + systemdQuote(from) + " ")
		if bytes.Count(data, old) == 0 {
			return data, false, nil
		}
		if bytes.Count(data, old) != 1 {
			return nil, false, errors.New("startup registration contains duplicate executable entries")
		}
		return bytes.Replace(data, old, []byte("ExecStart=:"+systemdQuote(to)+" "), 1), true, nil
	}
	if goos != "darwin" {
		return nil, false, fmt.Errorf("startup migration is not supported on %s", goos)
	}
	arguments, err := launchdProgramArguments(data)
	if err != nil || len(arguments) == 0 || arguments[0] != from {
		return data, false, err
	}
	key := []byte("<key>ProgramArguments</key>")
	start := bytes.Index(data, key)
	if start < 0 {
		return nil, false, errors.New("invalid launchd program arguments")
	}
	old, err := xml.Marshal(from)
	if err != nil {
		return nil, false, err
	}
	newValue, err := xml.Marshal(to)
	if err != nil {
		return nil, false, err
	}
	end := bytes.Index(data[start+len(key):], []byte("</array>"))
	index := bytes.Index(data[start+len(key):], old)
	if end < 0 || index < 0 || index >= end {
		return nil, false, errors.New("invalid launchd program arguments")
	}
	index += start + len(key)
	updated := append([]byte(nil), data[:index]...)
	updated = append(updated, newValue...)
	updated = append(updated, data[index+len(old):]...)
	return updated, true, nil
}

func (m *Manager) stopInstallation(ctx context.Context, userID int, remove bool) error {
	scopes := []bool{false}
	if m.SystemWide {
		scopes = append(scopes, true)
	}
	for _, system := range scopes {
		directory := filepath.Join(m.Home, ".config", "systemd", "user")
		pattern := "spynel-????????.service"
		if m.GOOS == "darwin" {
			directory = filepath.Join(m.Home, "Library", "LaunchAgents")
			pattern = "dev.spynel.workspace.????????.plist"
		}
		if system {
			directory = m.SystemUnitDirectory
			if m.GOOS == "darwin" {
				directory = m.SystemLaunchDirectory
			}
		}
		paths, err := filepath.Glob(filepath.Join(directory, pattern))
		if err != nil {
			return err
		}
		if len(paths) > 4096 {
			return errors.New("too many startup registrations to inspect safely")
		}
		for _, path := range paths {
			data, err := readRegistration(path)
			if err != nil {
				return err
			}
			matches, err := m.registrationMatches(data)
			if err != nil {
				return fmt.Errorf("inspect startup registration %s: %w", path, err)
			}
			if !matches {
				continue
			}
			name := filepath.Base(path)
			if m.GOOS == "darwin" {
				domain := "gui/" + strconv.Itoa(userID)
				if system {
					domain = "system"
				}
				service := domain + "/" + strings.TrimSuffix(name, ".plist")
				if _, err := runCommand(ctx, m.Log, "launchctl", "bootout", service); err != nil {
					// An installed plist need not have been loaded this login.
					if output, queryErr := runCommand(ctx, m.Log, "launchctl", "print", service); queryErr == nil || ctx.Err() != nil || !launchdServiceMissing(output, queryErr) {
						return fmt.Errorf("stop startup registration %s: %w", name, err)
					}
				}
				if remove {
					if err := os.Remove(path); err != nil {
						return err
					}
				}
				continue
			}
			args := []string{"--no-ask-password"}
			if !system {
				args = append(args, "--user")
			}
			runtimePath := filepath.Join("/run/user", strconv.Itoa(userID))
			if userID == os.Getuid() && os.Getenv("XDG_RUNTIME_DIR") != "" {
				runtimePath = os.Getenv("XDG_RUNTIME_DIR")
			}
			managerPath := filepath.Join(runtimePath, "systemd", "private")
			if system {
				managerPath = "/run/systemd/system"
			}
			_, managerErr := os.Stat(managerPath)
			if managerErr != nil && !errors.Is(managerErr, os.ErrNotExist) {
				return managerErr
			}
			run := func(arguments ...string) (string, error) {
				if managerErr != nil {
					return "", nil // No manager is running; remove the future registration.
				}
				if !system && os.Geteuid() == 0 && userID != 0 {
					prefix := []string{"-u", "#" + strconv.Itoa(userID), "--", "env", "XDG_RUNTIME_DIR=" + runtimePath, "systemctl"}
					return runCommand(ctx, m.Log, "sudo", append(prefix, arguments...)...)
				}
				return runCommand(ctx, m.Log, "systemctl", arguments...)
			}
			if _, err := run(append(args, "stop", name)...); err != nil {
				state, queryErr := run(append(args, "show", "--property=ActiveState", "--value", name)...)
				if queryErr != nil || strings.TrimSpace(state) != "inactive" && strings.TrimSpace(state) != "failed" {
					return fmt.Errorf("stop startup registration %s: %w", name, err)
				}
			}
			if !remove {
				continue
			}
			for _, target := range []string{"default.target.wants", "multi-user.target.wants"} {
				link := filepath.Join(directory, target, name)
				if destination, err := os.Readlink(link); err == nil && (destination == filepath.Join("..", name) || destination == path) {
					if err := os.Remove(link); err != nil {
						return err
					}
				}
			}
			if err := os.Remove(path); err != nil {
				return err
			}
			if _, err := run(append(args, "daemon-reload")...); err != nil {
				return err
			}
		}
	}
	return nil
}

func readRegistration(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxCommandOutput {
		return nil, fmt.Errorf("invalid startup registration %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxCommandOutput+1))
	if len(data) > maxCommandOutput {
		return nil, errors.New("startup registration exceeds size limit")
	}
	return data, err
}

func (m *Manager) registrationMatches(data []byte) (bool, error) {
	if m.GOOS == "linux" {
		for _, line := range strings.Split(string(data), "\n") {
			if !strings.HasPrefix(line, "ExecStart=") {
				continue
			}
			line = strings.TrimPrefix(line, "ExecStart=")
			line = strings.TrimPrefix(line, ":")
			if strings.HasPrefix(line, systemdQuote(m.Executable)+" ") {
				return true, nil
			}
			if m.NPMLauncher != "" && strings.Contains(line, " "+systemdQuote(m.NPMLauncher)+" "+systemdQuote("serve")+" ") {
				return true, nil
			}
		}
		return false, nil
	}
	args, err := launchdProgramArguments(data)
	return err == nil && (len(args) > 0 && args[0] == m.Executable || m.NPMLauncher != "" && len(args) > 1 && args[1] == m.NPMLauncher), err
}

func launchdProgramArguments(data []byte) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "key" {
			continue
		}
		var key string
		if err := decoder.DecodeElement(&key, &start); err != nil {
			return nil, err
		}
		if key != "ProgramArguments" {
			continue
		}
		var arguments struct {
			Values []string `xml:"string"`
		}
		if err := decoder.Decode(&arguments); err != nil {
			return nil, err
		}
		return arguments.Values, nil
	}
}
