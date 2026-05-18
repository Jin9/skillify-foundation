<!--
  AGENTS.md SKELETON — fill in, then delete every comment before delivering.
  LINE BUDGET: target 100-150 lines, HARD CAP 200 lines. Above 150 the gate
  WARNs; above 200 it FAILs. If it will not fit, do not raise the ceiling —
  move more rules into the curation-checklist deferral table instead.
  DEFERRAL FILTER: if a rule is enforceable by a lint rule, a CI check, or a
  branch-protection rule, it MUST NOT appear here — log it to the checklist.
  PAIR EVERY DON'T WITH A DO: no prohibition without a concrete alternative.
  Exactly the six sections below (+ optional Squad Roles). No seventh.
-->

# AGENTS.md

<!-- 1-2 sentences max: what this repo is, just enough for an agent to orient. -->
One-line project summary. Stack and versions: <stack@version>.

## Commands

<!-- Exact CLI with flags. Not "run the tests" — the literal command. -->
- Build:  `<exact build command>`
- Test:   `<exact test command with flags>`
- Lint:   `<exact lint + format-check command>`

## Testing

<!-- Single test, full suite, and what "passing" means here. -->
- Single test: `<command for one test>`
- Full suite (must be green before PR): `<command>`
- "Passing" means: <suite green + coverage gate; gate enforced in CI, not here>.

## Project structure

<!-- 3-5 bullets. Orientation only — not a file tree, not architecture history. -->
- `<dir>/` — <responsibility>
- `<dir>/` — <responsibility>
- `<dir>/` — <responsibility>

## Code style

<!-- ONE snippet + a pointer to the lint config. Do NOT restate lint rules. -->
Style is enforced by `<lint/format config file>`. Shape new code to match:

    <one short representative code snippet in the house style>

## Git workflow

<!-- Only the non-enforceable parts. Required checks / approvals are
     branch-protection rules — defer them, do not write them here. -->
- Branch naming: `<pattern>`. Commit convention: `<convention>`.
- PR body: <what changed + how verified>. Do not self-merge — request
  review from the module owner instead.

## Boundaries

<!-- Exactly three tiers. Every prohibited line pairs a don't with a do. -->

### Always allowed (no approval)
- Read any file; run the test/lint commands above; open a draft PR.

### Requires human approval first
- Schema migrations — propose in the PR description and wait; do not apply.
- Dependency major-version upgrades — open an issue and wait instead.

### Explicitly prohibited (with the allowed alternative)
- Do not push directly to `main` — open a PR instead.
- Do not edit CI workflow files to force a check green — fix the code
  instead, or escalate to a human if the check itself is wrong.
- Do not commit secrets — reference them via the existing config mechanism
  instead.

<!-- OPTIONAL: include the block below ONLY for multi-agent squad repos.
     Delete it entirely for a single-agent repo. The handoff message schema
     is NOT defined here — only the handoff-target is named; point to the
     squad's handoff contract owned by multi-agent-handoff-architect. -->

## Squad Roles
- coordinator: owns nothing; routes tasks, escalates to the human gate.
- <role>: owns `<paths>`; hands off to `<role>`; escalates on <trigger>.
- <role>: owns `<paths>`; hands off to `<role>`; escalates on <trigger>.
- human-gate: required for <schema migrations, major upgrades, prod deploys>.
- Handoff payload schema: see the squad's handoff contract (not defined here).
