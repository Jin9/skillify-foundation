---
name: authoring-scaffold-profile
description: >
  Author or modify an agent-scaffold profile file at profiles/NAME.sh. Sets
  STAGES, GATED_STAGES, and per-stage AGENT / MODEL / PROMPT_PREFIX env
  vars. Validates that every STAGE has a runner script, that GATED_STAGES
  is a subset of STAGES, that AGENT names are recognized by the LiteLLM
  config, and that prompt-prefix files exist when referenced. Use when the
  user says "new profile for X domain", "customize the scaffold for our
  squad", "create agent profile", "add a stage to the profile", "change
  the model for the plan stage", or "fork code.sh for our team". Refuses
  to overwrite an existing profile without confirmation. Do NOT use for
  editing .agent/config.sh in-place (that is the per-repo defaults file —
  edit cautiously and via direct PR), for editing AGENTS.md/CLAUDE.md, or
  for writing prompt bodies (use drafting-stage-prompt for those).
---

# Authoring an agent-scaffold profile

## Purpose

A profile in `profiles/<name>.sh` is sourced by the dispatcher to override
which stages run, which are gated, and which model handles each. This
skill writes those files consistently and validates them against the
scaffold's runner inventory and LiteLLM model config.

## When to use this skill

- The user says "new profile for X domain", "customize for our squad",
  "fork code.sh", or "make a literature variant".
- A postmortem action calls for changing a default in a profile.
- A research squad is onboarding a new domain that doesn't fit `code`,
  `literature`, or `dataset`.

Do NOT use this skill to:
- Edit `.agent/config.sh` directly. That is the per-repo defaults file
  and changes there should be a deliberate PR with team-lead review.
- Author prompt bodies. Profiles only set `*_PROMPT_PREFIX` strings;
  full prompts live in `prompts/library/` (use `drafting-stage-prompt`).
- Modify scaffold runner scripts. Profiles do not own runner code.

## Universal preamble

1. Confirm the working directory is an agent-scaffold checkout
   (`profiles/`, the scaffold's stage-runners directory, `.agent/config.sh` exist).
2. Ask for the profile slug (kebab-case, ≤ 32 chars, matches
   `^[a-z][a-z0-9-]{0,31}$`).
3. Compute the path: `profiles/<slug>.sh`. Refuse to overwrite without
   explicit user confirmation; offer `<slug>-v2.sh` instead.
4. Read existing profiles to mirror their style:
   - `profiles/code.sh` (default, full pipeline).
   - `profiles/literature.sh` (skips implement/test, prompt-prefixed).
   - `profiles/dataset.sh` (skips test).

## Core workflow

### 1 — Determine STAGES

Pick from built-in stages. Each name must match a runner script that the
scaffold ships under its stage-runners directory.

| Stage | Runner present in default scaffold | Skip when |
|---|---|---|
| `research` | yes | almost never |
| `plan` | yes | almost never |
| `critique` | yes | rarely |
| `implement` | yes (+ `implement-sandboxed.sh`) | literature-style profiles |
| `review` | yes | almost never |
| `test` | yes (auto-detects go/npm/pytest/cargo/bats) | dataset/literature |

Custom stages require a matching runner; this skill does **not** create
runners — escalate to the user with "you'll need to write a runner
under the scaffold's stage-runners directory first" and stop.

### 2 — Determine GATED_STAGES

Default: `implement` (unless the profile skips implement).

Heuristics:

- If `STAGES` contains `implement`, `GATED_STAGES` must contain `implement`.
- If the profile is for credential-adjacent / fintech / lending work,
  consider also gating `plan` so the user reviews the plan before any
  research budget is spent on critique.
- For literature profiles (no implement), `GATED_STAGES=""` is normal.

### 3 — Set per-stage AGENT/MODEL

Mirror the default scaffold's bindings unless the user has a reason to
deviate:

| Stage | Default AGENT | Default MODEL |
|---|---|---|
| research | gemini | gemini-research |
| plan | codex | codex-plan |
| critique | claude | claude-critique |
| implement | codex | codex-implement |
| review | claude | claude-review |
| test | (none — local) | (none — local) |

Validate that AGENT names exist in the user's LiteLLM config
(`docker/litellm/config.yaml`). If unknown, surface a warning and ask
the user to confirm before writing.

### 4 — Set PROMPT_PREFIX where helpful

Prefixes are short steering strings prepended at runtime — not full
prompts. Examples from `literature.sh`:

```
RESEARCH_PROMPT_PREFIX='Optimize for synthesis across many sources. Cite paper titles. '
PLAN_PROMPT_PREFIX='Treat the "plan" as a draft outline of the synthesis document, not code steps. '
```

Rules:

- Prefixes ≤ 200 chars. Long steering belongs in
  `prompts/library/<stage>/<topic>.md`.
- End the prefix with a trailing space so the runner concatenates cleanly.
- No PII, no real customer names, no credentials.
- No `claude` or `anthropic` brand names — prefixes should drive
  behavior, not name a vendor.

### 5 — Render the profile

Use `templates/profile.sh`. The shape is bash that the dispatcher
sources; do not change the variable names — the runners read them.

### 6 — Validate the file

Before printing the path, run these checks (described in
`references/validation-checklist.md`):

1. Every stage in `STAGES` has a runner script.
2. `GATED_STAGES` is a subset of `STAGES`.
3. Every AGENT name appears in LiteLLM's config (or user accepted a
   warning).
4. Variable names match the canonical set:
   `<STAGE>_AGENT`, `<STAGE>_MODEL`, `<STAGE>_PROMPT_PREFIX`.
5. The file passes `bash -n profiles/<slug>.sh` (syntax check).

### 7 — Surface to the user

Print the path and the launch command:

```bash
WORKFLOW_PROFILE=<slug> just workflow <id> "<goal>" <cap>
```

Suggest `git add profiles/<slug>.sh` and a PR.

## Output format

One bash script at `profiles/<slug>.sh`, executable. The file is the
skill's only artifact — no companion docs.

## Constraints

- DO NOT create runner scripts. If a custom stage is requested, stop
  and tell the user they need to author the runner first.
- DO NOT edit `.agent/config.sh`. That is the per-repo default; profiles
  layer on top.
- DO NOT inline full prompt bodies. Profiles set short prefixes; full
  prompts live in `prompts/library/` via `drafting-stage-prompt`.
- DO NOT skip `bash -n` syntax check before declaring done.
- DO NOT include `claude` or `anthropic` in the profile slug or prefix
  text.
- DO NOT overwrite an existing profile without explicit confirmation.

## Validation gate

Before the file ships:

1. `bash -n profiles/<slug>.sh` exits 0.
2. Every `STAGES` entry has a corresponding runner.
3. `GATED_STAGES ⊆ STAGES`.
4. AGENT/MODEL names match the LiteLLM config (or warning accepted).
5. No real PII, credentials, or customer names in any prefix.
6. File mode allows the dispatcher to source it (no need for `+x`,
   bash sources without exec bit; but file must exist and be readable).

## References

| Need | Reference |
|---|---|
| Built-in stages and runners | `references/stage-inventory.md` |
| Validation checklist (preflight) | `references/validation-checklist.md` |
| Prefix style and length rules | `references/prompt-prefixes.md` |

## Templates

- `templates/profile.sh` — fillable profile skeleton.
