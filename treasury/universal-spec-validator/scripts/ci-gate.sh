#!/usr/bin/env bash
# universal-spec-validator CI / pre-commit wrapper.
# Usage:  ci-gate.sh <spec-or-glob> [--baseline PATH] [--config PATH] [--out DIR]
# Exit code propagates the gate signal (0 pass, 1 blocked, 2 error) so it can
# be dropped into a CI job step or a pre-commit hook unchanged.
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec python3 "${HERE}/validate_spec.py" "$@"
