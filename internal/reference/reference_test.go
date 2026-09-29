package reference

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  Query
	}{
		{"John 3", Query{Book: "john", Chapter: 3}},
		{"John 3:16", Query{Book: "john", Chapter: 3, StartVerse: 16, EndVerse: 16}},
		{"John 3:16-21", Query{Book: "john", Chapter: 3, StartVerse: 16, EndVerse: 21}},
		{"1 John 1:1", Query{Book: "1-john", Chapter: 1, StartVerse: 1, EndVerse: 1}},
		{"Jn. 3:16", Query{Book: "john", Chapter: 3, StartVerse: 16, EndVerse: 16}},
		{"  Psalm   23  ", Query{Book: "psalms", Chapter: 23}},
		{"1 Cor. 13", Query{Book: "1-corinthians", Chapter: 13}},
		{"João 3:16", Query{Book: "john", Chapter: 3, StartVerse: 16, EndVerse: 16}},
		{"1 Coríntios 13", Query{Book: "1-corinthians", Chapter: 13}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := Parse(test.input)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if got != test.want {
				t.Fatalf("Parse(%q) = %#v, want %#v", test.input, got, test.want)
			}
		})
	}
}

func TestParseRejectsInvalidReferences(t *testing.T) {
	inputs := []string{
		"John",
		"John zero",
		"John 0",
		"John 3:",
		"John 3:0",
		"John 3:20-16",
		"John 3:16-",
		"John 3:16-18-20",
		"John 3:16:17",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if _, err := Parse(input); err == nil {
				t.Fatalf("Parse(%q) unexpectedly succeeded", input)
			}
		})
	}
}
