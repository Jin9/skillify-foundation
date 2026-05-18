# Boundaries and squad-roles contract

## The three-tier boundaries model

A widely-cited structural pattern defines three tiers of agent autonomy:
**actions the agent may always do, actions that require human approval
first, and actions that are explicitly prohibited**. This matters most in
agentic squads where the agent operates with elevated permissions (file
writes, git commits, API calls). The boundaries section of AGENTS.md MUST
be written as exactly these three tiers.

```
## Boundaries

### Always allowed (no approval)
- Read any file in the repo.
- Run the test/lint commands listed above.
- Create a feature branch and open a draft PR.

### Requires human approval first
- Schema migrations — propose in the PR description and wait; do not apply.
- Dependency major-version upgrades — open an issue and wait instead.
- Any change under `src/db/migrations/` — describe the intended change and
  request the module owner's sign-off instead of editing.

### Explicitly prohibited (with the allowed alternative)
- Do not push directly to `main` — open a PR instead.
- Do not modify CI workflow files to make a check pass — fix the code
  instead, or escalate to a human if the check is wrong.
- Do not commit secrets or credentials — reference them via the existing
  config mechanism instead.
```

Rule: **every prohibited-tier line pairs the "don't" with a concrete "do"**
on the same line or the adjacent bullet. A prohibition with no alternative
is the over-specification failure mode that makes agents conservative and
incomplete. `scripts/agents_md_gate.py` fails the candidate on any
unpaired prohibition.

The requires-human-approval tier implements a "human-on-the-loop" design:
agents work autonomously for routine decisions while escalating high-stakes
actions — schema migrations, dependency major upgrades, production deploys
— for explicit human sign-off.

## The optional squad-roles contract block

The standard six-section template was designed for a single agent. When
multiple agents with different scopes share a repository, AGENTS.md must
additionally express inter-agent contracts: who hands off to whom, what
triggers human escalation, and which agents have write access to which
paths. Without these, agents either duplicate effort or step on each
other's work (the multi-agent context-duplication failure also applies —
see `authorship-and-budget.md`).

Add the squad-roles block ONLY for multi-agent repos. Each role is one
line with four fields:

| Field | Meaning |
|---|---|
| `role` | the agent's name in the squad |
| `owned paths` | the path globs this role may write |
| `handoff-target` | which role it hands completed work to |
| `human-escalation triggers` | the actions that require the human gate |

```
## Squad Roles
- coordinator: owns nothing; routes tasks, spawns specialists, escalates
  to the human gate.
- backend-agent: owns `src/api/**`, `src/core/**`; hands off to
  tester-agent; escalates on schema migration.
- frontend-agent: owns `src/ui/**`; hands off to tester-agent; escalates
  on a public-API contract change.
- tester-agent: owns `tests/**`; hands off to coordinator; escalates on a
  flaky-test quarantine.
- human-gate: required for schema migrations, dependency major upgrades,
  production deploys.
```

## The handoff-schema boundary

The squad-roles block names *handoff-target* but **does not define the
handoff message schema** (payload shape, status codes, retry/contract
semantics). That schema is owned by `multi-agent-handoff-architect`. In
AGENTS.md, point to it (e.g. "handoff payload schema: see the squad's
handoff contract") and stop. Defining it here would both duplicate that
skill and blow the line budget.
