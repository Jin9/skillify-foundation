#!/usr/bin/env python3
"""Deterministic layout gate for the drawio-plus skill.

Parses an uncompressed draw.io / mxGraphModel file and reports layout
violations. It is the hard contract that a generated diagram must pass
before it is delivered.

Hard FAIL (blocks, exit 1):
  - overlap    : two boxes physically intersect
  - h-spacing  : side-by-side boxes closer than the horizontal minimum
  - v-spacing  : stacked boxes closer than the vertical minimum
  - padding    : a child sits closer than the minimum to its container edge

Advisory WARN (never blocks on its own):
  - through-box: an edge segment appears to cross an unrelated box. A static
                 parser cannot reproduce draw.io's auto-router, so this stays
                 a warning unless --strict AND the edge has explicit waypoints.
  - self-loop / no-geometry / compressed page / empty model: skipped, reported.

Exit codes:
  0  no hard FAIL (warnings allowed)
  1  at least one hard FAIL
  2  usage error, unreadable file, or no model found (distinct from a layout fail)

Stdlib only; Python 3.9 compatible.
"""
from __future__ import annotations

import argparse
import math
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple
import xml.etree.ElementTree as ET

# Layout standard — single source of truth, mirrored in references/layout-standards.md
GAP_H = 80.0          # minimum horizontal gap between side-by-side boxes
GAP_V = 60.0          # minimum vertical gap between stacked boxes
PAD = 40.0            # minimum padding inside a container/group
DEFAULT_START_SIZE = 30.0  # swimlane title-band height when style omits startSize
EDGE_INSET = 4.0      # shrink boxes when testing edge crossings (edges attach at borders)
MIN_PENETRATION = 8.0  # ignore corner grazes shorter than this when flagging crossings
EPS = 0.01            # float tolerance


class Cell:
    __slots__ = (
        "cid", "parent", "is_vertex", "is_edge", "source", "target",
        "style", "value", "has_geom", "x", "y", "w", "h",
        "points", "src_pt", "tgt_pt",
    )

    def __init__(self, cid: str, parent: str) -> None:
        self.cid = cid
        self.parent = parent
        self.is_vertex = False
        self.is_edge = False
        self.source: Optional[str] = None
        self.target: Optional[str] = None
        self.style = ""
        self.value = ""
        self.has_geom = False
        self.x = 0.0
        self.y = 0.0
        self.w = 0.0
        self.h = 0.0
        self.points: List[Tuple[float, float]] = []
        self.src_pt: Optional[Tuple[float, float]] = None
        self.tgt_pt: Optional[Tuple[float, float]] = None


def _f(value: Optional[str], default: float = 0.0) -> float:
    if value is None:
        return default
    try:
        return float(value)
    except ValueError:
        return default


def _style_num(style: str, key: str, default: float) -> float:
    for part in style.split(";"):
        part = part.strip()
        if part.startswith(key + "="):
            return _f(part.split("=", 1)[1], default)
    return default


def parse_cells(model: ET.Element) -> Dict[str, Cell]:
    cells: Dict[str, Cell] = {}
    for mc in model.iter("mxCell"):
        cid = mc.get("id")
        if cid is None:
            continue
        cell = Cell(cid, mc.get("parent", ""))
        cell.is_vertex = mc.get("vertex") == "1"
        cell.is_edge = mc.get("edge") == "1"
        cell.source = mc.get("source")
        cell.target = mc.get("target")
        cell.style = mc.get("style", "") or ""
        cell.value = mc.get("value", "") or ""
        geom = mc.find("mxGeometry")
        if geom is not None:
            cell.x = _f(geom.get("x"))
            cell.y = _f(geom.get("y"))
            cell.w = _f(geom.get("width"))
            cell.h = _f(geom.get("height"))
            cell.has_geom = cell.w > 0 and cell.h > 0
            for child in geom:
                if child.tag == "mxPoint":
                    as_attr = child.get("as")
                    pt = (_f(child.get("x")), _f(child.get("y")))
                    if as_attr == "sourcePoint":
                        cell.src_pt = pt
                    elif as_attr == "targetPoint":
                        cell.tgt_pt = pt
                elif child.tag == "Array" and child.get("as") == "points":
                    for p in child.findall("mxPoint"):
                        cell.points.append((_f(p.get("x")), _f(p.get("y"))))
        cells[cid] = cell
    return cells


def build_abs_pos(cells: Dict[str, Cell]):
    memo: Dict[str, Tuple[float, float]] = {}

    def abs_pos(cid: str) -> Tuple[float, float]:
        if cid in memo:
            return memo[cid]
        cell = cells.get(cid)
        if cell is None or not cell.has_geom:
            memo[cid] = (0.0, 0.0)
            return memo[cid]
        parent = cells.get(cell.parent)
        if parent is not None and parent.has_geom:
            pox, poy = abs_pos(cell.parent)
        else:
            pox, poy = 0.0, 0.0
        res = (pox + cell.x, poy + cell.y)
        memo[cid] = res
        return res

    return abs_pos


