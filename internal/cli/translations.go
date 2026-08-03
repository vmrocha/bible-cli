package cli

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/render"
)

// TranslationReader lists bundled translation metadata.
type TranslationReader interface {
	Translations(context.Context) ([]bible.Translation, error)
	Close() error
}

// TranslationReaderFactory opens translation metadata on demand.
type TranslationReaderFactory func(context.Context) (TranslationReader, error)

// WithTranslationReaderFactory enables translation discovery.
func WithTranslationReaderFactory(factory TranslationReaderFactory) Option {
	return func(configuration *configuration) {
		configuration.translationFactory = factory
	}
}

func newTranslationsCommand(
	factory TranslationReaderFactory,
	settings *outputSettings,
	isTerminal func(io.Writer) bool,
) *cobra.Command {
	return &cobra.Command{
		Use:   "translations",
		Short: "List bundled translations and attribution",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if factory == nil {
				return errors.New("translation metadata is unavailable")
			}
			reader, err := factory(command.Context())
			if err != nil {
				return err
			}
			translations, readErr := reader.Translations(command.Context())
			closeErr := reader.Close()
			if readErr != nil {
				return readErr
			}
			if closeErr != nil {
				return closeErr
			}
			return render.Translations(
				command.OutOrStdout(),
				translations,
				renderOptions(command, settings, isTerminal),
			)
		},
	}
}
