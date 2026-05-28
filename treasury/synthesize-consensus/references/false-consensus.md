# Reference: false consensus and the strongest-case-each-side rule

The failure mode that destroys a debate's value at the last step is
**laundering a real disagreement into a bland agreed sentence**. It makes
the output look confident and unanimous when the panel actually split. This
is the squad-brainstorm analogue of `extract-findings`' "do not collapse
contradictions" rule, applied to the final merge.

## The laundering anti-pattern

Symptom: a `live_disagreement` that existed through the final round is
written up in `## Where the panel agrees`, or as a single hedged sentence
("the panel broadly felt DDD is situational") that erases which conditions
each side actually argued and which evidence backed them.

Fix: if ≥1 panelist held a position in the final round that ≥1 other
panelist argued against and did not concede, it is a **live disagreement**.
It goes in `## Live disagreements` with the strongest case for *each* side,
each with its `[g-n]`. Only then may the moderator add a *lean* — clearly
labeled as the moderator's call, with a grounded rationale — or declare it
`genuinely_unresolved`.

## Strongest-case-each-side

For every live disagreement, write each side the way its best advocate
would, using its actual `[g-n]`. A side summarized weakly so the lean looks
obvious is a strawman and a rubric failure (the same anti-strawman principle
that governs the panel stages, now applied to the moderator).

A `moderator_lean` is legitimate only when the `debate_brief.rubric` +
`[g-n]` evidence favor one steelmanned side over the other steelmanned side.
"This side's model is usually better" is never a valid rationale.

## Reading the changelog: real vs fake convergence

When panelists end up agreeing, check *how* they got there before calling it
strong consensus:

- **Independent agreement** (`via: "independent"`) — they agreed from the
  opening round, no movement needed. Strong signal.
- **Principled convergence** (`via: "principled_convergence"`) — a position
  moved, and `final_positions[].changelog` shows the concession traces to a
  `major`/`decisive` objection with a real argument. Strong signal.
- **Capitulation-driven agreement** — the changelog shows `conceded`/
  `reversed` against only `minor` objections, or a concession to a claim no
  one strongly attacked. This is *social* agreement, not evidential. Note it
  in `convergence_summary` and DO NOT present it as robust consensus;
  `overall_confidence` should reflect the weakness.

A run where all three converge via capitulation is *less* trustworthy than
one with an honest, well-argued live disagreement. Say so plainly.
