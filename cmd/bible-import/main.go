package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/importer"
	"github.com/vmrocha/bible-cli/internal/storage"
	"github.com/vmrocha/bible-cli/internal/translation"
)

func main() {
	if err := newCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	var manifestPath string
	var archivePath string
	var outputPath string

	command := &cobra.Command{
		Use:           "bible-import",
		Short:         "Build a Bible CLI translation database",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if manifestPath == "" || archivePath == "" || outputPath == "" {
				return fmt.Errorf("--manifest, --archive, and --output are required")
			}
			return run(command.Context(), command, manifestPath, archivePath, outputPath)
		},
	}
	command.Flags().StringVar(&manifestPath, "manifest", "", "translation manifest JSON")
	command.Flags().StringVar(&archivePath, "archive", "", "publisher source ZIP")
	command.Flags().StringVar(&outputPath, "output", "", "new SQLite database path")

	return command
}

func run(ctx context.Context, command *cobra.Command, manifestPath, archivePath, outputPath string) error {
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("open translation manifest: %w", err)
	}
	manifest, err := translation.DecodeManifest(manifestFile)
	closeErr := manifestFile.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return fmt.Errorf("close translation manifest: %w", closeErr)
	}

	dataset, err := importer.LoadArchive(archivePath, manifest)
	if err != nil {
		return err
	}
	if err := storage.Generate(ctx, outputPath, manifest, dataset); err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		command.OutOrStdout(),
		"wrote %s (%d books, %d verses)\n",
		outputPath,
		len(dataset.Books),
		len(dataset.Verses),
	)
	return err
}
