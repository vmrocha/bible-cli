package storage

import "github.com/vmrocha/bible-cli/internal/canon"

func localizedBookName(translationID, bookID, fallback string) string {
	if translationID != "ptnvi" {
		return fallback
	}
	if name, ok := canon.PortugueseName(bookID); ok {
		return name
	}
	return fallback
}
