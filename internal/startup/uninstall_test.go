package startup

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestRemoveInstallation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix startup cleanup")
	}
	for _, platform := range []string{"linux", "darwin"} {
		for _, npm := range []bool{false, true} {
			t.Run(platform+"/npm="+strconv.FormatBool(npm), func(t *testing.T) {
				home := t.TempDir()
				tools := t.TempDir()
				log := filepath.Join(tools, "calls")
				t.Setenv("PATH", tools)
				t.Setenv("STARTUP_TEST_LOG", log)
				t.Setenv("XDG_RUNTIME_DIR", filepath.Join(home, "run"))
				if err := os.MkdirAll(filepath.Join(home, "run", "systemd", "private"), 0700); err != nil {
					t.Fatal(err)
				}
				for _, name := range []string{"systemctl", "launchctl"} {
					if err := os.WriteFile(filepath.Join(tools, name), []byte("#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$STARTUP_TEST_LOG\"\ncase \"$*\" in *stop*|*bootout*) [ -z \"${STARTUP_TEST_FAIL:-}\" ] || exit 1;; esac\ncase \"$*\" in *show*) echo active;; esac\n"), 0700); err != nil {
						t.Fatal(err)
					}
				}
				manager := &Manager{GOOS: platform, Home: home, Executable: filepath.Join(home, `install Ω $ ' " % \`, "spynel"), SystemWide: platform == "darwin", SystemLaunchDirectory: filepath.Join(home, "system")}
				if npm {
					manager.NPMLauncher = filepath.Join(home, `npm Ω $ ' " % \`, "spynel.js")
					manager.NodeExecutable = "/node"
				}
				manager.RunCommand = func(_ context.Context, name string, args ...string) (string, error) {
					if name == "launchctl" && args[0] == "print-disabled" {
						return "disabled services = {\n}\n", nil
					}
					if name == "systemctl" && strings.Contains(strings.Join(args, " "), "list-unit-files") {
						return args[len(args)-1] + " enabled enabled\n", nil
					}
					return "", nil
				}
				for _, path := range []string{manager.Executable, manager.NPMLauncher} {
					if path == "" {
						continue
					}
					if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0700); err != nil {
						t.Fatal(err)
					}
				}
				manager.NodeExecutable = manager.Executable
				cfg := startupTestConfig(t, filepath.Join(home, "workspace Ω"))
				if err := manager.Sync(cfg, true); err != nil {
					t.Fatal(err)
				}
				directory := filepath.Join(home, ".config", "systemd", "user")
				name := "spynel-" + workspaceID(cfg) + ".service"
				if platform == "darwin" {
					directory = manager.SystemLaunchDirectory
					name = "dev.spynel.workspace." + workspaceID(cfg) + ".plist"
				}
				path := filepath.Join(directory, name)
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				unrelated := filepath.Join(directory, strings.Replace(name, workspaceID(cfg), "12345678", 1))
				other := strings.ReplaceAll(string(data), "spynel", "different-program")
				if err := os.WriteFile(unrelated, []byte(other), 0600); err != nil {
					t.Fatal(err)
				}
				t.Setenv("STARTUP_TEST_FAIL", "1")
				if err := manager.RemoveInstallation(t.Context(), os.Getuid()); err == nil {
					t.Fatal("failed stop removed registration")
				}
				if _, err := os.Stat(path); err != nil {
					t.Fatal("registration lost on stop failure")
				}
				t.Setenv("STARTUP_TEST_FAIL", "")
				if err := manager.StopInstallation(t.Context(), os.Getuid()); err != nil {
					t.Fatal(err)
				}
				if retained, err := os.ReadFile(path); err != nil || string(retained) != string(data) {
					t.Fatalf("stop changed future startup registration: %v", err)
				}
				if platform == "linux" {
					if _, err := os.Readlink(filepath.Join(directory, "default.target.wants", name)); err != nil {
						t.Fatal("stop removed future startup enablement:", err)
					}
				}
				if err := manager.RemoveInstallation(t.Context(), os.Getuid()); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("registration remains")
				}
				if _, err := os.Stat(unrelated); err != nil {
					t.Fatal("unrelated registration removed")
				}
				calls, err := os.ReadFile(log)
				if err != nil {
					t.Fatal(err)
				}
				want := "--user stop " + name
				if platform == "darwin" {
					want = "bootout system/" + strings.TrimSuffix(name, ".plist")
				}
				if !strings.Contains(string(calls), want) {
					t.Fatalf("missing stop: %s", calls)
				}
				if platform == "linux" {
					if _, err := os.Lstat(filepath.Join(directory, "default.target.wants", name)); !errors.Is(err, os.ErrNotExist) {
						t.Fatal("startup enablement remains")
					}
					if err := os.RemoveAll(filepath.Join(home, "run")); err != nil {
						t.Fatal(err)
					}
					if err := manager.Sync(cfg, true); err != nil {
						t.Fatal(err)
					}
					t.Setenv("STARTUP_TEST_FAIL", "1")
					if err := manager.RemoveInstallation(t.Context(), os.Getuid()); err != nil {
						t.Fatal("offline registration cleanup:", err)
					}
				}
			})
		}
	}
}

