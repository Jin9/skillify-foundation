---
name: delegating-to-cli-models
description: >
  Drive external model CLIs as advisory sub-agents from one orchestrating agent: codex
  (GPT, gpt-5.6-sol) and agy (Antigravity/Gemini), plus optional headless self-delegation. Frame a
  scoped prompt, dispatch it headlessly under a watchdog, ground-truth the model that actually
  answered, then adjudicate the advisory output and record provenance. Use when the user asks
  to "ask codex/gpt", "get gemini's take via agy", "run gemini pro 3.1 high on this", "have agy
  digest the ResearchVault", "get a second opinion from another model", or "fan this out to codex
  and gemini". Modes: consult, research-digest, code-structure, fan-out. The orchestrating agent
  stays the decision-maker; CLI output is advisory only. Do NOT use for host-native agent
  pipelines (use composing-agent-pipelines), supervision and gate design (use
  agentic-workflow-design), the agent-scaffold platform (use orchestrating-agent-scaffold), or
  read-only vault Q&A (use research-vault-librarian).
---

# Delegating to CLI Models

## Purpose

Dispatch a scoped sub-task to an external command-line model — codex (GPT), agy
(Antigravity/Gemini), or an optional headless instance of the host agent — collect its output as
**advisory** input, then adjudicate it, verify it against primary sources, and record the verdict
with provenance. The orchestrating agent stays the single decision-maker.

## When to use this skill

- Use when the user asks to "ask codex/gpt this", "get gemini's take via agy", "run gemini pro 3.1
  high on X", "have agy digest the ResearchVault", "get a second opinion from another model", or
  "fan this out to codex and gemini and reconcile".
- Use when a task benefits from another model's grounding (vault + websearch), a structuring pass,
  or independent verification before the agent commits.
- Do NOT use for host-native sub-agent pipelines that need no external CLI — use
  `composing-agent-pipelines`.
- Do NOT use to design supervision, approval gates, or never-do policy — use
  `agentic-workflow-design`.
- Do NOT use to drive the just/state.json agent-scaffold platform — use
  `orchestrating-agent-scaffold`.
- Do NOT use for read-only questions about an existing research vault — use
  `research-vault-librarian`.

## Modes

| Mode | Use when | Dispatch | Output |
|------|----------|----------|--------|
| Consult | one-shot second opinion / judgment | one CLI | adjudicated answer + record |
| Research-digest | grounding from the vault + web | agy, frontier tier | adjudicated synthesis + record |
| Code-structure | hand a drafting/structuring/refactor pass to GPT | codex | reviewed artifact (advisory) + record |
| Fan-out | divergence-prone or high-stakes question | two or more CLIs in parallel | reconciled answer + per-CLI records |

## Universal preamble (run before every mode)

1. **Frame the job.** Write a self-contained prompt: state the task, the inputs (as absolute paths
   or pasted text), and the exact output shape. Pick the CLI by capability and cost tier (see
   `references/effort-and-cost-tiers.md`).
2. **Pre-flight the CLI.** Confirm it is on PATH, and choose the model by **tier**, not from memory.
   Do not gate the dispatch on a model-list lookup: for agy that call is server-backed and can hang
   indefinitely, and an unusable model identifier surfaces at dispatch anyway.
3. **Arm the watchdog.** macOS has no `timeout`; every dispatch goes through a wrapper in
   `scripts/` that kills a hung call. The CLI's own response timeout must sit *under* the watchdog,
   or the watchdog can never fire. Pick an effort tier appropriate to the task.

## Workflow

1. **Frame & route** per the preamble; choose the mode from the table.
2. **Dispatch** through the matching wrapper — never call the raw CLI inline:
   - `scripts/run_agy.sh "<MODEL>" "<prompt>" [neutral_dir]`
   - `scripts/run_codex.sh "<prompt>" [model]`
   - `scripts/run_headless_claude.sh "<prompt>" [model]`

   Each enforces the per-CLI contract (agy: per-run log, response timeout under the watchdog; codex:
   no model override unless you pin one; headless host agent: `--effort low`) plus a watchdog. See
   `references/cli-invocation-contracts.md`.
3. **Ground-truth the model.** Read the `PROVENANCE:` line the wrapper prints and require
   `backend_verified=yes` (agy) or `model_verified=yes` (codex). A CLI can accept a model
   identifier, answer normally, exit 0 — and have run a *different* model. A zero exit is not
   evidence, and the model's self-reported identity is worthless.
4. **Adjudicate (advisory only).** Treat the output as advice, not truth. Verify factual and
   citation claims against primary sources the agent can read directly; resolve conflicts; decide.
   For fan-out, reconcile per `references/fan-out-and-reconciliation.md`. See
   `references/adjudication-and-provenance.md`.
