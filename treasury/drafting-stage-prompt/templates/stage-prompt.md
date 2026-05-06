---
stage: <research|plan|critique|implement|review|test|custom>
topic: <kebab-slug>
model_affinity: <gemini|codex|claude|local|n/a>
last_validated: <YYYY-MM-DD>
outcome: <win|update>
---

# <Topic title — human-readable>

## Prompt

```
<verbatim prompt body — do not edit>
```

## Why it works

<one paragraph: the failure mode this prompt prevents, the constraint it
imposes on the model, the structural reason it survives the model's
default drift>

## Success metric

<one line: the observable signal that proved it worked — "critique
surfaced 3 P1 issues we hadn't seen", "plan listed every changed file
with line ranges", "implement produced a 14-line diff that compiled">

## Failure mode

<one line: the most likely way this stops working — "model upgrade
changes JSON schema preference", "longer goals overwhelm the constraint",
"works on Codex, drifts on Gemini">

## Invocation example

```bash
# Profile-driven: prompts/library/<stage>/<topic>.md is loaded by the
# stage runner when the user sets the corresponding env var, e.g.
PLAN_PROMPT_PREFIX="$(cat prompts/library/plan/<topic>.md)" \
  just workflow <id> "<goal>" 5.00
```

## Notes

<optional: edge cases, variants, related entries in the library>
