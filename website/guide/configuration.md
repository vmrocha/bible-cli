# Configuration

Bible CLI stores persistent display and translation preferences in a
versioned JSON file. On macOS and Linux, it resolves the path in this order:

1. `$BIBLE_TERMINAL_CONFIG_HOME/config.json`
2. `$XDG_CONFIG_HOME/bible-terminal/config.json`
3. `~/.config/bible-terminal/config.json`

The first two environment variables must be absolute paths.

## Manage preferences

Use the CLI rather than editing the file directly:

```console
bible config path
bible config show
bible config set plain true
bible config set color false
bible config set translation webp
bible config reset
```

Saved preferences are defaults. Command-line flags always take priority,
including `--plain=false` and `--no-color=false`.
