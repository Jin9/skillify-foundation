# Worked example — augment-diagrams (stack/layer)

Run: `capability-cost-model-router-20260515-023858`
Depth: `deep` · audience: `expert` · diagrams emitted: `1`

This example shows ONE body-section augmentation: the original paragraph
from `05-final_report.md`, the stack/layer diagram inserted immediately
after it, and the matching entry in `06-augment_notes.md`. The rest of
the report is unchanged.

## Input — paragraph in `05-final_report.md`

```
## Routing taxonomy

The contemporary taxonomy splits along two cuts: when the router decides
(pre-generation vs. post-generation cascade) and how it decides
(rule-based, embedding-similarity, learned classifier, LLM-judge) [1].
Pre-generation routers commit before any inference cost is paid;
post-generation cascades evaluate a cheap model's output and only escalate
when a scoring function rejects it [4]. Production systems sit on both
sides — Bedrock IPR routes pre-generation within a family [7], while
FrugalGPT's cascade is the canonical post-generation pattern [4].
```

`cited_sources[*]` for this section's `[n]` markers:
`[1] = taxonomy survey; [4] = FrugalGPT; [7] = AWS Bedrock IPR`. Their
`findings_supported[]` resolve to findings about the
pre-vs-post-generation cut, cascade scoring, and family-bounded routing.

## Decision

The section's load-bearing structure is a 2-axis decision space: **when
to decide** × **how to decide**. The prose describes it linearly, but a
stack/layer diagram makes the two cuts simultaneously scannable. The
scan test passes: a reader who sees only the diagram still gets "two
cuts, one pre-gen one post-gen, decision-style is orthogonal".

## Output — augmented paragraph in `augmented_report`

```
## Routing taxonomy

The contemporary taxonomy splits along two cuts: when the router decides
(pre-generation vs. post-generation cascade) and how it decides
(rule-based, embedding-similarity, learned classifier, LLM-judge) [1].
Pre-generation routers commit before any inference cost is paid;
post-generation cascades evaluate a cheap model's output and only escalate
when a scoring function rejects it [4]. Production systems sit on both
sides — Bedrock IPR routes pre-generation within a family [7], while
FrugalGPT's cascade is the canonical post-generation pattern [4].
```
```
┌──────────────────────────────────────────────────────────────┐
│ when ↓ \\ how →   rule   embed   classifier   LLM-judge       │
├──────────────────────────────────────────────────────────────┤
│ pre-generation    IPR    RouteLLM  Hybrid-LLM   NotDiamond    │
├──────────────────────────────────────────────────────────────┤
│ post-generation   ─       ─        FrugalGPT    LLM-Blender   │
└──────────────────────────────────────────────────────────────┘
```
```

Note: every entity (`IPR`, `RouteLLM`, `Hybrid-LLM`, `NotDiamond`,
`FrugalGPT`, `LLM-Blender`) appears as a cited finding in the section's
existing prose or in adjacent prose already cited via `[n]`. No `[n]`
inside the fenced block. ≤ 15 lines, ≤ 70 cols. No language tag.

## Output — `06-augment_notes.md` entry

```
### Section: Routing taxonomy
- decision: emit
- diagram_type: matrix (when × how)
- grounding_finding_ids: [f1, f4, f7, f3, f5, f6]
- scan_test: passed — two cuts visible simultaneously without prose
- placement: after first prose paragraph (offset recorded in Step 2)
```

## Why the scan test passes

A reader who scans only the matrix gets the same load-bearing claim the
paragraph carries: two orthogonal axes, one canonical occupant per cell
(plus an empty row for post-generation rule/embedding, which is itself a
signal — that part of the design space is sparse). The prose adds
mechanism; the diagram makes the **shape** of the design space scannable
in one second. That is the bar.
