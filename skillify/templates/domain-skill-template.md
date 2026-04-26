---
name: your-domain-skill-name
description: >
  Provides [domain] expertise for [task type]. Use when the user asks
  about "[trigger phrase 1]", "[trigger phrase 2]", or needs guidance
  on [domain-specific decisions].
# compatibility: claude-code, codex, opencode
# metadata:
#   author: your-name
#   version: 1.0.0
---

# Your Domain Skill Name

## Purpose

[One sentence: what domain expertise this skill provides and what decisions it enables.]

## When to use this skill

- Use when: [domain-specific trigger 1]
- Use when: [domain-specific trigger 2]
- Do NOT use when: [out-of-scope scenario]

## Decision framework

### Step 1: Assess context

1. Determine [key variable 1]: [how to determine it]
2. Determine [key variable 2]: [how to determine it]
3. Classify the request:
   - **Type A**: [criteria] → follow Workflow A
   - **Type B**: [criteria] → follow Workflow B
   - **Type C**: [criteria] → follow Workflow C

### Step 2: Apply domain rules

Before taking action, verify:
- [ ] [Domain rule 1 — e.g., compliance check]
- [ ] [Domain rule 2 — e.g., data validation]
- [ ] [Domain rule 3 — e.g., authorization check]

IF all checks pass → proceed to Step 3
ELSE → [describe fallback: flag for review, escalate, document reason]

### Step 3: Execute

Follow the appropriate workflow based on Step 1 classification.

#### Workflow A: [Name]
1. [Action 1]
2. [Action 2]
3. [Action 3]

#### Workflow B: [Name]
1. [Action 1]
2. [Action 2]

### Step 4: Document & audit

1. Log the decision made and rationale
2. Record which domain rules were applied
3. Generate audit summary

## Domain knowledge

For detailed reference material, consult:
- [Topic 1]: `references/[topic-1].md`
- [Topic 2]: `references/[topic-2].md`
- [Glossary/Schemas]: `references/[schemas].md`

> **Note**: Do NOT load all references upfront. Read only the reference relevant to the current classification from Step 1.

## Output format

[Describe what the agent produces — documents, reports, code with domain-specific requirements]

## Constraints

- DO NOT skip the domain rules check in Step 2
- DO NOT make [domain-specific unsafe action] without verification
- MUST document all decisions for audit trail
- MUST use terminology from `references/[glossary].md`

## Examples

### Example: [Common domain scenario]

**User says**: "[Example prompt]"

**Classification**: Type A

**Domain rules check**:
- [Rule 1]: ✅ passed
- [Rule 2]: ✅ passed

**Action**:
1. [Step 1 with domain context]
2. [Step 2 with domain context]

**Result**: [Domain-specific output]

## Directory structure

```
your-domain-skill-name/
├── SKILL.md
└── references/
    ├── [topic-1].md       # [Description of what this covers]
    ├── [topic-2].md       # [Description of what this covers]
    └── [schemas].md       # [Glossary, schemas, or data models]
```
