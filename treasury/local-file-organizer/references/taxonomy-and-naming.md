# Folder Taxonomy and Naming Rules

Default scheme. Override only when the user provides another convention.

## Default folder structure

```text
<target-root>/
  00_INBOX/
  10_PROJECTS/
  20_FINANCE/
  30_LEGAL/
  40_PERSONAL/
  50_TECH/
  60_REFERENCES/
  70_ARCHIVE/
  _REVIEW/
  _DUPLICATES/
```

## Naming the leaves

Prefer business meaning over file type:

```text
10_PROJECTS/client-a/contracts/
10_PROJECTS/client-a/requirements/
20_FINANCE/tax/2026/
50_TECH/scripts/
60_REFERENCES/articles/
```

Avoid extension-only buckets (`PDF/`, `Excel/`, `Images/`) at the top level. They are acceptable only inside `_REVIEW/` or `00_INBOX/` as a fallback.

## Confidence routing

| Confidence | Destination |
|------------|-------------|
| 0.90 – 1.00 | Final business folder |
| 0.80 – 0.89 | Final business folder, reason mentioned in report |
| 0.60 – 0.79 | `_REVIEW/<topic>/` |
| < 0.60      | Leave in place unless the user overrides |

## Renaming rules

- Preserve original filename by default.
- Rename only when it improves clarity or prevents conflict.
- Never strip identifiers: invoice number, project name, date, version.
- Use safe characters for macOS/Linux paths. Avoid characters that need shell quoting.
- Use lowercase or kebab-case only when the user explicitly requests normalization.

## Conflict handling

When two files map to the same target, keep both with a duplicate suffix:

```text
original.pdf
original__duplicate-001.pdf
original__duplicate-002.pdf
```

Never overwrite. If a target already exists at apply time, route the source to `_DUPLICATES/` with the same suffix scheme.
