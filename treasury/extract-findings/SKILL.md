---
name: extract-findings
description: >
  Stage 3 of the researcher workflow. Given a research_plan (from
  plan-research) and sources (from search-sources), pull out structured,
  source-grounded claims organized against the plan's sub-questions.
  Every claim carries a verbatim evidence excerpt and a source_id that
  must exist in the input sources — no source_id, no claim. Use when the
  researcher workflow has just finished search-sources and synthesize-report
  has not yet started, or when a caller passes a research_plan + sources
  pair and asks for "extract findings", "structured claims with citations",
  "evidence table", or "claim list against the plan". Do NOT use to run
  search, to write prose paragraphs (use synthesize-report), or to judge
  whether claims are true in the world (review-report owns adversarial
  checking). Output is JSON-shaped findings — never freeform notes.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # no external tools — closed-context extraction from input sources; no fetch
version: 0.1.2
stage: extract-findings
workflow: researcher
inputs:
  - { name: research_plan, type: object, required: true, source: stages.plan-research.research_plan, description: "from plan-research" }
  - { name: sources,       type: array,  required: true, source: stages.search-sources.sources,      description: "from search-sources" }
outputs:
  - { name: findings, type: array, description: "structured claims: {id, sub_question_id, claim, evidence, source_id, supporting_source_ids[], disputed_by[], confidence, note}" }
---

# Extract Findings

## Purpose

Convert a bag of gathered sources into a structured, citation-grounded
findings table that is already aligned to the research_plan's
sub-questions. The next stage (synthesize-report) consumes this output as
an outline-ready evidence map; it should not have to re-read the source
bodies to write the report.

This skill is the **faithfulness gate** of the workflow. It enforces:
every claim is anchored to an `evidence` excerpt and a `source_id`. Claims
without that anchor are dropped, not softened. Faithfulness (the citation
supports the claim as stated) is the goal here; **correctness** (the
claim is true in the world) is owned by `review-report` downstream.

## When to use this skill

- Stage 3 of the `researcher` workflow runs this skill automatically.
- Standalone invocation when a caller hands over a `research_plan` +
  `sources` pair and asks for an evidence table.
- Trigger phrases: "extract findings", "pull structured claims",
  "evidence table", "claim-source map", "findings against the plan".

Do NOT use this skill to:
- Plan research (use `plan-research`).
- Search for new sources (use `search-sources`). This skill cannot fetch
  anything — it works only with the sources the workflow already gathered.
- Compose the prose report (use `synthesize-report`).
- Verify whether claims are factually accurate against the world (use
  `review-report`).
- Edit `workflows/researcher.yaml` or any other workflow file.

## Inputs

The workflow contract (`workflows/researcher.yaml`, stage `extract-findings`)
passes exactly two inputs:

| Field | Required | Source | Shape |
|---|---|---|---|
| `research_plan` | yes | `stages.plan-research.research_plan` | object with `sub_questions[]`; each entry has at least `id` (stable, e.g. `sq-1`) and `question` (text). May also carry `rationale`, `priority`. |
| `sources` | yes | `stages.search-sources.sources` | list of objects; each entry has at least `id` (stable, e.g. `s-1`), `title`, `url`, and a body field (`text`, `content`, or `excerpt`) the model can read. May also carry `author`, `published_at`, `accessed_at`. |

Both inputs are passed as-is from upstream stages. The skill must not
fabricate fields that the upstream stages did not provide; if a source
arrived without a body, treat it as unreadable (see the validation gate).

## Outputs

Single output `findings`, a JSON-shaped object with the following shape:

```json
{
  "findings": [
    {
      "id": "f-1",
      "sub_question_id": "sq-1",
      "claim": "One declarative sentence.",
      "evidence": "Verbatim or close-paraphrase excerpt from the source body.",
      "source_id": "s-3",
      "supporting_source_ids": ["s-7"],
      "disputed_by": [],
      "confidence": "high",
      "note": null
    }
  ],
  "coverage": [
    {
      "sub_question_id": "sq-1",
      "finding_ids": ["f-1", "f-2"],
      "status": "covered",
      "gap_note": null
    }
  ],
  "unsupported_sources": [
    { "source_id": "s-12", "reason": "off-topic — no claims extracted" }
  ]
}
```

