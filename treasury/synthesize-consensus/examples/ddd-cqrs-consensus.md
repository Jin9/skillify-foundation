# Worked example: synthesize-consensus

The DDD `standard` debate, after round 1 revision. The moderator sees
attributed final positions + the cross-exam exchanges.

Final state (abbreviated): after round 1, P-gemini narrowed to "unproven
pending the stabilized-population base rate" (defended with g-1); P-codex and
P-claude both hold "true-conditional: lowers TCO under strategic-first
adoption", P-claude via the different-populations reconciliation. All three
agree the strategic core (g-3) is the only high-confidence anchor.

### Output — `convergence`

```json
{
  "proposition": "Domain-Driven Design lowers the total cost of ownership of enterprise systems more often than it raises it.",
  "panelists": ["P-codex","P-gemini","P-claude"],
  "agreements": [
    { "point": "The strategic core of DDD (bounded contexts, context maps) is the only high-confidence cost-relevant anchor.",
      "held_by": ["P-codex","P-gemini","P-claude"], "grounding": ["g-3"], "confidence": "high", "via": "principled_convergence" },
    { "point": "The two direct cost findings are opposite-signed and both medium-confidence; neither alone settles the proposition.",
      "held_by": ["P-codex","P-gemini","P-claude"], "grounding": ["g-1","g-2"], "confidence": "medium", "via": "independent" }
  ],
  "live_disagreements": [
    {
      "cq_id": "cq1",
      "summary": "Whether the proposition is 'true-conditional' (holds under strategic-first adoption) or 'unproven' (pending the size of the stabilized population).",
      "sides": [
        { "held_by": ["P-codex","P-claude"], "position": "True as a conditional: g-1 and g-2 plausibly describe different populations, so DDD lowers TCO for teams that adopt strategic design and gate tactical ceremony.",
          "strongest_case": "g-3 (high confidence) shows the strategic core is robust independent of tactical churn; g-1's saving is real for the population that stabilizes; nothing forces g-1 and g-2 onto the same teams.", "grounding": ["g-1","g-3"] },
        { "held_by": ["P-gemini"], "position": "Unproven, not false: the conditional only matters if the stabilized population is non-trivial, which the pack never sizes.",
          "strongest_case": "g-1 gives no base rate for reaching stabilization; g-2 shows a 58% majority rolling back; a 'true more often than not' claim carries a burden to quantify the enabling condition that the pack cannot discharge.", "grounding": ["g-1","g-2"] }
      ],
      "moderator_lean": "genuinely_unresolved",
      "lean_rationale": "Rubric r1 (grounded) and r3 (engages counter-evidence) are satisfied by BOTH steelmanned sides. The deciding fact — the size of the stabilized population — is exactly what no [g-n] item supplies (g-1 is conditional, g-2 is a rate over a possibly-different population). Per r1, a lean would require evidence the pack does not contain; therefore unresolved."
    }
  ],
  "overall_confidence": "medium",
  "confidence_rationale": "One high-confidence agreement (g-3) but the headline proposition turns on a base rate absent from the grounding pack; the panel converged on the structure of the disagreement, not on an answer.",
  "convergence_summary": "Partial convergence: all three aligned on the strategic-core anchor (principled — P-gemini moved via a g-1-defended concession, not capitulation); the cost-direction question stayed an honest live disagreement through round 1.",
  "grounding": "full",
  "no_favoritism_check": "cq1 was left genuinely_unresolved strictly on rubric r1 (no [g-n] sizes the stabilized population) — a verdict that goes AGAINST the P-codex+P-claude majority and against the side P-claude argued, demonstrating the lean tracks evidence, not identity. The single agreement upgrade to 'high' rests only on g-3's confidence tier, which is identity-independent."
}
```

### Output — `answer` (abbreviated markdown)

```markdown
# Does Domain-Driven Design lower total cost of ownership in enterprise systems more often than not?

## Bottom line
The panel could not settle this on the available evidence, and that is the
honest answer: DDD's strategic core clearly pays off [g-3], but whether the
*headline* cost claim holds depends on how often teams actually reach the
"model stabilized" state — a number the research never establishes. Treat
"DDD lowers TCO" as conditional and unproven at the population level, not as
a default.

## Where the panel agrees
- The strategic core (bounded contexts, context maps) is the only
  high-confidence cost-relevant anchor [g-3] — all three converged here.
- The two direct cost findings point opposite ways and are both only
  medium-confidence [g-1][g-2]; neither settles the question alone.

## Live disagreements
**Is the proposition "true-conditional" or "unproven"?** (cq1)
- *True-conditional (P-codex, P-claude):* g-1 and g-2 likely describe
  different populations; the strategic core is robust [g-3], so strategic-first
  adoption captures the saving [g-1].
- *Unproven (P-gemini):* the conditional only matters if the stabilized
  population is non-trivial; g-1 gives no base rate, g-2 shows 58% rollback
  [g-1][g-2].
- *Moderator: genuinely unresolved.* Both cases are well-grounded; the
  deciding fact (size of the stabilized population) is absent from the
  evidence, so taking a side would require inventing it.

## Confidence & caveats
Medium. Strong agreement on the strategic core; the headline proposition
turns on a missing base rate. Convergence on the strategic point was
principled (P-gemini moved on a defended argument, not capitulation).

## Grounding
- [g-1] Teams adopting DDD report lower long-run change cost once the model stabilizes — origin finding:f-12, confidence medium
- [g-2] A majority of surveyed teams abandoned strict DDD within 18 months citing ceremony — origin finding:f-31, confidence medium
- [g-3] DDD's bounded contexts map cleanly onto service boundaries — origin finding:f-44, confidence high
```

This is the ideal shape: a real live disagreement preserved with both
strongest cases, a lean that goes *against* the Claude panelist's side
purely on the rubric (the no-favoritism check makes that explicit and
verifiable), and confidence honestly capped by what the evidence can support.
