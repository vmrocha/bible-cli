# Installation

Bible CLI releases are single-binary archives for macOS and Linux on AMD64 and
ARM64. Downloads include checksums for verification and do not require a
GitHub account.

## Download a release

1. Open the [latest release](https://github.com/vmrocha/bible-cli/releases/latest).
2. Download the archive that matches your operating system and architecture.
3. Verify the archive against `checksums.txt`, then extract it and put `bible`
   on your `PATH`.

See the repository's [full installation guide](https://github.com/vmrocha/bible-cli/blob/main/docs/INSTALL.md)
for platform-specific commands, source builds, and checksum verification.

## Confirm the installation

```console
$ bible version
$ bible read "John 3:16"
```

## Shell completion

Generate completion scripts directly from the installed binary:

```console
bible completion bash
bible completion zsh
bible completion fish
bible completion powershell
```

The [installation guide](https://github.com/vmrocha/bible-cli/blob/main/docs/INSTALL.md#shell-completion)
has setup instructions for each shell.
