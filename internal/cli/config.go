package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/config"
	"github.com/vmrocha/bible-cli/internal/render"
)

// PreferenceStore persists CLI defaults.
type PreferenceStore interface {
	Path() string
	Load() (config.Preferences, error)
	Save(config.Preferences) error
	Reset() error
}

// WithPreferenceStore enables persistent CLI preferences.
func WithPreferenceStore(store PreferenceStore) Option {
	return func(configuration *configuration) {
		configuration.preferenceStore = store
	}
}

func newConfigCommand(
	store PreferenceStore,
	settings *outputSettings,
	isTerminal func(io.Writer) bool,
) *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Inspect and update persistent preferences",
	}
	command.AddCommand(newConfigPathCommand(store))
	command.AddCommand(newConfigShowCommand(store, settings, isTerminal))
	command.AddCommand(newConfigSetCommand(store))
	command.AddCommand(newConfigResetCommand(store))
	return command
}

func newConfigPathCommand(store PreferenceStore) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the effective configuration path",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if store == nil {
				return errors.New("configuration is unavailable")
			}
			_, err := fmt.Fprintln(command.OutOrStdout(), store.Path())
			return err
		},
	}
}

func newConfigShowCommand(
	store PreferenceStore,
	settings *outputSettings,
	isTerminal func(io.Writer) bool,
) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the effective persistent preferences",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if store == nil {
				return errors.New("configuration is unavailable")
			}
			preferences, err := store.Load()
			if err != nil {
				return err
			}
			return render.Preferences(
				command.OutOrStdout(),
				store.Path(),
				preferences,
				renderOptions(command, settings, isTerminal),
			)
		},
	}
}

func newConfigSetCommand(store PreferenceStore) *cobra.Command {
	return &cobra.Command{
		Use:       "set <preference> <value>",
		Short:     "Save a persistent preference",
		Args:      cobra.ExactArgs(2),
		ValidArgs: []string{"translation", "plain", "color"},
		RunE: func(command *cobra.Command, args []string) error {
			if store == nil {
				return errors.New("configuration is unavailable")
			}
			preferences, err := store.Load()
			if err != nil {
				return err
			}
			key := strings.ToLower(args[0])
			value, err := updatePreference(&preferences, key, args[1])
			if err != nil {
				return err
			}
			if err := store.Save(preferences); err != nil {
				return err
			}
			_, err = fmt.Fprintf(command.OutOrStdout(), "saved %s=%s\n", key, value)
			return err
		},
	}
}

func newConfigResetCommand(store PreferenceStore) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Remove saved preferences",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if store == nil {
				return errors.New("configuration is unavailable")
			}
			if err := store.Reset(); err != nil {
				return err
			}
			_, err := fmt.Fprintln(command.OutOrStdout(), "configuration reset")
			return err
		},
	}
}

func updatePreference(preferences *config.Preferences, key, rawValue string) (string, error) {
	switch key {
	case "translation":
		value := strings.ToLower(rawValue)
		if value == "webp" {
			value = "engwebp"
		}
		if value != "engwebp" {
			return "", fmt.Errorf("translation is not available: %s", rawValue)
		}
		preferences.Translation = value
		return value, nil
	case "plain":
		value, err := strictBoolean(rawValue)
		if err != nil {
			return "", fmt.Errorf("plain: %w", err)
		}
		preferences.Plain = value
		return fmt.Sprint(value), nil
	case "color":
		value, err := strictBoolean(rawValue)
		if err != nil {
			return "", fmt.Errorf("color: %w", err)
		}
		preferences.Color = value
		return fmt.Sprint(value), nil
	default:
		return "", fmt.Errorf("unknown preference %q; expected translation, plain, or color", key)
	}
}

func strictBoolean(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("expected true or false, got %q", value)
	}
}
