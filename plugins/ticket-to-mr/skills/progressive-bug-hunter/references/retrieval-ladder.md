# Retrieval ladder: when to escalate, when to stop

Distilled from research reports *Selective context retrieval for code agents*
and *Agentic grep vs. dense-vector RAG*. Core asymmetry: agentic lexical search
is an iterative read→grep→decide loop over the raw corpus that fails **loudly**
(empty result = clean signal); dense-vector RAG is a single fixed top-k
projection that fails **silently** (a plausible wrong neighbor that corrupts
reasoning) and cannot recover evidence the initial top-k dropped. For bug
localization the loud-failure property is the reason to start lexical.

## The ladder (climb only as far as needed)

- **Tier 0 — anchor.** Parse the failure (stack frames `file:line`, exception
  type, failing assertion, identifiers). These are the highest-signal seeds;
  oracle-quality localization is ~6× BM25 on SWE-bench, so a good anchor is
  worth more than any retriever.
- **Tier 1 — agentic grep (default for narrow fixes).** Exact-identifier /
  rare-token / structured-code lookups: `ripgrep`/`grep`, `find`, tree-sitter
  `list_code_definition_names`. Index-free, maintenance-free, cache-friendly.
  Parallelize independent searches — N greps in one turn ≈ one grep's latency
  while exploring N× the tree (~4× faster end-to-end). The bulk of real bugs
  resolve here; a top SWE-bench agent found grep/find sufficient and embeddings
  not the bottleneck (setting-dependent — small contained repo).
- **Tier 2 — structural (escalate when Tier 1 underdetermines).** The locus is
  known-ish but the *mechanism* spans definitions/callers/types. Follow the
  symbol/call/import/type graph; pull AST-aligned spans. See
  `symbol-graph-and-ast.md`. Graph-guided localization shows large, repeatedly
  reproduced gains (LocAgent up to ~92.7% file-level; SweRank SOTA).
- **Tier 3 — broad cross-file unknown (escalate only when locus unknown).**
  Hybrid lexical+structural widening; on a paraphrastic/semantic gap, dense
  retrieval *would* help — but this skill localizes with available search
  tools and does **not** build an index. Iterate (RepoCoder-style: feed each
  finding back as the next query; +10–30 pts over one-shot). Never one-shot a
  single broad query.

## Decision boundary (route by query/corpus class)
| signal | tier |
|--------|------|
| exact identifier, error string, stack frame | Tier 1 |
| "who calls / what breaks if" , regression blast radius | Tier 2 |
| paraphrastic "where is the logic that…", unknown locus | Tier 3 (iterate) |

## Stop conditions (budget tokens — retrieved tokens are not free)
- Stop climbing the moment the working set deterministically explains the
  failure (MoTCoder-style uncertainty gating cuts ~38% tokens at no pass@1 loss).
- Stop and report if Tier 3 iterations stop yielding new relevant spans
  (avoid the long-irrelevant-dump failure of grep).
- Never widen "just in case" — the repo-map is cheap, retrieved spans are not;
  long-context packing loses ~20% accuracy in the middle (lost-in-the-middle)
  and retrieval+short-context beats long-context on localization.
