#!/usr/bin/env bash
# run_agy.sh — dispatch a headless agy (Antigravity CLI) consult at an EXACT model label,
# under a macOS-safe watchdog, then ground-truth the backend from the agy CLI log.
#
# Usage:
#   run_agy.sh "<MODEL LABEL>" "<prompt>" [neutral_working_dir]
#
#   <MODEL LABEL>  Exact display label from `agy models` (e.g. "Gemini 3.1 Pro (High)").
#                  A wrong or omitted label SILENTLY falls back to the settings.json default
#                  model — this wrapper refuses unknown labels up front so that never happens.
#   <prompt>       The consult prompt. Keep it tightly scoped: agy --print is agentic and will
#                  wander into the workspace (shell history, dotfiles, repo) on vague prompts.
#   [neutral dir]  Optional. If given, agy runs from there to curb workspace-scanning. Use an
#                  empty, trusted directory. Pass absolute paths inside the prompt for files to read.
#
# Env:
#   TIMEOUT      Watchdog seconds (default 600). macOS has no `timeout`, so this wrapper
#                implements its own background-PID watchdog.
#   AGY_LOG_DIR  agy CLI log directory (default ~/.gemini/antigravity-cli/log).
#
# Output: agy's stdout/stderr, then a single PROVENANCE: line carrying the backend label parsed
#         from the newest log — agy injects a Gemini identity, so its self-report is worthless;
#         trust this line, not the model's words.
#
# Exit: 0 ok · 1 usage/label/precondition error · 124 watchdog timeout · else agy's own exit code.

set -euo pipefail

TIMEOUT="${TIMEOUT:-600}"
AGY_LOG_DIR="${AGY_LOG_DIR:-$HOME/.gemini/antigravity-cli/log}"

[ "$#" -ge 2 ] || { echo 'usage: run_agy.sh "<MODEL LABEL>" "<prompt>" [neutral_working_dir]' >&2; exit 1; }
LABEL="$1"; PROMPT="$2"; WORKDIR="${3:-}"

command -v agy >/dev/null 2>&1 || { echo "error: agy not found on PATH" >&2; exit 1; }

# 1. Verify the EXACT label (whole-line match) so we never get a silent Flash fallback.
if ! agy models 2>/dev/null | grep -Fxq -- "$LABEL"; then
  {
    echo "error: model label not found: \"$LABEL\""
    echo "valid labels (from \`agy models\`):"
    agy models 2>/dev/null | sed 's/^/  - /'
  } >&2
  exit 1
fi

# 2. Dispatch under a macOS-safe watchdog. --model MUST precede --print/the prompt, or Go's
#    flag parser stops at the first positional and silently drops it (-> default model).
OUT="$(mktemp "${TMPDIR:-/tmp}/agy-out.XXXXXX")"
trap 'rm -f "$OUT"' EXIT
(
  [ -n "$WORKDIR" ] && cd "$WORKDIR"
  exec agy --model "$LABEL" --print "$PROMPT"
) >"$OUT" 2>&1 </dev/null &
pid=$!
( sleep "$TIMEOUT"; kill -TERM "$pid" 2>/dev/null ) & watcher=$!

set +e
wait "$pid"; status=$?
set -e
kill "$watcher" 2>/dev/null || true
wait "$watcher" 2>/dev/null || true

cat "$OUT"

# 3. Ground-truth the backend from the newest log (do not trust agy's self-reported identity).
newest_log="$(ls -t "$AGY_LOG_DIR"/cli-*.log 2>/dev/null | head -1 || true)"
backend="(unverified)"
if [ -n "$newest_log" ]; then
  line="$(grep -F 'Propagating selected model override to backend' "$newest_log" 2>/dev/null | tail -1 || true)"
  if [ -n "$line" ]; then
    backend="$(printf '%s' "$line" | sed -n 's/.*label="\([^"]*\)".*/\1/p')"
    [ -n "$backend" ] || backend="(override line present, label unparsed)"
  else
    backend="(no backend-override line in newest log)"
  fi
fi
echo "PROVENANCE: cli=agy requested_label=\"$LABEL\" backend_label=\"$backend\" timeout=${TIMEOUT}s exit=${status} log=\"${newest_log:-none}\""

if [ "$status" -eq 143 ]; then
  echo "error: agy exceeded the ${TIMEOUT}s watchdog and was killed (SIGTERM)" >&2
  exit 124
fi
exit "$status"
