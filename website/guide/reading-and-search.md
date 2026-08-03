# Reading and search

Bible CLI accepts common book names and aliases. Quote a reference containing
punctuation or a verse range so your shell passes it as one argument.

## Read Scripture

```console
bible read John 3
bible read "John 3:16"
bible read "John 3:16-21"
bible read Jn 3:16 --plain
```

Move through chapters, including across book boundaries:

```console
bible read John 3 --next
bible read Matthew 1 --previous
```

## Search offline

Search returns verses containing every query token, ranked by relevance with a
stable canonical order as a tie-breaker.

```console
bible search "living water"
bible search "faith hope love" --limit 10
bible search "kingdom of God" --plain
```

Punctuation and case do not change matching. A plain search with no matches
writes no output, making it convenient to compose with other shell tools.

## Output for people and programs

Interactive output uses restrained terminal styling. Redirected output switches
to plain text automatically; `--plain` and `--no-color` let you choose
explicitly.

```console
bible read "Psalm 23"          # styled terminal output
bible read "Psalm 23" --plain  # stable tab-separated output
bible read "Psalm 23" | less   # automatically plain
```
