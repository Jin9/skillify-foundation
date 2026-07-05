# About agent skills

Source: https://docs.github.com/en/copilot/concepts/agents/about-agent-skills
Accessed: 2026-04-26
Category: codex-copilot / copilot skills concept
Provenance: normalized 2026-07-05 from raw HTML capture; original Accessed date preserved

## Why This Source Matters

GitHub's concept documentation for Copilot agent skills: discovery, storage locations, and behavior. Tier 1 official source for the Copilot host.


# About agent skills

Skills allow Copilot to perform specialized tasks.

## Who can use this feature?

Copilot cloud agent is available with the GitHub Copilot Pro, GitHub Copilot Pro+, GitHub Copilot Business and GitHub Copilot Enterprise plans. The agent is available in all repositories stored on GitHub, except repositories owned by managed user accounts and where it has been explicitly disabled.

GitHub Copilot CLI is available with all Copilot plans. If you receive Copilot from an organization, the Copilot CLI policy must be enabled in the organization's settings.
[Sign up for Copilot](https://github.com/features/copilot/plans?ref_product=copilot&ref_type=purchase&ref_style=button)

## In this article

 Note

 Agent skills work with Copilot cloud agent, the GitHub Copilot CLI, and agent mode in Visual Studio Code.

## About agent skills

Agent skills are folders of instructions, scripts, and resources that Copilot can load when relevant to improve its performance in specialized tasks. The Agent Skills specification is an [open standard](https://github.com/agentskills/agentskills), used by a range of different AI systems.

You can create your own skills to teach Copilot to perform tasks in a specific, repeatable way—or use skills shared online, for example in the [`anthropics/skills`](https://github.com/anthropics/skills) repository or GitHub's community-created [`github/awesome-copilot`](https://github.com/github/awesome-copilot) collection.

You can also use `gh skill` in GitHub CLI to discover and install skills from GitHub repositories. For more information, see Adding agent skills for GitHub Copilot.

Copilot supports:

- Project skills, stored in your repository (`.github/skills`, `.claude/skills`, or `.agents/skills`)
- Personal skills, stored in your home directory and shared across projects (`~/.copilot/skills`, `~/.claude/skills`, or `~/.agents/skills`)

Support for organization-level and enterprise-level skills is coming soon.

## Next steps

To create or add agent skills, see:

- Adding agent skills for GitHub Copilot
- Adding agent skills for GitHub Copilot CLI
- Copilot customization cheat sheet
