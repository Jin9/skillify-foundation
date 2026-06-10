#!/usr/bin/env python3
"""Run the decomposed banking brief pipeline with pluggable stage commands.

The runner is deliberately provider-agnostic. It writes each stage input file,
executes a user-supplied command, validates the JSON the command writes, and
mutates only the orchestrator-owned pipeline state.
"""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

try:
    import jsonschema
except Exception:  # pragma: no cover - exercised only in minimal Python envs
    jsonschema = None


STAGES = [
    {
        "key": "extract",
        "stage": "extraction",
        "skill_dir": "extract-brief-structure",
        "schema": "templates/stage1_extraction.schema.json",
        "output": "stage1_extraction.json"
    },
    {
        "key": "compliance",
        "stage": "compliance",
        "skill_dir": "evaluate-banking-compliance",
        "schema": "templates/stage2_compliance.schema.json",
        "output": "stage2_compliance.json"
    },
    {
        "key": "ambiguity",
        "stage": "ambiguity",
        "skill_dir": "sweep-ambiguities",
        "schema": "templates/stage3_ambiguity.schema.json",
        "output": "stage3_ambiguity.json"
    },
    {
        "key": "gherkin",
        "stage": "gherkin",
        "skill_dir": "generating-gherkin-acceptance-criteria",
        "schema": "templates/stage4_gherkin.schema.json",
        "output": "stage4_gherkin.json"
    }
]


def load_json(path: Path) -> Any:
    with path.open("r", encoding="utf-8") as handle:
        return json.load(handle)


