package cli

import (
	"context"
	"crypto/rand"
	"errors"
	"io"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/render"
)

// RandomReader selects one verse using caller-provided entropy.
type RandomReader interface {
	Random(context.Context, io.Reader) (bible.Passage, error)
	Close() error
}

// RandomReaderFactory opens random verse access on demand.
type RandomReaderFactory func(context.Context) (RandomReader, error)

// WithRandomReaderFactory enables random verse discovery.
func WithRandomReaderFactory(factory RandomReaderFactory) Option {
	return func(configuration *configuration) {
		configuration.randomFactory = factory
	}
}

func newRandomCommand(
	factory RandomReaderFactory,
	settings *outputSettings,
	isTerminal func(io.Writer) bool,
) *cobra.Command {
	return &cobra.Command{
		Use:   "random",
		Short: "Read one randomly selected verse",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if factory == nil {
				return errors.New("random verse reader is unavailable")
			}
			reader, err := factory(command.Context())
			if err != nil {
				return err
			}
			passage, randomErr := reader.Random(command.Context(), rand.Reader)
			closeErr := reader.Close()
			if randomErr != nil {
				return randomErr
			}
			if closeErr != nil {
				return closeErr
			}
			return render.Passage(
				command.OutOrStdout(),
				passage,
				renderOptions(command, settings, isTerminal),
			)
		},
	}
}
