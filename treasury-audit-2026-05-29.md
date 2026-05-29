# Treasury Skill Audit — all 92 skills under `treasury/`

- Date: `2026-05-29` · Mode: `Audit` (skillify, read-only) · **Refreshed after the 2026-05-29 validator fix + 4-skill remediation**
- Gate 1 below is recomputed **live** against the corrected `quick_validate.py` + `check_links.py`. Rubric scores are from the parallel per-skill audit; verdicts for skills changed this session are overridden to reflect the fixes.

## 1. Executive summary

- **Skills:** 92/92
- **Gate 1 (live, post-fix):** 91 pass · 1 fail — remaining: `planning-banking-tests` (banned run/ artifacts; out of remediation scope).
- **Rubric:** 90/92 pass (every dim ≥4/5 and total ≥40/50). Score range 45–50/50, mean 49.0.
- **Recommendations:** Ship: 91 · Refactor before ship: 1
- **Fixed this session (8):** `cross-examine`, `drafting-ba-stories`, `frame-debate`, `opening-debate-panel`, `report-debate`, `reporting-research-run`, `revise-positions`, `synthesize-consensus` — all now pass gate 1.
- **Security:** zero flags across all 92 skills.
- **Exceptions needing attention (2):** `planning-banking-tests`, `refactoring-go-services`.

## 2. Cross-cutting findings

**Refreshed 2026-05-29**, after a remediation session that fixed the skillify validator and refactored every skill the corrected gate surfaced as fixable. Gate-1 results below are recomputed live against the **corrected** validators. The treasury is healthy (audit mean **48.9/50**, **zero security flags**) and now sits at **91/92 passing gate 1** — the lone remaining failure (`planning-banking-tests`) is a real, previously-masked issue requiring folder cleanup, not a one-liner. Findings are grouped: **(A)** validator — fixed; **(B)** skills — fixed; **(C)** remaining open; **(D)** portfolio.

### A. Skillify validator — FIXED this session

The original audit's gate 1 ran against a buggy validator that both over- and under-reported. All four defects are now fixed in `skillify/scripts/` (stdlib-only preserved):

1. **False-positives on valid nested YAML** — `quick_validate.py` rejected nested mappings / block sequences (the `metadata:` and `inputs:`/`outputs:` blocks). Fixed: the parser now consumes nested blocks. 16 spuriously-failing skills now pass.
2. **Early-return masked structural checks** — `validate()` returned on the frontmatter "error" before the banned-doc / size / refs-depth checks. Fixed: those checks now always run. This is what newly (correctly) surfaces the `planning-banking-tests` banned READMEs in (C), and the trigger-language misses fixed in (B).
3. **Inline-comment angle-bracket false-positive** — the parser captured `# comment` text as part of a field value, so `allowed-tools: Read, Write  # reads <dir>` tripped the no-angle-brackets rule. Fixed: inline comments are stripped before the check.
4. **`check_links.py` prose false-positive** — the resource regex matched bare prose tokens. Fixed: it now requires backtick or markdown-link delimiters, so `synthesize-report`'s prose `examples/framing` is no longer flagged.

**Caveat the fix keeps:** the validator is a structural lint, not a full YAML parse — deep YAML validity inside nested blocks is intentionally not checked (no YAML dependency); genuine YAML errors are caught at skill-load time, and the three that existed were fixed in (B). The trigger-language check is strict (requires a literal "Use when/for/after/Triggers on/Activate when" marker).

### B. Skills — FIXED this session (8)

All eight now pass gate 1 (verified) and are recommended **Ship**:

- **`cross-examine`, `synthesize-consensus`** — quoted the unquoted `??` in the `inputs` flow-mapping `source:` value (genuinely invalid YAML) and added a literal trigger marker.
- **`revise-positions`** — quoted the unquoted `[panelist]` in `source:` (invalid YAML) and added a trigger marker.
- **`drafting-ba-stories`** — created the missing `references/example_TierRateStory.md` exemplar (a full `user-story-template` instance: numbered Business-Logic rules, an "Income → interest-rate tiers" decision table, a `> ⚠️ Note`); the broken link now resolves.
- **`frame-debate`, `opening-debate-panel`, `report-debate`, `reporting-research-run`** — added a literal trigger marker to each description (`", or when asked to"` → `". Use when asked to"`; `"Use as the panel-open stage"` → `"Use when the workflow runs the panel-open stage"`; `"Trigger phrases:"` → `"Triggers on:"`). These were the gate-1 trigger-language misses the early-return had masked; rubric quality was already high.

### C. Remaining open item

- **Banned docs / folder hygiene: `planning-banking-tests`.** 3 `README.md` files under `runs/`, plus 2.7M of run/audit detritus (`__pycache__`, ad-hoc `merge.py`/`diff.py`, baseline fixtures) and foreign-user absolute paths `/Users/IF640063` baked into several `references/*.md`. The `SKILL.md` + `references/` + `schemas/` core is excellent; the surrounding artifacts are not shippable. Needs a dedicated cleanup pass (move run/audit artifacts out of the skill; strip foreign paths) — deliberately left out of this remediation.

### D. Portfolio observations — OPEN

- **Near-duplicate pair: `model-selection` vs `model-selection-updated`** — identical section structure, both Ship individually; reconcile (promote one, retire/merge the other).
- **Borderline output contract: `refactoring-go-services`** — rubric-fails on Output Contract (3/5: deliverable is implied in-place edits) yet rated Ship. Tighten the deliverable/diff contract.
- **Documentation drift (flag only — out of scope):** `CLAUDE.md` says "34 skills / 7 categories"; actual is **92 / 9**.

