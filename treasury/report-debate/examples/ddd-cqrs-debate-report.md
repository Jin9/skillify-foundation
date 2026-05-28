# Worked example: report-debate — a rendered 00-debate_report.md

A complete `00-debate_report.md` for the DDD `standard` / `general` run
(metrics.json fully captured, 3/3 panelists, 1 exchange round).

---

```markdown
# Debate report — Domain-Driven Design lowers the total cost of ownership of enterprise systems more often than it raises it.

_rounds: standard · audience: general · grounding: full · debate_id: ddd-in-enterprise-systems-20260516-101859 · started_at: 2026-05-16T10:18:59Z · panel: codex 0.130.0 / gemini 0.42.0 / claude 2.1.143 · panel_order_seed: 4193_

## Panel

| Panelist | CLI | Tool version | Model | Turns | Final stance_delta |
|---|---|---|---|---|---|
| P-codex | codex exec | 0.130.0 | gpt-5.5 | 2 | unchanged |
| P-gemini | gemini -p | 0.42.0 | gemini-2.5-pro | 2 | narrowed |
| P-claude | claude -p | 2.1.143 | opus | 2 | unchanged |

## Per-turn stats

| Turn | CLI | Model | Status | Wall | Bytes |
|---|---|---|---|---|---|
| panel-open::P-codex::r0 | codex exec | gpt-5.5 | ok | 41s | 5120 |
| panel-open::P-gemini::r0 | gemini -p | gemini-2.5-pro | ok | 33s | 4730 |
| panel-open::P-claude::r0 | claude -p | opus | ok | 52s | 6010 |
| cross-examine::P-codex::r1 | codex exec | gpt-5.5 | ok | 38s | 3980 |
| cross-examine::P-gemini::r1 | gemini -p | gemini-2.5-pro | ok | 35s | 4120 |
| cross-examine::P-claude::r1 | claude -p | opus | ok | 47s | 4890 |
| revise-positions::P-codex::r1 | codex exec | gpt-5.5 | ok | 36s | 4310 |
| revise-positions::P-gemini::r1 | gemini -p | gemini-2.5-pro | ok | 40s | 5020 |
| revise-positions::P-claude::r1 | claude -p | opus | ok | 44s | 4660 |
| synthesize-consensus | claude -p | opus | ok | 71s | 7240 |
| **Total** | | | **10/10 ok** | **477s** | **50090** |

## Convergence map

- **P-codex**: opening "true-conditional (strategic-first)" → r1 unchanged → final "true-conditional". `unchanged`
- **P-gemini**: opening "evidence-says-no (base rates)" → r1 narrowed ⟲ → final "unproven pending stabilized-population base rate". `narrowed →`
- **P-claude**: opening "true-conditional (different populations)" → r1 unchanged → final "true-conditional". `unchanged`

_Partial convergence: all three aligned on the strategic-core anchor (principled — P-gemini moved via a g-1-defended concession, not capitulation); the cost-direction question stayed an honest live disagreement through round 1._

## Agreements

- The strategic core (bounded contexts, context maps) is the only high-confidence cost-relevant anchor — held by P-codex, P-gemini, P-claude — [g-3] — confidence high — via principled_convergence.
- The two direct cost findings are opposite-signed and both medium-confidence; neither settles the proposition alone — held by all three — [g-1][g-2] — confidence medium — via independent.

## Live disagreements

**cq1 — Is the proposition "true-conditional" or "unproven"?**
- *side 1 — P-codex, P-claude:* True as a conditional. Strongest case: g-3 (high confidence) shows the strategic core is robust independent of tactical churn; nothing forces g-1 and g-2 onto the same teams. [g-1][g-3]
- *side 2 — P-gemini:* Unproven, not false. Strongest case: g-1 gives no base rate for reaching stabilization; g-2 shows 58% rollback; the headline claim carries an unmet quantification burden. [g-1][g-2]
- **Moderator lean: genuinely_unresolved.** _Rubric r1 (grounded) and r3 (engages counter-evidence) are satisfied by BOTH steelmanned sides. The deciding fact — the size of the stabilized population — is exactly what no [g-n] item supplies; per r1 a lean would require evidence the pack does not contain; therefore unresolved._

## Transcript index

```
00-panel.json
01-debate_brief.json
01-grounding_pack.json
02-open__P-codex.json        (round 0, P-codex)
02-open__P-gemini.json       (round 0, P-gemini)
02-open__P-claude.json       (round 0, P-claude)
03-xexam__r1__P-codex.json   (round 1, P-codex critiques P-gemini,P-claude)
03-xexam__r1__P-gemini.json  (round 1, P-gemini critiques P-codex,P-claude)
03-xexam__r1__P-claude.json  (round 1, P-claude critiques P-codex,P-gemini)
04-revise__r1__P-codex.json  (round 1, P-codex)
04-revise__r1__P-gemini.json (round 1, P-gemini)
04-revise__r1__P-claude.json (round 1, P-claude)
05-final_answer.md
05-convergence.json
00-debate_report.md          (this file)
```

## Pipeline integrity

- Grounding-citation integrity: 3/3 `[g-n]` markers in 05-final_answer.md resolve to grounding_pack items — OK.
- Moderator-favoritism: leaned to P-claude's side on 0/1 contested questions (the only live disagreement was left genuinely_unresolved, against the P-codex+P-claude majority). no_favoritism_check (verbatim): "cq1 was left genuinely_unresolved strictly on rubric r1 (no [g-n] sizes the stabilized population) — a verdict that goes AGAINST the side P-claude argued, demonstrating the lean tracks evidence, not identity. The single agreement upgrade to 'high' rests only on g-3's confidence tier, which is identity-independent."
- Absent panelists: 0. All rounds ran with 3/3.
- Rounds executed: 1 exchange (cross-examine + revise) — matches `standard`.
- Convergence honesty: the strategic-core agreement is principled (P-gemini's changelog traces to a g-1-defended concession against a major objection); no capitulation flagged.

## Distillation

The panel could not settle this on the available evidence, and that is the
honest answer: DDD's strategic core clearly pays off [g-3], but whether the
headline cost claim holds depends on how often teams actually reach the
"model stabilized" state — a number the research never establishes. Treat
"DDD lowers TCO" as conditional and unproven at the population level, not as
a default.

## Notes
Report generated by report-debate@0.1.0 on 2026-05-16T10:31:12Z.
```

---

Note how the report is a faithful mirror: the moderator's
`genuinely_unresolved` verdict and `no_favoritism_check` are pasted verbatim;
the favoritism line is a *computed statistic* ("0/1") presented for human
audit, not a re-judgement; the metrics table sums only the captured turns;
the distillation is the `## Bottom line` copied exactly.
