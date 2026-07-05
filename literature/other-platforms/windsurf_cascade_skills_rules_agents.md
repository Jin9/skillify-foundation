# Windsurf Cascade: Skills, Rules, AGENTS.md, Memories, And Workflows

Source: https://docs.windsurf.com/windsurf/cascade/memories
Accessed: 2026-05-14
Category: other-platforms / windsurf activation modes
Provenance: header added 2026-07-05; body unchanged

## Why This Source Matters

Windsurf Cascade documentation with the clearest activation-mode taxonomy in the corpus (always_on / model_decision / glob / manual) across memories, rules, AGENTS.md, workflows, and skills. Tier 1 official source for that host.


Sources:
- https://docs.windsurf.com/windsurf/cascade/memories
- https://docs.windsurf.com/windsurf/cascade/skills
- https://docs.windsurf.com/windsurf/cascade/agents-md
- https://docs.windsurf.com/windsurf/cascade/workflows

Accessed: 2026-05-14  
Category: Other platforms / multi-surface agent customization

## Why This Source Matters

Windsurf provides one of the clearest taxonomies for when to use rules, `AGENTS.md`, workflows, skills, and memories. This is directly useful for Skillify because the main design failure in agent customization is putting the right content into the wrong surface.

## Core Data Points

- Memories are auto-generated context retrieved when Cascade thinks it is relevant.
- Rules are manually defined behavior guidance at global, workspace, or system level.
- `AGENTS.md` files provide location-scoped rules without frontmatter.
- Workflows are prompt templates for repeatable multi-step tasks and are manually invoked through slash commands.
- Skills are multi-step procedures bundled with supporting files such as scripts, templates, and checklists.
- Skills use progressive disclosure: only `name` and `description` are shown by default; full instructions and resources load on invocation.
- Workspace skills live in `.windsurf/skills/<skill-name>/`.
- Global skills live in `~/.codeium/windsurf/skills/<skill-name>/`.
- Cross-agent skill locations are also scanned, including `.agents/skills/` and `~/.agents/skills/`.
- If Claude Code config reading is enabled, `.claude/skills/` and `~/.claude/skills/` can also be scanned.
- Rules can be stored in `.windsurf/rules/`, global rules files, `AGENTS.md`, or enterprise system locations.
- Workspace rules support activation modes:
  - `always_on`
  - `model_decision`
  - `glob`
  - `manual`
- Root `AGENTS.md` is always on.
- Subdirectory `AGENTS.md` files are treated like glob rules for that directory.
- Windsurf recommends writing durable knowledge as rules or `AGENTS.md`, not relying only on auto-generated memories.

## Design Implications For Skillify

- This taxonomy should be part of Skillify's routing logic:
  - durable behavior -> rule or `AGENTS.md`,
  - location-specific guidance -> subdirectory `AGENTS.md`,
  - manual repeatable prompt -> workflow,
  - complex procedure with resources -> skill,
  - ephemeral user/project facts -> memory.
- `.agents/skills/` is a good portable target even for Windsurf because it is scanned as a cross-agent location.
- Skill descriptions matter because Windsurf uses them for model-driven invocation.
- `model_decision` rules mirror skill trigger behavior and can reduce context cost.
- Enterprise deployments need system-level paths and read-only policy controls.

## Practical Checklist

- Use `AGENTS.md` for simple directory-scoped conventions.
- Use `.windsurf/rules/` when activation logic needs `always_on`, `glob`, `manual`, or `model_decision`.
- Use `.windsurf/skills/` or `.agents/skills/` for procedures with support files.
- Use workflows for manual slash-command procedures that do not need bundled resources.
- Keep all rules and skill bodies concise; move detailed references into support files.
