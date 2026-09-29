package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/reference"
)

//go:embed engwebp.db
var embeddedWEBP []byte

var embeddedTranslations = map[string][]byte{
	"engwebp": embeddedWEBP,
}

type deserializer interface {
	Deserialize([]byte) error
}

// Reader provides read-only access to the database embedded in the binary.
type Reader struct {
	database      *sql.DB
	connection    *sql.Conn
	translationID string
}

// RegisterEmbedded adds a locally embedded translation during package
// initialization. It is intended for ignored, license-specific source files.
func RegisterEmbedded(id string, database []byte) {
	if id == "" || len(database) == 0 {
		panic("embedded translation requires an id and database")
	}
	if _, exists := embeddedTranslations[id]; exists {
		panic("embedded translation is already registered: " + id)
	}
	embeddedTranslations[id] = database
}

// OpenEmbedded loads the embedded database into a dedicated in-memory
// connection. It does not create files or require network access.
func OpenEmbedded(ctx context.Context, requested ...string) (*Reader, error) {
	if len(requested) > 1 {
		return nil, errors.New("open embedded database: at most one translation may be selected")
	}
	translationID := "engwebp"
	if len(requested) == 1 && requested[0] != "" {
		translationID = requested[0]
	}
	embeddedDatabase, ok := embeddedTranslations[translationID]
	if !ok {
		return nil, fmt.Errorf("translation is not embedded: %s", translationID)
	}

	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open embedded database: %w", err)
	}
	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)

	connection, err := database.Conn(ctx)
	if err != nil {
		database.Close()
		return nil, fmt.Errorf("connect to embedded database: %w", err)
	}

	if err := connection.Raw(func(driverConnection any) error {
		loader, ok := driverConnection.(deserializer)
		if !ok {
			return errors.New("SQLite driver does not support database deserialization")
		}
		return loader.Deserialize(embeddedDatabase)
	}); err != nil {
		connection.Close()
		database.Close()
		return nil, fmt.Errorf("load embedded database: %w", err)
	}
	if _, err := connection.ExecContext(ctx, "PRAGMA query_only = ON"); err != nil {
		connection.Close()
		database.Close()
		return nil, fmt.Errorf("protect embedded database: %w", err)
	}

	var schemaVersion int
	if err := connection.QueryRowContext(ctx, "PRAGMA user_version").Scan(&schemaVersion); err != nil {
		connection.Close()
		database.Close()
		return nil, fmt.Errorf("read embedded schema version: %w", err)
	}
	if schemaVersion != 2 {
		connection.Close()
		database.Close()
		return nil, fmt.Errorf("unsupported embedded schema version %d", schemaVersion)
	}
	var translationCount int
	if err := connection.QueryRowContext(
		ctx,
		"SELECT count(*) FROM translations WHERE id = ?",
		translationID,
	).Scan(&translationCount); err != nil {
		connection.Close()
		database.Close()
		return nil, fmt.Errorf("validate embedded translation: %w", err)
	}
	if translationCount != 1 {
		connection.Close()
		database.Close()
		return nil, fmt.Errorf("embedded database does not contain translation %s", translationID)
	}

	return &Reader{database: database, connection: connection, translationID: translationID}, nil
}

// Close releases the in-memory database.
func (reader *Reader) Close() error {
	return errors.Join(reader.connection.Close(), reader.database.Close())
}

// Read resolves and reads one chapter, verse, or inclusive verse range.
func (reader *Reader) Read(ctx context.Context, query reference.Query) (bible.Passage, error) {
	book, translation, err := reader.resolveBook(ctx, query.Book)
	if err != nil {
		return bible.Passage{}, err
	}

	rows, err := reader.connection.QueryContext(ctx, `
        SELECT chapter, verse, text
        FROM verses
        WHERE translation_id = ?
          AND book_id = ?
          AND chapter = ?
          AND (? = 0 OR verse BETWEEN ? AND ?)
        ORDER BY verse
    `, reader.translationID, book.ID, query.Chapter, query.StartVerse, query.StartVerse, query.EndVerse)
	if err != nil {
		return bible.Passage{}, fmt.Errorf("read passage: %w", err)
	}
	defer rows.Close()

	passage := bible.Passage{
		Translation: translation,
		Book:        book,
		Chapter:     query.Chapter,
		StartVerse:  query.StartVerse,
		EndVerse:    query.EndVerse,
	}
	for rows.Next() {
		var verse bible.Verse
		verse.BookID = book.ID
		if err := rows.Scan(&verse.Chapter, &verse.Number, &verse.Text); err != nil {
			return bible.Passage{}, fmt.Errorf("scan passage: %w", err)
		}
		passage.Verses = append(passage.Verses, verse)
	}
	if err := rows.Err(); err != nil {
		return bible.Passage{}, fmt.Errorf("read passage rows: %w", err)
	}

	if len(passage.Verses) == 0 {
		return bible.Passage{}, passageNotFound(book.Name, query)
	}
	if !query.IsChapter() && len(passage.Verses) != query.EndVerse-query.StartVerse+1 {
		return bible.Passage{}, passageNotFound(book.Name, query)
	}

	return passage, nil
}

