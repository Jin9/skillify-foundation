# Skillify — Custom Instructions for GitHub Copilot

> **Source**: Adapted from `skillify/SKILL.md` for GitHub Copilot.
> This is a generated file. Edit the canonical `SKILL.md` and re-run `platforms/install.sh` to regenerate.

## When To Activate

When the user asks to create, write, draft, refactor, update, review, or critique a SKILL.md or agent skill, follow the workflow below.

## Skill Authoring Rules

- Read source material before designing. Extract reusable behavior, not every detail.
- Prefer one compact SKILL.md file. Add `references/` only when it reduces context load.
- Keep SKILL.md body under 500 lines.
- Use forward-slash paths. Keep file references one level deep.
- Do not force a skill for one-off answers or generic knowledge. Recommend a document or checklist instead.

## Authoring Workflow

1. Read source material. For simple requests, summarize purpose and scope in one paragraph.
2. Decide if a skill is the right output. If not, recommend an alternative (document, script, checklist).
3. Choose a skill name in **gerund form**: `verb-ing-noun` (max 64 chars, lowercase + hyphens only).
4. Draft YAML frontmatter with `name` and a third-person `description` that includes trigger terms.
5. Draft body: Purpose → When To Use → Workflow → Decision Rules → Validation.
6. Set degree of freedom: high (text), medium (pseudocode/templates), or low (exact scripts).
7. Move large content to `references/` directory.
8. Run the validation loop before delivering.

## SKILL.md Template

```markdown
---
name: verb-ing-noun
description: Does the task. Use when the user asks for [triggers].
---

# Skill Title

## Purpose
[What job does this skill perform?]

## When To Use
- [Trigger or context]

## Workflow
1. [Action]
2. [Validation step]

## Validation
- [ ] Description is third-person with trigger terms
- [ ] Output is verifiable
```

## Validation Loop

Before delivering, confirm:
- Name is gerund form, lowercase, hyphens, ≤64 chars
- Description is third-person, specific, includes triggers, ≤1024 chars
- Body under 500 lines
- At least two triggers in When To Use
- Workflow includes a validation step
- No time-sensitive info
- No duplicated content between SKILL.md and references
