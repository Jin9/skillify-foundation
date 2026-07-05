---
name: discord-harness-extend
description: >-
  Add or modify one tool in this Discord LLM harness while preserving the gate-lives-outside-the-model
  invariant. Use when the user says "add a tool to the harness", "register a new @tool", "add a
  read-only or write tool", "wire a new capability behind the approval gate", or "modify a tool's
  safety class". Guides the full contract: write the @tool handler
  in tools/, re-validate side effects inside the handler (defense in depth), classify the tool in
  config.yaml as ALLOW/CONFIRM/DENY (unlisted fails closed to DENY), add a test_harness.py scenario,
  run the suite, and verify the audit chain. Do NOT use to run, monitor, or deploy the bot (use
  discord-harness-operate), to redesign the policy engine or approval mechanism, or for general Python
  refactoring.
---

# discord-harness-extend

Add or modify exactly one tool in the Discord LLM harness while preserving its core invariant: **the
authorization gate lives outside the model**. A tool's *behavior* and its *safety class* are registered
in two different places — miss either and the tool silently can't run, or runs unguarded.

## When to use

The agent has been asked to add a capability to the harness — a new `@tool`, a new operator-domain
tool, or a change to an existing tool's safety class. Read this whole file, then work the checklist as
todos.

Out of scope: running / monitoring / deploying the bot (use `discord-harness-operate`); redesigning the
policy engine or the approval gate; general refactoring.

## The contract in one sentence

Adding a tool = **write the handler** (in `tools/…` with `@tool(...)`) **and** **classify it** (in
`config.yaml` under `tools:` as ALLOW / CONFIRM / DENY) — both, or it can never run. Unlisted tools
fail closed to DENY (`policy.classify`, `policy.py:113`).

## Checklist

1. **Decide the safety class first.** Classify by reversibility and blast radius *before* writing code:
   ALLOW (reversible/sandboxed, auto-runs), CONFIRM (irreversible/control-plane, human approval
   button), DENY (never). See `references/safety-classes.md`.
2. **Write the handler.** Handlers go in `tools/db_ops.py` (the MySQL operator tools), or a new
   `tools/<x>.py` you also import in `tools/__init__.py`. Decorate with
   `@tool(name, description, parameters)` where `parameters` is a JSON Schema object. Signature:
   `def fn(args: dict, ctx: ToolContext) -> str` (sync or `async`). Start from
   `templates/tool-handler.py`; full contract in `references/tool-contract.md`.
3. **Re-validate inside the handler (defense in depth).** Even though the gate already approved a
   CONFIRM call, the handler MUST re-check every irreversible effect against the allowlist — never
   trust that the gate was correct. The MySQL re-validators in `policy.py`:
   - table identifier → `allowed_table(table)` (`policy.py:136`)
   - column identifier → `allowed_column(table, col)` (`policy.py:145`)
   - the single permitted database → `allowed_database(name)` (`policy.py:125`), enforced inside
     `db_ops._get_conn` so a write can never reach a database other than `db.database`.
   A guard raises `PermissionError`, which the loop reports as `BLOCKED by guard` (`agent.py:173`).
4. **Register a NEW module for import.** If you created a new `tools/<x>.py`, add `from . import <x>`
   to the import block in `tools/__init__.py:58`. The `db_ops` module is already imported — skip this
   when extending it.
5. **Classify in config.** Add `name: ALLOW|CONFIRM|DENY` under `tools:` in `config.yaml`. This is
   mandatory — an unlisted handler is silently DENY and is never even offered to the model
   (`agent.py:65`). Add any new allowlist entries the tool needs (`db.allowed_tables`,
   `db.allowed_columns`). There is a single config profile (`config.yaml`); `HARNESS_CONFIG` can point
   at an alternate policy file, but there is no second profile to mirror the class into.
6. **Never put authorization in the prompt.** Do not gate via `agent.SYSTEM_PROMPT` (`agent.py:29`) or
   the config `system_prompt`. The decision stays in `policy.py` + config; the prompt only steers
   tone/scope (`agent.py:79`).
7. **Add a test scenario.** Append a scripted-litellm scenario to `tests/test_harness.py` that asserts
   BOTH the side effect AND the policy outcome (auto-ran / gated-and-approved / declined / DENY /
   guard-blocked). Use the existing helpers, model on agent scenarios A–F. Start from
   `templates/test-scenario.py`; details in `references/test-recipe.md`.
8. **Run the suite.** `python tests/test_harness.py` — every line must print `OK` and it must end with
   `ALL HARNESS LOGIC TESTS PASSED`. No Discord token or API key is needed (litellm is mocked).
9. **Cross-check registration.** Run `scripts/check_registration.py` from the repo root (in the
   `.venv`): it prints the tool × class matrix for `config.yaml` and fails on an orphan (a `tools:`
   entry with no handler) or an invalid class string (which would silently become DENY). Confirm your
   new tool appears under the class you intended — not under DENY.
10. **Verify the audit chain.**
    `python -c "from audit import AuditLog; print(AuditLog('audit.jsonl').verify())"` → `True`.

## Hard rules (do not violate)

- **Both registrations or nothing.** Handler without a `config.yaml` class = dead tool (DENY). Class
  without a handler = orphan that is never offered.
- **Defense in depth is not optional.** A CONFIRM tool still re-validates in its handler. The gate and
  the handler guard are independent layers (`tools/db_ops.py:6`).
- **Keep `Optional[X]` / `Union[...]` as-is.** Do not churn pydantic model fields (`contracts.py:66`)
  or module-level type aliases (`tools/__init__.py:16`) to `X | None` — that form is intentional
  Python ≤3.9 back-compat and breaks pydantic field evaluation on old interpreters (the runtime target
  is now 3.13, but the `Optional`/`Union` form stays).
- **Secrets stay in `.env`; policy stays in `config.yaml`.** A new tool never reads secrets from yaml.

## Exit criteria

Handler written + re-validated, module imported (if new), classified in config, test scenario added,
`python tests/test_harness.py` green, `scripts/check_registration.py` clean, audit chain verifies.

## References

| Need | File |
|------|------|
| `@tool` decorator, ToolContext, schema, dispatch, files-to-touch table | `references/tool-contract.md` |
| ALLOW/CONFIRM/DENY decision + the defense-in-depth re-validators | `references/safety-classes.md` |
| test_harness.py mock helpers, agent scenarios A–F, how to add one | `references/test-recipe.md` |
| Handler skeletons (MySQL read + write) + config snippet | `templates/tool-handler.py` |
| Test scenario skeleton | `templates/test-scenario.py` |
| Registration cross-check (handlers vs config) | `scripts/check_registration.py` |
