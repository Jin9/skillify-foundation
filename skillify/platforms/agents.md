# Skill: Skillify

> **Source**: Adapted from `skillify/SKILL.md` for OpenAI Codex (AGENTS.md format).
> This is a generated file. Edit the canonical `SKILL.md` and re-run `platforms/install.sh` to regenerate.

## Purpose

Convert source material into a reusable agent skill (SKILL.md format) that is easy to discover, concise enough to load, and specific enough to guide reliable execution.

## When To Activate

- User asks to create, write, draft, refactor, update, or improve a SKILL.md.
- User asks to critique, review, or analyze a skill design.
- User provides source material (paper, template, workflow, notes, prompt) and wants it turned into a skill.
- Keywords: "skill", "SKILL.md", "agent skill", "reusable skill", "skill template", "authoring".

## Operating Rules

- Read source material before designing. Extract reusable behavior, not every detail.
- Challenge weak assumptions that affect discovery, scope, reliability, or maintenance. State safe assumptions briefly and proceed.
- Prefer one compact SKILL.md. Add `references/`, `scripts/`, or `assets/` only when they reduce context load or improve reliability.
- Do not create extra docs (README, changelogs, install guides) unless explicitly asked.
- Use forward-slash paths. Keep all file references one level deep from SKILL.md.
- Keep SKILL.md body under 500 lines. Move variant-specific content to `references/` before it grows large.
- Do not force a skill for one-off answers, simple prompts, temporary plans, or generic knowledge. Recommend a document, prompt snippet, script, or checklist instead.

## Source Reading Workflow

Apply only the lenses the task requires. For simple requests, compress steps 1–3 into one paragraph.

1. **Intake**: identify document type, purpose, audience, scope, and key sections.
2. **Extract**: pull requirements, constraints, risks, decision points, workflows, and examples.
3. **Understand**: rewrite the core idea in practical terms for the target agent and user.
4. **Analyze**: inspect structure, trade-offs, failure cases, and maintenance burden.
5. **Assumptions**: classify as safe / risky / needs-confirmation.
6. **Validate**: check feasibility, contradictions, ambiguity, and missing acceptance criteria.
7. **Synthesize**: design workflow, decision rules, progressive disclosure, and validation loop.

## Degree of Freedom

- **High freedom** — concise text instructions when multiple approaches are valid.
- **Medium freedom** — workflow steps or templates when a preferred shape exists.
- **Low freedom** — exact scripts and commands when the work is fragile or high-stakes.

## Authoring Checklist

Copy and track as you work:

```
Skill Authoring Progress:
- [ ] Read source material; compress Intake → Understand if task is simple
- [ ] Decide: is a skill the right output? (if not, recommend alternative)
- [ ] Choose skill name in gerund form: verb-ing-noun (max 64 chars, lowercase + hyphens)
- [ ] Draft frontmatter: name + third-person discovery description with trigger terms
- [ ] Draft body: Purpose → When To Use → Workflow → Decision Rules → Validation
- [ ] Set degree of freedom (high / medium / low) and match instruction style
- [ ] Move large or variant-specific content to references/
- [ ] Apply validation loop below
- [ ] Revise until all items pass
```

## SKILL.md Template

```markdown
---
name: verb-ing-noun
description: Does the reusable task. Use when the user asks for [triggers].
---

# Human-Readable Skill Title

## Purpose
[One or two sentences.]

## When To Use
- [Specific trigger or context]

## Workflow
1. [First reliable action]
2. [Second reliable action]
3. [Validation or review action]

## Decision Rules
- [Rule that prevents bad scope]
- [Rule for choosing between alternatives]

## Validation
- [ ] Description is third-person and includes trigger terms
- [ ] Output is verifiable
```

## Validation Loop

Before delivering, confirm all pass. If any fails, revise and re-check:

- [ ] `name` is gerund form, lowercase, hyphens only, ≤64 chars, no reserved words.
- [ ] `description` is non-empty, third-person, specific, includes trigger terms, ≤1024 chars.
- [ ] Body is under 500 lines.
- [ ] `When To Use` section exists with at least two concrete triggers.
- [ ] Workflow has clear steps and at least one validation step.
- [ ] No time-sensitive information embedded.
- [ ] No duplicated content between SKILL.md and reference files.
- [ ] All file references one level deep, forward slashes.
- [ ] Degree of freedom matches instruction style.
