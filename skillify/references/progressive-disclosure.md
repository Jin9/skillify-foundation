# Progressive Disclosure Architecture

The context window is a shared public good. Skills compete with the system prompt, conversation history, other skills' metadata, and the user's request. Progressive disclosure ensures the agent loads only what the current task requires.

See also: `workflow-patterns.md` for organizing Tier-2 instructions and `anti-patterns.md` for context-bloat failure modes.

## The 3-Tier Model

### Tier 1: Root Rules File (~100 lines, always loaded)

The highest-leverage file in the repo. Every token competes for attention.

- **Files**: `CLAUDE.md`, `AGENTS.md`, `.github/copilot-instructions.md`
- **Content**: Universal rules that apply to every task
- **Structure**: Why / What (project map) / How (always-apply rules) / Progressive Disclosure pointers
- **Rule**: If it doesn't apply to literally every task, push it down to Tier 2 or 3

### Tier 2: Skills (< 500 lines, loaded on demand)

Task-specific behavior loaded only when the agent determines the skill is relevant.

- **Files**: `SKILL.md` inside skill directories
- **Loading**: Agent reads the `description` field from frontmatter → decides whether to load the body
- **Content**: Step-by-step instructions, output formats, constraints, examples
- **Rule**: Keep under 500 lines. Point to Tier 3 for deep content.

### Tier 3: Deep References (unlimited, loaded as needed)

Reference material the agent reads only when it needs the full picture.

- **Files**: `references/`, `assets/`, `docs/agent-guides/`
- **Loading**: Agent follows pointers from Tier 2 SKILL.md
- **Content**: API docs, schemas, detailed guides, domain knowledge, templates
- **Rule**: Structure with table of contents for files over 100 lines

## The 3-Level Loading System

1. **Metadata (name + description)** — Always in context (~100 words). Used for trigger decisions.
2. **SKILL.md body** — Loaded when skill triggers (< 5,000 words). Contains instructions.
3. **Bundled resources** — Loaded as needed by the agent (unlimited). Scripts can execute without being read into context.

## When to Split Content

Move content from SKILL.md to references when:

- SKILL.md exceeds **500 lines** or **5,000 tokens**
- A section is only needed for a specific sub-task
- The same reference is useful across multiple skills
- The content is detailed documentation (API specs, schemas, policies)
- The content is a large example or walkthrough

## Splitting Patterns

### Pattern 1: High-level guide with references

Keep the quick-start in SKILL.md. Link to deep content.

```markdown
# PDF Processing

## Quick start
Extract text with pdfplumber: [code example]

## Advanced features
- **Form filling**: See `references/forms.md` for complete guide
- **API reference**: See `references/api-reference.md` for all methods
- **Examples**: See `references/examples.md` for common patterns
```

The agent loads `forms.md`, `api-reference.md`, or `examples.md` only when needed.

### Pattern 2: Domain-specific organization

For skills with multiple domains, organize by domain to avoid loading irrelevant context:

```
bigquery-skill/
├── SKILL.md (overview + navigation)
└── references/
    ├── finance.md    (revenue, billing metrics)
    ├── sales.md      (opportunities, pipeline)
    ├── product.md    (API usage, features)
    └── marketing.md  (campaigns, attribution)
```

When a user asks about sales metrics, the agent only reads `sales.md`.

### Pattern 3: Conditional details

Show basic content inline. Link to advanced content for edge cases.

```markdown
# DOCX Processing

## Creating documents
Use docx-js for new documents. See `references/docx-js.md`.

## Editing documents
For simple edits, modify the XML directly.

**For tracked changes**: See `references/redlining.md`
**For OOXML details**: See `references/ooxml.md`
```

## Anti-patterns

- **Do NOT deeply nest references.** Keep all reference files one level deep from SKILL.md.
- **Do NOT duplicate content across tiers.** A skill should point to a reference, not copy its content.
- **Do NOT stuff everything into Tier 1.** Every line in the root rules file competes for attention.
- **Do NOT create skills for trivially simple tasks.** A single CI workflow doesn't need its own skill.

## Key Principle

> The default assumption is that the agent is already very smart. Only add context the agent doesn't already have. Challenge each piece of information: "Does the agent really need this explanation?" and "Does this paragraph justify its token cost?"
