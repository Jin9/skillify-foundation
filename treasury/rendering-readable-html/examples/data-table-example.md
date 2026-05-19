# Example: CSV → single self-contained HTML (Pattern A)

A compact end-to-end trace. Shows the `{{CONTENT}}` fragment only, not a full
generated file (the skill produces the file; this anchors expected structure).

## Input the user provided

```
service,calls,error_rate
auth-api,18420,0.4%
ledger & "core",9310,1.1%
notify,540,0.0%
```

## Decisions

- Input type: rows/columns → **Pattern A (tabular data)**.
- `calls` is entirely numeric → header + cells get `class="num"`.
- `error_rate` is a percentage string → treated as text, left-aligned.
- `ledger & "core"` contains `&` and `"` → must be escaped.
- Output path: user gave none → default `./service-summary.html`.

## Expected `{{CONTENT}}` (escaped, semantic)

```html
<div class="table-wrap">
  <table>
    <caption>Service call summary (3 rows)</caption>
    <thead>
      <tr>
        <th scope="col">service</th>
        <th scope="col" class="num">calls</th>
        <th scope="col">error_rate</th>
      </tr>
    </thead>
    <tbody>
      <tr><td>auth-api</td><td class="num">18420</td><td>0.4%</td></tr>
      <tr><td>ledger &amp; &quot;core&quot;</td><td class="num">9310</td><td>1.1%</td></tr>
      <tr><td>notify</td><td class="num">540</td><td>0.0%</td></tr>
    </tbody>
  </table>
</div>
```

This fragment replaces `{{CONTENT}}` in `templates/page.html`;
`{{TITLE}}`/`{{HEADING}}` become `Service summary`, the `{{META}}` line is
removed (no provenance given).

## Verification

`python3 scripts/check_self_contained.py service-summary.html` exits `0`:
no script, no remote resource, doctype + utf-8 present. Note `&` rendered as
`&amp;` and `"` as `&quot;` — raw input never reaches the markup.
