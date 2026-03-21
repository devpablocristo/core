#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

discover_go_modules() {
  find "${ROOT_DIR}" -type f -path '*/go/go.mod' -printf '%h\n' | sed "s#^${ROOT_DIR}/##" | sort
}

while IFS= read -r module; do
  [[ -n "${module}" ]] || continue
  echo "==> go test ${module}"
  (
    cd "${ROOT_DIR}/${module}"
    go test ./...
  )
done < <(discover_go_modules)
