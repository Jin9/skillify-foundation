#!/usr/bin/env python3
"""
policy_match.py — command-safety matcher for the ticket-to-mr workflow.

Single source of truth is core/policy/policy.json. A command is classified by
trying, in order: deny -> confirm -> allow -> default. Deny always wins (the
first deny match short-circuits), so an overly-broad allow can never re-permit
something deny already forbade.

Two modes:

  --hook        Claude Code PreToolUse hook. Reads the hook JSON on stdin
                ({"tool_name": "...", "tool_input": {"command": "..."}}). If a
                run is NOT active it prints nothing and exits 0 (asserts no
                decision, so the host's normal permission flow proceeds). If a
                run IS active it maps the classification to a permissionDecision
                of allow|ask|deny and prints the hookSpecificOutput JSON.

  --command C   CLI/debug. Prints "DECISION rule_id: reason" for command C.

A run is "active" when env TICKET_TO_MR_ENFORCE is truthy OR a runs/.active
marker file exists (written by bin/gate_runner.py start). This keeps the
plugin-level hook from policing every Bash call in unrelated sessions.

Exit code is always 0 in --hook mode (the JSON carries the decision). In
--command mode: 0 allow, 1 deny, 2 confirm/ask.
"""

from __future__ import annotations

import argparse
import json
import os
import re
import sys
from pathlib import Path

PLUGIN_ROOT = Path(__file__).resolve().parents[2]
POLICY_PATH = PLUGIN_ROOT / "core" / "policy" / "policy.json"
ACTIVE_MARKER = PLUGIN_ROOT / "runs" / ".active"

_TRUTHY = {"1", "true", "yes", "on"}


def load_policy(path: Path = POLICY_PATH) -> dict:
    with path.open("r", encoding="utf-8") as fh:
        return json.load(fh)


def run_active() -> bool:
    if os.environ.get("TICKET_TO_MR_ENFORCE", "").strip().lower() in _TRUTHY:
        return True
    return ACTIVE_MARKER.exists()


def classify(command: str, policy: dict) -> tuple[str, str, str]:
    """Return (decision, reason, rule_id). decision in {deny, confirm, allow, ask}."""
    text = command or ""
    for bucket, decision in (("deny", "deny"), ("confirm", "confirm"), ("allow", "allow")):
        for rule in policy.get(bucket, []):
            if re.search(rule["pattern"], text, re.IGNORECASE):
                return decision, rule.get("reason", ""), rule.get("id", bucket)
    default = policy.get("default", "ask")
    return default, "No rule matched; default policy applies.", "default"


# Map internal classification to the host permissionDecision vocabulary.
_DECISION_TO_HOOK = {"deny": "deny", "confirm": "ask", "ask": "ask", "allow": "allow"}


def _emit_hook(decision: str, reason: str, rule_id: str) -> None:
    permission = _DECISION_TO_HOOK.get(decision, "ask")
    out = {
        "hookSpecificOutput": {
            "hookEventName": "PreToolUse",
            "permissionDecision": permission,
            "permissionDecisionReason": f"[ticket-to-mr:{rule_id}] {reason}".strip(),
        }
    }
    print(json.dumps(out))


def run_hook() -> int:
    try:
        payload = json.load(sys.stdin)
    except (json.JSONDecodeError, ValueError):
        # Malformed hook input: do not assert a decision, defer to host.
        return 0

    if payload.get("tool_name") != "Bash":
        return 0  # only police Bash
    if not run_active():
        return 0  # inactive: empty stdout, host's normal permissions apply

    command = (payload.get("tool_input") or {}).get("command", "")
    decision, reason, rule_id = classify(command, load_policy())
    _emit_hook(decision, reason, rule_id)
    return 0


def run_command(command: str) -> int:
    decision, reason, rule_id = classify(command, load_policy())
    print(f"{decision.upper()} {rule_id}: {reason}".rstrip())
    return {"allow": 0, "deny": 1}.get(decision, 2)


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="ticket-to-mr command-safety matcher")
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--hook", action="store_true", help="PreToolUse hook mode (stdin/stdout JSON)")
    group.add_argument("--command", metavar="CMD", help="classify a single command and print the decision")
    args = parser.parse_args(argv)

    if args.hook:
        return run_hook()
    return run_command(args.command)


if __name__ == "__main__":
    sys.exit(main())
