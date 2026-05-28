# Worked example: cross-examine (round 1)

Continuing the DDD `standard` debate. Round 1. This is the
`panelist = P-claude` instance critiquing its two peers `P-codex` and
`P-gemini` (the driver passed only those two as `peer_positions`; P-claude's
own opening is not in scope here).

Recall the openings (from `opening-debate-panel` example):
- **P-codex**: conditional-yes — strategic design pays, tactical ceremony is
  the risk.
- **P-gemini**: evidence-says-no — base rates favor the cost-raising side.

```json
{
  "author_panelist_id": "P-claude",
  "round_no": 1,
  "critiques": [
    {
      "target_panelist_id": "P-codex",
      "concessions": [
        "Agree the high-confidence anchor [g-3] is a strategic-design result, so 'strategic value is the durable part' is well-grounded.",
        "Agree tactical ceremony is the plausible locus of the abandonment signal [g-2]."
      ],
      "rebuttals": [
        {
          "target_claim": "DDD lowers TCO 'only when strategic design lands before the team burns out on tactical ceremony'.",
          "objection": "This treats the ordering (strategic-before-tactical) as a free variable, but nothing in the pack shows teams can choose that ordering — g-2 only reports that rollback happened, not that a strategic-first sequence was available or attempted.",
          "evidence_refs": ["g-2"],
          "logic_only": false,
          "severity": "major"
        }
      ],
      "unsupported_claim_flags": [
        "'Getting service boundaries wrong is the expensive failure mode' — P-codex listed this in uncited_assertions, so it is honestly flagged by them; not a violation, noted for the moderator."
      ],
      "strongest_point_from_target": "If the abandonment in g-2 is specifically tactical-pattern rollback while the strategic core (g-3) stays intact, then a strategic-only adoption captures g-1's savings without paying g-2's cost — a coherent reading of the proposition as true-conditional."
    },
    {
      "target_panelist_id": "P-gemini",
      "concessions": [
        "Agree the direct cost evidence is a medium-vs-medium tie pointing opposite ways [g-1][g-2]."
      ],
      "rebuttals": [
        {
          "target_claim": "Base rates favor the cost-raising side because g-2 reports a 58% majority rollback.",
          "objection": "Base-rate reasoning requires g-1 and g-2 to describe the SAME population. g-1 is explicitly conditioned on 'once the model stabilizes'; g-2's rollback teams may be exactly the teams whose model never stabilized. Treating 58% as the prior for g-1's population is a base-rate misapplication, not a valid tiebreak.",
          "evidence_refs": ["g-1","g-2"],
          "logic_only": false,
          "severity": "decisive"
        },
        {
          "target_claim": "When direct evidence ties, base rates should break the tie.",
          "objection": "The principle is fine in the abstract but here it smuggles in an unstated independence assumption between g-1 and g-2; flagged as a reasoning gap, no pack item needed.",
          "evidence_refs": [],
          "logic_only": true,
          "severity": "major"
        }
      ],
      "unsupported_claim_flags": [],
      "strongest_point_from_target": "Even granting the population objection, the burden is on the proposition's defenders to show the stabilized-population is not vanishingly small — P-gemini is right that an unquantified 'it works if it stabilizes' is weak if stabilization is rare."
    }
  ],
  "status": "ok"
}
```

Note: concessions stated first for both targets; one `decisive` evidential
rebuttal (the base-rate/population flaw — if right, P-gemini's stance can't
stand as written) and one `logic_only` rebuttal correctly marked; the
mandatory steelman is a *genuinely strong* restatement of each opponent's
best case, not a throwaway. P-claude never revises its own position here.
