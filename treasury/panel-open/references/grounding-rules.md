# Reference: grounding rules (shared by panel-open, cross-examine, revise-positions)

Every panel stage argues inside one **closed evidence world**: the
`grounding_pack` that `frame-debate` produced. This is the debate analogue
of squad-researcher's closed-world citation rule, and it is what stops a
verbose or confidently-hallucinating model from "winning" by importing
private facts.

## The [g-n] contract

- The only citable evidence is `grounding_pack.items[]`, addressed by its
  `id` in square brackets: `[g-3]`, `[g-7]`. Multiple: `[g-2][g-5]` (no
  comma, no range).
- Place the marker **adjacent to the claim it supports**, not at paragraph
  end. Citation recall collapses when markers cluster — the same
  fluency-vs-precision finding squad-researcher's `synthesize-report` cites.
- A `[g-n]` you reference MUST exist in `grounding_pack.items`. Inventing a
  `g-` id is the single worst failure in this workflow — it is an
  unfalsifiable fake citation. The Validation gate of every panel skill
  checks this.
- Diagrams/ASCII are not used in panel positions. Prose only.

## The honesty ledger (`uncited_assertions`)

Real arguments sometimes need a step the grounding pack does not cover (a
definition, a logical inference, a widely-known fact). That is allowed — but
it must be **declared**, not disguised:

- Any factual claim not traceable to a `[g-n]` item goes into the position's
  `uncited_assertions` array, verbatim or closely paraphrased.
- An honest "I assert X without grounding-pack support" is strictly better
  than a `[g-n]` that does not actually support X. The latter is detected
  and scored as a rubric failure (evidential-grounding criterion).
- Pure-logic rebuttals in `cross-examine` are legitimate but must be marked
  (`evidence_refs: []` + flagged as logic-only) — never dressed up as
  evidential.

## Out-of-scope guard

`grounding_pack.out_of_scope` lists what the original research deliberately
excluded. Do not drag the debate there: a position that wins on an
out-of-scope tangent is off-topic and loses the "addresses the proposition"
rubric criterion.

## Degraded grounding

If `grounding_pack.items` is empty (degraded mode, no `source_text`), the
debate is reasoning-only. Argue carefully, put essentially everything in
`uncited_assertions`, and do not fabricate `[g-n]` ids to look grounded.
`synthesize-consensus` already caps confidence at `medium` for degraded
runs; faking grounding does not raise it, it just corrupts the transcript.
