package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/buildinfo"
)

func newVersionCommand(info buildinfo.Info) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build information",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(
				command.OutOrStdout(),
				"bible %s\ncommit: %s\nbuilt: %s\n",
				info.Version,
				info.Commit,
				info.Date,
			)
			return err
		},
	}
}
