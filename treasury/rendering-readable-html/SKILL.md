---
name: rendering-readable-html
description: Render structured content the agent already has - tabular data, a markdown or plain-text report, or findings produced this session - into one clean, self-contained, static HTML file made for a human to read offline. Use when the user asks to "make a simple HTML view of this", "turn this into a human-readable HTML page", "visualize this data as a simple HTML report", "render this as a single self-contained HTML file", or "render this plan as a readable HTML page". Produces exactly one .html file with inline CSS only, zero JavaScript, and zero external or network resource requests, so it stays portable, printable, and offline. Do NOT use to build React, production web apps, multi-page sites, or interactive dashboards (use a frontend skill instead); do NOT use to compute domain analyses such as agent-spend or postmortems - this skill only renders content that already exists.
---

# Rendering Readable HTML

Turn content the agent already has into one clean, self-contained, static HTML
page a human can open offline, read comfortably, and print. Render the content
as given — change structure, never invent or drop data.

## When to use

- The user asks to "make a simple HTML view of this", "turn this into a
  human-readable HTML page", "visualize this data as a simple HTML report",
  "render this as a single self-contained HTML file", or "render this plan as a
  readable HTML page".
- There is concrete content to show: a table/dataset, a written report, or
  findings/comparisons/a plan the agent produced this session.

## When NOT to use

- Building a web app, React/component UI, multi-page site, or interactive
  dashboard — defer to a frontend skill.
- Computing a domain analysis (e.g. agent-spend, postmortems). Those skills own
  the data logic; this skill only renders content that already exists.
- Editing an existing HTML file's behavior, or producing anything that needs
  JavaScript, a build step, or external assets.

## Hard rules

1. Output is **exactly one `.html` file**. No companion `.css`/`.js`/asset
   files.
2. **Zero JavaScript**: no `script` element, no `on*` handler, no `javascript:`
   URL.
3. **Zero external resource requests**: no remote stylesheet, font, image,
   media, `@import`, or CSS `url(http...)`. Hyperlinks (`a href`) and `data:`
   URIs are allowed. Full list in `references/self-containment-rules.md`.
4. **Escape every inserted value** per the table in
   `references/layout-patterns.md`. The only HTML you author is structural
   markup; all leaf text is escaped input.
5. **Preserve content**: structure may change, data may not. Never make up,
   shorten away, or quietly drop input.
6. Build on `templates/page.html` — keep its doctype, `utf-8` charset,
   viewport, single inline `<style>`, and `@media print` block.

## Inputs

- The content to render (provided inline, by file path, or as session findings).
- Optional: desired output path/filename, page title, one source line.
  If not given, default the path to `./<kebab-title>.html` and derive a title
  from the content; do not block on these.

## How it works (at a glance)

Take content you already have, wrap it in a fixed HTML shell, prove it loads
with nothing from the network, and hand back one file. The only judgment call
is *what kind* of content it is — a table, a document, or findings like a plan
or comparison. The only gate is the self-containment check: don't report done
until it passes.

```text
content
   │
   ▼
detect input type ──▶ A  tabular
   │              ├──▶ B  document
   │              └──▶ C  findings
   ▼
escape every value
   │
   ▼
assemble from page.html template
   │
   ▼
self-contained check ──no──▶ fix line ─┐
   │ yes                                │
   │            ◀───────────────────────┘
   ▼
report path + pattern
```

## Workflow

1. **Confirm content and destination.** Entry: a render request. Identify the
   content source and resolve the output path (default `./<kebab-title>.html`;
   never overwrite an existing file without confirmation). Exit: known content
   + target path.
2. **Detect input type.** Apply Step 0 of `references/layout-patterns.md` to
   pick one primary pattern: A (tabular), B (document), or C (findings). Mixed
   input → ordered sections, each with its own pattern. Exit: chosen pattern(s).
3. **Build escaped `{{CONTENT}}`.** Construct the semantic HTML for the chosen
   pattern(s), HTML-escaping every value from the input/findings. Exit: a
   content fragment with no raw input in markup.
4. **Assemble the file.** Copy `templates/page.html`; substitute `{{TITLE}}`,
   `{{HEADING}}`, `{{CONTENT}}`, and either fill or delete the `{{META}}` line.
   Set `<html lang>` to the content language. Exit: one `.html` written, no
   `{{PLACEHOLDER}}` token remaining.
5. **Verify self-containment.** Run
   `python3 scripts/check_self_contained.py <output.html>`. If it exits
   non-zero, fix the flagged line and re-run until it exits `0`. Do not skip
   this step. Exit: checker exits `0`.
6. **Report.** State the written file path and the chosen pattern(s). Do not
   open a browser or add anything beyond the single file.

## Output contract

- Exactly one file: `<output.html>` (default `./<kebab-title>.html`).
- Self-contained: passes `scripts/check_self_contained.py` (doctype + `utf-8`
  charset, no `script`, no `on*`/`javascript:`, no remote resource load).
- Readable: single `<h1>`, semantic structure, ~80%-wide column, responsive/
  scrollable tables, print-clean, warm light theme — all from the template,
  unmodified in intent.
- No other files created or modified.

## Constraints & anti-patterns

- DO NOT add JavaScript, a CSS framework, web fonts, a CDN link, or a build
  step "to make it nicer".
- DO NOT paste raw input between tags or into attributes — always escape first.
- DO NOT recompute or re-derive the user's domain analysis; render what exists.
- DO NOT emit more than one file or a `README`/sidecar alongside it.
- MUST ALWAYS keep the template's doctype, charset, viewport, and print block.
- MUST ALWAYS finish on a passing `check_self_contained.py` run.

## Validation checklist

Before reporting done, verify:

- [ ] Exactly one `.html` file written; no companion files.
- [ ] `scripts/check_self_contained.py <output.html>` exited `0`.
- [ ] No `{{PLACEHOLDER}}` token remains.
- [ ] Spot-checked input containing `<`, `&`, or quotes appears as entities.
- [ ] Structure changed, content preserved — nothing invented or dropped.

## References

| Need | File |
|------|------|
| Input-type → semantic HTML mapping, Markdown subset, the escaping table | `references/layout-patterns.md` |
| Exact self-containment prohibitions, required structure, run checklist | `references/self-containment-rules.md` |
| Worked CSV → HTML trace (Pattern A) | `examples/data-table-example.md` |
| Worked plan → HTML trace (Pattern C) | `examples/plan-example.md` |

## Templates and scripts

- `templates/page.html` — the canonical single-file skeleton (do not weaken it).
- `scripts/check_self_contained.py` — automatic pass/fail check; the skill
  runs it on its own output before reporting completion.
