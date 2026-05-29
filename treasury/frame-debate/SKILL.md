---
name: frame-debate
description: >
  Stage 1 of the squad-brainstorm workflow. Read a finished squad-researcher
  run (05-final_report.md + 03-findings.json + 01-research_plan.json) and
  distill it into a neutral debate brief — a contestable proposition, the
  contested questions (seeded from findings' disputed_by pairs and the
  research plan's open questions), a scoring rubric (seeded from the plan's
  success_criteria), the panel roster echo, and a self-contained grounding
  pack the panel CLIs argue against. Degrades to a free-form topic +
  source_text brief when no researcher artifacts are present. Use as the
  frame-debate stage of workflows/brainstorm.yaml. Use when asked to "set up
  a debate on X", "turn this research into a debate brief", "frame the
  contested questions". One stage = one LLM call; no clarifying-question
  loops, no recursion. Do NOT use to argue a position (opening-debate-panel), critique
  one (cross-examine), revise one (revise-positions), or merge the debate
  (synthesize-consensus).
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: Read  # reads the squad-researcher run dir (read-only); emits JSON
version: 0.1.0
stage: frame-debate
workflow: brainstorm
inputs:
  - { name: research_run_dir, type: string, required: false, source: workflow.inputs.research_run_dir, description: "absolute path to a squad-researcher run dir; primary grounding source" }
  - { name: topic,            type: string, required: false, source: workflow.inputs.topic,            description: "required iff research_run_dir absent (degraded mode)" }
  - { name: source_text,      type: string, required: false, source: workflow.inputs.source_text,      description: "degraded-mode grounding text" }
  - { name: rounds,           type: string, required: false, source: workflow.inputs.rounds,           description: "quick | standard | deep; default standard" }
  - { name: audience,         type: string, required: false, source: workflow.inputs.audience,         description: "who the final answer is for; default general" }
  - { name: panel,            type: array,  required: true,  source: workflow.panel,                   description: "fixed 3-CLI roster from brainstorm.yaml" }
outputs:
  - { name: debate_brief,   type: object, description: "proposition, contested_questions[], rubric, panel echo, rounds/exchange_count, grounding flag" }
  - { name: grounding_pack, type: object, description: "self-contained [g-n] evidence the panel CLIs cite; no panelist reads squad-researcher's filesystem" }
---

# Skill: frame-debate

## Purpose

Convert a finished piece of research into a *contestable* debate setup. The
research report answered a question; this stage asks "where is that answer
actually uncertain, disputed, or under-argued?" and packages those tensions
so three independent LLM panelists can argue them out against one shared
evidence base.

`frame-debate` is the **upstream-leverage** stage (the squad-brainstorm
analogue of `plan-research`): a leading proposition, a missed real
controversy, or a biased rubric cascades into all three panelists and the
final synthesis. It **frames**; it never argues, critiques, or merges.

**Atomicity:** one stage, one LLM call. No clarifying-question loop (atomic
stages cannot pause), no recursion.

## Output contract

Emit two JSON objects. In workflow context, no prose around them.

### `debate_brief`

```jsonc
{
  "topic": "string",                     // echo (degraded) or derived from the report's H1
  "rounds": "quick|standard|deep",       // echo
  "exchange_count": 0,                   // rounds-1 -> {quick:0, standard:1, deep:2}
  "audience": "string",                  // echo
  "grounding": "full|degraded",
  "grounding_note": "string",            // what was / wasn't available; banner text for downstream
  "source_run": "string|null",           // researcher run-dir basename, or null in degraded mode
  "proposition": "string",               // ONE contestable sentence the panel takes positions on
  "supporting_propositions": ["string"], // optional sub-claims (deep only; [] otherwise)
  "contested_questions": [               // count gated by rounds (see Rounds dial)
    {
      "id": "cq1",                       // cq1..cqN, contiguous
      "question": "string",              // a genuine open/disputed question, NOT leading
      "why_contested": "string",         // e.g. "findings g-3 and g-7 disagree" / "plan open question"
      "seed": "finding_dispute|plan_open_question|synthesizer_gap|degraded_inferred",
      "evidence_refs": ["g-3","g-7"]     // ids into grounding_pack.items ([] only in degraded)
    }
  ],
  "rubric": [                            // seeded from research_plan.success_criteria + debate criteria
    { "id": "r1", "criterion": "string", "weight": 1.0 }   // weight meaningful only on deep
  ],
  "panel": [ /* echo of the input `panel` rows, byte-for-byte */ ],
  "steelman_directive": "string|null",   // deep-only: forces each panelist to argue the strongest opposing case once
  "brief_skipped": false
}
```

