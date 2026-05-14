# Open Agent Skills Specification

Source: https://openagentskills.dev/docs/specification  
Accessed: 2026-05-14  
Category: Open standard / skill format reference

## Why This Source Matters

This source is the most explicit cross-agent specification for the `SKILL.md` folder format. It is useful as a neutral baseline when designing skills that should work across Codex, Claude, GitHub Copilot, Windsurf, and other skills-compatible agents.

## Core Data Points

- A skill is a directory with a required `SKILL.md` file.
- `SKILL.md` must contain YAML frontmatter followed by Markdown instructions.
- Required frontmatter fields are `name` and `description`.
- The `name` field is constrained to 1-64 characters, lowercase letters/numbers/hyphens, no leading or trailing hyphen, no consecutive hyphens, and it should match the parent directory.
- The `description` field is constrained to 1-1024 characters and should describe both capability and activation context.
- Optional frontmatter fields include `license`, `compatibility`, `metadata`, and experimental `allowed-tools`.
- The body has no strict schema, but recommended sections include step-by-step instructions, input/output examples, and common edge cases.
- Optional directories include `scripts/`, `references/`, and `assets/`.
- Progressive disclosure is treated as a first-class design principle:
  - metadata is available at startup,
  - the full `SKILL.md` body is loaded when the skill activates,
  - supporting resources are loaded only as needed.
- The source recommends keeping `SKILL.md` under 500 lines and moving detailed material into referenced files.
- File references should be relative to the skill root and kept shallow.
- Validation is supported through the `skills-ref` reference library.

## Design Implications For Skillify

- Skillify should enforce the standard `name` constraints, not merely check that a name exists.
- Skillify should treat the `description` as the main retrieval surface, not as decorative metadata.
- `compatibility` is useful when a skill needs a specific host, package, network access, or runtime.
- `allowed-tools` should be treated as experimental: preserve it when present, but do not require it for portability.
- A portable skill generator should keep folder depth shallow: `SKILL.md`, `references/`, `templates/`, `scripts/`, and `assets/` are enough for most cases.
- A validator should catch parent-directory/name mismatches because some hosts silently fail to load invalid skills.

## Useful Checklist

- Does every skill folder contain exactly one `SKILL.md` at its root?
- Does `name` match the folder name?
- Is `description` specific enough to trigger the skill without reading the body?
- Is the body short enough to load cheaply?
- Are long examples, domain rules, and templates moved to support files?
- Are referenced files one level deep and linked from `SKILL.md`?
- Are scripts self-contained, deterministic, and documented where dependencies exist?

