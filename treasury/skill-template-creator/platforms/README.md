# Platform Setup Guide

This directory contains platform-specific configurations to install the `creating-skill-templates` skill globally across all supported AI coding agents.

## Supported Platforms

| Platform | Config Location | Format | Install Method |
|---|---|---|---|
| **Claude** (Claude Code) | `~/.claude/skills/creating-skill-templates/` | YAML frontmatter + markdown | Copy `SKILL.md` + `references/` |
| **Gemini** (Gemini CLI) | `~/.gemini/skills/creating-skill-templates/` | YAML frontmatter + markdown | Copy `SKILL.md` + `references/` |
| **Codex** (OpenAI Codex) | `AGENTS.md` in project root | Plain markdown, no frontmatter | Append or merge into `AGENTS.md` |
| **Copilot** (GitHub Copilot) | `.github/copilot-instructions.md` | Plain markdown | Append or merge into instructions |
| **Antigravity** (Google DeepMind) | `~/.gemini/antigravity/` + user rules | YAML frontmatter + markdown | Copy to Gemini skills + user rule |

## Quick Install

Run the install script from this directory:

```bash
bash install.sh
```

The script will:
1. Copy the canonical `SKILL.md` and `references/` to Claude and Gemini global skill directories.
2. Generate `AGENTS.md` (Codex-compatible) in the skill root.
3. Generate `.github/copilot-instructions.md` (Copilot-compatible) in the skill root.
4. Register the skill trigger in Antigravity's user rules.

## Manual Install

See the platform-specific files in this directory:
- [claude.md](claude.md) — Claude Code setup
- [gemini.md](gemini.md) — Gemini CLI setup
- [codex.md](codex.md) — OpenAI Codex setup
- [copilot.md](copilot.md) — GitHub Copilot setup
- [antigravity.md](antigravity.md) — Antigravity setup

## Architecture Decision

We use **one canonical source** (`SKILL.md` + `references/`) and **generate platform-specific wrappers** from it, rather than maintaining five separate copies. This ensures:
- Single source of truth for skill content
- No drift between platforms
- Platform-specific adaptations are minimal wrappers, not full rewrites
