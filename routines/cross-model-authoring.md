---
routine: cross-model-authoring
summary: Author a document across model tiers - digest sources via an advisory CLI, synthesize and draft at frontier tier, adversarially review via a second advisory consult, then compact.
version: 0.1.0
default_on_fail: stop
owner: chinnawat.w
tags: authoring, multi-model
requires: at least one external advisory CLI configured for the delegating-to-cli-models skill
---

Generic multi-model authoring pipeline: cheap large-context digestion feeds frontier-tier synthesis and drafting, an independent advisory model critiques the draft, and a small-tier pass compacts the final document. External model output is advisory only — the launching agent adjudicates every handoff.

## Inputs

- source_material: path to, or description of, the source corpus to digest (required)
- brief: one-paragraph statement of what the final document must achieve and for whom (required)

## Nodes

### Node: digest

- purpose: Digest the source corpus into a structured synthesis of claims, themes, and gaps
- executor: delegating-to-cli-models
- executor_mode: research-digest
- tier: mid
- inputs: user.source_material
- outputs: 01-digest.md
- gate: after
- notes: Advisory output - the agent adjudicates and records provenance before the digest is accepted as an artifact.

### Node: synthesize

- purpose: Extract principles and design the document architecture from the digest against the brief
- executor: agent-inline
- tier: frontier
- inputs: 01-digest.md, user.brief
- outputs: 02-synthesis.md
- gate: none

### Node: draft

- purpose: Write the full draft implementing the synthesis
- executor: agent-inline
- tier: frontier
- inputs: 02-synthesis.md
- outputs: 03-draft.md
- gate: after

### Node: adversarial-review

- purpose: Obtain an independent adversarial critique of the draft for gaps, ambiguity, and overclaiming
- executor: delegating-to-cli-models
- executor_mode: consult
- tier: frontier
- inputs: 03-draft.md
- outputs: 04-review.md
- gate: after
- on_fail: continue
- notes: Advisory second opinion from a different model family - a failed or unavailable consult must not kill the run; the human decides at the gate whether to proceed without it.

### Node: compact

- purpose: Apply accepted review points and compress the draft into the final copy-ready document
- executor: agent-inline
- tier: small
- inputs: 03-draft.md, 04-review.md
- outputs: 05-final.md
- gate: after
