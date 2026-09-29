package canon

import "testing"

func TestProtestantBooks(t *testing.T) {
	books := ProtestantBooks()
	if len(books) != 66 {
		t.Fatalf("got %d books, want 66", len(books))
	}

	for index, entry := range books {
		wantPosition := index + 1
		if entry.Book.Position != wantPosition {
			t.Fatalf("%s has position %d, want %d", entry.Book.Name, entry.Book.Position, wantPosition)
		}
		values := append([]string{entry.Book.ID, entry.Book.Name, entry.Book.SourceCode}, entry.Aliases...)
		for _, value := range values {
			got, ok := Resolve(value)
			if !ok || got != entry.Book.ID {
				t.Errorf("Resolve(%q) = %q, %t; want %q, true", value, got, ok, entry.Book.ID)
			}
		}
	}
}

func TestResolveNormalizesPunctuationAndSpacing(t *testing.T) {
	for _, value := range []string{"1 Jn", "1Jn", "1.Jn.", "1-john", "  1   JOHN  "} {
		got, ok := Resolve(value)
		if !ok || got != "1-john" {
			t.Errorf("Resolve(%q) = %q, %t; want 1-john, true", value, got, ok)
		}
	}
}

func TestPortugueseProtestantBooks(t *testing.T) {
	books := PortugueseProtestantBooks()
	if len(books) != 66 {
		t.Fatalf("got %d books, want 66", len(books))
	}
	if books[42].Book.Name != "João" || books[65].Book.Name != "Apocalipse" {
		t.Fatalf("unexpected Portuguese names: %q and %q", books[42].Book.Name, books[65].Book.Name)
	}
	for _, value := range []string{"João", "Joao"} {
		if got, ok := Resolve(value); !ok || got != "john" {
			t.Errorf("Resolve(%q) = %q, %t; want john, true", value, got, ok)
		}
	}
}
