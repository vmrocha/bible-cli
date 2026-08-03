package cli

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/reference"
	"github.com/vmrocha/bible-cli/internal/render"
)

// PassageReader is the read behavior required by the CLI.
type PassageReader interface {
	Read(context.Context, reference.Query) (bible.Passage, error)
	Navigate(context.Context, reference.Query, int) (reference.Query, error)
	Close() error
}

// ReaderFactory opens the offline Bible reader on demand.
type ReaderFactory func(context.Context) (PassageReader, error)

type configuration struct {
	readerFactory      ReaderFactory
	searchFactory      SearcherFactory
	translationFactory TranslationReaderFactory
	randomFactory      RandomReaderFactory
	preferenceStore    PreferenceStore
	isTerminal         func(io.Writer) bool
}

// Option configures optional root-command dependencies.
type Option func(*configuration)

// WithReaderFactory enables commands that access the embedded Bible.
func WithReaderFactory(factory ReaderFactory) Option {
	return func(configuration *configuration) {
		configuration.readerFactory = factory
	}
}

func newReadCommand(factory ReaderFactory, settings *outputSettings, isTerminal func(io.Writer) bool) *cobra.Command {
	var next bool
	var previous bool
	command := &cobra.Command{
		Use:     "read <reference>",
		Short:   "Read a chapter, verse, or verse range",
		Example: "  bible read John 3\n  bible read John 3 --next\n  bible read 'John 3:16'\n  bible read 'John 3:16-21' --plain",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			query, err := reference.Parse(strings.Join(args, " "))
			if err != nil {
				return err
			}
			if next && previous {
				return errors.New("--next and --previous cannot be used together")
			}
			if (next || previous) && !query.IsChapter() {
				return errors.New("--next and --previous require a chapter reference")
			}
			if factory == nil {
				return errors.New("Bible reader is unavailable")
			}

			reader, err := factory(command.Context())
			if err != nil {
				return err
			}
			if next || previous {
				direction := 1
				if previous {
					direction = -1
				}
				query, err = reader.Navigate(command.Context(), query, direction)
				if err != nil {
					_ = reader.Close()
					return err
				}
			}
			passage, readErr := reader.Read(command.Context(), query)
			closeErr := reader.Close()
			if readErr != nil {
				return readErr
			}
			if closeErr != nil {
				return closeErr
			}

			return render.Passage(command.OutOrStdout(), passage, renderOptions(command, settings, isTerminal))
		},
	}
	command.Flags().BoolVar(&next, "next", false, "read the next chapter")
	command.Flags().BoolVar(&previous, "previous", false, "read the previous chapter")
	return command
}

func renderOptions(command *cobra.Command, settings *outputSettings, isTerminal func(io.Writer) bool) render.Options {
	interactive := isTerminal != nil && isTerminal(command.OutOrStdout())
	plain := settings.plain || !interactive
	return render.Options{
		Plain: plain,
		Color: interactive && !plain && !settings.noColor,
	}
}
