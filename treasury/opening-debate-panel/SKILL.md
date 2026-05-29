---
name: opening-debate-panel
description: >
  Stage 2 of the squad-brainstorm workflow, run once per panel CLI (P-codex
  via `codex exec`, P-gemini via `gemini -p`, P-claude via `claude -p`).
  Given the debate brief and the shared grounding pack, the invoked panelist
  produces ONE independent, evidence-grounded opening position answering the
  proposition and every contested question, citing only [g-n] items from the
  grounding pack and naming the strongest counterargument to its own view.
  The panelist does NOT see any other panelist's output at this stage
  (independence prevents early herding). Use when the workflow runs the panel-open stage of
  workflows/brainstorm.yaml. Do NOT critique peers (cross-examine), revise
  (revise-positions), or merge (synthesize-consensus). One panelist, one
  call, one position object.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # the prompt itself calls no tools; the driver invokes the external CLI
version: 0.1.0
stage: panel-open
workflow: brainstorm
inputs:
  - { name: debate_brief,   type: object, required: true, source: stages.frame-debate.debate_brief,   description: "from frame-debate" }
  - { name: grounding_pack, type: object, required: true, source: stages.frame-debate.grounding_pack, description: "the closed evidence world; cite [g-n] only" }
  - { name: panelist,       type: object, required: true, source: fanout.panelist,                    description: "this instance's panel[] row: panelist_id, cli, model, tool_version" }
outputs:
  - { name: positions, type: array, description: "collected across fan-out: one position object per panelist, round 0" }
---

# Skill: Opening Debate Panel

## Purpose

Produce one panelist's **independent opening stance** on the debate. This is
the squad-brainstorm analogue of a debater's opening statement: a clear
position, argued from the shared evidence, with the strongest objection to it
named up front.

Independence is load-bearing. At this stage no panelist may see another's
view — this prevents the early herding / anchoring that collapses multi-model
debates into one model's framing. It is the debate analogue of
`extract-findings`' evidence-first discipline.

**Atomicity:** one panelist, one CLI call, one `position` object. The driver
runs this skill three times in parallel (once per `panel[]` row) and
collects the three objects into the `positions` array.

## When to use this skill

- Stage 2 fan-out dispatch from `workflows/brainstorm.yaml`.
- Prompt: "give your opening position on this debate brief".

Do NOT use this skill to:
- Critique another panelist — that is `cross-examine` (stage 3).
- Revise your position — that is `revise-positions` (stage 4).
- Merge the debate — that is `synthesize-consensus` (stage 5).

## Inputs

| Name | Type | Notes |
|------|------|-------|
| `debate_brief` | object | From `frame-debate`. The `proposition`, `contested_questions[]`, `rubric`, `rounds`. |
| `grounding_pack` | object | The **only** evidence you may cite. Reference items as `[g-n]`, exactly like researcher's `[n]`. No web, no outside knowledge presented as grounded. |
| `panelist` | object | Your own identity row: `{ panelist_id, cli, model, tool_version }`. Stamp it verbatim onto your output. You MUST NOT claim to be a different panelist or model. |

If `debate_brief.brief_skipped == true`, emit a position with
`status: "absent"` and `position: ""` — do not invent a debate.

## Output contract

Each fan-out instance emits exactly one object (the driver collects three
into `positions`):

```jsonc
{
  "panelist_id": "P-codex",            // verbatim from input panelist.panelist_id — identity stamp
  "cli": "codex exec",                 // echo of panelist.cli (provenance)
  "model": "gpt-5.5",              // echo of panelist.model
  "round": 0,                          // opening = round 0
  "stance": "string",                  // ONE sentence: your answer to debate_brief.proposition
  "position": "markdown",              // the argued case; [g-n] cites adjacent to claims
  "answers": [
    { "cq_id": "cq1", "answer": "string", "evidence_refs": ["g-3","g-7"] }
  ],
  "key_claims": [
    { "claim": "string", "grounding": ["g-2"], "confidence": "high|medium|low" }
  ],
  "strongest_counterargument_acknowledged": "string",  // the best case AGAINST your stance
  "uncited_assertions": ["string"],    // claims you make that are NOT backed by a [g-n] item (honesty ledger)
  "status": "ok"                       // "absent" only on brief_skipped / failure
}
```

