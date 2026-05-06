#!/usr/bin/env bash
# state_summary.sh
# Emit the canonical orchestrator status line:
#   <workflow_id>: <current_stage> <status> (<elapsed>m, $<spent>/$<cap>, src=<mix>)
#
# Reads .agent/state.json from the current scaffold checkout. Read-only.
# Exits non-zero if state.json is missing or unparseable.

set -euo pipefail

STATE="${STATE:-.agent/state.json}"
[[ -f "$STATE" ]] || { echo "state.json not found at $STATE" >&2; exit 1; }
command -v jq >/dev/null || { echo "jq not installed" >&2; exit 1; }

jq -r '
  def src_mix:
    [.stages[] | select(.cost_source) | .cost_source]
    | (map(select(. == "real")) | length) as $r
    | (map(select(. == "estimate")) | length) as $e
    | "\($r) real, \($e) estimate";

  . as $root
  | (.current_stage // "—") as $cur
  | (.stages[$cur] // {}) as $s
  | ($s.status // "—") as $st
  | ($s.started_epoch // 0) as $start
  | (now | floor) as $nowt
  | (if $start > 0 then (($nowt - $start) / 60 | floor) else 0 end) as $elapsed
  | (.total_cost_usd // 0) as $spent
  | (.spend_cap_usd // 0) as $cap
  | "\(.workflow_id // "—"): \($cur) \($st) (\($elapsed)m, $\($spent)/$\($cap), src=\(src_mix))"
' "$STATE"
