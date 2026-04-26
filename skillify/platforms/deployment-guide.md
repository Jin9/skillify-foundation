# Platform Deployment Guide

This directory contains platform-specific instructions for installing
`skillify` across AI coding agents.

## Supported Platforms

| Platform | Primary skill locations | Optional always-on file | Install method |
|---|---|---|---|
| Claude Code | `~/.claude/skills/skillify/`, `.claude/skills/skillify/` | `CLAUDE.md` | Copy full skill folder |
| OpenAI Codex | `~/.agents/skills/skillify/`, `$CODEX_HOME/skills/skillify/`, `.agents/skills/skillify/` | `AGENTS.md` | Copy full skill folder |
| GitHub Copilot | `~/.copilot/skills/skillify/`, `~/.agents/skills/skillify/`, `.github/skills/skillify/`, `.agents/skills/skillify/` | `.github/copilot-instructions.md` | Copy full skill folder |
| Gemini CLI | `~/.gemini/skills/skillify/`, `~/.agents/skills/skillify/`, `.gemini/skills/skillify/`, `.agents/skills/skillify/` | `GEMINI.md`, `AGENTS.md` | Copy full skill folder |
| Antigravity | `~/.gemini/antigravity/skills/skillify/`, `.agents/skills/skillify/` | `~/.gemini/GEMINI.md`, `.agents/rules/` | Copy full skill folder |

## Quick Install

Run from this directory:

```bash
bash install.sh
```

The script copies the full skill folder to user/global locations for the
supported hosts:

1. Claude Code: `~/.claude/skills/skillify/`
2. Shared Agent Skills: `~/.agents/skills/skillify/`
3. Codex compatibility path: `$CODEX_HOME/skills/skillify/`
4. Gemini native path: `~/.gemini/skills/skillify/`
5. Copilot native path: `~/.copilot/skills/skillify/`
6. Antigravity global path: `~/.gemini/antigravity/skills/skillify/`

It does not edit project rule files automatically. Use the prebuilt
wrappers only when you want always-on instructions:

- [agents.md](agents.md) - optional Codex `AGENTS.md` wrapper
- [copilot-instructions.md](copilot-instructions.md) - optional Copilot custom instructions wrapper

## Manual Install

See the platform-specific files in this directory:

- [claude.md](claude.md) - Claude Code setup
- [codex.md](codex.md) - OpenAI Codex setup
- [copilot.md](copilot.md) - GitHub Copilot setup
- [gemini.md](gemini.md) - Gemini CLI setup
- [antigravity.md](antigravity.md) - Antigravity setup

## Architecture Decision

Use one canonical source (`SKILL.md` plus support directories) and
platform-specific installation notes. This avoids drift while preserving
host-specific discovery paths.

Copy the whole skill folder, not just `SKILL.md`: this skill references
`references/`, `templates/`, `scripts/`, and `platforms/`.
