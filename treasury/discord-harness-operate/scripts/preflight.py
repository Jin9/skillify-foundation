#!/usr/bin/env python3
"""Read-only pre-flight for running the harness. Reports the PRESENCE (never the value) of the env vars
the active config needs, confirms the config parses, and warns if authorization is locked shut.

It reads .env only to learn which key NAMES have a non-empty value — secret values are never printed or
retained. Selects the profile via HARNESS_CONFIG, like the bot.

Run from the repo root (use the .venv so PyYAML is importable):
    .venv/bin/python .claude/skills/discord-harness-operate/scripts/preflight.py
Exit 0 if ready to start, 1 if a blocker is found.
"""

from __future__ import annotations

import os
import sys
from pathlib import Path

# models.work prefix -> required .env key (None = local model, no key needed). This harness wires
# exactly two routes: cloud Gemini and local Qwen via Ollama.
PROVIDER_KEY = {
    "gemini/": "GEMINI_API_KEY",
    "ollama_chat/": None,
    "ollama/": None,
}
SHOWABLE = {  # only these NON-secret key NAMES are echoed (never values)
    "DISCORD_TOKEN",
    "OWNER_ID",
    "CHANNEL_ID",
    "ROLE_ID",
    "GEMINI_API_KEY",
    "HARNESS_CONFIG",
    "AUDIT_PATH",
    "HEALTH_PORT",
}


def find_repo_root() -> Path:
    for d in [
        Path.cwd(),
        *Path.cwd().parents,
        Path(__file__).resolve().parent,
        *Path(__file__).resolve().parents,
    ]:
        if (d / "config.yaml").is_file() and (d / "bot.py").is_file():
            return d
    print(
        "error: could not locate the harness repo root (needs config.yaml + bot.py)",
        file=sys.stderr,
    )
    raise SystemExit(2)


def dotenv_keys(path: Path) -> set[str]:
    """Names with a non-empty value in .env. Values are tested for emptiness only, not stored."""
    present: set[str] = set()
    if not path.is_file():
        return present
    for line in path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        if v.strip():
            present.add(k.strip())
    return present


def main() -> int:
    root = find_repo_root()
    try:
        import yaml
    except Exception:  # noqa: BLE001
        print("error: PyYAML not available — run inside the repo's .venv", file=sys.stderr)
        return 2

    profile = os.environ.get("HARNESS_CONFIG") or "config.yaml"
    cfg_path = root / profile
    try:
        cfg = yaml.safe_load(cfg_path.read_text(encoding="utf-8")) or {}
    except Exception as exc:  # noqa: BLE001
        print(f"FAIL: cannot parse {profile}: {exc}", file=sys.stderr)
        return 1

    env_present = dotenv_keys(root / ".env") | {k for k in os.environ if os.environ.get(k)}
    problems: list[str] = []
    warnings: list[str] = []

    if "DISCORD_TOKEN" not in env_present:
        problems.append("DISCORD_TOKEN missing (.env) — bot.py exits without it")

    work = (cfg.get("models") or {}).get("work") or ""
    need = next((PROVIDER_KEY[p] for p in PROVIDER_KEY if work.startswith(p)), "UNKNOWN")
    if need == "UNKNOWN":
        warnings.append(
            f"models.work={work!r}: unrecognized provider prefix; check its key yourself"
        )
    elif need is None:
        warnings.append(
            f"models.work={work!r} is local (Ollama) — ensure `ollama serve` is running"
        )
    elif need not in env_present:
        problems.append(f"models.work={work!r} needs {need} in .env")

    users = (cfg.get("allowlist") or {}).get("users") or []
    if not users and "OWNER_ID" not in env_present:
        problems.append(
            "authorization locked: allowlist.users empty AND no OWNER_ID — nobody can use the bot"
        )

    if not cfg.get("tools"):
        warnings.append("no tools: classified — every tool fails closed to DENY")

    print(f"profile: {profile}   work model: {work or '(unset)'}")
    print(f"env keys present: {', '.join(sorted(env_present & SHOWABLE)) or '(none)'}")
    for w in warnings:
        print(f"  warn:  {w}")
    for p in problems:
        print(f"  BLOCK: {p}")
    if problems:
        print("FAIL: not ready to start — resolve the BLOCK item(s) above.", file=sys.stderr)
        return 1
    print("OK: pre-flight passed — ready to start.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