def write_json(path: Path, data: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8") as handle:
        json.dump(data, handle, indent=2, sort_keys=True)
        handle.write("\n")


def utc_now() -> str:
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat()


def validate_json(data: Any, schema_path: Path) -> list[str]:
    schema = load_json(schema_path)
    if jsonschema is None:
        errors = []
        for key in schema.get("required", []):
            if key not in data:
                errors.append(f"missing required field: {key}")
        return errors

    validator_cls = jsonschema.validators.validator_for(schema)
    validator_cls.check_schema(schema)
    validator = validator_cls(schema)
    errors = []
    for error in sorted(validator.iter_errors(data), key=lambda item: list(item.path)):
        path = ".".join(str(part) for part in error.absolute_path) or "$"
        errors.append(f"{path}: {error.message}")
    return errors


def parse_stage_commands(values: list[str]) -> dict[str, str]:
    commands: dict[str, str] = {}
    for value in values:
        if "=" not in value:
            raise ValueError(f"stage command must be key=command: {value}")
        key, command = value.split("=", 1)
        key = key.strip()
        if key not in {stage["key"] for stage in STAGES}:
            raise ValueError(f"unknown stage command key: {key}")
        if not command.strip():
            raise ValueError(f"empty command for stage: {key}")
        commands[key] = command
    return commands


def render_command(command: str, replacements: dict[str, Path | str]) -> str:
    rendered = command
    for key, value in replacements.items():
        rendered = rendered.replace("{" + key + "}", str(value))
    return rendered


def failure_payload(
    state: dict[str, Any],
    failed_stage: str,
    failure_code: str,
    failure_reason: str,
    evidence: Any | None = None
) -> dict[str, Any]:
    return {
        "pipeline_status": "failed",
        "failed_stage": failed_stage,
        "failure_code": failure_code,
        "failure_reason": failure_reason,
        "evidence": evidence,
        "state_path": state.get("state_path"),
        "created_at": utc_now()
    }


def write_failure(output_dir: Path, state: dict[str, Any], payload: dict[str, Any]) -> int:
    state["pipeline_status"] = "failed"
    state["errors"].append(payload)
    write_json(output_dir / "pipeline_state.json", state)
    write_json(output_dir / "failure.json", payload)
    return 1


def stage_input_payload(state: dict[str, Any], stage_name: str) -> dict[str, Any]:
    return {
        "stage": stage_name,
        "input": state["input"],
        "stage_outputs": state["stage_outputs"],
        "processing_metadata": {
            "pipeline_id": state["pipeline_id"],
            "stage_order": [stage["stage"] for stage in STAGES],
            "current_stage": stage_name
        }
    }


def assemble_output(state: dict[str, Any]) -> dict[str, Any]:
    extraction = state["stage_outputs"]["extraction"]
    compliance = state["stage_outputs"]["compliance"]
    ambiguity = state["stage_outputs"]["ambiguity"]
    gherkin = state["stage_outputs"]["gherkin"]

    p1_open_questions = [
        item for item in ambiguity.get("open_questions", [])
        if item.get("severity") == "P1"
    ]
    blocks_tl_handoff = bool(
        compliance.get("blocks_tl_handoff") or p1_open_questions
    )

    return {
        "output_type": "blocked_partial_brief" if blocks_tl_handoff else "brief",
        "blocks_tl_handoff": blocks_tl_handoff,
        "idempotency_key": state["input"]["idempotency_key"],
        "source_type": extraction.get("source_type"),
        "scope_kind": extraction.get("scope_kind"),
        "epics": extraction.get("epics", []),
        "stories": gherkin.get("stories", []),
        "out_of_scope_deferred": extraction.get("out_of_scope_deferred", []),
        "glossary_candidates": extraction.get("glossary_candidates", []),
        "pii_inventory": compliance.get("pii_inventory", []),
        "stakeholders": compliance.get("stakeholders", []),
        "legal_status_by_epic": compliance.get("legal_status_by_epic", []),
        "regulatory_dependencies": compliance.get("regulatory_dependencies", []),
        "governance_gaps": compliance.get("governance_gaps", []),
        "open_questions": ambiguity.get("open_questions", []),
        "assumptions_made": ambiguity.get("assumptions_made", []),
        "processing_metadata": {
            "pipeline_status": "complete",
            "pipeline_id": state["pipeline_id"],
            "stage_order": [stage["stage"] for stage in STAGES],
            "stage_attempts": state["stage_attempts"],
            "extraction": extraction.get("processing_metadata", {}),
            "ambiguity": ambiguity.get("processing_metadata", {}),
            "gherkin": gherkin.get("processing_metadata", {})
        }
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", required=True, help="Pipeline input JSON path.")
    parser.add_argument("--output-dir", required=True, help="Directory for state and outputs.")
    parser.add_argument(
        "--stage-command",
        action="append",
        default=[],
        help="Stage command in key=command form. Keys: extract, compliance, ambiguity, gherkin."
    )
    parser.add_argument(
        "--suite-root",
        help="Pipeline suite root (directory containing the four stage skill folders). Defaults to this script's grandparent."
    )
    args = parser.parse_args()

    script_dir = Path(__file__).resolve().parent
    suite_root = Path(args.suite_root).resolve() if args.suite_root else script_dir.parents[1]
    output_dir = Path(args.output_dir).resolve()
    output_dir.mkdir(parents=True, exist_ok=True)

    try:
        commands = parse_stage_commands(args.stage_command)
    except ValueError as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2

    missing = [stage["key"] for stage in STAGES if stage["key"] not in commands]
    if missing:
        print(f"error: missing stage command(s): {', '.join(missing)}", file=sys.stderr)
        return 2

    input_path = Path(args.input).resolve()
    try:
        input_data = load_json(input_path)
    except (OSError, json.JSONDecodeError) as exc:
        print(f"error: cannot load input JSON: {exc}", file=sys.stderr)
        return 2

    input_schema = script_dir.parent / "templates" / "pipeline_input.schema.json"
    input_errors = validate_json(input_data, input_schema)
    raw_idempotency_key = ""
    if isinstance(input_data, dict):
        raw_idempotency_key = str(input_data.get("idempotency_key", "invalid-input"))

    state: dict[str, Any] = {
        "pipeline_id": raw_idempotency_key[:8],
        "pipeline_status": "running",
        "state_path": str(output_dir / "pipeline_state.json"),
        "created_at": utc_now(),
        "input": input_data,
        "stage_outputs": {},
        "stage_attempts": [],
        "current_stage": None,
        "errors": []
    }
    write_json(output_dir / "pipeline_state.json", state)

    if input_errors:
        return write_failure(
            output_dir,
            state,
            failure_payload(state, "input", "input_validation_failed", "Pipeline input failed validation.", input_errors)
        )

    for stage in STAGES:
        key = stage["key"]
        stage_name = stage["stage"]
        stage_dir = suite_root / stage["skill_dir"]
        schema_path = stage_dir / stage["schema"]
        stage_input = output_dir / f"{stage_name}_input.json"
        stage_output = output_dir / stage["output"]
        state["current_stage"] = stage_name
        write_json(stage_input, stage_input_payload(state, stage_name))
        write_json(output_dir / "pipeline_state.json", state)

        command = render_command(commands[key], {
            "input": stage_input,
            "output": stage_output,
            "skill_dir": stage_dir,
            "stage": stage_name,
            "state": output_dir / "pipeline_state.json"
        })
        attempt = {
            "stage": stage_name,
            "command": command,
            "started_at": utc_now()
        }
        completed = subprocess.run(command, shell=True, text=True, capture_output=True)
        attempt["finished_at"] = utc_now()
        attempt["returncode"] = completed.returncode
        attempt["stdout"] = completed.stdout[-4000:]
        attempt["stderr"] = completed.stderr[-4000:]
        state["stage_attempts"].append(attempt)

        if completed.returncode != 0:
            return write_failure(
                output_dir,
                state,
                failure_payload(state, stage_name, "stage_command_failed", "Stage command exited non-zero.", attempt)
            )
        if not stage_output.exists():
            return write_failure(
                output_dir,
                state,
                failure_payload(state, stage_name, "stage_output_missing", "Stage command did not write the expected output file.", str(stage_output))
            )
        try:
            stage_data = load_json(stage_output)
        except json.JSONDecodeError as exc:
            return write_failure(
                output_dir,
                state,
                failure_payload(state, stage_name, "stage_output_invalid_json", "Stage output is not valid JSON.", str(exc))
            )

        errors = validate_json(stage_data, schema_path)
        if errors:
            return write_failure(
                output_dir,
                state,
                failure_payload(state, stage_name, "stage_schema_validation_failed", "Stage output failed schema validation.", errors)
            )

        state["stage_outputs"][stage_name] = stage_data
        write_json(output_dir / "pipeline_state.json", state)

        if stage_data.get("stage_status") == "blocked":
            return write_failure(
                output_dir,
                state,
                failure_payload(
                    state,
                    stage_name,
                    stage_data.get("failure_code", "stage_blocked"),
                    stage_data.get("failure_reason", "Stage returned blocked status."),
                    stage_data
                )
            )

    state["current_stage"] = None
    state["pipeline_status"] = "complete"
    output = assemble_output(state)
    output_errors = validate_json(output, script_dir.parent / "templates" / "pipeline_output.schema.json")
    if output_errors:
        return write_failure(
            output_dir,
            state,
            failure_payload(state, "assembly", "pipeline_output_validation_failed", "Assembled output failed validation.", output_errors)
        )

    write_json(output_dir / "pipeline_state.json", state)
    write_json(output_dir / "output.json", output)
    print(f"ok: wrote {output_dir / 'output.json'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
