package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
)

type stubRandomReader struct {
	passage bible.Passage
	err     error
	closed  bool
	source  io.Reader
}

func (reader *stubRandomReader) Random(_ context.Context, source io.Reader) (bible.Passage, error) {
	reader.source = source
	return reader.passage, reader.err
}

func (reader *stubRandomReader) Close() error {
	reader.closed = true
	return nil
}

func randomPassageFixture() bible.Passage {
	return bible.Passage{
		Translation: "WEBP",
		Book:        bible.Book{ID: "psalms", Name: "Psalms"},
		Chapter:     23,
		StartVerse:  1,
		EndVerse:    1,
		Verses: []bible.Verse{
			{BookID: "psalms", Chapter: 23, Number: 1, Text: "Yahweh is my shepherd."},
		},
	}
}

func TestRandomCommand(t *testing.T) {
	reader := &stubRandomReader{passage: randomPassageFixture()}
	factory := func(context.Context) (RandomReader, error) { return reader, nil }

	output, err := executeWithOptions(
		t,
		[]Option{WithRandomReaderFactory(factory)},
		"random", "--no-color",
	)
	if err != nil {
		t.Fatalf("execute random: %v", err)
	}
	if !strings.Contains(output, "Psalms 23:1 · WEBP") || !strings.Contains(output, "Yahweh is my shepherd.") {
		t.Fatalf("unexpected random output:\n%s", output)
	}
	if reader.source == nil {
		t.Fatal("random command did not provide entropy")
	}
	if !reader.closed {
		t.Fatal("random reader was not closed")
	}
}

func TestRandomCommandAutomaticallyUsesPlainOutputWhenRedirected(t *testing.T) {
	reader := &stubRandomReader{passage: randomPassageFixture()}
	factory := func(context.Context) (RandomReader, error) { return reader, nil }
	output := new(bytes.Buffer)
	command := New(testBuild, WithRandomReaderFactory(factory))
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"random"})

	if err := command.Execute(); err != nil {
		t.Fatalf("execute redirected random: %v", err)
	}
	if want := "Psalms 23:1\tYahweh is my shepherd.\n"; output.String() != want {
		t.Fatalf("unexpected redirected output: %q", output.String())
	}
}

func TestRandomCommandReportsReaderErrorAndCloses(t *testing.T) {
	reader := &stubRandomReader{err: errors.New("random failed")}
	factory := func(context.Context) (RandomReader, error) { return reader, nil }

	_, err := executeWithOptions(t, []Option{WithRandomReaderFactory(factory)}, "random")
	if err == nil || !strings.Contains(err.Error(), "random failed") {
		t.Fatalf("expected random error, got %v", err)
	}
	if !reader.closed {
		t.Fatal("random reader was not closed after an error")
	}
}
