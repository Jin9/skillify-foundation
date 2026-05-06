# Mode Playbooks

Per-mode prompt skeletons, parallelism, and exit gates. Read this when running any mode beyond the universal preamble in `SKILL.md`.

Each playbook follows the same shape:
- **Trigger**: phrasing that selects the mode.
- **Input**: artifact paths the mode reads.
- **Spawn**: how many agents and which `subagent_type`.
- **Prompt skeleton**: the exact prompt template handed to spawned agents.
- **Output**: artifact written.
- **Exit gate**: condition to mark `phases.<mode>.status = done`.

---

## Mode: Plan

- **Trigger**: "decompose", "break this task into sub-questions", "plan the phases".
- **Input**: user task prompt + detected domain.
- **Spawn**: 1 × `subagent_type: general-purpose`. (Earlier versions used `subagent_type: Plan`; switched because `Plan` is documented as having no Write tool. See `agent-spawning.md`.)
- **Prompt skeleton**:
  ```
  You are decomposing a task into 1–8 sub-questions for a multi-agent pipeline.

  Task: {{user_prompt}}
  Domain: {{domain}}
  Pipeline state directory: {{task_dir}}

  Produce a markdown artifact at the path {{task_dir}}/01-plan.md following the
  template at /Users/n7lab/Desktop/Projects/Jin9/skillify-foundation/treasury/composing-agent-pipelines/templates/01-plan.md

  Constraints:
  - 1–8 sub-questions. If the task is atomic, output exactly 1.
  - Each sub-question must be answerable by a single search/read pass over the codebase or docs.
  - State explicit success criteria for the overall task.
  - List assumptions you are making.
  - Do not spawn other agents. Do not edit anything outside {{task_dir}}/01-plan.md.
  - Do not invoke composing-agent-pipelines.
  ```
- **Manifest**: before spawn, `update_manifest.py --task-id <id> --phase plan --status running --add-agent-call "general-purpose|Plan: decompose task"`. After exit gate, `--status done`.
- **Output**: `01-plan.md`.
- **Exit gate**: file exists; contains a numbered sub-question list (1–8 items); contains a "Success criteria" section; user types `proceed` (hard gate).

---

## Mode: Gather

- **Trigger**: "gather sources", "scout the codebase", "collect evidence".
- **Input**: `01-plan.md`.
- **Spawn**: N × `subagent_type: general-purpose` (parallel), where N = number of sub-questions in the plan. **All N calls in one message** so the runtime parallelizes them. (Earlier versions used `Explore`; switched because `Explore` has no Write tool, so artifacts were not produced reliably.)
- **Prompt skeleton** (per sub-question):
  ```
  You are gathering evidence for sub-question {{n}} of a multi-agent pipeline.

  Sub-question: {{sub_question_text}}
  Pipeline state directory: {{task_dir}}
  Plan: {{task_dir}}/01-plan.md

  Search the codebase / docs and produce {{task_dir}}/02-evidence/q{{n}}.md
  following the template at /Users/n7lab/Desktop/Projects/Jin9/skillify-foundation/treasury/composing-agent-pipelines/templates/02-evidence.md

  Constraints:
  - Cite every finding with file_path:line_number or doc URL.
  - If no evidence found, state explicitly "No evidence found" with what you searched for.
  - Read-only on source files. Only write to the q{{n}}.md path.
  - Do not invoke composing-agent-pipelines.
  ```
- **Manifest**: before spawn, `update_manifest.py --task-id <id> --phase gather --status running` plus one `--add-agent-call "general-purpose|Gather qN: ..."` per sub-question. After exit gate, `--status done`.
- **Output**: `02-evidence/q1.md`, `q2.md`, ...
- **Exit gate**: every sub-question has a corresponding `qN.md` file; each file contains either ≥1 cited source or "No evidence found"; soft gate (proceed unless user halts).

---

## Mode: Analyze

