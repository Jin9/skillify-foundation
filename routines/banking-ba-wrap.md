---
routine: banking-ba-wrap
summary: Wrap the two internally gated banking BA runners as black-box nodes - brief pipeline bundle first, then the five-stage BA pipeline to a Tech-Lead handoff. Their internal gates still fire.
version: 0.1.0
default_on_fail: stop
owner: chinnawat.w
tags: banking, business-analysis
---

Registers the existing hardcoded banking runners as routine nodes. Each executor is itself a multi-stage, internally gated runner skill — a skill, not a routine, so the no-recursion rule holds. This routine only sequences the two and verifies their bundle artifacts; every internal gate of each runner still fires with its named human owner.

## Inputs

- raw_brief: the raw banking requirement or brief text, or a path to it (required)

## Nodes

### Node: brief

- purpose: Run the decomposed banking brief pipeline to an extraction, compliance, ambiguity, and Gherkin bundle
- executor: orchestrate-banking-brief-pipeline
- tier: mid
- inputs: user.raw_brief
- outputs: 01-banking-brief-bundle.md
- gate: after
- notes: The executor expects a pipeline-input JSON with raw_content and idempotency_key - the agent assembles it from user.raw_brief. The node artifact embeds the assembled output.json bundle and its blocks_tl_handoff flag; the executor's own stage validation and failure.json handling still apply.

### Node: ba-handoff

- purpose: Take the brief bundle through the five-stage BA pipeline to a Tech-Lead-ready handoff
- executor: running-ba-pipeline
- tier: mid
- inputs: 01-banking-brief-bundle.md
- outputs: 02-tl-handoff.md
- gate: after
- notes: Consumes the brief bundle as the raw requirement. The executor's own gates G1-G5 still fire with their named human owners; its map is documented to run at small-mid tier and stage-level model routing is internal to the executor. Artifact embeds the TL Handoff Bundle summary and pointers to its run directory.
