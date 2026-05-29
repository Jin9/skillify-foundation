# Output field definitions

Authoritative per-field rules for the `findings` output object emitted by
`extract-findings`. The JSON skeleton lives in `SKILL.md` (`## Outputs`);
this file defines every field's type and rule.

## Input field shapes

The workflow contract (`workflows/researcher.yaml`, stage `extract-findings`)
passes exactly two inputs, both as-is from upstream stages. The skill must not
fabricate fields the upstream stages did not provide; if a source arrived
without a body, treat it as unreadable (see the SKILL.md validation gate).

| Field | Required | Source | Shape |
|---|---|---|---|
| `research_plan` | yes | `stages.plan-research.research_plan` | object with `sub_questions[]`; each entry has at least `id` (stable, e.g. `sq-1`) and `question` (text). May also carry `rationale`, `priority`. |
| `sources` | yes | `stages.search-sources.sources` | list of objects; each entry has at least `id` (stable, e.g. `s-1`), `title`, `url`, and a body field (`text`, `content`, or `excerpt`) the model can read. May also carry `author`, `published_at`, `accessed_at`. |

## Per-finding fields

| Field | Required | Type | Rule |
|---|---|---|---|
| `id` | yes | string | Stable, monotonic (`f-1`, `f-2`, ...). Used by downstream stages to reference the finding. |
| `sub_question_id` | yes | string | Must match an `id` in `research_plan.sub_questions[]`, OR the literal `"unmapped"`. |
| `claim` | yes | string | One declarative sentence. ≤ 240 chars. No hedging beyond what the source uses. |
| `evidence` | yes | string | Verbatim excerpt or close paraphrase from the source body, ≤ 600 chars. Quoted material in double quotes. |
| `source_id` | yes | string | Must match an `id` in input `sources[]`. The primary source for this claim. |
| `supporting_source_ids` | no | string[] | Other source ids that independently support the claim. Use for corroboration. |
| `disputed_by` | no | string[] | Source ids that contradict the claim. When non-empty, also emit the opposing claim as a separate finding. |
| `confidence` | yes | enum | `high` \| `medium` \| `low`. Per the rubric in `references/confidence-rubric.md`. |
| `note` | no | string \| null | Caveats, missing nuance, unmapped-bucket explanation, or other context. |

## Per-coverage fields

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

## Why the output is a structured list (not freeform)

`extract-findings` emits a structured, extended findings list rather than
freeform notes:

- `synthesize-report` and `review-report` both need to address individual
  claims (cite them, flag them, drop them). Freeform notes force them to
  re-parse text. Structured rows are addressable.
- `sub_question_id` makes the output **outline-conditioned** (STORM's term):
  the synthesis stage gets a free outline from `coverage[]` and never has to
  reorganize raw notes.
- `supporting_source_ids` and `disputed_by` let corroboration and
  contradiction survive into synthesis instead of being collapsed.
- `confidence` is qualitative (`high|medium|low`), not numeric — see
  `references/confidence-rubric.md`. Verbalized numeric confidence from LLMs
  is poorly calibrated.

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
