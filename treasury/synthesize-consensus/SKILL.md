---
name: synthesize-consensus
description: >
  Stage 5 of the squad-brainstorm workflow. A single Claude moderator reads
  the attributed final panelist positions plus the cross-examination
  exchanges and produces the user-facing synthesized answer: settled
  agreements, live disagreements with the strongest case for each side, a
  calibrated confidence, and one integrated answer to the proposition framed
  for the audience. This is an OPEN moderator (it sees panelist_id) — so it
  MUST judge only on the rubric + grounding evidence, never by which model
  said what, and must emit an explicit no-favoritism self-check. Use when the workflow runs the
  synthesize-consensus stage of workflows/brainstorm.yaml. Single-pass; it
  merges, it does not re-debate. Do NOT use to argue or critique.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # closed-world merge over in-context positions
version: 0.1.0
stage: synthesize-consensus
workflow: brainstorm
inputs:
  - { name: debate_brief,     type: object, required: true,  source: stages.frame-debate.debate_brief,                 description: "proposition, contested questions, rubric, grounding flag" }
  - { name: grounding_pack,   type: object, required: true,  source: stages.frame-debate.grounding_pack,               description: "the same [g-n] closed world the panel used" }
  - { name: audience,         type: string, required: false, source: workflow.inputs.audience,                         description: "frames the final answer's tone/vocabulary" }
  - { name: final_positions,  type: array,  required: true,  source: "stages.revise-positions.revised_positions ?? stages.panel-open.positions", description: "attributed last-round positions (panelist_id visible)" }
  - { name: debate_exchanges, type: array,  required: true,  source: stages.cross-examine.critiques,                   description: "attributed critiques across all rounds ([] when rounds=quick)" }
outputs:
  - { name: answer,      type: string, description: "the single integrated final answer, audience-framed, [g-n] cited" }
  - { name: convergence, type: object, description: "agreements, live disagreements (strongest case each), confidence, no-favoritism self-check" }
---

# Skill: synthesize-consensus

## Purpose

Produce the one answer the user gets out of the whole debate. This is the
highest-leverage reasoning stage (top model tier — the squad-brainstorm
analogue of researcher's "reviewer ≥ writer": the merger must out-reason
three frontier panelists). It clusters the final positions into what the
panel actually settled, what it genuinely still disputes (with the strongest
case for each side, not a laundered false consensus), and a single
integrated answer to `debate_brief.proposition`.

**Open moderator.** Per the project decision there is no anonymization
layer: this stage sees `panelist_id` on every position and critique,
including that one panelist is `P-claude` and that this stage itself runs on
a Claude model. The bias control is therefore explicit and behavioral, not
structural: judge **only** on `debate_brief.rubric` + `[g-n]` evidence,
never by which model authored a position, and prove it with a recorded
`no_favoritism_check`.

**Atomicity:** one Claude call. It merges; it does not re-open the debate,
re-critique, or add new claims.

## When to use this skill

- Stage 5 dispatch from `workflows/brainstorm.yaml` — after the final
  revision round, or directly after `panel-open` when `rounds == quick`.
- Prompt: "synthesize the debate into one answer", "merge these positions".

Do NOT use this skill to:
- Argue or critique — stages 2–4.
- Re-run or extend the debate.
- Summarize run mechanics — that is `report-debate` (stage 6).

## Inputs

| Name | Type | Notes |
|------|------|-------|
| `debate_brief` | object | `proposition`, `contested_questions[]`, `rubric`, `grounding`, `grounding_note`. |
| `grounding_pack` | object | The same `[g-n]` world the panel used. The answer cites `[g-n]` only. |
| `audience` | string | Frames the final answer's tone/vocabulary (researcher audience-profile mapping). Does NOT change which positions are merged. |
| `final_positions` | array | The last-round position per panelist, attributed (`panelist_id` present). On `quick`, the `panel-open` openings. |
| `debate_exchanges` | array | All `cross-examine` critique objects, attributed. `[]` on `quick`. |

If `debate_brief.brief_skipped == true`, emit an empty answer and a
`convergence` with `overall_confidence: "low"` noting the skipped brief.

## Output contract

### `convergence`

```jsonc
{
  "proposition": "string",                 // echo of debate_brief.proposition
  "panelists": ["P-codex","P-gemini","P-claude"],  // attributed (open moderator)
  "agreements": [
    { "point": "string", "held_by": ["P-codex","P-claude"], "grounding": ["g-2"],
      "confidence": "high|medium|low", "via": "independent|principled_convergence" }
  ],
  "live_disagreements": [
    {
      "cq_id": "cq1",
      "summary": "string",
      "sides": [
        { "held_by": ["P-gemini"], "position": "string", "strongest_case": "string", "grounding": ["g-1"] },
        { "held_by": ["P-codex","P-claude"], "position": "string", "strongest_case": "string", "grounding": ["g-3"] }
      ],
      "moderator_lean": "side-1|side-2|genuinely_unresolved",
      "lean_rationale": "string"           // grounded reason tied to the rubric, or why it stays unresolved
    }
  ],
  "overall_confidence": "high|medium|low",
  "confidence_rationale": "string",
  "convergence_summary": "string",         // did the panel converge / partially / hold divergent over rounds?
  "grounding": "full|degraded",            // carried from debate_brief; gates hedging
  "no_favoritism_check": "string"          // MANDATORY: one line affirming the verdict was decided on rubric+[g-n], not panelist identity; cite the rubric criterion/[g-n] that drove each lean
}
```

