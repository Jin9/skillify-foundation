# OpenAI Codex - Skill Setup

## How It Works

Codex supports on-demand skills through `SKILL.md` directories and
always-on project instructions through `AGENTS.md`. Use a skill when the
workflow should load only for matching requests; use `AGENTS.md` for
repo-wide policy that should always apply.

Supported skill locations include the cross-agent `.agents/skills/` path
and Codex's local `$CODEX_HOME/skills/` path.

## Install

Run the cross-platform installer from this directory:

```bash
bash install.sh
```

Manual personal install from `skillify/platforms/`:

```bash
for DEST in "$HOME/.agents/skills/skillify" "${CODEX_HOME:-$HOME/.codex}/skills/skillify"; do
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

For project-scoped use, copy the skill folder to:

```text
<repo>/.agents/skills/skillify/
```

## Optional AGENTS.md Wrapper

If you want always-on project instructions instead of an on-demand skill,
append the prebuilt wrapper:

```bash
cat agents.md >> /path/to/your/project/AGENTS.md
```

Do not use the wrapper as a replacement for the `SKILL.md` skill when you
need trigger-based loading or referenced support files.

## Verify

Restart Codex if required by your host, then ask:
*"Create a skill for processing CSV files"*. The agent should load
`skillify`, or you can explicitly invoke `$skillify` where supported.

## Format Notes

- Keep `SKILL.md` frontmatter intact for skill discovery.
- Keep `AGENTS.md` short and repo-specific.
- Include the referenced support directories (`references/`, `templates/`, `scripts/`, `examples/`); `platforms/` stays behind — install notes are not part of the deployed skill.
