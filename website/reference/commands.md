# Command reference

## Reading

```console
bible read <reference> [--next | --previous]
```

Read a chapter, verse, or inclusive range. Examples: `John 3`, `John 3:16`,
and `John 3:16-21`.

## Search and discovery

```console
bible search <query> [--limit <1-100>]
bible books
bible random
bible translations
```

`search` works against the bundled offline text. `books` lists canonical names,
source codes, and accepted aliases. `random` chooses a verse uniformly from the
bundled translation.

## Preferences and shell integration

```console
bible config path | show | set <key> <value> | reset
bible completion <bash|zsh|fish|powershell>
bible version
```

## Global options

```console
--translation <id>
--plain
--no-color
--help
```

Use `bible <command> --help` for the complete, version-specific command help.
