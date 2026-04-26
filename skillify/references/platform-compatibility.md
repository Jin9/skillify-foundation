# Platform Compatibility Guide

Skills are an open standard. This guide covers platform-specific conventions.

## Compatibility Matrix

| Feature | Claude Code | Codex | Copilot | Gemini | OpenCode |
|---------|-------------|-------|---------|--------|----------|
| Skill file | `SKILL.md` | `SKILL.md` | N/A | `AGENTS.md` | `SKILL.md` |
| Root rules | `CLAUDE.md` | `AGENTS.md` | `.github/copilot-instructions.md` | `AGENTS.md` | `AGENTS.md` |
| Skill dir | `.claude/skills/` | `skills/` | N/A | `.gemini/` | `.agents/skills/` |
| Frontmatter | YAML required | YAML required | N/A | Headings | YAML |
| Auto-trigger | ✅ | ✅ | ❌ | ❌ | ✅ |
| MCP support | ✅ | ❌ | ❌ | ✅ | ❌ |

## Naming Conventions (Universal)

- **kebab-case**: `pdf-editor` ✅
- No spaces, underscores, or capitals
- Folder name must match the `name` field
- Claude: No "claude" or "anthropic" in names (reserved)
- Codex: Under 64 characters, prefer verb-led phrases

## Portability Strategies

### Strategy 1: Compatibility Field

```yaml
---
name: build-test-verify
description: Run lint, test, and build verification.
compatibility: claude-code, opencode, codex
---
```

### Strategy 2: Platform-Agnostic Body

- Say "Run the following command" not "Ask Claude to run..."
- Say "Read `references/api.md`" not "Use Claude Code's file reader..."
- Reference tool capabilities generically

## Adapting Across Platforms

### Claude Code → Codex
1. Move from `.claude/skills/` to `skills/`
2. Add `agents/openai.yaml` with UI metadata
3. Remove Claude-specific frontmatter (`context`, `agent`, `hooks`)

Codex UI metadata lives outside `SKILL.md` in `agents/openai.yaml` when packaging a skill for a Codex plugin or marketplace entry:

```yaml
display_name: Skill Display Name
short_description: One concise UI sentence.
default_prompt: Use this skill to [task].
```

Keep this metadata aligned with the frontmatter `description`, but do not copy the full trigger list into the UI chip.

### Claude Code → Copilot
1. Extract instructions into `.github/copilot-instructions.md`
2. Convert YAML triggers to natural language
3. Copilot instructions are always-on (no auto-triggering)

### Claude Code → Gemini
1. Convert to `AGENTS.md` format with markdown headings
2. Keep `references/` for progressive disclosure
3. Gemini reads `AGENTS.md` at any directory level

## Distribution

- Host in a public GitHub repo
- Add README.md at the **repo level** (NOT inside the skill folder)
- Include installation instructions per platform
- Organization deployment: Claude Code via Settings, Codex via `skills/` directory
