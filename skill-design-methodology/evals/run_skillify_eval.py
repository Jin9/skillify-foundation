#!/usr/bin/env python3
"""Deterministic before/after indicators for the Skillify skill.

The default comparison is:
- before: committed HEAD version of skillify/SKILL.md and mode-playbooks.md
- after: current working-tree skillify/SKILL.md and mode-playbooks.md
"""

from __future__ import annotations

import argparse
import json
import re
import shutil
import subprocess
import sys
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[2]
DEFAULT_PROMPTS = ROOT / "skill-design-methodology" / "evals" / "golden_prompts.json"
DEFAULT_REPORT = ROOT / "skill-design-methodology" / "evals" / "latest_report.md"
SKILL_REL = Path("skillify/SKILL.md")
MODE_REL = Path("skillify/references/mode-playbooks.md")


MODE_KEYWORDS = {
    "Create": ("create a skill", "write a skill.md", "design an agent skill", "build a skill"),
    "Refactor": ("refactor", "improve this skill", "fix this skill.md", "fix this skILL.md"),
    "Review": ("review", "feedback", "is this skill good"),
    "Audit": ("audit", "score", "rubric"),
    "Compress": ("compress", "reduce token", "shrink"),
    "Split": ("split", "does too much"),
    "Merge": ("merge", "combine"),
    "Adapt": ("adapt", "codex", "copilot", "gemini"),
}


@dataclass
class Version:
    name: str
    skill_text: str
    mode_text: str
    skill_dir: Path


def run(cmd: list[str], cwd: Path = ROOT) -> subprocess.CompletedProcess[str]:
    return subprocess.run(cmd, cwd=cwd, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)


def git_show(ref: str, rel: Path) -> str:
    proc = run(["git", "show", f"{ref}:{rel.as_posix()}"])
    if proc.returncode != 0:
        raise SystemExit(f"git show failed for {ref}:{rel}: {proc.stderr.strip()}")
    return proc.stdout


def parse_frontmatter(text: str) -> dict[str, str]:
    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        return {}
    try:
        close = lines[1:].index("---") + 1
    except ValueError:
        return {}
    fields: dict[str, str] = {}
    index = 1
    while index < close:
        line = lines[index]
        if not line.strip() or line.startswith((" ", "\t")) or ":" not in line:
            index += 1
            continue
        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip()
        if value in {">", "|", ">-", "|-"}:
            block: list[str] = []
            index += 1
            while index < close and (not lines[index] or lines[index].startswith((" ", "\t"))):
                block.append(lines[index].strip())
                index += 1
            fields[key] = " ".join(part for part in block if part)
            continue
        fields[key] = value.strip("\"'")
        index += 1
    return fields


def mode_rows(text: str) -> dict[str, dict[str, str]]:
    rows: dict[str, dict[str, str]] = {}
    for line in text.splitlines():
        if not line.startswith("|") or "---" in line or "Mode" in line:
            continue
        parts = [part.strip() for part in line.strip("|").split("|")]
        if len(parts) >= 3 and parts[0] in MODE_KEYWORDS:
            rows[parts[0]] = {"use_when": parts[1], "output": parts[2]}
    return rows


def section(text: str, heading: str) -> str:
    pattern = re.compile(rf"^## {re.escape(heading)}\s*$", re.MULTILINE)
    match = pattern.search(text)
    if not match:
        return ""
    start = match.end()
    next_match = re.search(r"^## ", text[start:], re.MULTILINE)
    end = start + next_match.start() if next_match else len(text)
    return text[start:end]


def has_generic_tree_platforms(text: str) -> bool:
    output = section(text, "Output format")
    return "├── platforms/" in output or "|-- platforms/" in output


