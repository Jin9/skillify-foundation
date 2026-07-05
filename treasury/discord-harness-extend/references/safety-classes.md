# Safety classes — ALLOW / CONFIRM / DENY

The class IS the authorization decision, and it lives in config, enforced by `policy.classify()`
(`policy.py:113`), never in the prompt. `contracts.py:17`:

```python
class SafetyClass(str, Enum):
    ALLOW = "ALLOW"      # reversible / sandboxed → auto-execute
    CONFIRM = "CONFIRM"  # irreversible / control-plane → sync human approval
    DENY = "DENY"        # outside allowlist → refuse
```

## How to choose

Decide BEFORE writing the handler. Ask: if the model calls this with worst-case arguments, what is the
blast radius, and is it reversible?

| Signal | Class |
|--------|-------|
| Read-only, sandboxed to an allowlist, no external side effect | ALLOW |
| Mutates the database destructively, spends money, irreversible | CONFIRM |
| Should never run in this deployment | DENY (or just leave it unlisted) |

In-repo examples: the reads `mysql_get_row` / `mysql_query_rows` = ALLOW; the validated upserts
`mysql_insert_row` / `mysql_update_row` are ALLOW too (parameterized, allowlist-checked); the
destructive `mysql_delete_row` = CONFIRM (`config.yaml:49`).

## Fail closed

`classify()` returns DENY for any name not in `tools:` — including a typo in the class value (it is a
case-sensitive exact match; `allow` becomes DENY). So a handler you forgot to list, or mis-cased, is
silently unavailable. `scripts/check_registration.py` surfaces both.

```python
def classify(tool_name: str) -> SafetyClass:
    cfg = load_config()
    raw = (cfg.get("tools") or {}).get(tool_name)
    try:
        return SafetyClass(raw)
    except (ValueError, TypeError):
        return SafetyClass.DENY
```

## What CONFIRM actually does

When the loop hits a CONFIRM tool (`agent.py:151`) it moves to AWAITING_APPROVAL and calls
`approval_cb(call)`. In the bot that posts a Discord Approve/Deny button only an allowlisted approver
can resolve (`bot.py:109`). Declined → the tool does not run; the model is told to continue without it.

## Defense in depth — re-validate in the handler

CONFIRM means a human said yes to THIS call, not that the arguments are safe. The handler re-checks
against the allowlist regardless. The MySQL guards live in `policy.py` and raise `PermissionError`:

| Effect | Guard | Anchor |
|--------|-------|--------|
| MySQL table identifier | `allowed_table(table)` | `policy.py:136` |
| MySQL column identifier | `allowed_column(table, col)` | `policy.py:145` |
| MySQL SELECT projection | `table_projection(table)` | `policy.py:158` |
| the single permitted database | `allowed_database(name)` | `policy.py:125` |

`allowed_database` is enforced inside `db_ops._get_conn` so no op can reach a database other than
`db.database`; `allowed_table` / `allowed_column` re-check every identifier inside each `db_ops` handler
before the parameterized SQL is built.

The loop turns a raised `PermissionError` into `BLOCKED by guard: …` and audits a `guard_block` event
(`agent.py:171`). Agent scenario F (`tests/test_harness.py:241`) proves an APPROVED `mysql_delete_row`
on an UN-allowlisted table is still blocked by the handler guard.

## The invariant, restated

The model never decides whether a tool runs. Adding a tool must not move any authorization into
`agent.SYSTEM_PROMPT` or the config `system_prompt` — those steer behavior, not permissions
(`agent.py:79`). If you catch yourself writing "the model should refuse to…", stop: encode it as a
class or an allowlist instead.
