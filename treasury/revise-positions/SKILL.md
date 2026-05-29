---
name: revise-positions
description: >
  Stage 4 of the squad-brainstorm workflow, run once per panel CLI per
  debate round (skipped when rounds=quick). The invoked panelist reads ONLY
  the critiques aimed at itself, then produces a revised position: concede
  where the rebuttal lands, defend with grounding-pack evidence where it does
  not, and record a per-objection changelog plus an honest stance_delta.
  Principled convergence is allowed and tracked; capitulation without reason
  and unprincipled flip-flopping are forbidden. Use when the workflow runs the revise-positions
  stage of workflows/brainstorm.yaml. One panelist, one round, one revised
  position object. Do NOT critique others (cross-examine) or merge the
  debate (synthesize-consensus).
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # the prompt calls no tools; the driver invokes the external CLI
version: 0.1.0
stage: revise-positions
workflow: brainstorm
inputs:
  - { name: debate_brief,       type: object,  required: true, source: stages.frame-debate.debate_brief,   description: "from frame-debate" }
  - { name: grounding_pack,     type: object,  required: true, source: stages.frame-debate.grounding_pack, description: "closed evidence world; cite [g-n] only" }
  - { name: panelist,           type: object,  required: true, source: fanout.panelist,                    description: "the revising panelist's identity row" }
  - { name: prior_position,     type: object,  required: true, source: "self.position_for[panelist]",        description: "this panelist's own latest position (round_no-1)" }
  - { name: critiques_received, type: array,   required: true, source: stages.cross-examine.critiques, description: "ONLY the critiques whose target_panelist_id == this panelist (driver-filtered)" }
  - { name: round_no,           type: integer, required: true, source: loop.round_no,                       description: "1-based debate round" }
outputs:
  - { name: revised_positions, type: array, description: "collected per round; the final round's set is the input to synthesis" }
---

# Skill: revise-positions

## Purpose

Each panelist updates its own position in light of the critiques it
received: principled convergence where rebuttals land, evidenced defense
where they do not. The per-objection `changelog` and the `stance_delta` make
movement *real and auditable* — so the convergence map in
`report-debate` reflects what actually happened, not asserted agreement.

This is where a healthy debate either converges (rebuttals landed and
positions move toward the evidence) or stays honestly divergent (defenses
held with grounding). Both are valid outcomes; an *unprincipled* move in
either direction is not.

**Atomicity:** one panelist, one round, one CLI call → one revised position.
Skipped when `rounds == quick`.

## When to use this skill

- Stage 4 fan-out dispatch, per round, after the matching `cross-examine`
  round, when `rounds != quick`.
- Prompt: "revise your position given the critiques aimed at you, round N".

Do NOT use this skill to:
- Critique other panelists — stage 3.
- Merge the debate — stage 5.
- Respond to critiques aimed at someone else (the driver filtered those out;
  you only see your own).

## Inputs