def flag_map(version: Version) -> dict[str, bool]:
    text = version.skill_text
    lower = text.lower()
    create = section(text, "Core workflow: Create").lower()
    preamble = section(text, "Universal preamble").lower()
    validation = section(text, "Validation gate").lower()
    constraints = section(text, "Constraints").lower()
    modes = mode_rows(text)
    return {
        "collect_3_prompts": "at least 3 concrete user prompts" in lower,
        "prompt_collection_in_create": "at least 3 concrete user prompts" in create,
        "prompt_collection_in_preamble": "at least 3 concrete user prompts" in preamble,
        "do_not_generate_target_artifact": (
            "do not generate the target artifact" in lower
            or "do not generate the downstream artifact" in lower
        ),
        "frontmatter_rules": "description under 1024" in lower or "description` must be under 1024" in lower,
        "validation_gate": "## validation gate" in lower,
        "validation_gate_script_pointer": "quick_validate.py" in validation and "check_links.py" in validation,
        "preserve_original": "preserve the original" in lower or "preserve_original" in lower,
        "review_audit_no_edit": "review and audit never edit" in lower or "review and audit do not edit" in lower,
        "rubric_gate": "validation-rubric.md" in lower and "40/50" in lower,
        "security_sweep": "security-checklist.md" in lower,
        "progressive_disclosure": "progressive-disclosure.md" in lower or "progressive disclosure" in lower,
        "mode_playbooks": "references/mode-playbooks.md" in lower,
        "negative_triggers": "do not use" in lower or "do not handle" in lower,
        "platform_compatibility": "platform-compatibility.md" in lower or "compatibility:" in lower,
        "repo_policy_boundary": "agents.md" in lower and "claude.md" in lower,
        "one_off_prompt_boundary": "one-off" in lower,
        "init_skill_wired": "init_skill.py <skill-name>" in text,
        "platforms_not_in_generic_tree": not has_generic_tree_platforms(text),
        "iteration_cap_single_sourced": "three passes" not in validation and "three improvement passes" in version.mode_text.lower(),
        "no_50_line_threshold": "50 lines" not in constraints,
        "no_deployment_guide_reference": "platforms/deployment-guide.md" not in text,
        "mode_count_8": len(modes) == 8,
    }


def classify_prompt(prompt: str) -> tuple[bool, str | None]:
    p = prompt.lower()
    if ("not a skill" in p or "build me the react" in p or "one-off" in p or "system prompt" in p):
        return False, None
    if "agents.md" in p or "claude.md" in p or "repository branching policy" in p:
        return False, None
    for mode, keywords in MODE_KEYWORDS.items():
        if any(keyword.lower() in p for keyword in keywords):
            return True, mode
    return False, None


def prompt_results(version: Version, prompts: list[dict[str, Any]]) -> list[dict[str, Any]]:
    flags = flag_map(version)
    modes = mode_rows(version.skill_text)
    results = []
    for item in prompts:
        predicted_use, predicted_mode = classify_prompt(item["prompt"])
        mode_supported = item["expected_mode"] is None or item["expected_mode"] in modes
        output_supported = True
        if item["expected_mode"] in modes and item.get("expected_contract"):
            output_supported = item["expected_contract"].lower().split(",")[0] in modes[item["expected_mode"]]["output"].lower()
        preserve_supported = all(flags.get(name, False) for name in item.get("must_preserve", []))
        passed = (
            predicted_use == item["expected_use"]
            and predicted_mode == item["expected_mode"]
            and mode_supported
            and output_supported
            and preserve_supported
        )
        results.append(
            {
                "id": item["id"],
                "expected_use": item["expected_use"],
                "expected_mode": item["expected_mode"],
                "predicted_use": predicted_use,
                "predicted_mode": predicted_mode,
                "mode_supported": mode_supported,
                "output_supported": output_supported,
                "preserve_supported": preserve_supported,
                "passed": passed,
            }
        )
    return results


def make_temp_skill_dir(base_skill_dir: Path, skill_text: str, mode_text: str) -> Path:
    temp_root = Path(tempfile.mkdtemp(prefix="skillify-eval-"))
    temp_skill = temp_root / "skillify"
    shutil.copytree(base_skill_dir, temp_skill)
    (temp_skill / "SKILL.md").write_text(skill_text, encoding="utf-8")
    (temp_skill / "references" / "mode-playbooks.md").write_text(mode_text, encoding="utf-8")
    return temp_skill


def validator_status(skill_dir: Path, scripts_dir: Path) -> dict[str, Any]:
    quick = run(["python3", str(scripts_dir / "quick_validate.py"), str(skill_dir)])
    links = run(["python3", str(scripts_dir / "check_links.py"), str(skill_dir)])
    banned = list(skill_dir.rglob("README.md"))
    return {
        "quick_validate": quick.returncode == 0,
        "quick_stdout": quick.stdout.strip(),
        "quick_stderr": quick.stderr.strip(),
        "check_links": links.returncode == 0,
        "links_stdout": links.stdout.strip(),
        "links_stderr": links.stderr.strip(),
        "banned_docs_absent": not banned,
    }


