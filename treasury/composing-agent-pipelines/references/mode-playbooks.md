# Mode Playbooks

Per-mode prompt skeletons, delegation shape, and exit gates. Read this before
running any mode beyond the universal preamble in `SKILL.md`.

Each playbook follows the same shape:
- **Trigger**: phrasing that selects the mode.
- **Input**: artifact paths the mode reads.
- **Delegation**: whether to use a host-native worker or run inline.
- **Prompt skeleton**: the prompt template handed to the worker or used inline.
- **Output**: artifact written.
- **Exit gate**: condition to mark `phases.<mode>.status = done`.

Resolve template paths relative to the skill folder. Do not hardcode absolute
local paths in worker prompts.

---

## Mode: Plan

- **Trigger**: "decompose", "break this task into sub-questions", "plan the phases".
- **Input**: user task prompt + detected domain.
- **Delegation**: 1 writer worker, or inline when delegation is unavailable.
- **Prompt skeleton**:
  ```text
  You are decomposing a task into 1-8 sub-questions for a multi-agent pipeline.

  Task: {{user_prompt}}
  Domain: {{domain}}
  Pipeline state directory: {{task_dir}}
  Template: templates/01-plan.md

  Produce {{task_dir}}/01-plan.md.

  Constraints:
  - 1-8 sub-questions. If the task is atomic, output exactly 1.
  - Each sub-question must be answerable by a single search/read pass over the codebase or docs.
  - State explicit success criteria for the overall task.
  - List assumptions you are making.
  - Do not delegate to other workers. Do not edit anything outside {{task_dir}}/01-plan.md.
  - Do not invoke composing-agent-pipelines.
  ```
- **Manifest**: before work, `update_manifest.py --task-id <id> --phase plan --status running --add-delegation "writer|Plan: decompose task"` (or `inline|Plan: decompose task`). After exit gate, `--status done`.
- **Output**: `01-plan.md`.
- **Exit gate**: file exists; contains a numbered sub-question list (1-8 items); contains a "Success criteria" section; user types `proceed` (hard gate).

---

## Mode: Gather

- **Trigger**: "gather sources", "scout the codebase", "collect evidence".
- **Input**: `01-plan.md`.
- **Delegation**: N writer workers in parallel when available, where N is the number of sub-questions. Sequential or inline fallback is allowed, but still write one file per sub-question.
- **Prompt skeleton** (per sub-question):
  ```text
  You are gathering evidence for sub-question {{n}} of a multi-agent pipeline.

  Sub-question: {{sub_question_text}}
  Pipeline state directory: {{task_dir}}
  Plan: {{task_dir}}/01-plan.md
  Template: templates/02-evidence.md

  Search the codebase / docs and produce {{task_dir}}/02-evidence/q{{n}}.md.

  Constraints:
  - Cite every finding with file_path:line_number or doc URL.
  - If no evidence found, state exactly "No evidence found" with what you searched.
  - Read-only on source files. Only write to the q{{n}}.md path.
  - Do not invoke composing-agent-pipelines.
  ```
- **Manifest**: before work, `update_manifest.py --task-id <id> --phase gather --status running` plus one `--add-delegation "writer|Gather qN: ..."` per sub-question. After exit gate, `--status done`.
- **Output**: `02-evidence/q1.md`, `q2.md`, ...
- **Exit gate**: every sub-question has a corresponding `qN.md`; each file contains at least one cited source or "No evidence found"; soft gate.

---

## Mode: Analyze

- **Trigger**: "synthesize", "find patterns", "what does the evidence say".
- **Input**: all `02-evidence/*.md`.
- **Delegation**: 1 writer worker, or inline.
- **Manifest**: before work, `update_manifest.py --task-id <id> --phase analyze --status running --add-delegation "writer|Analyze: synthesize evidence"`. After exit gate, `--status done`.
- **Prompt skeleton**:
  ```text
  You are synthesizing evidence into findings for a multi-agent pipeline.

  Pipeline state directory: {{task_dir}}
  Evidence files: {{task_dir}}/02-evidence/q1.md ... q{{N}}.md
  Plan: {{task_dir}}/01-plan.md
  Template: templates/03-analysis.md

  Read every evidence file. Produce {{task_dir}}/03-analysis.md.

  Constraints:
  - Each finding must cite at least one evidence file path.
  - Surface contradictions and gaps explicitly; do not paper over them.
  - No editing source files. Do not run tools beyond reading the evidence files.
  - Do not invoke composing-agent-pipelines.
  ```
- **Output**: `03-analysis.md`.
- **Exit gate**: file exists; every finding has at least one source reference; soft gate.

---

## Mode: Review

- **Trigger**: "critique", "find gaps", "devil's advocate", "challenge the analysis".
- **Input**: `03-analysis.md`.
- **Delegation**: 1 writer worker, or inline.
- **Manifest**: before work, `update_manifest.py --task-id <id> --phase review --status running --add-delegation "writer|Review: adversarial critique"`. After exit gate, `--status done`.
- **Prompt skeleton**:
  ```text
  You are an adversarial reviewer. Your job is to find weaknesses in the analysis.

  Pipeline state directory: {{task_dir}}
  Analysis: {{task_dir}}/03-analysis.md
  Template: templates/04-review.md

  Produce {{task_dir}}/04-review.md.

  Categorize each issue as P1 (would change the conclusion), P2 (would change a
  finding), or P3 (cosmetic). For each P1/P2 issue, propose a concrete fix.

  Constraints:
  - Read-only. Do not edit the analysis.
  - Do not run tools beyond reading the analysis and cited evidence files.
  - Do not invoke composing-agent-pipelines.
  ```
