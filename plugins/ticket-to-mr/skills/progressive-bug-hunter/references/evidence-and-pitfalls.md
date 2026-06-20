# Evidence base & failure modes

Distilled, cross-cutting, from all four source reports. Use to justify
escalation decisions and to avoid known traps. These are calibration priors —
never paste a benchmark number into the actual diagnosis as if it were
evidence about *this* bug.

## What the evidence supports (why the ladder is ordered this way)
- Good localization dominates everything: oracle retrieval ≈ 6× BM25 on
  SWE-bench → invest in anchors (Tier 0) before any retriever.
- Agentic grep is sufficient for the bulk of narrow fixes and fails loudly;
  a top SWE-bench agent found embeddings were not the bottleneck (caveat:
  contained repo — SWE-bench pulls the whole repo into a container).
- Graph-guided localization is the strongest, least-contested structural
  signal (LocAgent, SweRank, RANGER) → Tier 2 before broad widening.
- Retrieval + short context beats long-context packing on localization;
  lost-in-the-middle costs ~20% absolute in the middle of a long window →
  never "just load more".
- Parallel independent greps ≈ one grep's latency for N× coverage → batch
  Tier-1 queries.
- Iterative retrieve→generate adds +10–30 pts over one-shot (RepoCoder) →
  always feed findings back.

## Failure modes to detect and avoid
- **Silent dense miss:** a vector/semantic search returns a plausible-but-wrong
  neighbor; the model then confidently builds on it. Prefer loud lexical
  failure; if a semantic tool is used, corroborate every hit against concrete
  evidence.
- **Chunk-boundary loss:** a logical unit split across chunks looks irrelevant
  in each half; no downstream step can recover it. Pull AST-aligned whole
  units (see `symbol-graph-and-ast.md`).
- **Stale/index drift:** any pre-built index can be out of date (embedding
  upgrade → up to ~20% accuracy decline; previously-top docs invisible). This
  skill is index-free precisely to avoid this; if an environment index exists,
  treat it as a hint, verify against the live tree.
- **Static-graph overapproximation:** call graphs/slices are undecidable →
  conservative; a graph edge is a candidate, not proof. Confirm with the
  failing path/values.
- **Cache thrash:** volatile results in the cached prefix → rewrite every
  iteration, latency up. Keep static-first/dynamic-last (`cache-prefix-ordering.md`).
- **Contested-baseline trap:** the AST-chunking benefit is genuinely disputed
  (cAST gains vs. an 864-setting null result). Do not present a contested
  research claim as a fact in the report; the report states *this run's*
  evidence only.

## Scope guard
This skill *localizes and diagnoses*. It does not write the fix, open a PR,
build an index, review code generally, or profile performance — those are
different skills. The deliverable is a diagnosis + a minimal context bundle a
fixer consumes next.
