# Consult record — [task short name]

> One record per CLI call. One file per call; do not concatenate. Fill every field.

- **Timestamp:** [YYYY-MM-DD HH:MM:SS]
- **Mode:** consult | research-digest | code-structure | fan-out
- **CLI:** codex | agy | headless-claude
- **Requested model label:** [e.g. "Gemini 3.1 Pro (High)" | gpt-5.5 | opus]
- **Backend verified (from log/provenance):** [yes — backend_label="…" | no — explain]
- **Effort / tier:** [small | mid | frontier; flag or label used]
- **Watchdog:** [TIMEOUT=Ns; fired? yes/no] · **Exit:** [code]

## Prompt

[the exact prompt sent, or a path to it]

## Advisory output

[the CLI output, a snippet, or a path to the captured output]

## Adjudication (the agent's decision)

- **Verdict:** accept | modify | reject
- **Verification done:** [which claims/citations were checked against which primary sources]
- **Conflicts / dissent:** [for fan-out: where models disagreed and how it was resolved]
- **Reasoning:** [why this verdict; what was kept, modified, or discarded]

## Provenance line

[paste the PROVENANCE: line emitted by the wrapper]