### E. What's strong

Triggers (content), scope focus, workflow clarity, and security are uniformly high. No external HTTP, secret access, destructive commands, broad permissions, or vendor bias in any skill. Sibling-skill delegation is used well and sharpens boundaries.

## 3. Scorecard matrix

`Gate-1` is live against the fixed validators. `Rubric` passes only if every dimension ≥4/5 and total ≥40/50. ✔ = fixed this session.

### Banking, BA Delivery & Requirements (14)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `analyzing-banking-requirements` | PASS | 47/50 | ✅ | Ship | clean |
| `assembling-tl-handoff` | PASS | 49/50 | ✅ | Ship | clean |
| `assessing-ba-feasibility` | PASS | 49/50 | ✅ | Ship | clean |
| `checking-ba-governance` | PASS | 48/50 | ✅ | Ship | clean |
| `drafting-ba-stories` | PASS | 47/50 | ✅ | ✔ Ship | FIXED 2026-05-29: created references/example_TierRateStory.md; broken link resolved, passes gate 1. |
| `eliciting-banking-brief` | PASS | 50/50 | ✅ | Ship | gate-1 indented-frontmatter is validator limitation (valid YAML per frontmatter-guide); committed scripts/__pycache__ .pyc cruft; 2 delegation boundaries (design-review, implement-from-spec) name skills absent from catalog but are advisory negative-triggers, not links. |
| `evaluate-banking-compliance` | PASS | 50/50 | ✅ | Ship | clean |
| `extract-brief-structure` | PASS | 50/50 | ✅ | Ship | clean |
| `orchestrate-banking-brief-pipeline` | PASS | 49/50 | ✅ | Ship | clean (note: runner uses subprocess shell=True on operator-supplied stage commands by design, must be user-run) |
| `running-ba-pipeline` | PASS | 48/50 | ✅ | Ship | clean |
| `running-business-analysis-workflow` | PASS | 50/50 | ✅ | Ship | clean |
| `scoping-ba-intake` | PASS | 50/50 | ✅ | Ship | clean |
| `scoping-technical-requirements` | PASS | 48/50 | ✅ | Ship | clean |
| `sweep-ambiguities` | PASS | 50/50 | ✅ | Ship | clean |

### Architecture, Engineering Decisions & Planning (20)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `api-contract-design` | PASS | 49/50 | ✅ | Ship | clean |
| `architecting-fintech-systems` | PASS | 49/50 | ✅ | Ship | clean |
| `architecture-decision` | PASS | 48/50 | ✅ | Ship | clean |
| `data-modeling` | PASS | 49/50 | ✅ | Ship | clean |
| `defining-engineering-standards` | PASS | 49/50 | ✅ | Ship | clean |
| `delivery-planning` | PASS | 47/50 | ✅ | Ship | Description negative trigger 'Do NOT use for sizing complexity and unknown-risk before a date' is grammatically truncated (dangling fragment) — copy fix only. |
| `designing-tech-lead-handoff` | PASS | 49/50 | ✅ | Ship | clean |
| `domain-modeling` | PASS | 48/50 | ✅ | Ship | clean |
| `engineer-growth-planning` | PASS | 49/50 | ✅ | Ship | clean |
| `engineering-doc-planning` | PASS | 49/50 | ✅ | Ship | clean |
| `generate-ux-pack` | PASS | 47/50 | ✅ | Ship | gate1 FAIL-indented = validator/guide conflict (valid YAML); broken xref 'ba-elicit-from-raw' not in catalog; RATIONALE.md aux doc + ref<->ref duplication (both self-acknowledged, minor) |
| `integration-design` | PASS | 49/50 | ✅ | Ship | clean |
| `model-selection` | PASS | 47/50 | ✅ | Ship | clean; note model-selection-updated sibling exists in catalog — flag for cross-cutting near-dup review |
| `model-selection-updated` | PASS | 48/50 | ✅ | Ship | Near-duplicate suspicion: catalog also has model-selection; '-updated' suffix implies it supersedes the twin (cross-cutting reconcile, likely Retire/Merge older). |
| `pr-design-review` | PASS | 49/50 | ✅ | Ship | clean |
| `production-readiness` | PASS | 48/50 | ✅ | Ship | clean |
| `refactor-decision` | PASS | 49/50 | ✅ | Ship | clean |
| `risk-estimation` | PASS | 49/50 | ✅ | Ship | clean |
| `technical-debt-management` | PASS | 49/50 | ✅ | Ship | clean |
| `technical-feasibility` | PASS | 49/50 | ✅ | Ship | clean |

