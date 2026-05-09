Here are the best source links for creating a strong skill.md / agent-skill template.

Priority	Source	Why it matters
1	OpenAI — Agent Skills for Codex: https://developers.openai.com/codex/skills	Best official source for Codex skills. It says skills should be focused, use clear inputs/outputs, prefer instructions unless scripts are needed, and test trigger behavior.
2	OpenAI — Skill Creator example: https://github.com/openai/skills/blob/main/skills/.system/skill-creator/SKILL.md	Best concrete SKILL.md example from OpenAI. Good for learning structure, wording, and how to write modular skill instructions.
3	OpenAI — Codex Best Practices: https://developers.openai.com/codex/learn/best-practices	Useful for the bigger workflow: planning first, reusable guidance, validation loops, MCP, and automation.
4	OpenAI — AGENTS.md guide: https://developers.openai.com/codex/guides/agents-md	Not exactly skill.md, but important. Use AGENTS.md for repo-level persistent rules; use SKILL.md for reusable task workflows.
5	OpenAI — Codex Customization: https://developers.openai.com/codex/concepts/customization	Explains the relationship between AGENTS.md, skills, memory, MCP, and subagents. Good for designing your full agent system.
6	Anthropic — Complete Guide to Building Skills for Claude: https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf	Best official Claude skill guide. Covers planning, structure, testing, iteration, sharing, and troubleshooting.
7	Anthropic — Claude Code Best Practices: https://code.claude.com/docs/en/best-practices	Useful for how Claude behaves in real codebases: context setup, planning, testing, and iterative workflows.
8	Anthropic — Claude Code Subagents: https://code.claude.com/docs/en/sub-agents	Important if you split skills into BA / Tech Lead / Dev / QA agents. It explains how descriptions help Claude decide when to delegate.
9	AGENTS.md open format: https://agents.md	Good simple standard for repo-level agent instructions. Useful companion to skill.md.
10	Community Claude Skills notes: https://github.com/shanraisshan/claude-code-best-practice/blob/main/best-practice/claude-skills.md	Practical community notes about fields like description, when_to_use, and trigger behavior. Use as secondary reference, not primary source.

My recommended reading order:

1. OpenAI Agent Skills
2. OpenAI Skill Creator example
3. Anthropic Complete Guide to Building Skills
4. OpenAI AGENTS.md guide
5. Claude Code Best Practices
6. Claude Subagents

For your use case, the key principle is:

AGENTS.md = repo-wide rules
SKILL.md = repeatable workflow / domain expertise
Subagent = delegated role such as Architect, Backend Dev, Reviewer, QA
MCP = external tools such as GitHub, draw.io, Jira, DB, docs

--

Yes. Beyond the 10 links above, these are useful additional sources — especially agentskills.io.

Priority	Source	Why useful
1	AgentSkills.io — Best Practices: https://agentskills.io/skill-creation/best-practices	Good independent guide. It emphasizes keeping SKILL.md lean, under about 500 lines / 5,000 tokens, and moving deeper references into separate files.  ￼
2	Anthropic API Docs — Skill authoring best practices: https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices	Very useful official Claude source. It covers naming conventions, including using gerund-style skill names like processing-pdfs, testing-code, and writing-documentation.  ￼
3	Anthropic API Docs — Agent Skills overview: https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview	Explains how SKILL.md works as the instruction layer and how extra files can be loaded only when needed.  ￼
4	OpenAI API — Skills guide: https://developers.openai.com/api/docs/guides/tools-skills	Useful if you want skills that align with the open Agent Skills standard. It describes a skill as a versioned bundle with SKILL.md front matter plus instructions.  ￼
5	OpenAI Cookbook — Skills in API: https://developers.openai.com/cookbook/examples/skills_in_api	Practical implementation example for API/container execution. Good when you want skills with scripts/assets, not just Markdown instructions.  ￼
6	OpenAI Skills Catalog GitHub: https://github.com/openai/skills	Good real-world examples of skill folder structure and reusable task packaging.  ￼
7	mgechev/skills-best-practices: https://github.com/mgechev/skills-best-practices	Community guide focused on professional-grade skills, validation, and keeping the context window lean.  ￼
8	Spring AI Agentic Patterns — Agent Skills: https://spring.io/blog/2026/01/13/spring-ai-generic-agent-skills	Good LLM-agnostic explanation. Useful if you want skills to work across Codex, Claude, OpenCode, or custom agents.  ￼
9	Strapi — What Are Agent Skills and How To Use Them: https://strapi.io/blog/what-are-agent-skills-and-how-to-use-them	Beginner-friendly explanation of why skills are useful for repeated workflows and specialized knowledge.  ￼
10	Deep Dive SKILL.md — Medium: https://abvijaykumar.medium.com/deep-dive-skill-md-part-1-2-09fc9a536996	Practical article explaining skill folder structure: SKILL.md, scripts, references, and assets. Treat as secondary/community reference.  ￼

