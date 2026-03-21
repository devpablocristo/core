#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

bash "${ROOT_DIR}/scripts/validate-runtime-layout.sh"
bash "${ROOT_DIR}/scripts/validate-module-versions.sh"
bash "${ROOT_DIR}/scripts/test-go-modules.sh"
bash "${ROOT_DIR}/scripts/test-ai.sh"
