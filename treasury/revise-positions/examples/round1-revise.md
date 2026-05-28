# Worked example: revise-positions (round 1)

Continuing the DDD `standard` debate. Round 1 revision, `panelist = P-gemini`.

The driver passed P-gemini only the critiques whose
`target_panelist_id == "P-gemini"`. From the `cross-examine` example,
P-claude's critique of P-gemini contained:
- A **decisive** evidential rebuttal: base-rate reasoning needs g-1 and g-2
  to be the same population; g-1 is conditioned on stabilization, so 58% is
  not a valid prior for g-1's population.
- A **major** logic-only rebuttal: the tie-break principle smuggles an
  unstated independence assumption.
- Steelman: even granting the objection, defenders must show the stabilized
  population isn't vanishingly small.

(Assume P-codex's critique of P-gemini, also in `critiques_received`, added
one `minor` wording objection — shown conceded-as-reframe below.)

```json
{
  "panelist_id": "P-gemini", "cli": "gemini -p", "model": "gemini-2.5-pro",
  "round": 1,
  "stance": "On the current evidence the proposition is unproven, not false: the cost-raising reading only holds if the stabilized population is small, which the pack does not establish.",
  "position": "I withdraw the strong claim that base rates settle this against the proposition. P-claude is right that g-1 is explicitly conditioned on the model stabilizing [g-1], so g-2's 58% rollback rate [g-2] is not a valid prior over g-1's population — they may be different teams. What survives is weaker but still pointed: the proposition's defenders owe a quantification. g-1 gives no base rate for how often stabilization is reached, and the only high-confidence item [g-3] is about boundary quality, not cost. So the honest state is 'unproven pending the size of the stabilized population', not 'true more often than not' ...",
  "answers": [
    { "cq_id": "cq1", "answer": "Indeterminate on the pack: lowers TCO for the stabilized population, raises it for the rollback population, and the pack does not size either.", "evidence_refs": ["g-1","g-2"] },
    { "cq_id": "cq2", "answer": "Conceded toward P-codex/P-claude: g-3 makes the strategic core the only high-confidence anchor; my 'under-evidenced for cost specifically' point stands but is now a scoping caveat, not a rebuttal.", "evidence_refs": ["g-3"] },
    { "cq_id": "cq3", "answer": "Unchanged: the pack still does not isolate the abandonment cause; g-2 is a rate, not a trigger.", "evidence_refs": ["g-2"] }
  ],
  "key_claims": [
    { "claim": "The proposition is unproven (not false) absent a stabilization base rate.", "grounding": ["g-1","g-2"], "confidence": "medium" },
    { "claim": "The strategic core is the only high-confidence anchor.", "grounding": ["g-3"], "confidence": "high" }
  ],
  "changelog": [
    {
      "in_response_to": { "from_panelist_id": "P-claude", "objection": "Base-rate tiebreak requires g-1 and g-2 to be the same population; g-1 is conditioned on stabilization." },
      "action": "conceded",
      "what_changed": "Dropped 'base rates favor the cost-raising side'. The 58% is no longer treated as a prior over g-1's population.",
      "evidence_refs": []
    },
    {
      "in_response_to": { "from_panelist_id": "P-claude", "objection": "Tiebreak principle smuggles an unstated independence assumption (logic-only, major)." },
      "action": "conceded",
      "what_changed": "Removed the general 'ties break on base rates' principle from the argument entirely.",
      "evidence_refs": []
    },
    {
      "in_response_to": { "from_panelist_id": "P-claude", "objection": "Steelman: defenders must show the stabilized population is not vanishingly small." },
      "action": "defended",
      "what_changed": "",
      "evidence_refs": ["g-1"]
    },
    {
      "in_response_to": { "from_panelist_id": "P-codex", "objection": "Minor: 'dead heat' wording overstated the symmetry of g-1 vs g-2." },
      "action": "reframed",
      "what_changed": "Replaced 'dead heat' with 'opposite-signed, both medium-confidence, populations unestablished'."
    }
  ],
  "stance_delta": "narrowed",
  "uncited_assertions": ["A claim of the form 'X is true more often than not' carries a burden to quantify the base rate of X's enabling condition (argumentative burden principle, not a pack item)."],
  "status": "ok"
}
```

Why this is principled, not capitulation:
- Each `conceded` traces to a specific `decisive`/`major` objection that
  actually targeted P-gemini — not social agreement.
- The surviving point is **defended** with `[g-1]` (the burden-of-quantification
  argument), not abandoned.
- `stance_delta: narrowed` is consistent with the changelog (concessions on
  the tiebreak, defense of the weaker surviving claim) — a genuine
  convergence toward "unproven", which is where the evidence actually sits.