def count_quoted_triggers(description: str) -> int:
    return len(re.findall(r'"[^"]+"', description))


def metrics(version: Version, prompts: list[dict[str, Any]], validators: dict[str, Any]) -> dict[str, Any]:
    fields = parse_frontmatter(version.skill_text)
    desc = fields.get("description", "")
    flags = flag_map(version)
    prompt_rows = prompt_results(version, prompts)
    pass_count = sum(1 for row in prompt_rows if row["passed"])
    indicator_flags = {
        "quick_validate": validators["quick_validate"],
        "check_links": validators["check_links"],
        "banned_docs_absent": validators["banned_docs_absent"],
        "description_under_1024": len(desc) <= 1024,
        "trigger_count_at_least_10": count_quoted_triggers(desc) >= 10,
        "mode_count_8": flags["mode_count_8"],
        "target_artifact_boundary": flags["do_not_generate_target_artifact"],
        "repo_policy_boundary": flags["repo_policy_boundary"],
        "one_off_prompt_boundary": flags["one_off_prompt_boundary"],
        "review_audit_no_edit": flags["review_audit_no_edit"],
        "non_create_routes_to_playbooks": flags["mode_playbooks"],
        "prompt_collection_in_create": flags["prompt_collection_in_create"],
        "init_skill_wired": flags["init_skill_wired"],
        "platforms_not_in_generic_tree": flags["platforms_not_in_generic_tree"],
        "validation_gate_script_pointer": flags["validation_gate_script_pointer"],
        "iteration_cap_single_sourced": flags["iteration_cap_single_sourced"],
        "no_50_line_threshold": flags["no_50_line_threshold"],
        "no_deployment_guide_reference": flags["no_deployment_guide_reference"],
    }
    return {
        "line_count": len(version.skill_text.splitlines()),
        "approx_tokens": round(len(version.skill_text) / 4),
        "description_chars": len(desc),
        "quoted_trigger_count": count_quoted_triggers(desc),
        "negative_trigger_count": len(re.findall(r"\bdo not\b", version.skill_text, flags=re.IGNORECASE)),
        "mode_count": len(mode_rows(version.skill_text)),
        "golden_pass": pass_count,
        "golden_total": len(prompt_rows),
        "indicator_pass": sum(1 for value in indicator_flags.values() if value),
        "indicator_total": len(indicator_flags),
        "indicators": indicator_flags,
        "prompt_results": prompt_rows,
        "validators": validators,
    }


def status(value: Any) -> str:
    if isinstance(value, bool):
        return "PASS" if value else "FAIL"
    return str(value)


def delta(before: Any, after: Any) -> str:
    if isinstance(before, bool) and isinstance(after, bool):
        if before == after:
            return "same"
        return "improved" if after else "regressed"
    if isinstance(before, (int, float)) and isinstance(after, (int, float)):
        diff = after - before
        if diff == 0:
            return "same"
        sign = "+" if diff > 0 else ""
        return f"{sign}{diff}"
    return ""


