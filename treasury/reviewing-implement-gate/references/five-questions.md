# The five questions (PLAYBOOK §10)

Before typing `just approve implement`:

1. Did I read `plan.md` end to end on a screen wider than my phone?
2. Did the critique surface any blocker that the plan didn't address?
3. Is the cap tight enough that a runaway implement won't burn $50?
4. Is `IMPLEMENT_SANDBOXED=1` set if this repo touches credentials/data?
5. Am I ready to babysit the implement stage live (or accept that it'll
   happen unattended)?

If any answer is **no** or **I don't know** — `just reject implement
"<reason>"` and adjust.

## Why each question matters

| # | Failure mode it prevents |
|---|---|
| 1 | Approving on a phone notification without reading the plan. The single most common postmortem cause. |
| 2 | Critique surfaces a P1 (e.g., missing tx boundary) and the plan ignored it; implement ships the bug. |
| 3 | A loose cap turns a misfiring implement into a $50 incident before anyone notices. |
| 4 | Implement runs without proxy enforcement on a repo with credentials; egress to attacker-controlled domains is possible. |
| 5 | Unattended implement on a tight cap is fine; unattended implement on a loose cap with a real bug burns the budget silently. |

## Hard rejections

These conditions force a `reject` regardless of user judgement:

- Any **unaddressed** P1 issue from the critique.
- A repo that triggers any sandbox-required signal but launched without
  `IMPLEMENT_SANDBOXED=1`.
- `state.json.approval_pending != "implement"` (the gate is not actually
  pending).
- `critique.md` missing or empty.
