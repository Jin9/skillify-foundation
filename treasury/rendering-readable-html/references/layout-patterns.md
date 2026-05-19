# Layout patterns: input type → semantic HTML

The single source of truth for how to shape `{{CONTENT}}` in
`templates/page.html`. The escaping rule below is canonical — do not restate it
elsewhere; point here.

## The escaping rule (always)

Every value that originates from the input or from session findings MUST be
HTML-escaped before it enters the document, including text inside attributes:

| Char | Replace with |
|------|--------------|
| `&`  | `&amp;`  |
| `<`  | `&lt;`   |
| `>`  | `&gt;`   |
| `"`  | `&quot;` (attribute values) |
| `'`  | `&#39;` (attribute values) |

Escape `&` first. Never paste raw input between tags or into an attribute. The
only HTML you author is the structural markup in the patterns below; all leaf
text is escaped input.

## Step 0 — detect the input type

Pick exactly one primary layout. If the input mixes types, render sections in
input order, each with its own pattern.

| Signal | Input type | Pattern |
|--------|------------|---------|
| Rows/columns: CSV, TSV, JSON array of objects, Markdown table, a result set | **Tabular data** | A |
| Prose with headings/lists/code: a report, summary, analysis, README-like text | **Document** | B |
| Discrete labelled items the agent produced this session: findings, comparisons, a plan, pros/cons | **Findings** | C |

## Pattern A — tabular data

1. Derive columns from the header row or object keys (first object's keys; union
   if keys vary). Preserve source order.
2. Build one `<table>`:
   - `<caption>` = short, escaped description of what the table is.
   - `<thead>` one row of `<th scope="col">`; add `class="num"` to a header
     whose column is entirely numeric.
   - `<tbody>` one `<tr>` per record; numeric cells get `<td class="num">`.
   - Empty/missing cell → empty `<td>` (do not invent values).
3. Wrap the table: `<div class="table-wrap">…</div>` (keeps wide tables
   scrollable instead of overflowing on mobile / print).
4. Optional summary above the table when it adds signal — a `<dl class="summary">`
   with row count, column count, and any obvious totals you can compute from the
   given data only. Do not infer beyond the data.

## Pattern B — document

Convert this Markdown subset only. Anything outside it: render as an escaped
`<p>` (or `<pre>` if it is clearly preformatted). Do not invent structure.

| Markdown | HTML |
|----------|------|
| `# … / ## … / ### …` | `<h1>` reserved for the page `{{HEADING}}`; map document `#`→`<h2>`, `##`→`<h3>`, deeper→`<h4>` (never duplicate the page h1) |
| Paragraph (blank-line separated) | `<p>` |
| `- ` / `* ` / `1. ` lists | `<ul>` / `<ol>` with `<li>` |
| `**bold**` / `*italic*` | `<strong>` / `<em>` |
| `` `code` `` | `<code>` |
| fenced ```` ``` ```` block | `<pre><code>` (escape contents; keep newlines) |
| `> quote` | `<blockquote>` |
| `[text](url)` | `<a href="ESCAPED_URL">text</a>` — only `http`, `https`, `mailto`, or relative/`#` targets; drop `javascript:` and other schemes |
| `---` | `<hr>` |
| `\| a \| b \|` table | Pattern A |

Hyperlinks are allowed: an `<a href>` is navigation, not a load-time resource
request, so it does not break offline self-containment. Resource references
(images, CSS, JS, fonts, media) are not — see
`references/self-containment-rules.md`.

## Pattern C — findings

For each item the agent is presenting:

```
<section class="card">
  <h2>ESCAPED item title</h2>
  <p>ESCAPED summary line</p>
  <ul>
    <li>ESCAPED detail</li>
  </ul>
</section>
```

Use a top `<p>` for any framing sentence. For pros/cons or A-vs-B comparisons,
use Pattern A (two columns) instead of prose. Keep one concept per card.

## Images and media

Default: omit. If an image is essential and small, embed it as a `data:` URI
(`<img src="data:image/png;base64,…" alt="ESCAPED">`). Never reference a remote
image, stylesheet, font, script, or media URL — that breaks the offline,
single-file guarantee enforced by `scripts/check_self_contained.py`.

## Title, heading, meta

- `{{TITLE}}` (in `head`): a concise plain-text page title (escaped).
- `{{HEADING}}` (the page `<h1>`): the human-facing report title (escaped).
- `{{META}}`: one optional escaped provenance line (e.g. source filename + date).
  If unused, remove the entire `<p class="meta">` line from the template.
