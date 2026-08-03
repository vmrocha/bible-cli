package render

import (
	"bytes"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
)

func TestPassage(t *testing.T) {
	passage := bible.Passage{
		Translation: "WEBP",
		Book:        bible.Book{Name: "John"},
		Chapter:     3,
		StartVerse:  16,
		EndVerse:    17,
		Verses: []bible.Verse{
			{Chapter: 3, Number: 16, Text: "For God so loved…"},
			{Chapter: 3, Number: 17, Text: "For God didn’t send…"},
		},
	}

	var output bytes.Buffer
	if err := Passage(&output, passage, Options{}); err != nil {
		t.Fatalf("Passage: %v", err)
	}
	want := "John 3:16–17 · WEBP\n\n16  For God so loved…\n17  For God didn’t send…\n"
	if output.String() != want {
		t.Fatalf("unexpected output:\n%s", output.String())
	}
}

func TestPlainPassage(t *testing.T) {
	passage := bible.Passage{
		Book:    bible.Book{Name: "Luke"},
		Chapter: 17,
		Verses: []bible.Verse{
			{Chapter: 17, Number: 36, Text: ""},
			{Chapter: 17, Number: 37, Text: "They asked…"},
		},
	}

	var output bytes.Buffer
	if err := Passage(&output, passage, Options{Plain: true}); err != nil {
		t.Fatalf("Passage: %v", err)
	}
	want := "Luke 17:36\t\nLuke 17:37\tThey asked…\n"
	if output.String() != want {
		t.Fatalf("unexpected output:\n%s", output.String())
	}
}

func TestPassageRejectsEmptyContent(t *testing.T) {
	var output bytes.Buffer
	if err := Passage(&output, bible.Passage{}, Options{}); err == nil {
		t.Fatal("expected empty passage to return an error")
	}
}

func TestStyledPassage(t *testing.T) {
	passage := bible.Passage{
		Translation: "WEBP",
		Book:        bible.Book{Name: "John"},
		Chapter:     3,
		StartVerse:  16,
		EndVerse:    16,
		Verses:      []bible.Verse{{Chapter: 3, Number: 16, Text: "For God so loved…"}},
	}

	var output bytes.Buffer
	if err := Passage(&output, passage, Options{Color: true}); err != nil {
		t.Fatalf("Passage: %v", err)
	}
	want := "\x1b[1m\x1b[36mJohn 3:16\x1b[0m · \x1b[2mWEBP\x1b[0m\n\n" +
		"\x1b[1m\x1b[36m16\x1b[0m  For God so loved…\n"
	if output.String() != want {
		t.Fatalf("unexpected styled output:\n%q", output.String())
	}
}
