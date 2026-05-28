#!/usr/bin/env python3
"""
validate_dev_output.py — structural-conformance gate between Dev emit and QA-L1
dispatch.

Hard-rejects a Dev artifact (per backend service or the frontend) if any of the
six backend rules or four frontend rules fail. Failures route as `code_mismatch`
(cap 2) per docs/orchestrator.md § Routing Table.

Backend service rules (B2C E-Commerce Platform/backend/services/<svc>/):
  1. go.sum present (and optionally `go mod verify` clean if --run-tools).
  2. _test.go sibling exists for every FULL handler_*.go in app/<domain>/.
     A handler is FULL unless its body contains the literal "NOT_IMPLEMENTED_MVP"
     (the squad's documented 501 stub marker). Handlers with no body match are
     considered FULL.
  3. (advisory unless --run-tools) make lint clean.
  4. Layout matches backend/go-template: app/<domain>/, access/, router/,
     config/, migrations/ all exist.
  5. go.mod's `replace` directive for the common module points to ../../common
     (the dry-run #1 footgun was pointing at gitlab.com/.../common.git).
  6. Response envelope code literal in code is "SUCCESS" (or another contract
     literal documented in cross-cutting.response-envelope), NOT "0000".

Frontend rules (B2C E-Commerce Platform/frontend/):
  F1. No occurrences of `dangerouslySetInnerHTML` anywhere under src/.
  F2. Cookies set or cleared via response headers / next/headers cookies API
      include `HttpOnly`, `SameSite`, and `Secure` attributes (best-effort grep).
  F3. Idempotency-Key uses `crypto.randomUUID()` for the `checkout.commit` and
      `payment.simulate` flows.
  F4. (advisory) Route paths in api proxy handlers match TD's per-API spec MDs
      (full match not enforceable without parsing every TD; this rule lints
      for the few known mismatches from KI-13).

Usage:
  python3 validate_dev_output.py <backend_services_dir> [--frontend <frontend_dir>]
                                 [--run-tools]

Exit codes:
  0 = all hard rules pass.
  2 = at least one hard-reject violation.
  64 = usage error.
"""
from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path

# Top-level dirs required at the service root per docs/templates/dev.md
# § "Backend service mandatory structure". `access/` is nested under
# `app/<domain>/access/` (not at root) per the template tree.
EXPECTED_LAYOUT = ["app", "router", "config", "migrations"]
EXPECTED_DOMAIN_SUBDIR = "access"  # required under each app/<domain>/
STUB_MARKERS = (
    "NOT_IMPLEMENTED_MVP",
    "NOT_IMPLEMENTED",  # legacy
)
HANDLER_GLOB = "handler_*.go"

# Anti-patterns that should never appear in response envelope code.
BAD_ENVELOPE_LITERALS = (
    '"0000"',
    "'0000'",
)
GOOD_ENVELOPE_LITERAL = '"SUCCESS"'

FRONTEND_BAD_PATTERNS = [
    ("dangerouslySetInnerHTML", "F1: dangerouslySetInnerHTML present"),
]
FRONTEND_GOOD_REQUIRED = [
    ("crypto.randomUUID(", "F3: crypto.randomUUID() not used for Idempotency-Key"),
]
COOKIE_REQUIRED_ATTRS = ("HttpOnly", "SameSite", "Secure")

# KI-13 known route mismatches we hard-reject to prove the gate works.
KI13_BAD_ROUTES = [
    ("/api/proxy/order/list", "F4: KI-13 — should be /api/proxy/order/list-mine"),
]


def err(violations: list[str], scope: str, msg: str) -> None:
    violations.append(f"[{scope}] {msg}")


def has_stub_marker(text: str) -> bool:
    return any(m in text for m in STUB_MARKERS)


def find_handler_pairs(domain_dir: Path) -> list[tuple[Path, Path]]:
    """Return (handler_path, expected_test_path) for every handler_*.go."""
    pairs: list[tuple[Path, Path]] = []
    for handler in sorted(domain_dir.glob(HANDLER_GLOB)):
        if handler.name.endswith("_test.go"):
            continue
        test_path = handler.with_name(handler.stem + "_test.go")
        pairs.append((handler, test_path))
    return pairs


