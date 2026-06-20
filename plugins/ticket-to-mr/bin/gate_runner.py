#!/usr/bin/env python3
"""
gate_runner.py — reference orchestrator for the ticket-to-mr gate state machine.

Drives one Jira issue through gates 1..4. Each transition is recorded as a
hash-chained entry in runs/<trace>/log.jsonl (tamper-evident audit trail), and
each gate's checkpoint payload is validated against core/schemas/gate<N>-*.json
before it is accepted. Caps from core/policy/policy.json terminate the run on
breach. This is the *reference* implementation of the discipline the SKILL.md
describes — the agent follows the same gate order; this module makes the state,
the audit chain, and the caps real and inspectable.

Subcommands:
  start   --key KEY --trace TRACE        begin a run; write genesis log entry + runs/.active marker
  submit  --gate N --payload FILE        validate + record a gate checkpoint (enforces caps)
  approve --gate N --approver NAME       record a named human approval; advance state
  status  --trace TRACE                  print current state + verify the hash chain

Hash chain: entry_hash = sha256(prev_hash + event + canonical_json(data)). The
genesis prev_hash is 64 zeros. `ts` is stored on each entry but excluded from the
hash so a fixed-input replay is deterministic (testable).

Self-approval guard: an approver value of "agent"/"assistant"/"self"/"ai" (any
case) is rejected — gates require a named human.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
import time
from pathlib import Path

PLUGIN_ROOT = Path(__file__).resolve().parents[1]
POLICY_PATH = PLUGIN_ROOT / "core" / "policy" / "policy.json"
SCHEMA_DIR = PLUGIN_ROOT / "core" / "schemas"
RUNS_DIR = PLUGIN_ROOT / "runs"
ACTIVE_MARKER = RUNS_DIR / ".active"
GENESIS = "0" * 64

GATE_SCHEMA = {
    1: "gate1-analyze.json",
    2: "gate2-implement.json",
    3: "gate3-mr.json",
    4: "gate4-jira-comment.json",
}
_SELF_APPROVAL = {"agent", "assistant", "self", "ai", "bot", "claude", "model"}
_JSON_TYPES = {
    "object": dict,
    "array": list,
    "string": str,
    "number": (int, float),
    "boolean": bool,
}


# --------------------------------------------------------------------------- #
# hashing / log
# --------------------------------------------------------------------------- #
def canonical(data) -> str:
    return json.dumps(data, sort_keys=True, separators=(",", ":"), ensure_ascii=False)


def entry_hash(prev_hash: str, event: str, data) -> str:
    return hashlib.sha256((prev_hash + event + canonical(data)).encode("utf-8")).hexdigest()


def _log_path(trace: str) -> Path:
    return RUNS_DIR / trace / "log.jsonl"


def _state_path(trace: str) -> Path:
    return RUNS_DIR / trace / "state.json"


def read_log(trace: str) -> list[dict]:
    path = _log_path(trace)
    if not path.exists():
        return []
    return [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines() if line.strip()]


def append(trace: str, event: str, data: dict, ts: float | None = None) -> dict:
    """Append a hash-chained entry. `ts` is injectable for deterministic tests."""
    log = read_log(trace)
    prev = log[-1]["hash"] if log else GENESIS
    record = {
        "ts": ts if ts is not None else time.time(),
        "event": event,
        "data": data,
        "prev_hash": prev,
        "hash": entry_hash(prev, event, data),
    }
    path = _log_path(trace)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("a", encoding="utf-8") as fh:
        fh.write(json.dumps(record, ensure_ascii=False) + "\n")
    return record


def verify_chain(trace: str) -> tuple[bool, int]:
    """Return (ok, first_bad_index). first_bad_index == -1 when ok."""
    prev = GENESIS
    for i, rec in enumerate(read_log(trace)):
        if rec.get("prev_hash") != prev:
            return False, i
        if rec.get("hash") != entry_hash(prev, rec.get("event", ""), rec.get("data")):
            return False, i
        prev = rec["hash"]
    return True, -1


# --------------------------------------------------------------------------- #
# minimal JSON-Schema (draft-07 subset: type / required / properties)
# --------------------------------------------------------------------------- #
def validate(payload, schema) -> list[str]:
    errors: list[str] = []

    def _check(value, sch, path):
        expected = sch.get("type")
        if expected and expected in _JSON_TYPES:
            py = _JSON_TYPES[expected]
            ok = isinstance(value, py) and not (expected == "number" and isinstance(value, bool))
            if not ok:
                errors.append(f"{path or '<root>'}: expected {expected}, got {type(value).__name__}")
                return
        if sch.get("type") == "object":
            for req in sch.get("required", []):
                if not isinstance(value, dict) or req not in value:
                    errors.append(f"{path or '<root>'}: missing required field '{req}'")
            for prop, psch in sch.get("properties", {}).items():
                if isinstance(value, dict) and prop in value:
                    _check(value[prop], psch, f"{path}.{prop}" if path else prop)

    _check(payload, schema, "")
    return errors


def load_schema(gate: int) -> dict:
    return json.loads((SCHEMA_DIR / GATE_SCHEMA[gate]).read_text(encoding="utf-8"))


def load_caps() -> dict:
    return json.loads(POLICY_PATH.read_text(encoding="utf-8")).get("caps", {})


def check_caps(gate: int, payload: dict) -> list[str]:
    """Cap enforcement that applies at submit time (currently the gate-2 diff caps)."""
    caps = load_caps()
    breaches: list[str] = []
    if gate == 2:
        files = payload.get("files_changed", [])
        if isinstance(files, list) and len(files) > caps.get("diff_max_files", 1e9):
            breaches.append(f"diff_max_files exceeded: {len(files)} > {caps['diff_max_files']}")
        lines = payload.get("diff_lines", 0)
        if isinstance(lines, (int, float)) and lines > caps.get("diff_max_lines", 1e9):
            breaches.append(f"diff_max_lines exceeded: {lines} > {caps['diff_max_lines']}")
    return breaches


# --------------------------------------------------------------------------- #
# subcommands
# --------------------------------------------------------------------------- #
def cmd_start(args) -> int:
    state = {"key": args.key, "trace": args.trace, "gate": 1, "approved": [], "status": "open"}
    _state_path(args.trace).parent.mkdir(parents=True, exist_ok=True)
    _state_path(args.trace).write_text(json.dumps(state, indent=2), encoding="utf-8")
    RUNS_DIR.mkdir(parents=True, exist_ok=True)
    ACTIVE_MARKER.write_text(args.trace + "\n", encoding="utf-8")
    append(args.trace, "start", {"key": args.key, "trace": args.trace})
    print(f"OK: run {args.trace} started for {args.key}. Enforcement active.")
    print("Hint: export TICKET_TO_MR_ENFORCE=1  (the PreToolUse hook also honors runs/.active)")
    return 0


def cmd_submit(args) -> int:
    payload = json.loads(Path(args.payload).read_text(encoding="utf-8"))
    schema_errors = validate(payload, load_schema(args.gate))
    if schema_errors:
        print(f"FAIL: gate {args.gate} payload failed schema validation:", file=sys.stderr)
        for e in schema_errors:
            print(f"  - {e}", file=sys.stderr)
        return 1
    cap_breaches = check_caps(args.gate, payload)
    if cap_breaches:
        print(f"FAIL: gate {args.gate} payload breaches caps (STOP and escalate):", file=sys.stderr)
        for b in cap_breaches:
            print(f"  - {b}", file=sys.stderr)
        append(args.trace, f"gate{args.gate}_cap_breach", {"breaches": cap_breaches})
        return 1
    append(args.trace, f"gate{args.gate}_submit", payload)
    print(f"OK: gate {args.gate} checkpoint recorded for {args.trace}.")
    return 0


def cmd_approve(args) -> int:
    if args.approver.strip().lower() in _SELF_APPROVAL:
        print(f"FAIL: '{args.approver}' is not a named human; self-approval is forbidden.", file=sys.stderr)
        return 1
    state = json.loads(_state_path(args.trace).read_text(encoding="utf-8"))
    state.setdefault("approved", []).append({"gate": args.gate, "approver": args.approver})
    state["gate"] = min(args.gate + 1, 4)
    if args.gate >= 4:
        state["status"] = "complete"
    _state_path(args.trace).write_text(json.dumps(state, indent=2), encoding="utf-8")
    append(args.trace, f"gate{args.gate}_approve", {"approver": args.approver})
    print(f"OK: gate {args.gate} approved by {args.approver}. Next gate: {state['gate']}.")
    return 0


def cmd_status(args) -> int:
    if not _state_path(args.trace).exists():
        print(f"FAIL: no run found for trace {args.trace}", file=sys.stderr)
        return 1
    state = json.loads(_state_path(args.trace).read_text(encoding="utf-8"))
    ok, bad = verify_chain(args.trace)
    print(json.dumps(state, indent=2))
    print(f"chain: {'OK' if ok else f'TAMPERED at entry {bad}'} ({len(read_log(args.trace))} entries)")
    return 0 if ok else 1


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="ticket-to-mr gate state machine")
    sub = parser.add_subparsers(dest="cmd", required=True)

    p = sub.add_parser("start"); p.add_argument("--key", required=True); p.add_argument("--trace", required=True)
    p.set_defaults(func=cmd_start)
    p = sub.add_parser("submit"); p.add_argument("--trace", required=True)
    p.add_argument("--gate", type=int, required=True, choices=[1, 2, 3, 4]); p.add_argument("--payload", required=True)
    p.set_defaults(func=cmd_submit)
    p = sub.add_parser("approve"); p.add_argument("--trace", required=True)
    p.add_argument("--gate", type=int, required=True, choices=[1, 2, 3, 4]); p.add_argument("--approver", required=True)
    p.set_defaults(func=cmd_approve)
    p = sub.add_parser("status"); p.add_argument("--trace", required=True)
    p.set_defaults(func=cmd_status)

    args = parser.parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
