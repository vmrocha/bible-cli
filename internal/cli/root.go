package cli

import (
	"io"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/vmrocha/bible-cli/internal/buildinfo"
)

type outputSettings struct {
	plain   bool
	noColor bool
}

// New constructs the root command with the supplied build metadata.
func New(info buildinfo.Info, options ...Option) *cobra.Command {
	configuration := configuration{isTerminal: terminalWriter}
	for _, option := range options {
		option(&configuration)
	}
	settings := &outputSettings{}

	command := &cobra.Command{
		Use:           "bible",
		Short:         "Read the Bible from your terminal",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       info.Version,
	}

	command.SetVersionTemplate("bible {{.Version}}\n")
	command.PersistentFlags().BoolVar(&settings.plain, "plain", false, "emit stable plain output")
	command.PersistentFlags().BoolVar(&settings.noColor, "no-color", false, "disable terminal colors")
	command.PersistentPreRunE = func(selected *cobra.Command, _ []string) error {
		if configuration.preferenceStore == nil || skipsSavedPreferences(selected) {
			return nil
		}
		preferences, err := configuration.preferenceStore.Load()
		if err != nil {
			return err
		}
		if !command.PersistentFlags().Changed("plain") {
			settings.plain = preferences.Plain
		}
		if !command.PersistentFlags().Changed("no-color") {
			settings.noColor = !preferences.Color
		}
		return nil
	}
	command.AddCommand(newBooksCommand(settings, configuration.isTerminal))
	command.AddCommand(newReadCommand(configuration.readerFactory, settings, configuration.isTerminal))
	command.AddCommand(newSearchCommand(configuration.searchFactory, settings, configuration.isTerminal))
	command.AddCommand(newTranslationsCommand(configuration.translationFactory, settings, configuration.isTerminal))
	command.AddCommand(newRandomCommand(configuration.randomFactory, settings, configuration.isTerminal))
	command.AddCommand(newCompletionCommand())
	command.AddCommand(newConfigCommand(configuration.preferenceStore, settings, configuration.isTerminal))
	command.AddCommand(newVersionCommand(info))

	return command
}

func skipsSavedPreferences(command *cobra.Command) bool {
	for current := command; current != nil; current = current.Parent() {
		if current.Name() == "config" {
			return true
		}
	}
	return command.Name() == "completion" || command.Name() == "version"
}

type fileDescriptor interface {
	Fd() uintptr
}

func terminalWriter(writer io.Writer) bool {
	output, ok := writer.(fileDescriptor)
	if !ok {
		return false
	}
	return isatty.IsTerminal(output.Fd()) || isatty.IsCygwinTerminal(output.Fd())
}