- **Output**: `04-review.md`.
- **Exit gate**: file exists; issue list may be empty; hard cap of one Review pass per pipeline; soft gate.

---

## Mode: Validate

- **Trigger**: "fact-check", "verify claims", "validate against source".
- **Input**: `03-analysis.md` + `04-review.md` when present.
- **Delegation**: 1 writer worker, or inline. Do not parallelize; one artifact receives all verdicts.
- **Prompt skeleton**:
  ```text
  You are fact-checking claims from a multi-agent analysis.

  Pipeline state directory: {{task_dir}}
  Inputs:
  - Analysis: {{task_dir}}/03-analysis.md
  - Review (if exists): {{task_dir}}/04-review.md
  Template: templates/05-validation.md

  For each P1/P2 claim in the analysis and each P1/P2 issue from the review,
  produce a verdict in {{task_dir}}/05-validation.md.

  Verdict must be one of:
  - confirmed: source supports the claim verbatim or by clear implication.
  - refuted: source contradicts the claim.
  - unverifiable: source is missing, ambiguous, or does not address the claim.

  Constraints:
  - Read-only on inputs and source files; only write 05-validation.md.
  - Quote the relevant lines from the source. Do not paraphrase.
  - Mark unverifiable when source is missing or ambiguous.
  - Do not invoke composing-agent-pipelines.
  ```
- **Manifest**: before work, `update_manifest.py --task-id <id> --phase validate --status running --add-delegation "writer|Validate: fact-check claims"`. After exit gate, `--status done`. Add a `--note` summarizing verdict counts.
- **Output**: `05-validation.md`.
- **Exit gate**: every P1/P2 claim has a verdict; hard cap of one Validate pass per pipeline; if any P1 verdict is `refuted`, surface to user and pause before Decide.

---

## Mode: Decide

- **Trigger**: "recommend", "decide between", "which option", "choose".
- **Input**: `03-analysis.md` plus `04-review.md` and `05-validation.md` when present.
- **Delegation**: 1 writer worker, or inline.
- **Manifest**: before work, `update_manifest.py --task-id <id> --phase decide --status running --add-delegation "writer|Decide: produce recommendations"`. After exit gate, `--status done`.
- **Prompt skeleton**:
  ```text
  You are producing a recommendation from a multi-agent pipeline.

  Pipeline state directory: {{task_dir}}
  Inputs:
  - Analysis: {{task_dir}}/03-analysis.md
  - Review (if exists): {{task_dir}}/04-review.md
  - Validation (if exists): {{task_dir}}/05-validation.md
  Template: templates/06-decision.md

  Produce {{task_dir}}/06-decision.md.

  Required sections:
  - Options table with at least two options and their pros/cons/risks.
  - Recommendation naming one option.
  - Residual risks.
  - Trigger conditions for revisiting the decision.

  Constraints:
  - Recommendation must be supported by analysis findings; cite them.
  - If validation refuted a P1 claim, do not base the recommendation on that claim.
  - Read-only on inputs; only write 06-decision.md.
  - Do not invoke composing-agent-pipelines.
  ```
- **Output**: `06-decision.md`.
- **Exit gate**: file exists; options table has at least two rows; recommendation is cited; user types `proceed` (hard gate).

---

## Mode: Compact

- **Trigger**: "produce final report", "write the deliverable", "compact this".
- **Input**: all artifacts.
- **Delegation**: inline only.
- **Manifest**: before starting, `update_manifest.py --task-id <id> --phase compact --status running`. After both files exist, `--status done`.
- **Procedure**:
  1. Read `01-plan.md`, all `02-evidence/*.md`, `03-analysis.md`, and any of `04-review.md`, `05-validation.md`, `06-decision.md` that exist.
  2. Write `07-final.md` following `templates/07-final.md`.
  3. Write `summary.md`, a 100-200 word summary with the most important citations.
- **Exit gate**: both files exist; `07-final.md` is no more than 3000 words, or the overage is surfaced; every claim cites a source.

---

## Mode: Compose

- **Trigger**: "run the pipeline", "compose phases", "end-to-end" — or no mode hint.
- **Input**: user task + domain.
- **Delegation**: dispatches per the chosen domain shape and phase playbooks.
- **Procedure**:
  1. Run the universal preamble in `SKILL.md`.
  2. Print the chosen phase shape and wait for `proceed` or override.
  3. Run each phase in shape order. After each phase, update the manifest, print a one-line summary, and check the gate.
  4. Pause at hard gates (Plan, Decide). Resume only on user `proceed`.
  5. Halt if a phase fails its gate after one retry.
- **Output**: full pipeline directory.
- **Exit gate**: every phase in `phase_shape` has status `done` or `skipped`; final and summary files exist.