// Navigate resolves the chapter immediately before or after a chapter query,
// crossing book boundaries when necessary.
func (reader *Reader) Navigate(ctx context.Context, query reference.Query, direction int) (reference.Query, error) {
	if !query.IsChapter() {
		return reference.Query{}, errors.New("chapter navigation requires a chapter reference")
	}
	if direction != -1 && direction != 1 {
		return reference.Query{}, fmt.Errorf("invalid navigation direction %d", direction)
	}

	book, _, err := reader.resolveBook(ctx, query.Book)
	if err != nil {
		return reference.Query{}, err
	}
	lastChapter, err := reader.lastChapter(ctx, book.ID)
	if err != nil {
		return reference.Query{}, err
	}
	if query.Chapter > lastChapter {
		return reference.Query{}, passageNotFound(book.Name, query)
	}

	if direction > 0 {
		if query.Chapter < lastChapter {
			return reference.Query{Book: book.ID, Chapter: query.Chapter + 1}, nil
		}
		next, err := reader.bookAtPosition(ctx, book.Position+1)
		if errors.Is(err, sql.ErrNoRows) {
			return reference.Query{}, errors.New("already at the end of the Bible")
		}
		if err != nil {
			return reference.Query{}, err
		}
		return reference.Query{Book: next.ID, Chapter: 1}, nil
	}

	if query.Chapter > 1 {
		return reference.Query{Book: book.ID, Chapter: query.Chapter - 1}, nil
	}
	previous, err := reader.bookAtPosition(ctx, book.Position-1)
	if errors.Is(err, sql.ErrNoRows) {
		return reference.Query{}, errors.New("already at the beginning of the Bible")
	}
	if err != nil {
		return reference.Query{}, err
	}
	lastChapter, err = reader.lastChapter(ctx, previous.ID)
	if err != nil {
		return reference.Query{}, err
	}
	return reference.Query{Book: previous.ID, Chapter: lastChapter}, nil
}

func (reader *Reader) resolveBook(ctx context.Context, name string) (bible.Book, string, error) {
	candidateID := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	var book bible.Book
	var translation string
	err := reader.connection.QueryRowContext(ctx, `
        SELECT b.id, b.source_code, b.position, b.name, t.abbreviation
        FROM books AS b
        JOIN translations AS t ON t.id = b.translation_id
        WHERE b.translation_id = ?
          AND (lower(b.name) = lower(?) OR lower(b.source_code) = lower(?) OR b.id = ?)
    `, reader.translationID, name, name, candidateID).Scan(
		&book.ID,
		&book.SourceCode,
		&book.Position,
		&book.Name,
		&translation,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return bible.Book{}, "", fmt.Errorf("book not found: %s", name)
	}
	if err != nil {
		return bible.Book{}, "", fmt.Errorf("resolve book: %w", err)
	}
	book.Name = localizedBookName(reader.translationID, book.ID, book.Name)
	return book, translation, nil
}

func (reader *Reader) bookAtPosition(ctx context.Context, position int) (bible.Book, error) {
	var book bible.Book
	err := reader.connection.QueryRowContext(ctx, `
        SELECT id, source_code, position, name
        FROM books
        WHERE translation_id = ? AND position = ?
    `, reader.translationID, position).Scan(&book.ID, &book.SourceCode, &book.Position, &book.Name)
	if err != nil {
		return bible.Book{}, err
	}
	book.Name = localizedBookName(reader.translationID, book.ID, book.Name)
	return book, nil
}

func (reader *Reader) lastChapter(ctx context.Context, bookID string) (int, error) {
	var chapter int
	err := reader.connection.QueryRowContext(ctx, `
        SELECT max(chapter)
        FROM verses
        WHERE translation_id = ? AND book_id = ?
    `, reader.translationID, bookID).Scan(&chapter)
	if err != nil {
		return 0, fmt.Errorf("read last chapter for %s: %w", bookID, err)
	}
	return chapter, nil
}

func embeddedTranslationIDs() []string {
	ids := make([]string, 0, len(embeddedTranslations))
	for id := range embeddedTranslations {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func passageNotFound(book string, query reference.Query) error {
	if query.IsChapter() {
		return fmt.Errorf("chapter not found: %s %d", book, query.Chapter)
	}
	if query.StartVerse == query.EndVerse {
		return fmt.Errorf("verse not found: %s %d:%d", book, query.Chapter, query.StartVerse)
	}
	return fmt.Errorf(
		"verse range not found: %s %d:%d-%d",
		book,
		query.Chapter,
		query.StartVerse,
		query.EndVerse,
	)
}
