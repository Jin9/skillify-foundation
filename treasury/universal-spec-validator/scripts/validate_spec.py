#!/usr/bin/env python3
"""universal-spec-validator gate.

Deterministic, offline CI/pre-commit gate for agent specs. Applies three rule
sets (see ../references): portability (P*), command-safety (C*),
schema-evolution (E*, baseline-required). Writes spec-validation.json and
spec-validation.md. Exit: 0 = no blocking findings, 1 = blocking finding(s),
2 = usage/parse error.

No network, no command execution, no LLM. stdlib only (YAML config optional).
"""
import argparse, glob, json, re, sys
from pathlib import Path

SEV_GATE = {"critical": "block", "high": "block", "medium": "warn",
            "low": "warn", "info": "info"}
DESTRUCTIVE = [
    r"\brm\s+-rf\b", r":\(\)\s*\{.*\|\s*:.*\};:", r"\bcurl\b[^\n|]*\|\s*(sh|bash)",
    r"\bwget\b[^\n|]*\|\s*(sh|bash)", r"\beval\s*\(", r"\bsudo\b",
    r"\bchmod\s+777\b", r"\bgit\s+push\b[^\n]*(--force|-f)\b",
    r"\bDROP\s+TABLE\b", r"\bTRUNCATE\b", r'"shell"\s*:\s*true',
]
BROAD_SCOPE = re.compile(r'("(permissions|scope|scopes|role)"\s*:\s*'
                         r'("?\*"?|"all"|"admin"|"service-role"|"root"))',
                         re.I)
INJECTION_SRC = re.compile(r"(pr[_ ]?title|issue[_ ]?body|readme|"
                           r"support[_ ]?ticket|tool[_ ]?description|"
                           r"user[_ ]?input|comment)", re.I)
IRREVERSIBLE = re.compile(r"(deploy|publish|push|delete|drop|terraform\s+apply|"
                          r"prod(uction)?|iam|grant|revoke)", re.I)
FLOATING_VER = re.compile(r'(:latest\b|"\^|"~|@latest\b|"version"\s*:\s*"\*")')
VENDOR_KEYS = {"openai": '"type": "function"', "anthropic": '"input_schema"',
               "gemini": '"types.Schema"'}


def load_cfg(p):
    if not p:
        return {}
    txt = Path(p).read_text(encoding="utf-8")
    try:
        import yaml  # optional
        return yaml.safe_load(txt) or {}
    except Exception:
        try:
            return json.loads(txt)
        except Exception:
            return {}


def read_spec(p):
    raw = Path(p).read_text(encoding="utf-8")
    obj = None
    try:
        obj = json.loads(raw)
    except Exception:
        try:
            import yaml
            obj = yaml.safe_load(raw)
        except Exception:
            obj = None
    return raw, obj


def add(findings, axis, rid, sev, spec, locator, msg, ref):
    findings.append({"axis": axis, "id": rid, "severity": sev,
                     "gate": SEV_GATE[sev], "spec": spec, "locator": locator,
                     "message": msg, "rule_ref": ref})


def check_command_safety(raw, spec, F):
    for pat in DESTRUCTIVE:
        m = re.search(pat, raw, re.I)
        if m:
            add(F, "command-safety", "C1", "critical", spec,
                f"offset {m.start()}",
                f"destructive/unsafe command pattern: {m.group(0)[:60]!r}",
                "command-safety-rules.md#C1")
    m = BROAD_SCOPE.search(raw)
    if m:
        add(F, "command-safety", "C2", "high", spec, "permission scope",
            f"over-broad permission scope: {m.group(0)[:60]!r}; require an allowlist",
            "command-safety-rules.md#C2")
    if INJECTION_SRC.search(raw) and re.search(r"(command|cmd|exec|query|path)", raw, re.I):
        add(F, "command-safety", "C3", "high", spec, "untrusted field",
            "untrusted-origin field referenced near a command/query/path; declare sanitization",
            "command-safety-rules.md#C3")
    if IRREVERSIBLE.search(raw) and not re.search(
            r"(hitl|human[_ ]?in[_ ]?the[_ ]?loop|approval|requires_approval|confirm)", raw, re.I):
        add(F, "command-safety", "C4", "high", spec, "irreversible action",
            "irreversible/high-blast action with no HITL/approval annotation",
            "command-safety-rules.md#C4")
    m = FLOATING_VER.search(raw)
    if m:
        add(F, "command-safety", "C5", "medium", spec, "version ref",
            "floating/unpinned version reference; require exact pins",
            "command-safety-rules.md#C5")
    if IRREVERSIBLE.search(raw) and not re.search(
            r"(identity|service[_ ]?account|audit|action[_ ]?log)", raw, re.I):
        add(F, "command-safety", "C6", "medium", spec, "auditability",
            "autonomous surface without distinct identity / audit log",
            "command-safety-rules.md#C6")
    if re.search(r"(api[_-]?key|secret|password|token)", raw, re.I) and not re.search(
            r"(oidc|short[- ]?lived|sts|expir)", raw, re.I):
        add(F, "command-safety", "C7", "medium", spec, "credentials",
            "long-lived credential assumption; prefer short-lived OIDC",
            "command-safety-rules.md#C7")


