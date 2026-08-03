# Contributing

## Community standards

By participating in Bible CLI, you agree to follow the
[Code of Conduct](CODE_OF_CONDUCT.md). Report suspected vulnerabilities through
the private channels in the [Security Policy](SECURITY.md), not through a public
issue.

## Prerequisites

- Go 1.26 or newer
- Make

## Local workflow

Run the complete local validation suite before submitting a change:

```console
make check
```

The individual commands are also available:

```console
make fmt
make lint
make test
make build
./bin/bible version
```

Keep reference parsing, Bible reading behavior, storage, and terminal rendering
in separate packages. New behavior should include tests, and commands should
write data to standard output and diagnostics to standard error.

## Releases

Follow the complete checklist in [docs/RELEASING.md](docs/RELEASING.md).
Pushing a semantic-version tag such as `v0.1.0` runs the release workflow. It
reruns the full validation suite, cross-builds CGO-free macOS and Linux binaries
for AMD64 and ARM64, creates deterministic `tar.gz` archives, writes
`checksums.txt`, and publishes a GitHub release.

Create a release only from a reviewed, green commit on `main`, and do not reuse
or move a published release tag.

## Bible text and licensing

Do not commit Bible text or generated translation databases without documented
provenance, redistribution terms, attribution requirements, source version, and
a source checksum. Translation licensing must be reviewed before import work is
merged.
