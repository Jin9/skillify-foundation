# Reference: convergence discipline

The point of a multi-model debate is *principled* movement: positions should
converge when the evidence forces it and hold when it does not. Two
pathologies destroy the signal — this reference defines them precisely so
`revise-positions` can avoid them and `report-debate` can flag them.

## The four actions

| `action` | Use when | Required fields |
|----------|----------|-----------------|
| `conceded` | The rebuttal is correct and load-bearing; the claim cannot stand. | `what_changed` describes the removal; the claim is gone/withdrawn in the new `position`. |
| `partially_conceded` | The rebuttal is partly right; the claim survives only narrowed or caveated. | `what_changed` states the new, narrower claim. |
| `reframed` | The underlying insight is sound but the *phrasing* invited the objection. | `what_changed` states the reframing; the substantive position is unchanged. |
| `defended` | The objection fails on the evidence or logic. | `evidence_refs` non-empty (`[g-n]`); `what_changed == ""`. A defense is an argument, not a restatement. |

## Capitulation (forbidden)

Conceding a point that was **not** under real rebuttal pressure — no
`major`/`decisive` objection actually targeted it — purely to appear
agreeable or to converge faster. Symptoms: `stance_delta: reversed` with a
changelog full of `conceded` actions whose `in_response_to.objection` were
all `minor`; conceding a claim no peer attacked at all.

Why it is banned: it manufactures fake consensus. A debate whose value is
"three frontier models independently agree" is worthless if the agreement
was social, not evidential. The open moderator is told to discount
agreement that arrived via capitulation.

## Flip-flop (forbidden)

Reversing a position you **previously defended well** without a *new*
landed objection. Symptoms: round 1 `defended` claim X with `[g-3]`; round 2
`conceded` X with no new rebuttal introduced between the rounds. Movement
must be caused by an argument that appeared since your last position, not by
drift.

Why it is banned: it makes the convergence map noise. Stable, evidence-anchored
positions that only move on new pressure are what let synthesis trust the
final state.

## Principled convergence (the goal)

Movement is healthy when every `conceded`/`partially_conceded` traces to a
specific `major`/`decisive` objection in `critiques_received`, and every
`defended` carries real `[g-n]` (or correctly-marked logic) rebutting the
objection. A debate where all three positions narrow toward the same
evidence-supported answer through traceable concessions is the best
outcome — and it looks completely different from capitulation in the
changelog, which is exactly why the changelog is mandatory and audited.

## stance_delta calibration

- `unchanged` — only `defended`/`reframed` actions; substantive stance
  identical.
- `narrowed` — one or more `partially_conceded`; the claim is now more
  qualified/scoped.
- `broadened` — a peer's grounded point legitimately extended your position
  (rare; needs a changelog entry citing the peer).
- `reversed` — a `conceded` on a load-bearing claim flipped your answer to
  the proposition. Only valid against a `major`/`decisive` objection;
  `reversed` against only `minor` objections is the capitulation signature.
