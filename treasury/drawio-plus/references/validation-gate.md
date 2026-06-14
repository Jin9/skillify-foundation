# Validation gate

`scripts/validate_layout.py` is the deterministic layout gate. A diagram is not
done until it exits 0. Run it on the file you wrote:

```bash
python3 scripts/validate_layout.py path/to/diagram.drawio
```

Add `--strict` to promote through-box findings on edges that have explicit
waypoints from WARN to FAIL.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | no hard FAIL (warnings allowed) — deliver it |
| 1 | at least one hard FAIL — fix and re-run |
| 2 | usage error, unreadable file, or no model found |

It validates **each page independently** and supports bare `mxGraphModel` files
and uncompressed `mxfile`/`diagram` wrappers. A compressed `<diagram>` body is
reported as a WARN and skipped (this skill always writes uncompressed XML).

## Findings and fix recipes

| Finding | Severity | What it means | Fix |
|---------|----------|---------------|-----|
| `overlap` | FAIL | two boxes physically intersect | move one box to the next grid column/row using the formula in `references/layout-standards.md` |
| `h-spacing` | FAIL | side-by-side boxes closer than 80px | increase the column gap to >=80, i.e. set `x = X0 + c*(W+80)` |
| `v-spacing` | FAIL | stacked boxes closer than 60px | increase the row gap to >=60, i.e. set `y = Y0 + r*(H+60)`; space rows by the tallest box |
| `padding` | FAIL | a child is closer than 40px to its container edge | enlarge the lane (`laneW`/`laneH` formula) or move the child inward; remember the top inset is below `startSize` |
| `through-box` | WARN (FAIL with `--strict` + waypoints) | an edge segment appears to cross an unrelated box | add `mxPoint` waypoints so the edge elbows through the empty gap between lanes/rows |
| `self-loop` | WARN | edge whose source equals its target | usually fine; ignore unless unintended |
| `no-geometry` | WARN | edge endpoints could not be resolved | give the edge a real `source` and `target`, or explicit source/target points |
| `compressed-page` | WARN | a diagram page is compressed | re-export uncompressed if you need it validated |
| `empty-page` | WARN | a page has no boxes | remove the empty page or add content |

## How to read it

- Findings are prefixed `FAIL` or `WARN` and name the page, the cells involved,
  and the measured gap, e.g.
  `FAIL h-spacing [Architecture]: 'Web App' (fe-web) and 'API Service' (be-api) are 50px apart, need >=80px`.
- Resolve all FAILs. Resolve `through-box` WARNs too when feasible (add
  waypoints); the remaining WARNs (self-loop, no-geometry on intentional
  free-floating edges) are advisory.
- Re-run until the summary line reads `0 fail`.
