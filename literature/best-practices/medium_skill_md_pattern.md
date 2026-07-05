# The SKILL.md Pattern: How to Write AI Agent Skills That Actually Work | by Bibek Poudel

Source: https://bibek-poudel.medium.com/the-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee
Accessed: 2026-07-05
Category: best-practices / SKILL.md pattern
Provenance: re-captured 2026-07-05 and normalized from HTML (previous capture was a 139-word Cloudflare challenge stub)

## Why This Source Matters

Community deep-dive on writing skills that trigger reliably — its core reminder is that the description alone controls activation. Tier 5 community article; inspiration and phrasing evidence only.


# The SKILL.md Pattern: How to Write AI Agent Skills That Actually Work

Bibek Poudel

15 min read

·

Feb 26, 2026

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fvote%2Fp%2F72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&user=Bibek+Poudel&userId=439dc236ad07&source=---header_actions--72a3169dd7ee---------------------clap_footer------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fvote%2Fp%2F72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&user=Bibek+Poudel&userId=439dc236ad07&source=---header_actions--72a3169dd7ee---------------------clap_footer------------------)

--

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Frepost%2Fp%2F72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&user=Bibek+Poudel&userId=439dc236ad07&source=---header_actions--72a3169dd7ee---------------------repost_header------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Frepost%2Fp%2F72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&user=Bibek+Poudel&userId=439dc236ad07&source=---header_actions--72a3169dd7ee---------------------repost_header------------------)

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fbookmark%2Fp%2F72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&source=---header_actions--72a3169dd7ee---------------------bookmark_footer------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fbookmark%2Fp%2F72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&source=---header_actions--72a3169dd7ee---------------------bookmark_footer------------------)

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2Fplans%3Fdimension%3Dpost_audio_button%26postId%3D72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&source=---header_actions--72a3169dd7ee---------------------post_audio_button------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2Fplans%3Fdimension%3Dpost_audio_button%26postId%3D72a3169dd7ee&operation=register&redirect=https%3A%2F%2Fbibek-poudel.medium.com%2Fthe-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee&source=---header_actions--72a3169dd7ee---------------------post_audio_button------------------)

Press enter or click to view image in full size

If your skill does not trigger, it is almost never the instructions. It is the description.

That is the thing most people figure out after an hour of frustration. You write a SKILL.md, drop it in the right folder, ask the agent to use it, and nothing happens. You rewrite the instructions. Still nothing. The problem was never what you wrote inside the skill. It was the two lines at the top that the agent uses to decide whether to activate it at all.

In this guide I want to go through exactly how Agent Skills work under the hood, why most people write that part wrong, and then build four skills together from easy to complex so you can see the pattern click into place. By the end you will have a README writer, a git commit generator, a code reviewer, and a full MCP-powered sprint planner, all built from scratch.

## What Is an Agent Skill?

A Skill is not a plugin. It is not a script you wire up to an API. Think of it like writing an onboarding guide for a new team member. Instead of re-explaining your workflows and preferences in every conversation, you package them once and the agent picks them up automatically whenever your request matches.

At its core, a skill is just a folder:

```
your-skill-name/
├── SKILL.md          # Required: instructions + metadata
├── scripts/          # Optional: executable code the agent runs
├── references/       # Optional: docs loaded only when needed
└── assets/           # Optional: templates, images, fonts
```

The only required file is `SKILL.md`. Everything else is optional but becomes important as skills grow in complexity.

What makes this particularly useful right now is that the SKILL.md format is an **open standard**, published by Anthropic at agentskills.io in December 2025. It works across Claude Code, OpenAI Codex, and OpenClaw. While the format is standardized, each platform implements discovery and tooling slightly differently. Think shared language, not identical behavior. A skill that works on Claude Code will very likely work on Codex, but runtime behaviors like session snapshotting, tool permissions, and invocation modes differ between platforms.

## Where Skills Live

