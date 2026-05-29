# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repository is

This is **not an application** — it is a research-backed toolkit for engineering `SKILL.md` files (Agent Skills) for AI coding agents. There is no build, no server, no app runtime. The "product" is Markdown skill folders. Three layers, each grounded in the one before it:

- **`literature/`** — 54 curated source documents (~755K words) across the Anthropic, OpenAI, Copilot, open-standard, platform-rule, and academic agent ecosystems. The evidence base. Read-only research corpus; do not edit source files.
- **`skillify/`** — the skill-creator *meta-skill*. A single `SKILL.md` (8 modes) plus `references/`, `templates/`, `scripts/`, `examples/` synthesized from the literature. This is the engine that produces and audits all other skills.
- **`treasury/`** — 90 production skills built with `skillify`, grouped into 9 purpose groups. `treasury/README.md` is the **authoritative catalog** (the root `README.md` count can lag — trust `treasury/README.md`).

`skill-design-methodology/` exists locally but is **gitignored** (see `.gitignore`) — it is the design/research process that produced the toolkit, including `cross-model-skillification-pipeline.md`. Do not assume it is part of the published repo or rely on it being present for anyone else.

## Commands

All tooling is Python 3 standard-library only — no dependencies, no virtualenv. Run from the repo root.

```bash
# Scaffold a new skill folder skeleton
python3 skillify/scripts/init_skill.py <skill-name>

# Deterministic frontmatter + structure validation (validation gate 1)
python3 skillify/scripts/quick_validate.py treasury/<skill-folder>

# Verify every local references/ templates/ scripts/ link in a SKILL.md resolves
python3 skillify/scripts/check_links.py treasury/<skill-folder>
```

Validating a skill folder *is* the test suite for this repo — there is no global test runner. Some treasury skills ship their own `tests/` fixtures and `scripts/` self-checks (e.g. `treasury/eliciting-banking-brief/scripts/check_frame_rule_data_drift.py`); run those directly when touching that skill.

## The validation gate (hard contract)

Every skill change must pass, in order — this is enforced by `skillify/SKILL.md` and must be respected by any agent editing skills here:

1. `quick_validate.py` and `check_links.py` exit zero. **If gate 1 fails, do not write target files** — surface the error verbatim and stop.
2. `skillify/references/validation-rubric.md`: every dimension ≥ 4/5 **and** total ≥ 40/50.
3. `skillify/references/anti-patterns.md`: every applicable item absent or explicitly mitigated.
4. `skillify/references/security-checklist.md`: no unreviewed external HTTP, secret access, destructive command, broad permission, or vendor-bias risk.

`quick_validate.py` mechanically enforces: `name` is lowercase kebab-case, ≤64 chars, **must equal the folder name**, and must not contain `claude`/`anthropic`; `description` ≤1024 chars and must contain trigger language ("Use when" / "Use for" / "Use after" / "Triggers on" / "Activate when"); **no XML angle brackets in any frontmatter field**; `SKILL.md` ≤500 lines; `references/` files must be exactly one level deep; and no `README.md`/`CHANGELOG.md`/`INSTALLATION_GUIDE.md`/`QUICK_REFERENCE.md`/`CONTRIBUTING.md` anywhere inside a skill folder.

## Skill anatomy and progressive disclosure

A skill is a folder whose contents map onto a 3-tier context-loading model — content must live in **exactly one** tier:

- **Tier 1 (always in context):** `name` + `description` frontmatter, ~100 words.
- **Tier 2 (loaded when the skill triggers):** the `SKILL.md` body, < 5,000 words / ≤500 lines, imperative numbered steps with entry/exit conditions and an exact output contract.
- **Tier 3 (loaded on demand):** `references/` (deep guidance, one level deep), `templates/` (output skeletons), `scripts/` (deterministic automation for fragile/irreversible logic), `schemas/` (structured I/O), `examples/`, `assets/` (resources the agent should not read into context).

Do not duplicate guidance between `SKILL.md` and a reference; point to the reference instead. `skillify/platforms/` is skillify-internal cross-host install material and is **not** part of generated skill output.

## skillify's 8 modes

`skillify/SKILL.md` operates in: Create, Refactor, Review, Audit, Compress, Split, Merge, Adapt. All modes except Create are driven by `skillify/references/mode-playbooks.md`. Review and Audit are **read-only — they never edit files**. Refactor/Compress/Split/Merge preserve the original (diff or sibling dir) unless the user authorizes in-place overwrite.

## Constraints when working in this repo

These come from `skillify/SKILL.md` and govern any skill work here:

- When asked to build a skill that *would generate* artifact X, produce the **skill**, not artifact X.
- Do **not** edit repo-policy files (`AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`) as part of skill work unless the user explicitly asks for platform-policy adaptation.
- Do **not** invent trigger phrases, target users, or output contracts the user did not provide — ask once instead.
- Skills are authored to be **portable across hosts** (Claude Code, Codex, Copilot, Gemini/Antigravity). Avoid host-specific assumptions in skill bodies; use Adapt mode for platform-specific conventions.

## Cross-model authoring methodology

The toolkit itself was produced by a deliberate multi-model pipeline (documented in the gitignored `skill-design-methodology/cross-model-skillification-pipeline.md`): Gemini digests sources → Claude Opus synthesizes principles, designs lifecycle, and critiques for trigger/anti-pattern risk → GPT structures and polishes the final Markdown. When reasoning about *why* a skill is shaped a certain way, Claude Opus's role is the deep-synthesis and adversarial-critique step.

## Git workflow

Work happens on feature branches merged via PR (e.g. `Jin9/<topic>`). Git author of record is `chinnawat.w` (GitHub `Jin9`). Commit/push only when asked; branch off `main` first.
