#!/usr/bin/env bash
# run_agy.sh — dispatch a headless agy (Antigravity CLI) consult at a requested model, under a
# macOS-safe watchdog, then ground-truth the model the backend ACTUALLY ran from a per-run log.
#
# Usage:
#   run_agy.sh "<MODEL>" "<prompt>" [neutral_working_dir]
#
#   <MODEL>   agy accepts either form: a display label ("Gemini 3.1 Pro (High)") or a slug from
#             `agy models` (gemini-3.6-flash-low). They are NOT equally reliable — at least one
#             listed slug silently resolves to the settings.json default instead of erroring.
#             That is why this wrapper verifies the backend after the run and fails on a mismatch
#             (exit 3). Prefer the display-label form; verify, never assume.
#   <prompt>  The consult prompt. Keep it tightly scoped: agy --print is agentic and will wander
#             into the workspace (shell history, dotfiles, repo) on vague prompts.
#   [neutral dir]  Optional. If given, agy runs from there to curb workspace-scanning. Use an empty,
#             trusted directory. Pass absolute paths inside the prompt for files to read.
#
# Env:
#   TIMEOUT      Watchdog seconds (default 600, whole seconds, minimum 45). macOS has no `timeout`,
#                so this wrapper implements its own background-PID watchdog.
#   AGY_LOG_DIR  Directory for this run's log (default ~/.gemini/antigravity-cli/log).
#
# Why there is no `agy models` pre-flight: that call is server-backed and has been observed to hang
# indefinitely, and it ran BEFORE the watchdog was armed — so it could wedge the whole dispatch. agy
# already rejects an unrecognized model itself (exit 1, printing the valid labels), so the pre-flight
# bought nothing it can be trusted for. The post-run backend check below is the real guard.
#
# Output: agy's stdout/stderr, then a single PROVENANCE: line. agy injects a Gemini identity, so the
#         model's self-report is worthless — trust the PROVENANCE line, not the model's words.
#
# Exit: 0 ok · 1 usage/precondition error · 3 backend model mismatch (silent fallback detected)
#       · 124 watchdog or agy print-timeout · else agy's own exit code.

set -euo pipefail

TIMEOUT="${TIMEOUT:-600}"
AGY_LOG_DIR="${AGY_LOG_DIR:-$HOME/.gemini/antigravity-cli/log}"

[ "$#" -ge 2 ] || { echo 'usage: run_agy.sh "<MODEL>" "<prompt>" [neutral_working_dir]' >&2; exit 1; }
MODEL="$1"; PROMPT="$2"; WORKDIR="${3:-}"

command -v agy >/dev/null 2>&1 || { echo "error: agy not found on PATH" >&2; exit 1; }

case "$TIMEOUT" in
  ''|*[!0-9]*) echo "error: TIMEOUT must be whole seconds (got '$TIMEOUT')" >&2; exit 1;;
esac
[ "$TIMEOUT" -ge 45 ] || { echo "error: TIMEOUT must be >= 45 seconds" >&2; exit 1; }

# Normalize a model identifier for comparison: lowercase, keep only [a-z0-9]. This makes the slug
# and display-label forms comparable:
#   gemini-3.6-flash-low  ->  gemini36flashlow  <-  "Gemini 3.6 Flash (Low)"
# LC_ALL=C because BSD tr aborts with "Illegal byte sequence" on non-UTF-8 input.
norm() { printf '%s' "${1:-}" | LC_ALL=C tr 'A-Z' 'a-z' | LC_ALL=C tr -cd 'a-z0-9'; }

# agy's own print timeout fires BEFORE the watchdog so it can exit gracefully and flush the log.
# The value is a Go duration and MUST carry a unit — a bare integer aborts at flag-parse.
PRINT_TO=$(( TIMEOUT - 30 ))

