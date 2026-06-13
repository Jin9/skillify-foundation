---
name: rejected-task-replacement
description: >
  For people whose timesheet entries were REJECTED, derive replacement task descriptions sourced
  from Jira (read-only), written in each person's own style, and save them to a TaskForJO file
  after the user reviews them. Use when the user asks to "find Jira tasks to replace the rejected
  timesheet entries", "derive replacement tasks for the rejected people", "generate the TaskForJO
  file from Jira", or "suggest tasks for the rejected days". Asks for the Jira project / cloudId
  and the Fix Version per milestone, queries Stories/Tasks (Dev) or SIT Defects (Tester), reads
  each rejected person's PASS-day writing style, and produces replacement tasks at least equal in
  count to that person's rejected entries. Do NOT use to decide who is rejected
  (timesheet-individual-review) or to check JO totals (milestone-jo-check). Jira is read-only:
  never create or modify Jira issues, never fabricate ticket IDs, and never write the output file
  without human review.
---

# Replacement Tasks for Rejected Entries

Produce believable, in-style task descriptions for each rejected person, grounded in real Jira
issues, and save them only after the user reviews them.

## Inputs and intake

- The verdict file from timesheet-individual-review (e.g. `ReviewResult-JO13.txt`) listing who
  was rejected and on which dates, plus each person's Timesheet PDF.
- This skill starts at the Stage 2 → 3 boundary, which is a clean (compacted / fresh) entry
  point. First read `review-state.json` (read-only) if present — it gives the rejected list, JO
  scope, Jira cloudId, M-per-month, and the `styleFingerprints` captured during the review. Use
  the persisted fingerprint instead of re-reading a person's whole timesheet under a tight
  budget; only re-read PASS days if no fingerprint was recorded.
- If `review-state.json` has no `jiraCloudId`, ask the user for the **Jira project code or board
  URL** (derive the cloudId from the URL hostname). Example: `DGL` or a board URL.

## Procedure

1. **Per milestone, ask the Fix Version(s)** — state the period and the number of rejected tasks
   so the user can answer.
2. **Query Jira (read-only)** by role. JQL templates, the sub-task fallback, the 0-results
   fallback, and the `jq` extraction are in `references/jira-queries.md`.
3. **Read issue content to derive tasks**:
   - Dev: read the Story Description (Business Logic + Acceptance Criteria) → implementation tasks.
   - Tester: use the defect summary directly (it names the feature/scenario) → testing activities.
4. **Read the rejected person's Timesheet once, for two purposes**:
   - Duplicate guard: note tasks already present so suggestions do not repeat them.
   - Writing style: scan **PASS days only** (at least 5) to extract the person's style across the
     five dimensions in `references/writing-style.md`. Use that person's own style, not a role
     average.
5. **Derive tasks** by role and extracted style — see `references/writing-style.md` for the
   role-based task menus, the role-default fallback, and the `[Detail]` block rule.
6. **Tester only — re-check condition 4** on the *suggested* tasks against other same-role
   testers on the same/adjacent day; vary feature area or pick a different defect ID if too
   similar. See `references/writing-style.md` (continuity rules).
7. **Review with the user first**, then write the file.

## Human-approval gate (required before writing the output file)

Only after the user reviews the proposed tasks, write `TaskForJO<JO-number>.txt` (e.g.
`TaskForJO13.txt`) in the same folder.

```
=====M<milestone>=====
<person name>
- Task

<person name>
- Task
```

Rules: no blank line between a name and its first task; blank line when the person or milestone
changes; blank line between task items when that person uses a `[Detail]` block; the number of
tasks for a person must be at least the number of their rejected entries.

## Guardrails

- Jira access is **read-only** — NEVER create, edit, transition, or comment on any Jira issue.
- NEVER fabricate ticket IDs (for example `DGL-` numbers). Ticket references in a derived task
  must come from the actual query results; if you have no real ID, omit the reference.
- NEVER write or overwrite the `TaskForJO*.txt` file without the user's review and go-ahead.
- NEVER modify, move, or delete any source file (timesheets, verdict file, `review-state.json`,
  folders) — read them only; the `TaskForJO*.txt` deliverable is the only file you write.
- Write tasks in English, with no tags like `[FE]` / `[BAU]`, not duplicating the existing
  timesheet, in a style that reads as written by that person.
