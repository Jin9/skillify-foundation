# Prompt-prefix style

Profile-level prefixes are short steering strings the runner prepends to
its built-in stage prompt. They are not the place for full prompt bodies.

## Length and shape

- ≤ 200 characters. Anything longer belongs in
  `prompts/library/<stage>/<topic>.md` (use `drafting-stage-prompt`).
- End with a trailing space — the runner concatenates without inserting
  one. Without the trailing space, you get `Prefix.Builtin prompt...`.
- Single line. No embedded newlines.

## Tone and grammar

- Imperative mood: "Optimize for synthesis."
- Concrete focus area: "Cite paper titles."
- Domain hint: "Treat the plan as an outline, not code steps."
- No "please", no "you are an AI assistant", no role-play setup.

## Good examples (from `literature.sh`)

```bash
RESEARCH_PROMPT_PREFIX='Optimize for synthesis across many sources. Cite paper titles. '
PLAN_PROMPT_PREFIX='Treat the "plan" as a draft outline of the synthesis document, not code steps. '
CRITIQUE_PROMPT_PREFIX='Focus on missing perspectives, weak citations, and unsupported claims. '
REVIEW_PROMPT_PREFIX='Review the synthesis for completeness and citation accuracy, not code quality. '
```

## Anti-examples

| Bad | Why |
|---|---|
| `'You are an expert AI assistant. Please help write code...'` | role-play setup; verbose; brand-leaning. |
| `'Use Claude best practices'` | brand name; "best practices" is empty. |
| `'For each finding, output: …(80 lines of detail)…'` | long; belongs in library entry. |
| `'Optimize.'` | too vague; gives the model no signal. |

## When to escalate to the library

If the prefix you want to write would exceed 200 characters, stop and
hand off to `drafting-stage-prompt`. The library entry can be loaded at
runtime by the user's profile or per-workflow:

```bash
PLAN_PROMPT_PREFIX="$(cat prompts/library/plan/<topic>.md) " \
  WORKFLOW_PROFILE=<slug> just workflow <id> "<goal>" 5.00
```

This keeps profiles small and keeps long-form prompts in a place where
the Friday review can vet them.