def is_ancestor(cells: Dict[str, Cell], a_id: str, b_id: str) -> bool:
    """True if a_id is an ancestor (any depth) of b_id."""
    seen = set()
    cur = cells.get(b_id)
    cur = cells.get(cur.parent) if cur else None
    while cur is not None and cur.cid not in seen:
        if cur.cid == a_id:
            return True
        seen.add(cur.cid)
        cur = cells.get(cur.parent)
    return False


def _outcode(x: float, y: float, x0: float, y0: float, x1: float, y1: float) -> int:
    code = 0
    if x < x0:
        code |= 1
    elif x > x1:
        code |= 2
    if y < y0:
        code |= 4
    elif y > y1:
        code |= 8
    return code


def seg_rect_overlap_len(px: float, py: float, qx: float, qy: float,
                         rect: Tuple[float, float, float, float]) -> float:
    """Length of the portion of segment pq lying inside rect (0 if none).

    Cohen-Sutherland clip; used to measure how deeply an edge penetrates a box.
    """
    x0, y0, x1, y1 = rect
    ax, ay, bx, by = px, py, qx, qy
    oa = _outcode(ax, ay, x0, y0, x1, y1)
    ob = _outcode(bx, by, x0, y0, x1, y1)
    while True:
        if not (oa | ob):
            return math.hypot(bx - ax, by - ay)
        if oa & ob:
            return 0.0
        oc = oa or ob
        if oc & 8:
            if by == ay:
                return 0.0
            x = ax + (bx - ax) * (y1 - ay) / (by - ay)
            y = y1
        elif oc & 4:
            if by == ay:
                return 0.0
            x = ax + (bx - ax) * (y0 - ay) / (by - ay)
            y = y0
        elif oc & 2:
            if bx == ax:
                return 0.0
            y = ay + (by - ay) * (x1 - ax) / (bx - ax)
            x = x1
        else:
            if bx == ax:
                return 0.0
            y = ay + (by - ay) * (x0 - ax) / (bx - ax)
            x = x0
        if oc == oa:
            ax, ay = x, y
            oa = _outcode(ax, ay, x0, y0, x1, y1)
        else:
            bx, by = x, y
            ob = _outcode(bx, by, x0, y0, x1, y1)


def label(cell: Cell) -> str:
    name = cell.value.strip() or "(unnamed)"
    return "'{}' ({})".format(name, cell.cid)


