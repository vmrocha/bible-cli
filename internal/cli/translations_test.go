package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
)

type stubTranslationReader struct {
	translations []bible.Translation
	err          error
	closed       bool
}

func (reader *stubTranslationReader) Translations(context.Context) ([]bible.Translation, error) {
	return reader.translations, reader.err
}

func (reader *stubTranslationReader) Close() error {
	reader.closed = true
	return nil
}

func TestTranslationsCommand(t *testing.T) {
	reader := &stubTranslationReader{translations: []bible.Translation{{
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
		TrademarkNotice: "Trademark notice.",
		TextPolicy:      "Preserve text.",
	}}}
	factory := func(context.Context) (TranslationReader, error) { return reader, nil }

	output, err := executeWithOptions(
		t,
		[]Option{WithTranslationReaderFactory(factory)},
		"translations", "--no-color",
	)
	if err != nil {
		t.Fatalf("execute translations: %v", err)
	}
	for _, expected := range []string{"WEBP", "World English Bible", "Public domain", "Trademark notice."} {
		if !strings.Contains(output, expected) {
			t.Fatalf("translation output does not contain %q:\n%s", expected, output)
		}
	}
	if !reader.closed {
		t.Fatal("translation reader was not closed")
	}
}

func TestTranslationsCommandAutomaticallyUsesPlainOutputWhenRedirected(t *testing.T) {
	reader := &stubTranslationReader{translations: []bible.Translation{{
		ID:           "engwebp",
		Abbreviation: "WEBP",
		Name:         "World English Bible",
	}}}
	factory := func(context.Context) (TranslationReader, error) { return reader, nil }
	output := new(bytes.Buffer)
	command := New(testBuild, WithTranslationReaderFactory(factory))
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"translations"})

	if err := command.Execute(); err != nil {
		t.Fatalf("execute redirected translations: %v", err)
	}
	if !strings.HasPrefix(output.String(), "engwebp\tWEBP\tWorld English Bible\t") {
		t.Fatalf("unexpected redirected output: %q", output.String())
	}
	if strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("redirected output contains ANSI escapes: %q", output.String())
	}
}

func TestTranslationsCommandReportsReaderErrorAndCloses(t *testing.T) {
	reader := &stubTranslationReader{err: errors.New("metadata failed")}
	factory := func(context.Context) (TranslationReader, error) { return reader, nil }

	_, err := executeWithOptions(t, []Option{WithTranslationReaderFactory(factory)}, "translations")
	if err == nil || !strings.Contains(err.Error(), "metadata failed") {
		t.Fatalf("expected metadata error, got %v", err)
	}
	if !reader.closed {
		t.Fatal("translation reader was not closed after an error")
	}
}
