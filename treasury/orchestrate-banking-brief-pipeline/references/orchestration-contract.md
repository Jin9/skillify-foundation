# Orchestration Contract

## Stage Order

1. `extraction` via `extract-brief-structure`
2. `compliance` via `evaluate-banking-compliance`
3. `ambiguity` via `sweep-ambiguities`
4. `gherkin` via `generating-gherkin-acceptance-criteria`

The orchestrator must never skip a stage or continue after invalid stage output.

## Stage Command Placeholders

Stage commands are configured as `stage=command` pairs. The runner replaces:

- `{input}` with the stage input JSON path.
- `{output}` with the expected stage output JSON path.
- `{skill_dir}` with the stage skill folder path.
- `{stage}` with the canonical stage name.
- `{state}` with the current `pipeline_state.json` path.

The command must write strict JSON to `{output}`. The orchestrator owns validation and state mutation.

## State Shape

`pipeline_state.json` contains:

- `pipeline_status`
- `input`
- `stage_outputs`
- `stage_attempts`
- `current_stage`
- `errors`

Only the orchestrator mutates this state file. Stage skills receive read-only snapshots.

## Failure Behavior

Halt and write `failure.json` when:

- pipeline input fails validation
- a stage command is missing
- a stage command exits non-zero
- a stage output file is missing
- a stage emits invalid JSON
- a stage output fails its schema
- a stage returns `stage_status: "blocked"`

The failure payload must include `failed_stage`, `failure_code`, `failure_reason`, and enough evidence for a human BA or engineer to retry.

## Assembly Rules

After all stages pass, assemble `output.json` from validated stage payloads:

- carry `epics`, `stories`, `scope_kind`, source details, and extraction notes from extraction
- carry `pii_inventory`, `stakeholders`, `legal_status_by_epic`, `regulatory_dependencies`, `governance_gaps`, and `blocks_tl_handoff` from compliance
- carry `open_questions`, `assumptions_made`, and hidden sweep metadata from ambiguity
- carry final `stories.acceptance_criteria` and banking-grade rows from Gherkin

Set `output_type` to `blocked_partial_brief` when `blocks_tl_handoff` is true or any P1 open question exists. Otherwise set `output_type` to `brief`.
