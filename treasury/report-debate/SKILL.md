---
name: report-debate
description: >
  Stage 6 (meta) of the squad-brainstorm workflow. Reads prior-stage
  artifacts from the debate directory plus an optional metrics.json and
  produces a single 00-debate_report.md summarising the run: panel roster
  with tool versions, per-panelist / per-round stats, a convergence map
  (which panelist moved which way per round), a transcript index, pipeline
  integrity checks (incl. a moderator-favoritism check), a verbatim
  distillation of the final answer, and the live disagreements. Faithful
  summariser only — it never re-judges the debate. Graceful when
  metrics.json is absent. Triggers on: "generate debate report",
  "summarise this debate", "produce DEBATE_REPORT.md". Do NOT author the
  answer, critique, or re-merge.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: Read, Write  # reads <debate_dir> artifacts; writes only <debate_dir>/00-debate_report.md
version: 0.1.0
stage: report-debate
workflow: brainstorm
inputs:
  - { name: debate_dir,        type: string, required: true, source: workflow.runtime.debate_dir,             description: "absolute path to the debate dir (tmp/debates/<slug>-<ts>/)" }
  - { name: rounds,            type: string, required: true, source: workflow.inputs.rounds,                  description: "for the report header" }
  - { name: audience,          type: string, required: true, source: workflow.inputs.audience,               description: "for the report header" }
  - { name: debate_brief,      type: object, required: true, source: stages.frame-debate.debate_brief,        description: "proposition, contested questions, panel roster, grounding flag" }
  - { name: positions,         type: array,  required: true, source: stages.panel-open.positions,             description: "opening positions (round 0)" }
  - { name: critiques,         type: array,  required: true, source: stages.cross-examine.critiques,          description: "all rounds; [] when rounds=quick" }
  - { name: revised_positions, type: array,  required: true, source: stages.revise-positions.revised_positions, description: "all rounds; [] when rounds=quick" }
  - { name: answer,            type: string, required: true, source: stages.synthesize-consensus.answer,      description: "the final answer (verbatim distillation source)" }
  - { name: convergence,       type: object, required: true, source: stages.synthesize-consensus.convergence,  description: "agreements / live disagreements / no_favoritism_check" }
outputs:
  - { name: debate_report,    type: string, description: "markdown summary; ALSO saved to <debate_dir>/00-debate_report.md" }
  - { name: transcript_index, type: array,  description: "ordered debate artifacts with panelist + round attribution" }
---

# Skill: report-debate

## Purpose

Produce the canonical single-page debate summary — the file a user opens
first when they come back to a debate after a session reset. Self-contained,
scannable in under a minute, and it references every other artifact in the
debate directory.

