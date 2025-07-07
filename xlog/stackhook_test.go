package xlog

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

// helper to clear global registry between tests
func clearRegistry() {
	registryMu.Lock()
	storages = make(map[string]GIDStorage)
	registryMu.Unlock()
}

func checkLog() {
	log.Info().Msg("test log")
}

func TestStackHook(t *testing.T) {
	const hookStorageKey = "x-request-id"

	hook, err := RegisterHook(hookStorageKey, nil)
	require.NoError(t, err)

	log.Logger = zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Hook(HookCallerFunc{}).
		Hook(hook)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		SetValue(hookStorageKey, "123")
		checkLog()
		DeleteValue(hookStorageKey)
		checkLog()
	}()

	wg.Wait()
}

func TestSimpleMap(t *testing.T) {
	m := NewSimpleMap()
	key := uint64(42)
	value := "hello"

	m.Store(key, value)
	v, ok := m.Load(key)
	if !ok {
		t.Fatalf("expected key %d to exist", key)
	}
	if v.(string) != value {
		t.Errorf("expected value %q, got %q", value, v)
	}

	m.Delete(key)
	_, ok = m.Load(key)
	if ok {
		t.Errorf("expected key %d to be deleted", key)
	}
}

func TestRegisterHook(t *testing.T) {
	clearRegistry()
	name := "test-storage"

	hook, err := RegisterHook(name, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hook.name != name {
		t.Errorf("expected hook.name %q, got %q", name, hook.name)
	}

	// duplicate registration should fail
	_, err = RegisterHook(name, nil)
	if err == nil {
		t.Errorf("expected error on duplicate registration")
	}
}

func TestSetValueDeleteValue(t *testing.T) {
	clearRegistry()
	name := "foo"
	_, err := RegisterHook(name, nil)
	if err != nil {
		t.Fatalf("RegisterHook error: %v", err)
	}

	testVal := 12345
	if err := SetValue(name, testVal); err != nil {
		t.Fatalf("SetValue error: %v", err)
	}
	gid := GetGID()
	stored, ok := storages[name].Load(gid)
	if !ok {
		t.Fatalf("value not stored for gid %d", gid)
	}
	if stored.(int) != testVal {
		t.Errorf("expected stored %d, got %v", testVal, stored)
	}

	if err := DeleteValue(name); err != nil {
		t.Fatalf("DeleteValue error: %v", err)
	}
	_, ok = storages[name].Load(gid)
	if ok {
		t.Error("expected value deleted after DeleteValue")
	}
}

func TestHookIntegration(t *testing.T) {
	clearRegistry()
	name := "xreq"
	hook, err := RegisterHook(name, nil)
	if err != nil {
		t.Fatalf("RegisterHook error: %v", err)
	}

	// set value for current goroutine
	expected := "req-123"
	if err := SetValue(name, expected); err != nil {
		t.Fatalf("SetValue error: %v", err)
	}
	defer DeleteValue(name)

	// prepare logger writing JSON to buffer
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Hook(hook)
	logger.Info().Msg("testing integration")

	// parse JSON output
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal log JSON: %v", err)
	}

	// check that our field is present
	got, ok := out[name]
	if !ok {
		t.Fatalf("field %q not found in log output: %+v", name, out)
	}
	if got != expected {
		t.Errorf("expected field %q to be %q, got %v", name, expected, got)
	}
}

// GetGID test: ensure GetGID returns non-zero and stable within same goroutine
func TestGetGID(t *testing.T) {
	gid1 := GetGID()
	gid2 := GetGID()
	if gid1 == 0 {
		t.Error("expected non-zero gid")
	}
	if gid1 != gid2 {
		t.Errorf("expected same gid, got %d and %d", gid1, gid2)
	}
}

// Test concurrency: ensure separate goroutines have different GIDs
func TestConcurrentGID(t *testing.T) {
	var wg sync.WaitGroup
	gids := make(chan uint64, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		gids <- GetGID()
	}()
	go func() {
		defer wg.Done()
		gids <- GetGID()
	}()
	wg.Wait()
	close(gids)

	collected := []uint64{<-gids, <-gids}
	if collected[0] == collected[1] {
		t.Errorf("expected different GIDs for different goroutines, got %d twice", collected[0])
	}
}

func TestUnsafe(t *testing.T) {
	EnableUnsafe()

	TestStackHook(t)
	TestSimpleMap(t)
	TestRegisterHook(t)
	TestSetValueDeleteValue(t)
	TestHookIntegration(t)
	TestGetGID(t)
	TestConcurrentGID(t)
}