Before writing anything, you need to know where to put it. Each platform loads skills from specific locations, and the location defines the scope.

**Claude Code:**

Location Scope `~/.claude/skills/` Personal, available across all your projects `.claude/skills/` Project-level, shared with your team via git

**OpenAI Codex:**

Location Scope `~/.codex/skills/` User-level, applies to any repo you work in `.codex/skills/` Repo-level, checked into git

**OpenClaw:**

Location Scope `~/.openclaw/skills/` Global, available to all configured agents Per-agent workspace Scoped to a specific agent only

When two skills share the same name, the higher-precedence location wins. A project-level skill overrides a personal one with the same name. This lets teams define defaults that individuals can override for their own setups.

## How the Three-Level Loading System Works

This is the part most people skip, and it explains almost every problem with skills not triggering or consuming too much context.

Skills use **progressive disclosure**: a three-level loading system where content is pulled into context only as it is needed.

Press enter or click to view image in full size

**Level 1: Metadata (always loaded, ~100 tokens per skill)**

At startup, the agent reads only the `name` and `description` from every installed skill's YAML frontmatter. Nothing else. This compact listing goes into the system prompt so the agent knows what skills exist and when to use them. The practical implication: you can install many skills without a context penalty.

**Level 2: Instructions (loaded when triggered, under 5k tokens)**

When the agent decides a skill is relevant, it reads the full body of `SKILL.md` into context using a bash call. Only at this point do your actual instructions get loaded.

**Level 3: Referenced files and scripts (loaded on demand, effectively unlimited)**

If the skill body references other files, the agent reads those only when it needs them. Scripts can be executed without being read into context at all. This is what makes skills scalable: the token cost at idle is zero regardless of how much content you bundle.

Here is what this looks like in sequence for a real request:

```
1. Session starts
   --> Agent loads: name + description from every skill (~100 tokens each)
2. User asks: "Can you write a README for this project?"
   --> Agent reads: readme-writer/SKILL.md full body (Level 2)
3. SKILL.md references a style guide file
   --> Agent reads: readme-writer/references/style.md (Level 3)
4. SKILL.md includes a validation script
   --> Agent executes: scripts/validate.sh
   (runs without being read into context)
```

## Skills vs. Slash Commands

Both Claude Code and Codex support two invocation modes. The agent can activate a skill automatically when your request matches the description (implicit invocation), or you can call it directly (explicit invocation).

In **Claude Code**, skills appear in the slash command menu by default. You can invoke one directly with `/skill-name`, or just describe what you want and Claude will activate the relevant skill automatically:

```
# Direct invocation in Claude Code
/readme-writer
# Or just describe the task and Claude activates it automatically
Can you write a README for this project?
```

In **Codex CLI**, you mention a skill with `$` prefix or use the `/skills` selector:

```
# Explicit invocation in Codex CLI
$readme-writer document this project
# Or use the skill selector
/skills
```

Even though both platforms support explicit invocation, a well-written description still matters enormously. It is what drives automatic activation, the behavior that makes skills feel like a natural extension of how you already work rather than a command you have to remember to type.

## The Most Common Mistake

The `description` field in your YAML frontmatter is not for humans. It is the trigger condition the agent uses when deciding whether to activate your skill.

Here is the structure that works:

```
[What the skill does] + [When to use it, with specific trigger phrases]
```

Bad:

```
description: Helps with documents.
```

Also bad, because it describes what but not when:

```
description: Creates sophisticated multi-page documentation with advanced
formatting.
```

Good:

```
description: Creates and writes professional README.md files for software
projects. Use when user asks to "write a README", "create a readme",
"document this project", "generate project documentation", or "help me write
a README.md".
```

The agentskills.io spec defines these constraints:

- `name`: lowercase letters, numbers, and hyphens only, max 64 characters, must not start or end with a hyphen, no consecutive hyphens
- `description`: max 1024 characters, must describe both what the skill does and when to use it
- The file must be named exactly `SKILL.md`, case-sensitive
- Avoid XML angle brackets (`<` or `>`) in frontmatter as they can inject unintended instructions into the system prompt

