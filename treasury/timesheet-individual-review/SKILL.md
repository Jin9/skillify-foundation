---
name: timesheet-individual-review
description: >
  Review each person's daily timesheet entries, milestone by milestone, against four rules
  (repeated task over consecutive days, 2-day repeat pattern over the monthly limit, work
  outside the JO scope, and duplicate work between same-role people) and produce a per-person
  PASS or REJECT verdict file. Use when the user asks to "review each person's timesheet against
  the rules", "flag repeated or out-of-scope timesheet entries", "check timesheets for duplicate
  work between people", "produce the timesheet review result", or "create the ReviewResult file".
  Reads Timesheet PDFs per milestone, applies the rules in batches of 3 people, checkpoints to
  TempAnalysis.md, and — only after explicit human confirmation — writes ReviewResult-JO13.txt
  style verdict file. Do NOT use to check JO manday/amount totals (milestone-jo-check) or to
  derive Jira replacement tasks for rejected people (rejected-task-replacement). Never modifies
  source files and never writes the verdict file without human sign-off.
---

# Per-Person Timesheet Review

Decide PASS or REJECT for each person on each milestone against the four review conditions, then
record the result in a verdict file after the user confirms it.

## Inputs and intake

- `Timesheet_*.pdf` inside each milestone folder (a folder whose name suffix is `M<number>`).
  Each PDF page is one person for one month. Read `Role:` and focus on the Task/Details column.
- At the start, ask the user for the **JO scope** (the tasks/products considered valid for this
  JO) — needed for condition 3. Example answer: "DGL loan origination on Paotang/NEXT: Payday,
  Tungngern, Top up, Smart Money, 5Plus, 100K".

## Review conditions

The four conditions and how to count them are detailed in `references/review-conditions.md`.
Read it before judging. Summary:

1. Same task more than 3 consecutive working days (more than 95% similar) → REJECT.
2. 2-day repeat pattern over the monthly limit `max_allowed = round(3 * M / 20)` → REJECT.
3. Task outside the JO scope the user gave → flag.
4. Same-role people with more than 95% identical task on the same day (or 1 day apart) → REJECT.

## Procedure (per milestone, until all milestones are done)

1. Read this milestone's Timesheet PDF.
2. For each page, record the working days and the task list ordered by date.
3. Apply conditions 1–3 per person, **3 people at a time** to avoid output overflow:
   - Append each batch's findings to `TempAnalysis.md` in the same folder (the checkpoint),
     keyed by `(milestone, page)` so re-entry is idempotent.
   - For any person you REJECT, capture their 5-dimension writing style from PASS days now and
     record it in `review-state.json` under `styleFingerprints` (a cheap hand-off to
     rejected-task-replacement, which then need not re-read the whole timesheet).
   - Print the batch result for the user.
   - Continue straight to the next batch (do not wait for confirmation between batches).
4. Apply condition 4 across same-role people **within this milestone**.
5. Append this milestone's summary table to `TempAnalysis.md` and rewrite `review-state.json`
   (milestone status, pages done, M-per-month, rejected entries).
6. This is a **safe compaction barrier** — state is fully checkpointed. The agent cannot
   self-invoke `/compact`; rely on auto-compaction or let the operator run `/compact` here. Then
   move to the next milestone.

The on-screen results table and the verdict-file block format are in
`references/output-format.md`.

## Checkpointing & compaction (200k context window)

Running in a fixed window, treat **disk as the source of truth**. Compaction can only happen at a
barrier where all needed state is already on disk (after a batch, after a milestone summary, or at
the boundary into rejected-task-replacement) — never mid-person or mid-page. The agent cannot
self-invoke `/compact`; rely on auto-compaction (tunable via `CLAUDE_AUTOCOMPACT_PCT_OVERRIDE`) or
the operator running `/compact` at a barrier. On resume after a compaction, re-read
`review-state.json` + `TempAnalysis.md` before doing anything else.

**Recommended under 200k:** run one subagent per milestone (Workflow tool) so each starts with a
fresh, isolated context that reads `review-state.json` — no in-agent compaction needed. The full
protocol — `review-state.json` schema, keep/drop table, crash-safe append, verdict-freeze,
resume-verify, the 200k budget, and this workflow alternative — is in `references/checkpointing.md`.

## Human-approval gate (required before writing the verdict file)

After the summary table for all milestones, **stop and wait for the user to confirm** — they may
adjust results. Only after explicit confirmation, write the verdict file
`ReviewResult-JO<JO-number>.txt` (e.g. `ReviewResult-JO13.txt`) in the same folder, following
`references/output-format.md`.

## Guardrails

- NEVER write or overwrite the `ReviewResult-*.txt` file without explicit human confirmation.
- NEVER alter a PASS/REJECT verdict toward a desired outcome; the verdict must follow the
  evidence in the timesheet.
- NEVER modify, move, or delete any source file (Timesheet PDFs, folders). `TempAnalysis.md` and
  `review-state.json` are the only files you write to before the verdict gate.
- Verdict-freeze: once a verdict is recorded in `TempAnalysis.md`, never re-derive it from a
  compacted summary on resume — re-read the file. A paraphrase could flip REJECT to PASS.
- After any compaction, resume-verify: recompute milestones/pages done from the files and confirm
  the counts before continuing, so a dropped milestone cannot pass unnoticed.
- When condition-4 "who is the copy" is unclear, mark both people and let the user decide rather
  than guessing.
