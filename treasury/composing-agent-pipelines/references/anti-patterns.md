# Pipeline-Specific Anti-Patterns

These are failure modes specific to this skill. For general skill anti-patterns,
use the active `skillify` skill's `references/anti-patterns.md`.

## 1. Recursive Compose

Compose mode invoking Compose, or any phase invoking Compose, creates an unbounded loop and torches the context window.

**Symptom**: A delegated worker reads `SKILL.md` and decides to start another pipeline as part of its phase work.

**Fix**: Phase prompts must include the literal instruction "Do not invoke
`composing-agent-pipelines`." The composer's validation gate refuses recursive
Compose when a parent manifest or caller metadata shows this skill is already
inside another pipeline.

## 2. Infinite Review/Validate loops

If Review surfaces P1 issues, the natural temptation is to re-run Analyze → Review until none remain. This loop has no fixed point.

**Fix**: Hard cap of one Review pass and one Validate pass per pipeline run. After one pass, surface remaining issues to the user and let them decide:
- Accept the issues → proceed to Decide / Compact.
- Re-run a specific phase manually with the corrected context (counts as a new pipeline run).
- Halt the pipeline.

## 3. Over-decomposition on trivial tasks

Plan happily produces 8 sub-questions for a question that needed 1. Each sub-question fans out to a worker — wasted parallelism and context.

**Fix**: When the user prompt is under ~50 words and clearly single-domain, the composer should suggest skipping straight to the relevant phase mode (e.g., "this looks like a single-question gather; want to invoke Gather standalone instead?"). Plan's prompt also caps at 8 sub-questions and explicitly asks for 1 sub-question if the task is atomic.

## 4. Re-reading the entire pipeline directory each phase

Each phase only needs specific upstream artifacts. Reading the whole directory blows context.

**Fix**: Phase prompts include exact file paths the worker should read. The worker is instructed to NOT use `ls` or `find` on the pipeline directory.

## 5. Background Worker Zombies

Long-running Gather workers launched in background mode can be forgotten if the orchestrator crashes or the user halts.

**Fix**: Background is only allowed when the user explicitly asks. The manifest
records the background id under `phases.gather.delegations[].background_id`.
Before declaring Gather done, the composer must confirm all background workers
reported completion. If the user halts, stop outstanding background work using
the host's cancellation mechanism when available.

## 6. Worktree leakage

`isolation: worktree` creates a temporary git worktree. If the worker does not exit cleanly, the worktree persists.

**Fix**: Only Gather and Validate may use worktrees, and only on explicit user request. Each phase must `ExitWorktree` before writing its artifact. The composer's gate verifies the worktree is closed by checking for residual paths in the manifest.

## 7. Manifest drift

Multiple workers writing to `manifest.json` race; the orchestrator alone owns it.

**Fix**: Delegated workers are forbidden to write the manifest. They write only to the artifact path the orchestrator specified. The orchestrator updates the manifest after each phase completes.

## 8. Stale state collision

Reusing a `<task-id>` for a fresh task overwrites previous artifacts.

**Fix**: `init_pipeline.py` refuses to write to an existing directory unless `--resume` is passed. Resume requires a valid manifest and resumes from the next pending phase.

## 9. Citation-free claims

Phases that produce findings without citations have nothing to validate. The pipeline's epistemic integrity depends on traceability.

**Fix**: Each artifact template includes a "Sources" section. The validation gate at every phase rejects findings/claims/recommendations that lack a source reference. Empty results are acceptable (`"no evidence found"`) but unsourced assertions are not.

## 10. Compose without an announced shape

If the composer launches into Plan without printing the chosen domain shape, the user cannot intervene before context-expensive phases run.

**Fix**: Compose's universal preamble step 5 mandates announcing the phase shape and waiting for `proceed` or override. No exceptions.

## 11. PII or secrets in the manifest

`manifest.json` records the user prompt verbatim. If the prompt contains a token, key, or PII, it leaks to anyone reading the file.

**Fix**: `init_pipeline.py` scrubs candidate secrets from the recorded prompt: any token-like substring (24+ chars of base64 / hex), email addresses, and `password=`/`token=` patterns are replaced with `[redacted]` before persistence. Document this scrubbing rule in `state-convention.md`.
