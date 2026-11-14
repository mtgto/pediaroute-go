#!/usr/bin/env bash
# Download Wikipedia dumps and generate application data.
#
# Usage: scripts/gen-data.sh YYYYMMDD [lang...]
#   e.g. scripts/gen-data.sh 20261101
#        scripts/gen-data.sh 20261101 ja
#
# Dumps are saved in data/dumps/, output goes to data/<lang>/.
# Safe to re-run: completed downloads are skipped, interrupted ones are resumed.
set -euo pipefail

VERSION="${1:-}"
if [[ ! "$VERSION" =~ ^[0-9]{8}$ ]]; then
  echo "Usage: $0 YYYYMMDD [lang...]" >&2
  exit 1
fi
shift
LANGS=("$@")
if [[ ${#LANGS[@]} -eq 0 ]]; then
  LANGS=(en ja)
fi

BASE_URL="${DUMP_BASE_URL:-https://dumps.wikimedia.your.org}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DATA_DIR="${DATA_DIR:-$ROOT/data}"
DUMP_DIR="$DATA_DIR/dumps"
GEN="$ROOT/build/gen"

if [[ ! -x "$GEN" ]]; then
  echo "$GEN not found. Run 'make' first." >&2
  exit 1
fi
mkdir -p "$DUMP_DIR"

download() {
  local url="$1" dest="$2"
  if [[ -f "$dest" ]]; then
    echo "skip (already downloaded): $dest"
    return
  fi
  echo "download: $url"
  # Resume a partial download; rename only after curl succeeds.
  curl --fail --location --retry 5 --retry-delay 5 -C - -o "$dest.part" "$url"
  mv "$dest.part" "$dest"
}

for lang in "${LANGS[@]}"; do
  for kind in page pagelinks linktarget; do
    file="${lang}wiki-${VERSION}-${kind}.sql.gz"
    download "$BASE_URL/${lang}wiki/${VERSION}/$file" "$DUMP_DIR/$file"
  done
done

for lang in "${LANGS[@]}"; do
  out="$DATA_DIR/$lang"
  mkdir -p "$out"
  echo "generate: $lang -> $out"
  "$GEN" -lang "$lang" -version "$VERSION" \
    -ip "$DUMP_DIR/${lang}wiki-${VERSION}-page.sql.gz" \
    -il "$DUMP_DIR/${lang}wiki-${VERSION}-pagelinks.sql.gz" \
    -ilt "$DUMP_DIR/${lang}wiki-${VERSION}-linktarget.sql.gz" \
    -o "$out"
done
