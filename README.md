# Bible CLI

An offline-first command-line application for reading and searching the Bible
without leaving the terminal.

## Status

The Go CLI is under active development and includes an offline WEBP reader
backed by an embedded SQLite database. Reading output is styled when used in a
terminal and automatically switches to stable plain text when redirected. See
[docs/PLAN.md](docs/PLAN.md) for the product scope, architecture, milestones,
and acceptance criteria.

## Product direction

The initial interface will be command-oriented and friendly to both people and
shell pipelines:

```console
bible read John 3
bible read "John 3:16"
bible read "John 3:16-21" --plain
bible search "for God so loved"
bible books
bible random
```

The currently implemented commands include:

```console
bible read John 3
bible read "John 3:16"
bible read "John 3:16-21"
bible read Jn 3:16 --plain
bible read John 3 --next
bible read Matthew 1 --previous
bible read "1 Cor 13"
bible books
bible search "living water"
bible random
bible translations
bible config show
```

`bible books` lists all canonical names, source codes, and accepted aliases.
Chapter navigation crosses book boundaries, so moving forward from John 21 reads
Acts 1 and moving backward from Matthew 1 reads Malachi 4.

Output adapts to its destination:

```console
bible read "Psalm 23"          # styled terminal output
bible read "Psalm 23" --plain  # stable tab-separated output
bible read "Psalm 23" | less   # automatically plain, with no ANSI escapes
bible read "Psalm 23" --no-color # readable layout without terminal colors
```

`--plain`, `--no-color`, and `--translation` (`-t`) are global flags and may
appear before or after the subcommand. The translation flag accepts a bundled
translation id or registered abbreviation.

Search works entirely offline and returns verses containing every query token,
ranked by relevance with canonical Scripture order as a stable tie-breaker.
Matching words are emphasized in interactive colored output:

```console
bible search "living water"
bible search "faith hope love" --limit 10
bible search "kingdom of God" --plain
```

The default result limit is 20. Use `--limit` (or `-n`) to request between 1 and
100 results. Punctuation and case do not affect matching. A plain search with no
matches writes no output, making it safe to use in shell pipelines.
Match highlighting never changes `--no-color`, `--plain`, or redirected verse
text.

Discovery commands also work offline:

```console
bible random
bible random --plain
bible translations
```

`bible random` selects uniformly from all verses in the selected translation.
`bible translations` reports the bundled text edition, language, canon, source,
public-domain rights notice, trademark notice, and publisher text policy.

## Configuration

Bible CLI uses the same configuration convention on macOS and Linux. The
path is resolved in this order:

1. `$BIBLE_CLI_CONFIG_HOME/config.json`
2. `$XDG_CONFIG_HOME/bible-cli/config.json`
3. `~/.config/bible-cli/config.json`

The first two environment variables must contain absolute paths. Inspect and
change preferences with the CLI instead of editing JSON directly:

```console
bible config path
bible config show
bible config set plain true
bible config set color false
bible config set translation webp
bible config reset
```

Saved preferences provide defaults. Explicit command-line flags take priority,
including `--plain=false` and `--no-color=false`. Redirected output remains
plain even when the saved plain preference is false.

## Local licensed translations

Licensed translations can be embedded locally without adding their text to
Git. A local registration source file can call `storage.RegisterEmbedded` and
`config.RegisterTranslation` during package initialization, while both that
file and its database remain ignored. The tracked application then supports the
translation across reading, search, navigation, random selection, discovery,
and configuration without making the licensed text part of public releases.

For the locally registered Brazilian Portuguese NVI database:

```console
bible --translation nvi read "João 3:16"
bible -t ptnvi search amor
bible config set translation nvi
bible books --translation nvi
```

Only use and distribute a locally embedded translation as allowed by its
publisher authorization.

The first release should:

- work completely offline;
- start quickly and ship as a single executable;
- understand common Bible reference formats and book aliases;
- produce readable terminal output and clean redirected output;
- use Bible text that is legally redistributable; and
- leave room for additional languages and translations.

## Proposed stack

- Go
- Cobra for CLI command parsing
- SQLite with full-text search
- Lip Gloss for terminal presentation
- Bubble Tea later, if a full-screen reader proves useful

## Development

The project requires Go 1.26 or newer and Make. Run the complete local check:

```console
make check
./bin/bible version
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the development workflow and Bible
text licensing requirements.

## Installation

Install the latest stable release with Homebrew on macOS or Linux:

```console
brew install vmrocha/tap/bible-cli
```

Tagged releases ship checksummed single-binary archives for macOS and Linux on
Intel/AMD64 and ARM64. Public release downloads require no GitHub account, API
token, or GitHub CLI. See the [installation guide](docs/INSTALL.md) for archive
verification, Homebrew upgrades, source builds, and Bash, Zsh, Fish, and
PowerShell completion setup.

## License

Bible CLI's source code is licensed under the [MIT License](LICENSE).
Bible translations are separate works with their own copyright and
redistribution terms; no translation should be bundled until its license and
required attribution have been verified and documented. See
[Third-Party Notices](THIRD_PARTY_NOTICES.md) for the bundled World English
Bible text's provenance, public-domain status, trademark notice, and text
policy.
