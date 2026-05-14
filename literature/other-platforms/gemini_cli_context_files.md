# Gemini CLI Context Files: GEMINI.md

Source: https://google-gemini.github.io/gemini-cli/docs/cli/gemini-md.html  
Related source: https://google-gemini.github.io/gemini-cli/docs/core/memport.html  
Accessed: 2026-05-14  
Category: Other platforms / repo instructions and memory

## Why This Source Matters

Gemini CLI uses `GEMINI.md` as its standing instruction file. It is not a `SKILL.md` system, but it is an important adjacent pattern for platform adaptation because it shows how global, hierarchical, and imported instruction files are loaded.

## Core Data Points

- `GEMINI.md` provides reusable instructional context for the Gemini model.
- It can contain project instructions, persona guidance, coding style, or other persistent context.
- Gemini CLI loads context hierarchically:
  - global context from `~/.gemini/GEMINI.md`,
  - project root and ancestor context files,
  - subdirectory context files relevant to the current working path.
- Loaded context is concatenated and sent with prompts.
- `/memory show` displays the full loaded context.
- `/memory reload` or refresh-style commands force the CLI to rescan context files.
- `/memory add <text>` appends persistent memory to the global file.
- Context can be modularized with `@file.md` imports.
- Imports support relative and absolute paths, nested imports, circular import detection, access validation, and a default maximum depth.
- The context file name is configurable in settings, including lists such as `AGENTS.md`, `CONTEXT.md`, and `GEMINI.md`.

## Design Implications For Skillify

- Gemini adaptation should usually map portable skills into `AGENTS.md` or `GEMINI.md` guidance unless the host supports the Agent Skills standard directly.
- The import mechanism gives Gemini a progressive-disclosure-like pattern, but it is still standing context rather than relevance-triggered skill activation.
- Large `GEMINI.md` files should be modularized with imports, but imports must avoid cycles and excessive depth.
- A cross-platform adapter should emit a concise root instruction file plus imported modules for larger guidance.
- Because Gemini concatenates context, generic instructions can become noisy. Skillify should keep Gemini-oriented guidance scoped and specific.

## Practical Mapping

- `SKILL.md` frontmatter description -> short trigger note in `GEMINI.md`.
- `SKILL.md` body -> imported workflow file, referenced only when relevant.
- `references/` -> imported files or explicit file references.
- `scripts/` -> preserved scripts with clear instructions in the imported workflow.
- Platform-specific validation -> check loaded context with `/memory show`.

