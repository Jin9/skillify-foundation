# Claude Code — Global Skill Setup

## How It Works

Claude Code discovers skills from `~/.claude/skills/`. Each skill is a directory containing a `SKILL.md` with YAML frontmatter (`name` + `description`) and an optional `references/` directory.

Claude pre-loads all skill names and descriptions at startup. When a user request matches, Claude reads the full `SKILL.md` and loads reference files on demand.

## Install

```bash
# Create skill directory
mkdir -p ~/.claude/skills/skillify/references

# Copy canonical files
cp ../SKILL.md ~/.claude/skills/skillify/SKILL.md
cp ../references/*.md ~/.claude/skills/skillify/references/
```

## Verify

```bash
ls -la ~/.claude/skills/skillify/
# Should show: SKILL.md, references/
```

Then open Claude Code and ask: *"Create a skill for processing CSV files"* — the skill should trigger automatically.

## Permissions

If Claude Code requires explicit read permissions, add to `.claude/settings.local.json`:

```json
{
  "permissions": {
    "allow": [
      "Read(//Users/<username>/.claude/skills/**)"
    ]
  }
}
```

## Format Notes

- Claude Code uses the **exact same format** as the canonical `SKILL.md` — no adaptation needed.
- YAML frontmatter fields: `name` (lowercase, hyphens, ≤64 chars) and `description` (≤1024 chars).
- Reference files are loaded on-demand via filesystem access.