func TestMigrateInstallation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("standalone distribution excludes Windows")
	}
	for _, platform := range []string{"linux", "darwin"} {
		t.Run(platform, func(t *testing.T) {
			home := t.TempDir()
			runtimeDirectory := filepath.Join(home, "run")
			t.Setenv("XDG_RUNTIME_DIR", runtimeDirectory)
			if err := os.MkdirAll(filepath.Join(runtimeDirectory, "systemd", "private"), 0o700); err != nil {
				t.Fatal(err)
			}
			install := filepath.Join(home, "installation")
			oldExecutable := filepath.Join(install, "spynel")
			newExecutable := filepath.Join(install, "iris")
			unrelatedExecutable := filepath.Join(home, "other", "iris")
			for _, path := range []string{oldExecutable, newExecutable, unrelatedExecutable} {
				if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o700); err != nil {
					t.Fatal(err)
				}
			}
			var calls []string
			run := func(_ context.Context, name string, args ...string) (string, error) {
				calls = append(calls, strings.Join(append([]string{name}, args...), " "))
				if name == "launchctl" && args[0] == "print-disabled" {
					return "disabled services = {\n}\n", nil
				}
				if name == "systemctl" && strings.Contains(strings.Join(args, " "), "list-unit-files") {
					return args[len(args)-1] + " enabled enabled\n", nil
				}
				return "", nil
			}
			manager := &Manager{GOOS: platform, Home: home, Executable: oldExecutable, SystemLaunchDirectory: filepath.Join(home, "system"), RunCommand: run}
			cfg := startupTestConfig(t, filepath.Join(home, "workspace"))
			if err := manager.Sync(cfg, true); err != nil {
				t.Fatal(err)
			}
			unrelated := *manager
			unrelated.Executable = unrelatedExecutable
			otherConfig := startupTestConfig(t, filepath.Join(home, "other-workspace"))
			if err := unrelated.Sync(otherConfig, true); err != nil {
				t.Fatal(err)
			}
			directory := filepath.Join(home, ".config", "systemd", "user")
			name := "spynel-" + workspaceID(cfg) + ".service"
			otherName := "spynel-" + workspaceID(otherConfig) + ".service"
			if platform == "darwin" {
				directory = filepath.Join(home, "Library", "LaunchAgents")
				name = "dev.spynel.workspace." + workspaceID(cfg) + ".plist"
				otherName = "dev.spynel.workspace." + workspaceID(otherConfig) + ".plist"
			}
			otherPath := filepath.Join(directory, otherName)
			otherBefore, err := os.ReadFile(otherPath)
			if err != nil {
				t.Fatal(err)
			}
			calls = nil
			if err := manager.MigrateInstallation(t.Context(), oldExecutable, newExecutable); err != nil {
				t.Fatal(err)
			}
			if platform == "darwin" {
				domain := "gui/" + strconv.Itoa(os.Getuid())
				commands := strings.Join(calls, "\n")
				if !strings.Contains(commands, "launchctl bootout "+domain+"/"+strings.TrimSuffix(name, ".plist")) ||
					!strings.Contains(commands, "launchctl bootstrap "+domain+" "+filepath.Join(directory, name)) {
					t.Fatalf("launchd migration did not reload the owned job: %s", commands)
				}
				if strings.Contains(commands, strings.TrimSuffix(otherName, ".plist")) {
					t.Fatalf("launchd migration touched unrelated job: %s", commands)
				}
			}
			data, err := os.ReadFile(filepath.Join(directory, name))
			if err != nil {
				t.Fatal(err)
			}
			newManager := *manager
			newManager.Executable = newExecutable
			if oldMatch, err := manager.registrationMatches(data); err != nil || oldMatch {
				t.Fatalf("legacy registration still matches: %t %v", oldMatch, err)
			}
			if newMatch, err := newManager.registrationMatches(data); err != nil || !newMatch {
				t.Fatalf("Iris registration does not match: %t %v", newMatch, err)
			}
			otherAfter, err := os.ReadFile(otherPath)
			if err != nil || !bytes.Equal(otherAfter, otherBefore) {
				t.Fatalf("unrelated registration changed: %v", err)
			}
			if err := manager.MigrateInstallation(t.Context(), newExecutable, oldExecutable); err != nil {
				t.Fatal(err)
			}
			data, err = os.ReadFile(filepath.Join(directory, name))
			if err != nil {
				t.Fatal(err)
			}
			if oldMatch, err := manager.registrationMatches(data); err != nil || !oldMatch {
				t.Fatalf("legacy registration rollback does not match: %t %v", oldMatch, err)
			}
		})
	}
}

