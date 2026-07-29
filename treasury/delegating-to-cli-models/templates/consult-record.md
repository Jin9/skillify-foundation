# Consult record — [task short name]

> One record per CLI call. One file per call; do not concatenate. Fill every field.

- **Timestamp:** [YYYY-MM-DD HH:MM:SS]
- **Mode:** consult | research-digest | code-structure | fan-out
- **CLI:** codex | agy | headless-claude
- **Requested model:** [as passed, e.g. "Gemini 3.1 Pro (High)" | gpt-5.6-sol | opus | (config default)]
- **Model that answered:** [backend_label / resolved_model from the provenance line]
- **Verified:** [yes | MISMATCH — discard and re-dispatch | no — explain why it could not be checked]
- **Effort / tier:** [small | mid | frontier]
- **Watchdog:** [TIMEOUT=Ns] · **Timeout:** [none | watchdog | CLI response timeout] · **Exit:** [code]

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
