package main

import (
	"context"
	"fmt"
	"os"

	"github.com/vmrocha/bible-cli/internal/buildinfo"
	"github.com/vmrocha/bible-cli/internal/cli"
	"github.com/vmrocha/bible-cli/internal/config"
	"github.com/vmrocha/bible-cli/internal/storage"
)

func main() {
	configurationPath, err := config.DefaultPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	preferenceStore, err := config.NewStore(configurationPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	command := cli.New(
		buildinfo.Current(),
		cli.WithReaderFactory(func(ctx context.Context, translation string) (cli.PassageReader, error) {
			return storage.OpenEmbedded(ctx, translation)
		}),
		cli.WithSearcherFactory(func(ctx context.Context, translation string) (cli.Searcher, error) {
			return storage.OpenEmbedded(ctx, translation)
		}),
		cli.WithTranslationReaderFactory(func(ctx context.Context) (cli.TranslationReader, error) {
			return storage.OpenEmbedded(ctx)
		}),
		cli.WithRandomReaderFactory(func(ctx context.Context, translation string) (cli.RandomReader, error) {
			return storage.OpenEmbedded(ctx, translation)
		}),
		cli.WithPreferenceStore(preferenceStore),
	)
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
