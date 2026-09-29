package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const schemaVersion = 1

var translationAliases = map[string]string{
	"engwebp": "engwebp",
	"webp":    "engwebp",
}

// RegisterTranslation adds a locally available translation id and aliases.
// Registration is intended to happen during package initialization.
func RegisterTranslation(id string, aliases ...string) {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" {
		panic("translation id is required")
	}
	values := append([]string{id}, aliases...)
	for _, value := range values {
		key := strings.ToLower(strings.TrimSpace(value))
		if key == "" {
			panic("translation alias is required")
		}
		if existing, ok := translationAliases[key]; ok && existing != id {
			panic(fmt.Sprintf("translation alias %q is already registered for %s", key, existing))
		}
		translationAliases[key] = id
	}
}

// Preferences contains the user-configurable CLI defaults.
type Preferences struct {
	Version     int    `json:"version"`
	Translation string `json:"translation"`
	Plain       bool   `json:"plain"`
	Color       bool   `json:"color"`
}

type document struct {
	Version     *int    `json:"version"`
	Translation *string `json:"translation"`
	Plain       *bool   `json:"plain"`
	Color       *bool   `json:"color"`
}

// Defaults returns the preferences used when no configuration file exists.
func Defaults() Preferences {
	return Preferences{
		Version:     schemaVersion,
		Translation: "engwebp",
		Color:       true,
	}
}

// Validate checks that preferences are supported by this application version.
func (preferences Preferences) Validate() error {
	if preferences.Version != schemaVersion {
		return fmt.Errorf("unsupported configuration version %d", preferences.Version)
	}
	if normalized, ok := NormalizeTranslation(preferences.Translation); !ok || normalized != preferences.Translation {
		return fmt.Errorf("translation is not available: %s", preferences.Translation)
	}
	return nil
}

// NormalizeTranslation resolves a supported translation id or abbreviation.
func NormalizeTranslation(value string) (string, bool) {
	normalized, ok := translationAliases[strings.ToLower(strings.TrimSpace(value))]
	return normalized, ok
}

// Store persists one configuration file.
type Store struct {
	path string
}

// NewStore creates a store for an already resolved absolute path.
func NewStore(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("configuration path is required")
	}
	if !filepath.IsAbs(path) {
		return nil, errors.New("configuration path must be absolute")
	}
	return &Store{path: filepath.Clean(path)}, nil
}

// Path returns the configuration file location.
func (store *Store) Path() string {
	return store.path
}

// Load returns defaults when the configuration file does not exist.
func (store *Store) Load() (Preferences, error) {
	file, err := os.Open(store.path)
	if errors.Is(err, os.ErrNotExist) {
		return Defaults(), nil
	}
	if err != nil {
		return Preferences{}, fmt.Errorf("open configuration: %w", err)
	}

	var persisted document
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	decodeErr := decoder.Decode(&persisted)
	if decodeErr == nil {
		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			if err == nil {
				decodeErr = errors.New("multiple JSON values")
			} else {
				decodeErr = err
			}
		}
	}
	closeErr := file.Close()
	if decodeErr != nil {
		return Preferences{}, fmt.Errorf("decode configuration %s: %w", store.path, decodeErr)
	}
	if closeErr != nil {
		return Preferences{}, fmt.Errorf("close configuration: %w", closeErr)
	}
	if persisted.Version == nil {
		return Preferences{}, fmt.Errorf("decode configuration %s: version is required", store.path)
	}

	preferences := Defaults()
	preferences.Version = *persisted.Version
	if persisted.Translation != nil {
		preferences.Translation = *persisted.Translation
	}
	if persisted.Plain != nil {
		preferences.Plain = *persisted.Plain
	}
	if persisted.Color != nil {
		preferences.Color = *persisted.Color
	}
	if err := preferences.Validate(); err != nil {
		return Preferences{}, fmt.Errorf("validate configuration %s: %w", store.path, err)
	}
	return preferences, nil
}

// Save atomically replaces the configuration with owner-only permissions.
func (store *Store) Save(preferences Preferences) error {
	if err := preferences.Validate(); err != nil {
		return fmt.Errorf("save configuration: %w", err)
	}

	directory := filepath.Dir(store.path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create configuration directory: %w", err)
	}
	temporary, err := os.CreateTemp(directory, ".config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary configuration: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return fmt.Errorf("protect temporary configuration: %w", err)
	}
	encoder := json.NewEncoder(temporary)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(preferences); err != nil {
		temporary.Close()
		return fmt.Errorf("encode configuration: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("sync temporary configuration: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary configuration: %w", err)
	}
	if err := os.Rename(temporaryPath, store.path); err != nil {
		return fmt.Errorf("publish configuration: %w", err)
	}
	return nil
}

// Reset removes the configuration file. Missing files are already reset.
func (store *Store) Reset() error {
	if err := os.Remove(store.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("reset configuration: %w", err)
	}
	return nil
}