Some platforms add conventions on top of these. When in doubt, check the platform-specific docs alongside the base spec at agentskills.io/specification.

## Skill 1: README Writer

This is a practical first skill. Almost every developer has written a README at some point, usually by hand, usually inconsistently. This skill teaches the agent your preferred structure and writes the file to disk automatically.

## Setup

```
mkdir -p ~/.claude/skills/readme-writer
```

## SKILL.md

```

---
name: readme-writer
description: Creates and writes professional README.md files for software projects.
Use when user asks to "write a README", "create a readme", "document this project",
"generate project documentation", or "help me write a README.md". Works from a
project description, existing code, or both.
---

# README Writer

## Overview

Generate a complete, professional README.md file and write it to disk. The output
should be clear enough for a first-time contributor to understand the project,
set it up locally, and start contributing.

## Step 1: Gather project context

Look for context in the codebase before asking the user:

```bash
ls -la
cat package.json 2>/dev/null || cat pyproject.toml 2>/dev/null || \
  cat go.mod 2>/dev/null || echo "No manifest found"
ls .env.example .env.sample 2>/dev/null || echo "No env example found"
```

Gather:
- What does this project do? (1-2 sentence summary)
- What language and main frameworks does it use?
- How do you install and run it?
- Are there environment variables needed?
- Is there a LICENSE file?

## Step 2: Write the README

Use this structure. Only include sections that are relevant. Do not add empty sections.

```
# Project Name

One clear sentence describing what this project does and who it is for.

## Features
- Feature one (be specific)
- Feature two

## Prerequisites
List what needs to be installed. Include version requirements if important.

## Installation
Step-by-step setup. Every command must be copy-pasteable.

```bash
git clone https://github.com/username/project
cd project
npm install
```

## Configuration
If the project needs environment variables, show an example:

```bash
cp .env.example .env
```

Then explain each variable the user needs to set manually.

## Usage
Show the most common use case first.

```bash
npm run dev
```

## License
[MIT](LICENSE)
```

## Step 3: Write the file to disk

Once the content is ready, write it:

```bash
cat > README.md << 'EOF'
[full readme content]
EOF
```

Confirm it was written:

```bash
echo "README.md written: $(wc -l < README.md) lines"
```

## Step 4: Quality check

Before finishing, verify:
- [ ] No placeholder text like "[your description here]" remains
- [ ] Every command in the Installation section is accurate for this project
- [ ] Prerequisites match what the project actually needs
- [ ] License section matches the LICENSE file if one exists
```
```

## Test it

Go to any project folder and ask:

```
Can you write a README for this project?
```

The agent will inspect the codebase, write the README, save it as `README.md`, and confirm with a line count. No copy-pasting required.

## Skill 2: Git Commit Message Generator

This skill shows how to write trigger phrases that cover the different ways a developer might ask for the same thing.

### Setup

```
mkdir -p ~/.claude/skills/git-commit-writer
```

## SKILL.md

```

---
name: git-commit-writer
description: Generates standardized git commit messages following conventional commits spec.
Use when user asks to "write a commit message", "help me commit", "summarize my changes",
"what should my commit say", or "draft a commit". Analyzes staged diffs and change
descriptions to produce type(scope): description format messages.
---

# Git Commit Message Writer

## Format

```
type(scope): short description

[optional body]

[optional footer]
```

Allowed types: feat, fix, docs, style, refactor, test, chore, perf, ci, build

## Instructions

### Step 1: Get the diff

```bash
git diff --staged
```

If nothing is staged:

```bash
git diff HEAD
```

### Step 2: Analyze the changes

Look for:
- What files changed and what category they belong to
- Whether this adds new functionality (feat), fixes a bug (fix), or updates docs/config/tests
- The scope: which module, component, or area is affected

