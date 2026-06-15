# Fan-out and reconciliation

Fan-out dispatches the same question to two or more CLIs in parallel and reconciles their answers.
It costs more, so use it deliberately.

## When to fan out

- The question is **divergence-prone** (judgment, design trade-offs, ambiguous specs) and one
  model's answer is not trustworthy alone.
- The stakes are **high or hard to reverse**, and independent corroboration lowers risk.
- You want to **neutralize a single model's bias or hallucination** on a verifiable fact.

Do not fan out for cheap, low-stakes, or easily-verified-directly tasks — a single consult plus the
agent's own check is enough.

## How to dispatch

- Run each CLI through its wrapper. They are independent, so dispatch them concurrently and collect
  all outputs (and provenance lines) before reconciling.
- Use a diverse panel when it helps: different backends (GPT via codex, Gemini via agy) surface
  different failure modes. Record one consult-record per CLI.

## How to reconcile

- **Majority vote** for verifiable, factual extraction — take the consensus, flag dissent.
- **Reasoning-path validation** for judgment tasks — evaluate *why* each answer holds, not just the
  final string; prefer the soundest chain, even if it is the minority view.
- **Meta-synthesis** for open-ended work — the agent composes a single answer from the strongest
  parts of each, attributing which input contributed what.
- On genuine deadlock, widen the evidence (read primary sources directly) or escalate to the user;
  do not average two guesses into a third guess.

The reconciled answer is still the agent's decision, recorded with the per-CLI verdicts that fed it.
See `adjudication-and-provenance.md`.
