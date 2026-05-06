# Anti-pattern line

Single line, append-only to `prompts/library/_anti-patterns.md`:

```
- <YYYY-MM-DD> · <stage>/<topic-slug> — <one-line why it failed>
```

## Example

```
- 2026-05-06 · plan/jwt-validation — Gemini drifted to alternatives instead of writing the plan; use Codex.
```

## Rules

- One line. No headings, no sub-bullets, no code blocks.
- Concrete reason. Editorial reasons ("didn't work", "low quality") fail
  the next researcher who needs the lesson.
- No real PII / credentials / customer names.
- Slug should match an existing or attempted library entry; if not,
  prefer creating the entry as a `loss` with a stub body rather than
  leaving an orphan line.