### `answer` — markdown

```
# <proposition restated as a question>

## Bottom line
<the integrated verdict, first — busy readers stop here. 2–5 sentences.>

## Where the panel agrees
<bulleted; each with [g-n] and whether it was independent or via principled convergence>

## Live disagreements
<for each: the question, the strongest case for EACH side (never collapsed),
and the moderator's grounded lean or an explicit "genuinely unresolved">

## Confidence & caveats
<overall confidence + rationale; include debate_brief.grounding_note VERBATIM
when grounding == degraded>

## Grounding
<the [g-n] items actually cited, rendered like researcher's ## Sources:
`- [g-1] <claim> — origin <origin>, confidence <confidence>`>
```

Cite `[g-n]` adjacent to claims. `## Bottom line` comes before
`## Live disagreements` (busy-reader discipline, same as researcher's
Executive Summary).

## Procedure

Single Claude call:

1. **Read** the brief, grounding pack, and every attributed final position
   and exchange.
2. **Cluster** claims across panelists by `contested_questions`: which are
   now agreed (held by ≥2 panelists with compatible grounding) vs still
   disputed. Mark each agreement `via: "independent"` (panelists agreed from
   the start) or `"principled_convergence"` (a position moved via a
   traceable changelog — read `final_positions[].changelog` /
   `stance_delta`). Discount agreement that arrived via capitulation
   (changelog shows `conceded` against only `minor` objections) — note it,
   do not count it as strong consensus. See
   [`references/false-consensus.md`](references/false-consensus.md).
3. **For each live disagreement**, write the **strongest** case for *each*
   side (steelman both — never launder a real split into fake agreement).
   Take a `moderator_lean` only where the `rubric` + `[g-n]` support one
   side; else `genuinely_unresolved`. The `lean_rationale` must cite the
   rubric criterion and/or `[g-n]` item that decides it — not "P-claude was
   more convincing".
4. **Run the no-favoritism check.** Verify each `moderator_lean` is
   justified by rubric/evidence independent of who held it. Write
   `no_favoritism_check` stating this explicitly and naming, per lean, the
   rubric criterion or `[g-n]` that drove it. If any lean coincides with
   `P-claude`'s side, the rationale must be extra-explicit about the
   evidential basis (defense-in-depth against the known open-moderator
   risk). See [`references/no-favoritism.md`](references/no-favoritism.md).
5. **Integrate** into one `answer`, framed for `audience` (general /
   executive / expert / academic — closest-profile mapping, same as
   researcher's `synthesize-report`). Audience changes framing only.
6. **Calibrate confidence** (qualitative `high|medium|low`, never numeric).
   `grounding == degraded` caps `overall_confidence` at `medium` and the
   `grounding_note` is pasted verbatim into `## Confidence & caveats`.
7. **Self-check** the Validation gate; fix in place.
8. **Emit** `answer` (markdown) and `convergence` (object).

## Rounds dial

`quick`: merges 3 openings; `convergence_summary` is a snapshot (no trend),
`debate_exchanges` is `[]`. `standard`: one exchange's worth of movement to
read. `deep`: two rounds; add a short dissent paragraph in `## Live
disagreements` for any side the lean went against, and address every
`contested_questions` entry. `audience` frames the answer at every `rounds`.

## Constraints

- **Open moderator, but identity-blind judgment.** You see `panelist_id`;
  you must not let it influence the verdict. No "the Claude position was
  strongest" reasoning — only rubric + `[g-n]`.
- **No false consensus.** A live disagreement must surface as one, with both
  strongest cases. Never merge a real split into a single bland claim.
- **Single-pass.** Merge, do not re-debate, re-critique, or add new claims
  beyond the positions + grounding pack.
- **Closed evidence world.** `[g-n]` only; no new evidence.
- **Honest confidence.** Degraded grounding ⟹ confidence ≤ `medium` + the
  verbatim banner.
- **JSON + markdown only** in workflow context.

## Validation gate

- [ ] `answer` H1 restates `debate_brief.proposition` as a question;
      `## Bottom line` precedes `## Live disagreements`.
- [ ] Every `[g-n]` in `answer` exists in `grounding_pack.items`; no
      invented evidence.
- [ ] Every `debate_brief.contested_questions` entry is resolved as an
      `agreement` or appears in `live_disagreements`.
- [ ] Every `live_disagreements[]` has a non-empty `strongest_case` for
      **every** side (no collapsed split).
- [ ] `overall_confidence ∈ {high,medium,low}`; `≤ medium` and the verbatim
      `grounding_note` is present when `grounding == degraded`.
- [ ] `no_favoritism_check` is non-empty and names, per `moderator_lean`,
      the rubric criterion / `[g-n]` that drove it (not a panelist).
- [ ] `agreements[].via` set; capitulation-driven agreement is noted, not
      counted as strong consensus.
- [ ] No new claim appears that is not traceable to a position or a
      `[g-n]` item.

## Worked example

See [`examples/ddd-cqrs-consensus.md`](examples/ddd-cqrs-consensus.md): the
attributed DDD final round → an integrated answer that names a live
disagreement with both strongest cases, a rubric-grounded lean, and the
explicit no-favoritism check.

## References

- [`references/false-consensus.md`](references/false-consensus.md) — the
  laundering anti-pattern and the strongest-case-each-side rule; how to
  read the changelog for capitulation vs principled convergence.
- [`references/no-favoritism.md`](references/no-favoritism.md) — what the
  open-moderator self-check must assert and the P-claude defense-in-depth
  rule.
