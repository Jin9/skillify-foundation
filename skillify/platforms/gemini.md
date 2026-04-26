# Gemini CLI — Global Skill Setup

## How It Works

Gemini CLI discovers skills from `~/.gemini/skills/`. The format is identical to Claude Code: a `SKILL.md` with YAML frontmatter and optional `references/` directory.

Gemini pre-loads skill metadata at startup and reads the full skill body on demand when the task matches.

## Install

```bash
# Create skill directory
mkdir -p ~/.gemini/skills/skillify/references

# Copy canonical files
cp ../SKILL.md ~/.gemini/skills/skillify/SKILL.md
cp ../references/*.md ~/.gemini/skills/skillify/references/
```

## Verify

```bash
ls -la ~/.gemini/skills/skillify/
# Should show: SKILL.md, references/
```

Then open Gemini CLI and ask: *"Create a skill for processing CSV files"* — the skill should trigger automatically.

## Format Notes

- Gemini CLI uses the **exact same format** as Claude Code — no adaptation needed.
- The canonical `SKILL.md` works as-is.
- If you also use Gemini's global instructions (`~/.gemini/GEMINI.md`), no changes needed there unless you want to add a memory trigger.

## Optional: Add Memory Trigger

To make Gemini remember this skill exists across sessions, add to `~/.gemini/GEMINI.md`:

```markdown
- When the user mentions 'skill template' or 'create a skill', activate the 'skillify' skill to guide SKILL.md authoring using best practices.
```