def validate_model(model: ET.Element, page: str, strict: bool,
                   fails: List[str], warns: List[str]) -> None:
    cells = parse_cells(model)
    abs_pos = build_abs_pos(cells)
    boxes = [c for c in cells.values() if c.is_vertex and c.has_geom]

    if not boxes:
        warns.append("WARN empty-page: {} has no boxes to validate".format(page))
        return

    # --- Box overlap + spacing (HARD FAIL) ---
    for i in range(len(boxes)):
        for j in range(i + 1, len(boxes)):
            a, b = boxes[i], boxes[j]
            if is_ancestor(cells, a.cid, b.cid) or is_ancestor(cells, b.cid, a.cid):
                continue  # intentional containment is legal
            ax0, ay0 = abs_pos(a.cid)
            ax1, ay1 = ax0 + a.w, ay0 + a.h
            bx0, by0 = abs_pos(b.cid)
            bx1, by1 = bx0 + b.w, by0 + b.h
            overlap_x = ax0 < bx1 - EPS and bx0 < ax1 - EPS
            overlap_y = ay0 < by1 - EPS and by0 < ay1 - EPS
            if overlap_x and overlap_y:
                fails.append("FAIL overlap [{}]: {} and {} physically overlap".format(
                    page, label(a), label(b)))
            elif overlap_y:
                gap = max(ax0, bx0) - min(ax1, bx1)
                if gap < GAP_H - EPS:
                    fails.append(
                        "FAIL h-spacing [{}]: {} and {} are {:.0f}px apart, need >={:.0f}px".format(
                            page, label(a), label(b), gap, GAP_H))
            elif overlap_x:
                gap = max(ay0, by0) - min(ay1, by1)
                if gap < GAP_V - EPS:
                    fails.append(
                        "FAIL v-spacing [{}]: {} and {} are {:.0f}px apart, need >={:.0f}px".format(
                            page, label(a), label(b), gap, GAP_V))

    # --- Container padding (HARD FAIL) ---
    for c in boxes:
        parent = cells.get(c.parent)
        if parent is None or not parent.is_vertex or not parent.has_geom:
            continue
        start = _style_num(parent.style, "startSize", DEFAULT_START_SIZE) \
            if "swimlane" in parent.style else 0.0
        insets = (
            ("left", c.x),
            ("top", c.y - start),
            ("right", parent.w - (c.x + c.w)),
            ("bottom", parent.h - (c.y + c.h)),
        )
        for side, val in insets:
            if val < PAD - EPS:
                fails.append(
                    "FAIL padding [{}]: {} is {:.0f}px from {} edge of container {}, need >={:.0f}px".format(
                        page, label(c), val, side, label(parent), PAD))

    # --- Arrow-through-box (advisory WARN; --strict promotes explicit-waypoint edges) ---
    box_by_id = {c.cid: c for c in boxes}
    for e in cells.values():
        if not e.is_edge:
            continue
        if e.source and e.source == e.target:
            warns.append("WARN self-loop [{}]: edge {} skipped".format(page, label(e)))
            continue
        pox, poy = abs_pos(e.parent)
        start_pt = None
        if e.src_pt is not None:
            start_pt = (pox + e.src_pt[0], poy + e.src_pt[1])
        elif e.source in box_by_id:
            s = box_by_id[e.source]
            sx, sy = abs_pos(s.cid)
            start_pt = (sx + s.w / 2.0, sy + s.h / 2.0)
        end_pt = None
        if e.tgt_pt is not None:
            end_pt = (pox + e.tgt_pt[0], poy + e.tgt_pt[1])
        elif e.target in box_by_id:
            t = box_by_id[e.target]
            tx, ty = abs_pos(t.cid)
            end_pt = (tx + t.w / 2.0, ty + t.h / 2.0)
        if start_pt is None or end_pt is None:
            warns.append("WARN no-geometry [{}]: edge {} has unresolved endpoints, "
                         "routing not checked".format(page, label(e)))
            continue
        poly = [start_pt] + [(pox + wx, poy + wy) for (wx, wy) in e.points] + [end_pt]
        has_waypoints = bool(e.points)
        crossed: List[Cell] = []
        for box in boxes:
            if box.cid in (e.source, e.target):
                continue
            if is_ancestor(cells, box.cid, e.source or "") or \
                    is_ancestor(cells, box.cid, e.target or ""):
                continue  # passing through an endpoint's own container is fine
            bx, by = abs_pos(box.cid)
            rect = (bx + EDGE_INSET, by + EDGE_INSET,
                    bx + box.w - EDGE_INSET, by + box.h - EDGE_INSET)
            if rect[2] <= rect[0] or rect[3] <= rect[1]:
                continue
            hit = any(
                seg_rect_overlap_len(poly[k][0], poly[k][1], poly[k + 1][0], poly[k + 1][1], rect) > MIN_PENETRATION
                for k in range(len(poly) - 1)
            )
            if hit:
                crossed.append(box)
        for box in crossed:
            if strict and has_waypoints:
                fails.append("FAIL through-box [{}]: edge {} crosses box {}".format(
                    page, label(e), label(box)))
            else:
                suffix = "" if has_waypoints else " (auto-routed; add waypoints to confirm or avoid)"
                warns.append("WARN through-box [{}]: edge {} may cross box {}{}".format(
                    page, label(e), label(box), suffix))


def find_models(root: ET.Element, warns: List[str]) -> List[Tuple[ET.Element, str]]:
    models: List[Tuple[ET.Element, str]] = []
    if root.tag == "mxGraphModel":
        models.append((root, "page-1"))
        return models
    diagrams = list(root.iter("diagram"))
    if diagrams:
        for idx, diag in enumerate(diagrams, 1):
            name = diag.get("name") or "page-{}".format(idx)
            model = diag.find("mxGraphModel")
            if model is not None:
                models.append((model, name))
            elif (diag.text or "").strip():
                warns.append("WARN compressed-page: '{}' body is compressed and "
                             "cannot be validated; skipped".format(name))
        return models
    for idx, model in enumerate(root.iter("mxGraphModel"), 1):
        models.append((model, "page-{}".format(idx)))
    return models


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(
        description="Validate draw.io diagram layout (spacing, overlap, padding, routing).")
    parser.add_argument("path", help="path to a .drawio / mxGraphModel XML file")
    parser.add_argument("--strict", action="store_true",
                        help="promote through-box findings on edges with explicit "
                             "waypoints from WARN to FAIL")
    args = parser.parse_args(argv)

    path = Path(args.path)
    if not path.is_file():
        print("error: not a file: {}".format(path), file=sys.stderr)
        return 2
    try:
        root = ET.parse(str(path)).getroot()
    except ET.ParseError as exc:
        print("error: XML parse error: {}".format(exc), file=sys.stderr)
        return 2

    fails: List[str] = []
    warns: List[str] = []
    models = find_models(root, warns)
    if not models:
        for line in warns:
            print(line)
        print("error: no mxGraphModel found in {}".format(path), file=sys.stderr)
        return 2

    for model, page in models:
        validate_model(model, page, args.strict, fails, warns)

    for line in warns:
        print(line)
    for line in fails:
        print(line)
    print("validate_layout: {} fail, {} warn".format(len(fails), len(warns)))
    return 1 if fails else 0


if __name__ == "__main__":
    sys.exit(main())
