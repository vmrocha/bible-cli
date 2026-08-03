package importer

import (
	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/canon"
)

var protestantBooks = func() []bible.Book {
	entries := canon.ProtestantBooks()
	books := make([]bible.Book, len(entries))
	for index, entry := range entries {
		books[index] = entry.Book
	}
	return books
}()

var booksBySourceCode = func() map[string]bible.Book {
	books := make(map[string]bible.Book, len(protestantBooks))
	for _, book := range protestantBooks {
		books[book.SourceCode] = book
	}
	return books
}()
