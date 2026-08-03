package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/canon"
)

func TestBooks(t *testing.T) {
	var output bytes.Buffer
	if err := Books(&output, canon.ProtestantBooks(), Options{}); err != nil {
		t.Fatalf("Books: %v", err)
	}

	for _, expected := range []string{
		"Old Testament",
		" 1  Genesis",
		"39  Malachi",
		"New Testament",
		"40  Matthew",
		"66  Revelation",
		"Jn, Jhn",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("book output does not contain %q", expected)
		}
	}
}

func TestPlainBooks(t *testing.T) {
	var output bytes.Buffer
	if err := Books(&output, canon.ProtestantBooks(), Options{Plain: true}); err != nil {
		t.Fatalf("Books: %v", err)
	}
	if !strings.Contains(output.String(), "john\tJohn\tJOH\tJn,Jhn\n") {
		t.Fatalf("unexpected plain output:\n%s", output.String())
	}
}

func TestStyledBooks(t *testing.T) {
	var output bytes.Buffer
	if err := Books(&output, canon.ProtestantBooks()[:1], Options{Color: true}); err != nil {
		t.Fatalf("Books: %v", err)
	}
	want := "\x1b[1m\x1b[36mOld Testament\x1b[0m\n" +
		"\x1b[2m 1\x1b[0m  Genesis           GEN  \x1b[2mGen, Ge, Gn\x1b[0m\n"
	if output.String() != want {
		t.Fatalf("unexpected styled output:\n%q", output.String())
	}
}
