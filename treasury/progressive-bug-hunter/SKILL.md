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
  taken, a root-cause hypothesis with evidence, a P0–P4 severity rating, and
  a minimal context bundle for a fixer. STOPS at diagnosis. Do NOT use to write or apply the fix or
  open a PR (hand off to a coding skill), to build embeddings / a vector
  index, for general code review or refactoring, or for performance profiling.
---

# Progressive Bug Hunter

## Purpose

Find and explain a bug while reading as little of the codebase as possible.

Search the cheap way first — plain-text search (lexical search, e.g. `grep`) for the
exact names and error strings from the failure. Only move up to structural search
(following the symbol / call / AST graph — that is, which code defines or calls which)
when the cheap search isn't enough. Order every search so the unchanging part of the
prompt gets reused from one try to the next (the "prompt cache" stays warm), which keeps
each step fast and cheap.

The output is a diagnosis a later fixer can act on. This skill never edits code itself.

## When to use this skill

- Use when: "find where this bug is in the repo" / "localize this failing test / stack trace".
- Use when: "why is this error happening — track the root cause" / "diagnose this regression across the codebase".
- Use when: "debug this without loading the whole codebase".
- Do NOT use when: the ask is to write/apply the fix or open a PR (hand to a coding skill), to build an embedding/vector index, to do general code review or refactoring, or to profile performance.

## Inputs

A **failure signal** — proof that something is broken: a stack trace, a failing test's name
or output, an error message, or steps that reproduce the problem. Plus access to the repo
with `grep`/`ripgrep` (plain-text search), `find`, the ability to read files, and — if you
have one — a tool that lists code definitions (e.g. tree-sitter / `list_code_definition_names`).
No vector index (a pre-built semantic-search database) is needed or created.

## Workflow

1. **Anchor from the failure.** Pull the concrete clues ("anchors") out of the failure signal:
   `file:line` frames from the stack trace, the exception type, identifier names (functions,
   variables, classes), and the assertion that failed. Optionally run
   `python3 scripts/seed_from_trace.py <trace-file>` to parse the trace the same way every time
   and produce a ranked list of searches to try first. Write the anchors into the escalation log
   (`templates/escalation-ladder.md`) — a scratchpad that tracks the hunt.
2. **Tier 1 — plain-text search (the default, and the cheapest).** Run focused `grep`/`ripgrep`
   for the anchors (exact identifier names, exact error strings), and list code definitions with
   tree-sitter. Run searches that don't depend on each other at the same time (in parallel) when
   you can. Start here for small bugs: it's fast, needs no pre-built index, and fails *loudly* —
   an empty result is a clear "not here" answer, not a wrong guess. Read only the small spans the
   hits point to, nothing more. If this confirms the cause, jump to step 5.
3. **Tier 2 — follow how the code is wired (only if Tier 1 found *where* but not *why*).** Walk
   the symbol / call / import / type graph outward from the anchor — that is, follow who-calls-what:
   callers, callees, where things are defined, and (for a regression) the "blast-radius": which
   other code a change can reach and break. Grab whole units of code — a complete function or class
   ("AST-aligned spans") — not a fixed number of lines that might cut a function in half. See
   `references/symbol-graph-and-ast.md`.
4. **Tier 3 — wide search across many files (only when the location is still unknown).** If you
   still don't know where the bug lives, widen out with a mix of plain-text and structural search,
   and *repeat in a loop*: take what each search finds and feed it back in as the next search —
   don't try to nail it in one big query. Watch your token budget: stop widening once the set of
   files you're holding covers the bug's location. This skill *uses* search tools; it never builds
   a semantic/embedding index.
5. **Diagnose.** Write a root-cause hypothesis (your best explanation of *why* it breaks), tied to
   specific evidence — the exact code spans you read. Rank the suspect locations as `file:line`,
   most likely first. Note any evidence that argues against your hypothesis, and how confident you
   are. Give the bug **one** P0–P4 severity, based only on what the code/trace actually shows (does
   it crash? lose data? is a security-sensitive path reachable? how wide is the blast-radius?).
   Keep this severity separate from how confident you are in the diagnosis — see the bands below
   and `references/severity-rubric.md`.
