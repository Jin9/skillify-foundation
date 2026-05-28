---
name: sweep-ambiguities
description: >
  Detect linguistic ambiguities and hidden requirements in a structured banking
  brief using the eight ambiguity detectors and ten elicitation frames. Use when
  a user asks to "find the hidden requirements in this brief", "scan these
  stories for ambiguity", or "what questions are missing from this spec". Do
  NOT use for strict PII or Legal audit, story extraction, or Gherkin writing.
---

# Skill: Sweep Ambiguities

## Purpose

Surface unclear language, unresolved decisions, and missing requirement classes that should become open questions or explicit assumptions before TL handoff.

## Input

Accept a pipeline state JSON containing completed extraction and compliance results.

Before scanning, read `references/ambiguity-and-frame-rules.md`, then inspect `templates/stage3_ambiguity.schema.json`.

## Procedure

1. Run the eight linguistic detectors over titles, descriptions, source evidence, story context, scope, compliance notes, and customer-facing strings.
2. Preserve every meaningful ambiguity facet. Do not deduplicate a phrase across detector types when each detector exposes a different risk.
3. Assign severity with the documented floors: P1 for handoff blockers, P2 for decisions needed before planning, and P3 for assumptions that can proceed with review.
4. Apply all ten hidden-requirement frames. Record each frame as applied or skipped with a reason; all non-failure outputs must account for frames 1 through 10.
5. Emit hidden-frame gaps as `open_questions[]` when a human decision is required. Emit `assumptions_made[]` only when a defensible default exists and include a revisit trigger.
6. Populate `processing_metadata.hidden_requirements_sweep` with frame coverage and finding counts.
7. Return strict JSON matching `templates/stage3_ambiguity.schema.json`.

## Output Contract

Return only JSON with:

- `stage: "ambiguity"`
- `stage_status: "complete"` or `"blocked"`
- `open_questions[]`, `assumptions_made[]`
- `ambiguity_findings[]`
- `processing_metadata.hidden_requirements_sweep`

If sweep coverage is partial, emit a P2 open question explaining the missing frame or sub-topic.

## References

- `references/ambiguity-and-frame-rules.md` - detector definitions, severity rules, and hidden-requirement frames.
- `templates/stage3_ambiguity.schema.json` - exact stage output contract.
- `examples/stage3_example.json` - compact example of valid stage output.
