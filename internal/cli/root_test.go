package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"

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
		"--help",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("help output does not contain %q", expected)
		}
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
