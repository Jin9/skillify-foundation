# Platform Setup Guide — crafting-backend-code

## Supported Platforms

| Platform | Config Location | Format | Install |
|---|---|---|---|
| **Claude** | `~/.claude/skills/crafting-backend-code/` | YAML frontmatter + markdown | `install.sh` (auto) |
| **Gemini** | `~/.gemini/skills/crafting-backend-code/` | YAML frontmatter + markdown | `install.sh` (auto) |
| **Antigravity** | `~/.gemini/skills/crafting-backend-code/` | Shared with Gemini + user rule | `install.sh` (auto + manual rule) |
| **Codex** | `AGENTS.md` in project root | Plain markdown | Copy `agents.md` to project |
| **Copilot** | `.github/copilot-instructions.md` | Plain markdown | Copy to `.github/` |

## Quick Install

```bash
bash platforms/install.sh
```

This copies the canonical `SKILL.md` + `references/` to Claude and Gemini global skill directories, and provides instructions for Codex, Copilot, and Antigravity.

## Architecture

Single canonical source (`SKILL.md` + `references/`) with platform-specific wrappers:
- **Claude / Gemini / Antigravity** → identical `SKILL.md` with progressive disclosure via `references/`
- **Codex** → `agents.md` inlines decision rules and guardrails (no progressive disclosure support)
- **Copilot** → `copilot-instructions.md` condensed to ~1 page (token budget constraint)

## Files

- [install.sh](install.sh) — Cross-platform install script
- [agents.md](agents.md) — Codex-compatible AGENTS.md
- [copilot-instructions.md](copilot-instructions.md) — Copilot-compatible instructions
