# Claude Code - Skill Setup

## How It Works

Claude Code discovers skills from personal and project skill directories.
Each skill is a directory containing `SKILL.md` with YAML frontmatter
(`name` and `description`) plus optional supporting directories such as
`references/`, `templates/`, `scripts/`, and `assets/`.

Claude keeps skill metadata available and loads the full skill when the
request matches the description or the user explicitly invokes it.

## Install

Run the cross-platform installer from this directory:

```bash
bash install.sh
```

Manual personal install from `skillify/platforms/`:

```bash
DEST="$HOME/.claude/skills/skillify"
mkdir -p "$DEST"
cp ../SKILL.md "$DEST/"
for dir in references templates scripts platforms examples assets; do
  if [ -d "../$dir" ]; then
    mkdir -p "$DEST/$dir"
    cp -R "../$dir/." "$DEST/$dir/"
  fi
done
```

For project-scoped use, copy the same skill folder to
`.claude/skills/skillify/` inside the target repository.

## Verify

```bash
ls -la ~/.claude/skills/skillify/
# Should show: SKILL.md plus referenced support directories.
```

Then open Claude Code and ask: *"Create a skill for processing CSV files"*.
The skill should trigger automatically, or you can invoke it explicitly if
your host supports slash-based skill calls.

## Format Notes

- Claude Code uses the canonical `SKILL.md` format directly.
- Keep frontmatter portable: `name`, `description`, and only supported optional fields.
- Keep repository policy in `CLAUDE.md`; keep reusable workflows in `SKILL.md`.
