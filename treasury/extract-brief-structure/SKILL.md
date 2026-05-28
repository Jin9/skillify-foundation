---
name: extract-brief-structure
description: >
  Convert redacted raw BA input from Jira, Slack, meeting notes, emails, or
  mixed prose into a strict Epic and Story JSON skeleton. Use when a user asks
  to "parse this raw Jira ticket into a brief", "extract the epics and stories
  from this Slack thread", or "structure these notes". Do NOT use for
  compliance scans, ambiguity sweeps, or Gherkin acceptance criteria.
---

# Skill: Extract Brief Structure

## Purpose

Extract the first-stage structural skeleton for the decomposed banking brief pipeline: source type, scope kind, epics, stories, deferred scope, glossary candidates, and source evidence.

## Input

Accept either raw redacted BA text or a pipeline state JSON containing `raw_content`, `idempotency_key`, and optional `source_type`, `requested_at`, `tier_hint`, and `project_context`.

Before extracting, read `references/extraction-rules.md` and inspect `templates/stage1_extraction.schema.json`.

## Procedure

1. Classify the source type from durable markers. If no marker is reliable, set `source_type: "mixed"` and `parsing_mode: "generic_prose_fallback"`.
2. Strip any ground-truth or audit annotation block before semantic extraction. If stripping is ambiguous, stop with `stage_status: "blocked"` and `failure_code: "ground_truth_strip_failed"`.
3. Segment the input into workstreams, decisions, acceptance criteria, stakeholders, dates, deferred scope, and unresolved references.
4. Classify `scope_kind`. Use `multi-epic` only when there are at least three distinct workstreams with at least two acceptance criteria or business-rule signals each.
5. Emit one epic per workstream and one story per user-value slice. Split by workflow step, business rule variation, data class, role boundary, happy/error path, CRUD unit, or spike. Never split by UI/API/database layer alone.
6. Preserve unresolved or ambiguous items as extraction notes; do not silently choose a policy value.
7. Emit strict JSON matching `templates/stage1_extraction.schema.json`.

## Output Contract

Return only JSON with:

- `stage: "extraction"`
- `stage_status: "complete"` or `"blocked"`
- `source_type`, `parsing_mode`, `scope_kind`
- `epics[]`, `stories[]`, `out_of_scope_deferred[]`, `glossary_candidates[]`, `extraction_notes[]`
- `processing_metadata.ground_truth_stripped`

On blocked output, include `failure_code`, `failure_reason`, and omit fabricated epics or stories.

## References

- `references/extraction-rules.md` - source detection, scope classification, and story split rules.
- `templates/stage1_extraction.schema.json` - exact stage output contract.
- `examples/stage1_example.json` - compact example of valid stage output.
