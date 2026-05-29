---
name: reporting-research-run
description: >
  Stage 6 (meta) of the squad-researcher workflow. Reads prior-stage
  artifacts from the run directory plus an optional metrics.json sidecar,
  and produces a single 00-run_report.md summarising the run: per-stage
  stats table, content metrics, file inventory, a short distillation
  paragraph from the final report, and the top risks copied from
  review_notes. Use as the final stage of the researcher workflow, or as
  a standalone replay against an existing run directory. Triggers on:
  "generate run report", "summarise this run", "produce RUN_REPORT.md".
  Do NOT use to author research content, extract claims, re-review the
  report, or replace any earlier stage. Graceful when metrics.json is
  absent — the report still renders content-only.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: Read, Write  # reads <run_dir> artifacts; writes only <run_dir>/00-run_report.md
version: 0.1.1
stage: report-run
workflow: researcher
inputs:
  - { name: run_dir,          type: string, required: true, source: workflow.runtime.run_dir,             description: "absolute path to the run directory (tmp/runs/<topic-slug>-<ts>/)" }
  - { name: topic,            type: string, required: true, source: workflow.inputs.topic,                description: "the subject of the run" }
  - { name: depth,            type: string, required: true, source: workflow.inputs.depth,                description: "quick | standard | deep — used for the report header" }
  - { name: audience,         type: string, required: true, source: workflow.inputs.audience,             description: "audience used for the run — for the report header" }
  - { name: research_plan,    type: object, required: true, source: stages.plan-research.research_plan,   description: "from plan-research; used for sub-question count + dimensions" }
  - { name: sources,          type: array,  required: true, source: stages.search-sources.sources,        description: "from search-sources; used for source count + diversity stats" }
  - { name: findings,         type: array,  required: true, source: stages.extract-findings.findings,     description: "from extract-findings; used for findings + contradictions + coverage" }
  - { name: final_report,     type: string, required: true, source: stages.review-report.final_report,    description: "post-review markdown; used for word count + distillation paragraph" }
  - { name: cited_sources,    type: array,  required: true, source: stages.review-report.cited_sources,   description: "post-review citation set; used for citation integrity check" }
  - { name: review_notes,     type: string, required: true, source: stages.review-report.review_notes,    description: "from review-report; used for edit counts + top risks" }
outputs:
  - { name: run_report, type: string, description: "markdown summary of the run; ALSO saved to <run_dir>/00-run_report.md as a side-effect" }
---

# Skill: Reporting Research Run

## Purpose

Produce the canonical single-page run summary for a `researcher` workflow run. The report is a markdown artifact saved to `<run_dir>/00-run_report.md` and returned as `run_report`. It is the file a user opens first when they come back to a run after a session reset — so it must be self-contained, scannable in under a minute, and reference every other artifact in the run directory.

This is **stage 6 of 6** in the `researcher` workflow, and the only meta-stage: it never produces research content; it only summarises the run that already happened.

## When to use this skill

- Final stage dispatch from `workflows/researcher.yaml`.
- Standalone replay against an existing run directory whose prior stages already completed (e.g., to backfill a report on an older run).
- Trigger phrases: "generate run report", "summarise this run", "produce RUN_REPORT.md", "post-run summary".

## Do NOT use this skill to

- Author research prose, extract claims, gather sources, or do any other research work — those are stages 1-4.
- Re-review the report — that is `review-report` (stage 5).
- Make decisions, recommendations, or evaluations beyond what the prior stage artifacts already contain. Be a faithful summariser, not a critic.
- Modify any artifact other than `<run_dir>/00-run_report.md`.

## Input contract

### About `run_dir`

`run_dir` has a new `source:` kind compared to the other skills: `workflow.runtime.run_dir`. This means it is **orchestration state** (the path the orchestrator chose when it set up the run), not a workflow input nor a stage output. Every other input in this skill is either `workflow.inputs.<x>` or `stages.<id>.<output>` — the established conventions. Only `report-run` uses `workflow.runtime.<x>`, and only for `run_dir`.

The path is treated as authoritative for two purposes:
1. Locating the optional `metrics.json` sidecar.
2. Writing the output `00-run_report.md`.

### `metrics.json` schema (optional sidecar at `<run_dir>/metrics.json`)

If the orchestrator captured per-stage runtime data, it writes a file like this:

