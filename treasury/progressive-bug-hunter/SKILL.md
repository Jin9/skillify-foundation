---
name: progressive-bug-hunter
description: >
  Localize and diagnose a bug by progressively retrieving the MINIMAL
  sufficient code context: start with cheap agentic grep / structured search,
  escalate to symbol-graph / call-graph / AST-aware retrieval only when needed,
  order retrieval to preserve the prompt-cache prefix, and avoid whole-repo
  dumps and dense-vector RAG. Use when the user asks to "find where this bug
  is in the repo", "localize this failing test / stack trace", "why is this
  error happening — track the root cause", "debug this without loading the
  whole codebase", or "diagnose this regression across the codebase". Produces
  a markdown diagnosis report: ranked suspect file:line, the retrieval path
  taken, a root-cause hypothesis with evidence, and a minimal context bundle
  for a fixer. STOPS at diagnosis. Do NOT use to write or apply the fix or
  open a PR (hand off to a coding skill), to build embeddings / a vector
  index, for general code review or refactoring, or for performance profiling.
---

# Progressive Bug Hunter

## Purpose

Find and explain a bug while reading as little of the codebase as possible:
cheap lexical search first, structural (symbol/call/AST) escalation only when
the cheap tier is insufficient, retrieval ordered so the prompt cache stays
warm across iterations. Output is a diagnosis a downstream fixer can act on —
this skill never edits code.

## When to use this skill

- Use when: "find where this bug is in the repo" / "localize this failing test / stack trace".
- Use when: "why is this error happening — track the root cause" / "diagnose this regression across the codebase".
- Use when: "debug this without loading the whole codebase".
- Do NOT use when: the ask is to write/apply the fix or open a PR (hand to a coding skill), to build an embedding/vector index, to do general code review or refactoring, or to profile performance.

## Inputs

A failure signal: a stack trace, failing-test name/output, error message, or a reproduction. Plus repo access with `grep`/`ripgrep`, `find`, file read, and (if available) a tree-sitter / `list_code_definition_names` tool. No vector index is required or built.

## Workflow

1. **Anchor from the failure.** Extract concrete anchors from the signal: file:line frames, exception type, identifiers, failing assertion. Optionally run `python3 scripts/seed_from_trace.py <trace-file>` to deterministically parse frames and emit a prioritized seed-query plan. Record anchors in the escalation log (`templates/escalation-ladder.md`).
2. **Tier 1 — agentic lexical search (default).** Issue targeted `grep`/`ripgrep` for the anchors (exact identifiers, error strings); list definitions with tree-sitter. Run independent searches in parallel where possible. Prefer this for narrow fixes — it is faster, index-free, and fails *loudly* (an empty result is a clean signal). Read only the spans the hits point to. If the cause is confirmed, go to step 5.
3. **Tier 2 — structural escalation (only if Tier 1 underdetermines).** Follow the symbol / call / import / type graph from the anchor: callers, callees, definition sites, and — for a regression — the change blast-radius via reachability. Pull AST-aligned spans (whole function/class), not fixed line windows. See `references/symbol-graph-and-ast.md`.
4. **Tier 3 — broad cross-file search (only for genuine cross-file unknowns).** When the locus is not yet known, widen with hybrid lexical+structural search and iterate: feed each finding back as the next query (do not one-shot). Budget tokens — stop widening once the working set covers the locus. This skill *uses* search tools; it does not build an embedding index.
5. **Diagnose.** Form a root-cause hypothesis tied to specific evidence (the exact spans). Rank suspect locations as `file:line`. Note contradicting evidence and confidence.
6. **Emit the report** per `templates/diagnosis-report.md`. Stop. Do not modify code.

Throughout, order retrieval to preserve the prompt-cache prefix: keep the
stable prefix (system + tools + repo-map) ahead of the `cache_control`
breakpoint and append volatile grep/read results after it, deterministically
ordered, so each iteration of the hunt stays cheap. See
`references/cache-prefix-ordering.md`. Strategy and escalation criteria:
`references/retrieval-ladder.md`. Evidence/pitfalls: `references/evidence-and-pitfalls.md`.

## Output contract

A single markdown diagnosis report (default `bug-diagnosis.md`, or inline if the caller prefers), shaped by `templates/diagnosis-report.md`, containing:

- **Ranked suspects** — `path/file:line` entries, most-likely first, each with a one-line why.
- **Retrieval path** — the escalation log: each tier, the query issued, hit count, the decision (escalate / stop), and the cache-prefix-preserving order used.
- **Root-cause hypothesis** — the mechanism, tied to quoted evidence spans, with confidence and any contradicting evidence.
- **Minimal context bundle** — the smallest set of files/spans a fixer needs (paths + line ranges), nothing more.

No code is modified; no fix or PR is produced; no index is persisted.

## Constraints

- DO NOT edit code, write the fix, or open a PR — diagnosis only; hand off the bundle.
- DO NOT build or persist an embedding/vector index; this skill is index-free retrieval.
- DO NOT dump whole files or the whole repo into context — escalate tiers, pull AST spans.
- DO NOT one-shot a single broad query — iterate, feeding findings back (RepoCoder-style).
- DO NOT reorder retrieval in a way that breaks the cache prefix (volatile results stay after the breakpoint).
- DO NOT duplicate the retrieval theory here; it lives one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no XML, description under 1024 chars with triggers + negatives.
- [ ] Workflow escalates cheap→structural→broad and explicitly stops before any code edit.
- [ ] Output contract names the four report sections and the no-edit/no-index boundary.

## References

- Escalation ladder, when-grep-beats-vector, iterate-don't-one-shot: `references/retrieval-ladder.md`
- Symbol/call/dependency-graph traversal + AST-aware spans + impact analysis: `references/symbol-graph-and-ast.md`
- Prompt-cache-preserving retrieval ordering: `references/cache-prefix-ordering.md`
- Benchmark evidence, failure modes, anti-extrapolation: `references/evidence-and-pitfalls.md`
- Skeletons: `templates/escalation-ladder.md`, `templates/diagnosis-report.md`
