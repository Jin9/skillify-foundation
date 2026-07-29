#!/usr/bin/env bash
# run_codex.sh — dispatch a headless codex (Codex CLI, `codex exec`) consult under a macOS-safe
# watchdog, then ground-truth the model codex actually resolved from its own run header.
#
# Usage:
#   run_codex.sh "<prompt>" [model]
#     [model] is OPTIONAL and normally omitted. When omitted, no -m flag is passed and codex
#     resolves the model itself from ~/.codex/config.toml (plus any --profile or -c override).
#     That is deliberate: hardcoding a default here would silently override the user's configured
#     model and go stale the moment the account tier moves. Pass [model] only to pin a specific
#     model for this one call.
#
# Env:
#   TIMEOUT                       watchdog seconds (default 600, minimum 45). macOS has no `timeout`.
#   CODEX_ALLOW_UNVERIFIED_MODEL  set to 1 to bypass the known-rejected-model guard below.
#
# Output: codex's output, then a PROVENANCE: line carrying the model parsed from codex's own run
#         header — the resolved model, not an assumption.
# Exit: 0 ok · 1 usage/model error · 124 watchdog timeout · else codex's own exit code.

set -euo pipefail
TIMEOUT="${TIMEOUT:-600}"

[ "$#" -ge 1 ] || { echo 'usage: run_codex.sh "<prompt>" [model]' >&2; exit 1; }
PROMPT="$1"; MODEL="${2:-}"

command -v codex >/dev/null 2>&1 || { echo "error: codex not found on PATH" >&2; exit 1; }

case "$TIMEOUT" in
  ''|*[!0-9]*) echo "error: TIMEOUT must be whole seconds (got '$TIMEOUT')" >&2; exit 1;;
esac
[ "$TIMEOUT" -ge 45 ] || { echo "error: TIMEOUT must be >= 45 seconds" >&2; exit 1; }

# Model guard, only when the caller pinned a model. These are dated observations on one ChatGPT
# account, not permanent truths — hence the escape hatch.
if [ -n "$MODEL" ]; then
  case "$MODEL" in
    gpt-5|gpt-5-codex)
      if [ "${CODEX_ALLOW_UNVERIFIED_MODEL:-0}" != "1" ]; then
        echo "error: '$MODEL' is rejected by the service on this ChatGPT account (re-confirmed 2026-07-30)." >&2
        echo "       Omit the model argument to use your configured default, or set" >&2
        echo "       CODEX_ALLOW_UNVERIFIED_MODEL=1 to try it anyway." >&2
        exit 1
      fi
      echo "warning: dispatching '$MODEL' despite a recorded rejection (override in effect)" >&2;;
    gpt-5.5)
      echo "warning: '$MODEL' is a legacy model here (last used 2026-06); it may be withdrawn." >&2;;
  esac
fi

OUT="$(mktemp "${TMPDIR:-/tmp}/codex-out.XXXXXX")"
KILLED="${OUT}.killed"
trap 'rm -f "$OUT" "$KILLED"' EXIT

start=$(date +%s)
if [ -n "$MODEL" ]; then
  codex exec -m "$MODEL" "$PROMPT" >"$OUT" 2>&1 </dev/null &
else
  codex exec "$PROMPT" >"$OUT" 2>&1 </dev/null &
fi
pid=$!
( sleep "$TIMEOUT"; : >"$KILLED"; kill -TERM "$pid" 2>/dev/null ) & watcher=$!
set +e
wait "$pid"; status=$?
set -e
kill "$watcher" 2>/dev/null || true
wait "$watcher" 2>/dev/null || true
elapsed=$(( $(date +%s) - start ))

# Ground-truth: codex prints its resolved configuration as a header before the response.
# Take the FIRST match so a model-authored line can never impersonate the header.
resolved="$(LC_ALL=C grep -m1 '^model: ' "$OUT" 2>/dev/null | sed 's/^model: *//' || true)"
[ -n "$resolved" ] || resolved="(unverified)"

verified=no
if [ "$resolved" != "(unverified)" ]; then
  if [ -z "$MODEL" ] || [ "$resolved" = "$MODEL" ]; then verified=yes; else verified=MISMATCH; fi
fi

cat "$OUT"
echo "PROVENANCE: cli=codex requested_model=\"${MODEL:-(config default)}\" resolved_model=\"$resolved\"" \
     "model_verified=${verified} watchdog=${TIMEOUT}s elapsed=${elapsed}s exit=${status}"

if [ -e "$KILLED" ] || [ "$status" -eq 143 ] || [ "$status" -eq 137 ]; then
  echo "error: codex exceeded the ${TIMEOUT}s watchdog and was killed" >&2
  exit 124
fi

if [ "$verified" = "MISMATCH" ]; then
  echo "error: requested \"$MODEL\" but codex resolved \"$resolved\" — discard this output" >&2
  exit 3
fi

exit "$status"
