# Layout standards

The fixed geometry every diagram must follow. These constants are the single
source of truth and are enforced byte-for-byte by `scripts/validate_layout.py`.
Do not deviate "to fit" — split into more pages instead (see the bottom of this
file).

## Constants

| Name | Value | Meaning |
|------|-------|---------|
| `GAP_H` | 80px | minimum horizontal gap between side-by-side boxes |
| `GAP_V` | 60px | minimum vertical gap between stacked boxes |
| `PAD` | 40px | minimum padding between a container edge and its children |
| `startSize` | 30px | swimlane/lane title-band height (top inset) |

## Default node sizes

Use one size per node type so rows and columns line up. Widths/heights below are
what the grid formula assumes.

| Node type | width x height |
|-----------|----------------|
| System / service box | 160 x 60 |
| Process step | 160 x 60 |
| Database (cylinder) | 160 x 60 (skeleton) or 120 x 80 (standalone) |
| External system | 160 x 60 |
| Decision (rhombus) | 80 x 80 |
| Actor / user | 40 x 80 |
| Note / annotation | 160 x 50 |

When a row mixes a 60px-tall box with an 80px-tall one, space the row by the
tallest element so the `GAP_V` below it is still satisfied.

## Grid formula

Place every box on a row/column grid. Column index `c` and row index `r` are
0-based. `W` and `H` are the node size for that cell.

```
x = X0 + c * (W + GAP_H)
y = Y0 + r * (H + GAP_V)
```

- **Top-level boxes** (parent is the root cell `1`): `X0 = Y0 = 40` (a margin).
- **Boxes inside a lane/container**: coordinates are parent-relative, so the
  container origin is added automatically by draw.io. Use:

  ```
  x = PAD + c * (W + GAP_H)              = 40, 280, 520, ...
  y = startSize + PAD + r * (H + GAP_V)  = 70, 190, 310, ...   (inside a swimlane)
  ```

  Inside a plain (non-swimlane) container drop the `startSize` term: `y = PAD + r*(H+GAP_V)`.

### Worked example (inside a swimlane)

Two service boxes (160 x 60) side by side, two rows deep:

```
col 0 -> x = 40        col 1 -> x = 40 + 160 + 80 = 280
row 0 -> y = 70        row 1 -> y = 70 + 60 + 60 = 190
```

Horizontal gap = 280 - (40 + 160) = 80 -> exactly `GAP_H`. Vertical gap =
190 - (70 + 60) = 60 -> exactly `GAP_V`. Both pass the gate.

## Container / lane sizing

A lane (swimlane) holding `cols` columns and `rows` rows must be sized so the
40px padding and the internal gaps both hold:

```
laneW = 2*PAD + cols*W + (cols-1)*GAP_H
laneH = startSize + 2*PAD + rows*H + (rows-1)*GAP_V
```

For one 160 x 60 child: `laneW = 80 + 160 = 240`, `laneH = 30 + 80 + 60 = 170`.
That is the lane size used in `templates/grouped-architecture.drawio`.

Padding is checked on all four sides. The top inset is measured **below** the
title band, so a child's relative `y` must be at least `startSize + PAD` (70).

## Lane placement and flow direction

Pick one main flow direction per page and keep it:

- **Left-to-right (LR)** — request/response or pipeline systems (user ->
  frontend -> backend -> database). Lanes are columns placed left to right with
  `>= GAP_H` between them: `laneX[k] = laneX[k-1] + laneW[k-1] + GAP_H`. Edges
  exit right (`exitX=1;exitY=0.5`) and enter left (`entryX=0;entryY=0.5`).
- **Top-to-bottom (TB)** — layered/stack diagrams and process flows. Lanes (or
  rows of boxes) are stacked with `>= GAP_V` between them. Edges exit bottom
  (`exitX=0.5;exitY=1`) and enter top (`entryX=0.5;entryY=0`).

Group related components into the same lane; put unrelated domains in separate
lanes. Standard lane set for architecture: Frontend, Backend, Database,
External Services, Infrastructure.

## Routing

- Use orthogonal edges: `edgeStyle=orthogonalEdgeStyle`.
- If a straight route would cross an unrelated box, add explicit waypoints (an
  `Array as="points"` of `mxPoint`) so the edge elbows around it. Put the
  waypoint in the empty gap between lanes/rows.
- Label any edge whose relationship is not obvious from context.

## When to split into multiple pages

Prefer several readable pages over one dense page. Split when **any** holds:

- A single page would exceed ~12-16 boxes.
- Two domains have no edges between them (give each its own page).
- The laid-out canvas would exceed ~2200px on a side.

Each page is a separate `<diagram>` inside the same `<mxfile>`. **Validate every
page independently** — the gate already checks each page on its own.
