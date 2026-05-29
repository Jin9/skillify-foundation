# Revision rules

The canonical decision rules and hard constraints for Step 5 of the Procedure.
After steps 1–4 have collected every observation, apply the decision table to
each observation, then enforce the hard constraints over every applied
revision. These are the authoritative copies; the SKILL.md body, Anti-patterns,
and Constraints sections all point here.

## Step 5 decision table

These rules are intentionally restrictive — the goal is minimum-edit, not
improvement-for-its-own-sake.

| Observation | Action |
|-------------|--------|
| `verified` | None (record in notes ledger). |
| `unsupported_by_cited_source` + swap candidate available | **Swap** the citation. Prose unchanged. |
| `unsupported_by_cited_source` + no swap | **Drop** the citation AND reword the sentence to remove the unsupported specificity (e.g., remove the precise number, soften the assertion). If the sentence cannot survive without the unsupported claim, **delete** the sentence. |
| `citation_broken` | Same rules as `unsupported_by_cited_source`. |
| Uncited claim + add candidate | **Add** the citation. Prose unchanged. |
| Uncited claim + no candidate (`unresolvable_gap`) | Default: **flag-only**, record the gap. Optional: soften the assertion if the prose admits a hedged form without losing meaning. Do NOT delete unless the claim is clearly wrong. |
| `off_topic` section, minor | Flag-only. |
| `off_topic` section, major (whole subsection drifts) | Flag-only — note in review_notes. Do not delete entire subsections in a self-review pass. |
| `audience_mismatch_pervasive` | Flag-only, top-level note. |
| `audience_mismatch_minor` | Flag-only. |
| `internal_contradiction`, resolvable | Replace the weaker-evidence side; record both. |
| `internal_contradiction`, unresolvable | Flag-only. |
| `diagram_inconsistent`, minor (one label off, arrow direction wrong) | Flag-only. The reviewer MUST NOT rewrite a diagram's contents — that risks introducing claims not in findings. |
| `diagram_inconsistent`, severe (whole structure contradicts the prose, or the diagram sits in a forbidden section, or contains an `[n]` marker) | **delete_diagram** — remove the fenced block. Surrounding prose is unchanged. Record the deletion in `Applied changes`. |

## Hard constraints on revision

Canonical list — the Anti-patterns and Constraints sections of SKILL.md point
here. The skill's anti-patterns are the inverse of these constraints.

- **No new claim.** No revision may introduce a claim that is not already in
  either the draft or `findings`.
- **No new source.** No revision may emit a `source_id` that is not in
  `cited_sources ∪ {f.source_id for f in findings}`. No fabricated sources.
- **No paraphrase of passing prose.** Style edits for their own sake are
  forbidden — the documented destructive-rewrite failure mode.
- **No fenced-block edits.** The only permitted operation on a diagram is
  `delete_diagram` (whole-block removal). Rewriting a diagram's labels,
  arrows, or layout would smuggle in claims the synthesize stage didn't
  ground.
- **Same top-level structure.** The final report must end with the same
  top-level sections as the draft unless a section was entirely flagged for
  removal. No section deletion in self-review — flag instead.
- **Single-pass.** Exactly one review iteration. No loops, no recursive
  calls, no re-search. Surface remaining issues in `review_notes` and stop.
- **No hedging-to-placate.** Do not add "it should be noted that…" framing
  purely to satisfy the rubric; grounding quality is what is scored.
- **No scope-widening.** Do not fix formatting choices, section ordering, or
  prose voice that fall outside the observation set from steps 2–4.
- **No empty `review_notes`.** The verified-claims ledger is required even
  on a clean draft.
- **No silent citation drops.** Every drop appears in `Applied changes`.
- **Minimum-edit.** Every revision traces to a specific observation from
  steps 2–4; prose that passes every check is left untouched.
- **Audience re-targeting is out of scope.** Flag pervasive mismatch and
  stop; do not rewrite for a different audience in this pass.
