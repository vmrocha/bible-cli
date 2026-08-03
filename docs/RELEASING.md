# Releasing Bible CLI

Release tags publish checksummed, self-contained archives through GitHub
Actions. A release must come from a reviewed commit on `main`.

## Prepare

1. Confirm `main` is clean, synchronized with `origin/main`, and green in CI.
2. Choose a semantic version that has never been used.
3. Add `docs/releases/<version>.md` with user-facing release notes.
4. Confirm the installation examples and supported platforms are still
   accurate.
5. Run `make check` locally.

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
publishes them with `checksums.txt` and the matching release-notes file.

## Verify

1. Wait for the release workflow to complete successfully.
2. Confirm the GitHub release contains four archives and `checksums.txt`.
3. Download the archive for the current machine and verify its checksum.
4. Run `bible version` and confirm the version, commit, and build date.
5. Run `bible read "John 3:16"` without network access.

Never move or reuse a published tag. If publishing fails after a release becomes
visible, preserve the tag and diagnose the workflow before taking further
action.
