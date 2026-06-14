# Example: three-tier web app

A worked trace from request to a gate-passing diagram. It shows the analysis,
the grid math, an XML excerpt, and the validation result.

## Request

> Draw a clean architecture diagram of our shop: a React web app and a mobile
> app talk to an API service, which reads/writes a Postgres database and calls
> Stripe for payments. Redis caches sessions.

## Step 2 — analysis

Entity table:

| Entity | Node type | Group | Id |
|--------|-----------|-------|----|
| Web App | service | frontend | `fe-web` |
| Mobile App | service | frontend | `fe-mobile` |
| API Service | service | backend | `be-api` |
| Postgres | database | database | `db-postgres` |
| Stripe | external | external | `ext-stripe` |
| Redis | service | infrastructure | `infra-redis` |

Edge list (source, target, label):

- `fe-web` -> `be-api` : "calls"
- `fe-mobile` -> `be-api` : "calls"
- `be-api` -> `db-postgres` : "reads/writes"
- `be-api` -> `ext-stripe` : "charges"
- `be-api` -> `infra-redis` : "caches"

## Step 3-4 — layout and coordinates

LR flow, lanes left to right. Frontend has two stacked boxes, so its lane is two
rows tall; the others have one. Using `references/layout-standards.md`:

- Frontend lane: 1 column x 2 rows -> `laneH = 30 + 2*40 + 2*60 + 1*60 = 290`.
  Children at relative `(40, 70)` and `(40, 190)` (gap 60).
- One-box lanes: `240 x 170`, child at relative `(40, 70)`.
- Lane x positions (each `laneW = 240`, gap 80): `40, 360, 680, 1000, 1320`.

## Step 5 — assembled XML (excerpt)

Start from `templates/grouped-architecture.drawio`, rename ids/labels, add the
second frontend box, and wire the edges:

```xml
<mxCell id="lane-fe" value="Frontend" style="swimlane;startSize=30;html=1;collapsible=0;fillColor=none;strokeColor=#999999;fontStyle=1;" vertex="1" parent="1">
  <mxGeometry x="40" y="40" width="240" height="290" as="geometry" />
</mxCell>
<mxCell id="fe-web" value="Web App" style="rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;" vertex="1" parent="lane-fe">
  <mxGeometry x="40" y="70" width="160" height="60" as="geometry" />
</mxCell>
<mxCell id="fe-mobile" value="Mobile App" style="rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;" vertex="1" parent="lane-fe">
  <mxGeometry x="40" y="190" width="160" height="60" as="geometry" />
</mxCell>
<mxCell id="e-web-api" value="calls" style="edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;endArrow=block;jettySize=auto;exitX=1;exitY=0.5;entryX=0;entryY=0.5;" edge="1" parent="1" source="fe-web" target="be-api">
  <mxGeometry relative="1" as="geometry">
    <Array as="points"><mxPoint x="320" y="140" /></Array>
  </mxGeometry>
</mxCell>
```

The Database, External Services, and Infrastructure lanes follow the shipped
template unchanged except for ids/labels (`db-postgres`, `ext-stripe`,
`infra-redis`).

`be-api` connects to non-adjacent lanes (External, Infrastructure), so those
edges must route **over the top** to avoid crossing the Database lane. Send them
up into the margin above the lanes (lanes start at `y=40`) and across, on two
different heights so the two arrows do not overlap:

```xml
<mxCell id="e-api-stripe" value="charges" style="edgeStyle=orthogonalEdgeStyle;...;exitX=0.5;exitY=0;entryX=0.5;entryY=0;" edge="1" parent="1" source="be-api" target="ext-stripe">
  <mxGeometry relative="1" as="geometry">
    <Array as="points"><mxPoint x="480" y="25" /><mxPoint x="1120" y="25" /></Array>
  </mxGeometry>
</mxCell>
```

## Step 6 — gate

```bash
$ python3 scripts/validate_layout.py shop-architecture.drawio
validate_layout: 0 fail, 0 warn
```

Exit 0 — every box is on-grid, lanes are 80px apart, children have 40px padding,
and the labeled edges elbow through the inter-lane gaps. The diagram is delivered
as XML only.

## What would have failed

- Putting `fe-mobile` at relative `y=140` instead of `190` -> `v-spacing` FAIL
  (gap 10px, need 60).
- Sizing the Frontend lane at `170` while holding two rows -> `padding` FAIL on
  the lower box.
- A straight `be-api` -> `ext-stripe` edge (no waypoints) -> `through-box` WARN:
  it crosses the Database lane and `db-postgres`. Routing it over the top clears
  the warning.