My recommended extra reading order:

1. Anthropic Skill authoring best practices
2. AgentSkills.io best practices
3. OpenAI API Skills guide
4. OpenAI Cookbook Skills in API
5. OpenAI Skills Catalog examples
6. mgechev/skills-best-practices

For your case, I would use:

Official source of truth:
- OpenAI Codex Skills
- OpenAI API Skills
- Anthropic Agent Skills docs
Practical implementation examples:
- OpenAI Skills GitHub
- OpenAI Cookbook Skills in API
Community validation:
- AgentSkills.io
- mgechev/skills-best-practices
- Spring AI Agentic Patterns

agentskills.io is worth reading, but I would treat it as secondary guidance, not the source of truth. For production skill.md, anchor your template on OpenAI + Anthropic official docs, then use AgentSkills.io to improve structure and maintainability.

--

Yes. For deeper research, I would group the sources into official docs, spec/reference, examples/catalogs, repo-level instruction docs, agent-design theory, and community field reports.

A. Primary / official sources

Source	Link	Use for
OpenAI Codex — Agent Skills	https://developers.openai.com/codex/skills	Codex skill best practices, trigger behavior, focused skills.
OpenAI API — Tools / Skills	https://developers.openai.com/api/docs/guides/tools-skills	OpenAI API usage of skills and open skill format.
OpenAI Cookbook — Skills in API	https://developers.openai.com/cookbook/examples/skills_in_api	Practical API/container example.
OpenAI Skills GitHub	https://github.com/openai/skills	Official examples and folder structure.
OpenAI Skill Creator example	https://github.com/openai/skills/blob/main/skills/.system/skill-creator/SKILL.md	Concrete SKILL.md authoring pattern.
OpenAI Codex Best Practices	https://developers.openai.com/codex/learn/best-practices	Broader Codex workflow: planning, validation, MCP, automation.
OpenAI AGENTS.md guide	https://developers.openai.com/codex/guides/agents-md	Repo-level instructions and layering behavior.
OpenAI Codex Customization	https://developers.openai.com/codex/concepts/customization	Relationship between AGENTS.md, skills, MCP, memory, subagents.
Anthropic Agent Skills overview	https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview	Claude Skill concept, SKILL.md, extra files, Skills API.
Anthropic Skill authoring best practices	https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices	Naming, description quality, concise instructions, testing.
Anthropic Skills GitHub	https://github.com/anthropics/skills	Official Anthropic skill examples.
Claude Code Best Practices	https://code.claude.com/docs/en/best-practices	Real coding-agent workflow with Claude Code.
Claude Code Subagents	https://code.claude.com/docs/en/sub-agents	Useful when splitting Architect / Dev / Reviewer / QA.
Anthropic — Building Effective Agents	https://www.anthropic.com/research/building-effective-agents	Agentic patterns: chaining, routing, orchestrator-workers, evaluator-optimizer.
Anthropic — Effective Context Engineering for AI Agents	https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents	Context engineering, memory, context discipline.
GitHub Docs — About agent skills	https://docs.github.com/en/copilot/concepts/agents/about-agent-skills	GitHub Copilot skill behavior.
VS Code Docs — Use Agent Skills	https://code.visualstudio.com/docs/copilot/customization/agent-skills	Skill placement and usage in VS Code/Copilot.
GitHub Changelog — Manage agent skills with GitHub CLI	https://github.blog/changelog/2026-04-16-manage-agent-skills-with-github-cli/	gh skill, version pinning, supply-chain concerns.

B. Open standard / reference sites

