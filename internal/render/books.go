package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/vmrocha/bible-cli/internal/canon"
)

// Books writes the accepted book names, source codes, and aliases.
func Books(writer io.Writer, entries []canon.Entry, options Options) error {
	for _, entry := range entries {
		aliases := strings.Join(entry.Aliases, ", ")
		if options.Plain {
			if _, err := fmt.Fprintf(
				writer,
				"%s\t%s\t%s\t%s\n",
				entry.Book.ID,
				entry.Book.Name,
				entry.Book.SourceCode,
				strings.Join(entry.Aliases, ","),
			); err != nil {
				return err
			}
			continue
		}

		if entry.Book.Position == 1 {
			heading := styled("Old Testament", ansiBold+ansiCyan, options.Color)
			if _, err := fmt.Fprintln(writer, heading); err != nil {
				return err
			}
		}
		if entry.Book.Position == 40 {
			heading := styled("New Testament", ansiBold+ansiCyan, options.Color)
			if _, err := fmt.Fprintln(writer, "\n"+heading); err != nil {
				return err
			}
		}
		position := styled(fmt.Sprintf("%2d", entry.Book.Position), ansiDim, options.Color)
		bookAliases := styled(aliases, ansiDim, options.Color)
		if _, err := fmt.Fprintf(
			writer,
			"%s  %-17s %-3s  %s\n",
			position,
			entry.Book.Name,
			entry.Book.SourceCode,
			bookAliases,
		); err != nil {
			return err
		}
	}
	return nil
}
