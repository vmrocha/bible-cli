package reference

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vmrocha/bible-cli/internal/canon"
)

// Query is a parsed, not-yet-resolved Bible reference.
type Query struct {
	Book       string
	Chapter    int
	StartVerse int
	EndVerse   int
}

// IsChapter reports whether the query requests a complete chapter.
func (query Query) IsChapter() bool {
	return query.StartVerse == 0
}

// Parse accepts references such as "John 3", "John 3:16", and
// "John 3:16-21".
func Parse(input string) (Query, error) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) < 2 {
		return Query{}, fmt.Errorf("invalid reference %q: expected a book and chapter", input)
	}

	book := strings.Join(fields[:len(fields)-1], " ")
	locator := fields[len(fields)-1]
	chapterText, verseText, hasVerses := strings.Cut(locator, ":")

	chapter, err := positiveNumber(chapterText, "chapter", input)
	if err != nil {
		return Query{}, err
	}

	query := Query{
		Book:    normalizeBook(book),
		Chapter: chapter,
	}
	if !hasVerses {
		return query, nil
	}
	if verseText == "" || strings.Contains(verseText, ":") {
		return Query{}, fmt.Errorf("invalid reference %q: malformed verse", input)
	}

	startText, endText, hasRange := strings.Cut(verseText, "-")
	start, err := positiveNumber(startText, "verse", input)
	if err != nil {
		return Query{}, err
	}
	end := start
	if hasRange {
		if endText == "" || strings.Contains(endText, "-") {
			return Query{}, fmt.Errorf("invalid reference %q: malformed verse range", input)
		}
		end, err = positiveNumber(endText, "verse", input)
		if err != nil {
			return Query{}, err
		}
		if end < start {
			return Query{}, fmt.Errorf("invalid reference %q: range ends before it starts", input)
		}
	}

	query.StartVerse = start
	query.EndVerse = end
	return query, nil
}

func positiveNumber(value, part, input string) (int, error) {
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return 0, fmt.Errorf("invalid reference %q: %s must be a positive number", input, part)
	}
	return number, nil
}

func normalizeBook(book string) string {
	if canonical, ok := canon.Resolve(book); ok {
		return canonical
	}
	return book
}
