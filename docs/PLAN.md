# Bible CLI execution plan

## 1. Vision

Bible CLI should make opening, navigating, and searching Scripture from a
terminal as direct as opening a local text file. Reading must work offline, and
the CLI must compose naturally with standard shell tools.

## 2. Product principles

1. **Offline first:** reading and searching the bundled translation require no
   account, API key, or network connection.
2. **Fast path first:** common references such as `John 3`, `Jn 3:16`, and
   `John 3:16-21` should work without flags or prompts.
3. **Human and machine friendly:** interactive output may use styling, while
   redirected output is stable plain text.
4. **Licensing is a feature:** every distributed translation has recorded
   provenance, license terms, version, and attribution.
5. **Small, testable core:** reference parsing, navigation, storage, and
   rendering remain independent of the command framework.

## 3. MVP scope

### Included

- Read a whole chapter, one verse, or an inclusive verse range.
- Resolve canonical book names and a documented set of aliases.
- List books and chapter counts.
- Navigate to the previous or next chapter.
- Search verse text offline.
- Automatically disable decoration when output is redirected.
- Choose a translation when more than one is installed.
- Show translation attribution and application version information.
- Support macOS and Linux initially; keep Windows compatibility in the design.

### Deferred

- Full-screen interactive TUI.
- Accounts and cross-device synchronization.
- Notes, highlights, and reading plans.
- Remote translation marketplace or automatic downloads.
- Commentary, dictionaries, and other study resources.

## 4. Proposed command surface

```console
bible read John 3
bible read "John 3:16"
bible read "John 3:16-21"
bible read John 3 --next
bible search "living water"
bible books
bible random
bible translations
bible version
```

Global options should include `--translation`, `--plain`, `--no-color`, and
`--help`. Commands should use conventional exit codes and write diagnostics to
standard error.

## 5. Technical direction

### Runtime and interface

- **Go** for fast startup, portability, and single-binary distribution.
- **Cobra** for subcommands, help, completion, and argument validation.
- **Lip Gloss** for restrained terminal styling.
- **Bubble Tea** only after the command-based reader is complete and tested.

### Data

- Store normalized Bible content in SQLite.
- Embed the default database in the executable for the first release.
- Use SQLite FTS for local search.
- Keep data importers separate from runtime code so generated databases are
  deterministic and auditable.
- Record translation identifier, language, name, source URL, source version,
  license, attribution, and import checksum.

### Suggested layout

```text
cmd/bible/              application entry point
internal/reference/    reference parser and book aliases
internal/bible/        reading and navigation use cases
internal/search/       search behavior
internal/storage/      SQLite repositories
internal/render/       terminal and plain-text output
internal/config/       preferences and platform paths
data/                   source metadata and import tooling
docs/                   plans and architecture decisions
```

The domain and reference parser should not depend on Cobra, terminal styling,
or concrete database types.

## 6. Milestones

### M0 — Foundation

- Initialize the Go module and command entry point.
- Add formatting, linting, unit-test, and build commands.
- Add GitHub Actions for supported platforms.
- Document contribution and licensing expectations.

**Exit criteria:** a version command builds and passes CI on macOS and Linux.

### M1 — Licensed text pipeline

Progress: complete. ADR 0001 selects the World English Bible Protestant Edition
and its official verse-per-line source. A validated provenance manifest pins the
source snapshot. The deterministic VPL importer and normalized SQLite schema are
implemented, and the generated database is embedded in the reader.

- Select and verify a redistributable first translation.
- Document provenance, license, attribution, and source checksum.
- Define the normalized schema.
- Build a deterministic importer and validation checks.
- Generate and embed the initial SQLite database.

**Exit criteria:** the importer reproduces a validated database with expected
book, chapter, and verse totals, and all required attribution is present.

### M2 — Reading

Progress: complete. The validated WEBP database is embedded directly in the
executable. The reader supports chapters, individual verses, inclusive
same-chapter ranges, aliases for all 66 books, cross-book previous/next chapter
navigation, and book discovery. Interactive output uses restrained styling;
non-TTY output automatically switches to stable plain text without ANSI escape
sequences. Global `--plain` and `--no-color` flags provide explicit control.

- Implement the reference grammar and aliases.
- Implement chapter, verse, and range reads.
- Add previous/next navigation.
- Render styled terminal output and automatic plain output.

**Exit criteria:** documented reference forms return correct verses, invalid or
ambiguous references give useful errors, and piped output contains no ANSI
escape sequences.

### M3 — Search and discovery

Progress: complete. Offline full-text search uses a deterministic SQLite FTS5
index, relevance ranking, stable canonical tie-breaking, bounded result limits,
adaptive terminal/plain rendering, and exact interactive match highlighting.
Book and translation discovery are implemented, including complete attribution
metadata, and random selection is uniform across all bundled verses.

- Add full-text search.
- Add book listing, translation listing, and random verse commands.
- Define result limits and stable output behavior.

**Exit criteria:** search works offline, is covered by representative tests, and
returns results within an agreed performance budget.

### M4 — Preferences and distribution

Progress: complete. Tagged release automation builds checksummed macOS and
Linux archives for AMD64 and ARM64, installation and release procedures are
documented, and native macOS/Linux CI executes the matching packaged binary.
The CLI generates Bash, Zsh, Fish, and PowerShell completion scripts. macOS and
Linux share an XDG-compatible configuration path, and versioned persistent
preferences remember translation and display defaults.

- Add a consistent, documented configuration path. (complete)
- Remember the preferred translation and optional display preferences. (complete)
- Produce checksummed release binaries. (complete)
- Add installation documentation and shell completion. (complete)

**Exit criteria:** a new user can install a release artifact and read a verse
without installing Go or downloading data separately.

### M5 — Post-MVP validation

- Gather usage feedback on navigation and search.
- Decide whether bookmarks, reading history, or a TUI is the next priority.
- Add features only after defining their storage and compatibility guarantees.

## 7. Quality strategy

- Table-driven unit tests for parsing, especially numbered books, aliases,
  ranges, whitespace, and invalid input.
- Repository contract tests against a small fixture database.
- Golden tests for terminal and plain-text rendering.
- Import validation for duplicate or missing verse keys and unexpected counts.
- End-to-end smoke tests that build and invoke the binary.
- Benchmarks for startup, chapter reads, and search before setting budgets.

## 8. Initial decisions to record

Before M1 is complete, create short architecture decision records for:

1. first translation and its redistribution terms;
2. canonical book identifiers and ordering;
3. reference grammar and ambiguity rules;
4. embedded database versus external data files; and
5. compatibility policy for configuration and translation databases.

## 9. Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| Translation redistribution is not permitted | Require a license review and provenance record before importing text |
| Reference aliases become ambiguous across languages | Scope aliases by language and return explicit ambiguity errors |
| Embedded data makes binaries too large | Measure the first database; support optional external packs if needed |
| Styled output breaks scripts | Detect non-TTY output and provide explicit `--plain`/`--no-color` flags |
| Search semantics surprise users | Document tokenization and test representative phrases and punctuation |

## 10. First implementation backlog

1. Scaffold the Go module and `bible version` command.
2. Establish CI and local quality commands.
3. Research and document candidates for the first translation.
4. Write the reference grammar and its tests.
5. Define the database schema and importer contract.

The next planning review should happen after M2. At that point, actual usage of
the reader should guide search presentation and the post-MVP roadmap.