`position` is markdown prose. Put `[g-n]` markers adjacent to the claim they
support, never clustered at paragraph end (same precision rule as
`synthesize-report`). There is exactly one `answers` entry per
`debate_brief.contested_questions` entry.

## Procedure

Execute in a single CLI call:

1. **Read the brief and the grounding pack.** Note the `proposition`, every
   `contested_questions` entry, the `rubric` (you will be judged on it), and
   `grounding_pack.out_of_scope` (do not argue excluded scope).
2. **Decide your stance** — one defensible sentence answering the
   proposition. You may support, oppose, or qualify it; pick the position
   the evidence best supports, not a contrarian pose.
3. **Answer every contested question** with a grounded `answers` entry.
   Cite `[g-n]` ids that exist in `grounding_pack.items`. If the pack does
   not speak to a question, say so explicitly and reason carefully — and
   add the ungrounded part to `uncited_assertions`.
4. **Write the `position`** prose, rubric-aware, `[g-n]` adjacent to claims.
   See [`references/grounding-rules.md`](references/grounding-rules.md).
5. **Name the strongest counterargument** to your own stance honestly —
   the best case the other side could make. A weak self-assigned counter is
   a scored failure (anti-strawman; the dominant debate failure mode).
6. **List `uncited_assertions`** — any factual claim you make that is not
   backed by a `[g-n]` item. Do not smuggle ungrounded claims in as
   grounded; an honest empty-handed claim beats a fake citation.
7. **Stamp identity** — copy `panelist_id`, `cli`, `model` from the input
   `panelist` byte-for-byte. `round: 0`.
8. **Self-check** the Validation gate; fix in place.
9. **Emit** the single JSON object. No prose around it.

## Rounds dial

`rounds` (read from `debate_brief.rounds`) sets only opening length and
breadth:

| `rounds` | `position` length | Breadth |
|----------|-------------------|---------|
| `quick` | 400–600 words | Answer the proposition + each contested question tersely |
| `standard` | 600–1000 words | Argue each contested question with its evidence |
| `deep` | 1000–1800 words | Engage every contested question; if `debate_brief.steelman_directive` is set, include one paragraph arguing the strongest opposing case before refuting it |

`audience` is **not** consulted here — opening positions are technical and
internal; audience framing is applied once, at synthesis.

## Constraints

- **Independence.** No peer input exists at this stage. Do not speculate
  about, reference, or pre-empt other panelists.
- **Closed evidence world.** Cite `[g-n]` only. No web, no training-data
  facts presented as grounded. Ungrounded reasoning is allowed but must be
  flagged in `uncited_assertions` (same closed-world rule as
  `extract-findings` / `synthesize-report`).
- **Identity honesty.** Stamp the given `panelist_id`/`cli`/`model`
  verbatim. Never claim to be a different panelist or model.
- **One call.** No recursion, no clarifying questions.
- **JSON only** in workflow context. The driver parses the object.
- Do not argue anything in `grounding_pack.out_of_scope`.

## Validation gate

- [ ] `panelist_id`, `cli`, `model` equal the input `panelist` byte-for-byte.
- [ ] `round == 0`.
- [ ] `stance` is one sentence answering `debate_brief.proposition`.
- [ ] Exactly one `answers` entry per `debate_brief.contested_questions`.
- [ ] Every `evidence_refs` / `grounding` id exists in
      `grounding_pack.items`.
- [ ] `strongest_counterargument_acknowledged` is non-empty and is a real
      objection to *your* stance (not a strawman).
- [ ] Every factual claim not backed by a `[g-n]` item appears in
      `uncited_assertions`.
- [ ] `position` length within the Rounds dial band (±20%).

## Worked example

See [`examples/three-panelists-open.md`](examples/three-panelists-open.md):
the same brief answered independently by all three CLIs, showing three
differently-shaped but identity-stamped, grounding-true position objects.

## References

- [`references/grounding-rules.md`](references/grounding-rules.md) — the
  `[g-n]`-only closed-world contract and the honesty-ledger rule (shared
  verbatim with `cross-examine` and `revise-positions`).
