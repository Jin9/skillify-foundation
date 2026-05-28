#!/usr/bin/env python3
"""
package_skill.py — deterministic skill-bundle packager (replaces missing skillify skill).

Emits the skillify-format bundle described in docs/templates/skill-bundle.md when
the orchestrator reaches terminal state Done or ShipWithCaveats. Exits 0 without
output for Failed / HardFail.

Inputs (read-only):
  - <workflow_root>/.squad-run/state.json
  - <workflow_root>/.squad-run/ba.json
  - <workflow_root>/.squad-run/{ADVISORY,PLAN_NOTES,KNOWN_ISSUES}.md (each optional)
  - <workflow_root>/.squad-run/components/<svc>/{qa-l1.json,qa-l1.csv}
  - <workflow_root>/.squad-run/layer-2/{qa-l2.json,qa-l2.csv}
  - <workflow_root>/requirement/<EPIC>/<STORY>.md (mirrored)
  - <workflow_root>/design/{architecture,components,spec}/ (mirrored)
  - <repo_root>/docs/templates/*.md (copied verbatim into bundle templates/ dir)

Outputs (written under <workflow_root>/):
  - SKILL.md
  - references/{ba-acceptance-criteria.md,architecture/,requirement/,components/,spec/}
  - templates/ (copies of docs/templates/*.md)
  - tests/{qa-l1-<svc>.{json,csv},qa-l2.{json,csv}}
  - ADVISORY.md, PLAN_NOTES.md, KNOWN_ISSUES.md (if exist in .squad-run/)

Usage:
  python3 package_skill.py --workflow <workflow_root> [--repo-root <repo_root>]
                          [--state <state.json>] [--out <bundle_root>]

Exit codes:
  0 = bundle emitted (or no-op for Failed/HardFail).
  2 = required input missing or invalid.
  64 = usage error.
"""
from __future__ import annotations

import argparse
import datetime as _dt
import json
import shutil
import sys
from pathlib import Path

GENERATOR_TAG = "package_skill.py@v0.1.0"
TEMPLATE_VERSION = "0.1.0"

TERMINAL_BUNDLE = {"Done", "ShipWithCaveats"}
TERMINAL_NOOP = {"Failed", "HardFail"}


def fail(msg: str, code: int = 2) -> None:
    print(f"package_skill: FAIL {msg}", file=sys.stderr)
    sys.exit(code)


def short(s: str, n: int = 120) -> str:
    s = " ".join(s.strip().split())
    return s if len(s) <= n else s[: n - 1].rstrip() + "…"


def slugify(s: str) -> str:
    out = []
    last_dash = False
    for ch in s.lower():
        if ch.isalnum():
            out.append(ch)
            last_dash = False
        elif not last_dash:
            out.append("-")
            last_dash = True
    return "".join(out).strip("-")


def mirror_tree(src: Path, dst: Path) -> int:
    """Copy src into dst; return file count. Skips on missing src."""
    if not src.exists():
        return 0
    if dst.exists():
        shutil.rmtree(dst)
    shutil.copytree(src, dst)
    return sum(1 for _ in dst.rglob("*") if _.is_file())


def copy_file(src: Path, dst: Path) -> bool:
    if not src.exists():
        return False
    dst.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(src, dst)
    return True


def order_epics_by_dependency(requirement_dir: Path) -> list[Path]:
    """Return EPIC MDs in dependency order (topological). Falls back to alpha."""
    epics: dict[str, dict] = {}
    for epic_md in sorted(requirement_dir.glob("EPIC_*/EPIC_*.md")):
        text = epic_md.read_text(errors="replace")
        epic_id = epic_md.stem
        deps: list[str] = []
        # Parse YAML frontmatter for dependencies — best-effort line scan.
        if text.startswith("---"):
            fm_end = text.find("---", 3)
            if fm_end > 0:
                fm = text[3:fm_end]
                in_deps = False
                for ln in fm.splitlines():
                    s = ln.strip()
                    if s.startswith("dependencies:"):
                        in_deps = True
                        rest = s[len("dependencies:"):].strip()
                        if rest.startswith("[") and rest.endswith("]"):
                            for tok in rest[1:-1].split(","):
                                tok = tok.strip().strip("'\"")
                                if tok:
                                    deps.append(tok)
                            in_deps = False
                        continue
                    if in_deps:
                        if s.startswith("-"):
                            deps.append(s[1:].strip().strip("'\""))
                        elif s and not s.startswith(" "):
                            in_deps = False
        epics[epic_id] = {"path": epic_md, "deps": deps}

    ordered: list[str] = []
    visiting: set[str] = set()
    visited: set[str] = set()

    def visit(eid: str) -> None:
        if eid in visited or eid not in epics:
            return
        if eid in visiting:
            return  # cycle: drop
        visiting.add(eid)
        for d in epics[eid]["deps"]:
            visit(d)
        visiting.discard(eid)
        visited.add(eid)
        ordered.append(eid)

    for eid in epics:
        visit(eid)
    return [epics[e]["path"] for e in ordered]


