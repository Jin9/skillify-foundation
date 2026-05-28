# Worked example: opening-debate-panel (panel-open fan-out x3)

Same brief (the DDD `standard` example from `frame-debate`), answered
independently by all three CLIs. The driver dispatches this skill three
times in parallel; each instance sees only the brief + grounding pack + its
own `panelist` row, never a peer's output.

Proposition under debate: *"Domain-Driven Design lowers the total cost of
ownership of enterprise systems more often than it raises it."*

---

## Instance 1 — `panelist = P-codex`

```json
{
  "panelist_id": "P-codex", "cli": "codex exec", "model": "gpt-5.5",
  "round": 0,
  "stance": "DDD lowers total cost of ownership only when strategic design is adopted and tactical ceremony is kept optional.",
  "position": "The cost case for DDD is real but conditional. Strategic design — bounded contexts and context maps — shows consistent payoff: migration retrospectives report clean service boundaries [g-3], which is the expensive part to get wrong. The long-run change-cost reduction [g-1] is plausible but it is gated on the model stabilizing, which is exactly where the 18-month abandonment signal [g-2] bites. So the proposition holds, but only under a discipline that treats tactical patterns as discretionary ...",
  "answers": [
    { "cq_id": "cq1", "answer": "Lowers TCO when strategic design lands before the team burns out on tactical ceremony; raises it when teams adopt tactical patterns wholesale first.", "evidence_refs": ["g-1","g-2","g-3"] },
    { "cq_id": "cq2", "answer": "Yes — the durable value is strategic; g-3 is the only high-confidence item and it is a strategic-design result.", "evidence_refs": ["g-3"] },
    { "cq_id": "cq3", "answer": "Ceremony stops paying off when the domain is simple or volatile enough that the model never stabilizes (implied by g-2's rollback timeframe).", "evidence_refs": ["g-2"] }
  ],
  "key_claims": [
    { "claim": "Strategic design is the load-bearing source of DDD's cost benefit.", "grounding": ["g-3"], "confidence": "high" },
    { "claim": "Tactical-pattern ceremony drives the abandonment signal.", "grounding": ["g-2"], "confidence": "medium" }
  ],
  "strongest_counterargument_acknowledged": "If 58% of teams abandon strict DDD within 18 months [g-2], the average team may never reach the stabilization point where g-1's savings appear — making DDD a net cost for the median team.",
  "uncited_assertions": ["Getting service boundaries wrong is the expensive failure mode (general engineering claim, not in the pack)."],
  "status": "ok"
}
```

## Instance 2 — `panelist = P-gemini`

```json
{
  "panelist_id": "P-gemini", "cli": "gemini -p", "model": "gemini-2.5-pro",
  "round": 0,
  "stance": "On the surveyed evidence the proposition is not supported: DDD raises TCO at least as often as it lowers it.",
  "position": "The two cost-direction findings are both medium-confidence and they point opposite ways [g-1][g-2]. With a dead heat on the direct evidence, the tiebreaker is base rates: g-2 reports a *majority* (58%) rolling back within 18 months, while g-1's saving is conditional on a stabilization most teams in g-2 never reach. The one high-confidence item [g-3] is about boundary mapping, not cost ...",
  "answers": [
    { "cq_id": "cq1", "answer": "Raises TCO for the median team; the cost reduction is real but rare (conditional on stabilization that the majority in g-2 do not achieve).", "evidence_refs": ["g-1","g-2"] },
    { "cq_id": "cq2", "answer": "Strategic design has the only high-confidence support [g-3], but that is a boundary-quality result, not a TCO result — so 'value is strategic' is under-evidenced for the cost claim specifically.", "evidence_refs": ["g-3"] },
    { "cq_id": "cq3", "answer": "The pack does not isolate the cause of abandonment; g-2 gives the rate, not the trigger.", "evidence_refs": ["g-2"] }
  ],
  "key_claims": [
    { "claim": "Direct cost evidence is a medium-vs-medium tie pointing opposite ways.", "grounding": ["g-1","g-2"], "confidence": "medium" },
    { "claim": "Base-rate reasoning favors the cost-raising side.", "grounding": ["g-2"], "confidence": "medium" }
  ],
  "strongest_counterargument_acknowledged": "If abandonment is driven by tactical ceremony specifically (not strategic DDD), then teams that adopt only strategic design could capture g-1's savings without paying g-2's cost — which would flip the proposition back to true.",
  "uncited_assertions": ["When direct evidence ties, base rates should break the tie (a reasoning principle, not a pack item)."],
  "status": "ok"
}
```

## Instance 3 — `panelist = P-claude`

```json
{
  "panelist_id": "P-claude", "cli": "claude -p", "model": "opus",
  "round": 0,
  "stance": "The proposition is true but only as a conditional: DDD lowers TCO for teams that adopt strategic design and gate tactical ceremony on demonstrated need.",
  "position": "Both cost findings can be true simultaneously if they describe different populations. g-1 describes teams whose model stabilized; g-2 describes teams that rolled back tactical patterns — these need not be the same teams. The high-confidence boundary result [g-3] suggests the strategic core is robust independent of the tactical-ceremony churn ...",
  "answers": [
    { "cq_id": "cq1", "answer": "Lowers TCO conditional on strategic-first adoption; the disagreement in g-1/g-2 dissolves if they describe different populations.", "evidence_refs": ["g-1","g-2","g-3"] },
    { "cq_id": "cq2", "answer": "Yes — g-3 is the only high-confidence anchor and it is strategic; tactical patterns are where g-2's cost concentrates.", "evidence_refs": ["g-2","g-3"] },
    { "cq_id": "cq3", "answer": "Ceremony stops paying when tactical patterns are mandated before the domain model is stable; g-2's timeframe is consistent with premature tactical adoption.", "evidence_refs": ["g-2"] }
  ],
  "key_claims": [
    { "claim": "g-1 and g-2 are reconcilable as different populations, not a true contradiction.", "grounding": ["g-1","g-2"], "confidence": "medium" },
    { "claim": "The strategic core (g-3) is robust to tactical churn.", "grounding": ["g-3"], "confidence": "high" }
  ],
  "strongest_counterargument_acknowledged": "The different-populations reconciliation is unproven from the pack — g-1 and g-2 could be the *same* teams, in which case the proposition fails for the median team (P-gemini's base-rate reading).",
  "uncited_assertions": ["The population-overlap question is the crux and the pack does not resolve it (meta-observation, not a pack claim)."],
  "status": "ok"
}
```

Note: three genuinely different opening stances (conditional-yes /
evidence-says-no / reconcilable-yes), each grounded in the same `[g-n]`
world, each naming a real counter to itself, each identity-stamped. This is
the healthy starting state for a debate — divergence the later rounds will
test.
