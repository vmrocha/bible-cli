package storage

import (
	"context"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
)

func TestTranslations(t *testing.T) {
	reader, err := OpenEmbedded(context.Background())
	if err != nil {
		t.Fatalf("OpenEmbedded: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })

	translations, err := reader.Translations(context.Background())
	if err != nil {
		t.Fatalf("Translations: %v", err)
	}
	if len(translations) != len(embeddedTranslationIDs()) {
		t.Fatalf("Translations returned %d entries, want %d", len(translations), len(embeddedTranslationIDs()))
	}
	var got bible.Translation
	for _, translation := range translations {
		if translation.ID == "engwebp" {
			got = translation
			break
		}
	}
	if got.ID != "engwebp" ||
		got.Abbreviation != "WEBP" ||
		got.Name != "World English Bible" ||
		got.LanguageTag != "en-US" ||
		got.Edition != "Protestant Edition" ||
		got.RightsStatus != "public-domain" {
		t.Fatalf("unexpected translation: %#v", got)
	}
	if got.SourcePublisher == "" ||
		got.SourceHomepage == "" ||
		got.RightsNoticeURL == "" ||
		got.TrademarkNotice == "" ||
		got.TextPolicy == "" {
		t.Fatalf("translation attribution is incomplete: %#v", got)
	}
}