### `grounding_pack` — the closed evidence world for ALL panel CLIs

```jsonc
{
  "summary": "string",                   // <=400-word NEUTRAL digest of 05-final_report.md (or source_text)
  "items": [                             // dereferenced from 03-findings.json
    {
      "id": "g-1",                       // g-1..g-N, stable; panelists cite [g-1] like researcher's [n]
      "claim": "string",
      "evidence": "string",              // verbatim excerpt carried from the finding
      "confidence": "high|medium|low",   // carried from the finding
      "disputed": true,                  // mirrors findings.disputed_by being non-empty
      "origin": "finding:f-12"           // provenance back to the researcher artifact
    }
  ],
  "out_of_scope": ["string"]             // carried from research_plan.out_of_scope (panelists must not re-litigate)
}
```

### `brief_skipped` failure shape

Return this when neither grounding path resolves (no valid `research_run_dir`
AND empty/whitespace `topic`), or the topic is unintelligible after one
parse attempt:

```jsonc
{
  "topic": "<echo or empty>",
  "rounds": "<echo>", "exchange_count": 0, "audience": "<echo>",
  "grounding": "degraded",
  "grounding_note": "brief_skipped: <reason>",
  "source_run": null,
  "proposition": "", "supporting_propositions": [],
  "contested_questions": [], "rubric": [],
  "panel": [ /* echo */ ], "steelman_directive": null,
  "brief_skipped": true
}
```

