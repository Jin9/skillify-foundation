# Gemini CLI Custom Commands

Source: https://google-gemini.github.io/gemini-cli/docs/cli/custom-commands.html  
Accessed: 2026-05-14  
Category: Other platforms / reusable prompt commands

## Why This Source Matters

Custom commands are Gemini CLI's reusable prompt mechanism. They are closer to prompt files than to skills, but they are useful for adapting lightweight skills that do not need scripts, references, or automatic invocation.

## Core Data Points

- Custom commands let users save frequently used prompts as shortcuts.
- Commands can be global or project-local.
- Global commands live in `~/.gemini/commands/`.
- Project commands live in `<project-root>/.gemini/commands/`.
- Project commands override global commands with the same name.
- A command's name is derived from its file path relative to the commands directory.
- Directory nesting creates command namespaces.
- Commands can use arguments to parameterize prompts.
- TOML frontmatter can provide metadata such as a description.
- Commands are manually invoked, so they do not provide automatic relevance-triggered activation.

## Design Implications For Skillify

- Use Gemini custom commands for manual reusable workflows, especially when the workflow is a prompt template rather than a full skill bundle.
- Keep project commands under version control when they define team workflows.
- Prefer command namespaces for related workflows, such as `review/security` or `docs/release-notes`.
- A generated command should include a clear description so users can discover it quickly.
- Do not use custom commands for complex workflows that need many supporting files; use imported context files or a skills-compatible host instead.

## Adaptation Rule

When adapting a `SKILL.md` to Gemini CLI:

- If it is a manual checklist or reusable prompt, make a `.gemini/commands/<name>.md`.
- If it is durable repo policy, put it in `GEMINI.md` or `AGENTS.md`.
- If it requires scripts/templates, preserve a folder and reference those files from the command or context file.