### Step 3: Write the message

- Keep the subject line under 72 characters
- Use imperative mood: "add feature" not "added feature"
- Do not end the subject line with a period
- Add a body if the change needs more context than the subject allows

### Quality check

- [ ] Type is one of the allowed types
- [ ] Subject line is under 72 characters
- [ ] Imperative mood is used
- [ ] Scope is specific enough to be useful

## Examples

```
feat(auth): add OAuth2 login with Google

Implements Google OAuth2 flow using the existing session management
system. Users can now sign in with their Google account.

Closes #142
```

```
fix(api): handle null response from payment provider
```

```
docs(readme): update local setup instructions for Node 22
```
```
```

## Skill 3: Code Reviewer (Multi-File)

This skill shows when to split content across multiple files. The process stays in `SKILL.md`. The detailed criteria live in a reference file loaded only during an actual review. This is the right way to structure complex skills.

```
mkdir -p ~/.claude/skills/code-reviewer/references
```

## SKILL.md

```

---
name: code-reviewer
description: Conducts structured code reviews with categorized feedback. Use when user asks
to "review this code", "check my PR", "look over this function", or "give me feedback on
this implementation". Produces structured output with blocking issues separate from suggestions.
---

# Code Reviewer

## Review Process

### Step 1: Understand context

Before reviewing, establish:
- What is this code supposed to do?
- What language and framework is it using?
- Is this a new feature, a bug fix, or a refactor?

### Step 2: Run the review

For detailed review criteria by category, see [references/criteria.md](references/criteria.md).

Work through each category in order. Do not skip categories even if they seem unlikely to have issues.

### Step 3: Structure the output

```
## Summary
[2-3 sentence overview and overall assessment]

## Blocking Issues
[Issues that must be fixed: security vulnerabilities, logic errors, data loss risks.
If none, write "None found."]

## Suggestions
[Non-blocking improvements numbered. Include where, why, and how to fix each.]

## Positive Notes
[What the code does well. Always include at least one.]
```
```

```

### references/criteria.md

```
# Review Criteria

## Security (Check First)
- SQL injection: are user inputs parameterized?
- XSS: is output properly escaped before rendering?
- Auth checks: are protected routes actually protected?
- Secrets: are API keys or credentials hardcoded anywhere?
- Input validation: is validation happening server-side?

## Correctness
- Does the logic match the stated intent?
- Are edge cases handled: empty arrays, null values, zero, negative numbers?
- Are error states surfaced correctly?
- Are async operations awaited properly?

## Readability
- Can a new team member understand this in 5 minutes?
- Are variable and function names descriptive?
- Are functions doing one thing or multiple things?

## Performance
- Are there obvious N+1 query patterns?
- Are expensive operations inside loops that could be outside?

## Tests
- Are there tests for the new behavior?
- Are edge cases tested, not just the happy path?
```

The `SKILL.md` body stays under 40 lines. The detailed criteria live in `references/criteria.md` and are loaded only when a review is running. This keeps Level 2 lean while the agent still has access to everything it needs at Level 3.

## Skill 4: Linear Sprint Planner with MCP

This is a Category 3 skill: MCP Enhancement. The MCP server gives the agent access to Linear’s API. The skill gives it the knowledge of how to use that access reliably and consistently. Without the skill, users connect the MCP but still have to figure out every step. With the skill, the entire workflow runs from one sentence.

## Setup

```
mkdir -p ~/.claude/skills/linear-sprint-planner/references
```

## SKILL.md

```
---
name: linear-sprint-planner
description: Automates Linear sprint planning including cycle creation, backlog triage,
and task assignment. Use when user says "plan the sprint", "set up the next cycle",
"help me prioritize the backlog", or "create sprint tasks in Linear". Requires Linear
MCP server to be connected.
metadata:
  mcp-server: linear
  version: 1.0.0
---
# Linear Sprint Planner

## Prerequisites

