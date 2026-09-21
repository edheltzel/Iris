package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	root, err := os.MkdirTemp("", "iris-instance-test-")
	if err != nil {
		panic(err)
	}
	for name, value := range map[string]string{
		"HOME":            filepath.Join(root, "home"),
		"XDG_CONFIG_HOME": filepath.Join(root, "config"),
		"XDG_CACHE_HOME":  filepath.Join(root, "cache"),
	} {
		if err := os.Setenv(name, value); err != nil {
			panic(err)
		}
	}
	code := m.Run()
	_ = os.RemoveAll(root)
	os.Exit(code)
}

func testEnvironmentID(character string) string { return strings.Repeat(character, 64) }

func TestEnvironmentIDIsStablePrivateAndEnvironmentScoped(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	first, err := EnvironmentID()
	if err != nil {
		t.Fatal(err)
	}
	second, err := EnvironmentID()
	if err != nil || second != first || !validEnvironmentID(first) {
		t.Fatalf("stable environment ID = %q, %q, %v", first, second, err)
	}
	tokenPath := filepath.Join(home, ".agents", "Iris", "environment-token")
	token, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(first, strings.TrimSpace(string(token))) {
		t.Fatal("published environment ID exposed its private source token")
	}
	info, err := os.Stat(tokenPath)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("environment token permissions = %v, %v", info, err)
	}
	home2 := t.TempDir()
	t.Setenv("HOME", home2)
	separate, err := EnvironmentID()
	if err != nil || separate == first {
		t.Fatalf("separate configuration environment ID = %q, first %q, %v", separate, first, err)
	}
}

func TestEnvironmentIDMigratesLegacyNamespaces(t *testing.T) {
	token := strings.Repeat("ab", environmentTokenBytes)
	for _, name := range []string{"iris", "spynel"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			configDirectory, err := os.UserConfigDir()
			if err != nil {
				t.Fatal(err)
			}
			legacy := filepath.Join(configDirectory, name)
			if err := os.MkdirAll(legacy, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(legacy, "environment-token"), []byte(token+"\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			id, err := EnvironmentID()
			if err != nil || !validEnvironmentID(id) {
				t.Fatalf("migrated environment ID = %q, %v", id, err)
			}
			dest := filepath.Join(home, ".agents", "Iris", "environment-token")
			got, err := os.ReadFile(dest)
			if err != nil || strings.TrimSpace(string(got)) != token {
				t.Fatalf("migrated token = %q, %v", got, err)
			}
			if _, err := os.Stat(legacy); !os.IsNotExist(err) {
				t.Fatalf("legacy namespace still present: %v", err)
			}
		})
	}
}