### Implementation, Platform Templates & Code Review (12)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `crafting-backend-code` | PASS | 50/50 | ✅ | Ship | clean |
| `crafting-frontend-code` | PASS | 50/50 | ✅ | Ship | clean |
| `implement-backend-feature` | PASS | 49/50 | ✅ | Ship | gate1 FAIL-indented (validator vs frontmatter-guide conflict, not a real defect); two forward cross-refs design-backend-feature/generate-backend-fix not yet in catalog |
| `implement-frontend-feature` | PASS | 48/50 | ✅ | Ship | gate-1 FAIL-indented (nested metadata: mapping, valid YAML); 3 deferral-target sibling skills absent from catalog (by-design future, negative-trigger only) |
| `implementing-go-template-requirements` | PASS | 50/50 | ✅ | Ship | clean |
| `platform-common` | PASS | 49/50 | ✅ | Ship | clean |
| `platform-go-service` | PASS | 49/50 | ✅ | Ship | clean |
| `platform-go-service-common` | PASS | 50/50 | ✅ | Ship | clean |
| `refactoring-go-services` | PASS | 46/50 | ❌ | Ship | Output contract vague: no per-iteration report/diff format or naming; deliverable is implied in-place edits. |
| `review-backend-code` | PASS | 49/50 | ✅ | Ship | clean |
| `review-frontend-code` | PASS | 49/50 | ✅ | Ship | gate1 FAIL-indented (validator<->guide conflict, frontmatter is valid YAML); phantom delegation to analyze-frontend-performance (not in catalog) x4 |
| `review-rust-code` | PASS | 50/50 | ✅ | Ship | clean |

### Testing, QA & Validation (4)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `generating-gherkin-acceptance-criteria` | PASS | 50/50 | ✅ | Ship | clean |
| `planning-banking-tests` | FAIL | 45/50 | ❌ | Refactor before ship | 3 banned README.md under runs/ (latent gate-1, masked by frontmatter parser early-return); runs/ tree 2.7M incl __pycache__/.pyc + ad-hoc merge.py/diff.py; foreign-user paths /Users/IF640063 + stale name qa-plan-from-brief in refs/run/audit; gate-1 indented-frontmatter (validator<->guide conflict, frontmatter is valid YAML) |
| `testing-strategy` | PASS | 48/50 | ✅ | Ship | clean; minor: 'contract-testing' reads like a skill name but is not in catalog (inline note, not a file link, link-check passes) |
| `validating-banking-implementation` | PASS | 50/50 | ✅ | Ship | clean |

### Agent Orchestration & Workflow Infrastructure (10)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `agent-context-initializer` | PASS | 50/50 | ✅ | Ship | clean |
| `agentic-workflow-design` | PASS | 48/50 | ✅ | Ship | clean |
| `authoring-scaffold-profile` | PASS | 49/50 | ✅ | Ship | clean |
| `composing-agent-pipelines` | PASS | 50/50 | ✅ | Ship | clean |
| `developing-langgraph-workflows` | PASS | 49/50 | ✅ | Ship | clean |
| `drafting-stage-prompt` | PASS | 49/50 | ✅ | Ship | clean |
| `multi-agent-handoff-architect` | PASS | 50/50 | ✅ | Ship | clean |
| `orchestrating-agent-scaffold` | PASS | 50/50 | ✅ | Ship | clean |
| `orchestrating-openclaw-squad` | PASS | 50/50 | ✅ | Ship | clean |
| `reviewing-implement-gate` | PASS | 50/50 | ✅ | Ship | clean |

### Research, Debate & Knowledge Synthesis (14)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `augment-diagrams` | PASS | 50/50 | ✅ | Ship | gate-1 FAIL-indented = nested block-sequence inputs/outputs the validator cannot parse but frontmatter-guide permits; valid YAML, loads in real Claude Code |
| `cross-examine` | PASS | 49/50 | ✅ | ✔ Ship | FIXED 2026-05-29: quoted '??' in inputs flow-mapping + added 'Use when' trigger; valid YAML, passes gate 1. |
| `extract-findings` | PASS | 49/50 | ✅ | Ship | gate-1 FAIL is validator<->guide conflict (valid block-sequence YAML); not a real defect |
| `frame-debate` | PASS | 48/50 | ✅ | ✔ Ship | FIXED 2026-05-29: ', or when asked to' -> '. Use when asked to'; passes gate 1. |
| `opening-debate-panel` | PASS | 50/50 | ✅ | ✔ Ship | FIXED 2026-05-29: 'Use as the panel-open stage' -> 'Use when the workflow runs the panel-open stage'; passes gate 1. |
| `plan-research` | PASS | 48/50 | ✅ | Ship | gate-1 FAIL is the known validator-vs-guide nested-mapping conflict (inputs/outputs flow-seq); not a real defect. |
| `report-debate` | PASS | 49/50 | ✅ | ✔ Ship | FIXED 2026-05-29: 'Trigger phrases:' -> 'Triggers on:'; passes gate 1. |
| `reporting-research-run` | PASS | 48/50 | ✅ | ✔ Ship | FIXED 2026-05-29: 'Trigger phrases:' -> 'Triggers on:'; passes gate 1. |
| `research-vault-librarian` | PASS | 50/50 | ✅ | Ship | clean |
| `review-report` | PASS | 50/50 | ✅ | Ship | gate-1 FAIL is validator<->guide conflict (block-seq inputs/outputs), not a defect; validation-gate.md says 'four' but lists five sections (cosmetic) |
| `revise-positions` | PASS | 48/50 | ✅ | ✔ Ship | FIXED 2026-05-29: quoted '[panelist]' in inputs source + added trigger; frontmatter now valid YAML, passes gate 1. |
| `search-sources` | PASS | 48/50 | ✅ | Ship | gate-1 quick_validate fails ONLY on nested inputs/outputs block-sequences (valid YAML, guide-permitted); validator<->guide conflict, not a skill defect. |
| `synthesize-consensus` | PASS | 49/50 | ✅ | ✔ Ship | FIXED 2026-05-29: quoted '??' in inputs source + added trigger; frontmatter now valid YAML, passes gate 1. |
| `synthesize-report` | PASS | 50/50 | ✅ | Ship | RESOLVED 2026-05-29: check_links prose false-positive fixed at validator level; clean 50/50. |

