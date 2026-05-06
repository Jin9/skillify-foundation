# Failure triage

When a stage flips to `failed`, classify the failure before suggesting any
recovery action.

## Failure classes

| Class | Signal in log | Suggested next step |
|---|---|---|
| Stale-PID / runner died | `runner died (stale pid …)`; pid in state but not alive | `just resume` (dispatcher auto-recovers stale-running). |
| Cap tripped | `spend cap reached: $X / $Y — halting`; rc=2 | Run `just llm-spend 1` for cost truth. Cap raise needs explicit user request. |
| Model API rate-limit / network blip | `429`, `503`, `connection reset`, `timeout` | `just resume` after a brief wait. If repeats, escalate. |
| LiteLLM unreachable | `connection refused 127.0.0.1:4000` | `curl http://127.0.0.1:4000/health/readiness`; suggest `just llm-up`. |
| Sandbox egress denied | tinyproxy `403`, host not in allowlist | Hand off to `configuring-sandbox-allowlist`. Do not edit the filter from this skill. |
| Bad prompt / model misbehaved | model output empty, off-topic, or refused | Hand off to `drafting-stage-prompt` to update `prompts/library/<stage>/<topic>.md`. |
| Real bug in scaffold or runner | shell error, jq parse error, missing file | Read the runner script, surface the failing line; user fixes and `just resume`. |
| Approval log integrity failure | `just approvals-verify` shows tampered entries | Treat as security incident; escalate to team-lead. |

## Triage workflow

1. Cat the redacted tail:
   ```bash
   tail -n 80 .agent/stages/<stage>.log
   ```
2. Match the error against the table above. If multiple match, prefer the
   last error in the log.
3. Read `.agent/runlog.md`'s last 30 lines for context (was this the first
   failure? has resume been tried?).
4. Suggest **one** next step. Do not chain "try X then Y then Z".
5. If the failure spent > $1 (check `cost_usd` for the failed stage),
   surface the postmortem trigger and offer to hand off to
   `agent-workflow-postmortem`.

## Forbidden recovery actions

- Auto-running `just resume` without user say-so when cost > $0.50 has been
  spent on the failure.
- Editing `.agent/state.json` to flip `failed` → `pending`.
- Killing pids manually with `kill -9` instead of `just abort`.
- Skipping `just doctor` when the failure looks environmental.

## Escalation triggers

Stop and escalate to the user if:

- The same stage fails twice with the same root cause.
- Cap is tripped and the cost looks anomalously high (>2× the historical
  median for this profile).
- The approval log integrity check fails.
- Sandbox enforcement appears to be bypassed (a non-allowlisted host
  responded successfully).