def write_skill_md(workflow_root: Path,
                   bundle_root: Path,
                   state: dict,
                   ba: dict,
                   epics_ordered: list[Path],
                   components: list[str],
                   has_known_issues: bool) -> Path:
    name = slugify(state.get("run_id", "skill")).rsplit("-", 3)[0] or "skill"
    description = short(ba.get("purpose", state.get("run_id", "")), 120)
    domain = state.get("domain", "b2c-retail")
    started_at = state.get("started_at", _dt.date.today().isoformat())
    terminal = state.get("terminal_state", "Done")
    run_id = state.get("run_id", "")

    when_to_use_lines = []
    for epic_md in epics_ordered:
        eid = epic_md.stem
        when_to_use_lines.append(f"  - {eid}")
    if not when_to_use_lines:
        when_to_use_lines = ["  - (no EPICs found — verify requirement/ tree)"]

    not_for_lines = []
    for item in ba.get("out_of_scope", []) or []:
        if isinstance(item, str):
            not_for_lines.append(f"  - {short(item, 200)}")
        elif isinstance(item, dict) and "item" in item:
            not_for_lines.append(f"  - {short(str(item['item']), 200)}")
    if not not_for_lines:
        not_for_lines = ["  - (none recorded)"]

    fm_lines = [
        "---",
        f"name: {name}",
        f"description: {description}",
        "when_to_use:",
        *when_to_use_lines,
        "not_for:",
        *not_for_lines,
        f"domain: {domain}",
        f"template_version: {TEMPLATE_VERSION}",
        f"generated_at: {started_at}",
        f"generator: {GENERATOR_TAG}",
        f"source_run_id: {run_id}",
        f"terminal_state: {terminal}",
        "---",
        "",
    ]

    workflow_name = ba.get("workflow_name") or state.get("run_id", "Workflow")

    body: list[str] = []
    body.append(f"# {workflow_name}")
    body.append("")

    body.append("## Inputs")
    body.append("")
    in_scope = ba.get("in_scope", []) or []
    if in_scope:
        for item in in_scope[:30]:
            label = item if isinstance(item, str) else (
                item.get("item") or item.get("name") or json.dumps(item, ensure_ascii=False)
            )
            body.append(f"- {short(str(label), 200)}")
        if len(in_scope) > 30:
            body.append(f"- (+{len(in_scope) - 30} more — see "
                        f"`references/ba-acceptance-criteria.md`)")
    else:
        body.append("- (no in-scope items recorded — see "
                    "`references/ba-acceptance-criteria.md`)")
    body.append("")

    body.append("## Steps")
    body.append("")
    if epics_ordered:
        for i, epic_md in enumerate(epics_ordered, 1):
            rel = Path("references") / "requirement" / epic_md.parent.name / epic_md.name
            body.append(f"{i}. [{epic_md.stem}]({rel.as_posix()})")
    else:
        body.append("1. (no EPICs found)")
    body.append("")

    body.append("## Acceptance Criteria")
    body.append("")
    body.append("See [`references/ba-acceptance-criteria.md`](references/ba-acceptance-criteria.md).")
    body.append("")

    body.append("## Architecture")
    body.append("")
    arch_files = [
        ("components.json", "Component list"),
        ("contracts.json", "Integration contracts"),
        ("infra-summary.md", "Infra summary"),
        ("infra-topology.md", "Infra topology"),
        ("connectivity.md", "Connectivity"),
        ("observability-spec.md", "Observability spec"),
    ]
    for fname, label in arch_files:
        rel = f"references/architecture/{fname}"
        body.append(f"- [{label}]({rel})")
    body.append("- [ADRs](references/architecture/ADRs/)")
    body.append("")

    body.append("## Per-Component Specs")
    body.append("")
    for svc in components:
        td_rel = f"references/components/{svc}/td.json"
        erd_rel = f"references/components/{svc}/erd.md"
        body.append(f"- **{svc}** — [td.json]({td_rel}) · [erd.md]({erd_rel})")
    body.append("")

    body.append("## API Reference")
    body.append("")
    body.append("See [`references/spec/`](references/spec/) for one MD per endpoint, "
                "grouped by `<service>-<aggregate>/`.")
    body.append("")

    body.append("## Tests")
    body.append("")
    for svc in components:
        body.append(f"- {svc}: [qa-l1 csv](tests/qa-l1-{svc}.csv) · "
                    f"[qa-l1 json](tests/qa-l1-{svc}.json)")
    body.append("- system: [qa-l2 csv](tests/qa-l2.csv) · [qa-l2 json](tests/qa-l2.json)")
    body.append("")

    if has_known_issues:
        body.append("## Known Issues")
        body.append("")
        body.append("See [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) for the full caveat list "
                    "(this run terminated as `ShipWithCaveats`).")
        body.append("")

    body.append("## Advisories")
    body.append("")
    body.append("See [`ADVISORY.md`](ADVISORY.md).")
    body.append("")

    body.append("## Provenance")
    body.append("")
    body.append("Generated by squad-flow workflow squad. See `run-log.md` and "
                "`run-report.md` under `.squad-run/runs/<run_id>/` for the full "
                "execution narrative.")
    body.append("")

    skill_md_path = bundle_root / "SKILL.md"
    skill_md_path.write_text("\n".join(fm_lines + body))
    return skill_md_path


