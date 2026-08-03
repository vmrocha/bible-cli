package importer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/translation"
)

const maxVPLLineBytes = 1024 * 1024

// ParseVPL parses and validates an eBible.org verse-per-line text stream.
func ParseVPL(reader io.Reader, expected translation.Expected) (bible.Dataset, error) {
	var dataset bible.Dataset
	seenReferences := make(map[string]struct{}, expected.Verses)
	seenBooks := make(map[string]struct{}, expected.Books)

	var previous bible.Verse
	var previousBook bible.Book
	var firstReference string
	var lastReference string

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), maxVPLLineBytes)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		bookCode, remainder, ok := strings.Cut(line, " ")
		if !ok || bookCode == "" {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: missing book code", lineNumber)
		}

		reference, text, ok := strings.Cut(remainder, " ")
		if !ok || reference == "" {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: missing reference", lineNumber)
		}

		book, ok := booksBySourceCode[bookCode]
		if !ok {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: unknown book code %q", lineNumber, bookCode)
		}

		chapter, verseNumber, err := parseReference(reference)
		if err != nil {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: %w", lineNumber, err)
		}

		fullReference := bookCode + " " + reference
		if _, duplicate := seenReferences[fullReference]; duplicate {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: duplicate reference %s", lineNumber, fullReference)
		}
		seenReferences[fullReference] = struct{}{}

		if _, seen := seenBooks[book.ID]; !seen {
			if len(dataset.Books) > 0 && book.Position != previousBook.Position+1 {
				return bible.Dataset{}, fmt.Errorf(
					"parse VPL line %d: book %s follows %s out of canonical order",
					lineNumber,
					bookCode,
					previousBook.SourceCode,
				)
			}
			seenBooks[book.ID] = struct{}{}
			dataset.Books = append(dataset.Books, book)
			previousBook = book
		} else if previousBook.ID != book.ID {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: book %s is not contiguous", lineNumber, bookCode)
		}

		verse := bible.Verse{
			BookID:  book.ID,
			Chapter: chapter,
			Number:  verseNumber,
			// Some versification systems retain a reference whose source text is
			// intentionally empty. Preserve that payload instead of dropping the
			// address or inventing text.
			Text: text,
		}
		if len(dataset.Verses) > 0 && previous.BookID == verse.BookID {
			if err := validateNextReference(previous, verse); err != nil {
				return bible.Dataset{}, fmt.Errorf("parse VPL line %d: %w", lineNumber, err)
			}
		} else if chapter != 1 || verseNumber != 1 {
			return bible.Dataset{}, fmt.Errorf("parse VPL line %d: book %s must begin at 1:1", lineNumber, bookCode)
		}

		if firstReference == "" {
			firstReference = fullReference
		}
		lastReference = fullReference
		dataset.Verses = append(dataset.Verses, verse)
		previous = verse
	}
	if err := scanner.Err(); err != nil {
		return bible.Dataset{}, fmt.Errorf("scan VPL source: %w", err)
	}

	if len(dataset.Books) != expected.Books {
		return bible.Dataset{}, fmt.Errorf("validate VPL source: got %d books, want %d", len(dataset.Books), expected.Books)
	}
	if len(dataset.Verses) != expected.Verses {
		return bible.Dataset{}, fmt.Errorf("validate VPL source: got %d verses, want %d", len(dataset.Verses), expected.Verses)
	}
	if firstReference != expected.FirstReference {
		return bible.Dataset{}, fmt.Errorf("validate VPL source: first reference is %q, want %q", firstReference, expected.FirstReference)
	}
	if lastReference != expected.LastReference {
		return bible.Dataset{}, fmt.Errorf("validate VPL source: last reference is %q, want %q", lastReference, expected.LastReference)
	}

	return dataset, nil
}

func parseReference(reference string) (int, int, error) {
	chapterText, verseText, ok := strings.Cut(reference, ":")
	if !ok || chapterText == "" || verseText == "" || strings.Contains(verseText, ":") {
		return 0, 0, fmt.Errorf("invalid reference %q", reference)
	}

	chapter, err := strconv.Atoi(chapterText)
	if err != nil || chapter <= 0 {
		return 0, 0, fmt.Errorf("invalid chapter in reference %q", reference)
	}
	verse, err := strconv.Atoi(verseText)
	if err != nil || verse <= 0 {
		return 0, 0, fmt.Errorf("invalid verse in reference %q", reference)
	}

	return chapter, verse, nil
}

func validateNextReference(previous, current bible.Verse) error {
	if current.Chapter == previous.Chapter && current.Number == previous.Number+1 {
		return nil
	}
	if current.Chapter == previous.Chapter+1 && current.Number == 1 {
		return nil
	}
	return fmt.Errorf(
		"reference %d:%d does not follow %d:%d",
		current.Chapter,
		current.Number,
		previous.Chapter,
		previous.Number,
	)
}
