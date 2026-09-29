package cli

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/canon"
	"github.com/vmrocha/bible-cli/internal/render"
)

func newBooksCommand(settings *outputSettings, isTerminal func(io.Writer) bool) *cobra.Command {
	command := &cobra.Command{
		Use:   "books",
		Short: "List books and accepted aliases",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			books := canon.ProtestantBooks()
			if settings.translation == "ptnvi" {
				books = canon.PortugueseProtestantBooks()
			}
			return render.Books(
				command.OutOrStdout(),
				books,
				renderOptions(command, settings, isTerminal),
			)
		},
	}
	return command
}
