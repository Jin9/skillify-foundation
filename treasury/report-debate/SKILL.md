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

### `metrics.json` (optional sidecar at `<debate_dir>/metrics.json`)

All fields optional. Absent/partial → affected cells render `—` AND the
mandatory "metrics incomplete" banner appears; never fail on missing metrics.
Full schema + `metrics_incomplete` derivation rules: `references/compose-spec.md`.

## Procedure

One LLM call. Run sequentially; do not parallelise sub-tasks.

1. **Resolve `debate_dir`.** Scan for expected artifacts: `00-panel.json`,
   `01-debate_brief.json`, `01-grounding_pack.json`, `02-open__P-*.json`,
   `03-xexam__r*__P-*.json`, `04-revise__r*__P-*.json`, `05-final_answer.md`,
   `05-convergence.json`. (`01-grounding_pack.json` is the file gate 3 re-derives
   grounding-citation integrity from; both `01-` files are frame-debate outputs.)
   A missing artifact renders its row cell as `(artifact missing)`, never a
   failure.
2. **Load `metrics.json`.** Derive `metrics_incomplete` + `missing_count`.
   Rules: `references/compose-spec.md`.
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
4. **Compose `00-debate_report.md`** in the exact 13-section order — H1,
   metrics-incomplete banner, header line, degraded-grounding banner,
   `## Panel`, `## Per-turn stats`, `## Convergence map`, `## Agreements`,
   `## Live disagreements`, `## Transcript index`, `## Pipeline integrity`,
   `## Distillation`, `## Notes`. Full per-section content rules, both
   banners' verbatim wording, and the favoritism-tally phrasing:
   `references/compose-spec.md`. The `## Transcript index` block is also the
   `transcript_index` output; Distillation is a ≤300-word verbatim paste of
   the final answer's `## Bottom line`, truncated at a sentence boundary.
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

Seven failure modes to avoid: re-judging the lean, treating the favoritism
tally as a verdict override, failing on missing `metrics.json`, silent `—`
cells without the banner, omitting the transcript index, embedding full
positions, and paraphrasing the moderator. Full matrix (why-bad + the
do-instead for each): `references/compose-spec.md`.

## Notes for downstream

`report-debate` is the last stage; nothing reads `debate_report`
downstream. A future runtime may index `00-debate_report.md` across debates
for an aggregate dashboard — out of scope here.
