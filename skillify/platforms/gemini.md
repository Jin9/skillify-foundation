# Gemini CLI - Skill Setup

## How It Works

Gemini CLI supports skills as `SKILL.md` directories. Skills can live in
Gemini-native locations or in the shared `.agents/skills/` locations used
by multiple agent hosts.

Gemini uses `GEMINI.md` and `AGENTS.md` for always-on instructions. Keep
repeatable workflows in `SKILL.md` and use rule files only for background
policy.

## Install

Run the cross-platform installer from this directory:

```bash
bash install.sh
```

Manual personal install from `skillify/platforms/`:

```bash
for DEST in "$HOME/.gemini/skills/skillify" "$HOME/.agents/skills/skillify"; do
  mkdir -p "$DEST"
  cp ../SKILL.md "$DEST/"
  for dir in references templates scripts platforms examples assets; do
    if [ -d "../$dir" ]; then
      mkdir -p "$DEST/$dir"
      cp -R "../$dir/." "$DEST/$dir/"
    fi
  done
done
```

For project-scoped use, copy the skill folder to either location:

```text
<repo>/.gemini/skills/skillify/
<repo>/.agents/skills/skillify/
```

## Verify

```bash
gemini skills list
```

You can also use `/skills list` from a Gemini CLI session. Then ask:
*"Create a skill for processing CSV files"*. Gemini should activate the
skill after discovery and approval.

## Format Notes

- Gemini CLI uses the canonical `SKILL.md` format.
- `GEMINI.md` can mention that skill authoring should use `skillify`, but should not duplicate the full workflow.
- Include referenced support directories when copying a skill.