5. **Record.** Write one consult-record per CLI call from `templates/consult-record.md` (mode, CLI,
   requested model, verified model, prompt, advisory snippet, watchdog outcome, the agent's verdict,
   provenance).
6. **Emit** the adjudicated answer/artifact plus the record path(s).

## Output contract

- The adjudicated answer or artifact, clearly marked as the agent's decision (not the CLI's).
- One consult-record per CLI call — one file per call, do not concatenate.
- A provenance line per call: CLI, requested model, the model verified from the CLI's own log or
  header, verification status, watchdog/timeout outcome, exit code.

## Constraints

- **Advisory only.** External CLI output never auto-commits. The orchestrating agent verifies,
  adjudicates, and owns the result; record the verdict.
- **Never self-authorize a trust or permission bypass.** If a CLI is workspace-trust-gated (e.g.
  gemini-cli exits 55) or stalls on a permission prompt, ask the user to run it (e.g. via the
  session `!` prefix). Never pass `--skip-trust`, `--dangerously-skip-permissions`, or
  `--dangerously-bypass-approvals-and-sandbox`. See `references/cli-invocation-contracts.md`.
- **Never accept output whose model is unverified.** Discard and re-dispatch on a mismatch; a
  silent fallback is invisible in the answer itself.
- **Never point a spawned CLI at a TCC-protected path** (e.g. an iCloud vault). Route reads to a
  local mirror. See `references/vault-and-websearch-grounding.md`.
- **Never hardcode a model identifier in a dispatch path.** Route by tier, resolve at dispatch time,
  and confirm from the provenance line which model actually answered.
- **Scope agentic prompts.** agy `--print` wanders into the workspace on vague prompts; give it
  explicit paths and a neutral working directory, and never expose secrets to an agent that also
  has untrusted input and network egress (the lethal trifecta).
- Vendor-neutral: name CLIs as tools; "the agent" is the orchestrator. Do not steer toward a vendor.

## Troubleshooting

| Signal | Action |
|--------|--------|
| Provenance says `backend_verified=MISMATCH`, or the wrapper exits 3 | The CLI accepted the model identifier but ran another one. Discard the output — it is not from the model you asked for. Re-dispatch using the display-label form of the model. |
| agy returns weak / generic output but exits 0 | Same silent fallback, caught after the fact. Check the provenance line before reading the answer. |
| `agy models` hangs | It is a server-backed call. Never put it in the dispatch path; choose the model from the tier table instead. |
| Headless call hangs (no output, minutes) | Missing `--effort low` on a host-agent child, or no watchdog. Use `run_headless_claude.sh`. |
| CLI exits 55 / "workspace not trusted" | Trust gate. Ask the user to run it; do not self-authorize a bypass. |
| agy reads unrelated files / shell history | Prompt too broad. Tighten it, pass absolute paths, run from a neutral directory. |
| Spawned read of an iCloud path fails | TCC blocks spawned processes. Use the local mirror path. |
| codex errors on the model | Omit the model argument and let codex resolve its configured default. |
| Wrapper exits 124 | Either the watchdog or the CLI's own response timeout fired; the provenance line says which. Raise `TIMEOUT`, narrow the task, or split it. |

## References

| Need | Reference |
|------|-----------|
| Per-CLI flags, exact incantations, every gotcha, watchdog design | `references/cli-invocation-contracts.md` |
| Advisory-only doctrine, adjudication, verification, recording | `references/adjudication-and-provenance.md` |
| Vault mirror + websearch grounding, TCC, citation discipline | `references/vault-and-websearch-grounding.md` |
| When to fan out and how to reconcile divergent answers | `references/fan-out-and-reconciliation.md` |
| Effort flags and small/mid/frontier cost tiers per CLI | `references/effort-and-cost-tiers.md` |

## Templates and scripts

- `scripts/run_agy.sh` — watchdog-wrapped agy dispatch; per-run log, response timeout under the watchdog, and a backend check that fails on a silent model fallback.
- `scripts/run_codex.sh` — watchdog-wrapped `codex exec`; no model override unless you pin one, and the resolved model is read back from codex's own run header.
- `scripts/run_headless_claude.sh` — watchdog-wrapped headless host-agent child; forces `--effort low`.
- `templates/consult-record.md` — adjudication + provenance record skeleton (one per CLI call).
- `examples/vault-digest-walkthrough.md` — a real research-digest run: agy at the frontier tier over the ResearchVault mirror + websearch, then adjudication.
