# Cursor Rules And AGENTS.md Support

Source: https://docs.cursor.com/context/rules-for-ai  
Accessed: 2026-05-14  
Category: Other platforms / persistent AI rules

## Why This Source Matters

Cursor rules are a mature example of scoped persistent instructions. They are not the same as `SKILL.md`, but their activation model is relevant when Skillify adapts skill guidance into tools that use rule files instead of skills.

## Core Data Points

- Cursor rules provide reusable, persistent, prompt-level instructions for Agent and Inline Edit.
- Project rules are stored in `.cursor/rules` and are version-controlled.
- User rules are global and apply across the local Cursor environment.
- `AGENTS.md` is supported as a simple alternative for root-level project instructions.
- Legacy `.cursorrules` remains supported but is deprecated in favor of project rules.
- Rules are included in model context when they apply.
- `AGENTS.md` is plain Markdown without rule metadata.
- Current documented limitations for `AGENTS.md` include root-level-only placement, global project scope, and no built-in split across multiple files.

## Design Implications For Skillify

- Cursor adaptation should prefer `.cursor/rules` for scoped or conditional guidance.
- Use `AGENTS.md` for portable, simple root-level instructions.
- Do not map every skill to Cursor rules automatically. Reusable task workflows may be better as prompts or external docs unless the workflow should influence every matching edit.
- Since rule content is prompt-level context, it should be concise and concrete.
- Avoid putting long procedural skill bodies into always-on Cursor rules.

## Practical Mapping

- Repository conventions -> `AGENTS.md` or `.cursor/rules`.
- Directory-specific conventions -> `.cursor/rules` when scoping is needed.
- Multi-step workflow with references -> keep as separate documentation and link from a rule or prompt.
- Scripted workflow -> preserve scripts and provide a rule only for when to use them.

