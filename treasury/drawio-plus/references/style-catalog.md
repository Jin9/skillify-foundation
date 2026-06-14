# Style catalog

Copy these exact `style` strings onto the `mxCell` for each node type. Using the
catalog keeps colour, shape, and stroke consistent across every diagram. Do not
invent ad-hoc styles when a catalog entry fits.

## Node styles

| Node type | `style` string |
|-----------|----------------|
| System / service box | `rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;` |
| Database (cylinder) | `shape=cylinder3;whiteSpace=wrap;html=1;boundedLbl=1;backgroundOutline=1;fillColor=#d5e8d4;strokeColor=#82b366;` |
| External system | `rounded=1;whiteSpace=wrap;html=1;dashed=1;fillColor=#ffe6cc;strokeColor=#d79b00;` |
| Actor / user | `shape=umlActor;verticalLabelPosition=bottom;verticalAlign=top;html=1;outlineConnect=0;` |
| Process step | `rounded=0;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;` |
| Decision (rhombus) | `rhombus;whiteSpace=wrap;html=1;fillColor=#fff2cc;strokeColor=#d6b656;` |
| Note / annotation | `shape=note;whiteSpace=wrap;html=1;fillColor=#fff2cc;strokeColor=#d6b656;size=14;` |
| Lane / group container | `swimlane;startSize=30;html=1;collapsible=0;fillColor=none;strokeColor=#999999;fontStyle=1;` |

The colours are draw.io's standard palette swatches (blue = compute, green =
data, orange = external, grey = process/neutral, yellow = decision/note), so the
diagram reads correctly even in greyscale by shape alone.

## Edge styles

| Edge | `style` string |
|------|----------------|
| Orthogonal (default) | `edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;endArrow=block;startArrow=none;jettySize=auto;` |
| LR connection | add `exitX=1;exitY=0.5;entryX=0;entryY=0.5;` |
| TB connection | add `exitX=0.5;exitY=1;entryX=0.5;entryY=0;` |
| Bidirectional | add `startArrow=block;` |

Route around a box by adding waypoints to the edge geometry:

```xml
<mxCell id="e-web-api" value="calls" style="edgeStyle=orthogonalEdgeStyle;..." edge="1" parent="1" source="fe-web" target="be-api">
  <mxGeometry relative="1" as="geometry">
    <Array as="points">
      <mxPoint x="320" y="140" />
    </Array>
  </mxGeometry>
</mxCell>
```

Set the edge `value` to label the relationship when it is not obvious (for
example `reads/writes`, `publishes`, `on success`).

## ID and label naming convention

Deterministic ids make diffs and edits readable.

- **Group prefixes** (fixed): `fe` frontend, `be` backend, `db` database,
  `ext` external services, `infra` infrastructure.
- **Vertex id**: `<group>-<slug>` in lower kebab-case, for example `fe-web`,
  `be-api`, `db-orders`, `ext-stripe`, `infra-redis`. Lane container ids use
  `lane-<group>`, for example `lane-be`.
- **Edge id**: `e-<source-slug>-<target-slug>`, for example `e-web-api`.
- **Label** (`value`): a human title in Title Case, for example `Web App`,
  `Orders DB`, `Payment Gateway`. Keep labels short; wrapping is on
  (`whiteSpace=wrap`).

## XML hygiene (inherited hard rules)

- Never write XML comments inside the model.
- Escape `&`, `<`, `>`, `"` inside attribute values (`&amp;` `&lt;` `&gt;`
  `&quot;`).
- Every `mxCell` needs a unique `id`.
- Every edge needs a child `<mxGeometry relative="1" as="geometry" />`.
- Keep the two root cells `id="0"` and `id="1"`.

(Angle brackets are intentional here and in code fences; they are forbidden only
in the SKILL.md frontmatter.)
