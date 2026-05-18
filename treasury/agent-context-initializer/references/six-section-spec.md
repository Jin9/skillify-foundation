# The six-section spec

GitHub's analysis of 2,500+ repositories identified six core content areas
that consistently put AGENTS.md files in the top performance tier:
**commands, testing, project structure, code style, git workflow, and
boundaries**. The candidate AGENTS.md contains exactly these six (plus an
optional squad-roles block — see `boundaries-and-squad-roles.md`). No
seventh section. Each section below states what it MUST and MUST NOT carry,
with terse good/bad examples.

## 1. Commands

- MUST: the exact CLI invocations for build, test, lint — with flags and
  options, not just tool names. Missing-flag commands are the most common
  reason an agent fails CI on its first run.
- MUST NOT: prose describing what the tools are, or install narratives a
  human README owns.

Good:

```
- Build:  `make build`
- Test:   `pytest -v --cov=src tests/`
- Lint:   `ruff check . && ruff format --check .`
```

Bad: `Run pytest to test the project and use ruff for linting.`

## 2. Testing

- MUST: the test framework, how to run a single test, how to run the full
  suite, and what "a passing test" means for this project (e.g. coverage
  gate threshold stated, not enforced here).
- MUST NOT: a tutorial on writing tests, or test-content rules an agent can
  read from existing tests.

Good:

```
- Single test: `pytest tests/test_pricing.py::test_tier_rounding -q`
- Full suite (must be green before PR): `pytest -q`
```

Bad: `Make sure you write good tests with high coverage.`

## 3. Project structure

- MUST: 3-5 bullets — key modules and their responsibilities. Enough for an
  agent to know where to look; not a substitute for code navigation.
- MUST NOT: an exhaustive directory tree, file-by-file descriptions, or
  architecture history the agent can infer by reading the codebase
  (encyclopedic project documentation is a documented anti-pattern).

Good:

```
- `src/api/`   — HTTP handlers; thin, no business logic
- `src/core/`  — domain rules; the only place decisions live
- `src/db/`    — persistence; migrations in `src/db/migrations/`
- `tests/`     — mirrors `src/`; fixtures in `tests/conftest.py`
```

Bad: a 40-line tree dump of every file.

## 4. Code style

- MUST: one real code snippet that shows the house style, plus a pointer to
  the lint/format config file. One snippet beats three paragraphs.
- MUST NOT: restate the lint rules in prose (that rule belongs in the lint
  config; see `deferral-meta-principle.md`).

Good:

```
Style is enforced by `pyproject.toml` ([tool.ruff]). Shape to match:

    def price(order: Order) -> Money:
        # early-return on guard, no nested else
        if order.is_empty:
            return Money.zero()
        return _tier(order).apply(order.subtotal)
```

Bad: `Use 4-space indents, snake_case, max line length 88, no unused
imports, sort imports alphabetically...` (all of this is the lint config).

## 5. Git workflow

- MUST: PR format, commit-message convention, branch naming, who reviews
  what — only the parts not enforced by branch protection or CI.
- MUST NOT: rules a branch-protection setting already enforces (e.g.
  "require 1 approval" — that is a setting, not prose; defer it).

Good:

```
- Branch: `feat/<ticket>-<slug>`; commits: Conventional Commits.
- PR body: what changed + how verified. Do not self-merge — request review
  from the module owner instead.
```

Bad: `Never merge without passing CI.` (CI required-check is a
branch-protection rule — defer it; see `deferral-meta-principle.md`).

## 6. Boundaries

- MUST: the three-tier block (always-allowed / requires-human-approval /
  explicitly-prohibited). Every prohibition paired with an allowed
  alternative. Full model in `boundaries-and-squad-roles.md`.
- MUST NOT: a flat list of 15+ sequential "don'ts" with no positive
  alternative — this causes agents to over-explore, stay conservative, and
  produce incomplete work.

Good: `Do not modify files in `src/db/migrations/` — propose the migration
in the PR description and wait for human approval instead.`

Bad: `Never touch migrations. Never touch config. Never touch secrets.
Never delete files. Never...` (over-specification failure mode).

## Section-presence rule

`scripts/agents_md_gate.py` fails the candidate if any of the six section
headings is absent (case-insensitive match on the words *commands*,
*testing*, *project structure*, *code style*, *git workflow*,
*boundaries*). Keep the headings literally named.
