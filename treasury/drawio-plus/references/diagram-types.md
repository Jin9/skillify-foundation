# Diagram-type recipes

How to apply the universal standards (`references/layout-standards.md`) and styles
(`references/style-catalog.md`) to each diagram family. The spacing/overlap gate
in `scripts/validate_layout.py` applies to **every** type — boxes never overlap,
ever. What changes per type is the layout model and which styles dominate.

## Architecture / system

- **Layout**: grid + lanes. One lane per domain (Frontend, Backend, Database,
  External Services, Infrastructure). LR flow.
- **Start from**: `templates/grouped-architecture.drawio`.
- **Styles**: service boxes for compute, cylinder for data, dashed for external,
  process box for infra components.
- **Edges**: orthogonal, labeled with the relationship (`calls`, `reads/writes`,
  `publishes`).

## Flowchart / process

- **Layout**: single grid column, TB flow. Each step is one row.
- **Start from**: `templates/flow-grid.drawio`.
- **Decision branching**: place the decision rhombus (80 x 80) on its row; branch
  edges leave left/right (`exitX=0`/`exitX=1`, `exitY=0.5`) to boxes in adjacent
  columns, each labeled (`yes`/`no`). Keep `GAP_H` between the branch columns.
- **Styles**: process box for steps, rhombus for decisions, note for annotations.

## Entity-relationship (ER)

- **Layout**: grid of entity boxes; related entities adjacent, no lanes required.
- **Entities**: use the service-box style (or a list shape) titled with the entity
  name; attributes can go in the label with line breaks.
- **Edges**: orthogonal, labeled with cardinality (`1`, `1..*`, `0..1`) near each
  end; set `startArrow`/`endArrow` to convey direction. Keep the full `GAP_H`/
  `GAP_V` so cardinality labels do not collide.

## UML class

- **Layout**: grid of class boxes, grouped by package/module into rows.
- **Class box**: `style="rounded=0;whiteSpace=wrap;html=1;fillColor=#ffffff;strokeColor=#666666;"`
  with the class name in `value` (use a 3-compartment swimlane variant when
  fields/methods are needed: a vertical swimlane with two child rows).
- **Edges**: orthogonal; inheritance `endArrow=block;endFill=0;`, composition
  `startArrow=diamondThin;startFill=1;`.

## Network

- **Layout**: zones as lanes (DMZ, Internal, External), boxes for hosts/services.
- **Styles**: service box for nodes, dashed external for the internet/3rd-party,
  cylinder for storage. Optionally use draw.io network stencils via
  `shape=mscae/...` in the style; keep the standard sizes so the grid holds.
- **Edges**: orthogonal, labeled with protocol/port (`443/tcp`).

## Sequence

- **Layout**: the grid/lane model does **not** apply; use lifelines instead.
  Place actor/participant boxes (160 x 60) in a single row across the top, spaced
  by `GAP_H` (x = 40, 280, 520, ...). The overlap/spacing gate still validates
  that the participant boxes do not overlap.
- **Lifelines**: a dashed vertical edge dropping from each participant's bottom
  center. Messages are horizontal orthogonal edges between lifelines, ordered
  top-to-bottom, each labeled with the call.
- **Note**: the through-box check is advisory here; message arrows legitimately
  run between lifelines. Resolve real crossings, ignore lifeline grazes.

## Mockup / wireframe

- **Layout**: alignment + spacing + no-overlap, but no strict lanes. Snap UI
  blocks (header, nav, content, footer) to a coarse grid and keep `GAP_H`/`GAP_V`
  so elements never touch.
- **Styles**: plain process/box styles or draw.io mockup stencils
  (`shape=mxgraph.mockup...`). Keep box sizes consistent within a region.
- **Gate**: overlap/padding still enforced; routing checks rarely apply (few
  edges). This is the loosest type — when in doubt, favour generous spacing.

## Choosing fast

| Request mentions | Type | Layout model |
|------------------|------|--------------|
| services, API, database, infra, "architecture" | Architecture | grid + lanes (LR) |
| steps, "flow", decision, yes/no | Flowchart | grid column (TB) |
| tables, entities, relationships, cardinality | ER | entity grid |
| classes, inheritance, methods | UML class | class grid |
| hosts, zones, ports, firewall | Network | zone lanes |
| actors exchanging messages over time | Sequence | lifelines |
| screens, UI, pages, buttons | Mockup | coarse grid |
