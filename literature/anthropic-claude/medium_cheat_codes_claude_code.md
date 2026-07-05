# Agent Skills: The Cheat Codes for Claude Code

Source: https://medium.com/jonathans-musings/agent-skills-the-cheat-codes-for-claude-code-b8679f0c3c4d
Accessed: 2026-07-05
Category: anthropic-claude / community practice
Provenance: re-captured 2026-07-05 via content-extraction proxy (direct fetch is Cloudflare-gated); may be abridged relative to the original page; previous capture was a 139-word Cloudflare challenge stub

## Why This Source Matters

Practitioner field report on real Claude Code skill usage: concrete examples of a cross-model review skill, an output-compression skill, and company-internal CLI skills. Its "what to skill-ify" heuristics (repetitive workflows, high-context tasks, output formatting, safety rails) are useful community evidence for trigger design and scope selection. Tier 5 community source — inspiration and validation only.

---

_How I customized my AI coding workflow to be faster, cheaper, and more effective_

I spend most of my engineering hours inside Claude Code. It's become my primary interface for writing software — I've barely hand-written any code since joining Datadog. Discovering **skills** helped me unlock a lot of value.

Skills are modular instructions you can plug into Claude Code to extend its capabilities. Think of them like plugins or extensions, but for an AI agent. Some are community-built. Some you write yourself. And once you start customizing your agent workflow with the right skills, you'll wonder how you ever worked without them.

Here are three skills that fundamentally changed how I work.

## 1. The Codex Review Plugin: Your Second (and Third) Opinion

Code review has always been a bottleneck. You write something, you think it's good, you submit a PR, and then you wait for a teammate to find the obvious bug you missed.

The Codex Review Plugin gives Claude Code the ability to invoke OpenAI's Codex agent for code review — right from your terminal. But here's what makes it actually useful: it has two modes.

**Normal review** (`/codex:review`) does what you'd expect. It examines your changes, checks for bugs, suggests improvements, validates logic. Standard stuff.

**Adversarial review** (`/codex:adversarial-review`) is where it gets interesting. Instead of helpful suggestions, Codex actively tries to break your code. It probes edge cases. It questions architectural decisions. It plays devil's advocate on your design choices.

I now run adversarial review before most PRs. It's caught issues that would've made it through three rounds of human review — race conditions, missed null checks, assumptions that only held in happy-path scenarios.

The mental shift matters too. When you know something is going to attack your code, you plan more defensively from the start.

## 2. Caveman: 75% Token Savings, Zero Information Loss

Here's an uncomfortable truth about AI coding assistants: most of the tokens they generate are fluff.

_"Certainly! I'd be happy to help you refactor that function. Let me walk you through the changes step by step…"_

That preamble costs money. On Claude, you're paying per token — both input and output. A verbose agent that loves to explain itself can burn through your usage limits fast.

Enter Caveman.

Caveman is a skill that makes your agent communicate like… well, a caveman. Short sentences. No filler. Just the information you need.

Instead of:

> "I've successfully completed the refactoring of the authentication module. The changes include updating the token validation logic to handle edge cases more gracefully, adding appropriate error handling, and ensuring backwards compatibility with the existing API contracts."

You get:

> "Done. Token validation updated. Error handling added. API compatible."

Same information. 75% fewer tokens.

The first day feels weird. The second day feels efficient. By the third day, you realize all that verbosity was noise you were trained to tolerate.

Caveman also comes with a memory compression tool. Your `CLAUDE.md` file—which loads every session—can be compressed into caveman-speak too. Fewer input tokens means faster responses and lower costs.

## 3. Personal Datadog Skills: Making Observability Second Nature

This is where skills get personal.

I work at Datadog. We have a CLI tool called `pup` for interacting with our platform—querying metrics, tailing logs, managing dashboards. Powerful, but with a learning curve.

So I built custom skills that teach Claude Code how to use `pup` effectively. Now I can say things like:

* "Show me error rates for the experiments service over the last hour"
* "Tail logs from staging where status is 500"
* "How does our experimentation pipeline dashboard look in the past day? Can you group errors by organization ID?"

Claude constructs the right `pup` commands, runs them, and interprets the results. What used to require flipping between documentation tabs now happens conversationally.

The compounding effect is real. Every time I add a new workflow to my Datadog skills — how to query specific metric namespaces, how to filter by our team's service tags, how to format alerts for our on-call channel — my future self gets faster.

This is the point most people miss about agent skills: **they're not just for common tasks. The biggest wins come from encoding your specific environment, your team's conventions, your company's tools.**

## The Meta-Skill: Knowing What to Skill-ify

After a few months of building and using skills, a pattern emerged. The best candidates for skills are:

**Repetitive workflows** — Anything you do weekly that requires remembering syntax or options. Code review, deployment checklists, monitoring queries.

**High-context tasks** — Things that require knowledge your agent doesn't have by default. Your company's CLI tools, your team's naming conventions, your infrastructure quirks.

**Output formatting** — How you want information presented. Caveman is an extreme example, but even small tweaks to how your agent communicates can reduce cognitive load.

**Safety rails** — Things you want the agent to always check. The adversarial review is a safety rail. So is a skill that reminds Claude to run tests before committing.

## Getting Started

If you're using Claude Code (or similar AI coding tools), start simple:

1. **Install Caveman.** The token savings alone justify the 5-minute setup. Your wallet will thank you.
2. **Try the Codex review plugin.** Run adversarial review on your next PR. See what it catches.
3. **Identify one repetitive task** you do with an internal tool. Write a skill that teaches Claude how to do it. Even a rough first version will save you time.

Skills are how you stop using a generic AI assistant and start working with an AI assistant that knows _your_ job. The difference is everything.

---

Jonathan Fulton is a Staff Engineer at Datadog working on experimentation infrastructure. I write about AI-assisted engineering, A/B testing, and building software at velocity. Previously: Eppo, ID.me, Storyblocks.
