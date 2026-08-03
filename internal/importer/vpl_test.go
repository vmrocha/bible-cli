package importer

import (
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/translation"
)

func TestParseVPL(t *testing.T) {
	input := strings.Join([]string{
		"GEN 1:1 In the beginning…",
		"GEN 1:2 Verse text with  two spaces.",
		"GEN 1:3 ",
		"GEN 2:1 Next chapter.",
		"EXO 1:1 Second book.",
	}, "\n")
	expected := translation.Expected{
		Books:          2,
		Verses:         5,
		FirstReference: "GEN 1:1",
		LastReference:  "EXO 1:1",
	}

	dataset, err := ParseVPL(strings.NewReader(input), expected)
	if err != nil {
		t.Fatalf("ParseVPL: %v", err)
	}

	if len(dataset.Books) != 2 || dataset.Books[1].ID != "exodus" {
		t.Fatalf("unexpected books: %#v", dataset.Books)
	}
	if dataset.Verses[0].Text != "In the beginning…" {
		t.Fatalf("verse text changed: %q", dataset.Verses[0].Text)
	}
	if dataset.Verses[1].Text != "Verse text with  two spaces." {
		t.Fatalf("verse whitespace changed: %q", dataset.Verses[1].Text)
	}
	if dataset.Verses[2].Text != "" {
		t.Fatalf("empty verse gained text: %q", dataset.Verses[2].Text)
	}
}

func TestParseVPLRejectsInvalidSources(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "duplicate reference",
			input: "GEN 1:1 First.\nGEN 1:1 Again.",
			want:  "duplicate reference",
		},
		{
			name:  "missing verse",
			input: "GEN 1:1 First.\nGEN 1:3 Third.",
			want:  "does not follow",
		},
		{
			name:  "book out of order",
			input: "GEN 1:1 First.\nLEV 1:1 Third book.",
			want:  "out of canonical order",
		},
		{
			name:  "unknown book",
			input: "XYZ 1:1 Unknown.",
			want:  "unknown book code",
		},
		{
			name:  "verse range",
			input: "GEN 1:1-2 Joined.",
			want:  "invalid verse",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected := translation.Expected{
				Books:          1,
				Verses:         2,
				FirstReference: "GEN 1:1",
				LastReference:  "GEN 1:2",
			}
			_, err := ParseVPL(strings.NewReader(test.input), expected)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}
