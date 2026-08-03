package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Random selects one verse uniformly from the bundled WEBP translation.
func (reader *Reader) Random(ctx context.Context, source io.Reader) (bible.Passage, error) {
	if source == nil {
		return bible.Passage{}, errors.New("random source is required")
	}

	var count int64
	if err := reader.connection.QueryRowContext(ctx, `
        SELECT count(*)
        FROM verses
        WHERE translation_id = 'engwebp'
    `).Scan(&count); err != nil {
		return bible.Passage{}, fmt.Errorf("count verses: %w", err)
	}
	if count == 0 {
		return bible.Passage{}, errors.New("random verse unavailable: translation is empty")
	}

	selection, err := rand.Int(source, big.NewInt(count))
	if err != nil {
		return bible.Passage{}, fmt.Errorf("choose random verse: %w", err)
	}

	var passage bible.Passage
	var verse bible.Verse
	err = reader.connection.QueryRowContext(ctx, `
        SELECT
            b.id,
            b.source_code,
            b.position,
            b.name,
            t.abbreviation,
            v.chapter,
            v.verse,
            v.text
        FROM verses AS v
        JOIN books AS b
          ON b.translation_id = v.translation_id
         AND b.id = v.book_id
        JOIN translations AS t ON t.id = v.translation_id
        WHERE v.translation_id = 'engwebp'
        ORDER BY b.position, v.chapter, v.verse
        LIMIT 1 OFFSET ?
    `, selection.Int64()).Scan(
		&passage.Book.ID,
		&passage.Book.SourceCode,
		&passage.Book.Position,
		&passage.Book.Name,
		&passage.Translation,
		&verse.Chapter,
		&verse.Number,
		&verse.Text,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return bible.Passage{}, errors.New("random verse unavailable")
	}
	if err != nil {
		return bible.Passage{}, fmt.Errorf("read random verse: %w", err)
	}

	verse.BookID = passage.Book.ID
	passage.Chapter = verse.Chapter
	passage.StartVerse = verse.Number
	passage.EndVerse = verse.Number
	passage.Verses = []bible.Verse{verse}
	return passage, nil
}