### Security, Governance & Compliance (5)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `configuring-sandbox-allowlist` | PASS | 48/50 | ✅ | Ship | Refusal table appears in both SKILL.md (gate) and pattern-rules.md (rationale) with divergent columns — mild, intentional. |
| `devops-infrastructure-hardener` | PASS | 50/50 | ✅ | Ship | clean |
| `governance-policy-generator` | PASS | 50/50 | ✅ | Ship | clean |
| `reviewing-software-security` | PASS | 50/50 | ✅ | Ship | clean |
| `universal-spec-validator` | PASS | 50/50 | ✅ | Ship | clean |

### Observability, Cost & Incident Operations (6)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `authoring-workflow-postmortem` | PASS | 48/50 | ✅ | Ship | clean |
| `incident-response` | PASS | 49/50 | ✅ | Ship | clean |
| `observability-design` | PASS | 48/50 | ✅ | Ship | clean |
| `observability-telemetry-instrumenter` | PASS | 50/50 | ✅ | Ship | clean |
| `performance-cost-review` | PASS | 48/50 | ✅ | Ship | clean |
| `reviewing-agent-spend` | PASS | 50/50 | ✅ | Ship | clean |

### Code Analysis, Productivity & Publishing (7)

| Skill | Gate-1 | Score | Rubric | Recommendation | Notes |
|-------|--------|-------|--------|----------------|-------|
| `business-logic-extractor` | PASS | 50/50 | ✅ | Ship | clean |
| `generating-pseudocode` | PASS | 49/50 | ✅ | Ship | clean |
| `organizing-local-files` | PASS | 50/50 | ✅ | Ship | clean |
| `progressive-bug-hunter` | PASS | 50/50 | ✅ | Ship | clean |
| `publishing-git-review-requests` | PASS | 50/50 | ✅ | Ship | clean |
| `rendering-readable-html` | PASS | 50/50 | ✅ | Ship | clean |
| `thai-translator` | PASS | 48/50 | ✅ | Ship | clean |

## 4. Detailed reports (exceptions)

Full audit detail for every skill that currently fails a gate or scores below 40/50 (2 skills). Skills fixed this session are no longer exceptions and are omitted here.

### refactoring-go-services

#### Summary

- **Skill:** refactoring-go-services
- **Folder:** /Users/admin/Desktop/project/research/skillify-foundation/treasury/refactoring-go-services
- **Date:** 2026-05-29
- **Mode:** Audit (read-only)
- **Result:** Pass with fixes (rubric total 46/50, but Output Contract scored 3 so rubricPass = false)

A single-file SKILL.md (148 lines) defining an imperative, behavior-preserving incremental refactoring loop for Go DDD/CQRS microservices. Cleanly scoped and well-fenced from sibling skills. The only material gap is an under-specified output/deliverable contract.

#### Deterministic Checks

- **GATE 1: PASS.** `quick_validate.py treasury/refactoring-go-services` -> `ok` (EXIT=0). `check_links.py treasury/refactoring-go-services` -> `ok` (EXIT=0).
- SKILL.md = 148 lines (limit 500). Description ~349 chars (limit 1024). Name `refactoring-go-services` is lowercase kebab-case, == folder name, contains no `claude`/`anthropic`, no XML angle brackets in frontmatter. No README/CHANGELOG/aux docs in folder. No references/templates/scripts/schemas/examples subdirectories present.

#### Rubric Scores

| # | Dimension | Score | Note |
|---|-----------|-------|------|
| 1 | Trigger Quality | 5 | Exact user phrases ("clean up this Go service", "refactor this handler", "move logic into the domain layer", "remove Go service code smells") plus explicit negative triggers. |
| 2 | Scope Focus | 5 | Singular responsibility: incremental, behavior-preserving Go refactoring. Explicitly excludes features, bug fixes, frontend, greenfield. |
| 3 | Workflow Clarity | 5 | Section 5 is a numbered imperative 5-step loop with clear entry/exit (one smell -> one action -> minimal change -> verify -> repeat until smell extinguished). |
| 4 | Output Contract | 3 | Deliverable is implied (in-place refactored code aligned to layer mapping) with a verify step, but no exact format/paths/naming for what the agent returns per iteration (no diff/report structure). |
| 5 | Token Efficiency | 4 | Lean 148-line body, no bloat; single tier so all guidance is inline, but short enough not to be costly. |
| 6 | Conflict Risk | 5 | Portable; no repo-policy clashes. Postgres appears only as an illustrative adapter example, not a hardcoded rule. |
| 7 | Reusability | 5 | Drops into any Go DDD/CQRS service project; not bound to one repo or file. |
| 8 | Security & Safety | 5 | No scripts, no external HTTP, no secret/env access, no destructive commands, no broad permissions. |
| 9 | Frontmatter Correctness | 5 | Valid YAML, kebab name == folder, description < 1024 with trigger language, no XML angle brackets. |
| 10 | Progressive Disclosure | 4 | Content lives in exactly one tier with no duplication; minor deduction since the technique catalog + layer map could move to references/, but body is not bloated. |
| | **Total** | **46/50** | rubricPass = false (Output Contract = 3 < 4). |