func TestEnvironmentIDResolvesDualLegacyNamespaces(t *testing.T) {
	t.Run("prefers token source", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		config, err := os.UserConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		root := filepath.Join(home, ".agents", "Iris")
		if err := os.MkdirAll(root, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "existing"), []byte("kept"), 0o600); err != nil {
			t.Fatal(err)
		}
		iris := filepath.Join(config, "iris")
		spynel := filepath.Join(config, "spynel")
		for _, directory := range []string{iris, spynel} {
			if err := os.MkdirAll(directory, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(iris, "leftover"), []byte("not selected"), 0o600); err != nil {
			t.Fatal(err)
		}
		token := strings.Repeat("cd", environmentTokenBytes)
		if err := os.WriteFile(filepath.Join(spynel, "environment-token"), []byte(token+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := EnvironmentID(); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(root, "environment-token"))
		if err != nil || strings.TrimSpace(string(data)) != token {
			t.Fatalf("selected token = %q, %v", data, err)
		}
		if _, err := os.Stat(spynel); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("selected source remains: %v", err)
		}
		if _, err := os.Stat(filepath.Join(iris, "leftover")); err != nil {
			t.Fatalf("unselected source changed: %v", err)
		}
	})

	t.Run("rejects two tokens", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		config, err := os.UserConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		for index, name := range []string{"iris", "spynel"} {
			directory := filepath.Join(config, name)
			if err := os.MkdirAll(directory, 0o700); err != nil {
				t.Fatal(err)
			}
			token := strings.Repeat([]string{"ab", "cd"}[index], environmentTokenBytes)
			if err := os.WriteFile(filepath.Join(directory, "environment-token"), []byte(token+"\n"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := EnvironmentID(); err == nil || !strings.Contains(err.Error(), "both legacy environment identity sources contain tokens") {
			t.Fatalf("dual-token error = %v", err)
		}
		if _, err := os.Stat(filepath.Join(home, ".agents", "Iris", "environment-token")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("destination token unexpectedly published: %v", err)
		}
	})

	t.Run("resumes partial migration", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		config, err := os.UserConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		root := filepath.Join(home, ".agents", "Iris")
		legacy := filepath.Join(config, "spynel")
		for _, directory := range []string{filepath.Join(root, "processes"), filepath.Join(legacy, "processes")} {
			if err := os.MkdirAll(directory, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		token := strings.Repeat("ef", environmentTokenBytes)
		if err := os.WriteFile(filepath.Join(legacy, "environment-token"), []byte(token+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "processes", "123.json"), []byte("current"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(legacy, "processes", "123.json"), []byte("legacy"), 0o600); err != nil {
			t.Fatal(err)
		}
		for attempt := 1; attempt <= 2; attempt++ {
			if _, err := EnvironmentID(); err == nil || !strings.Contains(err.Error(), "migration entry conflicts with destination") {
				t.Fatalf("attempt %d conflict = %v", attempt, err)
			}
		}
		data, err := os.ReadFile(filepath.Join(root, "environment-token"))
		if err != nil || strings.TrimSpace(string(data)) != token {
			t.Fatalf("partially migrated token = %q, %v", data, err)
		}
		if _, err := os.Stat(legacy); err != nil {
			t.Fatalf("conflicting source removed: %v", err)
		}
	})

	t.Run("rejects destination token conflict", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		config, err := os.UserConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		root := filepath.Join(home, ".agents", "Iris")
		legacy := filepath.Join(config, "spynel")
		for _, directory := range []string{root, legacy} {
			if err := os.MkdirAll(directory, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(root, "environment-token"), []byte(strings.Repeat("ab", environmentTokenBytes)+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(legacy, "environment-token"), []byte(strings.Repeat("cd", environmentTokenBytes)+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := EnvironmentID(); err == nil || !strings.Contains(err.Error(), "legacy environment identity conflicts with destination") {
			t.Fatalf("destination-token conflict = %v", err)
		}
		if _, err := os.Stat(legacy); err != nil {
			t.Fatalf("conflicting source removed: %v", err)
		}
	})

	t.Run("merges tokenless leftovers", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		config, err := os.UserConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		for name, file := range map[string]string{"spynel": "from-spynel", "iris": "from-iris"} {
			directory := filepath.Join(config, name)
			if err := os.MkdirAll(directory, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(directory, file), []byte(name), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := EnvironmentID(); err != nil {
			t.Fatal(err)
		}
		root := filepath.Join(home, ".agents", "Iris")
		for _, file := range []string{"from-spynel", "from-iris"} {
			if _, err := os.Stat(filepath.Join(root, file)); err != nil {
				t.Fatalf("missing merged file %s: %v", file, err)
			}
		}
		for _, name := range []string{"spynel", "iris"} {
			if _, err := os.Stat(filepath.Join(config, name)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("legacy source %s remains: %v", name, err)
			}
		}
	})
}

func TestLeasePublishesEnvironmentAndLegacyOwnerRemainsFenced(t *testing.T) {
	state := t.TempDir()
	owner, err := NewWithEnvironmentID(state, testEnvironmentID("a"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	owner.now = func() time.Time { return now }
	token, _ := owner.NewToken()
	lease, acquired, err := owner.TryAcquire("127.0.0.1:10001", token)
	if err != nil || !acquired || lease.EnvironmentID != testEnvironmentID("a") {
		t.Fatalf("environment lease = %#v, %t, %v", lease, acquired, err)
	}

	lease.EnvironmentID = "obsolete-or-invalid"
	data, _ := json.Marshal(lease)
	if err := os.WriteFile(owner.leasePath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	contender, err := NewWithEnvironmentID(state, testEnvironmentID("b"))
	if err != nil {
		t.Fatal(err)
	}
	contender.now = func() time.Time { return now }
	legacy, err := contender.Current()
	if err != nil || legacy.EnvironmentID != "" || contender.CanTakeOver(legacy) {
		t.Fatalf("legacy lease safety = %#v, %v", legacy, err)
	}
	contenderToken, _ := contender.NewToken()
	if _, acquired, err := contender.TryAcquire("127.0.0.1:10002", contenderToken); err != nil || acquired {
		t.Fatalf("fresh legacy owner takeover = %t, %v", acquired, err)
	}
	contender.now = func() time.Time { return now.Add(StaleAfter) }
	replacement, acquired, err := contender.TryAcquire("127.0.0.1:10002", contenderToken)
	if err != nil || !acquired || replacement.EnvironmentID != testEnvironmentID("b") {
		t.Fatalf("stale legacy replacement = %#v, %t, %v", replacement, acquired, err)
	}
}

func TestRenewAndHandoffRestoreValidatedEnvironmentID(t *testing.T) {
	state := t.TempDir()
	owner, err := NewWithEnvironmentID(state, testEnvironmentID("a"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	owner.now = func() time.Time { return now }
	token, _ := owner.NewToken()
	lease, acquired, err := owner.TryAcquire("127.0.0.1:10001", token)
	if err != nil || !acquired {
		t.Fatalf("owner acquire = %#v, %t, %v", lease, acquired, err)
	}

	lease.EnvironmentID = "damaged"
	data, _ := json.Marshal(lease)
	if err := os.WriteFile(owner.leasePath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	now = now.Add(HeartbeatInterval)
	renewed, owned, err := owner.Renew(token)
	if err != nil || !owned || renewed.EnvironmentID != testEnvironmentID("a") {
		t.Fatalf("renewed environment = %#v, %t, %v", renewed, owned, err)
	}

	target, err := NewWithEnvironmentID(state, testEnvironmentID("a"))
	if err != nil {
		t.Fatal(err)
	}
	renewed.EnvironmentID = ""
	data, _ = json.Marshal(renewed)
	if err := os.WriteFile(owner.leasePath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	now = now.Add(HeartbeatInterval)
	handoff, handedOff, err := owner.Handoff(token, target.ID())
	if err != nil || !handedOff || handoff.EnvironmentID != testEnvironmentID("a") {
		t.Fatalf("handoff environment = %#v, %t, %v", handoff, handedOff, err)
	}
}

func TestOnlyOneConcurrentContenderAcquiresLease(t *testing.T) {
	state := t.TempDir()
	const contenders = 24
	elections := make([]*Election, contenders)
	for index := range elections {
		var err error
		elections[index], err = New(state)
		if err != nil {
			t.Fatal(err)
		}
	}
	start := make(chan struct{})
	var wait sync.WaitGroup
	results := make(chan Lease, contenders)
	for index, election := range elections {
		wait.Add(1)
		go func(index int, election *Election) {
			defer wait.Done()
			<-start
			token, err := election.NewToken()
			if err != nil {
				t.Errorf("token: %v", err)
				return
			}
			lease, acquired, err := election.TryAcquire("127.0.0.1:"+time.Now().Add(time.Duration(index)).Format("150405.000000000"), token)
			if err != nil {
				t.Errorf("acquire: %v", err)
				return
			}
			if acquired {
				results <- lease
			}
		}(index, election)
	}
	close(start)
	wait.Wait()
	close(results)
	var winners []Lease
	for lease := range results {
		winners = append(winners, lease)
	}
	if len(winners) != 1 {
		t.Fatalf("winners = %d, want exactly one: %#v", len(winners), winners)
	}
	current, err := elections[0].Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.InstanceID != winners[0].InstanceID || current.Token != winners[0].Token {
		t.Fatalf("published lease = %#v, winner = %#v", current, winners[0])
	}
}

func TestOnlyOneConcurrentContenderWinsStaleTakeover(t *testing.T) {
	state := t.TempDir()
	startTime := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	owner.now = func() time.Time { return startTime }
	ownerToken, _ := owner.NewToken()
	if _, acquired, err := owner.TryAcquire("127.0.0.1:9000", ownerToken); err != nil || !acquired {
		t.Fatalf("initial acquire = %t, %v", acquired, err)
	}

	const contenders = 24
	start := make(chan struct{})
	winners := make(chan Lease, contenders)
	var wait sync.WaitGroup
	for index := 0; index < contenders; index++ {
		candidate, err := New(state)
		if err != nil {
			t.Fatal(err)
		}
		candidate.now = func() time.Time { return startTime.Add(StaleAfter) }
		wait.Add(1)
		go func(index int, candidate *Election) {
			defer wait.Done()
			<-start
			token, _ := candidate.NewToken()
			lease, acquired, err := candidate.TryAcquire(fmt.Sprintf("127.0.0.1:%d", 10000+index), token)
			if err != nil {
				t.Errorf("takeover: %v", err)
				return
			}
			if acquired {
				winners <- lease
			}
		}(index, candidate)
	}
	close(start)
	wait.Wait()
	close(winners)
	count := 0
	for range winners {
		count++
	}
	if count != 1 {
		t.Fatalf("stale takeover winners = %d, want 1", count)
	}
	if _, owned, err := owner.Renew(ownerToken); err != nil || owned {
		t.Fatalf("stale owner renewed after takeover = %t, %v", owned, err)
	}
}

func TestStaleTakeoverFencesFormerOwner(t *testing.T) {
	state := t.TempDir()
	first, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	first.now = func() time.Time { return now }
	second.now = func() time.Time { return now }
	firstToken, _ := first.NewToken()
	firstLease, acquired, err := first.TryAcquire("127.0.0.1:10001", firstToken)
	if err != nil || !acquired {
		t.Fatalf("first acquire = %#v, %t, %v", firstLease, acquired, err)
	}
	if firstLease.PID != os.Getpid() || firstLease.InstanceID == "" || firstLease.Token == "" {
		t.Fatalf("owner identity was not persisted: %#v", firstLease)
	}
	info, err := os.Stat(filepath.Join(state, "runtime", "primary.json"))
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("lease permissions = %v, %v", info, err)
	}
	secondToken, _ := second.NewToken()
	if _, acquired, err := second.TryAcquire("127.0.0.1:10002", secondToken); err != nil || acquired {
		t.Fatalf("takeover before stale acquired = %t, err = %v", acquired, err)
	}
	now = now.Add(StaleAfter)
	if _, owned, err := first.Renew(firstToken); err != nil || owned {
		t.Fatalf("owner renewed an already stale term = %t, err = %v", owned, err)
	}
	secondLease, acquired, err := second.TryAcquire("127.0.0.1:10002", secondToken)
	if err != nil || !acquired {
		t.Fatalf("stale takeover = %#v, %t, %v", secondLease, acquired, err)
	}
	if _, owned, err := first.Renew(firstToken); err != nil || owned {
		t.Fatalf("former owner renewed = %t, err = %v", owned, err)
	}
	if err := first.Release(firstToken); err != nil {
		t.Fatal(err)
	}
	current, err := second.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.InstanceID != second.ID() || current.Token != secondToken {
		t.Fatalf("former owner disturbed successor: %#v", current)
	}
}

func TestRunWhileOwnerFencesFormerTermAfterStaleTakeover(t *testing.T) {
	state := t.TempDir()
	first, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 16, 9, 0, 0, 0, time.UTC)
	first.now = func() time.Time { return now }
	second.now = func() time.Time { return now }
	firstToken, _ := first.NewToken()
	if _, acquired, err := first.TryAcquire("127.0.0.1:10001", firstToken); err != nil || !acquired {
		t.Fatalf("first acquire = %t, %v", acquired, err)
	}

	now = now.Add(StaleAfter)
	secondToken, _ := second.NewToken()
	if _, acquired, err := second.TryAcquire("127.0.0.1:10002", secondToken); err != nil || !acquired {
		t.Fatalf("stale takeover = %t, %v", acquired, err)
	}
	formerRan := false
	if ran, err := first.RunWhileOwner(firstToken, func() error { formerRan = true; return nil }); err != nil || ran || formerRan {
		t.Fatalf("former term admission = ran %t action %t err %v", ran, formerRan, err)
	}
	currentRan := false
	if ran, err := second.RunWhileOwner(secondToken, func() error { currentRan = true; return nil }); err != nil || !ran || !currentRan {
		t.Fatalf("current term admission = ran %t action %t err %v", ran, currentRan, err)
	}
}

func TestRenewalKeepsLeaseFreshAndReleaseIsImmediate(t *testing.T) {
	election, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	election.now = func() time.Time { return now }
	token, _ := election.NewToken()
	if _, acquired, err := election.TryAcquire("127.0.0.1:10001", token); err != nil || !acquired {
		t.Fatalf("acquire = %t, %v", acquired, err)
	}
	now = now.Add(HeartbeatInterval)
	lease, owned, err := election.Renew(token)
	if err != nil || !owned || !lease.HeartbeatAt.Equal(now) || !validEnvironmentID(lease.EnvironmentID) {
		t.Fatalf("renew = %#v, %t, %v", lease, owned, err)
	}
	if err := election.Release(token); err != nil {
		t.Fatal(err)
	}
	if _, err := election.Current(); !os.IsNotExist(err) {
		t.Fatalf("lease after release error = %v", err)
	}
}

func TestCurrentDoesNotWaitForElectionMutationLock(t *testing.T) {
	state := t.TempDir()
	election, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	token, err := election.NewToken()
	if err != nil {
		t.Fatal(err)
	}
	want, acquired, err := election.TryAcquire("127.0.0.1:10001", token)
	if err != nil || !acquired {
		t.Fatalf("acquire = %#v, %t, %v", want, acquired, err)
	}

	lock, err := os.OpenFile(election.lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	if err := lockFile(lock); err != nil {
		t.Fatal(err)
	}
	locked := true
	defer func() {
		if locked {
			unlockFile(lock)
		}
	}()

	type currentResult struct {
		lease Lease
		err   error
	}
	result := make(chan currentResult, 1)
	go func() {
		lease, readErr := election.Current()
		result <- currentResult{lease: lease, err: readErr}
	}()

	select {
	case got := <-result:
		if got.err != nil {
			t.Fatal(got.err)
		}
		if got.lease.InstanceID != want.InstanceID || got.lease.Token != want.Token {
			t.Fatalf("current lease = %#v, want %#v", got.lease, want)
		}
	case <-time.After(time.Second):
		unlockFile(lock)
		locked = false
		<-result
		t.Fatal("lease discovery waited for the election mutation lock")
	}
}

func TestOwnerlessOperationSerializesPrimaryPublication(t *testing.T) {
	state := t.TempDir()
	cleanup, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	token, err := owner.NewToken()
	if err != nil {
		t.Fatal(err)
	}

	inside := make(chan struct{})
	release := make(chan struct{})
	cleanupDone := make(chan error, 1)
	go func() {
		ran, runErr := cleanup.RunWhileNoPrimaryLease(func() error {
			close(inside)
			<-release
			return nil
		})
		if runErr == nil && !ran {
			runErr = errors.New("ownerless operation did not run")
		}
		cleanupDone <- runErr
	}()
	<-inside

	ownerDone := make(chan error, 1)
	go func() {
		_, acquired, acquireErr := owner.TryAcquire("127.0.0.1:10001", token)
		if acquireErr == nil && !acquired {
			acquireErr = errors.New("primary was not acquired")
		}
		ownerDone <- acquireErr
	}()
	select {
	case err := <-ownerDone:
		t.Fatalf("primary publication crossed ownerless operation: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	if err := <-cleanupDone; err != nil {
		t.Fatal(err)
	}
	if err := <-ownerDone; err != nil {
		t.Fatal(err)
	}
}

func TestOwnerlessOperationFailsClosedWhenAnyPrimaryLeaseExists(t *testing.T) {
	state := t.TempDir()
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	cleanup, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	owner.now = func() time.Time { return start }
	token, _ := owner.NewToken()
	if _, acquired, err := owner.TryAcquire("127.0.0.1:10001", token); err != nil || !acquired {
		t.Fatalf("acquire = %t, %v", acquired, err)
	}
	cleanup.now = func() time.Time { return start.Add(StaleAfter) }
	called := false
	ran, err := cleanup.RunWhileNoPrimaryLease(func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if ran || called {
		t.Fatalf("ownerless operation ran beside stale lease: ran=%t called=%t", ran, called)
	}
}

func TestOwnerlessOperationWaitsThroughCleanReleaseFailoverGap(t *testing.T) {
	state := t.TempDir()
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	secondary, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	cleanup, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	owner.now = func() time.Time { return start }
	secondary.now = func() time.Time { return start.Add(RetryInterval) }
	cleanup.now = func() time.Time { return start.Add(RetryInterval / 2) }
	token, _ := owner.NewToken()
	if _, acquired, err := owner.TryAcquire("127.0.0.1:10001", token); err != nil || !acquired {
		t.Fatalf("owner acquire = %t, %v", acquired, err)
	}
	if err := owner.Release(token); err != nil {
		t.Fatal(err)
	}

	called := false
	ran, err := cleanup.RunWhileNoPrimaryLease(func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if ran || called {
		t.Fatalf("ownerless operation crossed release failover gap: ran=%t called=%t", ran, called)
	}

	secondaryToken, _ := secondary.NewToken()
	if _, acquired, err := secondary.TryAcquire("127.0.0.1:10002", secondaryToken); err != nil || !acquired {
		t.Fatalf("secondary acquire = %t, %v", acquired, err)
	}
}

func TestOwnerlessOperationRunsAfterCleanReleaseGraceExpires(t *testing.T) {
	state := t.TempDir()
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	cleanup, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	owner.now = func() time.Time { return start }
	token, _ := owner.NewToken()
	if _, acquired, err := owner.TryAcquire("127.0.0.1:10001", token); err != nil || !acquired {
		t.Fatalf("owner acquire = %t, %v", acquired, err)
	}
	if err := owner.Release(token); err != nil {
		t.Fatal(err)
	}
	cleanup.now = func() time.Time { return start.Add(OwnerlessCleanupGrace) }
	called := false
	ran, err := cleanup.RunWhileNoPrimaryLease(func() error {
		called = true
		return nil
	})
	if err != nil || !ran || !called {
		t.Fatalf("expired release fence = ran %t, called %t, err %v", ran, called, err)
	}
}

func TestOwnerlessOperationFailsClosedForMalformedReleaseFence(t *testing.T) {
	state := t.TempDir()
	cleanup, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(cleanup.releaseFencePath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cleanup.releaseFencePath, []byte("not-json\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	called := false
	ran, err := cleanup.RunWhileNoPrimaryLease(func() error {
		called = true
		return nil
	})
	if err == nil || ran || called {
		t.Fatalf("malformed release fence = ran %t, called %t, err %v", ran, called, err)
	}
}

func TestTargetedHandoffFencesOwnerAndExcludesOtherContenders(t *testing.T) {
	state := t.TempDir()
	now := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	target, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	other, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	owner.now = func() time.Time { return now }
	target.now = func() time.Time { return now }
	other.now = func() time.Time { return now }
	ownerToken, _ := owner.NewToken()
	if _, acquired, err := owner.TryAcquire("127.0.0.1:10001", ownerToken); err != nil || !acquired {
		t.Fatalf("owner acquire = %t, %v", acquired, err)
	}
	handoff, handedOff, err := owner.Handoff(ownerToken, target.ID())
	if err != nil || !handedOff || handoff.HandoffTo != target.ID() || handoff.HandoffAt != now || !validEnvironmentID(handoff.EnvironmentID) {
		t.Fatalf("handoff = %#v, %t, %v", handoff, handedOff, err)
	}
	if _, owned, err := owner.Renew(ownerToken); err != nil || owned {
		t.Fatalf("former owner renewed handoff term = %t, %v", owned, err)
	}
	if !target.CanTakeOver(handoff) || other.CanTakeOver(handoff) {
		t.Fatalf("takeover eligibility: target = %t, other = %t", target.CanTakeOver(handoff), other.CanTakeOver(handoff))
	}
	otherToken, _ := other.NewToken()
	if _, acquired, err := other.TryAcquire("127.0.0.1:10003", otherToken); err != nil || acquired {
		t.Fatalf("non-target acquire = %t, %v", acquired, err)
	}
	targetToken, _ := target.NewToken()
	lease, acquired, err := target.TryAcquire("127.0.0.1:10002", targetToken)
	if err != nil || !acquired || lease.InstanceID != target.ID() || lease.HandoffTo != "" {
		t.Fatalf("target acquire = %#v, %t, %v", lease, acquired, err)
	}
}

func TestTargetedHandoffExpiresWhenTargetDisappears(t *testing.T) {
	state := t.TempDir()
	now := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	owner, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	contender, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	owner.now = func() time.Time { return now }
	contender.now = func() time.Time { return now }
	ownerToken, _ := owner.NewToken()
	if _, acquired, err := owner.TryAcquire("127.0.0.1:10001", ownerToken); err != nil || !acquired {
		t.Fatalf("owner acquire = %t, %v", acquired, err)
	}
	handoff, handedOff, err := owner.Handoff(ownerToken, "missing-target-instance")
	if err != nil || !handedOff {
		t.Fatalf("handoff = %#v, %t, %v", handoff, handedOff, err)
	}
	contenderToken, _ := contender.NewToken()
	if _, acquired, err := contender.TryAcquire("127.0.0.1:10002", contenderToken); err != nil || acquired {
		t.Fatalf("acquire during reservation = %t, %v", acquired, err)
	}
	now = now.Add(HandoffTimeout)
	if !contender.CanTakeOver(handoff) {
		t.Fatal("expired handoff did not become eligible for fallback takeover")
	}
	if _, acquired, err := contender.TryAcquire("127.0.0.1:10002", contenderToken); err != nil || !acquired {
		t.Fatalf("acquire after reservation expiry = %t, %v", acquired, err)
	}
}

func TestMalformedLeaseCanBeRecovered(t *testing.T) {
	state := t.TempDir()
	election, err := New(state)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(state, "runtime", "primary.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not json\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	token, _ := election.NewToken()
	if _, acquired, err := election.TryAcquire("127.0.0.1:10001", token); err != nil || !acquired {
		t.Fatalf("recover malformed lease = %t, %v", acquired, err)
	}
}
