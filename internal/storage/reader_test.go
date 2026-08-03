package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/canon"
	"github.com/vmrocha/bible-cli/internal/reference"
)

func TestEmbeddedDatabase(t *testing.T) {
	digest := sha256.Sum256(embeddedWEBP)
	got := hex.EncodeToString(digest[:])
	const want = "5fb07654e06c6e6999d57b101b3d37d04e0b32d124ee6c1eac252cb542a354ea"
	if got != want {
		t.Fatalf("embedded database checksum is %s, want %s", got, want)
	}

	reader, err := OpenEmbedded(context.Background())
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	var books, verses, indexedVerses int
	if err := reader.connection.QueryRowContext(
		context.Background(),
		"SELECT (SELECT count(*) FROM books), (SELECT count(*) FROM verses), (SELECT count(*) FROM verses_fts)",
	).Scan(&books, &verses, &indexedVerses); err != nil {
		t.Fatalf("query embedded totals: %v", err)
	}
	if books != 66 || verses != 31103 || indexedVerses != verses {
		t.Fatalf("embedded totals are %d books, %d verses, and %d indexed verses", books, verses, indexedVerses)
	}

	if _, err := reader.connection.ExecContext(context.Background(), "DELETE FROM verses"); err == nil {
		t.Fatal("embedded database unexpectedly allowed a write")
	}
}

func TestRead(t *testing.T) {
	reader, err := OpenEmbedded(context.Background())
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	query, err := reference.Parse("Jn 3:16")
	if err != nil {
		t.Fatalf("parse reference: %v", err)
	}
	passage, err := reader.Read(context.Background(), query)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if passage.Book.Name != "John" || passage.Translation != "WEBP" {
		t.Fatalf("unexpected passage metadata: %#v", passage)
	}
	if len(passage.Verses) != 1 || !strings.HasPrefix(passage.Verses[0].Text, "For God so loved") {
		t.Fatalf("unexpected passage verses: %#v", passage.Verses)
	}
}

func TestReadRejectsMissingPassages(t *testing.T) {
	reader, err := OpenEmbedded(context.Background())
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	for _, input := range []string{"Notabook 1", "John 999", "John 3:999", "John 3:35-40"} {
		t.Run(input, func(t *testing.T) {
			query, err := reference.Parse(input)
			if err != nil {
				t.Fatalf("parse reference: %v", err)
			}
			if _, err := reader.Read(context.Background(), query); err == nil {
				t.Fatalf("Read(%q) unexpectedly succeeded", input)
			}
		})
	}
}

func TestNavigate(t *testing.T) {
	reader, err := OpenEmbedded(context.Background())
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	tests := []struct {
		input     string
		direction int
		want      reference.Query
	}{
		{"John 3", 1, reference.Query{Book: "john", Chapter: 4}},
		{"John 3", -1, reference.Query{Book: "john", Chapter: 2}},
		{"John 21", 1, reference.Query{Book: "acts", Chapter: 1}},
		{"Matthew 1", -1, reference.Query{Book: "malachi", Chapter: 4}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			query, err := reference.Parse(test.input)
			if err != nil {
				t.Fatalf("parse reference: %v", err)
			}
			got, err := reader.Navigate(context.Background(), query, test.direction)
			if err != nil {
				t.Fatalf("Navigate: %v", err)
			}
			if got != test.want {
				t.Fatalf("Navigate = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestNavigateRejectsBoundaries(t *testing.T) {
	reader, err := OpenEmbedded(context.Background())
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	for _, test := range []struct {
		input     string
		direction int
		want      string
	}{
		{"Genesis 1", -1, "beginning"},
		{"Revelation 22", 1, "end"},
		{"John 999", 1, "chapter not found"},
		{"John 3:16", 1, "requires a chapter"},
	} {
		t.Run(test.input, func(t *testing.T) {
			query, err := reference.Parse(test.input)
			if err != nil {
				t.Fatalf("parse reference: %v", err)
			}
			_, err = reader.Navigate(context.Background(), query, test.direction)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestNavigateAcrossEveryBookBoundary(t *testing.T) {
	ctx := context.Background()
	reader, err := OpenEmbedded(ctx)
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	books := canon.ProtestantBooks()
	for index := 0; index < len(books)-1; index++ {
		current := books[index].Book
		next := books[index+1].Book
		lastChapter, err := reader.lastChapter(ctx, current.ID)
		if err != nil {
			t.Fatalf("last chapter for %s: %v", current.Name, err)
		}

		forward, err := reader.Navigate(ctx, reference.Query{Book: current.ID, Chapter: lastChapter}, 1)
		if err != nil {
			t.Fatalf("navigate forward from %s: %v", current.Name, err)
		}
		if forward.Book != next.ID || forward.Chapter != 1 {
			t.Errorf("after %s is %#v, want %s 1", current.Name, forward, next.Name)
		}

		backward, err := reader.Navigate(ctx, reference.Query{Book: next.ID, Chapter: 1}, -1)
		if err != nil {
			t.Fatalf("navigate backward from %s: %v", next.Name, err)
		}
		if backward.Book != current.ID || backward.Chapter != lastChapter {
			t.Errorf("before %s is %#v, want %s %d", next.Name, backward, current.Name, lastChapter)
		}
	}
}
