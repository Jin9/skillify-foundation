# Walkthrough — research-digest mode (real run)

The research-digest run that grounded this skill, captured **2026-06-15**: dispatch agy at the
frontier tier to digest five ResearchVault reports (local mirror) plus a web pass, then adjudicate.
The framing and adjudication below are that run verbatim. The dispatch form and provenance lines in
sections 2 and 3 were **re-verified by live dispatch on 2026-07-30** against the current CLI
contract — the digest itself was not re-run, so section 4 remains the June output.

## 1. Frame & route

Goal: ground a skill about orchestrating headless model CLIs. Mode: research-digest. CLI: agy at the
frontier tier. Inputs: five `NN-final_report.md` files under the local mirror
`~/Desktop/project/research/literature/` (the iCloud original is TCC-blocked to spawned processes).
The prompt listed the exact absolute paths, forbade reading anything else, and demanded a
fixed-section Markdown digest.

## 2. Dispatch (through the wrapper, neutral CWD)

Use the **display-label** form for the frontier model — the equivalent slug is accepted without
error but silently runs the default model (see section 3).

```bash
TIMEOUT=600 bash scripts/run_agy.sh "Gemini 3.1 Pro (High)" "$(cat harvest-prompt.txt)" "$HOME/agy-harvest-scratch"
```

## 3. Ground-truth the model

```text
PROVENANCE: cli=agy requested_model="Gemini 3.1 Pro (High)" backend_label="Gemini 3.1 Pro (High)" backend_verified=yes print_timeout=90s watchdog=120s elapsed=11s timeout=none exit=0 log=".../cli-wrapper-20260730_013907-78960.log"
```

`backend_verified=yes` is the only thing that licenses reading the output: the call really ran on
the frontier model, and neither timeout fired.

This is not a formality. The same request written as a slug — a slug `agy models` itself lists —
produces a normal-looking answer and a zero exit while running a different, cheaper model:

```text
PROVENANCE: cli=agy requested_model="gemini-3.1-pro-high" backend_label="Gemini 3.6 Flash (High)" backend_verified=MISMATCH ... exit=0
error: SILENT FALLBACK — requested "gemini-3.1-pro-high" but the backend ran "Gemini 3.6 Flash (High)".
```

The wrapper exits 3 and the output is discarded. Nothing in the answer text would have revealed it.

## 4. Advisory output (excerpt, 2026-06-15 run)

```text
## Watchdog & timeout patterns
- On systems lacking GNU `timeout` (like macOS), use a background watchdog that kills a hung
  subprocess.
- Target the whole process group to avoid leaving zombie children.
...
## Citations
- agent-handoff-protocols-..., agent-command-safety-..., ... (vault reports)
- LLM Observability — https://openobserve.ai ; Multi-Agent Adjudication — https://medium.com ; ...
```

## 5. Adjudication (the agent's decision)

- **Verdict: accept with modifications.**
- **Verified:** the five cited vault folders exist and were read directly; the compounding-error
  figure checks out (0.85^10 ≈ 20%); the macOS-watchdog and advisory-adjudication guidance matches
  the project's own conventions.
- **Rejected as authoritative:** the web citations are bare domains (`medium.com`, `wandb.ai`, an
  unverified `architect-cli` repo). They were down-ranked to leads, not evidence, and not embedded
  as sources.
- **Kept as a note:** the process-group-termination idea — captured as a hardening note in
  `references/cli-invocation-contracts.md` rather than changing the already-working wrapper.
- **Result:** the grounded principles informed the references; CLI-specifics were taken from the
  authoritative project docs, not from the digest.

The raw digest and prompt were saved to a readable backup for audit. The committed skill content is
the agent's synthesis, not agy's text verbatim.
