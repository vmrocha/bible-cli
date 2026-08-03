# ADR 0001: Use the World English Bible Protestant Edition first

- Status: Accepted
- Date: 2026-07-15
- Decision owners: Bible CLI maintainers

## Context

Bible CLI needs a translation that can be embedded in release binaries,
read offline, searched, and redistributed without requiring an account or a
commercial agreement. The first importer also needs a stable, structured,
publisher-provided source.

The official eBible.org pages identify the World English Bible as public domain
and provide several editions and machine-readable downloads. The Protestant
Edition, identified by eBible.org as `engwebp` or `WEBP`, contains the 66-book
protocanon, uses American English, and renders the Tetragrammaton as `LORD` or
`GOD`. It is published as a subset of the World English Bible Updated.

## Decision

The first bundled translation will be the **World English Bible Protestant
Edition (WEBP)** with source identifier `engwebp`.

The importer will consume `engwebp_vpl.txt` from eBible.org's official
verse-per-line archive. That distribution intentionally contains Bible text
only and omits formatting, paragraph breaks, notes, introductions, and section
titles. Using the publisher's text-only export avoids treating those omissions
as transformations made by Bible CLI.

Every imported snapshot must be pinned by an archive SHA-256 checksum and
validated against the expected book and verse totals in its manifest. The
current manifest records the source retrieved on 2026-07-15; it is metadata,
not permission to silently accept future content at the mutable download URL.

## Rights and naming constraints

The publisher states that the World English Bible text is in the public domain
and may be copied and redistributed. The name "World English Bible" is a
trademark and may be used to identify faithful copies. If Bible CLI changes
the translation text or punctuation, the changed work must not be presented as
the World English Bible.

The importer may normalize reference keys and storage metadata, but it must
preserve the verse text from the selected archive member byte-for-byte after
UTF-8 decoding and removal of the record separator. Display wrapping and
terminal styling do not modify the stored translation.

Authoritative references:

- [World English Bible official site](https://worldenglish.bible/)
- [WEBP edition and public-domain notice](https://ebible.org/engwebp/copyright.htm)
- [WEBP formats and metadata](https://ebible.org/find/details.php?id=engwebp)
- [Official VPL archive](https://ebible.org/Scriptures/engwebp_vpl.zip)

## Consequences

- The MVP begins with English and the Protestant 66-book canon.
- Canon identifiers and ordering must remain translation-aware so later
  Catholic, ecumenical, and non-English editions are possible.
- Footnotes and paragraph structure are outside the first import because the
  selected official source is verse text only.
- Source updates require an explicit manifest change, checksum review, import
  validation, and visible release note.
- The application must expose the translation name, edition, source, and rights
  notice in its translation-information output.

## Alternatives considered

- **World English Bible Classic:** also public domain and stable, but includes a
  wider canon and uses `Yahweh`; that adds canon-selection decisions to the MVP.
- **World English Bible Updated:** public domain and modern, but includes a wider
  canon than the selected Protestant subset.
- **A popular copyrighted translation:** deferred because redistribution would
  require explicit permission and ongoing compliance work.
