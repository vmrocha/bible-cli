# ADR 0002: Use one XDG-compatible configuration path

- Status: accepted
- Date: 2026-07-16

## Context

Operating-system-native conventions would normally place application
configuration under `~/Library/Application Support` on macOS and under
`$XDG_CONFIG_HOME` on Linux. That difference makes shell scripts, dotfile
management, documentation, and switching between macOS and Linux less
predictable for a terminal-first application.

Bible CLI also needs a stable format for persistent defaults without
allowing saved values to silently defeat explicit command-line arguments.

## Decision

Use the same resolution order on macOS and Linux:

1. `$BIBLE_CLI_CONFIG_HOME/config.json`
2. `$XDG_CONFIG_HOME/bible-cli/config.json`
3. `~/.config/bible-cli/config.json`

Environment-provided directories must be absolute. The JSON document has an
explicit schema version and is decoded strictly. Writes use a temporary file,
owner-only permissions, and an atomic rename.

Saved preferences are defaults. Explicit flags take precedence, including
boolean negation such as `--plain=false` and `--no-color=false`. Commands for
finding and resetting the configuration remain available when the stored file
cannot be decoded.

## Consequences

- macOS and Linux users see identical paths and can share setup instructions.
- The default differs from the macOS graphical-application convention, which is
  acceptable for a terminal-focused program.
- `BIBLE_CLI_CONFIG_HOME` gives tests, portable installations, and users a
  direct application-specific override.
- Future schema changes must preserve version compatibility or provide a clear
  migration and error message.
