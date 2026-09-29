package cli

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/buildinfo"
)

var testBuild = buildinfo.Info{
	Version: "v0.1.0-test",
	Commit:  "abc1234",
	Date:    "2026-07-15T12:00:00Z",
}

func execute(t *testing.T, args ...string) (string, error) {
	return executeWithOptions(t, nil, args...)
}

func executeWithOptions(t *testing.T, options []Option, args ...string) (string, error) {
	t.Helper()

	output := new(bytes.Buffer)
	options = append(options, func(configuration *configuration) {
		configuration.isTerminal = func(io.Writer) bool { return true }
	})
	command := New(testBuild, options...)
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs(args)

	err := command.Execute()
	return output.String(), err
}

func TestHelp(t *testing.T) {
	output, err := execute(t, "--help")
	if err != nil {
		t.Fatalf("execute --help: %v", err)
	}

	for _, expected := range []string{
		"Read the Bible from your terminal",
		"books",
		"completion",
		"config",
		"read",
		"random",
		"search",
		"translations",
		"version",
		"--plain",
		"--no-color",
		"--translation",
		"--help",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("help output does not contain %q", expected)
		}
	}
}

func TestTranslationFlagIsNormalizedForFactories(t *testing.T) {
	var selected string
	reader := &stubReader{passage: bible.Passage{
		Book:    bible.Book{Name: "John"},
		Chapter: 3,
		Verses:  []bible.Verse{{Chapter: 3, Number: 16, Text: "Text"}},
	}}
	factory := func(_ context.Context, translation string) (PassageReader, error) {
		selected = translation
		return reader, nil
	}
	if _, err := executeWithOptions(t, []Option{WithReaderFactory(factory)}, "--translation", "webp", "read", "John", "3:16"); err != nil {
		t.Fatalf("execute translated read: %v", err)
	}
	if selected != "engwebp" {
		t.Fatalf("factory received translation %q, want engwebp", selected)
	}
}

func TestVersionFlag(t *testing.T) {
	output, err := execute(t, "--version")
	if err != nil {
		t.Fatalf("execute --version: %v", err)
	}

	if output != "bible v0.1.0-test\n" {
		t.Fatalf("unexpected version output: %q", output)
	}
}

func TestUnknownCommand(t *testing.T) {
	_, err := execute(t, "unknown")
	if err == nil {
		t.Fatal("expected unknown command to return an error")
	}
}
