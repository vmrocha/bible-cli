package storage

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/translation"
)

func TestGenerate(t *testing.T) {
	manifest, dataset := databaseFixture()
	path := filepath.Join(t.TempDir(), "nested", "translation.db")

	if err := Generate(context.Background(), path, manifest, dataset); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open generated database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	var text string
	err = database.QueryRow(`
        SELECT text FROM verses
        WHERE translation_id = 'example' AND book_id = 'genesis'
          AND chapter = 1 AND verse = 2
    `).Scan(&text)
	if err != nil {
		t.Fatalf("query generated verse: %v", err)
	}
	if text != "Second." {
		t.Fatalf("unexpected verse text: %q", text)
	}

	var schemaVersion int
	if err := database.QueryRow("PRAGMA user_version").Scan(&schemaVersion); err != nil {
		t.Fatalf("query schema version: %v", err)
	}
	if schemaVersion != 2 {
		t.Fatalf("schema version is %d, want 2", schemaVersion)
	}

	var indexed int
	if err := database.QueryRow("SELECT count(*) FROM verses_fts WHERE verses_fts MATCH 'Second'").Scan(&indexed); err != nil {
		t.Fatalf("query generated search index: %v", err)
	}
	if indexed != 1 {
		t.Fatalf("search index returned %d verses, want 1", indexed)
	}
}

func TestGenerateIsDeterministic(t *testing.T) {
	manifest, dataset := databaseFixture()
	directory := t.TempDir()
	firstPath := filepath.Join(directory, "first.db")
	secondPath := filepath.Join(directory, "second.db")

	if err := Generate(context.Background(), firstPath, manifest, dataset); err != nil {
		t.Fatalf("generate first database: %v", err)
	}
	if err := Generate(context.Background(), secondPath, manifest, dataset); err != nil {
		t.Fatalf("generate second database: %v", err)
	}
	first, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatalf("read first database: %v", err)
	}
	second, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatalf("read second database: %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("generated databases differ")
	}
}

func TestGenerateRejectsExistingOutput(t *testing.T) {
	manifest, dataset := databaseFixture()
	path := filepath.Join(t.TempDir(), "translation.db")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatalf("write existing output: %v", err)
	}

	err := Generate(context.Background(), path, manifest, dataset)
	if !errors.Is(err, ErrOutputExists) {
		t.Fatalf("expected ErrOutputExists, got %v", err)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read existing output: %v", err)
	}
	if string(contents) != "keep me" {
		t.Fatalf("existing output changed: %q", contents)
	}
}

func databaseFixture() (translation.Manifest, bible.Dataset) {
	manifest := translation.Manifest{
		SchemaVersion: 1,
		ID:            "example",
		Name:          "Example Bible",
		Abbreviation:  "EX",
		Language: translation.Language{
			Tag:  "en",
			Name: "English",
		},
		Edition:     "Test Edition",
		Canon:       "test",
		TextEdition: "1",
		Source: translation.Source{
			Publisher:     "Example",
			Homepage:      "https://example.com",
			ArchiveURL:    "https://example.com/source.zip",
			ArchiveMember: "source.txt",
			RetrievedAt:   "2026-07-15",
			ArchiveSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		Rights: translation.Rights{
			Status:     "public-domain",
			NoticeURL:  "https://example.com/rights",
			Trademark:  "Example naming notice.",
			TextPolicy: "Preserve source text.",
		},
		Expected: translation.Expected{
			Books:          1,
			Verses:         2,
			FirstReference: "GEN 1:1",
			LastReference:  "GEN 1:2",
		},
	}
	dataset := bible.Dataset{
		Books: []bible.Book{
			{ID: "genesis", SourceCode: "GEN", Position: 1, Name: "Genesis"},
		},
		Verses: []bible.Verse{
			{BookID: "genesis", Chapter: 1, Number: 1, Text: "First."},
			{BookID: "genesis", Chapter: 1, Number: 2, Text: "Second."},
		},
	}
	return manifest, dataset
}
