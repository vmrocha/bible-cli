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
			return render.Books(
				command.OutOrStdout(),
				canon.ProtestantBooks(),
				renderOptions(command, settings, isTerminal),
			)
		},
	}
	return command
}
