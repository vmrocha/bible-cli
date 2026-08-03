package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/translation"
	_ "modernc.org/sqlite"
)

// ErrOutputExists prevents an import from silently replacing reviewed data.
var ErrOutputExists = errors.New("output database already exists")

//go:embed schema.sql
var schemaSQL string

// Generate atomically creates a normalized SQLite translation database.
func Generate(ctx context.Context, outputPath string, manifest translation.Manifest, dataset bible.Dataset) error {
	if err := manifest.Validate(); err != nil {
		return err
	}
	if len(dataset.Books) != manifest.Expected.Books {
		return fmt.Errorf("generate database: got %d books, want %d", len(dataset.Books), manifest.Expected.Books)
	}
	if len(dataset.Verses) != manifest.Expected.Verses {
		return fmt.Errorf("generate database: got %d verses, want %d", len(dataset.Verses), manifest.Expected.Verses)
	}

	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("%w: %s", ErrOutputExists, outputPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect output database: %w", err)
	}

	directory := filepath.Dir(outputPath)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	temporary, err := os.CreateTemp(directory, ".bible-cli-*.db")
	if err != nil {
		return fmt.Errorf("create temporary database: %w", err)
	}
	temporaryPath := temporary.Name()
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary database: %w", err)
	}
	defer os.Remove(temporaryPath)

	if err := populate(ctx, temporaryPath, manifest, dataset); err != nil {
		return err
	}
	if err := os.Chmod(temporaryPath, 0o644); err != nil {
		return fmt.Errorf("set database permissions: %w", err)
	}
	if err := os.Link(temporaryPath, outputPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: %s", ErrOutputExists, outputPath)
		}
		return fmt.Errorf("publish output database: %w", err)
	}

	return nil
}

func populate(ctx context.Context, path string, manifest translation.Manifest, dataset bible.Dataset) error {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open temporary database: %w", err)
	}
	databaseClosed := false
	defer func() {
		if !databaseClosed {
			_ = database.Close()
		}
	}()
	database.SetMaxOpenConns(1)

	if _, err := database.ExecContext(ctx, `
        PRAGMA foreign_keys = ON;
        PRAGMA journal_mode = OFF;
        PRAGMA synchronous = OFF;
        PRAGMA page_size = 4096;
        PRAGMA auto_vacuum = NONE;
    `); err != nil {
		return fmt.Errorf("configure temporary database: %w", err)
	}

	transaction, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}
	defer transaction.Rollback()

	if _, err := transaction.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("create database schema: %w", err)
	}
	if err := insertTranslation(ctx, transaction, manifest); err != nil {
		return err
	}
	if err := insertBooks(ctx, transaction, manifest.ID, dataset.Books); err != nil {
		return err
	}
	if err := insertVerses(ctx, transaction, manifest.ID, dataset.Verses); err != nil {
		return err
	}
	if err := buildSearchIndex(ctx, transaction); err != nil {
		return err
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction: %w", err)
	}
	if _, err := database.ExecContext(ctx, "VACUUM"); err != nil {
		return fmt.Errorf("finalize database: %w", err)
	}
	if err := database.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	databaseClosed = true

	return nil
}

func insertTranslation(ctx context.Context, transaction *sql.Tx, manifest translation.Manifest) error {
	_, err := transaction.ExecContext(ctx, `
        INSERT INTO translations (
            id, name, abbreviation, language_tag, language_name, edition, canon,
            text_edition, source_publisher, source_homepage, source_archive_url,
            source_archive_sha256, source_retrieved_at, rights_status,
            rights_notice_url, trademark_notice, text_policy
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
		manifest.ID,
		manifest.Name,
		manifest.Abbreviation,
		manifest.Language.Tag,
		manifest.Language.Name,
		manifest.Edition,
		manifest.Canon,
		manifest.TextEdition,
		manifest.Source.Publisher,
		manifest.Source.Homepage,
		manifest.Source.ArchiveURL,
		manifest.Source.ArchiveSHA256,
		manifest.Source.RetrievedAt,
		manifest.Rights.Status,
		manifest.Rights.NoticeURL,
		manifest.Rights.Trademark,
		manifest.Rights.TextPolicy,
	)
	if err != nil {
		return fmt.Errorf("insert translation metadata: %w", err)
	}
	return nil
}

func insertBooks(ctx context.Context, transaction *sql.Tx, translationID string, books []bible.Book) error {
	statement, err := transaction.PrepareContext(ctx, `
        INSERT INTO books (translation_id, id, source_code, position, name)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("prepare book insert: %w", err)
	}
	defer statement.Close()

	for _, book := range books {
		if _, err := statement.ExecContext(ctx, translationID, book.ID, book.SourceCode, book.Position, book.Name); err != nil {
			return fmt.Errorf("insert book %s: %w", book.ID, err)
		}
	}
	return nil
}

func insertVerses(ctx context.Context, transaction *sql.Tx, translationID string, verses []bible.Verse) error {
	statement, err := transaction.PrepareContext(ctx, `
        INSERT INTO verses (translation_id, book_id, chapter, verse, text)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("prepare verse insert: %w", err)
	}
	defer statement.Close()

	for _, verse := range verses {
		if _, err := statement.ExecContext(
			ctx,
			translationID,
			verse.BookID,
			verse.Chapter,
			verse.Number,
			verse.Text,
		); err != nil {
			return fmt.Errorf("insert verse %s %d:%d: %w", verse.BookID, verse.Chapter, verse.Number, err)
		}
	}
	return nil
}

func buildSearchIndex(ctx context.Context, transaction *sql.Tx) error {
	if _, err := transaction.ExecContext(ctx, "INSERT INTO verses_fts(verses_fts) VALUES('rebuild')"); err != nil {
		return fmt.Errorf("build search index: %w", err)
	}
	return nil
}