def check_portability(raw, obj, spec, F, cfg):
    present = [v for v, key in VENDOR_KEYS.items() if key.replace(" ", "") in raw.replace(" ", "")]
    if len(present) == 1 and not re.search(r"(gateway|adapter|x-portable|normaliz)", raw, re.I):
        add(F, "portability", "P1", "high", spec, "tool-call dialect",
            f"single-vendor tool-call dialect ({present[0]}) with no normalization layer",
            "portability-rules.md#P1")
    if isinstance(obj, dict):
        def scan(node, path="$"):
            if isinstance(node, dict):
                if node.get("type") == "object" and "properties" in node:
                    if node.get("additionalProperties", None) in (None, True):
                        add(F, "portability", "P2", "medium", spec, path,
                            "object input without additionalProperties:false (vendor strict-mode drift)",
                            "portability-rules.md#P2")
                for k, v in node.items():
                    scan(v, f"{path}.{k}")
            elif isinstance(node, list):
                for i, v in enumerate(node):
                    scan(v, f"{path}[{i}]")
        scan(obj)
    if re.search(r"(embedding|prompt[_ ]?template|system[_ ]?prompt)", raw, re.I) and not re.search(
            r"(requires_capabilities|min_context_tokens|capability)", raw, re.I):
        add(F, "portability", "P4", "medium", spec, "artifact dependency",
            "model-specific prompt/embedding dependency without capability annotation",
            "portability-rules.md#P4")
    if re.search(r"(code[_ ]?execution|web[_ ]?access|file[_ ]?write|network)", raw, re.I) and not re.search(
            r"requires_capabilities", raw, re.I):
        add(F, "portability", "P5", "medium", spec, "capability gate",
            "implies privileged capability but declares no requires_capabilities for routing",
            "portability-rules.md#P5")


def _props(o):
    return o.get("properties", {}) if isinstance(o, dict) else {}


def _required(o):
    return set(o.get("required", []) or []) if isinstance(o, dict) else set()


def check_evolution(base_obj, cur_obj, spec, F, mode):
    if base_obj is None:
        add(F, "schema-evolution", "E0", "info", spec, "-",
            "schema-evolution checks skipped (no baseline)",
            "schema-evolution-rules.md")
        return
    bp, cp = _props(base_obj), _props(cur_obj)
    br, cr = _required(base_obj), _required(cur_obj)
    for name in bp:
        if name not in cp:
            add(F, "schema-evolution", "E1", "high", spec, f"$.{name}",
                f"field '{name}' removed (breaking under {mode})",
                "schema-evolution-rules.md")
        else:
            bt, ct = bp[name].get("type"), cp[name].get("type")
            if bt and ct and bt != ct:
                add(F, "schema-evolution", "E2", "high", spec, f"$.{name}",
                    f"field '{name}' type changed {bt}->{ct} (breaking)",
                    "schema-evolution-rules.md")
            be = bp[name].get("enum")
            ce = cp[name].get("enum")
            if be and ce and set(ce) < set(be):
                add(F, "schema-evolution", "E4", "high", spec, f"$.{name}",
                    f"enum narrowed for '{name}' (breaking)",
                    "schema-evolution-rules.md")
    for name in cr - br:
        if name not in bp and not (isinstance(cp.get(name), dict) and "default" in cp.get(name, {})):
            add(F, "schema-evolution", "E3", "high", spec, f"$.{name}",
                f"new required field '{name}' without default (breaking)",
                "schema-evolution-rules.md")
    for name in (cr & br.union(set(bp))) - br:
        add(F, "schema-evolution", "E5", "high", spec, f"$.{name}",
            f"optional field '{name}' became required (breaking)",
            "schema-evolution-rules.md")


