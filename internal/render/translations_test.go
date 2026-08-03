package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
)

func translationFixture() bible.Translation {
	return bible.Translation{
		ID:              "engwebp",
		Name:            "World English Bible",
		Abbreviation:    "WEBP",
		LanguageTag:     "en-US",
		LanguageName:    "English",
		Edition:         "Protestant Edition",
		Canon:           "protestant-66",
		TextEdition:     "2020 stable text edition",
		SourcePublisher: "eBible.org",
		SourceHomepage:  "https://ebible.org/engwebp/",
		RightsStatus:    "public-domain",
		RightsNoticeURL: "https://ebible.org/engwebp/copyright.htm",
		TrademarkNotice: "World English Bible is a trademark of eBible.org.",
		TextPolicy:      "Preserve the publisher-provided verse text.",
	}
}

func TestTranslations(t *testing.T) {
	var output bytes.Buffer
	if err := Translations(&output, []bible.Translation{translationFixture()}, Options{}); err != nil {
		t.Fatalf("Translations: %v", err)
	}
	want := "Available translations\n\n" +
		"WEBP · World English Bible\n" +
		"English (en-US) · Protestant Edition\n" +
		"Text edition: 2020 stable text edition\n" +
		"Canon: protestant-66\n" +
		"Rights: Public domain · https://ebible.org/engwebp/copyright.htm\n" +
		"Source: eBible.org · https://ebible.org/engwebp/\n" +
		"Trademark: World English Bible is a trademark of eBible.org.\n" +
		"Text policy: Preserve the publisher-provided verse text.\n"
	if output.String() != want {
		t.Fatalf("unexpected translation output:\n%s", output.String())
	}
}

func TestStyledTranslations(t *testing.T) {
	var output bytes.Buffer
	if err := Translations(&output, []bible.Translation{translationFixture()}, Options{Color: true}); err != nil {
		t.Fatalf("Translations: %v", err)
	}
	for _, expected := range []string{
		"\x1b[1m\x1b[36mAvailable translations\x1b[0m",
		"\x1b[1m\x1b[36mWEBP\x1b[0m",
		"\x1b[2mRights:\x1b[0m",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("styled output does not contain %q: %q", expected, output.String())
		}
	}
}

func TestPlainTranslations(t *testing.T) {
	translation := translationFixture()
	translation.TextPolicy = "Preserve\ttext\nexactly."

	var output bytes.Buffer
	if err := Translations(&output, []bible.Translation{translation}, Options{Plain: true}); err != nil {
		t.Fatalf("Translations: %v", err)
	}
	fields := strings.Split(strings.TrimSuffix(output.String(), "\n"), "\t")
	if len(fields) != 12 {
		t.Fatalf("plain output has %d fields, want 12: %q", len(fields), output.String())
	}
	if fields[0] != "engwebp" || fields[1] != "WEBP" || fields[7] != "public-domain" {
		t.Fatalf("unexpected plain fields: %#v", fields)
	}
	if fields[11] != "Preserve text exactly." {
		t.Fatalf("plain text policy was not sanitized: %q", fields[11])
	}
}

func TestTranslationsWithNoEntries(t *testing.T) {
	var output bytes.Buffer
	if err := Translations(&output, nil, Options{}); err != nil {
		t.Fatalf("Translations: %v", err)
	}
	if output.String() != "No translations available.\n" {
		t.Fatalf("unexpected empty output: %q", output.String())
	}
}
