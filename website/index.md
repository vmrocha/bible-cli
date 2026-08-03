---
layout: home

hero:
  name: Bible CLI
  text: Scripture in your terminal.
  tagline: Read and search the World English Bible offline, with a fast command-line interface that works naturally in shell pipelines.
  actions:
    - theme: brand
      text: Install Bible CLI
      link: /guide/installation
    - theme: alt
      text: Read the guide
      link: /guide/reading-and-search

features:
  - title: Offline by default
    details: Read, navigate, and search the bundled translation without an account, API key, or network connection.
  - title: Built for the shell
    details: Human-friendly terminal output becomes stable plain text automatically when redirected or piped.
  - title: One small binary
    details: Download a checksummed release for macOS or Linux and start reading without installing Go or a database.
---

## Start with a passage

```console
$ bible read "John 3:16-21"
$ bible search "living water"
$ bible random
```

Bible CLI understands common book aliases, lets you navigate between chapters,
and includes the tools you need to discover books, translations, and shell
completion.

[View releases on GitHub](https://github.com/vmrocha/bible-cli/releases)
