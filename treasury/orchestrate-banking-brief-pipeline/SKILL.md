---
name: orchestrate-banking-brief-pipeline
description: >
  Run the decomposed banking BA brief pipeline across extraction, compliance,
  ambiguity, and Gherkin stages with validated JSON handoffs and pluggable model
  commands. Use when a user asks to "run the decomposed banking brief pipeline",
  "chain the BA skills", or "orchestrate extraction compliance ambiguity and
  AC generation". Do NOT use for a single-stage skill run or for direct vendor
  API integration.
---

# Skill: Orchestrate Banking Brief Pipeline

## Purpose

Coordinate the decomposed banking brief skills while keeping deterministic validation, traceability, and failure handling outside the model prompts.

## Input

Accept a JSON file matching `templates/pipeline_input.schema.json`. The required fields are `raw_content` and `idempotency_key`; optional fields include `source_type`, `tier_hint`, `requested_at`, and `project_context`.

Before running, read `references/orchestration-contract.md`.

## Procedure

1. Validate the pipeline input before calling any stage.
2. Run the stages in this exact order: extraction, compliance, ambiguity, Gherkin.
3. For each stage, write a stage input JSON containing the original input, all prior validated stage outputs, and the current pipeline metadata.
4. Invoke the configured stage command with placeholders for `{input}`, `{output}`, `{skill_dir}`, `{stage}`, and `{state}`.
5. Load the stage output, validate it against that stage skill's schema, and append it to `pipeline_state.json`.
6. Halt on command failure, invalid JSON, schema failure, missing output file, or `stage_status: "blocked"`. Write `failure.json` with the stage name and evidence.
7. Assemble `output.json` only after all stages pass. Mark `blocks_tl_handoff` true when compliance blockers or P1 open questions remain.

## Output Contract

The orchestrator writes:

- `pipeline_state.json` after every stage attempt.
- `output.json` after all four stages complete.
- `failure.json` on terminal failure.

Run `scripts/run_elicitation_pipeline.py` for local execution. It is provider-agnostic and does not call a model API directly.

## References

- `references/orchestration-contract.md` - command placeholders, state shape, and failure behavior.
- `templates/pipeline_input.schema.json` - input contract.
- `templates/pipeline_output.schema.json` - assembled output contract.
- `examples/pipeline_input_example.json` - compact valid input example.
- `scripts/run_elicitation_pipeline.py` - deterministic runner.
