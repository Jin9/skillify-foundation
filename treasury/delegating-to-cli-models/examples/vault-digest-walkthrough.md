# Walkthrough — research-digest mode (real run)

A real research-digest run captured while building this skill: dispatch agy/"Gemini 3.1 Pro (High)"
to digest five ResearchVault reports (local mirror) plus a web pass, then adjudicate.

## 1. Frame & route

Goal: ground a skill about orchestrating headless model CLIs. Mode: research-digest. CLI: agy at the
frontier tier. Inputs: five `NN-final_report.md` files under the local mirror
`~/Desktop/project/research/literature/` (the iCloud original is TCC-blocked to spawned processes).
The prompt listed the exact absolute paths, forbade reading anything else, and demanded a
fixed-section Markdown digest.

## 2. Dispatch (through the wrapper, neutral CWD)

```bash
TIMEOUT=540 bash scripts/run_agy.sh "Gemini 3.1 Pro (High)" "$(cat harvest-prompt.txt)" "$HOME/agy-harvest-scratch"
```

## 3. Ground-truth the backend

```text
PROVENANCE: cli=agy requested_label="Gemini 3.1 Pro (High)" backend_label="Gemini 3.1 Pro (High)" timeout=540s exit=0 log=".../cli-20260615_223818.log"
```

`backend_label` equals the requested label — the call really ran on Gemini 3.1 Pro (High), not a
silent Flash fallback, and the watchdog did not fire.

## 4. Advisory output (excerpt)

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
