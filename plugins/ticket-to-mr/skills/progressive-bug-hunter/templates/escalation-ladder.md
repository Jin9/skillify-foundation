# Escalation log (working scratchpad — not the deliverable)

> Fill one row per retrieval action. Keep volatile content AFTER the cache
> breakpoint. Stop the moment the working set deterministically explains the
> failure (see references/retrieval-ladder.md → Stop conditions).

## Anchors (Tier 0)
- failure signal: <stack trace | failing test | error message | repro>
- frames (file:line, most-recent-call-last): <...>
- exception/type: <...>   identifiers: <...>   failing assertion: <...>

## Ladder
| # | tier | query / traversal issued | hits | read spans | decision (escalate ↑ / stop ✓) |
|---|------|--------------------------|------|------------|--------------------------------|
| 1 | T1 grep | `rg "<identifier>"` | 0–N | path:Lx-Ly | … |
| 2 | T1 grep | `rg "<error string>"` (parallel) | … | … | … |
| 3 | T2 graph | callers/callees of `<sym>` | … | … | … |
| 4 | T2 graph | blast-radius from `<changed sym>` | … | … | … |
| 5 | T3 broad | iterate: query from finding #k | … | … | … |

## Cache-order note
- stable prefix (unchanged all iterations): system + tools + repo-map
- breakpoint after prefix; rows above are appended AFTER it, deterministically ordered
- re-anchored breakpoint at row #__ (rolling window) — yes/no

## Stop rationale
- working set explains failure? <yes/no — why>
- token budget / diminishing returns? <note>
- impact / blast-radius for severity: <signal noted → tentative P0–P4 (see references/severity-rubric.md)>
