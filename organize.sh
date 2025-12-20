#!/usr/bin/env bash
set -euo pipefail

# Run from repo root, or pass the path to maasclient/testdata as first arg.
TESTDATA_DIR="${1:-maasclient/testdata}"

cd "$TESTDATA_DIR"

# 1) Create one folder per resource (based on the first token before "__")
# 2) Move files into that folder
# 3) Rename inside folder to drop the "<resource>__" prefix
#    e.g. bootresources__list__all.json -> bootresources/list__all.json
shopt -s nullglob

for f in *.json *.tmpl.json; do
  # Skip if no matches
  [[ -e "$f" ]] || continue

  resource="${f%%__*}"          # everything before first "__"
  rest="${f#${resource}__}"     # everything after "<resource>__"

  mkdir -p "$resource"

  # Move & rename (avoid overwriting accidentally)
  if [[ -e "$resource/$rest" ]]; then
    echo "ERROR: target exists: $resource/$rest (from $f)" >&2
    exit 1
  fi

  mv -- "$f" "$resource/$rest"
done

# Optional: show resulting tree (if 'tree' exists)
if command -v tree >/dev/null 2>&1; then
  echo
  tree -a .
else
  echo
  find . -maxdepth 2 -type f | sort
fi
