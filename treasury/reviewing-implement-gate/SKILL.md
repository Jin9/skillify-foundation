---
name: reviewing-implement-gate
description: >
  Walk a researcher through the five-question check before approving the
  implement gate in an agent-scaffold workflow. Reads .agent/stages/plan.md
  and .agent/stages/critique.md, confirms IMPLEMENT_SANDBOXED for
  credential-adjacent repos, sanity-checks the spend cap, and returns an
  approve/reject recommendation with the exact just command for the user to
  type. Use when the user says "should I approve implement", "review the
  implement plan", "is this safe to approve", "evaluate the gate", "five
  questions check", "ready to approve". Does NOT run just approve or just
  reject — only the user types those. Do NOT use for stages other than
  implement, for kicking off a workflow (use orchestrating-agent-scaffold),
  or for skills outside the agent-scaffold ecosystem.
---

# Reviewing the implement gate

## Purpose

The `implement` stage is the only gate in the default scaffold pipeline that
mutates source code. PLAYBOOK §10 demands a five-question check before any
human types `just approve implement`. This skill runs that check on the
agent's side, surfaces blockers, and never types the approve/reject command
itself.

## When to use this skill

- The agent-scaffold orchestrator surfaced an `implement` gate.
- The user asked "should I approve implement?", "is this plan safe?", or
  pasted a plan/critique pair for review.
- The user wants a written second opinion before approving.

Do NOT use this skill to:
- Approve, reject, or otherwise mutate `state.json`. The orchestrator skill
  does that *only at user instruction*; this skill does not call `just` at
  all beyond reads.
- Re-evaluate a gate that has already been approved. Once approved, this
  skill is a no-op.
- Run a generic security review. Hand off to `reviewing-software-security`
  for that. This skill consumes the critique that the security reviewer (or
  the scaffold's critique stage) already produced.

## Universal preamble

1. Confirm `.agent/state.json` has `approval_pending == "implement"` and
   `stages.implement.status == "blocked"`. If not, tell the user the gate
   is not pending and stop.
2. Read in order:
   - `.agent/stages/plan.md` — the implementation plan.
   - `.agent/stages/critique.md` — the critique (P1/P2/P3 issues).
   - `.agent/state.json` — `spend_cap_usd`, `total_cost_usd`,
     `stages.implement.cost_estimate` if set.
   - `IMPLEMENT_SANDBOXED` env var (or the launch command in `runlog.md`).
3. If any of the first two files is missing or empty, stop and ask the user
   whether the prior stage actually ran. Do not fabricate the missing
   content.

## Core workflow

Walk the five questions in order. Each must end in **yes / no / unsure**.
Any **no** or **unsure** is a blocker; surface it and recommend reject.

### Question 1 — Did I read `plan.md` end-to-end on a screen wider than my phone?

This is asked *of the human*, not of the agent. The agent's role: report the
plan length, summarize section headers, and ask the user to confirm they have
read it. If the plan is > 500 lines, suggest the user read it before
proceeding.

### Question 2 — Did the critique surface any blocker the plan didn't address?

Cross-walk:
- Extract every P1 (critical) and P2 (high) issue from `critique.md`.
- For each, search `plan.md` for an explicit mitigation.
- Mark each issue: `addressed` / `partial` / `unaddressed` / `accepted`.
- Any `unaddressed` P1 → reject. Any `unaddressed` P2 → recommend reject
  unless the user explicitly accepts.

### Question 3 — Is the cap tight enough that a runaway implement won't burn $50?

Read `state.json.spend_cap_usd` and `total_cost_usd`. The remaining budget
is `cap - total`. Heuristics:
- If remaining > $20 → flag as loose; suggest the user lower the cap or
  scope the goal tighter.
- If remaining < $0.50 → flag as too tight to even start; suggest
  `just raise-cap` (user types it).
- Otherwise → ok.

### Question 4 — Is `IMPLEMENT_SANDBOXED=1` set if this repo touches credentials/data?

Detect signals (any one triggers sandbox-required):
- `*.env`, `secrets/`, `config/credentials*` paths in plan.
- Mentions of: KYC, KMS, Vault, IAM, mTLS, JWT signing key, partner-bank,
  PII, PAN, NIK, account number, customer data.
- Plan touches `database/` migrations or `terraform/`.
- Repo is fintech-, lending-, or banking-tagged in `README.md` or
  `CLAUDE.md`.

If sandbox-required and the launch was unsandboxed, recommend reject and
tell the user to relaunch with `IMPLEMENT_SANDBOXED=1`.

### Question 5 — Am I ready to babysit the implement stage live (or accept that it'll happen unattended)?

Surface the user's answer back to them. If they're not ready and the cap
is loose (Q3 flag), recommend reject; if cap is tight, the loss is
bounded.

## Output format

Always emit `templates/gate-review.md` filled — the template is the single
source of the output shape: an **Inputs** block, the **P1/P2 cross-walk
table** (every P1 accounted for), the **five-question table**, the
recommendation, and the exact approve/reject commands. The block ends with
the exact command the user must type — never run it from the skill.

## Constraints

- DO NOT run `just approve` or `just reject` from this skill. Always print
  the command for the user to type.
- DO NOT skip the cross-walk in Q2. Auto-approval based on "the critique
  looked fine" is the failure mode this skill exists to prevent.
- DO NOT mark a P1 as `addressed` based on inference. Require explicit
  text in `plan.md` describing the mitigation.
- DO NOT proceed if `plan.md` or `critique.md` is missing — ask the user to
  re-run the prior stage.
- DO NOT extend scope to other gates. This skill is implement-only.

## Validation gate

Before printing the recommendation, verify:

1. The recommendation matches the worst answer across questions: any **no**
   or **unsure** in Q2/Q4 → reject; loose-cap or unsandboxed-required → reject.
2. Every P1 from `critique.md` is accounted for in the cross-walk.
3. The exact `just` command is printed and not run.

## References

| Need | Reference |
|---|---|
| The five questions in source form | `references/five-questions.md` |
| Sandbox-required signal list | `references/sandbox-signals.md` |
| Cross-walking critique to plan | `references/cross-walk.md` |

## Templates

- `templates/gate-review.md` — fillable review block.
