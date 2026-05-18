#!/usr/bin/env python3
"""Deterministic static-credential scanner for agent / CI/CD codebases.

Flags three Layer-1 / Layer-3 leakage patterns over a target path:

  (a) likely hard-coded secret literals (AWS AKIA keys, private-key headers,
      high-entropy/keyword-named literal assignments);
  (b) debug statements that print/log a secret-named variable;
  (c) secrets in env-var assignments inside agent / MCP config.

Stdlib only. Output is one `file:line: [code] message` finding per line,
sorted deterministically. Exit codes mirror the skillify validators:
  0 = clean, 1 = findings, 2 = usage / IO error.
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

SECRET_NAME = r"(?:secret|client_secret|token|api[_-]?key|apikey|password|passwd|private[_-]?key)"

# (a) hard-coded secret literals
AWS_AKIA_RE = re.compile(r"\bAKIA[0-9A-Z]{16}\b")
PRIVATE_KEY_RE = re.compile(r"-----BEGIN (?:RSA |EC |OPENSSH |DSA |PGP )?PRIVATE KEY-----")
HARDCODED_ASSIGN_RE = re.compile(
    r"(?i)\b" + SECRET_NAME + r"\b\s*[:=]\s*['\"][^'\"\n]{8,}['\"]"
)

# (b) debug statements emitting a secret-named variable
PY_LOG_RE = re.compile(
    r"(?i)\b(?:print|logging\.\w+|logger\.\w+|log\.\w+)\s*\([^)]*" + SECRET_NAME
)
JS_LOG_RE = re.compile(
    r"(?i)\bconsole\.(?:log|error|warn|info|debug)\s*\([^)]*" + SECRET_NAME
)

# (c) secret in an env-var assignment (shell / Dockerfile / compose / CI)
ENV_ASSIGN_RE = re.compile(
    r"(?i)(?:^|export\s+|ENV\s+)[A-Z0-9_]*" + SECRET_NAME.replace("[_-]?", "_?")
    + r"[A-Z0-9_]*\s*[:=]\s*['\"]?[^\s'\"]{6,}"
)

CONFIG_SUFFIXES = {".json", ".yaml", ".yml", ".toml", ".env", ".ini", ".cfg"}
CONFIG_HINTS = ("mcp", "agent", "compose", "dockerfile", ".env")
SCAN_SUFFIXES = {
    ".py", ".js", ".ts", ".tsx", ".jsx", ".sh", ".bash", ".rb", ".go",
    ".java", ".rs", ".json", ".yaml", ".yml", ".toml", ".env", ".ini",
    ".cfg", ".txt", "",
}
SKIP_DIRS = {".git", "node_modules", "__pycache__", ".venv", "venv", "dist", "build"}
PLACEHOLDER_RE = re.compile(
    r"(?i)(your[_-]?|example|placeholder|changeme|xxxx+|\.\.\.|<[^>]+>|\$\{?[A-Z_]+\}?)"
)


def is_config_like(path: Path) -> bool:
    name = path.name.lower()
    if path.suffix.lower() in CONFIG_SUFFIXES:
        return True
    return any(hint in name for hint in CONFIG_HINTS)


def scan_file(path: Path) -> list[tuple[int, str, str]]:
    findings: list[tuple[int, str, str]] = []
    try:
        text = path.read_text(encoding="utf-8", errors="replace")
    except OSError:
        return findings
    config_like = is_config_like(path)
    for lineno, line in enumerate(text.splitlines(), start=1):
        if AWS_AKIA_RE.search(line):
            findings.append((lineno, "A", "hard-coded AWS access key id (AKIA...)"))
        if PRIVATE_KEY_RE.search(line):
            findings.append((lineno, "A", "hard-coded private key header"))
        m = HARDCODED_ASSIGN_RE.search(line)
        if m and not PLACEHOLDER_RE.search(m.group(0)):
            findings.append((lineno, "A", "hard-coded secret literal assignment"))
        if PY_LOG_RE.search(line) or JS_LOG_RE.search(line):
            findings.append((lineno, "B", "debug statement logs a secret-named variable"))
        if config_like and ENV_ASSIGN_RE.search(line):
            seg = ENV_ASSIGN_RE.search(line)
            if seg and not PLACEHOLDER_RE.search(seg.group(0)):
                findings.append(
                    (lineno, "C", "secret in env-var assignment in agent/MCP config")
                )
    return findings


def iter_files(root: Path):
    if root.is_file():
        yield root
        return
    for path in sorted(root.rglob("*")):
        if path.is_dir():
            continue
        if any(part in SKIP_DIRS for part in path.parts):
            continue
        if path.suffix.lower() in SCAN_SUFFIXES or path.suffix == "":
            yield path


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Scan a path for static credentials, secret-logging, and config secrets."
    )
    parser.add_argument("--root", default=".", help="File or directory to scan")
    args = parser.parse_args()

    root = Path(args.root).expanduser().resolve()
    if not root.exists():
        print(f"error: path does not exist: {root}", file=sys.stderr)
        return 2

    all_findings: list[tuple[str, int, str, str]] = []
    for path in iter_files(root):
        for lineno, code, message in scan_file(path):
            all_findings.append((str(path), lineno, code, message))

    if not all_findings:
        print(f"ok: no static credentials found under {root}")
        return 0

    for file_str, lineno, code, message in sorted(all_findings):
        print(f"{file_str}:{lineno}: [{code}] {message}")
    print(f"error: {len(all_findings)} finding(s)", file=sys.stderr)
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
