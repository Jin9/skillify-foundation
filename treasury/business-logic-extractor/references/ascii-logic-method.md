# Rendering decision logic: ASCII flowchart + pseudo-code

How each extracted rule's control flow is rendered. This is a **rendering of
rules already extracted from cited code** — not new analysis, and **not** an
architecture / system-topology / sequence / ER diagram (those stay out of
scope, see `extraction-pitfalls.md`). It shows guard ordering, short-circuit /
early-return, and the path to each outcome that flat prose hides.

## Faithfulness derivation rule (non-negotiable)
The flowchart and pseudo-code are built **from the decision sites enumerated in
step 1 and their cited branches** — never from intuition about what the code
"probably" does. Therefore:
- **Every box, arrow, and pseudo-code line maps to a real branch with a
  `file:line`.** No box or arrow without a citation; no branch in the diagram
  that is absent from the code (a hallucinated edge is the same defect as an
  invented rule).
- **Diagram and pseudo-code enumerate the identical branch/outcome set**
  (parity). If they disagree, one is wrong — reconcile against the code, not
  against each other.
- **Preserve semantics exactly:** keep guard order, short-circuit, and
  early-return as written; never reorder or merge guards to simplify the
  picture.

## Pseudo-code is the provenance carrier of record
Diagram box labels are terse; pseudo-code is where provenance is machine-checked.
Write each branch as one line ending in a **fully-qualified** repo-relative
citation comment `# path/to/file.ext:line(-line)` (the form
`scripts/check_provenance.py` resolves). The rule-level `**Provenance:**` field
remains the per-rule anchor. Fence the block as ```text```.

## ASCII flowchart conventions (keep it minimal and portable)
- Vertical top-down flow: connect nodes with a `|` then a `v` arrowhead. Use
  `LR`-style horizontal `-->` only for the branch that exits to an outcome.
- **Entry / exit / final outcome:** wrap in round parens — `( validate promo )`,
  `( accept as valid )`.
- **Decision:** a `+--+` box holding the condition and a **bare** `:line`
  shorthand — `| now() > expiresAt? :74-77 |`. The taken branch leaves the box
  on a labelled arrow: `-- yes -->` / `-- no -->` (or the literal value).
- **Action / outcome:** square brackets — `[ reject: promo is expired ]`.
- Put only the **bare** `:line` (or `:line-line`) inside a box — the rule's
  `**Provenance:**` field and the pseudo-code name the file. Do NOT put a bare
  filename like `foo.go:12` in a box: the self-check would try to resolve it and
  fail. Use bare `:line`, or a full repo-relative path.
- One flowchart per rule. Do not collapse two rules into one chart.
- **Fence at column 0** — emit ```` ```ascii ```` (and the ```` ```text ````
  pseudo-code) as a top-level block, never indented under a list bullet or inside
  a blockquote, so the monospace stays un-nested and aligned. Introduce it with a
  bold line (e.g. `**Decision-logic flowchart**`) rather than a `-` list item.

## Non-branching rules
Render a minimal linear flow (`( entry ) --> [ single action ] --> ( exit )`) and
a one-line pseudo-code. Do not fabricate a decision box where the code has none.

## No external rendering
The flowchart is **plain ASCII text** — it needs no render service, diagram
engine, or special viewer, and displays identically in every editor, terminal,
and diff. The skill never calls a network endpoint (`curl` / `wget` / `fetch`).
This keeps the extractor offline and side-effect-free.

## Worked example — ordered-guard rule (table → flow)
A "first matching guard wins" validation (guards short-circuit; first hit returns):

```ascii
          ( validate promo )
                  |
                  v
   +-----------------------------+   yes
   | now() > expiresAt?  :74-77  | ------> [ reject: promo is expired ]
   +-----------------------------+
                  | no
                  v
   +-----------------------------+   yes
   | usages >= maxUsages? :79-82 | ------> [ reject: promo is fully used ]
   +-----------------------------+
                  | no
                  v
   +--------------------------------+ yes
   | cartAmt < minOrderAmt?  :84-87 | ---> [ reject: minimum order not met ]
   +--------------------------------+
                  | no
                  v
   +-----------------------------+   yes
   | memberUsageExists?  :89-103 | ------> [ reject: already used by member ]
   +-----------------------------+
                  | no
                  v
          ( accept as valid )  :72-106
```

```text
# validate(promo, cart, member) — first matching guard wins (short-circuit)
if now() > promo.expiresAt:               reject("promo is expired")         # promo-service/app/promo/handler_validate.go:74-77
elif promo.usages >= promo.maxUsages:     reject("promo is fully used")      # promo-service/app/promo/handler_validate.go:79-82
elif cart.amount < promo.minOrderAmount:  reject(minimum-order reason)       # promo-service/app/promo/handler_validate.go:84-87
elif memberUsageExists(promo, member):    reject("already used by member")   # promo-service/app/promo/handler_validate.go:89-103
else:                                     accept as valid                    # promo-service/app/promo/handler_validate.go:72-106
```

Each box and each pseudo-code line is a branch the code actually contains, cited
to its span; the diagram adds nothing the pseudo-code did not already prove.