Verify the Linear MCP server is connected. If not available, tell the user to connect
it in their MCP settings before continuing.

## Process

### Step 1: Gather current state

Fetch from Linear in sequence:
1. Current active cycle and completion percentage
2. All backlog issues (status: Backlog or Todo, not assigned to any cycle)
3. Team members and current workload
4. Any issues marked high priority

See [references/linear-api.md](references/linear-api.md) for pagination and rate limit handling.

### Step 2: Analyze capacity

- Count team members participating in the sprint
- Estimate points (default: 10 per person per week unless historical velocity exists)
- Subtract planned time off

Present a summary before proceeding:
Team capacity: - [N] engineers x [X] points = [Total] available - Carrying over: [X] points - Net new capacity: [X] points
### Step 3: Prioritize backlog

Sort in this order:
1. P0/P1 bugs and blockers (always include)
2. Items explicitly flagged by the user
3. Items that unblock other teams
4. Features by product priority
5. Tech debt

Do not exceed capacity by more than 10% unless the user asks.

### Step 4: Present for approval

Before creating anything in Linear, show the proposed plan:
Proposed Sprint [N]: [Date Range] Capacity: [X] points Issues: - [ISSUE-123] Fix payment timeout (P0) - 3 pts - [ISSUE-456] Add CSV export (P2) - 5 pts - [ISSUE-789] Refactor auth middleware - 2 pts Total: [X] / [X] points Shall I create this cycle and assign these issues?
Always wait for confirmation before making changes.

### Step 5: Create and confirm

Once approved:
1. Create the cycle with the agreed date range
2. Add each issue to the cycle
3. Update assignments if the user specified owners
4. Return a summary with the Linear cycle link

