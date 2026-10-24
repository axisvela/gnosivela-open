#!/usr/bin/env bash
# Validate ecosystem docs, versions, contracts and the partner registry.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PYTHON="${PYTHON:-python}"

"$PYTHON" "$ROOT/scripts/ecosystem-gate.py"
