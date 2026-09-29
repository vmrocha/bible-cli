package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/render"
)

const (
	defaultSearchLimit = 20
	maxSearchLimit     = 100
)

// Searcher is the full-text search behavior required by the CLI.
type Searcher interface {
	Search(context.Context, string, int) ([]bible.SearchResult, error)
	Close() error
}

// SearcherFactory opens the offline Bible search index on demand.
type SearcherFactory func(context.Context, string) (Searcher, error)

// WithSearcherFactory enables commands that access the embedded search index.
func WithSearcherFactory(factory SearcherFactory) Option {
	return func(configuration *configuration) {
		configuration.searchFactory = factory
	}
}

func newSearchCommand(factory SearcherFactory, settings *outputSettings, isTerminal func(io.Writer) bool) *cobra.Command {
	limit := defaultSearchLimit
	command := &cobra.Command{
		Use:     "search <terms>",
		Short:   "Search every verse offline",
		Example: "  bible search 'living water'\n  bible search 'faith hope love' --limit 10\n  bible search 'kingdom of God' --plain",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			query := strings.Join(strings.Fields(strings.Join(args, " ")), " ")
			if query == "" {
				return errors.New("search query cannot be empty")
			}
			if limit < 1 || limit > maxSearchLimit {
				return fmt.Errorf("--limit must be between 1 and %d", maxSearchLimit)
			}
			if factory == nil {
				return errors.New("Bible search is unavailable")
			}

			searcher, err := factory(command.Context(), settings.translation)
			if err != nil {
				return err
			}
			results, searchErr := searcher.Search(command.Context(), query, limit)
			closeErr := searcher.Close()
			if searchErr != nil {
				return searchErr
			}
			if closeErr != nil {
				return closeErr
			}

			return render.Search(
				command.OutOrStdout(),
				query,
				results,
				renderOptions(command, settings, isTerminal),
			)
		},
	}
	command.Flags().IntVarP(&limit, "limit", "n", defaultSearchLimit, "maximum number of results (1-100)")
	return command
}
