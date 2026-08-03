#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "usage: $0 <version> <commit> <archive-directory>" >&2
  exit 2
fi

version=$1
commit=$2
archive_directory=$3

case "$(uname -s)/$(uname -m)" in
  Darwin/x86_64) target=darwin_amd64 ;;
  Darwin/arm64) target=darwin_arm64 ;;
  Linux/x86_64) target=linux_amd64 ;;
  Linux/aarch64 | Linux/arm64) target=linux_arm64 ;;
  *)
    echo "unsupported smoke-test platform: $(uname -s)/$(uname -m)" >&2
    exit 2
    ;;
esac

archive="bible-cli_${version#v}_${target}.tar.gz"
archive_path="$archive_directory/$archive"
checksums_path="$archive_directory/checksums.txt"
if [[ ! -f $archive_path || ! -f $checksums_path ]]; then
  echo "release archive or checksum file is missing for $target" >&2
  exit 1
fi

(
  cd "$archive_directory"
  checksum_line=$(
    awk -v archive="$archive" '
      $2 == archive { print; found = 1 }
      END { if (!found) exit 1 }
    ' checksums.txt
  )
  if command -v sha256sum >/dev/null 2>&1; then
    printf '%s\n' "$checksum_line" | sha256sum --check -
  else
    printf '%s\n' "$checksum_line" | shasum -a 256 --check -
  fi
)

temp_root=${TMPDIR:-/tmp}
temp_root=${temp_root%/}
work=$(mktemp -d "$temp_root/bible-cli-smoke.XXXXXX")
trap 'rm -rf "$work"' EXIT
tar -xzf "$archive_path" -C "$work"
bible="$work/bible"
config_home="$work/config"

test -x "$bible"
test -f "$work/LICENSE"
test -f "$work/THIRD_PARTY_NOTICES.md"
grep -F "MIT License" "$work/LICENSE"
grep -F "World English Bible — Protestant Edition" "$work/THIRD_PARTY_NOTICES.md"
grep -F "public domain" "$work/THIRD_PARTY_NOTICES.md"
version_output=$("$bible" version)
grep -F "bible $version" <<<"$version_output"
grep -F "commit: $commit" <<<"$version_output"

config_path=$(BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" config path)
test "$config_path" = "$config_home/config.json"
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" config set plain true
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" config set color false
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" --plain config show |
  grep -F $'translation\tengwebp'
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" --plain config show |
  grep -F $'plain\ttrue'
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" --plain config show |
  grep -F $'color\tfalse'

BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" read "John 3:16" |
  grep -F "For God so loved the world"
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" search "God loved world" --limit 1 |
  grep -F "John 3:16"
BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" completion bash |
  grep -F "__start_bible"

BIBLE_CLI_CONFIG_HOME="$config_home" "$bible" config reset
test ! -e "$config_home/config.json"