def validate_backend_service(svc_dir: Path,
                             violations: list[str],
                             run_tools: bool) -> None:
    svc = svc_dir.name

    # Rule 4 — layout match
    for sub in EXPECTED_LAYOUT:
        if not (svc_dir / sub).exists():
            err(violations, svc, f"Rule 4: missing layout dir {sub!r}")

    # Rule 1 — go.sum present
    go_sum = svc_dir / "go.sum"
    if not go_sum.exists():
        err(violations, svc, "Rule 1: go.sum is missing")

    # Rule 5 — go.mod replace directive
    go_mod = svc_dir / "go.mod"
    if go_mod.exists():
        gm = go_mod.read_text()
        if "replace" in gm and "common" in gm:
            # Look for the common-replace line specifically
            replace_lines = [
                ln.strip() for ln in gm.splitlines()
                if ln.strip().startswith("replace ") and "common" in ln
            ]
            for ln in replace_lines:
                if "../../common" not in ln:
                    err(
                        violations,
                        svc,
                        f"Rule 5: go.mod replace directive for common does not "
                        f"point to ../../common: {ln!r}",
                    )
        else:
            # Some services may not need a replace if module is local-relative;
            # warn (not fail) to avoid false positives on simple services.
            pass
    else:
        err(violations, svc, "Rule 5: go.mod missing")

    # Rule 2 — _test.go siblings for FULL handlers in app/<domain>/
    # Rule 4b — each app/<domain>/ must contain an access/ subdir
    app_dir = svc_dir / "app"
    if app_dir.is_dir():
        domain_dirs = [d for d in app_dir.iterdir() if d.is_dir()]
        for domain_dir in domain_dirs:
            if not (domain_dir / EXPECTED_DOMAIN_SUBDIR).is_dir():
                err(
                    violations,
                    svc,
                    f"Rule 4b: missing {EXPECTED_DOMAIN_SUBDIR}/ subdir under "
                    f"app/{domain_dir.name}/",
                )
            for handler, test_path in find_handler_pairs(domain_dir):
                body = handler.read_text(errors="replace")
                if has_stub_marker(body):
                    continue  # 501 stub; tests not required
                if not test_path.exists():
                    err(
                        violations,
                        svc,
                        f"Rule 2: missing test sibling for FULL handler "
                        f"{handler.relative_to(svc_dir)} (expected "
                        f"{test_path.relative_to(svc_dir)})",
                    )

    # Rule 6 — response envelope code literal
    bad_hits: list[str] = []
    for go_file in svc_dir.rglob("*.go"):
        if go_file.name.endswith("_test.go"):
            continue
        try:
            text = go_file.read_text(errors="replace")
        except Exception:
            continue
        for bad in BAD_ENVELOPE_LITERALS:
            if bad in text:
                # Heuristic: only flag if it appears near the word "code" or "Code"
                # to reduce false positives on opaque "0000" strings used for ids.
                for match in re.finditer(re.escape(bad), text):
                    window = text[max(0, match.start() - 60):match.end() + 60]
                    if re.search(r"\b[Cc]ode\b", window):
                        bad_hits.append(
                            f"{go_file.relative_to(svc_dir)}: envelope code "
                            f"literal {bad} (expected {GOOD_ENVELOPE_LITERAL})"
                        )
                        break  # one hit per file is enough
    for hit in bad_hits:
        err(violations, svc, f"Rule 6: {hit}")

    # Rules 1b + 3 — opt-in tool runs
    if run_tools:
        try:
            r = subprocess.run(
                ["go", "mod", "verify"],
                cwd=svc_dir, capture_output=True, text=True, timeout=60,
            )
            if r.returncode != 0:
                err(
                    violations,
                    svc,
                    f"Rule 1b: go mod verify failed: {r.stderr.strip() or r.stdout.strip()}",
                )
        except FileNotFoundError:
            err(violations, svc, "Rule 1b: `go` toolchain not found in PATH")
        except subprocess.TimeoutExpired:
            err(violations, svc, "Rule 1b: go mod verify timed out (60s)")

        # `make lint` is service-local; only invoke if a Makefile is present.
        if (svc_dir / "Makefile").exists():
            try:
                r = subprocess.run(
                    ["make", "lint"],
                    cwd=svc_dir, capture_output=True, text=True, timeout=180,
                )
                if r.returncode != 0:
                    err(
                        violations,
                        svc,
                        f"Rule 3: make lint failed (see stderr).",
                    )
            except FileNotFoundError:
                err(violations, svc, "Rule 3: `make` not found in PATH")
            except subprocess.TimeoutExpired:
                err(violations, svc, "Rule 3: make lint timed out (180s)")


