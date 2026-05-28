# Reference: the open-moderator no-favoritism check

`squad-brainstorm` deliberately uses an **open** moderator (project
decision): `synthesize-consensus` sees `panelist_id` on every position and
critique, and the stage itself runs on a Claude model while one panelist is
`P-claude`. There is no anonymization layer. The bias control is therefore
behavioral and audited, not structural.

## What the check must assert

`convergence.no_favoritism_check` is a mandatory, non-empty string. It must:

1. State that every `moderator_lean` was decided on `debate_brief.rubric`
   criteria and `[g-n]` evidence, not on which model authored a side.
2. For **each** `live_disagreement` with a non-`genuinely_unresolved` lean,
   name the specific rubric criterion id and/or `[g-n]` item that drove the
   lean. "Side-2 wins on r3 (engages counter-evidence) + g-3" — not
   "side-2 was more convincing".
3. Confirm no agreement was upgraded, and no disagreement downgraded, on the
   basis of panelist identity.

Example (good):
> No-favoritism: cq1 leans side-2 on r2 (names ≥2 concrete payoff
> conditions) + g-3's high-confidence boundary result; cq3 stays
> genuinely_unresolved (no [g-n] isolates the cause). Leans were checked
> against held_by after deciding on evidence; the cq1 lean coincides with
> P-codex+P-claude but is justified solely by g-3's confidence tier, which
> would hold regardless of who argued it.

## The P-claude defense-in-depth rule

Whenever a `moderator_lean` coincides with the side `P-claude` held, the
`lean_rationale` and `no_favoritism_check` must be **extra explicit** about
the evidential basis: name the `[g-n]` and rubric criterion such that a
skeptical reader could verify the lean would be identical if P-gemini or
P-codex had argued that exact side with that exact evidence. This directly
counters the known risk of an open Claude moderator over-crediting the
Claude panelist.

If the honest reading is that the strongest evidence genuinely sits with the
side P-claude argued, that is a valid outcome — say so, and make the
evidence chain auditable. The rule is "make the basis checkable", not
"never agree with P-claude".

## What report-debate does with this

`report-debate`'s pipeline-integrity section recomputes whether the set of
`moderator_lean`s correlates with `P-claude`'s positions across all
contested questions and surfaces it as a check (e.g. "moderator leaned to
P-claude's side on 3/3 — review no_favoritism_check"). It does not re-judge;
it makes the pattern visible so a human can audit. A consistently
P-claude-favoring run with thin rationales is the signal this control exists
to catch.