This is **stage 6 of 6** and the only meta-stage. It never produces debate
content and never re-judges the verdict: it is a faithful mechanical
summariser of the run that already happened (the squad-brainstorm analogue
of researcher's `reporting-research-run`, and it inherits that skill's "faithful
summariser, not critic" constraint).

## When to use this skill

- Final stage dispatch from `workflows/brainstorm.yaml`.
- Standalone replay against an existing debate dir whose prior stages
  completed.
- Trigger phrases: "generate debate report", "summarise this debate",
  "produce DEBATE_REPORT.md".

## Do NOT use this skill to

- Author the answer, argue, critique, or re-merge — stages 2–5.
- Re-judge or second-guess the moderator's lean. Paste it verbatim.
- Modify any artifact other than `<debate_dir>/00-debate_report.md`.

## Input contract

### About `debate_dir`

`debate_dir` uses the `workflow.runtime.debate_dir` source kind — it is
**orchestration state** (the path the driver chose), not a workflow input
nor a stage output. This is the same new ref-kind researcher's `reporting-research-run`
introduced with `run_dir`. The path is authoritative for locating the
optional `metrics.json` and for writing `00-debate_report.md`.

### `metrics.json` schema (optional sidecar at `<debate_dir>/metrics.json`)

```json
{
  "debate_id": "ddd-in-enterprise-systems-20260516-101859",
  "started_at": "2026-05-16T10:18:59Z",
  "rounds": "standard",
  "turns": [
    { "stage": "panel-open", "panelist_id": "P-codex", "round": 0, "cli": "codex exec",
      "model": "gpt-5.5", "status": "ok", "exit": 0, "duration_s": 41, "bytes": 5120 }
  ],
  "metrics_incomplete": false
}
```

All fields optional. Absent/partial `metrics.json` → the affected cells
render `—` AND the mandatory "metrics incomplete" banner appears (identical
contract to researcher's `reporting-research-run`; never fail on missing metrics).

## Procedure

One LLM call. Run sequentially; do not parallelise sub-tasks.

1. **Resolve `debate_dir`.** Scan for expected artifacts: `00-panel.json`,
   `01-debate_brief.json`, `02-open__P-*.json`, `03-xexam__r*__P-*.json`,
   `04-revise__r*__P-*.json`, `05-final_answer.md`, `05-convergence.json`.
   A missing artifact renders its row cell as `(artifact missing)`, never a
   failure.
2. **Load `metrics.json`.** Set `metrics_incomplete = true` if the file is
   absent OR any expected turn (per the panel roster × stages × rounds) has
   no row or any of `status` / `duration_s` / `bytes` unset. Record
   `missing_count`.
3. **Compute content metrics** from inputs / artifacts:
   - Panel roster + tool versions + models from `00-panel.json` (or
     `debate_brief.panel`).
   - Per panelist: opening `stance` (round 0) vs final `stance`, and the
     `stance_delta` at each round (from `revised_positions[].stance_delta`).
   - Agreement / live-disagreement counts from `convergence`.
   - Grounding-citation integrity: every `[g-n]` in `05-final_answer.md`
     resolves to a `grounding_pack` item id.
   - `grounding: full|degraded` (from `debate_brief`).
   - Absent-panelist count (turns with `status: "absent"`).
   - **Moderator-favoritism check:** tally `convergence.live_disagreements[].moderator_lean`
     against each side's `held_by`; compute how often the lean coincided
     with the side `P-claude` held. This is a *surfaced statistic*, not a
     re-judgement — paste `convergence.no_favoritism_check` verbatim
     alongside it.
4. **Compose `00-debate_report.md`** in this exact order:
   1. `# Debate report — <debate_brief.proposition>`
   2. **metrics-incomplete banner** — only if `metrics_incomplete`; between
      H1 and the header line; exact wording:
      ```
      > ⚠️ **metrics incomplete** — <missing_count> turn row(s) missing or partial.
      > Re-run per `workflows/brainstorm.yaml` → Driver protocol §3 (per-CLI capture).
      > The report below renders content-only; `—` cells indicate missing telemetry.
      ```
   3. **Header line**: `_rounds: <rounds> · audience: <audience> · grounding: <full|degraded> · debate_id: <id> · started_at: <if metrics.json> · panel: codex <ver> / gemini <ver> / claude <ver> · panel_order_seed: <from 00-panel.json if present>_`
   4. **degraded-grounding banner** — verbatim `debate_brief.grounding_note`,
      only if `grounding == degraded`.
   5. `## Panel` — table: `Panelist · CLI · Tool version · Model · Turns · Final stance_delta`.
   6. `## Per-turn stats` — table, one row per `metrics.json.turns` entry
      keyed `<stage>::<panelist_id>::r<round>` + a **Total** row (sums
      `duration_s` / `bytes`; the report-debate turn itself is not
      counted). Missing cells render `—`.
   7. `## Convergence map` — for each panelist: opening stance → stance per
      round → final, with `stance_delta` arrows (`unchanged` / `narrowed →`
      / `broadened ↗` / `reversed ⟲`). End with a one-line verdict copied
      verbatim from `convergence.convergence_summary`.
   8. `## Agreements` — bulleted from `convergence.agreements`, with
      `held_by`, `[g-n]`, and `via` (independent / principled_convergence).
   9. `## Live disagreements` — from `convergence.live_disagreements`: each
      side's `held_by` + `strongest_case`, then the `moderator_lean` and
      `lean_rationale` **verbatim** (never re-judged).
   10. `## Transcript index` — fenced block, ordered, every artifact with
       round + panelist attribution. This is also the `transcript_index`
       output.
   11. `## Pipeline integrity` — bullets:
       - Grounding-citation integrity: `[g-n]` in `05-final_answer.md` ↔
         `grounding_pack` = OK / mismatched (flag only a real mismatch).
       - Moderator-favoritism: "leaned to P-claude's side on `<k>/<n>`
         contested questions" + the verbatim `no_favoritism_check`. State
         it as an observation for human audit; do NOT re-decide.
       - Absent panelists: count + which (and which rounds ran degraded).
       - Rounds executed vs the `rounds` dial expectation
         (`quick`→0 / `standard`→1 / `deep`→2 exchanges).
       - Convergence honesty: any agreement flagged
         capitulation-driven in `convergence` is echoed here.
   12. `## Distillation` — verbatim paste of the final answer's
       `## Bottom line` section (≤300 words, truncate cleanly at a sentence
       boundary; no rephrasing — same rule as `reporting-research-run`'s Distillation).
   13. `## Notes` — `Report generated by report-debate@<version> on <ISO 8601>.`
5. **Save** to `<debate_dir>/00-debate_report.md` (the `00-` prefix sorts it
   first).
6. **Return** the markdown as `debate_report` and the ordered artifact list
   as `transcript_index`.

## Constraints

- **Read-only on every artifact except `00-debate_report.md`.** Writes
  exactly one file.
- **No re-judging.** The moderator's lean/verdict is pasted verbatim. The
  favoritism tally is a surfaced statistic for human audit, never a
  correction of the verdict.
- **No new reasoning / no editorial reframing.** Distillation is a verbatim
  paste; agreements/disagreements are mechanical renders of `convergence`.
- **Graceful degradation.** Missing metrics → `—` cells + mandatory banner.
  Missing artifact → `(artifact missing)` cell. Never raise.
- **Single pass. Idempotent** (modulo the generated-on timestamp).

## Validation gate

The skill exits successfully only when ALL hold:

1. `<debate_dir>/00-debate_report.md` exists, parses, has the H1.
2. `## Per-turn stats` has one row per `metrics.json.turns` entry + a Total
   row (or the metrics-incomplete banner is present).
3. Grounding-citation integrity is re-derived from `05-final_answer.md` +
   `01-grounding_pack.json` and reported correctly.
4. The convergence map covers every panelist across every executed round.
5. The moderator-favoritism tally is computed from
   `convergence.live_disagreements` and `no_favoritism_check` is pasted
   verbatim (not paraphrased, not re-judged).
6. `## Transcript index` lists every file actually present (directory scan,
   not a hardcoded list).
7. Distillation ≤300 words, ends on a sentence boundary, verbatim.

If any check fails, fix in place and re-render; do not partial-emit.

## Worked example

See [`examples/ddd-cqrs-debate-report.md`](examples/ddd-cqrs-debate-report.md)
for a fully rendered `00-debate_report.md` of the DDD `standard` run: panel
table with tool versions, per-turn stats, convergence map with stance-delta
arrows, a live disagreement with the verbatim moderator lean, the
favoritism check, the transcript index, and the verbatim distillation.

## Anti-patterns

| Anti-pattern | Why bad | Instead |
|---|---|---|
| Re-judging the debate / "correcting" the lean | Out of scope; stage-5 work | Paste `moderator_lean`/`lean_rationale` verbatim |
| Treating the favoritism tally as a verdict override | It is an audit signal, not a re-decision | Surface the number + paste `no_favoritism_check` |
| Failing when metrics.json missing | Defeats graceful degradation | `—` cells + banner; never raise |
| Silent `—` cells without the banner | Hides the telemetry gap | Always emit the metrics-incomplete banner |
| Omitting the transcript index | The report is the entry point to the dir | Always list every artifact (dir scan) |
| Embedding full positions | Doubles disk, defeats the summary | Cap Distillation at ~300 words |
| Paraphrasing the moderator's words | Editorial drift | Verbatim paste only |

## Notes for downstream

`report-debate` is the last stage; nothing reads `debate_report`
downstream. A future runtime may index `00-debate_report.md` across debates
for an aggregate dashboard — out of scope here.