```json
{
  "run_id": "ddd-cqrs-20260514-161443",
  "started_at": "2026-05-14T16:14:43Z",
  "stages": {
    "plan-research":     { "tokens": 26416,  "tool_uses": 4,  "duration_ms": 72487 },
    "search-sources":    { "tokens": 122714, "tool_uses": 59, "duration_ms": 360196 },
    "extract-findings":  { "tokens": 57252,  "tool_uses": 5,  "duration_ms": 146230 },
    "synthesize-report": { "tokens": 95641,  "tool_uses": 26, "duration_ms": 459368 },
    "review-report":     { "tokens": 87218,  "tool_uses": 9,  "duration_ms": 351059 }
  }
}
```

All fields are optional. If `metrics.json` is absent or any stage's row is missing, the corresponding cells in the report render as `—` (em-dash) and the report still renders. Do not fail on missing metrics.

## Procedure

One LLM call. Run sequentially; do not parallelise sub-tasks.

1. **Resolve `run_dir`.** Confirm the path exists and contains the expected stage artifacts (`01-research_plan.json`, `02-sources.json`, `03-findings.json`, `04-draft_report.md`, `04-cited_sources.json`, `05-final_report.md`, `05-cited_sources.json`, `05-review_notes.md`). If any is missing, render the corresponding row's "Key output" cell as `(artifact missing)` rather than failing.

2. **Load `metrics.json`** at `<run_dir>/metrics.json` if it exists. Parse stages. Record which stage rows are present. Set `metrics_incomplete = true` if the file is absent **or** any of the 6 stage rows (`plan-research`, `search-sources`, `extract-findings`, `synthesize-report`, `review-report`, `report-run`) is missing or has any of `tokens` / `tool_uses` / `duration_ms` unset. Record `missing_count` = number of stage rows that fail this check (0–6).

3. **Compute content metrics** from the in-memory inputs (or by re-reading the run-dir artifacts, equivalently):
   - From `research_plan`: thesis question (verbatim), number of sub-questions, coverage dimensions used.
   - From `sources`: total count, total unique queries, source-type histogram, top domain + its share.
   - From `findings`: total count; confidence histogram (high / medium / low); contradictions count (number of `disputed_by` pairs); coverage statuses per sub-question (count of covered / partial / uncovered); `unsupported_sources` count.
   - From `final_report`: total word count (body, excluding the Sources list), section count.
   - From `cited_sources`: count.
   - From `review_notes`: claims_verified; edits_applied broken down by change_type (swap / add / drop / revise_sentence / delete_sentence / resolve_contradiction); flags_raised count; top risks list (verbatim, truncated to 3).

4. **Compose `00-run_report.md`** with the following sections, in this exact order:
   1. `# Run report — <topic>` (H1, topic verbatim).
   2. **Metrics-incomplete banner** — render this block ONLY if `metrics_incomplete` is true (from step 2); otherwise emit nothing for this slot:
      ```
      > ⚠️ **metrics incomplete** — <missing_count>/6 stage rows missing or partial.
      > Re-run per `workflows/researcher.yaml` → Driver protocol §2 (per-stage append).
      > The report below renders content-only; `—` cells indicate missing telemetry.
      ```
      The banner appears between the H1 and the header line. Do not raise; do not skip subsequent sections; do not editorialise — the wording above is the contract.
   3. **Header line**: `_depth: <depth> · audience: <audience> · run_id: <run_id parsed from run_dir basename> · started_at: <if metrics.json> · stages: 6/6_`.
   4. `## Per-stage stats` — a markdown table with one row per stage plus a final **Total** row. Columns:
      - `Stage` (1. plan-research, 2. search-sources, ...)
      - `Tokens` (from `metrics.stages.<id>.tokens` or `—`)
      - `Tools` (from `metrics.stages.<id>.tool_uses` or `—`)
      - `Wall` (formatted as `MM:SSs` or `Ss` from `duration_ms`, or `—`)
      - `Key output` (the content metric specific to that stage; see below)
   5. `## Pipeline integrity` — bullet list:
      - Citation bijection: `<cited_count>/<sources_count>` sources cited; `[n] ↔ cited_sources` integrity = OK / mismatched (only flag a real mismatch).
      - Contradictions surfaced: `<count>` (list them by finding pair if ≤3, else summarise).
      - Partial-coverage sub-questions: `<list of q-ids>` or `none`.
      - Review edits applied: total + by change_type.
      - Unsupported sources: `<count>` (and which sources, if ≤3).
   6. `## File inventory` — fenced code block listing every file in the run dir with size in KB. Match the format of the inventory I print after each run.
   7. `## Distillation` — paste the first 2–3 paragraphs of the final report's `## Executive Summary` section verbatim (no rephrasing, no editorial framing). Cap at ~300 words; truncate cleanly at a sentence boundary if needed.
   8. `## Top risks` — copy review_notes' "Top risks" section verbatim if present (the review-report skill puts these in a stable section). If review_notes lacks a "Top risks" section, render `_(none recorded by review-report)_`.
   9. `## Notes` — free-form footer with one line: `Report generated by report-run@<version> on <ISO 8601 timestamp>.`

