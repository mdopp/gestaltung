#!/bin/sh
# Copy the shared look into a program that is not written in Go, at build
# time (ADR 0001). Pin a tag or commit; never load it from a service at
# runtime.
#
#   scripts/sync.sh <dest-dir> [ref]
#
# Without ref it copies from this checkout; with ref (a tag or commit) it
# fetches that version from GitHub. Writes <dest-dir>/gestaltung.css.
set -eu
dest=${1:?usage: sync.sh <dest-dir> [ref]}
ref=${2:-}
mkdir -p "$dest"
out="$dest/gestaltung.css"
tmp="$out.tmp"
if [ -z "$ref" ]; then
  here=$(cd "$(dirname "$0")/.." && pwd)
  cat "$here/tokens.css" > "$tmp"
  printf '\n' >> "$tmp"
  cat "$here/components.css" >> "$tmp"
else
  base="https://raw.githubusercontent.com/mdopp/gestaltung/$ref"
  curl -fsSL "$base/tokens.css" > "$tmp"
  printf '\n' >> "$tmp"
  curl -fsSL "$base/components.css" >> "$tmp"
fi
mv "$tmp" "$out"
echo "wrote $out${ref:+ (mdopp/gestaltung@$ref)}"