func TestMigrateInstallationRollbackReloadsAfterCancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("standalone distribution excludes Windows")
	}
	home := t.TempDir()
	oldExecutable := filepath.Join(home, "installation", "spynel")
	newExecutable := filepath.Join(home, "installation", "iris")
	for _, path := range []string{oldExecutable, newExecutable} {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	manager := &Manager{
		GOOS:       "darwin",
		Home:       home,
		Executable: oldExecutable,
		RunCommand: func(_ context.Context, name string, args ...string) (string, error) {
			if name == "launchctl" && args[0] == "print-disabled" {
				return "disabled services = {\n}\n", nil
			}
			return "", nil
		},
	}
	cfg := startupTestConfig(t, filepath.Join(home, "workspace"))
	if err := manager.Sync(cfg, true); err != nil {
		t.Fatal(err)
	}
	registration := filepath.Join(home, "Library", "LaunchAgents", "dev.spynel.workspace."+workspaceID(cfg)+".plist")
	before, err := os.ReadFile(registration)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	bootouts := 0
	rollbackBootstrapped := false
	manager.RunCommand = func(commandCtx context.Context, name string, args ...string) (string, error) {
		if name == "launchctl" && args[0] == "bootout" {
			bootouts++
			if bootouts == 1 {
				cancel()
			}
			return "", nil
		}
		if name == "launchctl" && args[0] == "bootstrap" {
			if err := commandCtx.Err(); err != nil {
				return "", err
			}
			rollbackBootstrapped = true
		}
		return "", nil
	}
	err = manager.MigrateInstallation(ctx, oldExecutable, newExecutable)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("migration error = %v, want context cancellation", err)
	}
	after, err := os.ReadFile(registration)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("rollback did not restore the original launchd registration")
	}
	if bootouts != 2 || !rollbackBootstrapped {
		t.Fatalf("rollback reload = bootouts %d, bootstrapped %t", bootouts, rollbackBootstrapped)
	}
}
