---
name: drawio-plus
description: Generate clean, standardized, non-overlapping draw.io / diagrams.net diagrams of any kind (architecture, flow, ER, class, network, sequence, mockup) whose boxes never overlap and whose arrows route around boxes, following fixed spacing and grid standards (80px horizontal gap, 60px vertical gap, 40px container padding) and verified by a deterministic layout check before delivery. Use when the user asks for a clean, standardized, readable, well-spaced, or non-overlapping diagram, to follow diagram layout standards, to group components into frontend, backend, database, external, and infrastructure lanes, or to validate diagram layout. Do NOT use for a quick throwaway sketch where layout does not matter, or for exporting an existing diagram to PNG, SVG, or PDF or running the draw.io CLI; defer those to the basic drawio skill.
compatibility: claude-code, codex, copilot, gemini, antigravity
---

# drawio-plus

Generate standardized draw.io XML where boxes never overlap, arrows route around
boxes, and the layout follows fixed spacing and grid rules. The diagram is not
done until `scripts/validate_layout.py` exits 0.

## When to use

- The user wants a clean, standardized, readable, well-spaced, or non-overlapping
  diagram.
- The user asks to follow diagram layout standards or consistent spacing.
- The user wants components grouped into frontend / backend / database / external
  / infrastructure lanes.
- The user wants the diagram layout validated.

## When NOT to use

- A quick throwaway sketch where layout does not matter — use the basic `drawio`
  skill.
- Exporting a diagram to PNG, SVG, or PDF, or running the draw.io CLI — the basic
  `drawio` skill owns export.
- Editing the content/behaviour of an existing diagram rather than producing a
  standardized layout.

## Hard rules

1. **Spacing minimums**: at least 80px horizontal gap and 60px vertical gap
   between boxes; at least 40px padding inside any container. Constants live in
   `references/layout-standards.md` and are enforced by the gate.
2. **Grid only**: every box sits on a row/column from the formula in
   `references/layout-standards.md`. No ad-hoc coordinates.
3. **One flow direction per page** (left-to-right or top-to-bottom), chosen up
   front and kept consistent.
4. **Group related, separate unrelated** into distinct lanes/containers.
5. **Standard styles only**: take each node and edge style from
   `references/style-catalog.md`; do not invent styles when one fits.
6. **Orthogonal edges** with waypoints so an arrow never runs through an
   unrelated box; label any edge whose relationship is not obvious.
7. **Consistent naming**: ids and labels follow `references/style-catalog.md`.
8. **Anchor on a template**: start from a skeleton in `templates/` whenever the
   request maps to grouped tiers or a linear flow.
9. **Gate before done**: deliver only after `scripts/validate_layout.py` exits 0.

## Inputs

- The architecture or flow to draw (inline text or a file path).
- Optional: output path (default `./<name>.drawio`; never overwrite an existing
  file without confirming), flow direction (LR or TB), and whether the user wants
  an explanation alongside the XML.

## How it works (at a glance)

```text
architecture / flow description
        |
        v
identify entities, relationships, groups, flow direction
        |
        v
pick diagram type + layout model        <- references/diagram-types.md
        |
        v
assign grid coordinates                 <- references/layout-standards.md
        |
        v
apply node + edge styles                <- references/style-catalog.md
        |
        v
assemble XML from a templates/ skeleton
        |
        v
validate_layout.py --FAIL--> widen gaps / add waypoints / split page --+
        | pass                                                         |
        |  <-----------------------------------------------------------+
        v
too dense (>~12-16 boxes)? --yes--> split into pages, re-validate each
        | no
        v
return ONLY the .drawio XML
```

## Workflow

1. **Confirm intent and destination.**
   - Entry: a standardized-diagram request.
   - Resolve the output path (default `./<name>.drawio`; do not overwrite without
     confirmation) and whether an explanation is wanted.
   - Exit: known target and stated output contract.

2. **Analyze the architecture or flow.**
   - Entry: the description.
   - List every entity and classify each by node type (system/service, database,
     external, actor, process, decision, note) and by group (frontend, backend,
     database, external, infrastructure). List relationships as source, target,
     and label.
   - Exit: an entity table and an edge list.