### Resolving the workflow's open question

`workflows/researcher.yaml` asks: *findings shape — list of
`{claim, evidence, source_id, confidence}` or freeform?* This skill picks
**structured list, extended**. Rationale:

- `synthesize-report` and `review-report` both need to address individual
  claims (cite them, flag them, drop them). Freeform notes force them to
  re-parse text. Structured rows are addressable.
- Adding `sub_question_id` makes the output **outline-conditioned**
  (STORM's term): the synthesis stage gets a free outline from
  `coverage[]` and never has to reorganize raw notes.
- Adding `supporting_source_ids` and `disputed_by` lets corroboration
  and contradiction survive into synthesis instead of being collapsed.
- `confidence` is qualitative (`high|medium|low`), not numeric — see
  [`references/confidence-rubric.md`](references/confidence-rubric.md) for
  the rubric. Verbalized numeric confidence from LLMs is poorly calibrated.

### Per-finding fields

| Field | Required | Type | Rule |
|---|---|---|---|
| `id` | yes | string | Stable, monotonic (`f-1`, `f-2`, ...). Used by downstream stages to reference the finding. |
| `sub_question_id` | yes | string | Must match an `id` in `research_plan.sub_questions[]`, OR the literal `"unmapped"`. |
| `claim` | yes | string | One declarative sentence. ≤ 240 chars. No hedging beyond what the source uses. |
| `evidence` | yes | string | Verbatim excerpt or close paraphrase from the source body, ≤ 600 chars. Quoted material in double quotes. |
| `source_id` | yes | string | Must match an `id` in input `sources[]`. The primary source for this claim. |
| `supporting_source_ids` | no | string[] | Other source ids that independently support the claim. Use for corroboration. |
| `disputed_by` | no | string[] | Source ids that contradict the claim. When non-empty, also emit the opposing claim as a separate finding. |
| `confidence` | yes | enum | `high` \| `medium` \| `low`. Per the rubric below. |
| `note` | no | string \| null | Caveats, missing nuance, unmapped-bucket explanation, or other context. |

### Per-coverage fields

| Field | Required | Type | Rule |
|---|---|---|---|
| `sub_question_id` | yes | string | Mirrors a `research_plan.sub_questions[].id`. |
| `finding_ids` | yes | string[] | Finding ids that address this sub-question. Empty when no source covered it. |
| `status` | yes | enum | `covered` (≥1 finding), `partial` (some aspects unaddressed; explain in `gap_note`), `uncovered` (no findings). |
| `gap_note` | no | string \| null | One sentence on what's missing. Required when `status != covered`. |

`coverage[]` MUST contain one row per sub-question in the plan, in plan
order. This is the outline `synthesize-report` will consume.

`unsupported_sources[]` lists sources that produced zero findings —
often a smell (off-topic source, broken body, model overlooked it).
Surface for downstream review; do not silently drop them.

## Confidence rubric

Qualitative ordinal scale: `high | medium | low`. The bucket is chosen by
the strength of the evidence relative to the claim. Full bucket
definitions, the calibration rationale (why ordinal, not numeric), and the
downstream-consumer notes live in
[`references/confidence-rubric.md`](references/confidence-rubric.md).
Never emit a numeric confidence score.

## Procedure

1. **Validate inputs.** Confirm `research_plan.sub_questions[]` exists
   and each entry has an `id`. Confirm `sources[]` is non-empty and each
   entry has an `id`, a `url` or `title`, and a body field. If a source
   lacks a body, add it to `unsupported_sources[]` with reason
   `"no body — could not extract"` and skip it in step 3. If
   `research_plan` is missing or empty, stop and surface the failure;
   this skill cannot generate an outline.

2. **Establish the iteration order.** Iterate sub-question by
   sub-question, in `research_plan.sub_questions[]` order. For each
   sub-question, scan all sources for passages that bear on it. This
   outline-conditioned pass (STORM's approach) is preferable to a
   source-by-source pass because it produces a directly-usable outline
   for `synthesize-report` with no reshuffling.

3. **Extract evidence first, then claims.** For each sub-question, scan
   the available sources and capture short evidence excerpts (≤ 600
   chars, in double quotes when verbatim) before writing any claim. This
   sequence is load-bearing: writing the claim first and then attaching
   a citation produces "post-rationalized" citations (the dominant
   citation-faithfulness failure in vanilla RAG, ~57% of cases per
   published studies). Evidence-first prevents that.

4. **Compose each claim from its evidence.** Each claim is one
   declarative sentence (≤ 240 chars) that the `evidence` excerpt
   directly supports. Preserve hedging from the source — if the source
   says "may reduce", the claim must not say "reduces". If the
   evidence does not support the claim as stated, rewrite the claim;
   do not weaken the evidence.

5. **Cap findings per sub-question.** Default ceiling: 7 findings per
   sub-question. If you have more candidates, rank by centrality to the
   sub-question and drop the tail. Over-extraction is a known failure
   mode ("information overload" in deep-research evals — more findings
   correlates with lower downstream factual accuracy).

6. **Detect corroboration and contradiction.**
   - When two or more sources independently support the same claim,
     emit one finding and list the others in `supporting_source_ids`.
     This is the strongest signal for `confidence: high`.
   - When two sources contradict each other on the same point, emit
     **both** claims as separate findings, each with the opposing
     source listed in `disputed_by`. Do not collapse contradictions
     into a hedged single claim; the synthesis stage and the user need
     to know the question is contested.

7. **Assign confidence per the rubric.** Apply
   [`references/confidence-rubric.md`](references/confidence-rubric.md) and
   tag each finding. When uncertain between two buckets, pick the lower one.

8. **Handle unmapped findings.** A claim that genuinely doesn't fit any
   sub-question but is too important to drop (e.g., it reframes the
   research question) goes to `sub_question_id: "unmapped"` with a
   one-sentence `note` explaining why it was retained. Use sparingly —
   unmapped findings should be rare; if you have many, the
   research_plan is probably incomplete and that's a `review-report`
   concern, not a place to dump notes.

9. **Build `coverage[]`.** One row per sub-question in plan order.
   Status:
   - `covered` — ≥1 finding addresses the sub-question.
   - `partial` — finding(s) present but key aspects unaddressed; write
     a `gap_note`.
   - `uncovered` — zero findings; `finding_ids: []`; write a `gap_note`
     naming what's missing or which sources should have addressed it.

10. **Build `unsupported_sources[]`.** Any source that produced zero
    findings gets a row with a one-line reason: off-topic,
    body-too-short, paywalled-stub, duplicate-of, etc. Empty list is
    fine and expected for a focused source set.

11. **Final pass — anti-pattern sweep.** Before emitting, check
    [`references/anti-patterns.md`](references/anti-patterns.md) and confirm
    none apply.

## Output contract

Emit a single JSON object matching the schema in `## Outputs`. The
object is the value bound to the workflow's `findings` output.

- Top-level keys: `findings`, `coverage`, `unsupported_sources`.
- All `source_id` and `supporting_source_ids` / `disputed_by` entries
  must reference an `id` present in input `sources[]`. Unknown ids are
  a hard fail.
- All `sub_question_id` entries must reference an `id` present in
  input `research_plan.sub_questions[]`, or the literal `"unmapped"`.
- `findings[]` may be empty only if every sub-question is `uncovered`;
  in that case `coverage[]` documents the gap and the workflow should
  loop back to `search-sources` rather than proceed to synthesis.

## Worked example

See `examples/four-day-week.md` for a complete worked input → output
illustration. It shows: evidence-first extraction, corroboration via
`supporting_source_ids`, contradictions surfaced as separate findings with
`disputed_by` pointers, the qualitative confidence rubric in use, and a
`partial`-coverage `gap_note` for an under-covered sub-question.

## Anti-patterns

Sweep [`references/anti-patterns.md`](references/anti-patterns.md) before
emitting `findings` — 10 failure modes covering invented source ids, post-
rationalized citations, stripped hedging, collapsed contradictions,
over-extraction, numeric confidence, "unmapped" dumping, silent source
skips, faithfulness↔correctness confusion, and sub-question reordering. If
any applies to a finding, fix the finding before emitting.

## Constraints

- DO NOT fetch new sources. This skill is a closed-context extractor;
  it cannot run search and must not pretend it did.
- DO NOT modify the `research_plan` or `sources` inputs. They are
  immutable upstream contracts.
- DO NOT emit findings whose `source_id` is not in input `sources[]`.
- DO NOT emit findings whose `sub_question_id` is not in input
  `research_plan.sub_questions[]` and is not the literal `"unmapped"`.
- DO NOT write prose paragraphs or commentary — output is the JSON
  object only. Synthesis is `synthesize-report`'s job.
- DO NOT exceed 7 findings per sub-question without explicit user
  override in the calling context.
- DO NOT use numeric confidence scores. The rubric is ordinal.
- DO NOT silently drop sources. Unproductive sources go in
  `unsupported_sources[]` with a reason.
- DO NOT touch `workflows/researcher.yaml` or any file outside
  `skills/extract-findings/`.

## Validation gate

Before emitting the output object, confirm:

1. `findings[]` is present (may be empty only if all sub-questions are
   `uncovered`).
2. Every `finding.source_id` is in `sources[].id`. Every entry in
   `supporting_source_ids` and `disputed_by` is in `sources[].id`.
3. Every `finding.sub_question_id` is in
   `research_plan.sub_questions[].id` or equals `"unmapped"`.
4. Every `finding.confidence` ∈ `{high, medium, low}`.
5. Every `finding.evidence` is a non-empty string ≤ 600 chars; verbatim
   quotes are wrapped in double quotes.
6. Every `finding.claim` is one sentence ≤ 240 chars and does not
   strengthen the hedging of its `evidence`.
7. `coverage[]` has one row per sub-question in `research_plan.sub_questions[]`
   in plan order; `status ∈ {covered, partial, uncovered}`; `gap_note`
   present whenever `status ≠ covered`.
8. `unsupported_sources[]` lists every source that produced no findings,
   each with a one-sentence `reason`.
9. No source_id appears in both `supporting_source_ids` and
   `disputed_by` for the same finding.
10. Anti-pattern sweep
    ([`references/anti-patterns.md`](references/anti-patterns.md)) — every
    applicable item is absent or explicitly mitigated.

If any gate fails, fix the offending row before emitting; do not
emit a partial output and rely on `review-report` to catch it — gate
failures here corrupt the input to two downstream stages.

## References

- [`references/confidence-rubric.md`](references/confidence-rubric.md) — the
  qualitative bucket definitions, why-ordinal-not-numeric rationale, and
  downstream-consumer notes.
- [`references/anti-patterns.md`](references/anti-patterns.md) — the
  10-row failure-mode sweep used by Procedure step 11 and Validation gate
  item 10.

Prior-art that shaped the rubric and procedure (Elicit, STORM, FACTUM, FACTS
Grounding, PICO, CER, deep-research surveys) is not loaded into the model
context.

## Troubleshooting

| Signal | Action |
|---|---|
| `sources` list contains an entry without a body field | Add to `unsupported_sources[]` with reason `"no body — could not extract"`. Do not fabricate body content. |
| Same finding extractable from many sources | Emit one finding; list the rest in `supporting_source_ids`. Upgrade `confidence` to `high` if ≥2 independent sources. |
| Two sources contradict each other on the same point | Emit two findings, each with the opposing source in `disputed_by`. Do not collapse. |
| Sub-question has no covering source | `coverage[].status: uncovered`; `finding_ids: []`; `gap_note` names what's missing. Flag clearly — workflow may want to loop back to `search-sources`. |
| Source clearly off-topic | Add to `unsupported_sources[]` with reason `"off-topic — no claims relevant to plan"`. |
| Hitting the 7-findings-per-sub-question cap repeatedly | The sub-question is probably too broad. Emit only the top 7 by centrality; add a `note` on the lowest-ranked one flagging the over-population for `review-report` to consider. |
| Findings that don't fit any sub-question but feel important | Sparingly: emit with `sub_question_id: "unmapped"` and a `note` explaining why retained. If you have ≥3 unmapped findings, that's a research-plan-coverage signal — flag in the highest-priority `coverage` row's `gap_note`. |
