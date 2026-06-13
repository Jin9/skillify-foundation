# Checkpointing & Compaction (200k context window)

Principle: **disk is the source of truth; the context window is scratch.** Compaction must only
happen at a point where everything needed to resume is already written to a file. Follow this and
compaction can never lose work.

Compaction is **not agent-triggerable**: the agent cannot run `/compact` on itself. It happens
either by **auto-compaction** (fires near ≈95% of the window; trigger it earlier by lowering the
`CLAUDE_AUTOCOMPACT_PCT_OVERRIDE` env var, e.g. to 60) or by the **operator typing `/compact`** at
a barrier. So the job here is to make every checkpoint barrier safe to compact at.

## Compact only at safe barriers

| Safe barrier | On disk before compaction | Safe to compact? |
|--------------|---------------------------|------------------|
| After a 3-person batch | batch findings appended to `TempAnalysis.md` | yes |
| After a milestone summary | milestone table + `review-state.json` rewritten | yes (best) |
| Stage 2 → 3 boundary | `ReviewResult-JO*.txt` written + confirmed | yes (cleanest) |
| Mid-person / mid-PDF-page | nothing durable yet | never |

## review-state.json (the resume token)

Rewrite this file in the same folder at every barrier:

```json
{
  "joId": "JO13",
  "scope": "DGL loan origination on Paotang/NEXT: Payday, Tungngern, Top up, ...",
  "jiraCloudId": "xxx.atlassian.net",
  "milestones": { "M1": "done", "M2": "in_progress", "M3": "pending" },
  "pagesDone": { "M2": [1, 2, 3, 4, 5, 6] },
  "mandayPerMonth": { "M1": 21, "M2": 22 },
  "rejectedEntries": [ { "milestone": "M1", "page": 15, "name": "Chattarika", "dates": ["19/11"] } ],
  "styleFingerprints": { "Chattarika": { "case": "Title Case", "ticketRef": "generic series", "verbs": "Executed/Opened/Retested", "qualifier": "scope in parens", "lineStructure": "single line" } },
  "schemaVersion": 1
}
```

On resume after a compaction, read `review-state.json` + `TempAnalysis.md` first — that fully
reconstructs where you are in O(1) reads. Do not re-derive anything already recorded.

## Keep vs. drop across the boundary

| Keep (write to disk, must survive) | Drop (let compaction discard) |
|------------------------------------|-------------------------------|
| Per-person PASS/REJECT verdicts | Raw PDF page text |
| Rejected dates + short reasons | Intermediate reasoning / narration |
| JO scope, M-per-month, Jira cloudId | Raw Jira JSON (after `jq` extract) |
| Per-person style fingerprint (for Stage 3) | |
| Provenance: `milestone + page` per finding | |

Raw PDF text is the biggest context hog — read one milestone's timesheet at a time, and within a
large milestone read in page-ranges. Never hold all milestones' PDFs resident at once.

## Crash-safety

- Append to `TempAnalysis.md`; never rewrite it — a crash mid-write cannot corrupt earlier findings.
- Key every finding by `(milestone, page)`. On resume, if a key already exists, skip it. Re-entry
  is then idempotent: no double-count, no re-judge.

## Verdict-freeze and resume-verify (safety)

- **Verdict-freeze**: once a verdict is in `TempAnalysis.md`, it is frozen. On resume, re-read it;
  never re-derive a verdict from a compacted summary — a paraphrase could flip REJECT to PASS.
- **Resume-verify**: after every compaction, recompute counts from the files (milestones done,
  pages done vs PDF page count) and print a one-line resume summary; confirm it matches before
  continuing. This catches compaction silently dropping a milestone.

## Budget (200k)

- Reserve ~30–40k for system + skill + state files.
- Process one milestone, 3 people at a time → peak resident PDF text stays in the tens of k.
- Set `CLAUDE_AUTOCOMPACT_PCT_OVERRIDE` lower (e.g. 60) so auto-compaction fires earlier — at a
  barrier, not mid-task — or have the operator run `/compact` at the barrier. The checkpoint is
  already written, so the compaction can never fail to fit.
- **Loop cap**: if a milestone reverts `done → in_progress` in `review-state.json`, stop and
  escalate — that is a compaction-induced redo loop, not normal flow.

## Workflow alternative (recommended under 200k)

Instead of one long-lived agent compacting between milestones, run **one subagent per milestone**
(the Workflow tool). Each `agent()` call gets a fresh, isolated context: it reads
`review-state.json`, reviews its milestone, appends to `TempAnalysis.md`, rewrites the state, and
returns a small result. Because each step starts near-empty and state lives on disk, there is **no
growing context to compact** and no reliance on in-agent compaction. This is the cleanest way to
stay under the window; the verdict-freeze, idempotent keying, and resume-verify rules above still
apply (each subagent re-reads the state, never re-derives a frozen verdict).
