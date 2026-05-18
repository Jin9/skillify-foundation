# Bug diagnosis: <one-line failure summary>

> The deliverable. Diagnosis only — no code changes, no fix, no PR, no index.

## Ranked suspects
1. `path/to/file.ext:LINE` — <why this is the most likely site, one line>
2. `path/to/other.ext:LINE-LINE` — <why>
3. `…` — <why>   (most-likely first; include confidence: high/med/low)

## Root-cause hypothesis
<The mechanism of the bug in 2–5 sentences, tied to the evidence below.>

Evidence:
- `path:Lx-Ly` — <quoted/described span and what it shows>
- `path:Lx-Ly` — <…>
Contradicting / unexplained: <anything the hypothesis does not yet account for, or "none">
Confidence: <high | medium | low> — <why>

## Retrieval path taken
<Copy the escalation ladder: tiers climbed, queries issued, hit counts, the
escalate/stop decision at each step, and the cache-prefix-preserving order
used. This makes the diagnosis auditable and reproducible.>

## Minimal context bundle (hand-off to a fixer)
The smallest sufficient set — nothing more:
- `path/to/file.ext` lines `A-B` — <role in the bug>
- `path/to/other.ext` lines `C-D` — <role>
- relevant test: `path/to/test.ext::test_name`

## Out of scope (by design)
Fix/patch, PR, regression test authoring, and any code edit are NOT included —
hand the bundle to a coding skill.