Source	Link	Use for
AgentSkills.io home	https://agentskills.io/home	Open format overview.
AgentSkills.io best practices	https://agentskills.io/skill-creation/best-practices	Lean SKILL.md, references, trigger design.
AGENTS.md open format	https://agents.md	Simple repo-level instruction standard.
Awesome Agent Skills — VoltAgent	https://github.com/VoltAgent/awesome-agent-skills	Large catalog of community skills.
Awesome Agent Skills — scienceaix	https://github.com/scienceaix/agentskills	Curated references and learning list.
Anthony Fu skills repo	https://github.com/antfu/skills	Good example of curated skills and metadata conventions.

C. Community / practical implementation sources

Source	Link	Use for
mgechev/skills-best-practices	https://github.com/mgechev/skills-best-practices	Practical professional-grade skill guidance.
Shanraisshan Claude Code best practice — Claude Skills	https://github.com/shanraisshan/claude-code-best-practice/blob/main/best-practice/claude-skills.md	Claude-specific skill notes and field practice.
Soham Kamani — Mastering AI Agent Skills	https://www.sohamkamani.com/ai/what-are-ai-agents-and-how-to-make-your-own/	Explains folder conventions across Codex, Claude, Copilot, OpenCode.
Groff.dev — Implementing CLAUDE.md and Agent Skills	https://www.groff.dev/blog/implementing-claude-md-agent-skills	Good 3-tier pattern: global file, repo file, skills.
Medium — Deep Dive SKILL.md	https://abvijaykumar.medium.com/deep-dive-skill-md-part-1-2-09fc9a536996	Practical explanation of skill folder structure.
Medium — The SKILL.md Pattern	https://bibek-poudel.medium.com/the-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee	Good reminder that description controls trigger behavior.
Medium — Agent Skills: Cheat Codes for Claude Code	https://medium.com/jonathans-musings/agent-skills-the-cheat-codes-for-claude-code-b8679f0c3c4d	Practical examples of custom Claude Code skills.
Strapi — What Are Agent Skills and How To Use Them	https://strapi.io/blog/what-are-agent-skills-and-how-to-use-them	Beginner-friendly explanation.
Spring AI — Agent Skills pattern	https://spring.io/blog/2026/01/13/spring-ai-generic-agent-skills	LLM-agnostic framing for generic agent skills.
Snyk — Top Claude Skills for Developers	https://snyk.io/articles/top-claude-skills-developers/	Developer-oriented examples: security, Terraform, PR review.
DataCamp — Top Agent Skills	https://www.datacamp.com/blog/top-agent-skills	Broad catalog-style overview.

D. Sources for adjacent concepts

These are not only about SKILL.md, but they help you design better skills.

Topic	Sources to study
Repo-level instructions	OpenAI AGENTS.md guide, agents.md
Claude project memory	Claude Code best practices, CLAUDE.md articles
Subagents / squad design	Claude Code Subagents, Anthropic Building Effective Agents
MCP / external tools	OpenAI Codex customization, Anthropic tool-use docs, MCP docs
Context management	Anthropic Effective Context Engineering, OpenAI Codex best practices
Supply-chain risk	GitHub gh skill changelog, GitHub agent skills docs
Skill discoverability	Anthropic skill best practices, AgentSkills.io best practices, SKILL.md Pattern article

Recommended research order

1. OpenAI Codex Agent Skills
2. Anthropic Agent Skills overview
3. Anthropic Skill authoring best practices
4. AgentSkills.io home + best practices
5. OpenAI Skills GitHub examples
6. Anthropic Skills GitHub examples
7. OpenAI AGENTS.md guide
8. GitHub Copilot Agent Skills docs
9. Claude Code Best Practices
10. Claude Subagents
11. Anthropic Building Effective Agents
12. Anthropic Effective Context Engineering
13. Community examples: mgechev, shanraisshan, sohamkamani, groff.dev

My filter for source trust

Tier	Use
Tier 1: Official docs	Use as source of truth for syntax, folder paths, behavior, API support.
Tier 2: Open standard / reference	Use for portability across Claude, Codex, Copilot, Cursor, OpenCode.
Tier 3: Official examples	Use to copy structure and patterns.
Tier 4: Community guides	Use for tricks, naming, trigger behavior, and real-world mistakes.
Tier 5: Blog/listicle/catalog sources	Use for inspiration only; verify before adopting.

For your case, I’d anchor the template on OpenAI + Anthropic + AgentSkills.io, then use community sources only to improve naming, folder layout, and trigger behavior.