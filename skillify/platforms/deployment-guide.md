# Platform Deployment Guide

This directory contains platform-specific instructions for installing
`skillify` across AI coding agents.

## Supported Platforms

| Platform | Primary skill locations | Optional always-on file | Install method |
|---|---|---|---|
| Claude Code | `~/.claude/skills/skillify/`, `.claude/skills/skillify/` | `CLAUDE.md` | Copy full skill folder |
| OpenAI Codex | `~/.agents/skills/skillify/`, `.agents/skills/skillify/`; `$CODEX_HOME/skills/skillify/` only for legacy compatibility | `AGENTS.md` | Copy full skill folder |
| GitHub Copilot | `~/.copilot/skills/skillify/`, `~/.agents/skills/skillify/`, `.github/skills/skillify/`, `.agents/skills/skillify/` | `.github/copilot-instructions.md` | Copy full skill folder |
| Gemini CLI | `~/.agents/skills/skillify/`, `.agents/skills/skillify/` | `GEMINI.md`, `AGENTS.md` | Copy full skill folder |
| Antigravity | `~/.gemini/antigravity/skills/skillify/`, `.agents/skills/skillify/` | `~/.gemini/GEMINI.md`, `.agents/rules/` | Copy full skill folder |

## Quick Install

Run from this directory:

```bash
bash install.sh
```

By default the script copies the full skill folder to these user/global
locations:

1. Claude Code: `~/.claude/skills/skillify/`
2. Shared Agent Skills (read by Codex): `~/.agents/skills/skillify/`
3. Antigravity global path: `~/.gemini/antigravity-cli/skills/skillify/` (falls back to `~/.gemini/antigravity/skills/skillify/`)

Each target is cleanly refreshed (the `skillify/` folder is replaced), while
any other skills already installed alongside it are left untouched.

GitHub Copilot's native path is **not** installed by default. To also install
to `~/.copilot/skills/skillify/`, run:

```bash
INSTALL_COPILOT=1 bash install.sh
```

By default the script also skips `$CODEX_HOME/skills/skillify/` because current
Codex sessions also read `~/.agents/skills/skillify/`; installing both creates
duplicate `skillify` entries. For legacy Codex clients that only read
`$CODEX_HOME/skills`, run:

```bash
INSTALL_CODEX_COMPAT=1 bash install.sh
```

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
