# Adjudication and provenance

External CLI output is **advisory**. The orchestrating agent is the single decision-maker: it
verifies, adjudicates, and owns the committed result. This file defines how.

## Why advisory-only

- **Errors compound.** Across a multi-step chain, per-step reliability multiplies: ten 85%-reliable
  steps succeed end-to-end only about 20% of the time (0.85^10 ≈ 0.197). Unverified hand-offs
  amplify small errors into large ones, so adjudication at each boundary is the brake.
- **Separation of concerns.** Sub-agents advise; only the primary controller commits state. Keep
  advisory roles strictly separate from authoritative tool-execution.

## How to adjudicate

1. **Read the provenance first.** Confirm the call actually ran on the requested model (backend
   label from the log) and did not hit the watchdog. A degraded or fallback model invalidates the
   advice before you even read it.
2. **Verify claims against primary sources.** For facts and citations, check them directly — read
   the cited files, re-run the check, or compare against authoritative docs. Do not promote a claim
   because it is plausible or confidently stated. Treat bare-domain or unopenable citations as
   leads, not evidence.
3. **Resolve conflicts.** When advice conflicts with the agent's own analysis or another CLI's
   output, prefer what is verifiable. Score against an explicit rubric (correctness, evidence,
   reproducibility) rather than prose vibes — an "LLM-as-judge" pass with stated criteria beats a
   gut call.
4. **Calibrate autonomy to reversibility.** Auto-accept low-stakes, easily reversible outputs. For
   high-blast-radius or irreversible actions, require human review before committing.
5. **Decide and mark it.** Emit the result as the agent's decision, not the CLI's, with an explicit
   accept / modify / reject verdict and the reason.

## What to record (provenance)

Write one consult-record per CLI call from `templates/consult-record.md` — one file per call, never
concatenated:

- Timestamp, mode, CLI, requested model label, **backend label verified from the log**, effort.
- The prompt (or its path) and the advisory output (snippet or path).
- Watchdog outcome and exit code.
- The agent's verdict (accept / modify / reject) and the reasoning.

Recording exists for replay and audit: capture enough to answer *why* a decision was made and on
*what* model output. Redact secrets and PII before persisting verbose output. For retries of
side-effecting work, bind an idempotency key to the action so a watchdog-triggered retry does not
duplicate writes.
