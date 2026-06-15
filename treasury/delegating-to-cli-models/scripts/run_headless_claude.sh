#!/usr/bin/env bash
# run_headless_claude.sh — dispatch a headless `claude -p` child under a macOS-safe watchdog,
# ALWAYS forcing --effort low.
#
# Why --effort low is mandatory: a headless child inherits the parent session's high effort and
# hangs (9+ min, zero stdout) on substantial prompts. Environment variables (CLAUDE_EFFORT,
# MAX_THINKING_TOKENS) do NOT override it — only the --effort flag does.
#
# Usage:
#   run_headless_claude.sh "<prompt>" [model]
#     [model] defaults to opus. Prefer opus or haiku for headless driver stages; sonnet is
#     slower/less reliable headless and triggers a warning here.
#
# Env:
#   TIMEOUT  watchdog seconds (default 600). macOS has no `timeout`.
#   EFFORT   override the effort level (default low). Raising this re-opens the hang risk.
#
# Output: claude's response, then a PROVENANCE: line.
# Exit: 0 ok · 1 usage error · 124 watchdog timeout · else claude's own exit code.

set -euo pipefail
TIMEOUT="${TIMEOUT:-600}"
EFFORT="${EFFORT:-low}"

[ "$#" -ge 1 ] || { echo 'usage: run_headless_claude.sh "<prompt>" [model]' >&2; exit 1; }
PROMPT="$1"; MODEL="${2:-opus}"

command -v claude >/dev/null 2>&1 || { echo "error: claude not found on PATH" >&2; exit 1; }
case "$MODEL" in
  *sonnet*) echo "warning: sonnet is slower/less reliable headless; prefer opus or haiku" >&2;;
esac

OUT="$(mktemp "${TMPDIR:-/tmp}/claude-out.XXXXXX")"
trap 'rm -f "$OUT"' EXIT
claude -p "$PROMPT" --effort "$EFFORT" --model "$MODEL" >"$OUT" 2>&1 </dev/null &
pid=$!
( sleep "$TIMEOUT"; kill -TERM "$pid" 2>/dev/null ) & watcher=$!
set +e
wait "$pid"; status=$?
set -e
kill "$watcher" 2>/dev/null || true
wait "$watcher" 2>/dev/null || true

cat "$OUT"
echo "PROVENANCE: cli=claude model=\"$MODEL\" effort=\"$EFFORT\" timeout=${TIMEOUT}s exit=${status}"

if [ "$status" -eq 143 ]; then
  echo "error: claude exceeded the ${TIMEOUT}s watchdog and was killed (SIGTERM)" >&2
  exit 124
fi
exit "$status"