- **Trigger**: "synthesize", "find patterns", "what does the evidence say".
- **Input**: all `02-evidence/*.md`.
- **Spawn**: 1 × `subagent_type: general-purpose`.
- **Manifest**: before spawn, `update_manifest.py --task-id <id> --phase analyze --status running --add-agent-call "general-purpose|Analyze: synthesize evidence"`. After exit gate, `--status done`.
- **Prompt skeleton**:
  ```
  You are synthesizing evidence into findings for a multi-agent pipeline.

  Pipeline state directory: {{task_dir}}
  Evidence files: {{task_dir}}/02-evidence/q1.md ... q{{N}}.md
  Plan: {{task_dir}}/01-plan.md

  Read every evidence file. Produce {{task_dir}}/03-analysis.md
  following the template at /Users/n7lab/Desktop/Projects/Jin9/skillify-foundation/treasury/composing-agent-pipelines/templates/03-analysis.md

  Constraints:
  - Each finding must cite ≥1 evidence file path (e.g. 02-evidence/q3.md#section).
  - Surface contradictions and gaps explicitly; do not paper over them.
  - No editing source files. Do not run tools beyond reading the evidence files.
  - Do not invoke composing-agent-pipelines.
  ```
- **Output**: `03-analysis.md`.
- **Exit gate**: file exists; every finding has ≥1 source ref; soft gate.

---

## Mode: Review

- **Trigger**: "critique", "find gaps", "devil's advocate", "challenge the analysis".
- **Input**: `03-analysis.md`.
- **Spawn**: 1 × `subagent_type: general-purpose` (with adversarial framing). The agent is instructed not to edit the analysis — only to write `04-review.md`.
- **Manifest**: before spawn, `update_manifest.py --task-id <id> --phase review --status running --add-agent-call "general-purpose|Review: adversarial critique"`. After exit gate, `--status done`.
- **Prompt skeleton**:
  ```
  You are an adversarial reviewer. Your job is to find weaknesses in the
  analysis below. Do not be polite — be rigorous.

  Pipeline state directory: {{task_dir}}
  Analysis: {{task_dir}}/03-analysis.md

  Produce {{task_dir}}/04-review.md following the template at
  /Users/n7lab/Desktop/Projects/Jin9/skillify-foundation/treasury/composing-agent-pipelines/templates/04-review.md

  Categorize each issue as P1 (would change the conclusion), P2 (would change a
  finding), or P3 (cosmetic). For each P1/P2 issue, propose a concrete fix.

  Constraints:
  - Read-only. Do not edit the analysis.
  - Do not run tools beyond reading the analysis and any evidence files it cites.
  - Do not invoke composing-agent-pipelines.
  ```
- **Output**: `04-review.md`.
- **Exit gate**: file exists; issue list (may be empty); **hard cap: this mode runs at most once per pipeline**; soft gate (proceed unless user halts).

---

## Mode: Validate

- **Trigger**: "fact-check", "verify claims", "validate against source".
- **Input**: `03-analysis.md` + `04-review.md` (Review optional).
- **Spawn**: 1 × `subagent_type: general-purpose`. The single agent fact-checks all P1/P2 claims sequentially and writes a single `05-validation.md`. (Earlier versions specified N parallel `Explore` agents; collapsed to 1 agent for two reasons: `Explore` has no Write tool, and parallel writes to the same file would race.)
- **Prompt skeleton**:
  ```
  You are fact-checking claims from a multi-agent analysis.

  Pipeline state directory: {{task_dir}}
  Inputs:
  - Analysis: {{task_dir}}/03-analysis.md
  - Review (if exists): {{task_dir}}/04-review.md

  For each P1/P2 claim in the analysis (and each P1/P2 issue from the review),
  produce a verdict in {{task_dir}}/05-validation.md following the template at
  /Users/n7lab/Desktop/Projects/Jin9/skillify-foundation/treasury/composing-agent-pipelines/templates/05-validation.md

  Verdict must be one of:
  - confirmed: source supports the claim verbatim or by clear implication.
  - refuted: source contradicts the claim.
  - unverifiable: source is missing, ambiguous, or doesn't address the claim.

  Constraints:
  - Read-only on inputs and source files; only write 05-validation.md.
  - Quote the relevant lines from the source. Do not paraphrase.
  - Mark unverifiable when source is missing or ambiguous, not when you simply disagree.
  - Do not invoke composing-agent-pipelines.
  ```
