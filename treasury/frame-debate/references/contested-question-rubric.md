# Reference: what makes a question genuinely contestable

`frame-debate` lives or dies on the quality of `contested_questions`. A weak
question wastes three frontier-model panelist calls per round. Use this
rubric in Procedure step 5.

## A question is debate-worthy iff ALL hold

1. **Two defensible answers exist** given the grounding pack. If the evidence
   forces one answer, it is settled — not contested. Move it to the
   `summary`, not to `contested_questions`.
2. **It is not leading.** The phrasing must not encode the answer.
   - Bad: "Why is the modular monolith the right default for most teams?"
   - Good: "When does a modular monolith stop being the right default —
     team size, deploy independence, or domain volatility?"
3. **It is grounded.** It ties to ≥1 `grounding_pack.items` entry via
   `evidence_refs` (full mode). A question no finding speaks to is
   speculation; reserve `seed: synthesizer_gap` for tensions the report
   *raises but under-argues*, not for invented topics.
4. **It is decidable by argument.** A panelist could change another's mind
   with evidence + reasoning. Pure value/taste questions ("is elegant code
   worth more than fast code?") are not debate-worthy here.
5. **It is single-axis.** One question = one tension. Split compound
   questions ("is X faster AND cheaper AND safer?") into separate `cq`s or
   pick the load-bearing axis.

## Converting a `disputed_by` finding pair → one neutral question

This is the highest-value seed (`seed: finding_dispute`) — the squad-researcher
`extract-findings` stage already did the work of detecting the contradiction.

1. Take finding `f-A` (claim) and the finding(s) in its `disputed_by` (`f-B`).
2. Both become grounding-pack items (`g-A`, `g-B`) with `disputed: true`.
3. Write the question as the *axis of disagreement*, attributing neither side:
   - `f-A`: "DDD reduces long-term maintenance cost in enterprise systems."
   - `f-B`: "DDD's ceremony increases cost and most teams abandon it."
   - → `cq`: "Does Domain-Driven Design lower or raise total cost of
     ownership in enterprise systems, and what determines which?"
   - `why_contested`: "grounding items g-A and g-B directly disagree"
   - `evidence_refs`: `["g-A","g-B"]`
4. Never resolve the dispute in the question. The panel resolves it.

## Priority order when you have more candidates than the Rounds dial allows

`finding_dispute` > `plan_open_question` > `synthesizer_gap` >
`degraded_inferred`. Within a tier, prefer questions whose evidence items
have `confidence: low` or `disputed: true` — those are where a debate adds
the most over the original report. Drop the rest (do not pad to hit the
upper bound; the dial maxima are caps, not targets).
