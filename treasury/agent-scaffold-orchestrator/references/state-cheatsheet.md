# `.agent/state.json` cheatsheet

Mirror of the core fields the orchestrator skill reads. Authoritative source:
`docs/STATE_SCHEMA.md` in the scaffold checkout.

## Top-level

| Field | Type | Read for |
|---|---|---|
| `schema_version` | int | Drift detection. If unfamiliar, run `just doctor`. |
| `workflow_id` | string | `^[a-z][a-z0-9-]{0,31}$`. |
| `goal` | string | One sentence, ≤ 500 chars, no newlines. |
| `current_stage` | string\|null | Most recently dispatched stage. |
| `started_at` / `started_epoch` | iso / int | Always together. |
| `approval_pending` | string\|null | Stage waiting on a human gate. |
| `total_cost_usd` | number | Sum of completed stages' `cost_usd`. |
| `spend_cap_usd` | number | Hard cap; dispatcher exits 2 when total ≥ cap. |

## Per-stage

`.stages.<name>` is one of:

```jsonc
{
  "status": "pending|running|blocked|done|failed|rejected|aborted",
  "agent": "gemini|codex|claude|...",
  "model": "gemini-research|codex-plan|...",
  "started_at": "2026-...", "started_epoch": 1746514800,
  "ended_at":   "2026-...", "ended_epoch":   1746515520,
  "duration_s": 720,
  "cost_usd": 0.34,
  "cost_source": "real|estimate",
  "pid": 12345,                     // background mode only
  "zellij_tab": "stage:critique",   // zellij mode only
  "approved": true,
  "approved_at": "2026-...",
  "approved_by": "username",
  "error": "...",                   // set on failed
  "rejection_reason": "...",        // set on rejected
  "block_reason": "..."             // set on blocked (mid-stage)
}
```

## Reads (never lock)

```bash
# canonical status line
just status

# raw fetch
jq -r '.current_stage' .agent/state.json

# elapsed seconds
jq -r '
  .stages[.current_stage] as $s
  | (now | floor) - ($s.started_epoch // 0)
' .agent/state.json

# cap delta
jq -r '"\(.total_cost_usd)/\(.spend_cap_usd)"' .agent/state.json
```

## Writes

The orchestrator skill never writes. State mutations go through
`scripts/dispatch.sh` via `just`. If you find yourself wanting to write,
stop and surface the wanted change to the user.

## Backup discipline

Every mutation snapshots to `.agent/state.json.bak`. If `state.json` looks
corrupted (jq parse error, missing fields), prefer:

```bash
diff .agent/state.json.bak .agent/state.json
```

…before suggesting any recovery. Often the backup is intact and the user
can `cp .agent/state.json.bak .agent/state.json && just resume`.