#### Findings

| Severity | File | Finding | Fix |
|----------|------|---------|-----|
| Medium | SKILL.md (Sec 5 step 4 / overall) | No explicit output contract: the skill never states what the agent must return per iteration (e.g., a smell-and-action log, a unified diff, a "behavior preserved: how verified" note) or in what format. Verification is described qualitatively ("State how the isolated behavior works now") rather than as a required deliverable shape. | Add a short Output Contract section: per iteration emit (a) detected smell + chosen action, (b) the minimal diff, (c) the verification performed (tests run / build status). Define stop condition output. |
| Low | SKILL.md lines 117-139 (Appendix) | Pseudo-code inconsistency: BEFORE signature uses `balance float64`; AFTER calls `a.Balance.Subtract(Money{...})`, implying `Balance` is `Money`, but `Money` is only introduced as a free struct and the field-type transition is never shown. Illustrative only, not functional. | Show `Account.Balance Money` and a `Money.Subtract` method, or annotate the example as schematic. |
| Info | SKILL.md (whole) | All deep guidance (technique catalog, layer mapping) is inline. Acceptable at 148 lines, but if expanded, move detail to references/ to keep Tier 2 lean. | Optional: extract Sections 3-4 into references/ if the body grows. |

#### Anti-Pattern Sweep

| Anti-pattern | Status | Notes |
|--------------|--------|-------|
| 1. Do-Everything scope | Pass | One responsibility: incremental Go refactoring; adjacent tasks explicitly excluded. |
| 2. Weak/Vague triggers | Pass | Concrete user phrases + negative triggers. |
| 3. Context bloat | Pass | 148 lines, no oversized inline dumps. |
| 4. Repo-policy mixing | Pass | No AGENTS/CLAUDE.md rules embedded; does not edit policy files. |
| 5. Non-step-by-step workflow | Pass | Section 5 is a numbered imperative loop. |
| 6. Overriding user intent | Pass | Guardrails defer to user (no feature creep, stop when "good enough"). |
| 7. Generating target output instead of a skill | Pass | This is a skill defining a refactoring workflow, not a refactor of a specific file. |
| 8. README/aux doc inside folder | Pass | Folder contains only SKILL.md. |
| 9. Trigger-phrase absence | Pass | "Use when" + explicit example phrases present. |
| 10. Stale / context competition / near-duplicate | Pass | Distinct from refactor-decision (whether/scope), crafting-backend-code / implement-backend-feature / implementing-go-template-requirements (new code from spec), platform-go-service* (scaffold w/ repo conventions), review-backend-code (verification), technical-debt-management (prioritize/negotiate). No blatant overlap; noted only for cross-cutting review. |
| 11. Duplicated cross-tier content | Pass | Single tier; nothing duplicated. |
| 12. Hardcoded platform assumptions | Pass | DDD/CQRS/Fowler are methodology, not a repo; Postgres is one illustrative adapter example only. |

#### Security Sweep

- **External HTTP:** Clear — no network calls; no scripts.
- **Secret/env access:** Clear — no env vars, credentials, or secret reads.
- **Destructive commands:** Clear — no shell/file-deletion logic; refactoring guidance only.
- **Broad permissions:** Clear — no permission requests; no scripts/ to scope.
- **Vendor bias:** Clear — Postgres/Kafka/JWT named only as generic illustrative examples, no vendor steering.

#### Delta List

1. (Required to reach rubricPass) Add an explicit Output Contract section specifying the per-iteration deliverable: detected smell + chosen refactor action, the minimal diff produced, and the verification evidence (tests/build). Define the stop-condition output.
2. (Recommended) Fix the appendix pseudo-code so the `Account.Balance` field type transition to `Money` is shown, or mark the snippet as schematic.
3. (Optional) If Sections 3-4 grow, relocate the technique catalog and layer mapping to references/ to preserve a lean Tier 2 body.

#### Final Recommendation

**Ship.** The skill is secure, well-triggered, well-scoped, portable, and gate-1 clean (46/50). It is flagged as an exception only because Output Contract scored 3, which drops rubricPass below the all-dimensions->=4 bar. The single medium finding (add an explicit per-iteration output contract) is a small, low-risk addition that does not block shipping; address it in the next light pass.

---

### planning-banking-tests

#### Summary

- **Skill**: planning-banking-tests
- **Folder**: /Users/admin/Desktop/project/research/skillify-foundation/treasury/planning-banking-tests
- **Date**: 2026-05-29
- **Mode**: Audit (read-only)
- **Result**: Pass with fixes (the authored SKILL.md + references + schemas + shipping scripts are strong; the folder ships non-skill run/audit artifacts and trips gate 1)

The skill converts a completed BA brief (from `eliciting-banking-brief` v1.2+) into a structured QA test plan: canonical `output.json` plus a deterministic markdown tree rendered by `scripts/render_test_plan_tree.py`. The core authoring is high-quality. The defects are folder-hygiene: committed run output, research artifacts, foreign-machine absolute paths, a stale skill name, and a gate-1 frontmatter parse failure.

#### Deterministic Checks

