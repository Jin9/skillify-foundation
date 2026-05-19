---
authored_by: chinnawat.w (git author of record; GitHub Jin9)
phase_owner: Codex/GPT-5.5 (orchestration)
pipeline_phase: spec (0–7)
status: accepted
provenance_dated: 2026-05-19
---

# Cross-Model Skillify Upgrade Plan

Goal: upgrade `skillify` with a delegated research, analysis, implementation, review, and compaction workflow across Codex/GPT, Gemini, and Claude Code.

## Model Defaults

- Codex/GPT: `GPT-5.5` with `xhigh` reasoning.
  - Owns orchestration, baseline inventory, final merge authority, deterministic validation, and compaction.
- Gemini: `gemini-3.1-pro-preview`.
  - Owns broad source discovery, large-context reading, source tiering, and cross-platform synthesis.
- Claude Code: Opus 4.7 with max effort.
  - Owns deep skill architecture, contradiction detection, refactor recommendation, and adversarial review.

Command templates for execution:

```sh
gemini --model gemini-3.1-pro-preview --approval-mode plan --prompt "..."
claude -p --model opus --effort max --permission-mode plan "..."
```

If the local Claude CLI requires a concrete Opus 4.7 model string instead of the `opus` alias, use the configured Opus 4.7 model ID.

## Output Contract

Primary target:
- `~/.codex/skills/skillify/SKILL.md` (Codex skills install dir; resolve against `$HOME`/`$CODEX_HOME`, do not pin an OS account)

Research artifacts:
- `skill-design-methodology/research/sources.md`
- `skill-design-methodology/research/research-notes.md`
- `skill-design-methodology/research/analysis.md`
- `skill-design-methodology/research/recommendation.md`
- `skill-design-methodology/research/review.md`

Final artifacts:
- revised `skillify/SKILL.md`
- optional focused references under `skillify/references/` only when they reduce `SKILL.md` context load
- validation notes in the final assistant response, not as permanent human-facing docs inside the skill folder

## Cross-Model Delegation Rules