| Name | Type | Notes |
|------|------|-------|
| `debate_brief` | object | `rubric` is still the scoring lens. |
| `grounding_pack` | object | Defense evidence. Closed-world `[g-n]`. |
| `panelist` | object | Your identity row. Stamp it verbatim. |
| `prior_position` | object | Your own latest position (the driver supplies the correct round's artifact). |
| `critiques_received` | array | Only the critique entries whose `target_panelist_id` is you. You never revise based on a critique aimed at a peer. |
| `round_no` | integer | 1-based. Echo it. |

## Output contract

Same shape as an `opening-debate-panel` (`panel-open` stage) position (so synthesis takes a uniform input)
**plus** a `changelog` and `stance_delta`:

```jsonc
{
  "panelist_id": "P-gemini",
  "cli": "gemini -p", "model": "gemini-2.5-pro",
  "round": 1,
  "stance": "string",                  // may equal or differ from prior_position.stance
  "position": "markdown",              // the REVISED case; [g-n] cites adjacent to claims
  "answers": [ { "cq_id": "cq1", "answer": "string", "evidence_refs": ["g-2"] } ],
  "key_claims": [ { "claim": "string", "grounding": ["g-2"], "confidence": "high|medium|low" } ],
  "changelog": [
    {
      "in_response_to": { "from_panelist_id": "P-codex", "objection": "string" },
      "action": "conceded|partially_conceded|reframed|defended",
      "what_changed": "string",        // "" only when action == "defended"
      "evidence_refs": ["g-5"]         // REQUIRED (non-empty) when action == "defended"
    }
  ],
  "stance_delta": "unchanged|narrowed|broadened|reversed",   // feeds the convergence map
  "uncited_assertions": ["string"],
  "status": "ok"
}
```

## Procedure

Single CLI call:

1. **Read** your `prior_position` and every entry in `critiques_received`.
2. **For each received objection**, decide an `action`
   (see [`references/convergence-discipline.md`](references/convergence-discipline.md)):
   - `conceded` — the rebuttal lands; remove/withdraw the claim.
   - `partially_conceded` — narrow or caveat the claim.
   - `reframed` — the underlying point survives but needs restating to
     answer the objection.
   - `defended` — the objection fails; rebut it with `[g-n]` evidence
     (`evidence_refs` non-empty, `what_changed: ""`). A bare "I disagree"
     is not a defense.
   Severity guidance: a `major`/`decisive` rebuttal MUST get a changelog
   entry; on `deep`, *every* received rebuttal of severity ≥ `major` needs
   one.
3. **Rewrite the `position`** so it actually reflects the concessions —
   a `conceded` point must be absent or softened in the new prose, not
   merely annotated in the changelog.
4. **Set `stance_delta` honestly** vs `prior_position.stance`
   (`unchanged` / `narrowed` / `broadened` / `reversed`). It must be
   consistent with the changelog (all-`defended` ⟹ usually `unchanged`;
   any `conceded` ⟹ not `unchanged`).
5. **Stamp identity** from `panelist`; set `round` = `round_no`.
6. **Self-check** the Validation gate; fix in place.
7. **Emit** the single JSON object. No prose around it.

## Rounds dial

Runs `debate_brief.exchange_count` times (`quick` 0 / `standard` 1 /
`deep` 2). On `deep`, a changelog entry is required for every received
rebuttal of severity ≥ `major`. `audience` is unused.

## Constraints

- **Respond only to critiques aimed at you.** `critiques_received` is
  already filtered; do not invent peer critiques.
- **Principled convergence only.** Forbidden: conceding a point under no
  real rebuttal pressure just to look agreeable (capitulation); reversing a
  previously well-defended point with no *new* landed objection
  (flip-flop). Both are scored failures and visible to the moderator.
- **Closed evidence world.** Defenses cite `[g-n]`; ungrounded steps go in
  `uncited_assertions`.
- **Identity honesty.** Stamp the given `panelist_id` verbatim.
- **One call. JSON only** in workflow context.

## Validation gate

- [ ] `panelist_id` matches `panelist.panelist_id` (byte-for-byte);
      `round == round_no`.
- [ ] Every `changelog.in_response_to.objection` traces to an actual
      objection in `critiques_received`.
- [ ] `action == "defended"` ⟹ `evidence_refs` non-empty and
      `what_changed == ""`.
- [ ] `action == "conceded"` ⟹ the conceded claim is absent or softened in
      the new `position` (not just annotated).
- [ ] Every received rebuttal with `severity ∈ {major, decisive}` has a
      changelog entry (always on `deep`; at minimum the decisive ones on
      `standard`).
- [ ] `stance_delta` is consistent with the changelog.
- [ ] Every `[g-n]` id used exists in `grounding_pack.items`.

## Worked example

See [`examples/round1-revise.md`](examples/round1-revise.md): `P-gemini`
conceding one P-claude rebuttal, defending another with `[g-n]`, ending at
`stance_delta: narrowed`.

## References

- [`references/convergence-discipline.md`](references/convergence-discipline.md)
  — the concede/partially-concede/reframe/defend decision table; precise
  definitions of capitulation vs principled convergence vs flip-flop.
- Closed-world / honesty-ledger rules shared with `opening-debate-panel`:
  the shared grounding rules owned by `opening-debate-panel`.
