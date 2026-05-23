# Self-containment rules

The output is ONE `.html` file that opens correctly with no network, no build
step, and no companion files — on any machine, offline, and when printed.

## Hard prohibitions (the file MUST NOT contain)

1. **Any `script` element** — inline or external. Zero JavaScript. No
   `<script>` data islands either; keep the promise simple.
2. **Any inline event-handler attribute** (`onclick`, `onload`, `on*`) — that
   is JavaScript.
3. **Any `javascript:` URL** anywhere (`href`, `src`, CSS).
4. **Any external resource request**: a remote (`http://`, `https://`, or
   protocol-relative `//`) URL used to *load* a resource — stylesheet, script,
   image, `srcset` candidate, web font, `@import`, CSS `url(...)`, `iframe`,
   `embed`, `object`, `audio`/`video`/`source`/`track`, or SVG `use`/`image`.
5. **External or `@import`-ed CSS / web fonts.** Styling is the one inline
   `<style>` block plus optional `style="…"` attributes. System font stacks
   only.
6. **Companion files.** No sidecar `.css`, `.js`, or asset files. One file.

## Explicitly allowed

- **Hyperlinks**: `<a href="https://…">` and `mailto:` — these navigate on
  click, they do not fetch on load, so they keep the file self-contained.
- **`data:` URIs** for an essential embedded image.
- Relative/`#` fragment links within the page.
- CSS-only theming and interactivity-free niceties already in the template
  (`color-scheme`, `@media print`, `:focus-visible`, scrollable table
  wrapper).

## Required structure

- Starts with `<!doctype html>` (case-insensitive).
- Has `<meta charset="utf-8">`.
- Has `<meta name="viewport" content="width=device-width, initial-scale=1">`.
- `<html lang="…">` set (default `en`; change to match content language).
- A `<title>` and a single `<h1>`.
- Print-readable (the template's `@media print` block satisfies this; keep it).

## Verification checklist (the skill runs this before reporting done)

- [ ] Exactly one file was written, extension `.html`.
- [ ] `python3 scripts/check_self_contained.py <output.html>` exits `0`.
- [ ] Every interpolated value is HTML-escaped (spot-check input that contains
      `<`, `&`, or quotes — it must appear as entities, not raw markup).
- [ ] No `{{PLACEHOLDER}}` token remains in the file.
- [ ] No data was invented, summarised away, or silently dropped versus the
      input; structure changed, content preserved.
- [ ] Opens correctly from `file://` with no console errors and no missing
      resources (the checker proves this deterministically).

If `check_self_contained.py` exits non-zero, fix the flagged line and re-run
until it passes; do not report the task complete on a failing check.
