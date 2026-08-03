package render

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Search writes ranked verse matches in human-readable or stable plain form.
func Search(writer io.Writer, query string, results []bible.SearchResult, options Options) error {
	if options.Plain {
		for _, result := range results {
			if _, err := fmt.Fprintf(
				writer,
				"%s %d:%d\t%s\n",
				result.Book.Name,
				result.Chapter,
				result.Verse,
				result.Text,
			); err != nil {
				return err
			}
		}
		return nil
	}

	if len(results) == 0 {
		_, err := fmt.Fprintf(writer, "No results for %q.\n", query)
		return err
	}

	title := styled("Search results", ansiBold+ansiCyan, options.Color)
	translation := styled(results[0].Translation, ansiDim, options.Color)
	if _, err := fmt.Fprintf(writer, "%s · %s\n", title, translation); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%d %s for %q\n\n", len(results), matchLabel(len(results)), query); err != nil {
		return err
	}
	for _, result := range results {
		reference := fmt.Sprintf("%s %d:%d", result.Book.Name, result.Chapter, result.Verse)
		text, err := highlightSearchMatches(result.Text, result.Highlights, options.Color)
		if err != nil {
			return fmt.Errorf("render search result %s: %w", reference, err)
		}
		styledReference := styled(reference, ansiBold+ansiCyan, options.Color)
		if _, err := fmt.Fprintf(writer, "%s  %s\n", styledReference, text); err != nil {
			return err
		}
	}
	return nil
}

func highlightSearchMatches(text string, highlights []bible.TextRange, enabled bool) (string, error) {
	if !enabled || len(highlights) == 0 {
		return text, nil
	}

	var output strings.Builder
	last := 0
	for _, highlight := range highlights {
		if highlight.Start < last || highlight.End <= highlight.Start || highlight.End > len(text) {
			return "", fmt.Errorf("invalid highlight range [%d:%d]", highlight.Start, highlight.End)
		}
		if !utf8.ValidString(text[last:highlight.Start]) || !utf8.ValidString(text[highlight.Start:highlight.End]) {
			return "", fmt.Errorf("highlight range [%d:%d] splits UTF-8 text", highlight.Start, highlight.End)
		}
		output.WriteString(text[last:highlight.Start])
		output.WriteString(styled(text[highlight.Start:highlight.End], ansiBold+ansiYellow, true))
		last = highlight.End
	}
	output.WriteString(text[last:])
	return output.String(), nil
}

func matchLabel(count int) string {
	if count == 1 {
		return "match"
	}
	return "matches"
}
