# GitHub Copilot: Adding Agent Skills

Source: https://docs.github.com/en/copilot/how-tos/use-copilot-agents/cloud-agent/create-skills  
Related source: https://docs.github.com/copilot/how-tos/copilot-cli/customize-copilot/add-skills  
Accessed: 2026-05-14  
Category: GitHub Copilot / skill creation and distribution

## Why This Source Matters

GitHub's current skill documentation shows how the Agent Skills standard is implemented across Copilot cloud agent, GitHub Copilot CLI, and VS Code agent mode. It also adds operational details that are not always present in generic skill guides, especially around installation, provenance, and supply-chain review.

## Core Data Points

- Agent skills work with Copilot cloud agent, GitHub Copilot CLI, and agent mode in VS Code.
- A skill is a folder containing `SKILL.md` and optional resources such as supplementary Markdown files or scripts.
- Project skills can be stored in `.github/skills`, `.claude/skills`, or `.agents/skills`.
- Personal skills can be stored in `~/.copilot/skills`, `~/.claude/skills`, or `~/.agents/skills`.
- Each skill gets its own subdirectory, and the subdirectory name should be lowercase with hyphens.
- `SKILL.md` requires YAML frontmatter:
  - `name`: required, unique, lowercase, hyphen-separated.
  - `description`: required, explains what the skill does and when Copilot should use it.
  - `license`: optional.
- Shared skills should be inspected before use because skills may include hidden instructions or executable scripts.
- `gh skill` supports discovery, preview, installation, updating, and publishing.
- `gh skill preview` lets a user inspect the file tree and rendered `SKILL.md` before installing.
- `gh skill install` can target a specific agent host and scope.
- Installed skills can include provenance metadata in frontmatter, including source repository, ref, and tree SHA.
- Pinned versions are skipped by update flows.
- Publishing flows can validate against the Agent Skills specification and check remote repository settings.

## Design Implications For Skillify

- Generated skills should be safe to install from source control and should not hide behavior in obscure support files.
- Provenance metadata is worth preserving when adapting or refactoring an installed skill.
- Skillify should distinguish authoring a local skill from publishing a skill repository.
- A production audit should include supply-chain review: scripts, remote links, provenance, pinning, and license.
- If a skill is intended for Copilot, `.github/skills` is a first-class project path, but `.agents/skills` improves portability.

## Operational Checklist

- Before installing a community skill, inspect `SKILL.md`, file tree, scripts, and any remote references.
- Pin critical third-party skills where repeatability matters.
- Prefer project scope for team-shared workflow skills and user scope for personal habits.
- Keep `license` explicit when a skill is intended for sharing.
- Verify that all support files referenced by `SKILL.md` are actually included in the folder.

