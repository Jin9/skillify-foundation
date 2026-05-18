#!/usr/bin/env python3
"""governance-policy-generator: default-deny Rego linter (deterministic).

Validates the skill's OWN generated .rego before it is delivered:
  1. it declares a default-deny posture: a `default allow = false`
     (whitespace-tolerant) MUST be present;
  2. it never declares `default allow = true` anywhere;
  3. no `allow` rule has an empty / unconditioned body (an allow block with
     no condition silently grants everything).

No network, no LLM, no policy evaluation, no code execution — stdlib only.
Exit 0 = ok, 1 = policy violation, 2 = usage error.

Usage:
  python3 rego_deny_lint.py templates/policy.rego
"""
import argparse
import re
import sys
from pathlib import Path

# `default allow = false` / `default  allow=false` — whitespace tolerant.
DEFAULT_FALSE = re.compile(r"^\s*default\s+allow\s*=\s*false\s*$")
DEFAULT_TRUE = re.compile(r"^\s*default\s+allow\s*=\s*true\s*$")
# Start of a block-form rule: `allow {` or `allow = true {` (brace at EOL).
ALLOW_BLOCK_START = re.compile(r"^\s*allow\b[^{]*\{\s*$")


def strip_comment(line: str) -> str:
    # Rego line comments start with `#`; no string-aware parsing needed for
    # this lint (we only inspect rule structure, not literals).
    return line.split("#", 1)[0]


def lint(path: Path) -> list[str]:
    problems: list[str] = []
    lines = path.read_text(encoding="utf-8", errors="replace").splitlines()
    code = [strip_comment(ln) for ln in lines]

    has_default_false = any(DEFAULT_FALSE.match(ln) for ln in code)
    default_true_lines = [i + 1 for i, ln in enumerate(code) if DEFAULT_TRUE.match(ln)]

    if not has_default_false:
        problems.append(
            "no `default allow = false` found — policy must be default-deny"
        )
    for ln_no in default_true_lines:
        problems.append(
            f"line {ln_no}: `default allow = true` is forbidden (not default-deny)"
        )

    # Scan block-form allow rules for an empty / unconditioned body.
    i = 0
    n = len(code)
    while i < n:
        if ALLOW_BLOCK_START.match(code[i]):
            start = i + 1
            body: list[str] = []
            j = i + 1
            closed = False
            while j < n:
                stripped = code[j].strip()
                if stripped == "}":
                    closed = True
                    break
                if stripped:
                    body.append(stripped)
                j += 1
            if not closed:
                problems.append(
                    f"line {start}: `allow` block is not closed with `}}`"
                )
            elif not body:
                problems.append(
                    f"line {start}: `allow` rule has an empty/unconditioned "
                    f"body — this grants everything"
                )
            i = j + 1
            continue
        i += 1

    return problems


def main() -> int:
    ap = argparse.ArgumentParser(
        description="Lint a .rego policy for default-deny posture and "
        "unconditioned allow rules."
    )
    ap.add_argument("policy", help="path to the .rego policy file")
    a = ap.parse_args()

    policy_path = Path(a.policy)
    if not policy_path.is_file():
        print(f"error: policy file not found: {a.policy}", file=sys.stderr)
        return 2

    problems = lint(policy_path)
    if problems:
        print(f"FAIL: {policy_path} ({len(problems)} problem(s)):")
        for p in problems:
            print(f"  - {p}")
        print("FAIL — fix the above before delivering the policy.")
        return 1

    print(f"PASS: {policy_path}")
    print("OK — default-deny posture present and every allow rule is conditioned.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
