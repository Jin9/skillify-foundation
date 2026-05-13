# Trigger conditions

A postmortem is **required** within 48 hours of any of these.

## 1. Failed stage with > $1 spent

```bash
jq -r '
  .stages
  | to_entries
  | map(select(.value.status == "failed" and (.value.cost_usd // 0) > 1))
  | .[].key
' .agent/state.json
```

If the output is non-empty, this trigger fired.

## 2. Rejected implement gate

```bash
jq -r '.stages.implement.status' .agent/state.json
```

If `rejected`, this trigger fired. Read
`.agent/stages.implement.rejection_reason` and `.agent/approvals.md` for
the recorded reason.

## 3. Cap overrun

The dispatcher exits `rc=2` and emits `spend cap reached: $X / $Y`.
Confirm via:

```bash
jq -r '"\(.total_cost_usd) >= \(.spend_cap_usd)"' .agent/state.json
```

If `total_cost_usd >= spend_cap_usd`, this trigger fired regardless of
whether the cap was raised after the fact.

## 4. Sandbox egress failure

The dispatcher's pre-implement check verifies sandbox enforcement by
attempting to fetch a non-allowlisted host. If that curl unexpectedly
succeeded, the implement runner aborts and ntfy reports it. Look for:

- Recent `.agent/runlog.md` entry mentioning sandbox enforcement failure.
- `docker compose -f docker/sandbox-compose.yml logs proxy` showing a
  successful proxy_pass for a host that should have been blocked.

## Verification before drafting

| Trigger | Verification command |
|---|---|
| Failed > $1 | jq filter above |
| Rejected gate | jq filter above + `cat .agent/approvals.md \| tail -20` |
| Cap overrun | jq filter above + `just llm-spend 24` |
| Sandbox failure | scaffold's egress test logs |

If verification fails (e.g., user says "rejected gate" but state shows
`done`), stop and ask the user to clarify. Do not draft a postmortem on
a phantom trigger.

## Multiple triggers

A workflow can fire multiple triggers (e.g., rejected gate + cap
overrun). Author one postmortem covering all of them; the **Severity**
header lists each as a comma-separated value.