- `quick_validate.py treasury/planning-banking-tests` → **EXIT 1 (FAIL)**. Errors are ONLY "unexpected indented frontmatter line" for each key under the nested `metadata:` mapping (`version`, `stage_type`, `input_schema`, `output_schema`, `banking_grade`, `recommended_temperature`, `tier_review_levels`, `expected_duration_p95_seconds`, `max_retries_recommended`). Per the binding gate-1/frontmatter calibration this is the validator's known inability to parse nested mappings / object/list fields that `frontmatter-guide.md` explicitly permits; the frontmatter is valid YAML and loads in real Claude Code. Recorded as FAIL but NOT used to push Frontmatter Correctness below 4.
- `check_links.py treasury/planning-banking-tests` → **EXIT 0 (ok)**. All `references/`, `schemas/`, `scripts/` links in SKILL.md resolve (including `references/v1.1-role-boundaries.md` and `tests/assertions/`).
- **Latent secondary gate-1 violation (masked)**: `quick_validate.py` returns early when frontmatter parse errors exist (lines 89-90), so it never reaches its `rglob` banned-doc scan (line 138). Three `README.md` files exist under `runs/.../test-plan-*/` and `runs/.../test-plan-normalized/`. The instant the frontmatter is reshaped to satisfy the parser, gate 1 will fail again on "banned human-facing doc inside skill folder". This is a real defect hidden behind the parser short-circuit.

#### Rubric Scores

| Dimension | Score | Note |
|---|---|---|
| Trigger Quality | 5 | Tier-1 desc names exact intents (Stage 4c after BA Stage 1; per-story roster from a structured brief; sign-off from governance gaps) + 4 explicit negative triggers. |
| Scope Focus | 5 | Singular: BA brief → test PLAN, never test code; defers code-gen, defect-filing, framework choice, coverage-on-running-systems. |
| Workflow Clarity | 5 | 12 imperative numbered steps, each loads its reference with entry/exit + failure codes; TM-01..TM-12 failure table. |
| Output Contract | 5 | Exact: `output.json` per `schemas/output.json` (discriminated `output_type`), named deterministic renderer with exact dir/file names, reduced 3-file failure tree. |
| Token Efficiency | 5 | Body 154 lines, lean; all depth in 11 references + 2 schemas + 2 scripts. No inline duplication. |
| Conflict Risk | 4 | Portable body (compatibility: claude-code/codex/opencode, no host hardcoding, sibling-skill boundaries are healthy). Minus 1: references + run-artifact scripts carry a foreign user's absolute paths and a stale skill name (non-portable committed content, though no AGENTS/CLAUDE.md clash and no inter-skill instruction fight). |
| Reusability | 4 | Domain-scoped (banking/fintech QA grounded in `eliciting-banking-brief` output) — guards keep that at 4-5. Held at 4: the shipped folder is not a clean drop-in because `runs/` harness scripts hardcode one machine's paths + an old skill name and must be pruned first. |
| Security & Safety | 5 | All Python stdlib, deterministic, no network / secret / env / subprocess / destructive ops. Domain logic enforces no-real-PII (AP-Q7) and no-tipping-off (AP-Q11). |
| Frontmatter Correctness | 4 | Valid YAML; kebab `name` == folder; description 609 chars with "Use when" triggers; no XML angle brackets. Held at 4 (not 5): it does trip the repo's own gate-1 tool and needs a frontmatter reshape or validator fix to clear the hard contract. |
| Progressive Disclosure | 3 | Body/references/schemas/scripts tiers are clean — BUT the folder also commits `runs/` (generated output trees incl. 3 banned README.md + `__pycache__/*.pyc`) and `audit/phases/` (research artifacts). That content belongs to no tier, and the READMEs are a latent gate-1 banned-doc violation. |
| **Total** | **45/50** | rubricPass = false (Progressive Disclosure = 3 < 4). |

#### Findings

