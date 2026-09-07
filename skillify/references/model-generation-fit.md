# Fit for a Model Generation

Skills are per-generation artifacts: a line that is load-bearing on one model generation is cruft on the next, and a line the next generation needs may be missing from a skill that worked before. Encode behavior by capability tier, never by model name, and re-audit at every release (`lifecycle-and-iteration.md`, Model release row). This file holds the calibration principles; the reusable block they produce is `templates/operating-contract.md`, and the greppable signals are in `scripts/cruft_scan.py`.

See also: `anti-patterns.md` items 5, 6, and 12, `workflow-patterns.md` Cross-Cutting Techniques, and the Re-baseline sub-flow in `mode-playbooks.md`.

## 1. Specificity matched to fragility

Prescriptive text is for operations where exactly one sequence is safe. Everywhere else, state the goal, the one or two real constraints with their reasons, and what done looks like, and let the agent plan the method: current models plan natively, and a hand-written script for judgment work lowers output quality.

- Judgment work (review, analysis, design, writing): goal, reasoned constraints, verification, stop condition, in prose. Number steps only where a later step consumes an earlier step's output.
- Fragile or irreversible work (destructive commands, auth flows, compliance steps, deployments): numbered steps with entry and exit conditions, exact commands, and `scripts/` for the parts that must not vary.
- The deletion test for any sentence: if removing it would change neither what is legal nor how success is measured, it is strategy coaching, and the agent's own plan is usually better. Delete it.
- When an eval or a set of collected trigger prompts exists, A/B the skill with its step scaffolding removed before deciding it was needed. Removal is a hypothesis, not a conclusion.

## 2. Instruction priority and stops

Current models weight instructions found in skills and repo-policy files heavily. When a skill line conflicts with what the user asked, the agent may pause, ask for a confirmation nobody wanted, leave requested work unfinished, or follow the skill's rule instead of the request. Every multi-step skill therefore states the tie-break: the user's request in the session wins over the skill; repo-policy files outrank the skill's defaults; and if a skill line is set aside for that reason, the agent says which line and quotes it, so the author can fix it.

Stops are limited to two cases: a destructive or irreversible action, and a genuine change to the scope the user set. "Can you", "help me", and "please" are instructions, not invitations to ask whether to proceed. A skill asks at most one question per run, and only when a required input is missing and cannot be taken from the request, pasted text, or the workspace. When the user asked a question rather than for a change, the assessment is the deliverable. Before ending a turn, the agent checks its last paragraph: a plan or a promise is work not yet done.

## 3. Verification scope

Make self-verification explicit and name the evidence: an exit code, a validator's PASS line, a re-read of the file that was written. A claim of progress is reported only when a tool result from the session backs it, and a check that was not run is not described. Scale checks to reversibility: no new tests for reversible, low-impact edits; broader testing only when a failure justifies it; scratch checks need not be kept. A fresh-context verifier sub-agent outperforms self-critique for the final check. Verification text is never removed during a re-baseline, whatever else is loosened.

## 4. Delegation policy

Say when delegation is desirable, not only that it is allowed. Independent sub-tasks with no shared state are parallel sub-agent work; the final check is fresh-context sub-agent work; sequential or judgment-heavy work stays in the main thread. Batch independent tool calls in one turn. Prefer asynchronous fan-out with a bounded wait over spawn-and-block. Messages to sub-agents stand alone: goal, inputs, expected output shape, stop condition. Sub-agent and external-model output is advisory; the lead reconciles it and owns the result.

## 5. Effort tiers instead of model names

Two levers, kept separate: the model-cost tier of a step (`small` for mechanical extraction and formatting, `mid` for structuring and routine passes, `frontier` for ambiguous design and hard judgment) and the effort hint (`low`, `medium`, `high`). Effort is the primary quality-for-cost control on current models, and a frontier tier at low effort often beats a smaller tier at high effort on judgment steps; raise effort only for hard, verifiable steps. Tiers and hints live in the skill body, as a `model-cost tier:` sentence or heading suffix or as an inline `[tier/effort]` tag on the steps that differ. They stay out of portable frontmatter (the host-specific `model:` and `effort:` fields are documented in `frontmatter-guide.md` and are not portable), and a model name never appears in a rule: names rot within months, tiers do not. A dated "e.g." beside a tier in a data table is acceptable; a rule that says "use model X" is not.

## 6. Reporting and style

Current models narrate less between tool calls and format less in chat than earlier generations did. Ask for what you want instead of suppressing what you fear: an opening line saying what will be done and which files will be touched, and a closing recap that stands on its own (what changed, what was verified, what remains). Remove narration suppressors ("hold all findings for the final response", "don't narrate") and anti-formatting rules ("never use bullets"); they now strip output the reader wanted. Prose carries behavior and its reasons; structure carries reference data. Each real constraint travels with its "because", and emphasis is a scoped fix for one demonstrably under-weighted instruction, not a register. Prefer literal phrasing to stock connective phrases; a banned-word list is itself a dated pattern, so state the desired style positively in one line.

## 7. Memory and on-the-fly updates

Long-running skills name where the agent records learnings (a notes file in the working directory, never inside the skill folder) and tell it to consult that file in later sessions: one lesson per entry, corrections and confirmed approaches alike, with why they mattered. Current models are good at improving a skill from what they learn mid-task; let them propose the edit, and route it through Refactor's preservation rules rather than editing the installed copy silently.

## 8. Keep list

An audit that only says "delete" hurts the users who follow it most diligently. These stay even when a scanner flags them:

1. Context is never cruft: audience, product, environment facts, quality bar, constraints and their reasons. Too-short skills produce generic output because the agent fills gaps with safe defaults.
2. Cruft is not length. Never justify a deletion by line count alone.
3. Fragile operations keep exact scripts and "do not modify this command" wording.
4. Tool and script contract detail stays and often grows: parameters, limits, failure modes, what is not returned.
5. Prohibitions against current, demonstrated failures stay, with their reason beside them.
6. Trigger and routing text (the frontmatter `description`, a "do not use for" line) may carry calibrated urgency, because skills under-trigger; body text explains rather than shouts.
7. Format-pinning examples on genuinely format-sensitive outputs stay, labelled illustrative.
8. Working redundancy is not cruft; consolidate only when duplicates disagree. One deliberate recap of the few key constraints is a reasonable pattern; scattered repetition is not.
9. Re-baselining adds text too: matching a skill to a new generation sometimes means adding guidance for that generation's failure modes (a stop condition, a progress line, a verification step).

## 9. Current-generation profile (verified 2026-09-07)

Vendor-neutral traits of the frontier generation this file was calibrated against. Re-verify against the vendors' own prompting guides at each release; the sources are recorded in the methodology research notes for the 2026-09 corpus cohort.

- Follows instructions closely and literally, including instructions in skills and repo-policy files; pauses or diverges when those sources disagree with the request.
- Plans natively; scripted step choreography and "think step by step" scaffolds lower quality rather than raise it.
- Narrates less between tool calls and formats less in chat; narration suppressors and anti-formatting rules are now harmful.
- Runs parallel sub-agents dependably; asynchronous delegation outperforms spawn-and-wait.
- Treats effort as the main quality-for-cost lever; low effort on a frontier tier is often competitive with smaller tiers at high effort.
- May refuse a request that asks it to reproduce its own reasoning; ask for the verified result and its evidence instead.
- Tends to scope-creep at high effort (unrequested refactors, extra tests, whole-file rewrites) unless the skill states the scope discipline.
- Performs better with a memory surface and with the reason behind a request stated up front.
