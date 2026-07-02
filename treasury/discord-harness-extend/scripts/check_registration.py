#!/usr/bin/env python3
"""Cross-check the harness tool registry against each config profile's `tools:` classes.

The harness decouples a tool's behavior (a @tool handler in tools/) from its authorization (an
ALLOW/CONFIRM/DENY entry under `tools:` in config). Miss the config entry and the tool is silently
DENY; mis-case the class (e.g. `allow`) and it is ALSO silently DENY. This script makes the resulting
tool x class matrix visible per profile, and FAILS on the two unambiguous mistakes:

  * orphan  — a `tools:` entry whose name has no registered handler (typo / removed handler)
  * invalid — a class value that is not exactly ALLOW / CONFIRM / DENY (silently becomes DENY)

A registered handler being DENY (unlisted) in a profile is NOT an error — the default profile
intentionally leaves operator tools unlisted so they fail closed. The matrix lets you eyeball that your
new tool landed in the class you intended.

Run from anywhere (it locates the repo root). Use the repo's .venv so pydantic/yaml are importable:
    .venv/bin/python .claude/skills/discord-harness-extend/scripts/check_registration.py
"""

from __future__ import annotations

import os
import sys
from pathlib import Path

VALID = {"ALLOW", "CONFIRM", "DENY"}
PROFILES = ["config.yaml"]


def find_repo_root() -> Path:
    for d in [
        Path.cwd(),
        *Path.cwd().parents,
        Path(__file__).resolve().parent,
        *Path(__file__).resolve().parents,
    ]:
        if (d / "config.yaml").is_file() and (d / "tools" / "__init__.py").is_file():
            return d
    print(
        "error: could not locate the harness repo root (needs config.yaml + tools/__init__.py)",
        file=sys.stderr,
    )
    raise SystemExit(2)


def main() -> int:
    root = find_repo_root()
    sys.path.insert(0, str(root))
    try:
        import yaml
        import tools  # registers every handler via import side effects
        import policy
    except Exception as exc:  # noqa: BLE001
        print(
            f"error: import failed ({type(exc).__name__}: {exc}). Run with the repo's .venv python.",
            file=sys.stderr,
        )
        return 2

    handlers = set(tools.REGISTRY)
    problems = 0

    for profile in PROFILES:
        path = root / profile
        if not path.is_file():
            continue
        raw = (yaml.safe_load(path.read_text(encoding="utf-8")) or {}).get("tools") or {}
        os.environ["HARNESS_CONFIG"] = profile
        policy.load_config.cache_clear()  # re-read under the selected profile

        by_class: dict[str, list[str]] = {"ALLOW": [], "CONFIRM": [], "DENY": []}
        for name in sorted(handlers):
            by_class[policy.classify(name).value].append(name)

        orphans = sorted(k for k in raw if k not in handlers)
        invalid = sorted(f"{k}={raw[k]!r}" for k in raw if raw[k] not in VALID)

        print(f"profile: {profile}  ({len(handlers)} registered handlers)")
        for cls in ("ALLOW", "CONFIRM", "DENY"):
            print(f"  {cls:<7}: {', '.join(by_class[cls]) or '(none)'}")
        print(f"  orphans (config entry, no handler): {', '.join(orphans) or '(none)'}")
        print(f"  invalid class value (-> silently DENY): {', '.join(invalid) or '(none)'}")
        print()
        problems += len(orphans) + len(invalid)

    if problems:
        print(f"FAIL: {problems} registration problem(s) found.", file=sys.stderr)
        return 1
    print(
        "OK: every config profile's tools: classes resolve to a registered handler and a valid class."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
