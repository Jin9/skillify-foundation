---
name: creating-skill-templates
description: Creates, refactors, and reviews SKILL.md files from papers, templates, workflows, domain notes, code, and existing prompts. Use when the user asks to create, update, critique, or synthesize an agent skill; convert source material into a reusable skill; apply skill-authoring best practices; or explore skill design alternatives.
---

> **⚠️ SUPERSEDED**: This skill has been superseded by [`skillify`](../../skillify/SKILL.md), which provides a formalized 8-mode architecture, 4-gate validation system, and complete progressive-disclosure infrastructure. The unique `platforms/` deployment kit from this skill has been merged into `skillify/platforms/`. This file is retained as archival reference only.

# Creating Skill Templates

## Purpose

Convert source material into a reusable agent skill that is easy to discover, concise enough to load, and specific enough to guide reliable execution.

## When To Use

- User asks to create, write, or draft a SKILL.md.
- User asks to refactor, update, or improve an existing skill.
- User asks to critique, review, or analyze a skill design.
- User provides source material (paper, template, workflow, notes, prompt) and wants it turned into a skill.
- Trigger terms: "skill", "SKILL.md", "agent skill", "reusable skill", "skill template", "authoring".

## Operating Rules

- Read source material before designing. Extract reusable behavior — not every detail.
- Challenge weak assumptions that affect discovery, scope, reliability, or maintenance. State safe assumptions briefly and proceed.
- Prefer one compact `SKILL.md`. Add `references/`, `scripts/`, or `assets/` only when they reduce context load or improve reliability.
- Do not create extra docs (`README.md`, changelogs, install guides) unless explicitly asked.
- Use forward-slash paths. Keep all file references one level deep from `SKILL.md`.
- Keep `SKILL.md` body under 500 lines. Move variant-specific or detailed content to `references/` before it grows large.
- Do not force a skill for one-off answers, simple prompts, temporary plans, or generic knowledge the model already has. Recommend a document, prompt snippet, script, or checklist instead.

## Source Reading Workflow

Apply only the lenses the task requires. For simple "create a skill" requests, compress steps 1–3 into one paragraph and jump to authoring.

1. **Intake**: identify document type, purpose, audience, scope, and key sections.
2. **Extract**: pull requirements, constraints, risks, decision points, workflows, and examples.
3. **Understand**: rewrite the core idea in practical terms for the target agent and user.
4. **Analyze**: inspect structure, trade-offs, failure cases, and maintenance burden.
5. **Assumptions**: classify as safe / risky / needs-confirmation. Convert dangerous assumptions into questions or skill constraints.
6. **Validate**: check feasibility, contradictions, ambiguity, and missing acceptance criteria.
7. **Synthesize**: design workflow, decision rules, progressive disclosure structure, and validation loop.

See [references/paper-reading.md](references/paper-reading.md) for the full action list per lens.

## Degree of Freedom

- **High freedom** — concise text instructions when multiple approaches are valid and judgment matters.
- **Medium freedom** — workflow steps, pseudocode, or templates when a preferred shape exists but adaptation is expected.
- **Low freedom** — exact scripts and commands when the work is fragile, repetitive, high-stakes, or mechanically verifiable.

## Authoring Workflow

Copy this checklist into your response and check items off as you complete them:

```
Skill Authoring Progress:
- [ ] Read source material; compress Intake → Understand if task is simple
- [ ] Decide: is a skill the right output? (if not, recommend alternative)
- [ ] Choose skill name in gerund form: `verb-ing-noun` (max 64 chars, lowercase + hyphens)
- [ ] Draft frontmatter: name + third-person discovery description with trigger terms
- [ ] Draft body: Purpose → When To Use → Workflow → Decision Rules → Validation
- [ ] Set degree of freedom (high / medium / low) and match instruction style
- [ ] Move large or variant-specific content to references/
- [ ] Apply review checklist (see references/authoring-checklist.md)
- [ ] Revise until all checklist items pass
```

## SKILL.md Template

Delete sections that do not earn their tokens:

```markdown
---
name: verb-ing-noun
description: Does the reusable task and names when to use it. Use when the user asks for [trigger terms], [file types], [domain tasks], or [related actions].
---

# Human-Readable Skill Title

## Purpose

[One or two sentences: what job does this skill perform?]

## When To Use

- [Specific trigger or context]
- [File types, tools, domain phrases, or user phrases if relevant]

## Workflow

1. [First reliable action]
2. [Second reliable action]
3. [Validation or review action — always include a feedback loop]

## Decision Rules

- [Rule that prevents bad scope or wrong execution]
- [Rule for choosing between alternatives]
- [Rule for handling ambiguity: state assumption, proceed]

## References

- **[Topic]**: See [references/topic.md](references/topic.md) when [condition].

## Validation

Before sending output, confirm:
- [ ] Description is third-person and includes trigger terms
- [ ] Output is verifiable (checklist, script result, or user review step)
- [ ] Examples present when output style matters
```

## Concrete Example

**Bad skill** (vague description, no trigger terms, no workflow, no validation):

```markdown
---
name: helper
description: Helps with documents.
---
# Helper
Use this skill to help with documents.
```

**Good skill** (specific, discoverable, executable, validates output):

```markdown
---
name: generating-commit-messages
description: Generates descriptive git commit messages from staged diffs or change descriptions. Use when the user asks for help writing commit messages, reviewing staged changes, or following conventional commit format.
---

# Generating Commit Messages

## Purpose
Turn a git diff or plain-language change description into a conventional commit message.

## Workflow
1. Read the diff or change description provided.
2. Identify the change type: feat / fix / chore / docs / refactor / test.
3. Write: `type(scope): brief summary` then one-line detail if needed.

## Validation
- [ ] Type matches the actual change (not always `chore`)
- [ ] Summary is under 72 characters
- [ ] Detail line added when scope is ambiguous
```

## Validation Loop

Before delivering the skill, run this check. If any item fails, revise and re-check:

- [ ] `name` is gerund form, lowercase, hyphens only, ≤ 64 chars, no reserved words.
- [ ] `description` is non-empty, third-person, specific, includes trigger terms, ≤ 1024 chars.
- [ ] Body is under 500 lines.
- [ ] `When To Use` section exists with at least two concrete triggers.
- [ ] Workflow has clear sequential steps and at least one validation or feedback step.
- [ ] No time-sensitive information embedded (or isolated in a legacy section).
- [ ] No content duplicated between `SKILL.md` and reference files.
- [ ] All file references are one level deep and use forward slashes.
- [ ] Scripts, if present, handle errors explicitly and produce useful diagnostics.
- [ ] Degree of freedom matches instruction style (text / pseudocode / exact commands).

See full authoring and testing checklist: [references/authoring-checklist.md](references/authoring-checklist.md).