5. **Save** the composed markdown to `<run_dir>/00-run_report.md`. The `00-` prefix ensures it sorts first in directory listings — the user sees it before any other artifact.

6. **Return** the same markdown string as `run_report`.

### Per-stage "Key output" cell values

For the per-stage table's "Key output" column, render these strings (filling in numbers from the artifacts):

| Stage | Key output cell |
|---|---|
| 1. plan-research | `<N> sub-questions, <M> dimensions` |
| 2. search-sources | `<N> sources, <Q> queries, top domain <D%>` |
| 3. extract-findings | `<N> findings (<H> hi / <M> med / <L> lo), <C> contradictions, <U> unsupported` |
| 4. synthesize-report | `<W>-word draft, <S> sections, <K>/<N> sources cited` |
| 5. review-report | `<E> edits applied, <F> flags, <R> claims verified` |
| 6. report-run | `RUN_REPORT.md` (self-reference) |

The Total row sums tokens, tools, and wall across the 5 prior stages (report-run's own metrics are not double-counted; if present in metrics.json they appear in row 6 only).

## Output contract

- **Return value**: `run_report` (string) — the full markdown report.
- **Side effect**: `<run_dir>/00-run_report.md` is written (or overwritten) with the same content.

The file is the load-bearing output. A caller that ignores the return value but reads the file later will see the same report.

## Constraints

- **Read-only on every other artifact.** This skill writes exactly one file: `<run_dir>/00-run_report.md`. It MUST NOT modify any other file in the run directory or anywhere else in the project.
- **No web access.** No new findings, no new sources, no fact-checking against the world.
- **No editorial reframing.** The Distillation section is a verbatim paste from the final report. The Top risks are verbatim from review_notes. The numbers in the table are computed; everything else is mechanical.
- **Single pass.** No recursion, no re-call.
- **Graceful degradation.** Missing or partial `metrics.json` → `—` cells AND the "metrics incomplete" banner described in Procedure step 4.2, not failure. Missing artifact → `(artifact missing)` cell, not failure. The banner is mandatory whenever any of the 6 stage rows is missing or partial — silent `—` cells alone are NOT acceptable.
- **Idempotent.** Running the skill twice produces the same output (modulo the "Report generated by…" timestamp).

## Validation gate

The skill exits successfully only when ALL of the following hold:

1. `<run_dir>/00-run_report.md` exists and is a valid markdown file (parses, has the H1 header).
2. The per-stage table has exactly 6 rows + a Total row.
3. The citation bijection check matches reality: re-derive it from `<run_dir>/05-final_report.md` and `<run_dir>/05-cited_sources.json` and confirm the reported `<cited_count>` and integrity flag are correct.
4. The file inventory lists every file actually present in the run dir (use a directory scan, not a hardcoded list).
5. The Distillation section is ≤300 words and ends at a sentence boundary.

If any check fails, fix in-place and re-render; do not partial-emit.

## Worked example

See `examples/ddd-cqrs-run-report.md` for a complete rendered `00-run_report.md`.
It shows: the per-stage stats table with a Total row, pipeline-integrity bullets
including contradictions, the file inventory, a verbatim Distillation pasted
from the final report's Executive Summary, Top risks copied from `review_notes`,
and the standard footer line.

## Anti-patterns

| Anti-pattern | Why bad | What to do instead |
|---|---|---|
| Rewriting the final report's prose | Out of scope; that's stage 4-5 work | Paste the Executive Summary verbatim into Distillation |
| Adding new risks or commentary | Inserts unsupported opinion | Only copy what review_notes already states |
| Failing when metrics.json missing | Defeats the graceful-degradation requirement | Render `—` cells + banner; never raise |
| Silently rendering `—` cells without the banner | Hides the telemetry gap from a reader who only scans the table | Always emit the "metrics incomplete" banner from Procedure step 4.2 when any stage row is missing/partial |
| Recomputing token usage from scratch | Token data isn't in any LLM-call output | Pull only from metrics.json |
| Omitting the file inventory | The report is the entry point to the run dir | Always list every artifact |
| Embedding the full final report | Doubles disk usage, defeats the summary purpose | Cap Distillation at ~300 words |

## Notes for downstream

`report-run` is the last stage; nothing reads `run_report` downstream. A future runtime may index `00-run_report.md` across runs for an aggregate-runs dashboard, but that's out of scope for this skill.
