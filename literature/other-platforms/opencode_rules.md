# OpenCode Rules And AGENTS.md

Source: https://dev.opencode.ai/docs/rules/  
Accessed: 2026-05-14  
Category: Other platforms / AGENTS.md instructions

## Why This Source Matters

OpenCode uses `AGENTS.md` as its primary custom-instruction mechanism and supports external instruction references through configuration. It provides a useful pattern for adapting Skillify output to agents that do not use `SKILL.md` folders.

## Core Data Points

- OpenCode custom instructions are provided through `AGENTS.md`.
- `/init` can generate or update an `AGENTS.md` file by scanning the project.
- Project-level `AGENTS.md` should be committed for team consistency.
- Global instructions live at `~/.config/opencode/AGENTS.md`.
- OpenCode has Claude Code compatibility via `CLAUDE.md` fallbacks.
- Precedence order:
  - local files by traversing up from the current directory (`AGENTS.md`, `CLAUDE.md`),
  - global `~/.config/opencode/AGENTS.md`,
  - Claude Code global fallback unless disabled.
- If both `AGENTS.md` and `CLAUDE.md` exist in the same category, `AGENTS.md` wins.
- `opencode.json` can list instruction files and supports glob patterns.
- External references can be included manually in `AGENTS.md` with `@path` style references.

## Design Implications For Skillify

- OpenCode adaptation should emit `AGENTS.md` as the portable baseline.
- Use `opencode.json` when a repository needs multiple instruction files or globbed package-level guidance.
- Keep project `AGENTS.md` concise because it becomes standing context.
- Use external references for larger guidance that the agent can fetch when relevant.
- Preserve `CLAUDE.md` only as a fallback during migration; do not duplicate content across both files when `AGENTS.md` is the primary target.

## Practical Mapping

- Repo conventions -> root `AGENTS.md`.
- Package-specific guidance -> package `AGENTS.md` or files listed in `opencode.json`.
- Reusable workflow -> referenced Markdown file or command-like instruction.
- Skills with scripts -> keep scripts under version control and add concise usage instructions to `AGENTS.md`.

