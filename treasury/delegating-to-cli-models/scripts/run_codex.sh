#!/usr/bin/env bash
# run_codex.sh — dispatch a headless codex (Codex CLI, `codex exec`) consult under a
# macOS-safe watchdog, guarding the model so an unsupported one fails fast instead of mid-run.
#
# Usage:
#   run_codex.sh "<prompt>" [model]
#     [model] defaults to gpt-5.5. On a ChatGPT-account login, gpt-5 and gpt-5-codex are
#     REJECTED by the service — this wrapper refuses them up front with a clear message.
#
# Env:
#   TIMEOUT  watchdog seconds (default 600). macOS has no `timeout`.
#
# Output: codex's output, then a PROVENANCE: line.
# Exit: 0 ok · 1 usage/model error · 124 watchdog timeout · else codex's own exit code.

set -euo pipefail
TIMEOUT="${TIMEOUT:-600}"

[ "$#" -ge 1 ] || { echo 'usage: run_codex.sh "<prompt>" [model]' >&2; exit 1; }
PROMPT="$1"; MODEL="${2:-gpt-5.5}"

command -v codex >/dev/null 2>&1 || { echo "error: codex not found on PATH" >&2; exit 1; }

case "$MODEL" in
  gpt-5|gpt-5-codex)
    echo "error: codex on a ChatGPT account rejects '$MODEL'; use gpt-5.5" >&2; exit 1;;
  gpt-5.5) ;;
  *) echo "warning: unverified codex model '$MODEL' (known-good on this setup: gpt-5.5)" >&2;;
esac

OUT="$(mktemp "${TMPDIR:-/tmp}/codex-out.XXXXXX")"
trap 'rm -f "$OUT"' EXIT
codex exec -m "$MODEL" "$PROMPT" >"$OUT" 2>&1 </dev/null &
pid=$!
( sleep "$TIMEOUT"; kill -TERM "$pid" 2>/dev/null ) & watcher=$!
set +e
wait "$pid"; status=$?
set -e
kill "$watcher" 2>/dev/null || true
wait "$watcher" 2>/dev/null || true

cat "$OUT"
echo "PROVENANCE: cli=codex model=\"$MODEL\" timeout=${TIMEOUT}s exit=${status}"

if [ "$status" -eq 143 ]; then
  echo "error: codex exceeded the ${TIMEOUT}s watchdog and was killed (SIGTERM)" >&2
  exit 124
fi
exit "$status"
