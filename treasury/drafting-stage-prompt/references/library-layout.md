# Library layout

```
prompts/library/
├── README.md                       # squad-maintained index (optional)
├── _anti-patterns.md               # one-line loss notes (created on first loss)
├── research/
│   ├── auth-systems.md
│   └── kafka-consumer-pitfalls.md
├── plan/
│   ├── jwt-validation.md
│   └── schema-migration.md
├── critique/
│   └── ddd-aggregate-boundaries.md
├── implement/
│   └── go-handler-rewrite.md
├── review/
│   └── concurrency-walkthrough.md
└── test/
    └── bats-locking.md
```

## Naming rules

- Stage directory matches a built-in scaffold stage or a stage declared
  in a profile.
- Filename matches `^[a-z0-9-]{1,40}\.md$` (topic slug only — no spaces,
  no caps, no underscores, no claude/anthropic).
- Update entries stay at the same path; new versions append `-v2`,
  `-v3`, never `-v1.5`.

## Index file (`prompts/library/README.md`)

Optional but recommended. The squad keeps it as a flat index of titles
to filenames. This skill does not maintain the README — keep it human-
maintained at the Friday meeting.

## What lives outside `library/`

- `*_PROMPT_PREFIX` env vars in `profiles/<name>.sh` — these are short
  prefixes prepended at runtime, not full prompts. Edit profiles via
  `authoring-scaffold-profile`.
- Per-workflow prompts the user drafts on the fly — these stay
  workflow-local and are only promoted to the library when they win at
  the Friday review.

## Git hygiene

- Library entries are committed. PII / credential leaks are a security
  incident.
- `prompts/library/private/` is gitignored — use it for entries that
  reference squad-internal context (customer names redacted, IDs
  rotated). This skill should refuse to write there unless the user
  explicitly says "private".