def validate_frontend(frontend_dir: Path, violations: list[str]) -> None:
    if not frontend_dir.is_dir():
        return  # nothing to validate
    scope = frontend_dir.name

    src_files = list(frontend_dir.rglob("*.ts")) \
              + list(frontend_dir.rglob("*.tsx")) \
              + list(frontend_dir.rglob("*.js")) \
              + list(frontend_dir.rglob("*.jsx"))

    src_files = [
        f for f in src_files
        if "/node_modules/" not in str(f) and "/.next/" not in str(f)
    ]

    # F1 — no dangerouslySetInnerHTML
    for f in src_files:
        try:
            text = f.read_text(errors="replace")
        except Exception:
            continue
        for bad, label in FRONTEND_BAD_PATTERNS:
            if bad in text:
                err(violations, scope, f"{label} in {f.relative_to(frontend_dir)}")

        # F4 — KI-13 known route mismatches
        for bad_route, label in KI13_BAD_ROUTES:
            if bad_route in text:
                err(violations, scope,
                    f"{label} in {f.relative_to(frontend_dir)}")

    # F2 — cookie attributes (best-effort: any line setting or clearing a
    # session-style cookie should mention HttpOnly + SameSite + Secure)
    cookie_callsites: list[tuple[Path, str]] = []
    cookie_re = re.compile(
        r"(setAccessCookie|setRefreshCookie|clearAccessCookie|clearRefreshCookie|"
        r"cookies\(\)\.set|cookies\(\)\.delete)"
    )
    for f in src_files:
        try:
            text = f.read_text(errors="replace")
        except Exception:
            continue
        for m in cookie_re.finditer(text):
            window = text[max(0, m.start() - 200):m.end() + 600]
            cookie_callsites.append((f, window))
    for f, window in cookie_callsites:
        missing = [a for a in COOKIE_REQUIRED_ATTRS if a not in window]
        if missing:
            err(
                violations,
                scope,
                f"F2: cookie call in {f.relative_to(frontend_dir)} missing "
                f"attrs {missing}",
            )

    # F3 — Idempotency-Key uses crypto.randomUUID()
    saw_idempotency = any(
        "Idempotency-Key" in (f.read_text(errors="replace") if f.exists() else "")
        for f in src_files
    )
    if saw_idempotency:
        any_uuid = any(
            "crypto.randomUUID(" in f.read_text(errors="replace")
            for f in src_files if f.exists()
        )
        if not any_uuid:
            err(
                violations,
                scope,
                "F3: Idempotency-Key referenced but crypto.randomUUID() not used",
            )


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("backend_services_dir",
                   help="path to backend/services/ (one folder per service)")
    p.add_argument("--frontend", default=None,
                   help="path to frontend/ root (optional)")
    p.add_argument("--run-tools", action="store_true",
                   help="actually invoke `go mod verify` and `make lint` "
                        "(default: skip; static checks only)")
    args = p.parse_args()

    backend_dir = Path(args.backend_services_dir).resolve()
    if not backend_dir.is_dir():
        print(f"validate_dev_output: not a directory: {backend_dir}", file=sys.stderr)
        return 64

    violations: list[str] = []
    services = sorted(d for d in backend_dir.iterdir() if d.is_dir())
    if not services:
        print(f"validate_dev_output: no services found under {backend_dir}",
              file=sys.stderr)
        return 2

    for svc_dir in services:
        validate_backend_service(svc_dir, violations, args.run_tools)

    if args.frontend:
        validate_frontend(Path(args.frontend).resolve(), violations)

    if violations:
        print("validate_dev_output: FAIL", file=sys.stderr)
        for v in violations:
            print(f"  {v}", file=sys.stderr)
        return 2

    print(
        f"validate_dev_output: PASS "
        f"({len(services)} backend services{' + frontend' if args.frontend else ''})"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
