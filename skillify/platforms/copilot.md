# GitHub Copilot - Skill Setup

## How It Works

GitHub Copilot supports agent skills as `SKILL.md` directories for
agent workflows. Copilot also supports always-on repository instructions
through `.github/copilot-instructions.md` and path-specific
`.github/instructions/*.instructions.md` files.

Use a skill for reusable, trigger-based workflows. Use Copilot
instructions only for concise policy that should apply to every Copilot
interaction in a repository.

## Install

Run the cross-platform installer from this directory:

```bash
bash install.sh
```

Manual personal install from `skillify/platforms/`:

```bash
for DEST in "$HOME/.copilot/skills/skillify" "$HOME/.agents/skills/skillify"; do
  mkdir -p "$DEST"
  cp ../SKILL.md "$DEST/"
  for dir in references templates scripts examples assets; do
    if [ -d "../$dir" ]; then
      mkdir -p "$DEST/$dir"
      cp -R "../$dir/." "$DEST/$dir/"
    fi
  done
done
```

For project-scoped use, copy the skill folder to either location:

```text
<repo>/.github/skills/skillify/
<repo>/.agents/skills/skillify/
```

## Optional Custom Instructions

If you want always-on repository guidance, append the prebuilt compact
wrapper:

```bash
mkdir -p /path/to/your/project/.github
cat copilot-instructions.md >> /path/to/your/project/.github/copilot-instructions.md
```

Keep this wrapper short. Long workflows belong in `SKILL.md` so they load
only when relevant.

## Verify

Open a Copilot agent host for the target scope and ask:
*"Create a skill for processing CSV files"*. The agent should follow the
`skillify` workflow when skills are enabled for that host.

## Format Notes

- Agent skills use the canonical `SKILL.md` format.
- `.github/copilot-instructions.md` is always-on context, not skill discovery.
- Include referenced support directories when copying a skill.
