# Antigravity (Google DeepMind) — Skill Setup

## How It Works

Antigravity runs on the Gemini infrastructure and uses the same skill format as Gemini CLI (`~/.gemini/skills/`). Additionally, Antigravity supports **user rules** defined in its settings that act as persistent behavioral triggers across all conversations.

The skill installation has two parts:
1. **Skill files** — installed to Gemini's global skill directory (shared with Gemini CLI).
2. **User rule** — a trigger rule in Antigravity's user settings that activates the skill when relevant keywords are detected.

## Install

### Step 1: Install Skill Files (shared with Gemini)

```bash
# Create skill directory (reuses Gemini path)
mkdir -p ~/.gemini/skills/skillify/references

# Copy canonical files
cp ../SKILL.md ~/.gemini/skills/skillify/SKILL.md
cp ../references/*.md ~/.gemini/skills/skillify/references/
```

### Step 2: Add User Rule

Add the following to your Antigravity user rules (Settings → User Rules):

```
When the user asks to create, author, write, refactor, review, or critique a SKILL.md or agent skill, activate the 'skillify' skill from ~/.gemini/skills/skillify/ and follow its authoring workflow, validation loop, and review checklist.
```

## Verify

Open Antigravity and ask: *"Create a skill for processing CSV files"* — the user rule should trigger skill activation, and the agent should follow the authoring workflow with the validation loop.

## Format Notes

- Antigravity uses the **exact same SKILL.md format** as Gemini CLI — no adaptation needed.
- The user rule acts as a **persistent memory trigger** that ensures the skill is activated even when the default skill discovery mechanism doesn't fire.
- Antigravity's `~/.gemini/GEMINI.md` global instructions are also respected and can contain additional memory triggers.
