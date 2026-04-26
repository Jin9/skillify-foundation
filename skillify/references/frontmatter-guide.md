# Frontmatter Field Guide

The YAML frontmatter is how agents decide whether to load your skill. It appears in the system prompt, so every character matters.

## Required Fields

### `name` (required)

- **Format**: kebab-case only. No spaces, underscores, or capitals.
- **Length**: Under 64 characters.
- **Must match**: The skill folder name.
- **Examples**: `pdf-editor` ✅, `frontend-design` ✅, `My Cool Skill` ❌, `notion_project_setup` ❌

### `description` (required)

- **Length**: Under 1024 characters.
- **Must include**:
  1. What the skill does
  2. When to use it (specific trigger phrases)
  3. Key capabilities
- **No XML angle brackets** (`<` or `>`). They can inject instructions into the system prompt.
- **Structure**: `[What it does] + [When to use it] + [Key capabilities]`

## Recommended Fields

### `when_to_use` (Claude Code)

Additional context for when the agent should invoke the skill. Appended to `description` in the skill listing. Counts toward the 1,536-character cap combined with `description`.

### `compatibility`

1–500 characters. Indicates environment requirements:
- Intended product (e.g., `claude-code, opencode`)
- Required system packages
- Network access needs

### `metadata` (optional)

Custom key-value pairs. Suggested fields:
```yaml
metadata:
  author: YourName
  version: 1.0.0
  mcp-server: your-mcp-server
```

## Advanced Fields (Claude Code)

| Field | Type | Description |
|-------|------|-------------|
| `argument-hint` | string | Hint shown during autocomplete (e.g., `[issue-number]`) |
| `arguments` | string/list | Named positional arguments for `$name` substitution |
| `disable-model-invocation` | boolean | Set `true` to prevent auto-invocation |
| `user-invocable` | boolean | Set `false` to hide from `/` menu (background knowledge only) |
| `allowed-tools` | string | Tools allowed without permission prompts |
| `model` | string | Model override (e.g., `haiku`, `sonnet`, `opus`) |
| `effort` | string | Effort level override (`low`, `medium`, `high`, `xhigh`, `max`) |
| `context` | string | Set to `fork` for isolated subagent context |
| `agent` | string | Subagent type when `context: fork` (default: `general-purpose`) |
| `paths` | string/list | Glob patterns limiting auto-activation to matching files |
| `hooks` | object | Lifecycle hooks scoped to the skill |
| `shell` | string | Shell for command blocks (`bash` or `powershell`) |

## Advanced Fields (Codex)

| Field | Type | Description |
|-------|------|-------------|
| `metadata.short-description` | string | Short description for UI display |
| `agents/openai.yaml` | file | UI metadata: `display_name`, `short_description`, `default_prompt` |

## Security Restrictions

- No XML angle brackets (`<` `>`) in frontmatter
- No "claude" or "anthropic" in `name` (reserved)
- Frontmatter appears in the system prompt — malicious content could inject instructions

---

## Good vs. Bad Description Examples

### ✅ Good — specific and actionable

```yaml
description: >
  Analyzes Figma design files and generates developer handoff documentation.
  Use when user uploads .fig files, asks for "design specs",
  "component documentation", or "design-to-code handoff".
```

### ✅ Good — includes trigger phrases

```yaml
description: >
  Manages Linear project workflows including sprint planning, task creation,
  and status tracking. Use when user mentions "sprint", "Linear tasks",
  "project planning", or asks to "create tickets".
```

### ✅ Good — clear value proposition with negative triggers

```yaml
description: >
  Advanced data analysis for CSV files. Use for statistical modeling,
  regression, clustering. Do NOT use for simple data exploration
  (use data-viz skill instead).
```

### ❌ Bad — too vague

```yaml
description: Helps with projects.
```

### ❌ Bad — missing triggers

```yaml
description: Creates sophisticated multi-page documentation systems.
```

### ❌ Bad — too technical, no user triggers

```yaml
description: Implements the Project entity model with hierarchical relationships.
```

### ❌ Bad — exceeds character limit or contains XML

```yaml
description: "<skill>This is a very long description that also uses XML tags...</skill>"
```