3. **Choose the diagram type and layout strategy.**
   - Entry: the entity and edge lists.
   - Pick the diagram type and its layout model from `references/diagram-types.md`,
     choose the flow direction (LR or TB), and assign each entity to a lane and a
     (row, column). Decide up front whether one page suffices or the system splits
     into logical sections.
   - Exit: a grid assignment for every entity and a page plan.

4. **Compute coordinates.**
   - Entry: the grid assignment.
   - Apply the x/y formulas and container sizing from
     `references/layout-standards.md` — absolute geometry for top-level lanes,
     parent-relative geometry for children.
   - Exit: numeric x, y, width, height for every box.

5. **Apply styles and assemble the XML.**
   - Entry: the coordinates.
   - Start from `templates/grouped-architecture.drawio` (lanes) or
     `templates/flow-grid.drawio` (linear flow). Set each box style and consistent
     id/label from `references/style-catalog.md`; route edges orthogonally, adding
     waypoints where an arrow would otherwise cross a box.
   - Keep the two root cells, write no XML comments, and escape special characters.
   - Exit: one well-formed `.drawio` file.

6. **Validate the layout (gate).**
   - Entry: the assembled file.
   - Run `python3 scripts/validate_layout.py <output.drawio>`. On a non-zero exit,
     apply the matching fix recipe in `references/validation-gate.md` (widen gaps,
     add waypoints, or split the page) and re-run until it exits 0. Do not skip.
   - Exit: the gate exits 0.

7. **Report per the output contract.**
   - Entry: a passing gate.
   - Return only the final XML (or the written path); add a short explanation only
     if the user asked for one.
   - Exit: XML delivered, no extra files.

## Output contract

- When the user asks for a diagram file or content, return **only** the final
  draw.io XML — a complete `mxGraphModel` document directly importable into
  diagrams.net (or the written `.drawio` path if a file was requested). No
  surrounding prose by default.
- Include a short explanation **only when asked**, placed after the XML.
- The XML is well-formed, has no XML comments, uses catalog styles, and passes
  `scripts/validate_layout.py` with exit 0.
- Create no companion files; export to PNG/SVG/PDF is out of scope.

## Constraints and anti-patterns

- DO NOT place boxes off-grid or with sub-minimum gaps to make a diagram "fit".
- DO NOT route an edge straight through an unrelated box; add waypoints.
- DO NOT invent ad-hoc styles when a catalog style fits.
- DO NOT cram an unreadable single page when splitting into pages is clearer.
- DO NOT write XML comments or drop the root cells `id="0"` and `id="1"`.
- MUST finish on a passing gate.

## Validation checklist

- [ ] Exactly one `.drawio` file written (no companions).
- [ ] `scripts/validate_layout.py` exited 0.
- [ ] Every box is on-grid with at least 80px horizontal / 60px vertical gaps.
- [ ] Every container has at least 40px padding around its children.
- [ ] Edges are orthogonal and labeled where the relationship is not obvious.
- [ ] Ids and labels follow the naming convention.
- [ ] The XML imports cleanly into diagrams.net.

## References

| Need | File |
|------|------|
| Spacing constants, grid formula, lane sizing, page-split rules | `references/layout-standards.md` |
| Exact mxStyle per node/edge type + id/label naming | `references/style-catalog.md` |
| Per-type layout recipes (architecture, flow, ER, class, network, sequence, mockup) | `references/diagram-types.md` |
| What the gate checks, exit codes, fix recipes | `references/validation-gate.md` |
| Worked end-to-end example | `examples/three-tier-web-app.md` |

## Templates and scripts

- `templates/grouped-architecture.drawio` — five-lane skeleton (frontend,
  backend, database, external, infrastructure), passes the gate as shipped.
- `templates/flow-grid.drawio` — linear top-to-bottom process skeleton, passes the
  gate as shipped.
- `scripts/validate_layout.py` — deterministic layout gate (overlap, spacing,
  padding as FAIL; arrow-through-box as WARN). Run it on the file you wrote.
