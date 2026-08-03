package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name       string
		configHome string
		xdgHome    string
		userHome   string
		want       string
	}{
		{
			name:       "application override wins",
			configHome: filepath.Join(string(filepath.Separator), "custom", "bible"),
			xdgHome:    filepath.Join(string(filepath.Separator), "xdg"),
			userHome:   filepath.Join(string(filepath.Separator), "home", "reader"),
			want:       filepath.Join(string(filepath.Separator), "custom", "bible", "config.json"),
		},
		{
			name:     "XDG home is shared across macOS and Linux",
			xdgHome:  filepath.Join(string(filepath.Separator), "xdg"),
			userHome: filepath.Join(string(filepath.Separator), "home", "reader"),
			want:     filepath.Join(string(filepath.Separator), "xdg", "bible-cli", "config.json"),
		},
		{
			name:     "home fallback uses dot config",
			userHome: filepath.Join(string(filepath.Separator), "home", "reader"),
			want:     filepath.Join(string(filepath.Separator), "home", "reader", ".config", "bible-cli", "config.json"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ResolvePath(test.configHome, test.xdgHome, test.userHome)
			if err != nil {
				t.Fatalf("ResolvePath: %v", err)
			}
			if got != test.want {
				t.Fatalf("ResolvePath = %q, want %q", got, test.want)
			}
		})
	}
}

func TestResolvePathRejectsRelativeLocations(t *testing.T) {
	for _, test := range []struct {
		name       string
		configHome string
		xdgHome    string
		userHome   string
		want       string
	}{
		{"config override", "relative", "", "/home/reader", "BIBLE_CLI_CONFIG_HOME"},
		{"XDG override", "", "relative", "/home/reader", "XDG_CONFIG_HOME"},
		{"user home", "", "", "relative", "user home"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := ResolvePath(test.configHome, test.xdgHome, test.userHome)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestResolvePathRequiresAHome(t *testing.T) {
	if _, err := ResolvePath("", "", ""); err == nil {
		t.Fatal("ResolvePath unexpectedly accepted an empty home")
	}
}
