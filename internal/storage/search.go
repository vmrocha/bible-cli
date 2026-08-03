package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vmrocha/bible-cli/internal/bible"
)

const (
	highlightStart = "\ue000"
	highlightEnd   = "\ue001"
)

// Search finds verses containing every searchable token, ordered by relevance
// with canonical Scripture order as a deterministic tie-breaker.
func (reader *Reader) Search(ctx context.Context, query string, limit int) ([]bible.SearchResult, error) {
	if limit <= 0 {
		return nil, errors.New("search limit must be positive")
	}
	match, err := searchExpression(query)
	if err != nil {
		return nil, err
	}

	rows, err := reader.connection.QueryContext(ctx, `
        SELECT
            b.id,
            b.source_code,
            b.position,
            b.name,
            t.abbreviation,
            v.chapter,
            v.verse,
            v.text,
            highlight(verses_fts, 0, ?, ?) AS highlighted_text,
            bm25(verses_fts) AS relevance
        FROM verses_fts AS f
        JOIN verses AS v ON v.id = f.rowid
        JOIN books AS b
          ON b.translation_id = v.translation_id
         AND b.id = v.book_id
        JOIN translations AS t ON t.id = v.translation_id
        WHERE verses_fts MATCH ?
        ORDER BY relevance, b.position, v.chapter, v.verse
        LIMIT ?
    `, highlightStart, highlightEnd, match, limit)
	if err != nil {
		return nil, fmt.Errorf("search verses: %w", err)
	}
	defer rows.Close()

	results := make([]bible.SearchResult, 0, limit)
	for rows.Next() {
		var result bible.SearchResult
		var highlightedText string
		var relevance float64
		if err := rows.Scan(
			&result.Book.ID,
			&result.Book.SourceCode,
			&result.Book.Position,
			&result.Book.Name,
			&result.Translation,
			&result.Chapter,
			&result.Verse,
			&result.Text,
			&highlightedText,
			&relevance,
		); err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		highlights, err := parseHighlights(result.Text, highlightedText)
		if err != nil {
			return nil, fmt.Errorf(
				"parse search highlights for %s %d:%d: %w",
				result.Book.Name,
				result.Chapter,
				result.Verse,
				err,
			)
		}
		result.Highlights = highlights
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read search result rows: %w", err)
	}
	return results, nil
}

func parseHighlights(text, highlighted string) ([]bible.TextRange, error) {
	var plain strings.Builder
	plain.Grow(len(highlighted))
	highlights := make([]bible.TextRange, 0)
	start := -1

	for len(highlighted) > 0 {
		switch {
		case strings.HasPrefix(highlighted, highlightStart):
			if start >= 0 {
				return nil, errors.New("nested highlight marker")
			}
			start = plain.Len()
			highlighted = highlighted[len(highlightStart):]
		case strings.HasPrefix(highlighted, highlightEnd):
			if start < 0 {
				return nil, errors.New("closing highlight marker without an opening marker")
			}
			if start == plain.Len() {
				return nil, errors.New("empty highlighted range")
			}
			highlights = append(highlights, bible.TextRange{Start: start, End: plain.Len()})
			start = -1
			highlighted = highlighted[len(highlightEnd):]
		default:
			_, size := utf8.DecodeRuneInString(highlighted)
			plain.WriteString(highlighted[:size])
			highlighted = highlighted[size:]
		}
	}
	if start >= 0 {
		return nil, errors.New("unclosed highlight marker")
	}
	if plain.String() != text {
		return nil, errors.New("highlighted text differs from publisher text")
	}
	return highlights, nil
}

func searchExpression(query string) (string, error) {
	tokens := strings.FieldsFunc(query, func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsNumber(character)
	})
	if len(tokens) == 0 {
		return "", errors.New("search query must contain a letter or number")
	}

	phrases := make([]string, len(tokens))
	for index, token := range tokens {
		phrases[index] = fmt.Sprintf("%q", token)
	}
	return strings.Join(phrases, " AND "), nil
}
