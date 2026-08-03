#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 5 ]]; then
  echo "usage: $0 <version> <commit> <build-date> <source-epoch> <output-directory>" >&2
  exit 2
fi

version=$1
commit=$2
build_date=$3
source_epoch=$4
output=$5

if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
  echo "version must be a semantic version beginning with v" >&2
  exit 2
fi
if [[ ! $commit =~ ^[0-9a-fA-F]{7,64}$ ]]; then
  echo "commit must be a hexadecimal Git object ID" >&2
  exit 2
fi
if [[ ! $source_epoch =~ ^[0-9]+$ ]]; then
  echo "source epoch must be an integer" >&2
  exit 2
fi
if [[ -z $build_date || $build_date =~ [[:space:]] ]]; then
  echo "build date must be a non-empty value without whitespace" >&2
  exit 2
fi

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
mkdir -p "$output"
output=$(cd "$output" && pwd)
if compgen -G "$output/bible-cli_*.tar.gz" >/dev/null || [[ -e $output/checksums.txt ]]; then
  echo "output directory already contains release artifacts: $output" >&2
  exit 2
fi
work=$(mktemp -d "${TMPDIR:-/tmp}/bible-cli-release.XXXXXX")
trap 'rm -rf "$work"' EXIT

export CGO_ENABLED=0
export LC_ALL=C
cd "$root"

targets=(
  darwin/amd64
  darwin/arm64
  linux/amd64
  linux/arm64
)

for target in "${targets[@]}"; do
  goos=${target%/*}
  goarch=${target#*/}
  name="bible-cli_${version#v}_${goos}_${goarch}"
  stage="$work/$name"
  mkdir -p "$stage"

  GOOS=$goos GOARCH=$goarch go build \
    -trimpath \
    -ldflags "-s -w \
      -X github.com/vmrocha/bible-cli/internal/buildinfo.version=$version \
      -X github.com/vmrocha/bible-cli/internal/buildinfo.commit=$commit \
      -X github.com/vmrocha/bible-cli/internal/buildinfo.date=$build_date" \
    -o "$stage/bible" \
    ./cmd/bible

  cp LICENSE THIRD_PARTY_NOTICES.md "$stage/"

  tar \
    --sort=name \
    --mtime="@$source_epoch" \
    --owner=0 \
    --group=0 \
    --numeric-owner \
    -C "$stage" \
    -cf - \
    LICENSE THIRD_PARTY_NOTICES.md bible |
    gzip -n >"$output/$name.tar.gz"
done

(
  cd "$output"
  sha256sum bible-cli_*.tar.gz >checksums.txt
)
