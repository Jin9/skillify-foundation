---
name: cross-examine
description: >
  Stage 3 of the squad-brainstorm workflow, run once per panel CLI per
  debate round (skipped entirely when rounds=quick). The invoked panelist
  reads the OTHER two panelists' latest positions and produces a structured,
  rubric-anchored critique of each — concessions it grants, evidence-backed
  rebuttals (cite [g-n] or mark logic-only), unsupported-claim flags, and a
  mandatory steelman of each opponent's strongest point. It does NOT revise
  its own position here (that is revise-positions). Use as the cross-examine
  stage of workflows/brainstorm.yaml. One panelist, one round, one critiques
  object targeting its two peers. Do NOT open a position (opening-debate-panel) or
  merge the debate (synthesize-consensus).
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # the prompt calls no tools; the driver invokes the external CLI
version: 0.1.0
stage: cross-examine
workflow: brainstorm
inputs:
  - { name: debate_brief,   type: object,  required: true, source: stages.frame-debate.debate_brief,   description: "from frame-debate" }
  - { name: grounding_pack, type: object,  required: true, source: stages.frame-debate.grounding_pack, description: "closed evidence world for rebuttals; cite [g-n] only" }
  - { name: panelist,       type: object,  required: true, source: fanout.panelist,                    description: "the critiquing panelist's identity row" }
  - { name: peer_positions, type: array,   required: true, source: stages.revise-positions.revised_positions ?? stages.panel-open.positions, description: "the OTHER two panelists' latest positions this round (driver filters self out)" }
  - { name: round_no,       type: integer, required: true, source: loop.round_no,                       description: "1-based debate round" }
outputs:
  - { name: critiques, type: array, description: "collected: one critique object per (author, round); each targets the author's two peers" }
---

# Skill: cross-examine

## Purpose

Each panelist adversarially tests its two peers' latest positions against
the `debate_brief.rubric` and the shared `[g-n]` evidence. The goal is to
find where a peer's argument is unsupported, internally inconsistent, or
weaker than it looks — and to do so *honestly*, granting what is right
before attacking what is wrong.

Critique only. The panelist does **not** revise its own position here —
that separation is load-bearing (it is the debate analogue of
`extract-findings`' evidence-before-claim discipline). A panelist that
critiques and self-revises in one call anchors on its prior answer and
produces shallow rebuttals.

**Atomicity:** one panelist, one round, one CLI call → one `critiques`
object covering both peers. Skipped entirely when `rounds == quick`.

## When to use this skill

- Stage 3 fan-out dispatch, per debate round, when `rounds != quick`.
- Prompt: "critique the other panelists' positions for round N".

Do NOT use this skill to:
- State or revise your own position — stages 2 / 4.
- Merge the debate — stage 5.
- Critique yourself (you target only the *other* panelists).

## Inputs

| Name | Type | Notes |
|------|------|-------|
| `debate_brief` | object | The `rubric` is the scoring lens for every critique. |
| `grounding_pack` | object | Rebuttal evidence. Still closed-world `[g-n]`. |
| `panelist` | object | The critiquing panelist's identity row. Stamp `author_panelist_id` from `panelist.panelist_id`. |
| `peer_positions` | array | Exactly the other two panelists' latest positions. The driver removes your own; never critique yourself. |
| `round_no` | integer | 1-based. Echo into the output. |

## Output contract

Each fan-out instance emits one object (driver collects three into
`critiques`):

```jsonc
{
  "author_panelist_id": "P-claude",   // who is critiquing — from panelist.panelist_id
  "round_no": 1,
  "critiques": [
    {
      "target_panelist_id": "P-codex",          // who is being critiqued (the join key for revise-positions)
      "concessions": ["string"],                // points you GRANT this peer (state these first)
      "rebuttals": [
        {
          "target_claim": "string",             // quoted/closely paraphrased from the target's position
          "objection": "string",
          "evidence_refs": ["g-4"],             // [g-n] ids; [] ⇒ logic-only (must set logic_only:true)
          "logic_only": false,
          "severity": "minor|major|decisive"
        }
      ],
      "unsupported_claim_flags": ["string"],    // target claims with no grounding and not in its uncited_assertions
      "strongest_point_from_target": "string"   // MANDATORY steelman of this opponent's best point
    }
    // exactly one entry per OTHER panelist (2 entries on a 3-panel run)
  ],
  "status": "ok"                                // "absent" only on failure
}
```

## Procedure

Single CLI call:

1. **Read** the brief (esp. the `rubric`), the grounding pack, and the two
   `peer_positions`.
2. **For each peer**, in `peer_positions` order:
   a. **Concessions first.** List the points you grant. A critique with no
      concessions is almost always strawmanning — forced honest engagement
      is the anti-domination control.
   b. **Rebuttals.** For each weak point: quote the `target_claim`, state
      the `objection`, and back it with `[g-n]` evidence. A rebuttal with
      no grounding is legitimate only as *logic*: set `evidence_refs: []`,
      `logic_only: true`. Never dress a logic-only rebuttal as evidential
      (see the shared grounding rules owned by `opening-debate-panel`).
      Rate `severity` per
      [`references/critique-severity.md`](references/critique-severity.md).
    c. **Flag unsupported claims** the target presented as grounded but did
       not back and did not list in its own `uncited_assertions`.
    d. **Steelman.** State that peer's single strongest point in its
       strongest form (`strongest_point_from_target`). This field is
       mandatory and scored — a missing or weak steelman is a failure.
3. **Stamp** `author_panelist_id` from `panelist`; echo `round_no`.
4. **Self-check** the Validation gate; fix in place.
5. **Emit** the single JSON object. No prose around it.

## Rounds dial

`rounds` controls how many times this stage runs (via
`debate_brief.exchange_count`): `quick` → not run at all; `standard` → once
(round 1); `deep` → twice (rounds 1 and 2). On `deep`, additionally: produce
at least one rebuttal per `debate_brief.contested_questions` entry where the
target took a position. `audience` is unused.

## Constraints

- **Critique only.** No self-revision, no new position text for yourself.
- **Targets are only the other panelists.** Never self-critique.
- **Closed evidence world.** `[g-n]` for evidential rebuttals; logic-only
  rebuttals must be marked.
- **Mandatory steelman per target.** No winning by strawman.
- **Identity honesty.** Stamp the given `author_panelist_id` verbatim.
- **One call. JSON only** in workflow context.

## Validation gate

- [ ] `author_panelist_id == panelist.panelist_id` (byte-for-byte).
- [ ] `critiques[]` has exactly (panel size − 1) entries.
- [ ] Every `target_panelist_id` ≠ `author_panelist_id` and appears in
      `peer_positions`.
- [ ] Every `rebuttals[].evidence_refs` id exists in
      `grounding_pack.items`; `evidence_refs == []` ⟺ `logic_only == true`.
- [ ] Every target entry has a non-empty `strongest_point_from_target`.
- [ ] Every target entry has ≥1 `concessions` OR an explicit one-line note
      in `concessions` stating why none could honestly be granted.
- [ ] `round_no` echoes the input.

## Worked example

See [`examples/round1-cross-exam.md`](examples/round1-cross-exam.md):
`P-claude` critiquing `P-codex` and `P-gemini` in round 1 — concession +
grounded rebuttal + logic-only rebuttal + mandatory steelman.

## References

- [`references/critique-severity.md`](references/critique-severity.md) —
  the `minor` / `major` / `decisive` rubric and when a rebuttal is
  logic-only vs evidential.
- Closed-world / honesty-ledger rules are shared with `opening-debate-panel`:
  the shared grounding rules owned by `opening-debate-panel`.
