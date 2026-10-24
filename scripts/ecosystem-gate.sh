#!/usr/bin/env bash
# Validate ecosystem docs, versions, contracts and the partner registry.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PYTHON="${PYTHON:-}"
if [[ -z "$PYTHON" ]]; then
  if command -v python >/dev/null 2>&1; then
    PYTHON=python
  elif command -v python3 >/dev/null 2>&1; then
    PYTHON=python3
  else
    echo "python or python3 is required" >&2
    exit 1
  fi
fi

"$PYTHON" "$ROOT/scripts/ecosystem-gate.py" "$@"
