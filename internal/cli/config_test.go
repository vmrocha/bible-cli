package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/config"
)

type stubPreferenceStore struct {
	path        string
	preferences config.Preferences
	loadErr     error
	saved       *config.Preferences
	reset       bool
}

func (store *stubPreferenceStore) Path() string {
	return store.path
}

func (store *stubPreferenceStore) Load() (config.Preferences, error) {
	return store.preferences, store.loadErr
}

func (store *stubPreferenceStore) Save(preferences config.Preferences) error {
	store.saved = &preferences
	store.preferences = preferences
	return nil
}

func (store *stubPreferenceStore) Reset() error {
	store.reset = true
	return nil
}

func preferenceStore() *stubPreferenceStore {
	return &stubPreferenceStore{
		path:        "/home/reader/.config/bible-cli/config.json",
		preferences: config.Defaults(),
	}
}

func TestConfigPath(t *testing.T) {
	store := preferenceStore()
	output, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, "config", "path")
	if err != nil {
		t.Fatalf("execute config path: %v", err)
	}
	if output != store.path+"\n" {
		t.Fatalf("unexpected path output: %q", output)
	}
}

func TestConfigShowPlain(t *testing.T) {
	store := preferenceStore()
	store.preferences.Plain = true
	store.preferences.Color = false
	output, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, "--plain", "config", "show")
	if err != nil {
		t.Fatalf("execute config show: %v", err)
	}
	for _, expected := range []string{
		"path\t" + store.path,
		"translation\tengwebp",
		"plain\ttrue",
		"color\tfalse",
	} {
		if !strings.Contains(output, expected+"\n") {
			t.Errorf("config output does not contain %q: %q", expected, output)
		}
	}
}

func TestConfigSet(t *testing.T) {
	for _, test := range []struct {
		key   string
		value string
		check func(config.Preferences) bool
	}{
		{"translation", "webp", func(p config.Preferences) bool { return p.Translation == "engwebp" }},
		{"plain", "true", func(p config.Preferences) bool { return p.Plain }},
		{"color", "false", func(p config.Preferences) bool { return !p.Color }},
	} {
		t.Run(test.key, func(t *testing.T) {
			store := preferenceStore()
			output, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, "config", "set", test.key, test.value)
			if err != nil {
				t.Fatalf("execute config set: %v", err)
			}
			if store.saved == nil || !test.check(*store.saved) {
				t.Fatalf("preference was not saved: %#v", store.saved)
			}
			if !strings.HasPrefix(output, "saved "+test.key+"=") {
				t.Fatalf("unexpected set output: %q", output)
			}
		})
	}
}

func TestConfigSetRejectsInvalidValues(t *testing.T) {
	for _, arguments := range [][]string{
		{"config", "set", "plain", "yes"},
		{"config", "set", "translation", "unknown"},
		{"config", "set", "unknown", "true"},
	} {
		store := preferenceStore()
		if _, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, arguments...); err == nil {
			t.Fatalf("expected %v to fail", arguments)
		}
		if store.saved != nil {
			t.Fatalf("invalid preference was saved for %v", arguments)
		}
	}
}

func TestSavedPreferencesApplyAndFlagsOverrideThem(t *testing.T) {
	store := preferenceStore()
	store.preferences.Plain = true
	store.preferences.Color = false

	output, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, "books")
	if err != nil {
		t.Fatalf("execute books with saved preferences: %v", err)
	}
	if !strings.Contains(output, "john\tJohn\tJOH\tJn,Jhn") {
		t.Fatalf("saved plain preference was not applied: %q", output)
	}

	output, err = executeWithOptions(
		t,
		[]Option{WithPreferenceStore(store)},
		"--plain=false",
		"--no-color=false",
		"books",
	)
	if err != nil {
		t.Fatalf("execute books with explicit overrides: %v", err)
	}
	if !strings.Contains(output, "Old Testament") || !strings.Contains(output, "\x1b[") {
		t.Fatalf("explicit false flags did not override saved preferences: %q", output)
	}
}

func TestConfigurationCanBeResetWhenLoadingFails(t *testing.T) {
	store := preferenceStore()
	store.loadErr = errors.New("configuration is damaged")

	if _, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, "books"); err == nil {
		t.Fatal("expected a normal command to report the load error")
	}
	output, err := executeWithOptions(t, []Option{WithPreferenceStore(store)}, "config", "reset")
	if err != nil {
		t.Fatalf("execute config reset: %v", err)
	}
	if !store.reset || output != "configuration reset\n" {
		t.Fatalf("configuration was not reset: reset=%v output=%q", store.reset, output)
	}
}