mkdir -p "$AGY_LOG_DIR"
RUN_LOG="$AGY_LOG_DIR/cli-wrapper-$(date +%Y%m%d_%H%M%S)-$$.log"
OUT="$(mktemp "${TMPDIR:-/tmp}/agy-out.XXXXXX")"
KILLED="${OUT}.killed"
trap 'rm -f "$OUT" "$KILLED"' EXIT

# Every flag precedes --print, and --print carries the prompt as its value.
start=$(date +%s)
(
  if [ -n "$WORKDIR" ]; then cd "$WORKDIR"; fi
  exec agy --log-file "$RUN_LOG" --model "$MODEL" --print-timeout "${PRINT_TO}s" --print "$PROMPT"
) >"$OUT" 2>&1 </dev/null &
pid=$!
( sleep "$TIMEOUT"; : >"$KILLED"; kill -TERM "$pid" 2>/dev/null ) & watcher=$!

set +e
wait "$pid"; status=$?
set -e
kill "$watcher" 2>/dev/null || true
wait "$watcher" 2>/dev/null || true
elapsed=$(( $(date +%s) - start ))

# Classify a timeout by evidence, not by guessing a signal number: agy may trap SIGTERM, and its own
# print-timeout exits 1 with "timeout waiting for response" rather than 143.
timed_out=0; reason=none
if [ -e "$KILLED" ]; then
  timed_out=1; reason="watchdog(${TIMEOUT}s)"
elif [ "$status" -eq 143 ] || [ "$status" -eq 137 ]; then
  timed_out=1; reason="signal(${status})"
elif [ "$status" -ne 0 ] && grep -qiF 'timeout waiting for response' "$OUT" 2>/dev/null; then
  timed_out=1; reason="agy-print-timeout(${PRINT_TO}s)"
fi

# Ground-truth the backend from THIS run's log. A per-run --log-file is what makes concurrent
# fan-out dispatches safe; scanning a shared directory for the newest log is a race.
backend="(unverified)"; verified=no
if [ -s "$RUN_LOG" ]; then
  line="$(LC_ALL=C grep -F 'Propagating selected model override to backend' "$RUN_LOG" 2>/dev/null | tail -1 || true)"
  if [ -n "$line" ]; then
    backend="$(printf '%s' "$line" | sed -n 's/.*label="\([^"]*\)".*/\1/p')"
    [ -n "$backend" ] || backend="(override line present, label unparsed)"
  else
    backend="(no backend-override line in log)"
  fi
fi

nb="$(norm "$backend")"; nm="$(norm "$MODEL")"
if [ -n "$nb" ] && [ -n "$nm" ] && [ "$backend" != "(unverified)" ]; then
  # Prefix-tolerant: a backend label may carry a qualifier the slug omits
  # (slug claude-sonnet-4-6 -> label "Claude Sonnet 4.6 (Thinking)").
  case "$nb" in
    "$nm"*) verified=yes;;
    *) case "$nm" in
         "$nb"*) verified=yes;;
         *) verified=MISMATCH;;
       esac;;
  esac
fi

cat "$OUT"

echo "PROVENANCE: cli=agy requested_model=\"$MODEL\" backend_label=\"$backend\"" \
     "backend_verified=${verified} print_timeout=${PRINT_TO}s watchdog=${TIMEOUT}s" \
     "elapsed=${elapsed}s timeout=${reason} exit=${status} log=\"$RUN_LOG\""

if [ "$timed_out" -eq 1 ]; then
  echo "error: agy timed out (${reason}); output above may be truncated or empty" >&2
  exit 124
fi

if [ "$verified" = "MISMATCH" ]; then
  echo "error: SILENT FALLBACK — requested \"$MODEL\" but the backend ran \"$backend\"." >&2
  echo "       agy accepted the value without erroring and used its default model instead." >&2
  echo "       Re-dispatch with the display-label form, and discard the output above." >&2
  exit 3
fi

if [ "$verified" != "yes" ] && [ "$status" -eq 0 ]; then
  echo "warning: could not verify the backend model from $RUN_LOG — treat the output as unattributed" >&2
fi

exit "$status"
