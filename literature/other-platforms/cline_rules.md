# Cline Rules

Source: https://docs.cline.bot/features/cline-rules/overview  
Accessed: 2026-05-14  
Category: Other platforms / persistent AI rules

## Why This Source Matters

Cline recognizes several rule formats and supports conditional activation. This makes it a useful comparison point for skills because it separates always-on project guidance from path-triggered guidance.

## Core Data Points

- Rules are Markdown or text files that provide persistent instructions across conversations.
- Primary workspace rules live in `.clinerules/`.
- Global rules live in the user's Cline rules directory, commonly under `~/Documents/Cline/Rules`.
- Cline recognizes multiple rule sources:
  - `.clinerules/`
  - `.cursorrules`
  - `.windsurfrules`
  - `AGENTS.md`
- All detected rule types are visible in the Rules panel and can be toggled.
- Workspace rules and global rules can both apply; workspace rules take precedence when conflicts occur.
- Rule files should be focused on a single concern.
- Conditional rules can use YAML frontmatter with `paths`.
- Path-based rules activate when matching files are in the user's message, open tabs, visible editor panes, edited files, or pending operations.
- Invalid YAML can fail open, making raw content visible to help with debugging.

## Design Implications For Skillify

- A path-conditional rule is a good target for narrow conventions, such as frontend, backend, tests, or docs.
- Skillify should not generate one monolithic `.clinerules` file. Multiple focused files are easier to toggle, debug, and maintain.
- Cline's support for `AGENTS.md` reinforces `AGENTS.md` as a portable base layer.
- Conditional activation based on file paths is a rule-system equivalent of progressive disclosure.
- Path triggers should be explicit and tested; vague prompts may not activate the intended rule.

## Practical Patterns

- `.clinerules/frontend.md` for UI conventions.
- `.clinerules/backend.md` for service/API conventions.
- `.clinerules/testing.md` with `paths` for `**/*.test.ts`, `**/*.spec.ts`, and test folders.
- `.clinerules/docs.md` for Markdown and documentation standards.

## Skillify Adaptation Rule

Map durable, file-scoped conventions to Cline rules. Keep complex task procedures as skills, prompts, or linked workflow docs rather than always-on instructions.

