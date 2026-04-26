---
name: your-skill-name
description: >
  What this skill does in one sentence. Use when the user asks to
  "trigger phrase 1", "trigger phrase 2", or works with [file types].
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
- Do NOT use when: [negative trigger — what to avoid]

## Core workflow

1. **[Step 1 Name]**: [Imperative action the agent must take]
   - Expected input: [what the agent receives]
   - Expected output: [what this step produces]
2. **[Step 2 Name]**: [Imperative action the agent must take]
3. **[Step 3 Name]**: [Imperative action the agent must take]

## Output format

[Describe exactly what the agent should produce — specific files, formats, or deliverables]

```
output-directory/
├── [file 1]
├── [file 2]
└── [optional subdirectory]/
```

## Constraints & anti-patterns

- DO NOT: [Action to avoid — be specific]
- DO NOT: [Another action to avoid]
- MUST ALWAYS: [Non-negotiable requirement]
- MUST ALWAYS: [Another non-negotiable requirement]

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

## Validation checklist

Before finalizing output, verify:
- [ ] [Quality check 1]
- [ ] [Quality check 2]
- [ ] [Quality check 3]

## References

- For [deep topic]: See `references/[filename].md`
- For [another topic]: See `references/[filename].md`