def write_acceptance_md(bundle_root: Path, ba: dict) -> Path:
    out = bundle_root / "references" / "ba-acceptance-criteria.md"
    out.parent.mkdir(parents=True, exist_ok=True)
    lines: list[str] = ["# Acceptance Criteria",
                        "",
                        "Mirrored from `.squad-run/ba.json`.",
                        ""]

    purpose = ba.get("purpose")
    if purpose:
        lines.append("## Purpose")
        lines.append("")
        lines.append(str(purpose).strip())
        lines.append("")

    acs = ba.get("acceptance_criteria", []) or []
    lines.append(f"## Acceptance Criteria ({len(acs)})")
    lines.append("")
    for i, ac in enumerate(acs, 1):
        if isinstance(ac, str):
            lines.append(f"{i}. {ac}")
            continue
        ac_id = ac.get("id") or ac.get("ac_id") or f"AC-{i:03d}"
        given = ac.get("given") or ac.get("Given") or ""
        when = ac.get("when") or ac.get("When") or ""
        then = ac.get("then") or ac.get("Then") or ""
        lines.append(f"### {ac_id}")
        lines.append("")
        if given: lines.append(f"- **Given** {given}")
        if when:  lines.append(f"- **When** {when}")
        if then:  lines.append(f"- **Then** {then}")
        lines.append("")

    edges = ba.get("edge_cases", []) or []
    if edges:
        lines.append(f"## Edge cases ({len(edges)})")
        lines.append("")
        for e in edges:
            lines.append(f"- {e if isinstance(e, str) else json.dumps(e, ensure_ascii=False)}")
        lines.append("")

    succ = ba.get("success_conditions", []) or []
    if succ:
        lines.append("## Success conditions")
        lines.append("")
        for s in succ:
            lines.append(f"- {s if isinstance(s, str) else json.dumps(s, ensure_ascii=False)}")
        lines.append("")

    out.write_text("\n".join(lines))
    return out