| Severity | File | Finding | Fix |
|---|---|---|---|
| High | SKILL.md frontmatter (`metadata:` block, lines 13-22) | Gate-1 FAIL: nested `metadata:` mapping + inline object/list fields are unparseable by `quick_validate.py`. Valid YAML, but violates the hard validation gate. | Either flatten frontmatter to the validator-parseable subset, or land the known validator/guide reconciliation; until then the skill cannot clear gate 1. |
| High | runs/001-ecommerce-multi-epic/test-plan-11111111/README.md; .../test-plan-11111111-r2/README.md; .../test-plan-normalized/README.md | Three `README.md` inside the skill folder = banned-doc violation. Currently masked because the frontmatter parse error makes `quick_validate.py` return before its rglob scan. Will fire once frontmatter is fixed. | Remove the entire `runs/` directory from the skill folder (it is generated example output, not skill content); relocate to a gitignored examples/fixtures area outside the skill if retention is needed. |
| Medium | runs/ (whole tree, ~2.7M incl. `__pycache__/merge.cpython-314.pyc`, `merge.py`, `baseline/diff.py`) | Committed run output + ad-hoc transformation/diff scripts + a compiled `.pyc` bloat the skill folder and live in no progressive-disclosure tier. | Delete `runs/` and `__pycache__` from the published skill; if a worked example is desired, ship a single minimal fixture under `examples/`, not a full multi-epic render tree. |
| Medium | audit/phases/phase-a1-skill-genesis.md | Research/genesis artifact committed inside the skill folder; carries `skill_ref: qa-plan-from-brief` (stale name) and `/Users/IF640063/...` plan + holdout paths. Not part of any skill tier. | Move design/audit phase docs out of the skill folder (e.g., into the gitignored design-methodology area). |
| Medium | runs/.../merge.py; runs/.../baseline/diff.py | Hardcode `/Users/IF640063/Desktop/example/qa-bootstrap-kit/...` (a different developer's machine) and the old skill name `qa-plan-from-brief`. Non-portable; would not run for anyone else. (Note: these are run-harness scripts, not the shipping `scripts/`.) | Remove with the `runs/` tree; do not ship machine-specific harness scripts in the skill. |
| Low | references/nfr-derivation.md; references/compliance-test-patterns.md; references/test-data-design.md; references/v1.1-role-boundaries.md | Contain `/Users/IF640063/...` absolute paths in prose/examples — leak one developer's filesystem into shipped reference content. | Replace foreign absolute paths with relative or placeholder paths. |
| Low | schemas/input.json + output.json (`$id`), various run/output.json (`created_by`) | `$id` uses `https://qa-bootstrap-kit/skills/planning-banking-tests/...` and several generated outputs stamp `created_by: qa-plan-from-brief-v1.0.0` — residual old project/skill naming. Cosmetic; schemas are valid draft-07 and load fine. | Optional: align `$id` / `created_by` strings to the current skill name on next revision. |

#### Anti-Pattern Sweep

| Anti-pattern | Status | Notes |
|---|---|---|
| 1. Do-Everything scope | Pass | Tight single responsibility (brief → test plan); explicit deferrals. |
| 2. Weak/Vague triggers | Pass | Exact intents + 4 negative triggers in tier-1 desc and "When to use". |
| 3. Context bloat | Mitigated | Body lean and well-tiered; bloat is the committed `runs/` (2.7M) + `audit/`, which are non-skill artifacts to be pruned, not inline body bloat. |
| 4. Repo-policy mixing | Pass | No AGENTS/CLAUDE.md edits; no repo-policy content embedded. |
| 5. Non-step-by-step workflow | Pass | 12 imperative numbered steps with entry/exit + failure codes. |
| 6. Overriding user intent | Pass | Refuses to invent thresholds (AP-Q9), emits TBD-pending-OQ; asks BA to resolve. |
| 7. Generating target output instead of a skill | Pass | This IS a skill; it explicitly refuses to emit test code (AP-Q1). |
| 8. README/aux doc inside folder | Fail | 3 `README.md` under `runs/...` (latent gate-1 banned-doc, masked by frontmatter parser early-return). |
| 9. Trigger-phrase absence | Pass | "Use when" present (3x) + "Do NOT use when". |
| 10. Stale-skill / context competition / near-duplicate | Pass | Distinct from siblings (`generating-gherkin-acceptance-criteria`, `validating-banking-implementation`, `testing-strategy`); clean upstream/downstream boundaries. No catalog near-duplicate. |
| 11. Duplicated cross-tier content | Pass | SKILL.md points to references rather than restating; no body/reference duplication observed. |
| 12. Hardcoded platform assumptions | Mitigated | SKILL.md + shipping `scripts/` are host-portable (use `Path(__file__)`); but committed `runs/` harness scripts + some references hardcode one machine's absolute paths + a stale skill name. Confined to non-shipping artifacts; prune to fully clear. |

#### Security Sweep

- **External HTTP**: Clear — no `requests`/`urllib`/`socket`/URLs in any Python; only HTTP appears as JSON-Schema `$schema`/`$id` identifiers (not fetched).
- **Secret / env access**: Clear — no `os.environ`/`getenv`/secrets/tokens. The only `password` hit is a PII-field-name string literal in `normalize_fragments.py`; `token` hits are local tokenizer variable names in `diff.py`.
- **Destructive commands**: Clear — no `rm`/`shutil.rmtree`/`os.remove`/`os.unlink`/`subprocess`/`os.system`. Scripts only read inputs and write into an explicit `--output-dir`/`--out-dir`.
- **Broad permissions**: Clear — scoped filesystem I/O to caller-provided paths; deterministic, no clock/uuid/random at render time.
- **Vendor bias**: Clear — no vendor lock-in; `compatibility: [claude-code, codex, opencode]`, host-neutral body.

#### Delta List

1. (Required, blocks gate 1) Resolve the frontmatter parse failure: flatten the `metadata:` mapping to the validator-parseable subset, or land the documented validator/guide reconciliation. SKILL.md must make `quick_validate.py` exit 0.
2. (Required, blocks gate 1 once #1 lands) Delete the entire `runs/` directory from the skill folder — it carries 3 banned `README.md`, a `__pycache__/*.pyc`, and ad-hoc `merge.py`/`diff.py` with foreign-machine paths.
3. (Required) Remove `audit/phases/` from the skill folder; relocate genesis/audit docs to the gitignored design area.
4. (Recommended) Purge `/Users/IF640063/...` absolute paths from the four references that contain them; use relative/placeholder paths.
5. (Recommended) Replace residual `qa-plan-from-brief` naming and the `qa-bootstrap-kit` `$id` host with the current skill name in schemas/output stamps on next revision.
6. (Recommended) If a worked example is valuable, ship one minimal fixture under `examples/` instead of a full multi-epic render tree.

#### Final Recommendation

**Refactor before ship.** The authored skill — SKILL.md (12 well-conditioned steps, exact output contract, failure-mode table), 11 references, two valid draft-07 schemas, and two clean deterministic stdlib scripts — is genuinely strong and would score 5s across the board on its own. It is held back entirely by folder hygiene: a 2.7M `runs/` output tree (3 banned READMEs + `.pyc` + harness scripts), an `audit/` research artifact, foreign-user absolute paths, a stale skill name, and a gate-1 frontmatter parse failure. None of these touch the skill's logic or security; all are removable. After the deltas above, this skill is a Ship candidate.

---

## 5. Appendix — live gate-1 result for all skills

| Skill | Gate-1 | Detail |
|-------|--------|--------|
| `analyzing-banking-requirements` | PASS | ok |
| `assembling-tl-handoff` | PASS | ok |
| `assessing-ba-feasibility` | PASS | ok |
| `checking-ba-governance` | PASS | ok |
| `drafting-ba-stories` | PASS | ok |
| `eliciting-banking-brief` | PASS | ok |
| `evaluate-banking-compliance` | PASS | ok |
| `extract-brief-structure` | PASS | ok |
| `orchestrate-banking-brief-pipeline` | PASS | ok |
| `running-ba-pipeline` | PASS | ok |
| `running-business-analysis-workflow` | PASS | ok |
| `scoping-ba-intake` | PASS | ok |
| `scoping-technical-requirements` | PASS | ok |
| `sweep-ambiguities` | PASS | ok |
| `api-contract-design` | PASS | ok |
| `architecting-fintech-systems` | PASS | ok |
| `architecture-decision` | PASS | ok |
| `data-modeling` | PASS | ok |
| `defining-engineering-standards` | PASS | ok |
| `delivery-planning` | PASS | ok |
| `designing-tech-lead-handoff` | PASS | ok |
| `domain-modeling` | PASS | ok |
| `engineer-growth-planning` | PASS | ok |
| `engineering-doc-planning` | PASS | ok |
| `generate-ux-pack` | PASS | ok |
| `integration-design` | PASS | ok |
| `model-selection` | PASS | ok |
| `model-selection-updated` | PASS | ok |
| `pr-design-review` | PASS | ok |
| `production-readiness` | PASS | ok |
| `refactor-decision` | PASS | ok |
| `risk-estimation` | PASS | ok |
| `technical-debt-management` | PASS | ok |
| `technical-feasibility` | PASS | ok |
| `crafting-backend-code` | PASS | ok |
| `crafting-frontend-code` | PASS | ok |
| `implement-backend-feature` | PASS | ok |
| `implement-frontend-feature` | PASS | ok |
| `implementing-go-template-requirements` | PASS | ok |
| `platform-common` | PASS | ok |
| `platform-go-service` | PASS | ok |
| `platform-go-service-common` | PASS | ok |
| `refactoring-go-services` | PASS | ok |
| `review-backend-code` | PASS | ok |
| `review-frontend-code` | PASS | ok |
| `review-rust-code` | PASS | ok |
| `generating-gherkin-acceptance-criteria` | PASS | ok |
| `planning-banking-tests` | FAIL | quick_validate: banned human-facing doc inside skill folder: runs/001-ecommerce-multi-epic/test-plan-11111111/README.md · quick_validate: banned human-facing doc inside skill folder: runs/001-ecommerce-multi-epic/test-plan-11111111-r2/README.md · quick_validate: banned human-facing doc inside skill folder: runs/001-ecommerce-multi-epic/test-plan-normalized/README.md |
| `testing-strategy` | PASS | ok |
| `validating-banking-implementation` | PASS | ok |
| `agent-context-initializer` | PASS | ok |
| `agentic-workflow-design` | PASS | ok |
| `authoring-scaffold-profile` | PASS | ok |
| `composing-agent-pipelines` | PASS | ok |
| `developing-langgraph-workflows` | PASS | ok |
| `drafting-stage-prompt` | PASS | ok |
| `multi-agent-handoff-architect` | PASS | ok |
| `orchestrating-agent-scaffold` | PASS | ok |
| `orchestrating-openclaw-squad` | PASS | ok |
| `reviewing-implement-gate` | PASS | ok |
| `augment-diagrams` | PASS | ok |
| `cross-examine` | PASS | ok |
| `extract-findings` | PASS | ok |
| `frame-debate` | PASS | ok |
| `opening-debate-panel` | PASS | ok |
| `plan-research` | PASS | ok |
| `report-debate` | PASS | ok |
| `reporting-research-run` | PASS | ok |
| `research-vault-librarian` | PASS | ok |
| `review-report` | PASS | ok |
| `revise-positions` | PASS | ok |
| `search-sources` | PASS | ok |
| `synthesize-consensus` | PASS | ok |
| `synthesize-report` | PASS | ok |
| `configuring-sandbox-allowlist` | PASS | ok |
| `devops-infrastructure-hardener` | PASS | ok |
| `governance-policy-generator` | PASS | ok |
| `reviewing-software-security` | PASS | ok |
| `universal-spec-validator` | PASS | ok |
| `authoring-workflow-postmortem` | PASS | ok |
| `incident-response` | PASS | ok |
| `observability-design` | PASS | ok |
| `observability-telemetry-instrumenter` | PASS | ok |
| `performance-cost-review` | PASS | ok |
| `reviewing-agent-spend` | PASS | ok |
| `business-logic-extractor` | PASS | ok |
| `generating-pseudocode` | PASS | ok |
| `organizing-local-files` | PASS | ok |
| `progressive-bug-hunter` | PASS | ok |
| `publishing-git-review-requests` | PASS | ok |
| `rendering-readable-html` | PASS | ok |
| `thai-translator` | PASS | ok |

