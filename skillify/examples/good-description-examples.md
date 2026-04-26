# Good Description Examples

The `description` field in YAML frontmatter is the most important part of a skill. It determines when the agent loads the skill. These examples demonstrate the pattern: **[What] + [When/Triggers] + [Key capabilities]**.

---

## Category 1: Document & Asset Creation

### ✅ Frontend Design Skill
```yaml
description: >
  Create distinctive, production-grade frontend interfaces with high design
  quality. Use when building web components, pages, artifacts, posters,
  or applications. Triggers on "build a UI", "create a landing page",
  "design a dashboard", or "make a web app".
```

### ✅ DOCX Processing Skill
```yaml
description: >
  Comprehensive document creation, editing, and analysis with support for
  tracked changes, comments, formatting preservation, and text extraction.
  Use when working with professional documents (.docx files) for creating
  new documents, modifying content, working with tracked changes, or
  adding comments.
```

---

## Category 2: Workflow Automation

### ✅ Skill Creator (Meta-Skill)
```yaml
description: >
  Interactive guide for creating new skills. Walks the user through use case
  definition, frontmatter generation, instruction writing, and validation.
  Use when the user says "create a skill", "write a SKILL.md", or
  "help me build a skill".
```

### ✅ Sprint Planning Skill
```yaml
description: >
  Manages Linear project workflows including sprint planning, task creation,
  and status tracking. Use when user mentions "sprint", "Linear tasks",
  "project planning", or asks to "create tickets".
```

---

## Category 3: MCP Enhancement

### ✅ Sentry Code Review Skill
```yaml
description: >
  Automatically analyzes and fixes detected bugs in GitHub Pull Requests
  using Sentry error monitoring data via MCP. Use when reviewing PRs with
  Sentry issues, or when the user asks to "fix Sentry errors" or
  "review PR with Sentry data".
```

### ✅ Customer Onboarding Skill
```yaml
description: >
  End-to-end customer onboarding workflow for PayFlow. Handles account
  creation, payment setup, and subscription management. Use when user says
  "onboard new customer", "set up subscription", or "create PayFlow account".
```

---

## Category 4: Domain Expertise

### ✅ Financial Compliance Skill
```yaml
description: >
  Payment processing with regulatory compliance checks. Validates transactions
  against sanctions lists, jurisdiction rules, and risk assessments before
  processing. Use when handling payments, "process transaction", or
  "compliance check on payment".
```

### ✅ Database Migration Skill
```yaml
description: >
  Manages Alembic database migrations with safety checks. Generates migration
  files, validates schema changes, and runs upgrades. Use when the user asks
  to "create a migration", "update the schema", or "run alembic upgrade".
```

---

## Category 5: Code Quality

### ✅ Build-Test-Verify Skill
```yaml
description: >
  Run lint, test, and build verification commands for the project.
  Use when validating changes, running tests, or checking builds.
```

### ✅ Self-Review Checklist Skill
```yaml
description: >
  Quality gate to run before finishing work. Checks for convention drift,
  missing tests, and common mistakes. Use after completing code changes,
  before committing or creating a PR.
```

---

## With Negative Triggers

### ✅ Advanced Data Analysis (with negative trigger)
```yaml
description: >
  Advanced data analysis for CSV files. Use for statistical modeling,
  regression, clustering. Do NOT use for simple data exploration
  (use data-viz skill instead).
```

### ✅ Backend API Skill (with scope boundary)
```yaml
description: >
  Generates REST API endpoints following the project's FastAPI patterns.
  Use when creating new API routes, adding CRUD endpoints, or implementing
  pagination. Do NOT use for frontend components or database migrations
  (use frontend-design or alembic-migrations instead).
```

---

## Common Mistakes

### ❌ Too vague
```yaml
description: Helps with projects.
```

### ❌ Missing triggers
```yaml
description: Creates sophisticated multi-page documentation systems.
```

### ❌ Too technical, no user phrases
```yaml
description: Implements the Project entity model with hierarchical relationships.
```

### ❌ Contains XML (security risk)
```yaml
description: "<skill>Generates React components</skill>"
```

### ❌ No "when to use" information
```yaml
description: A comprehensive tool for managing cloud infrastructure deployments across multiple providers with support for rollbacks.
```