- Gemini and Claude own their assigned phase artifacts and may propose patches.
- Codex/GPT reviews all delegated artifacts before applying final repo edits.
- No model may silently overwrite another model's artifact; changes must be recorded in the next phase output.
- Unsupported model opinions lose to official sources, current repo constraints, validation scripts, and explicit user instructions.
- Keep permanent artifacts inside `skill-design-methodology/research/` or the target skill folder only.
- Do not create `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, `CONTRIBUTING.md`, or other human-facing docs inside the skill folder.

## Workflow

### Phase 0: Codex/GPT Baseline Inventory

Owner:
- Codex/GPT `GPT-5.5` with `xhigh` reasoning.

Inputs:
- current `skillify/SKILL.md`
- existing `skillify/references/`, `templates/`, `scripts/`, `platforms/`, `assets/`, and `examples/`
- local methodology files under `skill-design-methodology/`

Task:
- Read the current `skillify` implementation.
- Enumerate existing support files.
- Identify stable constraints: triggers, modes, output contracts, validation scripts, banned docs, and platform compatibility.
- Record baseline findings for use by Gemini and Claude.

Output:
- baseline section inside `research-notes.md`

Gate:
- Do not request model edits until the current skill boundary and validation assets are clear.

### Phase 1: Gemini Researcher

Owner:
- Gemini `gemini-3.1-pro-preview`.

Inputs:
- baseline notes
- `literature/skill-design-methodology/agent-skill-source-index.md`
- local methodology files

Task:
- Collect authoritative sources for agent skill design.
- Prioritize official OpenAI, Anthropic, GitHub Copilot, Gemini, and stable open-format references.
- Treat blogs, community repositories, and listicles as secondary evidence only.
- Record source trust tier and concrete influence on `skillify`.

Output:
- `research/sources.md`

Required fields per source:
- source name
- link or local file path
- trust tier
- relevant rule
- expected impact on `skillify`
- confidence

Gate:
- Separate official docs from community guidance.
- Every later recommendation must trace back to either current repo evidence or a listed source.

### Phase 2: Gemini Summarizer

Owner:
- Gemini `gemini-3.1-pro-preview`.

Inputs:
- `research/sources.md`
- existing methodology files
- baseline notes

Task:
- Extract reusable rules from sources and local methodology files.
- Remove duplicate, vague, or low-confidence advice.
- Group findings by trigger design, progressive disclosure, validation, safety, portability, and compaction.
- Keep the output operational enough to become skill instructions.

Output:
- `research/research-notes.md`

Gate:
- Avoid broad essays.
- Mark source-backed claims separately from Gemini synthesis.

### Phase 3: Claude Architecture Analyst

Owner:
- Claude Code Opus 4.7 with max effort.

Inputs:
- `research/sources.md`
- `research/research-notes.md`
- current `skillify` files

Task:
- Compare research findings against the current `skillify` implementation.
- Identify strengths to preserve, defects to fix, contradictions, scope risks, and overreach.
- Decide whether each improvement belongs in `SKILL.md`, `references/`, `templates/`, or `scripts/`.
- Challenge assumptions from Gemini's synthesis.

Output:
- `research/analysis.md`

Required sections:
- strengths to preserve
- defects to fix
- contradictions or weak assumptions
- scope boundaries
- compaction opportunities
- validation risks

Gate:
- Reject changes that make `skillify` generate target artifacts instead of skill files.
- Reject changes that duplicate guidance already covered by references.

### Phase 4: Claude Decision Maker

Owner:
- Claude Code Opus 4.7 with max effort.

Inputs:
- `research/analysis.md`
- current `skillify` files
- validation assets under `skillify/references/` and `skillify/scripts/`

Task:
- Produce a decision-complete refactor recommendation.
- Define the minimum patch that improves reliability without expanding scope.
- Choose which content should remain, move, split, shorten, or be left untouched.
- List rejected edits and why they should not be applied.

Output:
- `research/recommendation.md`

Required sections:
- recommended edits
- rejected edits and rationale
- expected behavior after refactor
- validation checklist
- proposed patch notes, if useful

Gate:
- Recommendation must be implementable in one focused pass.
- Do not recommend banned human-facing docs inside the skill folder.

### Phase 5: Codex/GPT Implementer

Owner:
- Codex/GPT `GPT-5.5` with `xhigh` reasoning.

Inputs:
- `research/recommendation.md`
- `research/analysis.md`
- current `skillify` files

Task:
- Apply the approved refactor to `skillify`.
- Keep `SKILL.md` concise, imperative, and self-sufficient after trigger.
- Move long mode-specific or optional material to one-level-deep reference files only when it reduces context load.
- Preserve existing scripts and templates unless the recommendation identifies a concrete defect.
- Resolve conflicts by favoring source-backed decisions and validation requirements.

Output:
- revised `skillify` files

Gate:
- Do not apply unsupported style-only churn.
- Preserve user edits and avoid unrelated refactors.

### Phase 6: Claude Reviewer

Owner:
- Claude Code Opus 4.7 with max effort.

Inputs:
- revised `skillify` files
- `research/recommendation.md`
- validation assets

Task:
- Review the changed skill against the validation rubric, anti-patterns, and security checklist.
- Challenge trigger quality, mode selection, output contracts, platform compatibility, and compaction.
- Identify blocking and non-blocking findings.

Output:
- `research/review.md`

Required sections:
- blocking findings
- non-blocking improvements
- validation concerns
- final accept or reject verdict

Gate:
- Blocking findings must be fixed before final compaction.

### Phase 7: Codex/GPT Compact And Validate

Owner:
- Codex/GPT `GPT-5.5` with `xhigh` reasoning.

Inputs:
- Claude review
- revised `skillify` files
- validation scripts

Task:
- Fix blocking review findings.
- Remove duplicated guidance and unnecessary wording.
- Keep frontmatter specific and within required limits.
- Keep `SKILL.md` focused on mode routing, core workflow, output contracts, and gates.
- Run deterministic validation where applicable.

Output:
- compact final `skillify`
- final assistant report with changed files, validation results, and remaining risks

Gate:
- Final `SKILL.md` must remain triggerable, compact, and self-sufficient for normal execution.
- Deterministic validation must pass or the final report must state the exact blocker.

## Execution Rules

- Prefer local methodology files first, then official sources, then secondary sources.
- Use Gemini for breadth and long-context synthesis.
- Use Claude Code Opus 4.7 for deep critique and architecture judgment.
- Use Codex/GPT for repo-safe integration and final validation.
- Preserve user edits and avoid unrelated refactors.
- Validate after editing with existing `skillify` scripts where applicable.
- Report changed files, validation results, and remaining risks at the end.