def render_report(before: dict[str, Any], after: dict[str, Any], before_ref: str) -> str:
    lines = [
        "# Skillify Before/After Eval Report",
        "",
        f"Before source: `git show {before_ref}:skillify/...`",
        "After source: working tree `skillify/...`",
        "",
        "## Summary Indicators",
        "",
        "| Indicator | Before | After | Delta |",
        "|---|---:|---:|---|",
        f"| SKILL.md lines | {before['line_count']} | {after['line_count']} | {delta(before['line_count'], after['line_count'])} |",
        f"| Approx tokens | {before['approx_tokens']} | {after['approx_tokens']} | {delta(before['approx_tokens'], after['approx_tokens'])} |",
        f"| Description chars | {before['description_chars']} | {after['description_chars']} | {delta(before['description_chars'], after['description_chars'])} |",
        f"| Quoted trigger count | {before['quoted_trigger_count']} | {after['quoted_trigger_count']} | {delta(before['quoted_trigger_count'], after['quoted_trigger_count'])} |",
        f"| Mode count | {before['mode_count']} | {after['mode_count']} | {delta(before['mode_count'], after['mode_count'])} |",
        f"| Golden prompt pass rate | {before['golden_pass']}/{before['golden_total']} | {after['golden_pass']}/{after['golden_total']} | {delta(before['golden_pass'], after['golden_pass'])} |",
        f"| Structural indicator pass rate | {before['indicator_pass']}/{before['indicator_total']} | {after['indicator_pass']}/{after['indicator_total']} | {delta(before['indicator_pass'], after['indicator_pass'])} |",
        "",
        "## Structural Indicators",
        "",
        "| Indicator | Before | After | Result |",
        "|---|---|---|---|",
    ]
    for name in before["indicators"]:
        b = before["indicators"][name]
        a = after["indicators"][name]
        lines.append(f"| `{name}` | {status(b)} | {status(a)} | {delta(b, a)} |")
    lines.extend(
        [
            "",
            "## Golden Prompt Results",
            "",
            "| Prompt ID | Before | After | Expected mode |",
            "|---|---|---|---|",
        ]
    )
    before_prompts = {row["id"]: row for row in before["prompt_results"]}
    for row in after["prompt_results"]:
        b = before_prompts[row["id"]]
        expected = row["expected_mode"] or "reject"
        lines.append(f"| `{row['id']}` | {status(b['passed'])} | {status(row['passed'])} | {expected} |")
    lines.extend(
        [
            "",
            "## Validator Output",
            "",
            f"- Before quick_validate: {status(before['validators']['quick_validate'])} `{before['validators']['quick_stdout'] or before['validators']['quick_stderr']}`",
            f"- Before check_links: {status(before['validators']['check_links'])} `{before['validators']['links_stdout'] or before['validators']['links_stderr']}`",
            f"- After quick_validate: {status(after['validators']['quick_validate'])} `{after['validators']['quick_stdout'] or after['validators']['quick_stderr']}`",
            f"- After check_links: {status(after['validators']['check_links'])} `{after['validators']['links_stdout'] or after['validators']['links_stderr']}`",
            "",
        ]
    )
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description="Run deterministic Skillify before/after indicators.")
    parser.add_argument("--before-ref", default="HEAD", help="Git ref to use as before baseline.")
    parser.add_argument("--skill-dir", default=str(ROOT / "skillify"), help="Working-tree skill directory.")
    parser.add_argument("--prompts", default=str(DEFAULT_PROMPTS), help="Golden prompts JSON file.")
    parser.add_argument("--report", default=str(DEFAULT_REPORT), help="Markdown report output path.")
    args = parser.parse_args()

    skill_dir = Path(args.skill_dir).resolve()
    prompts = json.loads(Path(args.prompts).read_text(encoding="utf-8"))
    before_skill = git_show(args.before_ref, SKILL_REL)
    before_mode = git_show(args.before_ref, MODE_REL)
    after_skill = (skill_dir / "SKILL.md").read_text(encoding="utf-8")
    after_mode = (skill_dir / "references" / "mode-playbooks.md").read_text(encoding="utf-8")

    before_temp = make_temp_skill_dir(skill_dir, before_skill, before_mode)
    try:
        before_version = Version("before", before_skill, before_mode, before_temp)
        after_version = Version("after", after_skill, after_mode, skill_dir)
        before_validators = validator_status(before_temp, skill_dir / "scripts")
        after_validators = validator_status(skill_dir, skill_dir / "scripts")
        before_metrics = metrics(before_version, prompts, before_validators)
        after_metrics = metrics(after_version, prompts, after_validators)
        report = render_report(before_metrics, after_metrics, args.before_ref)
        report_path = Path(args.report)
        report_path.parent.mkdir(parents=True, exist_ok=True)
        report_path.write_text(report, encoding="utf-8")
        print(report)
        print(f"Report written: {report_path}")
        regressed = (
            after_metrics["indicator_pass"] < before_metrics["indicator_pass"]
            or after_metrics["golden_pass"] < before_metrics["golden_pass"]
        )
        return 1 if regressed else 0
    finally:
        shutil.rmtree(before_temp.parent, ignore_errors=True)


if __name__ == "__main__":
    raise SystemExit(main())
