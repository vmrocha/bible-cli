# Releasing Bible CLI

Release tags publish checksummed, self-contained archives through GitHub
Actions. A release must come from a reviewed commit on `main`.

## Homebrew credentials

The `HOMEBREW_TAP_TOKEN` repository secret must contain a fine-grained GitHub
personal access token restricted to `vmrocha/homebrew-tap`, with repository
contents read and write access. The workflow uses this token only to update the
tap after a stable release; the tag-bound `GITHUB_TOKEN` cannot write to a
different repository.

## Prepare

1. Confirm `main` is clean, synchronized with `origin/main`, and green in CI.
2. Choose a semantic version that has never been used.
3. Add `docs/releases/<version>.md` with user-facing release notes.
4. Confirm the installation examples and supported platforms are still
   accurate.
5. Confirm the `HOMEBREW_TAP_TOKEN` repository secret can write contents in
   `vmrocha/homebrew-tap`.
6. Run `make check` locally.

Pull-request CI builds all four archives once and executes the matching packaged
binary on native macOS and Linux runners. The smoke test verifies checksums,
embedded version metadata, isolated persistent configuration, reading, search,
and shell completion.

## Publish

Create and push an annotated tag from the reviewed `main` commit:

```console
git tag -a v0.1.0 -m "Bible CLI v0.1.0"
git push origin v0.1.0
```

The release workflow reruns `make check`, creates deterministic archives, and
publishes them with `checksums.txt` and the matching release-notes file. For a
stable tag, it then updates `Formula/bible-cli.rb` in
`vmrocha/homebrew-tap` with the new source URL and SHA-256. Prerelease tags do
not update Homebrew.

## Verify

1. Wait for the release workflow to complete successfully.
2. Confirm the GitHub release contains four archives and `checksums.txt`.
3. Download the archive for the current machine and verify its checksum.
4. Run `bible version` and confirm the version, commit, and build date.
5. Run `bible read "John 3:16"` without network access.
6. Confirm the tap formula references the new tag, then run
   `brew update && brew upgrade bible-cli` and `brew test bible-cli`.

Never move or reuse a published tag. If publishing fails after a release becomes
visible, preserve the tag and diagnose the workflow before taking further
action.
