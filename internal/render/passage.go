package render

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Passage writes a Bible passage in human-readable or stable plain form.
func Passage(writer io.Writer, passage bible.Passage, options Options) error {
	if len(passage.Verses) == 0 {
		return errors.New("render passage: no verses")
	}
	if options.Plain {
		for _, verse := range passage.Verses {
			if _, err := fmt.Fprintf(
				writer,
				"%s %d:%d\t%s\n",
				passage.Book.Name,
				verse.Chapter,
				verse.Number,
				verse.Text,
			); err != nil {
				return err
			}
		}
		return nil
	}

	title := styled(passageTitle(passage), ansiBold+ansiCyan, options.Color)
	translation := styled(passage.Translation, ansiDim, options.Color)
	if _, err := fmt.Fprintf(writer, "%s · %s\n\n", title, translation); err != nil {
		return err
	}
	width := len(strconv.Itoa(passage.Verses[len(passage.Verses)-1].Number))
	for _, verse := range passage.Verses {
		number := styled(fmt.Sprintf("%*d", width, verse.Number), ansiBold+ansiCyan, options.Color)
		if _, err := fmt.Fprintf(writer, "%s  %s\n", number, verse.Text); err != nil {
			return err
		}
	}
	return nil
}

func passageTitle(passage bible.Passage) string {
	if passage.StartVerse == 0 {
		return fmt.Sprintf("%s %d", passage.Book.Name, passage.Chapter)
	}
	if passage.StartVerse == passage.EndVerse {
		return fmt.Sprintf("%s %d:%d", passage.Book.Name, passage.Chapter, passage.StartVerse)
	}
	return fmt.Sprintf(
		"%s %d:%d–%d",
		passage.Book.Name,
		passage.Chapter,
		passage.StartVerse,
		passage.EndVerse,
	)
}
