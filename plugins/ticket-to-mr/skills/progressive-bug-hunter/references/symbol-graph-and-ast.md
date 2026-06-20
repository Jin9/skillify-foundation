# Tier 2: symbol-graph traversal & AST-aware spans

Distilled from research report *Symbol-graph retrieval & AST-aware chunking*.
Use at Tier 2 when the locus is roughly known but the mechanism spans
definitions, callers, callees, or types. Graph-augmented localization is the
**least-contested** structural win in the evidence (LocAgent up to ~92.7%
file-level localization, +12% issue resolution at ~86% lower cost; SweRank SOTA
on SWE-Bench-Lite/LocBench; RANGER strong on CodeSearchNet/CrossCodeEval).

## Traversal pattern (no index build required)
Anchor symbol → expand along structural edges, rerank, read only the expanded
set:
```
anchor (from Tier 1) ─▶ definition site
        │
        ├─▶ callers      (who invokes it — upstream of the failure)
        ├─▶ callees       (what it invokes — downstream of the failure)
        ├─▶ type/interface (contract the value must satisfy)
        └─▶ imports        (cross-file edges)
   read AST-aligned spans of the few structurally-necessary nodes only
```
Build edges from whatever is available in the environment: tree-sitter
`list_code_definition_names` + grep for refs, an existing SCIP/LSIF index, or
language tooling. Per-file/tree-sitter graphs are cheap (a tree-sitter KG
indexed 28M-LOC Linux in ~3 min; ~4× incremental) — but this skill does not
persist one; it traverses on demand.

## AST-aware spans, not line windows
Pull whole semantic units (function/class/method bodies) — cAST-style: split
oversized AST nodes, merge small siblings within a budget. Fixed-line windows
split functions mid-body and lose the structural integrity that explains the
bug. Prepend the enclosing class/module as structural context for the span.

## Regression localization = impact / blast-radius
For "this used to work" bugs, the graph's second use is change-impact
analysis: from the suspected/changed symbol, traverse reachability to find what
the change can actually affect (distinguish "exists in the dependency tree"
from "actually executed on this path"). Combining two impact techniques beats
either alone — corroborate a suspect via both caller-traversal and
diff/blame-scoped reachability.

## Hard ceiling — do not over-trust the graph
Exact static call graphs are undecidable, so static call graphs and slices are
conservative **overapproximations**: edges may exist that no real run takes.
Treat graph-derived suspects as ranked candidates to confirm with concrete
evidence (the failing path / values), never as proof. The AST-chunking *boundary*
advantage is genuinely contested (cAST gains vs. an 864-setting null result,
confounded by embedding model and context budget) — so rely on graph
*traversal* for localization, and treat AST-span packing as a context-hygiene
practice, not a correctness claim. Never extrapolate a benchmark number into
this run's diagnosis.
