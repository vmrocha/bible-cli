package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/config"
)

func TestPreferences(t *testing.T) {
	var output bytes.Buffer
	if err := Preferences(&output, "/home/reader/config.json", config.Defaults(), Options{}); err != nil {
		t.Fatalf("Preferences: %v", err)
	}
	want := "Configuration\n" +
		"Path: /home/reader/config.json\n" +
		"Translation: engwebp\n" +
		"Plain output: false\n" +
		"Color: true\n"
	if output.String() != want {
		t.Fatalf("unexpected preferences output: %q", output.String())
	}
}

func TestStyledPreferences(t *testing.T) {
	var output bytes.Buffer
	if err := Preferences(&output, "/config.json", config.Defaults(), Options{Color: true}); err != nil {
		t.Fatalf("Preferences: %v", err)
	}
	if !strings.Contains(output.String(), "\x1b[1m\x1b[36mConfiguration\x1b[0m") {
		t.Fatalf("configuration heading was not styled: %q", output.String())
	}
}

func TestPlainPreferencesSanitizesFields(t *testing.T) {
	var output bytes.Buffer
	preferences := config.Defaults()
	if err := Preferences(&output, "/path\twith\nspaces", preferences, Options{Plain: true}); err != nil {
		t.Fatalf("Preferences: %v", err)
	}
	want := "path\t/path with spaces\n" +
		"translation\tengwebp\n" +
		"plain\tfalse\n" +
		"color\ttrue\n"
	if output.String() != want {
		t.Fatalf("unexpected plain preferences output: %q", output.String())
	}
}
