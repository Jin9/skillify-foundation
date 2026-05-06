---
name: drafting-stage-prompt
description: >
  Curate stage prompts for the agent-scaffold prompt library at
  prompts/library/STAGE/TOPIC.md. Captures the Friday prompt-review
  ritual: a stage prompt with a why-it-works paragraph, success metric,
  failure mode, model affinity, and an invocation example. Use when the
  user says "add this prompt to the library", "save this prompt", "curate
  stage prompt", "library entry", "prompt-review ritual", "this prompt
  worked great", or after a workflow finishes and the user wants to
  preserve a prompt that surprised them. Refuses to overwrite without
  confirmation. Do NOT use for editing AGENTS.md / CLAUDE.md, for writing
  prompts outside prompts/library/, or for retroactively rewriting historic
  workflow runs (use agent-workflow-postmortem for that).
---

# Drafting a stage prompt

## Purpose

The Friday prompt-review ritual moves prompts that won (or lost) into
`prompts/library/<stage>/<topic>.md`. This skill writes those entries
consistently so they survive future Fridays without re-litigating tone or
structure.

## When to use this skill

- User says "save this prompt", "add to library", "prompt-review entry".
- A workflow just shipped a clear win or loss the user wants captured.
- The user is preparing for the Friday prompt-review meeting and wants to
  draft entries in advance.

Do NOT use this skill to:
- Rewrite repo-policy files (`AGENTS.md`, `CLAUDE.md`, `.github/`).
- Save prompts outside `prompts/library/<stage>/`. Profiles use
  `*_PROMPT_PREFIX` env vars; profile authoring belongs to
  `authoring-scaffold-profile`.
- Author postmortems. A prompt entry is not a postmortem.

## Universal preamble

1. Confirm the working directory is an agent-scaffold checkout
   (`prompts/library/` exists, `profiles/` exists). If not, stop.
2. Ask for inputs if missing — never invent:
   - **Stage** — one of `research`, `plan`, `critique`, `implement`,
     `review`, `test`, or a custom stage from a profile.
   - **Topic slug** — kebab-case, ≤ 40 chars, e.g., `jwt-validation`.
   - **Prompt body** — the actual prompt the user wants to save.
   - **Why it worked / failed** — one paragraph.
   - **Outcome class** — `win` or `loss`.
3. Compute the target path:
   `prompts/library/<stage>/<topic>.md`. Refuse to overwrite if it exists
   without explicit user confirmation; offer `<topic>-v2.md` as an
   alternative.

## Modes

| Mode | Trigger | What this skill does |
|---|---|---|
| Win | "save this winning prompt", "this worked great" | Write a `win` entry — full prompt, why-it-worked, success metric, model affinity. |
| Loss | "this prompt was a disaster", "what we don't do" | Write a one-line note in the squad's `prompts/library/_anti-patterns.md` (creates if missing). Does **not** create a full entry. |
| Update | "update the library entry for <topic>" | Read existing entry, append a `## Updated YYYY-MM-DD` block with new findings. |

The Win mode is the default; the Friday ritual produces wins more often
than full losses (losses become one-liners, not entries).

## Win-mode workflow

1. **Validate stage and topic.** Stage must exist in `profiles/*.sh`'s
   STAGES (or be a built-in scaffold stage). Topic slug must match
   `^[a-z0-9-]{1,40}$`.
2. **Draft the entry** using `templates/stage-prompt.md`. Sections:
   - Frontmatter: stage, topic, model affinity, last-validated date.
   - Prompt body — copy verbatim from user; do not edit.
   - Why it works — one paragraph.
   - Success metric — one line; the observable signal that proves it
     worked (e.g., "critique surfaced 3 P1 issues we hadn't seen").
   - Failure mode — one line; the most likely way it stops working.
   - Model affinity — which model in the scaffold (gemini/codex/claude)
     this prompt was written for.
   - Invocation example — short snippet showing how the runner uses it.
3. **Confirm path.** Print the destination. If the file exists, print a
   diff and ask the user: overwrite, save as `-v2`, or abort.
4. **Write.** Use the Write tool. After writing, print the path and
   suggest `git add prompts/library/<stage>/<topic>.md`.

## Loss-mode workflow

1. Append one line to `prompts/library/_anti-patterns.md` (create if
   missing) in the form:
   ```
   - <date> · <stage>/<topic-slug> — <one-line why it failed>
   ```
2. Do not write a full entry. The point is to keep losses cheap so the
   squad doesn't waste a future workflow rediscovering them.

## Output format

Win-mode produces `prompts/library/<stage>/<topic>.md` matching
`templates/stage-prompt.md`. Loss-mode appends one line to
`prompts/library/_anti-patterns.md`. Update-mode reads then appends a
dated update block — never silently rewrites earlier sections.

## Constraints

- DO NOT edit the user's prompt body. Save it verbatim.
- DO NOT save prompts that contain real PII, real credentials, or real
  customer names. Refuse and ask the user to redact.
- DO NOT overwrite an existing entry without explicit confirmation. Offer
  `-v2` instead.
- DO NOT save outside `prompts/library/<stage>/`. The library shape is
  load-bearing for the Friday review.
- DO NOT include `claude` or `anthropic` in the topic slug — that bias is
  exactly what the platform-neutral library is meant to avoid.
- DO NOT mark a prompt as "validated" without a success metric the user
  named. Made-up metrics defeat the ritual.

## Validation gate

Before writing:

1. The target path resolves under `prompts/library/<stage>/`.
2. Stage exists in a profile or is a scaffold built-in.
3. Topic slug matches `^[a-z0-9-]{1,40}$`.
4. The entry has a non-empty prompt body, why-it-works paragraph, success
   metric, and failure mode.
5. No real PII or credentials in the body (string-match against common
   patterns; ask user if unsure).

## References

| Need | Reference |
|---|---|
| Library directory layout & naming rules | `references/library-layout.md` |
| Model affinity hints per stage | `references/model-affinity.md` |
| Anti-pattern note format and examples | `references/anti-patterns.md` |

## Templates

- `templates/stage-prompt.md` — full library entry skeleton.
- `templates/anti-pattern-line.md` — one-line loss note.
