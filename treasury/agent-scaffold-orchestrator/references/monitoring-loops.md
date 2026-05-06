# Monitoring cadence and summary shapes

## When to read state

- On every user-facing response that mentions the workflow.
- After any `just <verb>` invocation that mutates state (init, next, abort,
  resume, raise-cap, approve, reject).
- Never on a fixed timer. Read on demand.

## Canonical status line

Every response leads with:

```
<workflow_id>: <current_stage> <status> (<elapsed>m, $<spent>/$<cap>, src=<real|estimate>)
```

Generate via `scripts/state_summary.sh` (see Templates).

## State fields you will read most

| Field | Purpose |
|---|---|
| `.workflow_id` | Identifier for log lines, ntfy, follow-ups. |
| `.current_stage` | Which stage entry to read first. |
| `.stages[$s].status` | `pending` / `running` / `blocked` / `done` / `failed` / `rejected` / `aborted`. |
| `.stages[$s].started_epoch` | Compute elapsed seconds against `epoch_now`. |
| `.stages[$s].cost_usd`, `.cost_source` | Per-stage spend; `real` from LiteLLM, `estimate` is the fallback. |
| `.stages[$s].pid` | Background-mode runner pid. Use `kill -0` to test liveness. |
| `.stages[$s].zellij_tab` | If set, suggest `just goto <stage>` for live view. |
| `.approval_pending` | If non-null, a gate is waiting on the user. |
| `.total_cost_usd`, `.spend_cap_usd` | Cap math. Compare with `bc` or jq. |

## Tail-the-log shape

When status alone isn't enough, append the redacted tail:

```bash
tail -n 20 .agent/stages/<current>.log
```

Logs are redacted at write time, so it is safe to surface inline. Do not
copy raw lines longer than 200 chars without trimming.

## Approval-pending response

When `.approval_pending` is non-null:

```
🔒 Gate: <stage>
   Read:    just show-stage <prior>
   Approve: just approve <stage>          (you must type this)
   Reject:  just reject <stage> "<reason>"
```

For the `implement` gate, hand off to `reviewing-implement-gate` instead of
giving the user the raw approve/reject choice.

## Failure-pending response

When the most recent stage is `failed`:

1. Lead with the canonical status line.
2. Quote the last 5–10 lines of `.agent/stages/<stage>.log`.
3. Do **not** suggest `just resume` automatically — describe the failure
   class first (transient / prompt / real bug), then suggest one fix path.

## Throttling

The user typing "status" twice in 30 seconds gets one fresh read; do not
re-read between if `state.json`'s mtime is unchanged. Surface "(unchanged
since N seconds ago)" so the user knows to be patient.
