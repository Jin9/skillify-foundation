---
name: your-skill-name
description: >
  What this skill does in one sentence. Use when the user asks to
  "trigger phrase 1", "trigger phrase 2", "trigger phrase 3", or works with [file types].
  Do NOT use for [negative trigger — what this skill should NOT handle].
# Optional fields (uncomment as needed):
# when_to_use: Additional trigger context
# compatibility: claude-code, codex, opencode
# metadata:
#   author: your-name
#   version: 1.0.0
---

# Your Skill Name

## Purpose

[A one-sentence description of the goal of this skill.]

## When to use this skill

- Use when: [trigger 1 — specific user intent or phrase]
- Use when: [trigger 2 — specific user intent or phrase]
- Use when: [trigger 3 — specific user intent or phrase]
- Do NOT use when: [negative trigger — what to avoid]

## Approach

Goal: [what a finished result looks like, in one sentence]
Constraints: [the one or two real constraints, each with its reason]
Done when: [the evidence that the job is complete]

### Ordered steps

[Keep this block only where a later step consumes an earlier step's output or a step is fragile or irreversible; describe everything else under Approach. Tag steps whose tier or effort differs from the contract default.]

1. [small/low] **[Step 1 Name]**: [Imperative action the agent must take]
   - Expected input: [what the agent receives]
   - Expected output: [what this step produces]
2. [mid/medium] **[Step 2 Name]**: [Imperative action the agent must take]
3. [frontier/high] **[Step 3 Name]**: [Judgment step: state the outcome and how to verify it, not the method]

## Output format

[Describe exactly what the agent should produce — specific files, formats, or deliverables]

```
output-directory/
├── [file 1]
├── [file 2]
└── [optional subdirectory]/
```

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

- [The one or two real constraints, each with its reason: "Do not X; Y would happen."]
- [Second constraint], because [reason].

## Examples

### Example 1: [Common scenario name]

**User says**: "[Example prompt the user would type]"

**Action**:
1. [Step 1 the agent takes]
2. [Step 2 the agent takes]

**Result**: [Expected output description]

### Example 2: [Edge case or variant]

**User says**: "[Another example prompt]"

**Action**:
1. [Step 1]
2. [Step 2]

**Result**: [Expected output]

## Verification

Before claiming the result, check and quote:
- [ ] [Evidence 1: the exit code, PASS line, or re-read file you will look at]
- [ ] [Evidence 2]
- [ ] [Evidence 3]

## References

- For [deep topic]: See `references/[filename].md`
- For [another topic]: See `references/[filename].md`