See [references/error-handling.md](references/error-handling.md) if any API calls fail.
```

## references/linear-api.md

```
# Linear API Patterns
## Fetching backlog
Paginate in batches of 50. Check pageInfo.hasNextPage and use the after cursor.
Filter for: status [Backlog, Todo], cycle: null.
## Creating a cycle
Required: name, startsAt, endsAt, teamId
## Adding issues
Use issueUpdate mutation to set cycleId. Batch where possible.
## Rate limiting
On 429: wait 1 second, retry once. If it fails again, report to user and continue.
```

## references/error-handling.md

```
# Error Handling
## MCP connection errors
Tell the user: "Linear MCP appears to be disconnected. Please reconnect before
running sprint planning." Do not proceed.
## Missing data
Continue with available data. Note what is missing in the final summary.
## Cycle creation failure
Do not attempt to add issues. Report the error and suggest checking permissions.
```

## Test it

```
Help me plan the next sprint in Linear
```

The agent fetches live data, proposes a plan, waits for approval, and executes. One sentence triggers the entire workflow.

One important thing to internalize before moving on: skills do not guarantee execution. The model still decides whether to follow the instructions. Think of them as structured guidance that dramatically increases consistency, not deterministic automation. If the model goes off-script, the fix is almost always improving the instructions or the description, not debugging runtime behavior.

## How All Four Skills Relate

Here is a visual of the four skills and where they sit on the complexity spectrum:

Press enter or click to view image in full size

Start with `git-commit-writer`. When you need file output, use the `readme-writer` pattern. When your skill gets too long, split it like `code-reviewer`. When you need external tool coordination, use the `linear-sprint-planner` pattern.

## The `allowed-tools` Field

One frontmatter field most people never use is `allowed-tools`. It restricts which tools the agent can call when a skill is active.

This is useful for read-only skills where you do not want the agent accidentally writing or executing anything:

```
---
name: log-analyzer
description: Analyzes application log files to identify errors and patterns. Use when
user says "check the logs", "what errors are in my logs", or "analyze this log file".
allowed-tools: Read, Grep, Glob
---
# Log Analyzer
1. Use Glob to find log files: *.log, logs/*.log, /var/log/*.log
2. Use Grep to search for error patterns: ERROR, FATAL, Exception, Traceback
3. Group errors by type and frequency
4. Summarize with the most frequent errors first
```

With this active, the agent cannot execute shell commands, write files, or make any external calls. A simple way to add safety guarantees to observational skills.

Note: `allowed-tools` is marked experimental in the Agent Skills spec, and support varies between agent implementations. It is well-supported in Claude Code today.

## Sharing Skills with Your Team

The cleanest approach is to commit project skills to your repo:

```
mkdir -p .claude/skills/readme-writer
# add SKILL.md, then:
git add .claude/skills/
git commit -m "Add readme-writer skill for team"
git push
```

When teammates pull the repo, the skill is immediately available with no separate installation step. In Codex the equivalent path is `.codex/skills/`.

## Debugging When a Skill Does Not Trigger

**1. Check your description**

The most common issue. Add more specific trigger phrases that match how users actually phrase requests.

**2. Check your file path**

```
# Claude Code personal skill
ls ~/.claude/skills/your-skill/SKILL.md
# Claude Code project skill
ls .claude/skills/your-skill/SKILL.md
# Codex
ls ~/.codex/skills/your-skill/SKILL.md
```

**3. Check YAML syntax**

Invalid YAML silently prevents loading. Frontmatter must start on line 1 with `---` and close with another `---`.

**4. Restart the session**

Skills are snapshotted at session start. Edits made during a running session require a restart to take effect.

**5. Run in debug mode**

```
claude --debug
```

**6. Test with explicit invocation**

```
Use the readme-writer skill to document this project
```

If it works when invoked explicitly but not automatically, the description needs more specific trigger phrases.

## Security

Skills can bundle executable code and instructions that control agent behavior. That same power makes malicious skills dangerous.

Install skills only from trusted sources. Before installing any community skill, read every file in the folder, especially anything in `scripts/`. Pay attention to instructions that tell the agent to make outbound network calls or send data to external services.

Cisco researchers have warned about skills being used for silent data exfiltration via prompt injection. Security audits scanning thousands of community skills have found a meaningful fraction with critical vulnerabilities including credential theft and malware. Snyk has published findings on this specifically. ClawHub has a VirusTotal integration you can use to check skills before installing, but manual review is still worthwhile for anything with broad permissions.

Practical rules:

- Never install a skill that asks you to paste secrets into chat
- Always read `scripts/` before installing
- Be skeptical of skills that make outbound network calls in setup instructions
- Prefer skills from official sources: Anthropic’s or OpenAI’s skills repos, or your own team’s

## Wrapping Up

The three-level loading system is the core concept. Everything else follows from it. Level 1 is your trigger. Level 2 is your runbook. Level 3 is your reference library. Get those right and you can build anything from a single Markdown file to a multi-step workflow coordinating external APIs and executing code.

The SKILL.md format is also no longer just a Claude feature. OpenAI adopted it for Codex. OpenClaw uses it as its core plugin format. Skills you write today are portable across all three platforms. That portability is worth investing in.

Feel free to modify and expand the skills built in this guide to suit your specific workflows. The official skills repos are good places to study well-written examples from both Anthropic and OpenAI:

- Anthropic skills: [https://github.com/anthropics/skills](https://github.com/anthropics/skills)
- OpenAI skills: [https://github.com/openai/skills](https://github.com/openai/skills)

*References:*

- *Anthropic Agent Skills Overview: *[*https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview*](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- *Claude Code Skills Guide: *[*https://code.claude.com/docs/en/skills*](https://code.claude.com/docs/en/skills)
- *OpenAI Codex Skills: *[*https://developers.openai.com/codex/skills*](https://developers.openai.com/codex/skills)
- *Agent Skills Open Standard: *[*https://agentskills.io/what-are-skills*](https://agentskills.io/what-are-skills)
- *The Complete Guide to Building Skills for Claude: *[*https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf*](https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf)
- *Equipping Agents for the Real World: *[*https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills*](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills)
