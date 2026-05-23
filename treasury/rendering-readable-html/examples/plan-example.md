# Example: a plan → single self-contained HTML (Pattern C)

A compact end-to-end trace. Shows the `{{CONTENT}}` fragment only, not a full
generated file (the skill produces the file; this anchors expected structure).
The point of this example: a dense, LLM-style plan becomes a *readable* page by
changing only structure and styling — every word is preserved.

## Input the user provided

```
Plan: harden the auth service before Q3

1. Replace session cookies with JWT.
   - Migrate 12 endpoints across auth & billing.
   - Risk: encoded token must stay <8KB (legacy gateway cap).
2. Add a rate-limit middleware.
   - 100 req/min per IP; respond 429 with a Retry-After header.
   - Ship behind the rl_v2 feature flag, default off.
3. Rotate signing keys from HS256 -> RS256.
   - Publish JWKS at /.well-known/jwks.json.
   - Old keys valid for a 14-day overlap, then revoked.
```

## Decisions

- Input type: discrete labelled steps the agent produced this session →
  **Pattern C (findings)**.
- One numbered step → one `<section class="card">`; the step heading becomes its
  `<h2>`, its sub-bullets become a `<ul>`. The lead line becomes the page
  `{{HEADING}}`, so no extra framing `<p>` is needed.
- **Content preserved verbatim.** Every step and every sub-bullet is rendered
  as-is — "readable" means cleaner layout, not shorter text. Rewording or
  dropping any line would violate hard rule #5 ("never summarise away").
- Escaping (the only chars here that need it): `auth & billing` → `&amp;`;
  `<8KB` → `&lt;`; the `>` in `HS256 -> RS256` → `&gt;`.
- Output path: user gave none → default `./auth-hardening-plan.html`.

## Expected `{{CONTENT}}` (escaped, semantic)

```html
<section class="card">
  <h2>1. Replace session cookies with JWT.</h2>
  <ul>
    <li>Migrate 12 endpoints across auth &amp; billing.</li>
    <li>Risk: encoded token must stay &lt;8KB (legacy gateway cap).</li>
  </ul>
</section>

<section class="card">
  <h2>2. Add a rate-limit middleware.</h2>
  <ul>
    <li>100 req/min per IP; respond 429 with a Retry-After header.</li>
    <li>Ship behind the rl_v2 feature flag, default off.</li>
  </ul>
</section>

<section class="card">
  <h2>3. Rotate signing keys from HS256 -&gt; RS256.</h2>
  <ul>
    <li>Publish JWKS at /.well-known/jwks.json.</li>
    <li>Old keys valid for a 14-day overlap, then revoked.</li>
  </ul>
</section>
```

This fragment replaces `{{CONTENT}}` in `templates/page.html`;
`{{TITLE}}`/`{{HEADING}}` become `Plan: harden the auth service before Q3`, the
`{{META}}` line is removed (no provenance given).

## Verification

`python3 scripts/check_self_contained.py auth-hardening-plan.html` exits `0`:
no script, no remote resource (the only URL, `/.well-known/jwks.json`, is a
relative path), doctype + utf-8 present. Note `&` rendered as `&amp;`, `<` as
`&lt;`, and `>` as `&gt;` — raw input never reaches the markup, and all three
steps with every sub-bullet survive intact.
