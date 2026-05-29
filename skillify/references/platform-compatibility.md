# Platform Compatibility Guide

Skills are an open standard: a directory with `SKILL.md` plus optional
`references/`, `templates/`, `scripts/`, `assets/`, `examples/`, and
platform metadata. Repo instruction files such as `AGENTS.md`,
`CLAUDE.md`, `GEMINI.md`, and `.github/copilot-instructions.md` are
always-on policy files, not substitutes for an on-demand skill.

## Skill vs Adjacent Surfaces

A `SKILL.md` is one of several host configuration surfaces. Keep them distinct; do not fold one into another:

- **Skill** — a repeatable, triggered *workflow*. This is what skillify engineers.
- **Custom agent / subagent** — a persistent *persona* with its own tool permissions and model preference. It answers "who should act," not "how to perform this task."
- **Always-on rules** (`AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.github/copilot-instructions.md`) — durable repository *conventions*.
- **Memory** — ephemeral, session-scoped *facts*.

A request about a persona, tool permissions, or model choice is custom-agent work, not a `SKILL.md`.

## Compatibility Matrix

| Platform | User skill locations | Project skill locations | Always-on rules | Activation |
|---|---|---|---|---|
| Claude Code | `~/.claude/skills/<name>/` | `.claude/skills/<name>/` | `CLAUDE.md` | Implicit by `description` or explicit `/name` |
| OpenAI Codex | `~/.agents/skills/<name>/`, `$CODEX_HOME/skills/<name>/` | `.agents/skills/<name>/` | `AGENTS.md` | Implicit by `description` or explicit `$name` |
| GitHub Copilot | `~/.copilot/skills/<name>/`, `~/.agents/skills/<name>/` | `.github/skills/<name>/`, `.agents/skills/<name>/` | `.github/copilot-instructions.md`, `AGENTS.md` | Implicit by `description` in agent hosts |
| Gemini CLI | `~/.gemini/skills/<name>/`, `~/.agents/skills/<name>/` | `.gemini/skills/<name>/`, `.agents/skills/<name>/` | `GEMINI.md`, `AGENTS.md` | Implicit by `description` after skill activation consent |
| Antigravity | `~/.gemini/antigravity/skills/<name>/` | `.agents/skills/<name>/` | `~/.gemini/GEMINI.md`, `.agents/rules/` | Implicit by `description`; rules can reinforce triggers |

## Universal Format

- Use `SKILL.md` with YAML frontmatter.
- Required fields: `name` and `description`.
- Use lowercase kebab-case names under 64 characters.
- Keep the folder name aligned with the `name` field.
- Do not use reserved vendor terms such as `claude` or `anthropic` in skill names.
- Put platform UI metadata outside `SKILL.md` when a host requires it.

## Portability Rules

- Keep `SKILL.md` platform-neutral: say "Run the command" instead of "Ask Claude to run the command."
- Put host-specific installation notes in `platforms/`, not in the reusable workflow.
- Prefer `.agents/skills/<name>/` for repository-scoped cross-agent skills.
- Copy the whole skill folder, not only `SKILL.md`; referenced `scripts/`, `templates/`, and `platforms/` must travel with it.
- Keep repo policy in always-on rule files. Keep repeatable workflows in `SKILL.md`.

## Activation Mechanisms

Beyond `description`-based implicit loading, some hosts add explicit activation controls. When adapting a skill, map to the host's mechanism instead of assuming description-only loading:

- **Windsurf** uses `activation_mode`: `always_on`, `model_decision` (description-driven, the default), `glob` (path-matched), or `manual`.
- **Cline** uses a `paths:` glob list to scope activation to matching files; Claude Code's `paths:` field gives the same path-conditional effect.

Treat these as host-specific rule-surface controls, not portable frontmatter; isolate them per the Adapt mode.

## Adapting Across Platforms

### Claude Code

1. Install the skill directory under `~/.claude/skills/` or `.claude/skills/`.
2. Keep YAML frontmatter in `SKILL.md`.
3. Remove or isolate Claude-only fields before sharing with hosts that do not support them.

### OpenAI Codex

1. Install the skill directory under `.agents/skills/`, `~/.agents/skills/`, or `$CODEX_HOME/skills/`.
2. Use `AGENTS.md` only for always-on project instructions.
3. For Codex app/plugin packaging, put UI metadata in `agents/openai.yaml`:

```yaml
interface:
  display_name: Skill Display Name
  short_description: One concise UI sentence.
  default_prompt: Use this skill to [task].
```

Keep this metadata aligned with the frontmatter `description`, but do not
copy the full trigger list into the UI chip.

### GitHub Copilot

1. Install agent skills under `.github/skills/`, `.agents/skills/`,
   `~/.copilot/skills/`, or `~/.agents/skills/`.
2. Keep `.github/copilot-instructions.md` for short always-on repository guidance.
3. Use skills for detailed workflows that should load only when relevant.

### Gemini CLI

1. Install skills under `.gemini/skills/`, `.agents/skills/`,
   `~/.gemini/skills/`, or `~/.agents/skills/`.
2. Keep `GEMINI.md` for always-on background rules.
3. Use `gemini skills list` or `/skills list` to verify discovery.

### Antigravity

1. Install global skills under `~/.gemini/antigravity/skills/`.
2. Install workspace skills under `.agents/skills/`.
3. Use Rules only to reinforce behavior that should apply broadly; do not duplicate full skill instructions there.

## Distribution

- Host reusable skills in a public or internal Git repository.
- Keep human docs at the repository level, not inside the skill folder.
- Include per-platform installation instructions in `platforms/`.
- Prefer plugins or host package managers when distributing beyond one repo or one machine.