- **Manifest**: before spawn, `update_manifest.py --task-id <id> --phase validate --status running --add-agent-call "general-purpose|Validate: fact-check claims"`. After exit gate, `--status done`. Add a `--note` summarizing verdict counts (e.g. `"11 confirmed, 0 refuted, 0 unverifiable"`).
- **Output**: `05-validation.md` (one section per claim).
- **Exit gate**: every P1/P2 claim has a verdict; **hard cap: this mode runs at most once per pipeline**; if any P1 verdict is `refuted`, surface to user and pause before Decide.

---

## Mode: Decide

- **Trigger**: "recommend", "decide between", "which option", "choose".
- **Input**: `03-analysis.md` (+ `04-review.md` and `05-validation.md` if they exist).
- **Spawn**: 1 × `subagent_type: general-purpose`.
- **Manifest**: before spawn, `update_manifest.py --task-id <id> --phase decide --status running --add-agent-call "general-purpose|Decide: produce recommendations"`. After exit gate, `--status done`.
- **Prompt skeleton**:
  ```
  You are producing a recommendation from a multi-agent pipeline.

  Pipeline state directory: {{task_dir}}
  Inputs:
  - Analysis: {{task_dir}}/03-analysis.md
  - Review (if exists): {{task_dir}}/04-review.md
  - Validation (if exists): {{task_dir}}/05-validation.md

  Produce {{task_dir}}/06-decision.md following the template at
  /Users/n7lab/Desktop/Projects/Jin9/skillify-foundation/treasury/composing-agent-pipelines/templates/06-decision.md

  Required sections:
  - Options table (≥2 options with pros/cons/risks)
  - Recommendation (one option)
  - Residual risks (what could still go wrong)
  - Trigger conditions (when to revisit this decision)

  Constraints:
  - Recommendation must be supported by analysis findings; cite them.
  - If validation refuted a P1 claim, do not base the recommendation on that claim.
  - Read-only on inputs; only write 06-decision.md.
  - Do not invoke composing-agent-pipelines.
  ```
- **Output**: `06-decision.md`.
- **Exit gate**: file exists; options table has ≥2 rows; recommendation cited; user types `proceed` (hard gate).

---

## Mode: Compact

- **Trigger**: "produce final report", "write the deliverable", "compact this".
- **Input**: all artifacts.
- **Spawn**: none. The main agent (you) does this inline.
- **Manifest**: before starting, `update_manifest.py --task-id <id> --phase compact --status running`. After both files exist, `--status done`. (No `--add-agent-call` because nothing is spawned.)
- **Procedure**:
  1. Read `01-plan.md`, all `02-evidence/*.md`, `03-analysis.md`, and any of `04-review.md`, `05-validation.md`, `06-decision.md` that exist.
  2. Write `07-final.md` following the template at `templates/07-final.md`. Required sections: TL;DR, Sub-questions & answers, Findings, Recommendation (if Decide ran), Validation summary (if Validate ran), Open questions, Sources.
  3. Write `summary.md` — a 100–200 word TL;DR with the most important citations.
- **Exit gate**: both files exist; word count of `07-final.md` ≤ 3000 (or surface the overage to the user); every claim cites a source.

---

## Mode: Compose

- **Trigger**: "run the pipeline", "compose phases", "end-to-end" — or no mode hint at all.
- **Input**: user task + domain.
- **Spawn**: dispatches per the chosen domain shape from `domain-shapes.md`.
- **Procedure**:
  1. Run the universal preamble in `SKILL.md`.
  2. Print the chosen phase shape and wait for `proceed` or override.
  3. Run each phase in shape order. After each phase, update the manifest, print a one-line summary, and check the gate.
  4. Pause at hard gates (Plan, Decide). Resume on user `proceed`.
  5. Halt if a phase fails its gate after one retry. Do not auto-recover.
- **Output**: full pipeline directory.
- **Exit gate**: every phase in `phase_shape` has status `done` or `skipped`; final and summary files exist.
