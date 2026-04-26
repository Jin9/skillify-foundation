# GitHub Copilot — Skill Setup via copilot-instructions.md

## How It Works

Copilot reads custom instructions from two locations:
1. **Repository-wide**: `.github/copilot-instructions.md` — applies to all Copilot interactions in that repo.
2. **Path-specific**: `.github/instructions/<name>.instructions.md` — uses YAML `applyTo` to scope to specific files/paths.

Copilot does **not** support YAML frontmatter for skill discovery or global `~/` skill directories. Instructions are injected as system context for chat and inline suggestions.

## Install

Copy the generated instructions file into any repository:

```bash
mkdir -p /path/to/your/project/.github
cp copilot-instructions.md /path/to/your/project/.github/copilot-instructions.md
```

Or append to an existing file:

```bash
echo "" >> /path/to/your/project/.github/copilot-instructions.md
cat copilot-instructions.md >> /path/to/your/project/.github/copilot-instructions.md
```

## Format Differences from Claude/Gemini

| Feature | Claude / Gemini | Copilot |
|---|---|---|
| Discovery | YAML `name` + `description` | Always loaded for the repo |
| Scope | Global skill directory | Per-repo `.github/` |
| Progressive disclosure | Reference files loaded on demand | All content loaded at once |
| Token budget | 500 lines per SKILL.md | ~1–2 pages recommended |
| Frontmatter | Required (`name`, `description`) | Optional (only `applyTo` for path scoping) |

## Generated File

See [copilot-instructions.md](copilot-instructions.md) for the Copilot-compatible version.

## Verify

Place the file in `.github/copilot-instructions.md`, open VS Code with Copilot, and ask in chat: *"Create a skill for processing CSV files"* — the instructions should guide Copilot's response.

## Limitations

- Copilot does not support on-demand file loading. The full instruction content is injected every time.
- Keep instructions concise (~1–2 pages) to avoid diluting Copilot's focus.
- For complex skills, consider using Copilot's `AGENTS.md` support (agentic mode) instead.
