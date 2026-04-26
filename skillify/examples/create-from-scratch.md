# Create From Scratch Walkthrough

This example shows the expected trajectory for Create mode. It is a worked trace, not a template to load for every task.

## User Request

> Create a skill for generating release notes from merged pull requests.

## Step 0: Mode Detection

- Mode: Create
- Output contract: new skill folder named `generating-release-notes`
- Required examples: collect at least 3 trigger prompts before naming and drafting

## Concrete Trigger Prompts

1. "Create release notes for version 2.4.0 from the merged PRs."
2. "Summarize this milestone into customer-facing release notes."
3. "Generate changelog bullets from these GitHub pull requests."

## Boundary

- Owns: turning merged PR or milestone data into release-note drafts
- Does not own: publishing releases, tagging versions, deploying builds, or changing Git history
- Degree of freedom: Medium, because release-note style varies but the output structure is predictable

## Planned Resources

```text
generating-release-notes/
├── SKILL.md
├── references/
│   └── release-note-style.md
└── templates/
    └── release-notes-template.md
```

## Draft Frontmatter

```yaml
---
name: generating-release-notes
description: >
  Generates release notes and changelog drafts from merged pull requests,
  milestones, or change summaries. Use when the user asks to "create release
  notes", "summarize this milestone", "generate changelog bullets", or prepare
  customer-facing release notes. Do NOT use for publishing releases, tagging
  versions, or deploying builds.
---
```

## Draft Workflow

1. Read the provided PRs, milestone, or change summaries.
2. Classify changes as Features, Fixes, Breaking Changes, Deprecations, Documentation, or Internal.
3. Rewrite technical details into user-facing language.
4. Preserve links to PRs or issues when provided.
5. Output release notes using `templates/release-notes-template.md`.
6. Validate that every included claim traces to a source item.

## Validation

```text
scripts/quick_validate.py generating-release-notes
scripts/check_links.py generating-release-notes
```

Rubric expectation:

- Trigger Quality: 5/5, because exact user phrases are present.
- Scope Focus: 5/5, because publishing and deployment are excluded.
- Workflow Clarity: 4/5 or higher, depending on final output details.
- Output Contract: 4/5 or higher once the template path and sections are fixed.

## Final Delivery Note

Report the created files, validation results, and any open assumptions about release-note tone or target audience.
