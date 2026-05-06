# Anti-pattern note format

When a prompt loses badly enough that the squad wants the next researcher
to avoid it, the loss lands as a one-liner in
`prompts/library/_anti-patterns.md`. No full entry. No why-it-worked
paragraph. The point is to keep losses cheap.

## Format

```
- <ISO date> · <stage>/<topic-slug> — <one-line why it failed>
```

## Examples

```
- 2026-04-19 · plan/jwt-validation — Gemini drifted into a 1200-word essay; planning needs Codex.
- 2026-04-26 · critique/migration-safety — Claude refused on safety grounds; over-broad framing triggered the refusal.
- 2026-05-03 · research/loan-domain — prompt asked for "everything" → 4 files of low-signal output for $2.10.
- 2026-05-06 · implement/handler-rewrite — un-sandboxed run leaked an internal hostname into a comment.
```

## Rules

- Lines are append-only. Do not edit historic entries; add a new line if
  context changes.
- Slug must match an existing or once-attempted library entry. New slugs
  create confusion.
- Reason is concrete, not editorial — "Gemini drifted into essay"
  beats "didn't work".
- No real PII, real customer names, or real credentials in the reason.

## Promotion to a full entry

A loss can become a win after rewriting. When that happens, write the
full library entry with a fresh `last_validated` date and **leave** the
historic anti-pattern line in place. The audit trail matters more than
the symmetry.
