# Worked example: frame-debate

A `standard` / `general` brief built from a real squad-researcher run, plus a
degraded-mode mini-example.

---

## Example A — full grounding (primary path)

### Input

```
research_run_dir: /Users/.../squad-researcher/tmp/runs/ddd-in-enterprise-systems-20260514-161443
topic:            (unused — derived from the report)
rounds:           standard
audience:         general
panel:            [ {P-codex,codex exec,...}, {P-gemini,gemini -p,...}, {P-claude,claude -p,...} ]
```

Relevant slices the skill read:

`05-final_report.md` (H1): `# Domain-Driven Design in enterprise systems`

`03-findings.json` (abbreviated):
```json
[
  { "id": "f-12", "claim": "Teams adopting DDD report lower long-run change cost once the model stabilizes.",
    "evidence": "Across the surveyed case studies, change-failure rate dropped ...", "confidence": "medium",
    "disputed_by": ["f-31"] },
  { "id": "f-31", "claim": "A majority of surveyed teams abandoned strict DDD within 18 months citing ceremony.",
    "evidence": "58% of respondents reported rolling back tactical patterns ...", "confidence": "medium",
    "disputed_by": ["f-12"] },
  { "id": "f-44", "claim": "DDD's bounded contexts map cleanly onto service boundaries in microservice migrations.",
    "evidence": "Migration retrospectives consistently cite context maps ...", "confidence": "high",
    "disputed_by": [] }
]
```

`01-research_plan.json` (abbreviated):
```json
{ "thesis_question": "When does Domain-Driven Design pay off in enterprise systems?",
  "success_criteria": [
    "Names at least two concrete conditions under which DDD pays off",
    "Addresses the abandonment/ceremony counter-evidence",
    "Distinguishes strategic DDD from tactical patterns" ],
  "out_of_scope": ["event sourcing implementation detail", "specific ORM choice"] }
```

### Output — `debate_brief`

```json
{
  "topic": "Domain-Driven Design in enterprise systems",
  "rounds": "standard",
  "exchange_count": 1,
  "audience": "general",
  "grounding": "full",
  "grounding_note": "Grounded in squad-researcher run ddd-in-enterprise-systems-20260514-161443 (3 findings carried; 1 disputed pair).",
  "source_run": "ddd-in-enterprise-systems-20260514-161443",
  "proposition": "Domain-Driven Design lowers the total cost of ownership of enterprise systems more often than it raises it.",
  "supporting_propositions": [],
  "contested_questions": [
    { "id": "cq1",
      "question": "Does adopting DDD lower or raise total cost of ownership in enterprise systems, and what determines which outcome a team gets?",
      "why_contested": "grounding items g-1 and g-2 directly disagree (f-12 vs f-31)",
      "seed": "finding_dispute",
      "evidence_refs": ["g-1","g-2"] },
    { "id": "cq2",
      "question": "Is the value of DDD mostly in strategic design (bounded contexts, context maps) rather than tactical patterns?",
      "why_contested": "plan success_criterion 'distinguish strategic from tactical' is only partially met by g-3",
      "seed": "plan_open_question",
      "evidence_refs": ["g-3"] },
    { "id": "cq3",
      "question": "Under what team/domain conditions does DDD's modeling ceremony stop paying for itself?",
      "why_contested": "report raises the 18-month abandonment signal but does not isolate the cause",
      "seed": "synthesizer_gap",
      "evidence_refs": ["g-2"] }
  ],
  "rubric": [
    { "id": "r1", "criterion": "Position is grounded in [g-n] items, cited adjacent to the claim", "weight": 1.0 },
    { "id": "r2", "criterion": "Names at least two concrete conditions under which DDD pays off", "weight": 1.0 },
    { "id": "r3", "criterion": "Engages the abandonment/ceremony counter-evidence rather than ignoring it", "weight": 1.0 },
    { "id": "r4", "criterion": "Distinguishes strategic DDD from tactical patterns", "weight": 1.0 },
    { "id": "r5", "criterion": "Internally consistent; the strongest opposing case is acknowledged", "weight": 1.0 }
  ],
  "panel": [
    { "panelist_id": "P-codex",  "cli": "codex exec", "tool": "codex-cli",   "tool_version": "0.130.0", "model": "gpt-5.5" },
    { "panelist_id": "P-gemini", "cli": "gemini -p",  "tool": "gemini-cli",  "tool_version": "0.42.0",  "model": "gemini-2.5-pro" },
    { "panelist_id": "P-claude", "cli": "claude -p",  "tool": "claude-code", "tool_version": "2.1.143", "model": "opus" }
  ],
  "steelman_directive": null,
  "brief_skipped": false
}
```

### Output — `grounding_pack`

```json
{
  "summary": "The report finds DDD's payoff is conditional: strategic design (bounded contexts, context maps) shows consistent benefit, while tactical pattern ceremony correlates with abandonment. Net effect on cost of ownership is contested across the case base.",
  "items": [
    { "id": "g-1", "claim": "Teams adopting DDD report lower long-run change cost once the model stabilizes.",
      "evidence": "Across the surveyed case studies, change-failure rate dropped ...", "confidence": "medium", "disputed": true, "origin": "finding:f-12" },
    { "id": "g-2", "claim": "A majority of surveyed teams abandoned strict DDD within 18 months citing ceremony.",
      "evidence": "58% of respondents reported rolling back tactical patterns ...", "confidence": "medium", "disputed": true, "origin": "finding:f-31" },
    { "id": "g-3", "claim": "DDD's bounded contexts map cleanly onto service boundaries in microservice migrations.",
      "evidence": "Migration retrospectives consistently cite context maps ...", "confidence": "high", "disputed": false, "origin": "finding:f-44" }
  ],
  "out_of_scope": ["event sourcing implementation detail", "specific ORM choice"]
}
```

Note: `cq1` is the load-bearing question — it came straight from the
`disputed_by` pair, neither side baked into the phrasing. `cq3` is a
`synthesizer_gap`: the report *raised* the abandonment signal but did not
argue out the cause, so it is fair debate fuel (not invented).

---

## Example B — degraded mode (mini)

### Input

```
research_run_dir: (empty / not a valid run dir)
topic:            "Should small teams default to a modular monolith over microservices?"
source_text:      (none)
rounds:           quick
audience:         general
```

### Output — `debate_brief` (abbreviated)

```json
{
  "topic": "Should small teams default to a modular monolith over microservices?",
  "rounds": "quick", "exchange_count": 0, "audience": "general",
  "grounding": "degraded",
  "grounding_note": "⚠️ degraded grounding — debated from topic with no squad-researcher findings; conclusions are unverified and confidence is capped at medium.",
  "source_run": null,
  "proposition": "Small teams should default to a modular monolith and adopt microservices only on a demonstrated need.",
  "supporting_propositions": [],
  "contested_questions": [
    { "id": "cq1", "question": "What concrete signal should flip a small team from monolith to services?",
      "why_contested": "no grounding; genuine open design question", "seed": "degraded_inferred", "evidence_refs": [] }
  ],
  "rubric": [
    { "id": "r1", "criterion": "Position is internally consistent and names a concrete decision rule", "weight": 1.0 },
    { "id": "r2", "criterion": "Acknowledges the strongest opposing case", "weight": 1.0 },
    { "id": "r3", "criterion": "States its assumptions explicitly (degraded: no evidence base)", "weight": 1.0 }
  ],
  "panel": [ /* echoed */ ],
  "steelman_directive": null,
  "brief_skipped": false
}
```

`grounding_pack`: `summary` = the single ungrounded-debate sentence,
`items: []`, `out_of_scope: []`.
