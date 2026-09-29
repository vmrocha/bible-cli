package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
)

type stubSearcher struct {
	results []bible.SearchResult
	err     error
	closed  bool
	query   string
	limit   int
}

func (searcher *stubSearcher) Search(_ context.Context, query string, limit int) ([]bible.SearchResult, error) {
	searcher.query = query
	searcher.limit = limit
	return searcher.results, searcher.err
}

func (searcher *stubSearcher) Close() error {
	searcher.closed = true
	return nil
}

func TestSearchCommand(t *testing.T) {
	searcher := &stubSearcher{results: []bible.SearchResult{{
		Translation: "WEBP",
		Book:        bible.Book{Name: "John"},
		Chapter:     4,
		Verse:       10,
		Text:        "He would have given you living water.",
	}}}
	factory := func(context.Context, string) (Searcher, error) { return searcher, nil }

	output, err := executeWithOptions(
		t,
		[]Option{WithSearcherFactory(factory)},
		"search", "living", "water", "--limit", "7", "--no-color",
	)
	if err != nil {
		t.Fatalf("execute search: %v", err)
	}
	if searcher.query != "living water" || searcher.limit != 7 {
		t.Fatalf("search request was query %q with limit %d", searcher.query, searcher.limit)
	}
	if !strings.Contains(output, "John 4:10") || !strings.Contains(output, "living water") {
		t.Fatalf("unexpected search output:\n%s", output)
	}
	if !searcher.closed {
		t.Fatal("searcher was not closed")
	}
}

func TestSearchCommandAutomaticallyUsesPlainOutputWhenRedirected(t *testing.T) {
	searcher := &stubSearcher{results: []bible.SearchResult{{
		Book:    bible.Book{Name: "John"},
		Chapter: 4,
		Verse:   10,
		Text:    "Living water.",
		Highlights: []bible.TextRange{
			{Start: 0, End: 6},
			{Start: 7, End: 12},
		},
	}}}
	factory := func(context.Context, string) (Searcher, error) { return searcher, nil }
	output := new(bytes.Buffer)
	command := New(testBuild, WithSearcherFactory(factory))
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"search", "living water"})

	if err := command.Execute(); err != nil {
		t.Fatalf("execute redirected search: %v", err)
	}
	if want := "John 4:10\tLiving water.\n"; output.String() != want {
		t.Fatalf("unexpected redirected output: %q", output.String())
	}
	if strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("redirected output contains ANSI escapes: %q", output.String())
	}
}

func TestSearchCommandHighlightsInteractiveMatches(t *testing.T) {
	searcher := &stubSearcher{results: []bible.SearchResult{{
		Translation: "WEBP",
		Book:        bible.Book{Name: "John"},
		Chapter:     4,
		Verse:       10,
		Text:        "Living water.",
		Highlights: []bible.TextRange{
			{Start: 0, End: 6},
			{Start: 7, End: 12},
		},
	}}}
	factory := func(context.Context, string) (Searcher, error) { return searcher, nil }

	output, err := executeWithOptions(t, []Option{WithSearcherFactory(factory)}, "search", "living water")
	if err != nil {
		t.Fatalf("execute search: %v", err)
	}
	for _, expected := range []string{
		"\x1b[1m\x1b[33mLiving\x1b[0m",
		"\x1b[1m\x1b[33mwater\x1b[0m",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("interactive output does not contain %q: %q", expected, output)
		}
	}
}

func TestSearchCommandRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{
		{"search", "   "},
		{"search", "water", "--limit", "0"},
		{"search", "water", "--limit", "101"},
	} {
		if _, err := execute(t, args...); err == nil {
			t.Fatalf("execute %v unexpectedly succeeded", args)
		}
	}
}

func TestSearchCommandReportsSearchErrorAndCloses(t *testing.T) {
	searcher := &stubSearcher{err: errors.New("search failed")}
	factory := func(context.Context, string) (Searcher, error) { return searcher, nil }

	_, err := executeWithOptions(t, []Option{WithSearcherFactory(factory)}, "search", "water")
	if err == nil || !strings.Contains(err.Error(), "search failed") {
		t.Fatalf("expected search error, got %v", err)
	}
	if !searcher.closed {
		t.Fatal("searcher was not closed after an error")
	}
}
