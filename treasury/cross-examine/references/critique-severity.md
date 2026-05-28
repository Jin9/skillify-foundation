# Reference: critique severity + logic-only vs evidential

`cross-examine` rates every rebuttal. The severity is what
`revise-positions` and `synthesize-consensus` use to decide whether a
position must move — so calibrate honestly, not theatrically.

## Severity scale

| `severity` | Meaning | Effect downstream |
|------------|---------|-------------------|
| `minor` | The point is weak or imprecise but the conclusion survives. Wording, a missing caveat, a non-load-bearing over-claim. | `revise-positions` may address it; not required to. |
| `major` | A load-bearing claim is unsupported, misuses a `[g-n]` item, or has a real logical gap. The conclusion is at risk but recoverable with work. | `revise-positions` MUST respond (concede / partially concede / defend with new `[g-n]`). On `deep`, a changelog entry is required. |
| `decisive` | If correct, the target's stance cannot stand as written — a contradiction with a high-confidence `[g-n]` item, a self-contradiction, or an out-of-scope win. | Target must concede, reframe, or mount a grounded defense; an ignored decisive rebuttal is flagged by `report-debate`'s integrity check. |

Inflating `minor` to `decisive` to "win" is itself a rubric failure
(evidential-grounding + internal-consistency criteria) and is visible to the
open moderator in synthesis. Under-rating a real `decisive` to be agreeable
corrupts the convergence map. Rate what the evidence supports.

## Evidential vs logic-only rebuttals

- **Evidential** — `evidence_refs` names ≥1 `[g-n]` item that, read as
  written, contradicts or undercuts the `target_claim`. `logic_only: false`.
  Strongest form of rebuttal.
- **Logic-only** — the objection is a reasoning flaw (invalid inference,
  equivocation, base-rate neglect, non-sequitur) with no grounding-pack
  item behind it. `evidence_refs: []`, `logic_only: true`. Fully
  legitimate, often `major`, occasionally `decisive` (a genuine
  self-contradiction is decisive without any evidence). It MUST be marked —
  a logic-only rebuttal presented as evidential is detected and scored
  against the author.
- **Not allowed** — citing a `[g-n]` that does not actually support the
  objection (fake evidential). Worse than an honest logic-only.

## Concessions are not optional

A target entry with zero `concessions` is treated as presumptive
strawmanning. If a peer's position is genuinely all wrong (rare), say so in
one explicit `concessions` line ("No point survives: every key_claim rests
on g-2 which it misreads — see rebuttals") rather than leaving the array
empty. The mandatory `strongest_point_from_target` still applies even then:
state the best version of their case before you dismantle it.
