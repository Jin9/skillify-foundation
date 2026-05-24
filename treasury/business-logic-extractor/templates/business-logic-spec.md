# Business-logic specification: <scope>

> Descriptive of EXISTING behavior only. Every rule cites real code. No code,
> requirements, tests, or architecture/topology diagrams are produced here —
> the only diagram is a per-rule decision-logic ASCII flowchart rendering the cited
> rules. Companion file: `loss-ledger.md` (mandatory).

- scope: <modules / services / paths in scope>
- sources cross-referenced: requirements=<path|none> · traces=<path|none>
- extracted: <YYYY-MM-DD> · overall confidence: <high|medium|low>

## Rules

### R1 — <short rule name>
- **Statement (EARS):** When <trigger>, the <system/component> shall <response>.
  _(or Given <state>, When <event>, Then <outcome>.)_
- **Provenance:** `path/to/file.ext:120-138`  ‹quote the load-bearing span›
- **Edge cases:** <enumerated; each cited `path:line`> — or `none found`
- **Requirement xref:** `<REQ-ID>` (confidence: <h/m/l>, status: proposed|confirmed) — or `none`
- **Trace xref:** `trace_id=… span_id=…` evidence: `exercised|not-observed|dead` — or `none`
- **Confidence:** <high|medium|low> — <one-line why>

**Decision-logic flowchart** (ASCII; built only from the cited branches — no invented edge. One flowchart per rule. For a non-branching rule, a minimal `( entry ) --> [ action ] --> ( exit )` flow). Keep the fence at column 0 — NOT indented under a list bullet, so the monospace block stays un-nested and aligned:

```ascii
        ( <entry> )
             |
             v
  +--------------------------+  yes
  | <condition A?>    :line  | ------> [ <outcome A> ]
  +--------------------------+
             | no
             v
  +--------------------------+  yes
  | <condition B?>    :line  | ------> [ <outcome B> ]
  +--------------------------+
             | no
             v
        ( <final outcome> )
```

**Pseudocode** (provenance carrier — every line ends in a resolvable `# path:line`; must enumerate the SAME branches as the flowchart, in code order):

```text
if <condition A>:        <outcome A>      # path/to/file.ext:line-line
elif <condition B>:      <outcome B>      # path/to/file.ext:line
else:                    <final outcome>  # path/to/file.ext:line
```

### R2 — …
<repeat; one rule per decision/validation/calculation/state-transition/guard.
Never merge an edge case away for brevity.>

## Entities & invariants
<Named domain entities and always-true constraints (EARS Ubiquitous form),
each cited.>

## Open questions for the spec owner
<Anything the code makes ambiguous — list, do not resolve by guessing.>
