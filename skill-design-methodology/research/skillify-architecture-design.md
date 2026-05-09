---
document_role: phase-2-architecture-design
pipeline_stage: 2 of 5 (Claude Opus 4.7 Max — synthesize principles + design skill architecture)
upstream: ../../literature/skill-design-methodology/agent-skill-design-principles.md, ../cross-model-skillification-pipeline.md, ../../literature/**
downstream: phase-3 GPT-5.5 xHigh — convert this design into the final SKILL.md framework
status: draft v1 (2026-04-26)
---

# Skillify — Skill-Creator Architecture Design

This document is the Phase-2 deliverable defined by `cross-model-skillification-pipeline.md`. It is **not** the final `SKILL.md`. It is the architecture spec that the next pipeline stage will compile into a clean, copy-ready `SKILL.md` framework.

Reading order for downstream agents:
1. Section 1 — synthesized principles (the rules every skill-creator output must satisfy)
2. Section 2 — the skill-creator's own architecture (what `/skillify` *is*)
3. Section 3 — operating modes and per-mode workflow shape
4. Section 4 — audit of the existing `/skillify` draft with a concrete delta list
5. Section 5 — explicit hand-off contract to Phase 3
6. Section 6 — open questions for the human reviewer

---

## 1. Synthesized Principles

These eleven principles are the load-bearing constraints. Every later design choice traces back to one of them. Provenance is noted in brackets: `[methodology]` = `literature/skill-design-methodology/agent-skill-design-principles.md`; `[openai-creator]` = `literature/openai/openai_skill_creator.md`; `[skillify]` = the existing `/skillify` references.

### P1 — Progressive disclosure is the controlling architecture
Three loading tiers, in strict cost order. `[methodology §1, openai-creator §"Progressive Disclosure", skillify/references/progressive-disclosure.md]`

| Tier | Artifact | Budget | When loaded |
|------|----------|--------|-------------|
| 1 | YAML frontmatter (`name`, `description`) | ~100 words, always in context | Trigger decision |
| 2 | `SKILL.md` body | <5,000 words, <500 lines | After skill triggers |
| 3 | `references/`, `scripts/`, `assets/` | unlimited | Loaded by agent on demand; scripts can run without being read |

Implication: **anything that doesn't have to be in Tier 2 must move to Tier 3.** Variant playbooks, deep API docs, schemas, and large examples are Tier-3 by default.

### P2 — Frontmatter is the only thing that decides whether the skill loads
`description` is the trigger surface. It must contain three components in order: **what the skill does · when to use it (literal trigger phrases) · key capabilities**. `[methodology §2, openai-creator §"Frontmatter", skillify/references/frontmatter-guide.md]`

Hard rules:
- `name`: kebab-case, ≤64 chars, matches folder name, no `claude` / `anthropic` prefix.
- `description`: ≤1024 chars, no XML angle brackets (injection vector), trigger phrases verbatim from how a user actually speaks.
- "When to use this skill" sections inside the body are an anti-pattern — the body is only read *after* the trigger fires.

### P3 — Match degree of freedom to task fragility
Three bands. Choose deliberately. `[methodology §4, openai-creator §"Set Appropriate Degrees of Freedom"]`

- **High** (text instructions): creative work, multiple valid paths, context-dependent decisions.
- **Medium** (pseudocode / parameterized scripts): preferred pattern with allowed variation.
- **Low** (rigid scripts, few parameters): fragile, irreversible, or strictly sequential operations.

Rule: "Code is deterministic; language interpretation isn't." Critical validations and irreversible side effects belong in scripts, not prose.

### P4 — Skills are for agents, not humans
A skill folder must contain only files needed for the agent to do the job. **No** `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, or `CONTRIBUTING.md` inside the skill directory. Repo-level docs live at the repo root, not inside the skill. `[methodology §3, openai-creator §"What to Not Include", skillify/references/anti-patterns.md #8]`

### P5 — One responsibility per skill
A skill that "does everything" trains the agent to load it for nothing in particular. Prefer narrow, verb-led scope. Split when the description starts needing the word "and" twice. `[skillify/references/anti-patterns.md #1, skillify/references/validation-rubric.md #2]`

Naming convention (Anthropic best practice): gerund-style for activity skills (`processing-pdfs`, `testing-code`), verb-led phrases for action skills (`onboard-customer`, `rotate-pdf`). Namespace by tool when it improves trigger clarity (`gh-address-comments`, `linear-create-issue`). `[agent-skill-source-index.md priority 2]`

### P6 — Imperative, numbered, step-by-step
Bodies are instruction sets, not essays. Numbered imperative steps with explicit entry/exit conditions outperform bulleted theory. `[methodology §5, skillify/references/anti-patterns.md #5, skillify/references/validation-rubric.md #3]`

### P7 — Output contract must be explicit
The skill must state the exact deliverable shape (file list, file format, naming, location). "Produce good code" is a 1/5 contract; "Output a `jest` test file at `src/__tests__/<name>.test.ts` with one `describe` per public function" is a 5/5 contract. `[skillify/references/validation-rubric.md #4]`

### P8 — Information lives in exactly one place
No duplication across tiers. `SKILL.md` *points to* references; it does not copy them. Repo-policy belongs in `AGENTS.md` / `CLAUDE.md`, not in a portable skill. `[skillify/references/anti-patterns.md #4 and #11]`

### P9 — References are one level deep
All Tier-3 files link directly from `SKILL.md`. No nested `references/sub/...` chains. Files >100 lines need a table of contents. Files >10k words should expose grep patterns in `SKILL.md` so the agent can target sections. `[methodology §1, openai-creator §"Important guidelines"]`

### P10 — Validation is a gate, not a suggestion
Skills must pass a structured rubric before deployment. Critical fields (frontmatter validity, name format, missing triggers, XML in description) should be checked deterministically by a script — not left to the agent's judgment. `[methodology §5/§6, openai-creator §"Step 5: Validate", skillify/references/validation-rubric.md]`

### P11 — Skills are living documents
Iterate on real traces. Under-triggering → add triggers. Over-triggering → add negative triggers / narrow scope. Bloat → split or move to references. Inconsistent execution → tighten degree of freedom or add scripts. `[methodology §5 step 5, openai-creator §"Step 6: Iterate", skillify/SKILL.md §Step 6]`

---

## 2. Skillify Architecture

### 2.1 Identity

```yaml
name: skillify
type: meta-skill
classification: workflow-automation (creates other skills)
output: skill-files-only (never the target code)
```

Skillify is a **meta-skill**: it produces other skills. This single fact governs three rigid invariants:

- **Invariant A (Output discipline)**: Output MUST be skill files. If the user asks "create a skill that generates React components," skillify outputs a `react-components/` skill folder, not React components. `[skillify/references/anti-patterns.md #7]`
- **Invariant B (Self-conformance)**: Skillify must itself satisfy every principle in §1. The meta-skill is the reference implementation.
- **Invariant C (Portability)**: Skillify outputs must work across Claude Code, Codex, Copilot, Gemini, and OpenCode. Platform-specific accommodations live in the `compatibility` field and the `platform-compatibility.md` reference, not in the body.

### 2.2 Operating modes

The current `/skillify/SKILL.md` description advertises eight modes. They are real, distinct workflows with different inputs, validation gates, and outputs. They must be discoverable from the body, not just promised in the description.

| # | Mode | Trigger phrases | Input | Output | Primary degree of freedom |
|---|------|-----------------|-------|--------|---------------------------|
| 1 | **Create** | "create a skill", "write a SKILL.md", "build a skill for X" | User intent + concrete examples | Full skill folder | High → Low (collapses as design firms up) |
| 2 | **Refactor** | "refactor my skill", "improve this skill", "fix this SKILL.md" | Existing skill folder + pain point | Modified skill folder | Medium |
| 3 | **Review** | "review my skill", "is this skill good", "feedback on my SKILL.md" | Existing skill | Review report (no edits) | High |
| 4 | **Audit** | "audit my skill", "score my skill", "check this skill against the rubric" | Existing skill | Scored rubric + delta list | Low (rubric is rigid) |
| 5 | **Compress** | "shrink this skill", "reduce token cost", "compress my SKILL.md" | Existing skill (over budget) | Slimmer skill, deeper Tier-3 | Medium |
| 6 | **Split** | "split this skill", "this skill is doing too much" | Over-broad skill | Two or more focused skills | Medium |
| 7 | **Merge** | "merge these skills", "combine X and Y skills" | Two overlapping skills | One unified skill (or rejected with reason) | Medium |
| 8 | **Adapt** | "adapt this for Codex", "make this work in Copilot/Gemini" | Skill + target platform | Platform-adapted skill | Low (mechanical translation) |

Modes are *not* a menu the user picks from explicitly. The skill-creator infers the mode from the user's phrasing (the trigger phrases above) and confirms ambiguous cases with a single clarifying question.

### 2.3 Mode → workflow shape

Each mode follows a different shape, but they share a universal preamble (Step 0 — mode detection + intake) and a universal exit (Validation gate + iteration prompt).

```
              ┌─ Step 0: Mode detection + intake (universal)
              │
              ▼
   ┌──────────┴──────────┬──────────┬───────────┬──────────┬─────────┬────────┬────────┐
   │                     │          │           │          │         │        │        │
 CREATE              REFACTOR    REVIEW       AUDIT      COMPRESS   SPLIT    MERGE    ADAPT
   │                     │          │           │          │         │        │        │
   │ 6-step              │ 4-step   │ 3-step    │ 3-step   │ 4-step  │ 5-step │ 4-step │ 3-step
   │ pipeline            │ pipeline │ rubric+   │ rubric   │ tier-   │ scope  │ scope  │ platform
   │ (existing           │ (diff +  │ narrative │ scoring  │ migrate │ map +  │ check  │ matrix
   │  Steps 1–6)         │  retest) │  feedback │ + report │ to T3   │ extract│ + dedup│ + transform
   │                     │          │           │          │         │        │        │
   └──────────┬──────────┴──────────┴───────────┴──────────┴─────────┴────────┴────────┘
              ▼
      Validation gate (rubric ≥40/50, frontmatter check, anti-pattern sweep)
              ▼
      Iterate-or-ship prompt
```

### 2.4 Canonical directory layout

```
skillify/
├── SKILL.md                          # Tier 2 — body, ≤500 lines
├── references/                       # Tier 3 — loaded on demand, one level deep
│   ├── frontmatter-guide.md          # ✓ exists — all YAML fields, good/bad examples
│   ├── anti-patterns.md              # ✓ exists — 12 named anti-patterns
│   ├── validation-rubric.md          # ✓ exists — 10-dimension scoring
│   ├── progressive-disclosure.md     # ✓ exists — 3-tier model + splitting patterns
│   ├── workflow-patterns.md          # ✓ exists — 5 structural patterns
│   ├── platform-compatibility.md     # ✓ exists — Claude/Codex/Copilot/Gemini matrix
│   ├── security-checklist.md         # ✓ exists — pre-enable audit
│   ├── mode-playbooks.md             # ✗ NEW — per-mode operating instructions
│   └── lifecycle-and-iteration.md    # ✗ NEW — what to do after first ship
├── templates/                        # Tier 3 — boilerplate output skeletons
│   ├── basic-skill-template.md       # ✓ exists
│   ├── mcp-skill-template.md         # ✓ exists
│   ├── domain-skill-template.md      # ✓ exists
│   └── audit-report-template.md      # ✗ NEW — output shape for Audit/Review modes
├── scripts/                          # Tier 3 — deterministic checks (NEW DIR)
│   ├── init_skill.py                 # ✗ NEW — boilerplate generator (P3 Low-freedom)
│   ├── quick_validate.py             # ✗ NEW — frontmatter + structure validator (P10)
│   └── check_links.py                # ✗ NEW — verify all references/* pointers resolve
├── examples/                         # Tier 3 — worked walkthroughs
│   ├── good-description-examples.md  # ✓ exists
│   ├── skill-audit-walkthrough.md    # ✓ exists
│   └── create-from-scratch.md        # ✗ NEW — end-to-end Create mode trace
└── research/skillify-architecture-design.md # this file (Phase-2 design — NOT shipped with skill)
```

The design file itself is a build artifact. Phase 3 may either keep it (renamed) as `references/architecture.md`, or omit it from the shipped skill if the synthesized SKILL.md is self-sufficient. Default recommendation: **omit from shipped skill** to honor P4 (no human-facing meta-docs).

### 2.5 File contracts — what each file owns

To honor P8 (information lives in exactly one place), each file has a single owned topic. No topic appears in two files.

| File | Owned topic | Forbidden topic |
|------|-------------|-----------------|
| `SKILL.md` | Mode detection + universal preamble + per-mode entry pointer + validation gate | Per-mode deep workflows (point to mode-playbooks.md) |
| `references/frontmatter-guide.md` | All YAML fields, security restrictions on frontmatter, good/bad description examples | Workflow steps |
| `references/anti-patterns.md` | Named anti-patterns with fix recipe | Rubric scoring |
| `references/validation-rubric.md` | 10-dimension rubric + scoring template | Anti-pattern catalog |
| `references/progressive-disclosure.md` | 3-tier model + when/how to split content | Per-mode workflow |
| `references/workflow-patterns.md` | 5 structural patterns for arbitrary skills | Mode playbooks for skillify itself |
| `references/platform-compatibility.md` | Cross-platform matrix + adaptation recipes | Frontmatter security |
| `references/security-checklist.md` | Pre-enable safety audit | Validation rubric scoring |
| `references/mode-playbooks.md` *(new)* | Per-mode workflow detail (Refactor through Adapt) | Anti-patterns / rubric |
| `references/lifecycle-and-iteration.md` *(new)* | Post-ship iteration heuristics, retirement criteria | Initial creation |
| `templates/*` | Output skeletons | Workflow narrative |
| `scripts/init_skill.py` *(new)* | Boilerplate folder generation | Validation |
| `scripts/quick_validate.py` *(new)* | Deterministic structure + frontmatter checks | Boilerplate |
| `scripts/check_links.py` *(new)* | Reference link integrity | Frontmatter |
| `examples/*` | Worked traces | Reference material |

---

## 3. Workflow Design

### 3.1 Universal preamble (Step 0 — every mode runs this)

```
1. Detect mode from user phrasing.
   - Map the request to the trigger-phrase table in §2.2.
   - If two modes match (e.g., "improve and shrink"), default to the more
     conservative one (Refactor over Compress) and note both in the plan.
   - If no mode matches, ask one disambiguating question; do not guess.

2. Establish the input.
   - Create: ask for concrete user-trigger examples (3+ preferred).
   - Refactor/Review/Audit/Compress/Split/Merge/Adapt: read the existing
     SKILL.md and enumerate references/ files before proposing changes.

3. Confirm the output contract.
   - State explicitly what files will be produced or modified.
   - For destructive modes (Refactor, Compress, Split, Merge), confirm
     before overwriting; produce a diff or new directory unless told to
     overwrite.
```

### 3.2 Per-mode workflows

Only **Create** lives in `SKILL.md` body (it is the most common mode and the most instruction-heavy). The other seven live in `references/mode-playbooks.md` with the body holding only a one-paragraph entry pointer per mode.

#### Create (in SKILL.md body)
The existing 6-step workflow in `/skillify/SKILL.md` is correct:

1. Understand the target task (concrete examples → trigger phrases)
2. Plan reusable contents (scripts / references / assets)
3. Define the trigger condition (frontmatter)
4. Write the SKILL.md body (imperative numbered)
5. Validate (rubric + structure check)
6. Iterate (real-trace feedback loop)

Two refinements to add in Phase 3:
- **Step 1 must include "elicit at least 3 concrete user prompts"** before naming the skill. Without examples, scope drifts. `[openai-creator §"Step 1: Understanding the Skill with Concrete Examples"]`
- **Step 5 must run `scripts/quick_validate.py` if available** — deterministic check before agent-judgment rubric scoring. `[methodology §5 step 4]`

#### Refactor (in mode-playbooks.md)
1. Read current SKILL.md + all references; build inventory.
2. Diff against rubric — flag every dimension scoring <4/5.
3. Apply targeted fixes; do not rewrite passing dimensions.
4. Re-score; show before/after rubric table.

#### Review (in mode-playbooks.md)
1. Read SKILL.md + references; do not modify anything.
2. Score rubric; highlight top 3 risks.
3. Produce review report (template: `templates/audit-report-template.md`).

#### Audit (in mode-playbooks.md)
1. Run `scripts/quick_validate.py` for deterministic checks.
2. Score full rubric.
3. Generate machine-readable report (rubric scores + delta list + security findings).

Note: Review is narrative + qualitative; Audit is structured + quantitative. Keep them distinct.

#### Compress (in mode-playbooks.md)
1. Measure: line count, token estimate, references count.
2. Identify candidates for migration to Tier 3 (anything over ~50 lines that isn't core workflow).
3. Move content; replace with one-line pointer.
4. Re-validate; ensure no broken references.

#### Split (in mode-playbooks.md)
1. Map current scope; list every distinct task the skill handles.
2. Cluster tasks by domain; propose 2+ focused skills.
3. Generate new skill folders; allocate references appropriately.
4. Update each new SKILL.md description with sharper triggers + negative triggers pointing to siblings.

#### Merge (in mode-playbooks.md)
1. Compute scope overlap; reject merge if overlap <70% (recommend keeping separate).
2. Unify trigger phrase set; resolve naming conflict.
3. Combine references; deduplicate.
4. Score combined skill against P5 (one responsibility); abort if it now violates.

#### Adapt (in mode-playbooks.md)
1. Read source platform conventions vs. target (use `references/platform-compatibility.md`).
2. Translate frontmatter, file paths, tool references.
3. Preserve workflow logic; replace platform-specific tool names with generic equivalents.
4. Add `compatibility:` field listing both platforms if dual-target.

### 3.3 Validation gate

Every mode exits through this gate. Failure means the skill is not shipped.

```
Gate A — Deterministic (script):
  □ scripts/quick_validate.py exits 0
  □ Frontmatter is valid YAML, delimited by ---
  □ name is kebab-case, ≤64 chars, matches folder
  □ description ≤1024 chars, no XML angle brackets
  □ No README/CHANGELOG/INSTALLATION_GUIDE inside skill folder
  □ All references/* and templates/* links resolve

Gate B — Rubric (agent judgment, scored):
  □ All 10 dimensions ≥ 4/5
  □ Total ≥ 40/50

Gate C — Anti-pattern sweep (agent + checklist):
  □ Walk all 12 entries in references/anti-patterns.md
  □ For each, confirm absence with a one-line note

Gate D — Security (agent + checklist):
  □ Walk references/security-checklist.md
  □ Flag any external HTTP, env-var access, destructive commands,
    or vendor-bias in scripts/ or body
```

If Gate A fails: do not advance; surface the script's error verbatim.
If Gate B/C/D flags issues: report, do not auto-fix; ask the user before changes.

### 3.4 Iteration loop

After deployment, the skill enters lifecycle. `references/lifecycle-and-iteration.md` (new) holds the runbook:

- **Under-trigger signal**: skill description matched the user intent but the agent didn't load it. Action: add the missed phrasing to `description`.
- **Over-trigger signal**: skill loaded for an unrelated task. Action: add a negative trigger (`Do NOT use for ...`) or sharpen the positive triggers.
- **Bloat signal**: SKILL.md grew past 500 lines or 5,000 tokens. Action: run Compress mode.
- **Drift signal**: agent execution diverges from intent across multiple traces. Action: tighten degree of freedom (move from text → pseudocode → script).
- **Retirement signal**: skill hasn't earned its context cost in 90 days. Action: archive — every loaded skill competes for the context window.

---

## 4. Audit of Existing /skillify

Scoring the current draft against the §3.3 gates and the §1 principles.

### 4.1 What is already strong (preserve)

1. **Description structure** in current `SKILL.md` correctly leads with verb phrases ("Create, refactor, review...") and includes literal trigger phrases. Good P2 conformance.
2. **Negative trigger present**: `Do NOT use for generating target code — output MUST be skill files only.` Good P2 + Invariant A reinforcement.
3. **Six-step workflow** is imperative + numbered. Good P6 conformance.
4. **Reference table at the bottom of SKILL.md** is the right pointer pattern (Tier 2 → Tier 3). Good P9.
5. **Constraints section uses DO NOT rules** in imperative form. Good P6.
6. **Output format is explicit** with directory tree. Good P7.
7. **No README/CHANGELOG/etc. in the skill folder.** Good P4.
8. **Frontmatter is valid YAML, kebab-case, no XML.** Passes Gate A static checks.
9. **Reference set is well-bounded** — seven references each owning a distinct topic. Good P8 conformance, except for partial overlap noted below.
10. **Templates are differentiated** (basic / mcp / domain). Each addresses a real shape rather than being a single one-size-fits-all.

### 4.2 Gaps vs. principles

| # | Gap | Principle violated | Severity |
|---|-----|--------------------|----------|
| G1 | The SKILL.md describes 8 modes but its body documents only `Create`. Modes 2–8 are not discoverable from the body. | P5 (scope discoverability), P6 (clear workflows) | High |
| G2 | No `scripts/` directory exists; no `quick_validate.py`, no `init_skill.py`. Validation relies entirely on agent judgment. | P3 (low-freedom for fragile checks), P10 (deterministic gate) | High |
| G3 | Step 1 of Create workflow says "Ask: 'What specific tasks should this skill handle?'" but doesn't enforce eliciting concrete user examples first. | P11 (real examples) — but actually a Create-mode quality issue | Medium |
| G4 | Validation rubric exists as a reference but is not bound to a deterministic step in the workflow — "Score 4+/5" is left to agent judgment with no checklist artifact. | P10 (validation as a gate) | Medium |
| G5 | The Create workflow's Step 5 quick checklist duplicates parts of `validation-rubric.md` and `anti-patterns.md`. | P8 (one place) | Low |
| G6 | No troubleshooting section in `SKILL.md` itself (the `basic-skill-template.md` recommends one). | P6/P11 (iteration support) | Low |
| G7 | Examples cover descriptions and audit, but no end-to-end "from-scratch Create" walkthrough. | P11 (concrete examples) | Medium |
| G8 | `agents/openai.yaml` UI metadata is not mentioned anywhere — Codex deployments will be missing UI chips. | Codex portability (P2 cross-platform) | Low |
| G9 | The directory contract in SKILL.md output format lists `scripts/` and `assets/` as optional but the skill's own directory has neither — implicit message is "do as I say, not as I do." Self-conformance issue. | Invariant B (self-conformance) | Medium |
| G10 | No `mode-playbooks.md` for Modes 2–8. Without it, advertised modes are aspirational. | P5/P7 | High |
| G11 | The frontmatter description does not explicitly list the negative-trigger pattern for adjacent skills (e.g., do NOT use for AGENTS.md authoring). Could over-trigger. | P2 (trigger sharpness) | Low |
| G12 | The `references/` set includes structurally similar pairs that should cross-reference: `validation-rubric.md` ↔ `anti-patterns.md`, `progressive-disclosure.md` ↔ `workflow-patterns.md`. Missing pointers cause redundant agent reads. | P9 (efficient navigation) | Low |

### 4.3 Concrete delta list (for Phase 3)

Sorted by leverage. Implement top-down.

1. **[G1, G10] Add `references/mode-playbooks.md`** containing the workflow shape for Refactor, Review, Audit, Compress, Split, Merge, Adapt (see §3.2). In `SKILL.md` body, add a short "## Modes" section with one-paragraph entry-pointer per non-Create mode that links to the playbook.

2. **[G2] Create `scripts/` directory** with three files:
   - `init_skill.py` — generate boilerplate (folder + SKILL.md skeleton + chosen resource subdirs).
   - `quick_validate.py` — frontmatter + folder-structure check, exit non-zero on failure.
   - `check_links.py` — verify every `references/*` and `templates/*` pointer in SKILL.md resolves.

3. **[G3] Tighten Create Step 1** to require eliciting ≥3 concrete user-trigger examples before naming the skill. Update `SKILL.md` and add an example in `examples/create-from-scratch.md`.

4. **[G4] Bind the rubric to Step 5** by producing a filled scoring template as an artifact of every Create/Refactor/Audit run. Add `templates/audit-report-template.md` (matches the rubric's "Quick Scoring Template").

5. **[G7] Add `examples/create-from-scratch.md`** showing a full trajectory: user prompt → elicitation → mode detection → workflow execution → validation → final shipped skill.

6. **[G6] Add a short "Troubleshooting" section to SKILL.md** with the four iteration signals from §3.4 (under/over-trigger, bloat, drift). This is the iteration loop summary; full lifecycle goes in `references/lifecycle-and-iteration.md`.

7. **[G8] Document `agents/openai.yaml`** in `references/platform-compatibility.md` (Codex section). Add to Adapt-mode playbook.

8. **[G9] Self-conformance pass:** add `scripts/`, optionally `assets/`, to the actual `/skillify` directory. The skill should look like its own canonical output.

9. **[G5, G12] Consolidate cross-references:**
   - Remove the inline checklist from SKILL.md Step 5; replace with "Apply `references/validation-rubric.md`."
   - Add a one-line "see also" pointer at the top of each reference linking to its sibling pair.

10. **[G11] Sharpen description** with one explicit negative trigger covering the adjacent surface most likely to over-trigger. Candidate: `Do NOT use for editing AGENTS.md / CLAUDE.md repo-policy files (those are not skills).`

---

## 5. Hand-Off Contract to Phase 3 (GPT-5.5 xHigh)

Phase 3's job is to produce the final `SKILL.md` (and stub the new files identified above). This section gives Phase 3 explicit directives.

### 5.1 What to keep verbatim from the current /skillify
- Frontmatter description structure (refine wording, do not change shape).
- Six-step Create workflow (refinements per §4.3 deltas 3 and 4).
- Output-format directory tree (extend with `scripts/` per delta 8).
- Constraints section (add the one negative trigger from delta 10).
- Reference table (extend rows for the two new references from deltas 1 and 6).

### 5.2 What to refactor
- **Step 5 of Create**: replace inline checklist with a call to `scripts/quick_validate.py` (Gate A) followed by a pointer to `references/validation-rubric.md` (Gate B). Result: tighter SKILL.md, removed duplication.
- **Modes section**: lift the description's mode list into a body-level "## Modes" section with a 1–2 sentence pointer per mode. Modes 2–8 link to `references/mode-playbooks.md#<mode>`.
- **Iteration section**: replace current Step 6 prose with a four-bullet signal table; full runbook goes in `references/lifecycle-and-iteration.md`.

### 5.3 What to add
- New `references/mode-playbooks.md` (8 H2 sections, one per mode, ≤80 lines each).
- New `references/lifecycle-and-iteration.md` (signal → action runbook).
- New `templates/audit-report-template.md` (rubric scores + delta list shape).
- New `scripts/init_skill.py`, `scripts/quick_validate.py`, `scripts/check_links.py`.
- New `examples/create-from-scratch.md`.

### 5.4 Hard constraints on the final SKILL.md (Phase 3 must satisfy)
- **Length**: ≤500 lines, ideally ≤350.
- **Frontmatter description**: ≤1024 characters, includes ≥6 trigger phrases verbatim, includes ≥1 negative trigger, no XML.
- **Body structure (recommended order)**:
  1. Purpose (1 sentence)
  2. When to use this skill (3–5 lines)
  3. Modes (table + one-paragraph pointer per non-Create mode)
  4. Universal preamble (Step 0)
  5. Core workflow — Create (Steps 1–6, with refinements)
  6. Validation gate (4 gates, brief)
  7. Output format (directory tree)
  8. Constraints (DO NOT list)
  9. Troubleshooting (4 iteration signals)
  10. References table
  11. Templates list
  12. Scripts list
- **No** narrative prose outside the sections above.
- **Self-conformance**: every section must obey a §1 principle. If a section can't trace to a principle, cut it.

### 5.5 Phase-4 (Claude Opus 4.7 Max — review) checklist seed
After Phase 3 ships its draft, the Phase-4 reviewer should re-run §3.3 Gates A–D against the new SKILL.md and check:
- Did Phase 3 honor all §5.4 constraints?
- Are the new references (mode-playbooks, lifecycle-and-iteration) self-contained and one-level-deep?
- Does the new `scripts/quick_validate.py` actually catch the failures listed in Gate A?
- Are all delta items from §4.3 either implemented or explicitly deferred with rationale?

---

## 6. Open Questions / Deferred Decisions

These are flagged for the human reviewer. Each blocks a specific later choice.

- **Q1.** Should `/skillify` ship with the design document (this file) included as `references/architecture.md`, or omit it per P4 (no human-facing meta-docs)? *Recommendation:* omit. Re-introduce only if a future Refactor mode needs the rationale.
- **Q2.** Does the validation gate need a CI hook (e.g., a pre-commit script) or only an in-skill script? *Recommendation:* in-skill script only — CI is repo-policy, belongs in `AGENTS.md`, not in the portable skill.
- **Q3.** For Adapt mode, should we generate one combined skill with `compatibility: claude-code, codex, ...` or separate folders per platform? *Recommendation:* combined when frontmatter and body work unmodified across platforms; separate when paths or tool names differ structurally.
- **Q4.** Should `mode-playbooks.md` be one file or eight files? *Recommendation:* one file with H2-per-mode (P9 — references one level deep, easier navigation, single read for cross-mode refactors).
- **Q5.** The pipeline doc envisions Phase 5 (GPT-5.5 High) compressing the result. Is the §5.4 length cap (≤500 lines, ideally ≤350) the *post-compression* target or the Phase-3 target? *Recommendation:* Phase-3 target ≤500; Phase-5 target ≤350.

---

## Appendix A — Provenance Map

| Principle / Section | Primary source | Secondary sources |
|---------------------|----------------|-------------------|
| P1 Progressive disclosure | methodology §1 | openai-creator §"Progressive Disclosure"; skillify/references/progressive-disclosure.md |
| P2 Frontmatter triggers | methodology §2 | openai-creator §"Frontmatter"; skillify/references/frontmatter-guide.md |
| P3 Degrees of freedom | methodology §4 | openai-creator §"Set Appropriate Degrees of Freedom" |
| P4 No human-facing docs | methodology §3 | openai-creator §"What to Not Include"; skillify/references/anti-patterns.md #8 |
| P5 One responsibility | skillify/references/anti-patterns.md #1 | skillify/references/validation-rubric.md #2 |
| P6 Imperative numbered | methodology §5 | skillify/references/anti-patterns.md #5 |
| P7 Output contract | skillify/references/validation-rubric.md #4 | — |
| P8 One place per topic | skillify/references/anti-patterns.md #4, #11 | — |
| P9 References one level deep | methodology §1 | openai-creator §"Important guidelines" |
| P10 Validation as gate | methodology §5/§6 | openai-creator §"Step 5: Validate" |
| P11 Living document | methodology §5 step 5 | openai-creator §"Step 6: Iterate" |
| Mode taxonomy | skillify/SKILL.md description | This synthesis (no single prior source) |
| Validation gates A–D | This synthesis | methodology §6 pre-flight checklist; skillify/references/validation-rubric.md |
| Iteration signals | This synthesis | skillify/SKILL.md §Step 6; skillify/references/anti-patterns.md #10 |

## Appendix B — Action verb vocabulary

From `literature/skill-design-methodology/agent-reading-task-taxonomy.md` Compact Action Verb List, the verbs that must appear in skillify body workflows (in imperative form): **Read, Identify, Extract, Analyze, Classify, Validate, Synthesize, Design, Apply, Review, Improve, Document.** Other verbs from the taxonomy may appear in mode-playbooks.md per mode.

## Appendix C — What this design *does not* settle

- The exact wording of the final `description` (Phase 3 owns this; constraint is §5.4).
- The Python idiom for `init_skill.py` and `quick_validate.py` (Phase 3 / Phase 4 own; constraint is exit-code semantics in §3.3 Gate A).
- Visual / UX decisions for `agents/openai.yaml` chips (defer to user when Codex deployment is in scope).
