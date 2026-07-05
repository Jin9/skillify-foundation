# Antigravity - Skill Setup

## How It Works

Antigravity supports global skills in the Antigravity-specific Gemini
directory and workspace skills in the shared `.agents/skills/` directory.
It also supports persistent rules, but rules should reinforce activation
or workspace policy rather than duplicate the full skill.

The recommended install has two choices:

1. **Global skill** - available across Antigravity workspaces.
2. **Workspace skill** - committed or copied into one repository.

## Install

Run the cross-platform installer from this directory:

```bash
bash install.sh
```

Manual global install from `skillify/platforms/`:

```bash
DEST="$HOME/.gemini/antigravity-cli/skills/skillify"
mkdir -p "$DEST"
cp ../SKILL.md "$DEST/"
for dir in references templates scripts examples assets; do
  if [ -d "../$dir" ]; then
    mkdir -p "$DEST/$dir"
    cp -R "../$dir/." "$DEST/$dir/"
  fi
done
```

For workspace-scoped use, copy the skill folder to:

```text
<repo>/.agents/skills/skillify/
```

## Optional Rule

Add this rule only if Antigravity does not reliably activate the skill:

```text
When the user asks to create, author, write, refactor, review, audit, compress, split, merge, or adapt a SKILL.md or agent skill, use the skillify skill and follow its workflow and validation gates.
```

Use global rules in `~/.gemini/GEMINI.md` and workspace rules in
`.agents/rules/`, according to the scope you want.

## Verify

Open Antigravity and ask: *"Create a skill for processing CSV files"*.
The agent should activate `skillify`; if not, add the optional rule above.

## Format Notes

- Antigravity uses the canonical `SKILL.md` format for skills.
- `.agents/skills/` is the portable workspace path.
- Rules are optional activation helpers, not the source of skill behavior.
