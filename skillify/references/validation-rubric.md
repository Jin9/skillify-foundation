# Validation Rubric for Skills

Use this rubric to audit and review `SKILL.md` files. A production-ready skill scores 4+/5 in every category.

See also: `anti-patterns.md` for named failure modes and `security-checklist.md` for safety checks that inform dimension 8.

## 1. Trigger Quality (1-5)
- **1/5**: No description, or vague ("Helps with coding").
- **3/5**: Mentions the general domain but lacks specific trigger phrases.
- **5/5**: Highly specific. Contains exact user intents or phrases ("Use when the user asks to format CSV data").

## 2. Scope Focus (1-5)
- **1/5**: Attempts to do completely unrelated tasks (e.g., API creation + CSS styling).
- **3/5**: Related tasks, but too broad (e.g., "Fullstack React Developer").
- **5/5**: Singular, focused responsibility (e.g., "React Component Testing").

## 3. Workflow Clarity (1-5)
- **1/5**: Rambling paragraphs of theory.
- **3/5**: Bulleted list, but steps are vague or out of order.
- **5/5**: Imperative, numbered steps with clear entry and exit conditions.

## 4. Output Contract (1-5)
- **1/5**: Doesn't specify what the agent should produce.
- **3/5**: Vague output ("Produce good code").
- **5/5**: Exact deliverables ("Output a `jest` test file with 100% coverage").

## 5. Token Efficiency (1-5)
- **1/5**: Includes thousands of lines of documentation directly in `SKILL.md`.
- **3/5**: Moderately long, could be compressed.
- **5/5**: Lean `SKILL.md`. Uses `references/` for deep knowledge and `templates/` for code.

## 6. Conflict Risk (1-5)
- **1/5**: Hardcodes repo-specific rules (e.g., "Always use `my-company-cli`") that conflict with `AGENTS.md`.
- **3/5**: Minor overlapping instructions with other skills.
- **5/5**: Fully portable and isolated.

## 7. Reusability (1-5)
- **1/5**: Only works for one specific file in one specific repository.
- **3/5**: Works across the repository, but not easily portable to other projects.
- **5/5**: Can be dropped into any standard project using the target technology.

## 8. Security & Safety (1-5)
- **1/5**: Contains destructive commands, exfiltration patterns, or hardcoded secrets.
- **3/5**: No obvious risks, but uses broad tool permissions or unchecked external calls.
- **5/5**: No destructive commands, no external requests outside the workflow, no product bias, permissions are scoped.

## 9. Frontmatter Correctness (1-5)
- **1/5**: Missing frontmatter, wrong delimiters, or invalid YAML.
- **3/5**: Valid YAML but missing trigger phrases or using a non-kebab-case name.
- **5/5**: Valid YAML, kebab-case name matching folder, description under 1024 chars with triggers, no XML angle brackets.

## 10. Progressive Disclosure (1-5)
- **1/5**: All content in a single massive `SKILL.md` with no references.
- **3/5**: Some content split out, but deep docs remain inline.
- **5/5**: Lean `SKILL.md` body. Deep docs in `references/`. Clear pointers. Content lives in exactly one tier.

---

## Quick Scoring Template

```
Skill: [name]
Date: [YYYY-MM-DD]

| Dimension               | Score | Notes |
|--------------------------|-------|-------|
| 1. Trigger Quality       |  /5   |       |
| 2. Scope Focus           |  /5   |       |
| 3. Workflow Clarity       |  /5   |       |
| 4. Output Contract        |  /5   |       |
| 5. Token Efficiency       |  /5   |       |
| 6. Conflict Risk          |  /5   |       |
| 7. Reusability            |  /5   |       |
| 8. Security & Safety      |  /5   |       |
| 9. Frontmatter Correctness|  /5   |       |
| 10. Progressive Disclosure |  /5   |       |
| **TOTAL**                 |  /50  |       |

Pass threshold: every dimension at least 4/5 AND total at least 40/50.
```
