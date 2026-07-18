#!/usr/bin/env python3
"""Shared deterministic parser for routine definition files.

Imported by validate_routine.py and check_run.py from the same directory.
Python 3 stdlib only. references/routine-format.md mirrors this parser;
when the two disagree, this parser is the ground truth.
"""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from pathlib import Path

FRONTMATTER_REQUIRED = ("routine", "summary", "version")
FRONTMATTER_OPTIONAL = ("default_on_fail", "owner", "tags", "requires")
NODE_REQUIRED = ("purpose", "executor", "tier", "inputs", "outputs", "gate")
NODE_OPTIONAL = ("executor_mode", "on_fail", "notes")
TIERS = ("small", "mid", "frontier")
GATES = ("none", "before", "after", "both")
ON_FAIL = ("stop", "continue")
NO_INPUT = "none"

KEBAB_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
NODE_HEADING_RE = re.compile(r"^### Node: (.+?)\s*$")
BULLET_RE = re.compile(r"^- ([a-z_]+): (.*)$")


def split_csv(value: str) -> list[str]:
    return [part.strip() for part in value.split(",") if part.strip()]


@dataclass
class Node:
    node_id: str
    line: int
    fields: dict[str, str] = field(default_factory=dict)

    def inputs(self) -> list[str]:
        return split_csv(self.fields.get("inputs", ""))

    def outputs(self) -> list[str]:
        return split_csv(self.fields.get("outputs", ""))


@dataclass
class Routine:
    path: Path
    frontmatter: dict[str, str] = field(default_factory=dict)
    inputs: dict[str, str] = field(default_factory=dict)
    nodes: list[Node] = field(default_factory=list)
    nodes_heading_count: int = 0
    errors: list[str] = field(default_factory=list)

    @property
    def name(self) -> str:
        return self.frontmatter.get("routine", "")


def parse_routine(path: Path) -> Routine:
    routine = Routine(path=path)
    try:
        text = path.read_text(encoding="utf-8")
    except (OSError, UnicodeDecodeError) as exc:
        routine.errors.append(f"unreadable file: {exc}")
        return routine

    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        routine.errors.append("file must start with a --- frontmatter block")
        return routine
    close = None
    for index, line in enumerate(lines[1:], start=1):
        if line.strip() == "---":
            close = index
            break
    if close is None:
        routine.errors.append("frontmatter must close with ---")
        return routine

    for offset, line in enumerate(lines[1:close], start=2):
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue
        if line.startswith((" ", "\t")):
            routine.errors.append(f"line {offset}: frontmatter must be flat - no indented lines")
            continue
        if ":" not in line:
            routine.errors.append(f"line {offset}: invalid frontmatter line: {stripped}")
            continue
        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip().strip("\"'")
        if value in {"", ">", "|", ">-", "|-", ">+", "|+"}:
            routine.errors.append(f"line {offset}: frontmatter key '{key}' must carry a flat scalar value")
            continue
        if key in routine.frontmatter:
            routine.errors.append(f"line {offset}: duplicate frontmatter key: {key}")
            continue
        routine.frontmatter[key] = value

    section = None
    current_node = None
    for index, line in enumerate(lines[close + 1 :], start=close + 2):
        if line.startswith("## "):
            heading = line[3:].strip()
            if heading == "Inputs":
                section = "inputs"
            elif heading == "Nodes":
                section = "nodes"
                routine.nodes_heading_count += 1
            else:
                section = None
            current_node = None
            continue
        node_match = NODE_HEADING_RE.match(line)
        if node_match:
            if section != "nodes":
                routine.errors.append(f"line {index}: '### Node:' heading outside the '## Nodes' section")
                current_node = None
                continue
            current_node = Node(node_id=node_match.group(1).strip(), line=index)
            routine.nodes.append(current_node)
            continue
        bullet = BULLET_RE.match(line)
        if bullet is None:
            continue
        key, value = bullet.group(1), bullet.group(2).strip()
        if section == "inputs":
            if key in routine.inputs:
                routine.errors.append(f"line {index}: duplicate input key: {key}")
            else:
                routine.inputs[key] = value
        elif section == "nodes" and current_node is not None:
            if key in current_node.fields:
                routine.errors.append(f"line {index}: node '{current_node.node_id}' duplicate key: {key}")
            else:
                current_node.fields[key] = value
    return routine
