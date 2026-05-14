# VS Code Custom Agents, Prompt Files, And Skills

Source: https://code.visualstudio.com/docs/copilot/customization/custom-agents  
Related source: https://code.visualstudio.com/docs/copilot/customization/agent-skills  
Accessed: 2026-05-14  
Category: GitHub Copilot / VS Code customization

## Why This Source Matters

VS Code now distinguishes three related customization surfaces: custom agents, prompt files, and agent skills. This distinction is important for Skillify because many weak skills are actually better expressed as repo instructions, prompt files, or role agents.

## Core Data Points

- Use custom agents when the desired behavior is a persistent persona with specific tool permissions, model preferences, or handoffs between roles.
- Use prompt files for one-off or reusable prompts that do not need tool restrictions or bundled resources.
- Use agent skills for portable, reusable capabilities that may include scripts, examples, and supporting files.
- VS Code provides diagnostics for loaded custom agents, prompt files, instruction files, and skills.
- VS Code skills are stored in the same broad locations used by Copilot:
  - project: `.github/skills/`, `.claude/skills/`, `.agents/skills/`
  - personal: `~/.copilot/skills/`, `~/.claude/skills/`, `~/.agents/skills/`
- VS Code skill frontmatter includes the standard `name` and `description`, plus additional host-specific fields:
  - `argument-hint`
  - `user-invocable`
  - `disable-model-invocation`
  - experimental `context`
- Invalid skill names can fail to load silently.
- Skills can be exposed as slash commands.
- Extensions can contribute skills through the `chatSkills` contribution point.

## Design Implications For Skillify

- Skillify should classify the customization target before generating a skill:
  - standing coding conventions -> instructions or rules,
  - fixed reusable prompt -> prompt file,
  - role with tools/handoffs -> custom agent,
  - repeatable workflow with resources -> skill.
- For VS Code compatibility, Skillify can optionally emit host-specific fields, but should keep the portable core valid.
- `user-invocable: false` is useful for background knowledge skills that should activate automatically but not clutter slash menus.
- `disable-model-invocation: true` is useful for manual-only workflows, such as release checklists or destructive operations.
- Diagnostics should be part of troubleshooting guidance because silent load failure is a realistic failure mode.

## Decision Rule

If the content is mostly "who the agent should be", it belongs in a custom agent. If it is "what to do for this repeatable task, with files/scripts/templates", it belongs in a skill. If it is "how this repo works", it belongs in repo instructions or rules.

