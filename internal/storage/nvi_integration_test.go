package storage

import (
	"context"
	"testing"

	"github.com/vmrocha/bible-cli/internal/reference"
)

func TestLocalNVIIntegration(t *testing.T) {
	if _, ok := embeddedTranslations["ptnvi"]; !ok {
		t.Skip("local NVI database is not embedded")
	}

	ctx := context.Background()
	reader, err := OpenEmbedded(ctx, "ptnvi")
	if err != nil {
		t.Fatalf("OpenEmbedded(ptnvi): %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	query, err := reference.Parse("João 3:16")
	if err != nil {
		t.Fatalf("parse Portuguese reference: %v", err)
	}
	passage, err := reader.Read(ctx, query)
	if err != nil {
		t.Fatalf("read NVI passage: %v", err)
	}
	if passage.Translation != "NVI" || passage.Book.Name != "João" || len(passage.Verses) != 1 || passage.Verses[0].Text == "" {
		t.Fatalf("unexpected NVI passage metadata: %#v", passage)
	}

	results, err := reader.Search(ctx, "amor", 2)
	if err != nil {
		t.Fatalf("search NVI: %v", err)
	}
	if len(results) != 2 || results[0].Translation != "NVI" {
		t.Fatalf("unexpected NVI search results: %#v", results)
	}
}