Downstream stages treat `brief_skipped: true` as a hard stop; the workflow
returns an empty `answer_markdown` (mirrors researcher's `plan_skipped`).

## Procedure

Execute all of the following in a single LLM call.

1. **Resolve the grounding source.** If `research_run_dir` is a non-empty
   path, check it contains `05-final_report.md`, `03-findings.json`, and
   `01-research_plan.json`. All present → `grounding: full`. Otherwise fall
   back: if `topic` is non-empty → `grounding: degraded`; else emit the
   `brief_skipped` shape and stop. Never invent a topic.

2. **(full) Read the three researcher artifacts** with the Read tool
   (read-only — never write into the researcher run dir):
   - `05-final_report.md` → source for the neutral `summary` and the
     `proposition` (the report's central claim, restated as one contestable
     sentence — not as a settled fact).
   - `03-findings.json` → every finding becomes a `grounding_pack.items[]`
     entry with a fresh stable `g-` id and `origin: "finding:<f-id>"`.
     Carry `claim`, `evidence` (verbatim), `confidence`. Set
     `disputed: true` when the finding's `disputed_by` is non-empty.
   - `01-research_plan.json` → `success_criteria` seed the `rubric`;
     `out_of_scope` carries into `grounding_pack.out_of_scope`;
     `thesis_question` cross-checks the `proposition`.

3. **(degraded) Build a thin pack** from `topic` + optional `source_text`.
   `summary` = a neutral digest of `source_text`, or — if no `source_text` —
   a one-sentence statement that the debate is ungrounded and positions are
   reasoning-only. `items` derived from `source_text` passages (still `g-n`
   ids) or `[]`. `out_of_scope: []`. See
   [`references/degraded-mode.md`](references/degraded-mode.md).

4. **Write the `proposition`.** One declarative, *contestable* sentence a
   reasonable panelist could argue for or against. Reject phrasings the
   evidence has already settled, and reject loaded/leading wording. If the
   research is genuinely one-sided (no real debate), still emit, but say so
   in `grounding_note` — the panel will largely concur and that is a valid,
   honest outcome.

5. **Seed `contested_questions`** in priority order (see
   [`references/contested-question-rubric.md`](references/contested-question-rubric.md)):
   1. `findings[].disputed_by` pairs → `seed: finding_dispute` (the
      controversy was already identified by `extract-findings`; convert each
      to a neutral question, `evidence_refs` = the two `g-` ids).
   2. `research_plan` unmet/partial `success_criteria` and open
      thesis tensions → `seed: plan_open_question`.
   3. Synthesizer-inferred tensions the report under-argues →
      `seed: synthesizer_gap`.
   4. (degraded only) inferred from `topic`/`source_text` →
      `seed: degraded_inferred`.
   Reject leading or single-answer questions. Count gated by the Rounds dial.

6. **Seed the `rubric`** from `research_plan.success_criteria` plus the four
   standard debate criteria: evidential grounding (cites `[g-n]`), internal
   consistency, directly addresses the proposition, engages the strongest
   counter. 3 criteria on `quick`, 5 on `standard`, weighted on `deep`.

7. **Set `exchange_count` = rounds − 1** (`quick`→0, `standard`→1,
   `deep`→2). On `deep`, set a `steelman_directive` string.

8. **Echo `panel` verbatim** into `debate_brief.panel`.

9. **Self-check** against the Validation gate. Fix in place; no second call.

10. **Emit** `debate_brief` and `grounding_pack`. JSON only in workflow
    context.

## Rounds dial

| `rounds` | `contested_questions` | `rubric` | Extras |
|----------|----------------------|----------|--------|
| `quick` | ≤ 3 | 3 criteria, no weights | no steelman |
| `standard` | 3–6 | 5 criteria, no weights | no steelman |
| `deep` | 6–10 | 5+ criteria, weighted | `steelman_directive` set; `supporting_propositions` allowed |

If `rounds` is missing or unrecognized, treat it as `standard`.

`audience` is **echoed only**. It shapes the final answer (consumed by
`synthesize-consensus` via the researcher audience-profile mapping), never
which questions the panel debates.

## Constraints

- DO NOT make a second LLM call. Atomic stage = one call.
- DO NOT ask the user clarifying questions. Assume; record assumptions in
  `grounding_note`.
- DO NOT frame a proposition the evidence has already decided, or use
  leading / loaded question phrasing. Neutral framing only.
- DO NOT invent evidence, claims, or `g-` items not present in the
  researcher artifacts (full) or `source_text` (degraded).
- DO NOT write into the squad-researcher run dir or anywhere outside this
  stage's outputs — `allowed-tools: Read` is read-only by contract.
- DO NOT modify `workflows/brainstorm.yaml`.
- DO NOT emit prose around the JSON in workflow context.

## Validation gate

Before returning, verify all of:

- [ ] `grounding ∈ {full, degraded}` and matches what actually resolved.
- [ ] `proposition` is exactly one contestable sentence (not a settled fact,
      not leading), unless `brief_skipped`.
- [ ] `contested_questions` count matches the Rounds dial; every entry is a
      genuine open question with a non-empty `why_contested`.
- [ ] Every `contested_questions[].evidence_refs` id exists in
      `grounding_pack.items` (or `seed == degraded_inferred` with `[]`).
- [ ] Every `rubric` item is testable; count/weights match the Rounds dial.
- [ ] `panel` is echoed byte-for-byte from the input `panel`.
- [ ] `exchange_count == rounds - 1` (0 / 1 / 2).
- [ ] `grounding_pack.items[].id` are contiguous `g-1..g-N`; every `origin`
      traces to a real finding (full mode).
- [ ] `brief_skipped` is `false` unless the failure shape was emitted.

## Worked example

See [`examples/ddd-cqrs-debate-brief.md`](examples/ddd-cqrs-debate-brief.md)
for a complete worked input → output: a `standard` brief built from a real
squad-researcher run (a `disputed_by` pair → a contested question, a
`success_criterion` → a rubric item), plus a degraded-mode mini-example.

## References

- [`references/contested-question-rubric.md`](references/contested-question-rubric.md)
  — what makes a question genuinely contestable vs. leading; how to convert a
  `disputed_by` finding pair into one neutral question.
- [`references/degraded-mode.md`](references/degraded-mode.md) — the exact
  fallback procedure and the `grounding_note` banner wording.
