# Squad-researcher workflow context (context-only)

This is **explanatory background only**, so the librarian can interpret reports
and validate them against the gates that produced them. It does **not** authorize
designing, adding, or modifying pipeline stages. The authoritative in-vault
record is the `## Provenance & verification` section of `<vault>/index.md`;
re-read it for live version/count facts before quoting numbers.

## Boundary — the workflow lives outside this vault

The squad-researcher workflow itself — its `skills/*/SKILL.md` stage
definitions, plans, and `tmp/runs/` working dirs — lives in the external Claude
Code project / skills directory, **not** in this vault. This vault holds only
the *published final reports*, the `maps/` MOCs, the `index.md` catalog, and the
provenance summary. Any `skills/…`, `plans/…`, or `tmp/runs/…` path is relative
to that external workspace. If asked to change the pipeline, decline (Mode D is
context-only) and point the user to that external workspace.

## The 6-stage pipeline (v0.7.1)

A report is the output of a 6-stage pipeline, run at `depth=deep`:

1. **plan-research** — scope the topic and plan the search.
2. **search-sources** — gather candidate sources.
3. **extract-findings** — pull findings with source/citation IDs.
4. **synthesize-report** — write the report from findings.
5. **review-report** — quality/verification pass.
6. **report-run** — finalize and publish the run.

The 172 reports were produced across 6 batch phases (May 2026).

## Per-stage re-verification gates

Every run passed an independent re-verification. A report you inspect should
satisfy all of these (and a librarian audit can spot-check them):

- JSON validity.
- Schema conformance.
- Source / finding / citation-id integrity.
- Prose↔cited **bijection** (every `[n]` resolves to a listed source; every
  listed source is cited).
- Closed-world citations (uncited collected sources are dropped — hence the
  legitimate gaps in `## Sources` numbering).
- The **deep word band** (the gate evaluates raw `wc -w` of the final report).
- **No `[n]` citations inside fenced code / diagram blocks** (diagrams are plain
  text only).

These map directly to the citation-discipline gates in the vault `CLAUDE.md`;
when auditing a report, that `CLAUDE.md` remains the source of truth for the
exact thresholds.

## Frontier LLM Pros / Cons batch

A domain-3 subsection: 8 deep runs (8/8 PASS), ~4,010 average words, generated
2026-05-16 (UTC). Subject = Anthropic Claude Opus 4.7 vs OpenAI GPT-5.5
(per-run tier resolution) vs Google Gemini 3.1 Pro. Slugs are `frontier-llm-*`,
listed `1.`–`8.` under `### Frontier LLM Pros / Cons (8 runs)`.

## Open workflow-backlog items

`index.md` records three open improvements to the workflow itself. The librarian
can *explain* these and note where a report reflects (or is limited by) them,
but cannot implement them here:

1. **Per-stage reasoning effort** (low/med/high/xHigh/max) — not possible from
   inside `Agent` today; needs a runtime that calls the Anthropic API with
   `thinking.budget_tokens`.
2. **`recommended_model:` advisory field** in each stage `SKILL.md` frontmatter,
   cross-checked against the workflow YAML.
3. **Depth-gated stages** — e.g., skip `review-report` when `depth=quick`.

If the user wants to act on any of these, that work happens in the external
workflow workspace, not in this vault and not through this skill.