def main():
    ap = argparse.ArgumentParser(description="universal-spec-validator CI gate")
    ap.add_argument("specs", nargs="+", help="spec file(s) or glob(s)")
    ap.add_argument("--baseline", help="previously-shipped spec file (enables evolution checks)")
    ap.add_argument("--config", help=".spec-validator.yaml")
    ap.add_argument("--out", default=".", help="output dir (default .)")
    ap.add_argument("--format", default="both", choices=["json", "md", "both"])
    a = ap.parse_args()

    cfg = load_cfg(a.config)
    mode = cfg.get("compatibility_mode", "backward")
    overrides = cfg.get("gate_overrides", {}) or {}
    fail_on = set(cfg.get("fail_on", ["critical", "high"]))
    ignores = cfg.get("ignore", []) or []

    paths = []
    for s in a.specs:
        paths += [p for p in glob.glob(s, recursive=True) if Path(p).is_file()]
    if not paths:
        print("error: no spec files matched", file=sys.stderr)
        sys.exit(2)

    base_raw, base_obj = (read_spec(a.baseline) if a.baseline else (None, None))
    F = []
    for p in sorted(set(paths)):
        try:
            raw, obj = read_spec(p)
        except Exception as e:
            print(f"error: cannot read {p}: {e}", file=sys.stderr)
            sys.exit(2)
        check_command_safety(raw, p, F)
        check_portability(raw, obj, p, F, cfg)
        if a.baseline:
            check_evolution(base_obj, obj, p, F, mode)

    # apply ignore + gate_overrides + fail_on
    kept = []
    for f in F:
        if any(f["id"] == ig.get("rule") and
               glob.fnmatch.fnmatch(f["spec"], ig.get("locator", "*"))
               for ig in ignores):
            continue
        if f["id"] in overrides:
            f["gate"] = overrides[f["id"]]
        elif f["severity"] in fail_on and f["gate"] != "info":
            f["gate"] = "block"
        kept.append(f)

    block = sum(1 for f in kept if f["gate"] == "block")
    warn = sum(1 for f in kept if f["gate"] == "warn")
    info = sum(1 for f in kept if f["gate"] == "info")
    by_axis = {ax: sum(1 for f in kept if f["axis"] == ax)
               for ax in ("portability", "command-safety", "schema-evolution")}
    exit_code = 1 if block else 0
    report = {"summary": {"specs": len(set(paths)), "block": block,
                          "warn": warn, "info": info, "by_axis": by_axis},
              "exit_code": exit_code, "findings": kept}

    out = Path(a.out)
    out.mkdir(parents=True, exist_ok=True)
    if a.format in ("json", "both"):
        (out / "spec-validation.json").write_text(
            json.dumps(report, indent=2), encoding="utf-8")
    if a.format in ("md", "both"):
        lines = ["# Spec validation report", "",
                 f"- specs: {len(set(paths))}  block: {block}  "
                 f"warn: {warn}  info: {info}",
                 f"- exit_code: {exit_code}  (0 pass / 1 blocked / 2 error)",
                 f"- by axis: {by_axis}", ""]
        for sev in ("critical", "high", "medium", "low", "info"):
            rows = [f for f in kept if f["severity"] == sev]
            if not rows:
                continue
            lines.append(f"## {sev.upper()}")
            for f in rows:
                lines.append(f"- [{f['gate']}] **{f['axis']}/{f['id']}** "
                             f"`{f['spec']}` @ `{f['locator']}` — {f['message']} "
                             f"({f['rule_ref']})")
            lines.append("")
        (out / "spec-validation.md").write_text("\n".join(lines), encoding="utf-8")

    print(f"universal-spec-validator: {block} blocking, {warn} warn, "
          f"{info} info across {len(set(paths))} spec(s) -> exit {exit_code}")
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
