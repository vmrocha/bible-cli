# Installing Bible CLI

## Homebrew

Install the latest stable release from the project tap on macOS or Linux:

```console
brew install vmrocha/tap/bible-cli
```

The formula builds from the immutable tagged source archive and installs the
`bible` executable plus Bash, Zsh, and Fish completions. Upgrade or remove it
with the usual Homebrew commands:

```console
brew upgrade bible-cli
brew uninstall bible-cli
```

## Release archives

Tagged releases provide single-binary archives for:

- macOS on Intel (`darwin_amd64`)
- macOS on Apple silicon (`darwin_arm64`)
- Linux on x86-64 (`linux_amd64`)
- Linux on ARM64 (`linux_arm64`)

Download the archive for your platform and `checksums.txt` from the matching
[GitHub release](https://github.com/vmrocha/bible-cli/releases). For
example, once the repository is public, download version `0.1.0` for an Apple
silicon Mac without a GitHub account or GitHub CLI:

```console
curl --fail --location --remote-name \
  https://github.com/vmrocha/bible-cli/releases/download/v0.1.0/bible-cli_0.1.0_darwin_arm64.tar.gz
curl --fail --location --remote-name \
  https://github.com/vmrocha/bible-cli/releases/download/v0.1.0/checksums.txt
grep 'bible-cli_0.1.0_darwin_arm64.tar.gz$' checksums.txt | shasum -a 256 -c -
tar -xzf bible-cli_0.1.0_darwin_arm64.tar.gz
install -d ~/.local/bin
install -m 0755 bible ~/.local/bin/bible
```

On Linux, use `sha256sum --check -` in place of the `shasum` verification
command.

Ensure `~/.local/bin` is on your `PATH`, then verify the embedded release
metadata and read a verse:

```console
bible version
bible read "John 3:16"
```

Release archives contain the complete offline WEBP database. No separate text
download, account, API key, or network connection is needed after installation.

## Shell completion

Bible CLI generates completion scripts without requiring additional files.

### Bash

```console
mkdir -p ~/.local/share/bash-completion/completions
bible completion bash > ~/.local/share/bash-completion/completions/bible
```

### Zsh

```console
mkdir -p ~/.zfunc
bible completion zsh > ~/.zfunc/_bible
```

Add `fpath=(~/.zfunc $fpath)` to `~/.zshrc` before `compinit` is called, then
start a new shell.

### Fish

```console
mkdir -p ~/.config/fish/completions
bible completion fish > ~/.config/fish/completions/bible.fish
```

PowerShell completion is also available through `bible completion powershell`.

## Build from source

Building from source requires Go 1.26 or newer and Make:

```console
git clone https://github.com/vmrocha/bible-cli.git
cd bible-cli
make check
install -d ~/.local/bin
install -m 0755 bin/bible ~/.local/bin/bible
```
