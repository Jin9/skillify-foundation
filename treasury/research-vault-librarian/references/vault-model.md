# Vault operating procedure (librarian)

The vault's own `CLAUDE.md` is the **source of truth** for the cross-file
invariant, slug conventions, report anatomy, and citation gates. This file does
**not** restate those rules — it adds only the *operating procedure* for
querying the corpus and assembling a consistent intake proposal. Whenever a rule
is needed, re-read the live `<vault>/CLAUDE.md`.

## Graph shape

`index.md` → 8 `maps/<domain>.md` MOCs → report fans. A report's wikilink lives
in **exactly one** MOC (its own domain). `index.md` links the 8 MOCs and carries
the summary table + provenance; it does not list individual reports.

## Query traversal (Mode A)

1. **`index.md`** — read the domain summary table and the `## Domain maps`
   section. Map the user's topic to one or a few of the 8 domains.
2. **`maps/<domain>.md`** — each MOC opens with
   `# <Domain Title>` then `> Domain map (N reports) · part of [[index|Research Vault — Index]]`.
   Bullets are `- [[reports/<domain>/<slug>|Display]]`, optionally suffixed
   ` 🎯`. Scan bullet **Display** names to shortlist; do not open every report.
3. **Domain 3 only** — `maps/llm-model-selection-frontier-comparison.md` splits
   into `### Model Selection Study (65 topics)` (slugs prefixed `llm-`) and
   `### Frontier LLM Pros / Cons (8 runs)` (slugs prefixed `frontier-llm-`,
   numbered `1.`–`8.`). Use the subsection to disambiguate study vs. cross-vendor
   batch.
4. **`reports/<domain>/<slug>.md`** — open only the shortlisted files. Report
   structure: `# Title` → optional `**Depth:** … | **Audience:** … | **Date:** …`
   → `## Executive Summary` → `## Background` → `## Methodology` →
   `## Key Findings` → topic sections → `## Synthesis` →
   `## Limitations and Open Questions` → `## Sources`.
5. **Dashboard** — for LLM pros/cons or model-recommendation questions, also use
   the generated view `dashboard/llm-pros-cons.html` and
   `dashboard/_data/frontier.json` / `dashboard/_data/recommendations.json`.
   This is a regenerable synthesis, **not** part of the report corpus or counts;
   when it disagrees with a report, the report wins.

### The `🎯` FOCUS marker

` 🎯` on a MOC bullet flags a FOCUS_BACKLOG topic (agentic delivery squad). When
a query could be answered by several reports, prefer FOCUS-marked ones and say
they are current-focus. Never strip or move a `🎯` when proposing a patch.

### Quoting `## Sources` — closed-world numbering

`## Sources` lists only sources actually cited in the prose, so numbering
legitimately has gaps (`1, 2, 3, 5, 6 …`). When you quote or cite a source,
reproduce its existing number **verbatim**; never renumber to "close" a gap and
never infer that a gap is an error. (CLAUDE.md is the authority on why.)

## Intake / revision proposal (Mode C)

You only *propose*. Deliver one bundle the user applies. Before assembling it,
re-read `<vault>/CLAUDE.md` for the authoritative invariant, slug rules, and
citation gates, then satisfy this operational checklist:

1. **Domain placement** — choose the single domain folder + its matching MOC.
   A report appears in its own domain's MOC only. Moving a report between domain
   folders changes its wikilink and the MOC it lives in — flag this explicitly.
2. **Slug** — keep slug filename stable; `# Title` and the curated MOC Display
   name may legitimately differ from the slug and from each other. Do not
   propose renames to "align" them.
3. **Report file** — full `reports/<domain>/<slug>.md` content matching the
   report anatomy above and the CLAUDE.md citation gates: no YAML frontmatter;
   bracketed-integer inline citations; closed-world prose↔Sources bijection;
   `n. Title — Publisher — URL` source lines; **no `[n]` inside fenced/diagram
   blocks**; within the deep word band CLAUDE.md specifies.
4. **MOC bullet** — exact line, placed in the correct MOC (and correct `###`
   subsection for domain 3):

   ```
   - [[reports/<domain>/<slug>|Display Name]]
   ```

   Append ` 🎯` only if it is a FOCUS topic. Preserve the MOC's
   `> Domain map (N reports) · …` header — its `N` must change with the count.
5. **Catalog deltas** — the exact `index.md` summary-table row delta for that
   domain and the `**Total**` delta; the MOC header `N` delta. (For domain 3 the
   table note "Model Selection Study 65 + Frontier Pros/Cons 8" also moves.)
6. **Expected post-state** — state what `scripts/check_vault_consistency.sh`
   should print after the user applies the bundle: every check PASS and the new
   `files == links == header-total`.

Removal is the same checklist in reverse (delete file + its single MOC bullet +
decrement the MOC header `N`, the table row, and the total). Never propose a
report file change without its matching MOC + catalog deltas in the same bundle.
