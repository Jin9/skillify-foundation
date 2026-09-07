---
name: your-domain-skill-name
description: >
  Provides [domain] expertise for [task type]. Use when the user asks
  about "[trigger phrase 1]", "[trigger phrase 2]", "[trigger phrase 3]",
  or needs guidance on [domain-specific decisions].
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
- Use when: [domain-specific trigger 3]
- Do NOT use when: [out-of-scope scenario]

## Decision framework

### Step 1: Assess context · model-cost tier: [small | mid | frontier]

1. Determine [key variable 1]: [how to determine it]
2. Determine [key variable 2]: [how to determine it]
3. Classify the request:
   - **Type A**: [criteria] → follow Workflow A
   - **Type B**: [criteria] → follow Workflow B
   - **Type C**: [criteria] → follow Workflow C

### Step 2: Apply domain rules · model-cost tier: [small | mid | frontier]

[Keep this gate only when the checks are real business or compliance constraints; it is the one step that must not be skipped.]

Before taking action, verify:
- [ ] [Domain rule 1 — e.g., compliance check]
- [ ] [Domain rule 2 — e.g., data validation]
- [ ] [Domain rule 3 — e.g., authorization check]

IF all checks pass → proceed to Step 3
ELSE → [describe fallback: flag for review, escalate, document reason]

### Step 3: Execute · model-cost tier: [small | mid | frontier]

Follow the appropriate workflow based on Step 1 classification.

#### Workflow A: [Name]
1. [Action 1]
2. [Action 2]
3. [Action 3]

#### Workflow B: [Name]
1. [Action 1]
2. [Action 2]

### Step 4: Document & audit · model-cost tier: [small | mid | frontier]

1. Log the decision made and rationale
2. Record which domain rules were applied
3. Generate audit summary

## Domain knowledge

For detailed reference material, consult:
- [Topic 1]: `references/[topic-1].md`
- [Topic 2]: `references/[topic-2].md`
- [Glossary/Schemas]: `references/[schemas].md`

> Read only the reference that matches the Step 1 classification; the others cost context without changing the decision.

## Output format

[Describe what the agent produces — documents, reports, code with domain-specific requirements]

## Operating contract

[Filled from the skillify operating-contract template; keep the heading and the key prefixes exact, and delete this note.]

- Instruction priority: the user's request in this session takes precedence over this skill; repo-policy files (AGENTS.md, CLAUDE.md, or the host equivalent) take precedence over this skill's defaults. If following a line here would make you pause, ask for permission, leave requested work unfinished, or diverge from what the user asked, follow the user, say which line you set aside, and quote it.
- Autonomy: "can you", "help me", and "please" are instructions. Once the inputs above are present, act; do not ask for confirmation of work the user already authorized. Ask at most one question per run, and only when a required input is missing and cannot be taken from the request, pasted text, or the workspace; otherwise state the assumption in your opening line and proceed.
- Stop conditions: stop and ask only before [the skill's irreversible action, e.g. overwriting an existing file], or when finishing would change the scope the user set. When the user asks a question rather than for a change, the assessment is the deliverable. Before ending your turn, check your last paragraph: if it is a plan or a promise, do that work now.
- Verification: before claiming success, check [the specific evidence: exit code, PASS line, re-read of the written file] and quote it in the recap. Do not describe a check you did not run. Do not add tests for reversible, low-impact changes.
- Delegation: [none | which parts may run as parallel sub-agents and what each returns]. Batch independent tool calls; prefer asynchronous fan-out over spawn-and-wait; write messages to other agents so they stand alone.
- Progress: open with one line saying what you are about to do and which files you will touch; close with a recap that stands on its own (what changed, what was verified, what remains).
- Model-cost tier: [small | mid | frontier] by default; steps that differ are tagged inline as [tier/effort].

## Constraints

- Run the Step 2 domain-rules check before any action; it is the compliance gate and the one step that must not be skipped.
- Verify [domain-specific unsafe action] against [source] before doing it; [consequence if wrong].
- Record each decision and the rules applied; the audit trail is a deliverable.
- Use terminology from `references/[glossary].md` so outputs match the domain's vocabulary.

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
