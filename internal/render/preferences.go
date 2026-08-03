package render

import (
	"fmt"
	"io"

	"github.com/vmrocha/bible-cli/internal/config"
)

// Preferences writes the effective persistent configuration.
func Preferences(writer io.Writer, path string, preferences config.Preferences, options Options) error {
	if options.Plain {
		for _, field := range []struct {
			key   string
			value any
		}{
			{"path", plainMetadataField(path)},
			{"translation", plainMetadataField(preferences.Translation)},
			{"plain", preferences.Plain},
			{"color", preferences.Color},
		} {
			if _, err := fmt.Fprintf(writer, "%s\t%v\n", field.key, field.value); err != nil {
				return err
			}
		}
		return nil
	}

	heading := styled("Configuration", ansiBold+ansiCyan, options.Color)
	if _, err := fmt.Fprintln(writer, heading); err != nil {
		return err
	}
	for _, field := range []struct {
		label string
		value any
	}{
		{"Path", path},
		{"Translation", preferences.Translation},
		{"Plain output", preferences.Plain},
		{"Color", preferences.Color},
	} {
		label := styled(field.label+":", ansiDim, options.Color)
		if _, err := fmt.Fprintf(writer, "%s %v\n", label, field.value); err != nil {
			return err
		}
	}
	return nil
}