6. **Write the report** using `templates/diagnosis-report.md`. Then stop. Do not change any code.

Throughout, order your searches so the prompt cache stays warm (re-runs stay cheap). In plain
terms: put the parts that never change (the system text, the tool definitions, the repo map)
*first*, before the `cache_control` cache marker, and add the changing search/read results *after*
it, always in the same order. The unchanging front part then gets reused instead of re-billed each
round. See `references/cache-prefix-ordering.md`. For the full strategy and the rules on when to
climb to the next tier: `references/retrieval-ladder.md`. For the evidence behind this order and
the traps to avoid: `references/evidence-and-pitfalls.md`.

## Severity bands (P0–P4)

Rate the bug on this scale, where P0 is the worst. Full definitions and how to pick a band live
in `references/severity-rubric.md`:

- **P0 — Critical** — the system goes down, data is lost or corrupted, or there's a security break
  on a main path the app actually runs. Must be fixed before release.
- **P1 — High** — a core feature is broken with no way around it, and it affects a lot of other
  code (wide blast-radius).
- **P2 — Medium** — something works badly but there's a workaround; the damage stays in a small
  area (bounded blast-radius).
- **P3 — Low** — a minor or edge-case problem (only happens in unusual situations) that affects
  very little.
- **P4 — Trivial** — cosmetic, or barely matters.

## Output contract

One Markdown diagnosis report (by default a file named `bug-diagnosis.md`, or written straight
into the reply if the caller prefers), following `templates/diagnosis-report.md`. It contains:

- **Severity** — the one P0–P4 band for the bug, plus a one-line note on its impact. This is a
  separate thing from how confident you are about each suspect.
- **Ranked suspects** — `path/file:line` locations, most-likely first, each with a one-line reason.
- **Retrieval path** — the search trail (the escalation log): for each tier, what you searched for,
  how many hits, the decision (climb to the next tier / stop), and the cache-friendly order you kept.
- **Root-cause hypothesis** — how the bug actually happens, backed by quoted snippets of the code
  you read, with your confidence level and any evidence that points the other way.
- **Minimal context bundle** — the smallest set of files and line ranges a fixer needs, and nothing
  extra.

No code is changed; no fix or PR is created; no search index is saved.

## Constraints

- DO NOT edit code, write the fix, or open a PR — this skill only diagnoses; hand the bundle off
  to a fixer.
- DO NOT build or save an embedding/vector index (a semantic-search database) — this skill
  searches without one.
- DO NOT dump whole files or the whole repo into context — climb the tiers and pull only the
  whole-function/class spans (AST spans) you need.
- DO NOT try to find it in one big search — repeat in a loop, feeding each finding back in as the
  next search (the "RepoCoder" pattern).
- DO NOT reorder your searches in a way that breaks the cache prefix — the changing results must
  stay *after* the cache marker.
- DO NOT copy the retrieval theory into this file — it lives one level down in `references/`.
- DO NOT rate the severity higher than the code/trace evidence supports. If you'd need production
  or business context to settle the band, say so instead of guessing.

## Validation

- [ ] Frontmatter `name` matches the folder name, is kebab-case, has no angle brackets, and the
  description is under 1024 characters with trigger phrases plus "do NOT use" cases.
- [ ] The workflow goes cheap → structural → broad in that order and clearly stops before editing
  any code.
- [ ] The output contract lists the report sections (severity plus the four core sections) and
  states the no-edit / no-index boundary.
- [ ] The report gives exactly one P0–P4 severity, grounded in evidence, kept separate from confidence.

## References

- When to climb each tier, why `grep` often beats vector search, and why to loop instead of
  one-shot: `references/retrieval-ladder.md`
- Following the symbol / call / dependency graph, pulling whole-unit (AST) spans, and working out
  a change's blast-radius: `references/symbol-graph-and-ast.md`
- Ordering searches so the prompt cache stays warm: `references/cache-prefix-ordering.md`
- The evidence behind this approach, the ways it can fail, and not over-claiming from it:
  `references/evidence-and-pitfalls.md`
- The P0–P4 bands, how to pick one, and severity vs. confidence: `references/severity-rubric.md`
- Fill-in skeletons: `templates/escalation-ladder.md`, `templates/diagnosis-report.md`