def write_spec_index(bundle_root: Path, spec_root: Path) -> None:
    if not spec_root.is_dir():
        return
    idx = bundle_root / "references" / "spec" / "INDEX.md"
    idx.parent.mkdir(parents=True, exist_ok=True)
    lines = ["# API Reference Index", ""]
    for grp in sorted(p for p in spec_root.iterdir() if p.is_dir()):
        lines.append(f"## {grp.name}")
        lines.append("")
        for md in sorted(grp.glob("*.md")):
            lines.append(f"- [{md.stem}]({grp.name}/{md.name})")
        lines.append("")
    idx.write_text("\n".join(lines))


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--workflow", required=True, help="workflow root")
    p.add_argument("--repo-root", default=None,
                   help="repo root (default: parent of workflow root)")
    p.add_argument("--state", default=None,
                   help="state.json path (default: <workflow>/.squad-run/state.json)")
    p.add_argument("--out", default=None,
                   help="bundle output root (default: workflow root itself)")
    args = p.parse_args()

    workflow_root = Path(args.workflow).resolve()
    if not workflow_root.is_dir():
        fail(f"workflow root is not a directory: {workflow_root}")

    repo_root = Path(args.repo_root).resolve() if args.repo_root else workflow_root.parent
    state_path = Path(args.state).resolve() if args.state else (
        workflow_root / ".squad-run" / "state.json"
    )
    bundle_root = Path(args.out).resolve() if args.out else workflow_root

    if not state_path.exists():
        fail(f"state.json not found at {state_path}")

    try:
        state = json.loads(state_path.read_text())
    except json.JSONDecodeError as exc:
        fail(f"state.json invalid JSON: {exc}")

    terminal = state.get("terminal_state", "")
    if terminal in TERMINAL_NOOP:
        print(f"package_skill: NOOP terminal_state={terminal!r} (no bundle emitted)")
        return 0
    if terminal not in TERMINAL_BUNDLE:
        fail(
            f"terminal_state={terminal!r} not in {TERMINAL_BUNDLE}; "
            f"refusing to package mid-run state"
        )

    ba_path = workflow_root / ".squad-run" / "ba.json"
    if not ba_path.exists():
        fail(f"ba.json not found at {ba_path}")
    try:
        ba = json.loads(ba_path.read_text())
    except json.JSONDecodeError as exc:
        fail(f"ba.json invalid JSON: {exc}")

    components = list(state.get("components", []))
    bundle_root.mkdir(parents=True, exist_ok=True)

    # Mirror major design + requirement subtrees into references/.
    refs_root = bundle_root / "references"
    refs_root.mkdir(exist_ok=True)
    file_count = 0
    file_count += mirror_tree(workflow_root / "design" / "architecture",
                              refs_root / "architecture")
    file_count += mirror_tree(workflow_root / "design" / "components",
                              refs_root / "components")
    file_count += mirror_tree(workflow_root / "design" / "spec",
                              refs_root / "spec")
    file_count += mirror_tree(workflow_root / "requirement",
                              refs_root / "requirement")

    # Build per-aggregate spec index.
    write_spec_index(bundle_root, refs_root / "spec")

    # Mirror authored design templates into templates/ (so the bundle is itself
    # consumable by another squad instance).
    repo_templates_dir = repo_root / "docs" / "templates"
    if repo_templates_dir.is_dir():
        templates_dst = bundle_root / "templates"
        if templates_dst.exists():
            shutil.rmtree(templates_dst)
        templates_dst.mkdir(parents=True)
        for md in sorted(repo_templates_dir.glob("*.md")):
            shutil.copy2(md, templates_dst / md.name)
            file_count += 1

    # Mirror per-component QA + system QA into tests/.
    tests_dst = bundle_root / "tests"
    tests_dst.mkdir(exist_ok=True)
    for svc in components:
        svc_dir = workflow_root / ".squad-run" / "components" / svc
        for ext in ("json", "csv"):
            src = svc_dir / f"qa-l1.{ext}"
            if copy_file(src, tests_dst / f"qa-l1-{svc}.{ext}"):
                file_count += 1
    layer2 = workflow_root / ".squad-run" / "layer-2"
    for ext in ("json", "csv"):
        src = layer2 / f"qa-l2.{ext}"
        if copy_file(src, tests_dst / f"qa-l2.{ext}"):
            file_count += 1

    # Verbatim copies of advisory + plan-notes + known-issues.
    sr = workflow_root / ".squad-run"
    has_known_issues = False
    for fname in ("ADVISORY.md", "PLAN_NOTES.md", "KNOWN_ISSUES.md"):
        if copy_file(sr / fname, bundle_root / fname):
            if fname == "KNOWN_ISSUES.md":
                has_known_issues = True
            file_count += 1

    # Acceptance-criteria projection from ba.json.
    write_acceptance_md(bundle_root, ba)

    # Order EPICs by dependency for the SKILL.md ## Steps section.
    epics_ordered = order_epics_by_dependency(workflow_root / "requirement")

    skill_md = write_skill_md(
        workflow_root, bundle_root, state, ba,
        epics_ordered, components, has_known_issues,
    )

    print(
        f"package_skill: PASS bundle at {bundle_root} "
        f"(SKILL.md: {skill_md.name}; "
        f"{file_count} files mirrored; "
        f"terminal={terminal})"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
