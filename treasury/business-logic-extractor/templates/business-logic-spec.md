# Business-logic specification: <scope>

> Descriptive of EXISTING behavior only. Every rule cites real code. No code,
> requirements, tests, or diagrams are produced here. Companion file:
> `loss-ledger.md` (mandatory).

- scope: <modules / services / paths in scope>
- sources cross-referenced: requirements=<path|none> · traces=<path|none>
- extracted: <YYYY-MM-DD> · overall confidence: <high|medium|low>

## Rules

### R1 — <short rule name>
- **Statement (EARS):** When <trigger>, the <system/component> shall <response>.
  _(or Given <state>, When <event>, Then <outcome>.)_
- **Provenance:** `path/to/file.ext:120-138`  ‹quote the load-bearing span›
- **Edge cases:** <enumerated; each cited `path:line`> — or `none found`
- **Decision table** (only if the rule branches):

  | input / condition | … | outcome | provenance |
  |---|---|---|---|
  | <cond A true>  | … | <result> | `file:line` |
  | <cond A false> | … | <result> | `file:line` |

- **Requirement xref:** `<REQ-ID>` (confidence: <h/m/l>, status: proposed|confirmed) — or `none`
- **Trace xref:** `trace_id=… span_id=…` evidence: `exercised|not-observed|dead` — or `none`
- **Confidence:** <high|medium|low> — <one-line why>

### R2 — …
<repeat; one rule per decision/validation/calculation/state-transition/guard.
Never merge an edge case away for brevity.>

## Entities & invariants
<Named domain entities and always-true constraints (EARS Ubiquitous form),
each cited.>

## Open questions for the spec owner
<Anything the code makes ambiguous — list, do not resolve by guessing.>
