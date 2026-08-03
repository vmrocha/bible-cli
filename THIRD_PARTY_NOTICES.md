# Third-Party Notices

Bible CLI's source code is licensed under the MIT License. The bundled
Bible text described below is a separate work and is not licensed under the
project's MIT License.

## World English Bible — Protestant Edition

Bible CLI bundles the World English Bible, Protestant Edition:

- Application and source identifier: `engwebp`
- Abbreviation: `WEBP`
- Language: American English (`en-US`)
- Canon: Protestant 66-book canon
- Text edition: 2020 stable text edition
- Publisher and source: [eBible.org](https://ebible.org/engwebp/)

eBible.org identifies the World English Bible text as public domain. The name
“World English Bible” is a trademark of eBible.org and may be used to identify
faithful copies of the translation. If the actual translation text is changed,
the resulting work must not be identified as the World English Bible.

Bible CLI preserves the publisher-provided verse text after UTF-8 decoding
and removal of the source record separator. Reference identifiers, storage
metadata, terminal styling, and display wrapping are handled separately.

Official rights and edition information:

- [Public-domain and trademark notice](https://ebible.org/engwebp/copyright.htm)
- [WEBP edition details and downloads](https://ebible.org/find/details.php?id=engwebp)
- [World English Bible official site](https://worldenglish.bible/)

### Bundled source snapshot

The embedded database was generated from the official eBible.org
verse-per-line distribution recorded in
[`data/translations/engwebp/manifest.json`](data/translations/engwebp/manifest.json):

- Archive: `https://ebible.org/Scriptures/engwebp_vpl.zip`
- Archive member: `engwebp_vpl.txt`
- Retrieved: 2026-07-15
- SHA-256: `f5c7bd6d09cf5b9ddd188726f6a4dac0096228c4478401853de20a5d845d72a7`

The checksum identifies the exact publisher archive used to generate the
bundled database; the download URL itself may change over time.
