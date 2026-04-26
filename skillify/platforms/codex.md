# OpenAI Codex — Skill Setup via AGENTS.md

## How It Works

Codex discovers project instructions from `AGENTS.md` files placed at the repository root or within subdirectories. Unlike Claude/Gemini, Codex does **not** use YAML frontmatter or a global `~/.codex/skills/` directory.

Instead, the skill content is embedded directly into an `AGENTS.md` file using plain markdown. Codex resolves `AGENTS.md` hierarchically — root-level files provide global defaults, subdirectory files override for specific scopes.

## Install

Copy the generated `AGENTS.md` into any repository root where you want the skill available:

```bash
cp agents.md /path/to/your/project/AGENTS.md
```

Or append to an existing `AGENTS.md`:

```bash
echo "" >> /path/to/your/project/AGENTS.md
cat agents.md >> /path/to/your/project/AGENTS.md
```

## Format Differences from Claude/Gemini

| Feature | Claude / Gemini | Codex |
|---|---|---|
| Discovery | YAML frontmatter `name` + `description` | Markdown headings in `AGENTS.md` |
| Scope | Global (`~/.claude/skills/`) | Per-repo (`AGENTS.md` at root) |
| Reference files | `references/*.md` loaded on demand | Inline in `AGENTS.md` or linked files |
| Progressive disclosure | Native (filesystem-based) | Manual (use `## Section` headers) |

## Generated File

See [agents.md](agents.md) for the Codex-compatible version of this skill.

## Verify

Place the `AGENTS.md` in a project, open Codex, and ask: *"Create a skill for processing CSV files"* — the agent should follow the embedded authoring workflow.
