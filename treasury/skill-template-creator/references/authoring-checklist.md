# Skill Authoring Checklist

Use this as the complete pre-publish quality gate. All items in **Core quality** must pass. Apply other sections only when relevant.

## Core Quality

- [ ] Description is specific, third-person, and includes trigger terms (file types, user phrases, domain words).
- [ ] Description includes both **what** the skill does and **when** to use it.
- [ ] `name` is gerund or action-noun form, lowercase, hyphens only, ≤ 64 chars, no reserved words (`anthropic`, `claude`).
- [ ] `SKILL.md` body is under 500 lines.
- [ ] `When To Use` section has at least two concrete triggers.
- [ ] Workflow has clear sequential steps and at least one explicit validation or feedback loop.
- [ ] Examples are concrete when output style matters (not just abstract descriptions).
- [ ] Alternatives narrowed to a sensible default with escape hatch noted.
- [ ] No time-sensitive information embedded (use `## Old Patterns` section if legacy context needed).
- [ ] Consistent terminology throughout — choose one term per concept and stick to it.
- [ ] No duplicated content between `SKILL.md` and reference files.
- [ ] All file references are one level deep from `SKILL.md` and use forward slashes.
- [ ] Progressive disclosure applied: essential content in `SKILL.md`, detail in `references/`.

## Code and Scripts (apply when skill includes scripts)

- [ ] Scripts solve the problem directly — no punting to the model for error handling.
- [ ] All error conditions handled explicitly with useful diagnostic messages.
- [ ] No magic numbers — all constants documented with rationale.
- [ ] Required packages listed and verified available in the execution environment.
- [ ] Execution intent is explicit: "Run `x.py`" vs "Read `x.py` for the algorithm."
- [ ] No Windows-style paths — all paths use forward slashes.
- [ ] Validation / verification step included for any destructive or batch operation.

## Testing (apply before sharing widely)

- [ ] At least three evaluation scenarios created covering the core use cases.
- [ ] Tested with Haiku (does it provide enough guidance?) and Sonnet/Opus (does it avoid over-explaining?).
- [ ] Tested with real usage scenarios, not only synthetic test cases.
- [ ] Team feedback incorporated if the skill will be used across multiple people.

## Naming Convention Reference

| Pattern | Example | Notes |
|---|---|---|
| Gerund (preferred) | `creating-skill-templates` | Clearly describes the activity |
| Noun phrase (acceptable) | `skill-template-creation` | Slightly less discoverable |
| Action-oriented (acceptable) | `create-skill-templates` | Imperative form |
| **Avoid** | `helper`, `utils`, `tools` | Vague, not discoverable |
| **Avoid** | `claude-tools`, `anthropic-helper` | Reserved words |
