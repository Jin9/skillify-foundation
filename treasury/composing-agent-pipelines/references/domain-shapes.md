# Domain → Phase Shape Mapping

When Compose runs, it selects a phase subset based on the detected task domain. The composer prints the chosen shape upfront so the user can override before phases start.

## The four domains

| Domain trigger phrasing | Domain key |
|-------------------------|------------|
| "research", "investigate", "what does the codebase do for X" | `research` |
| "review", "audit", "is this safe", "find bugs" | `code-review` |
| "design", "plan how to build", "how should we implement" | `implementation` |
| "decide", "choose between", "which option", "trade-offs" | `decision` |

## Phase shapes

Each domain runs a specific subset of phases in this order:

| Domain | Phases (in order) | Skipped | Rationale |
|--------|-------------------|---------|-----------|
| `research` | Plan → Gather → Analyze → Review → Validate → Compact | Decide | Output is a memo with citations, not a recommendation. |
| `code-review` | Plan → Gather → Analyze → Review → Validate → Decide → Compact | (none) | All phases — adversarial review and fact-checking are non-negotiable for code that ships. |
| `implementation` | Plan → Gather → Analyze → Decide → Compact | Review, Validate | Default skips quality gates; opt-in for high-stakes designs (see Override below). |
| `decision` | Plan → Gather → Analyze → Decide → Validate → Compact | Review | Validation is run *after* Decide to confirm the recommendation's claims; standalone Review is replaced by the Decide phase's options-comparison. |

## Override protocol

After detecting the domain, the composer prints:

```
Domain: code-review
Phase shape: Plan → Gather → Analyze → Review → Validate → Decide → Compact
Override? Reply with: 'shape: <phase1>,<phase2>,...' or 'proceed'.
```

Wait for `proceed` or a custom shape before running Plan. Custom shapes must:
- Start with `plan` (no exceptions — every pipeline begins with decomposition).
- End with `compact` (final artifact is mandatory).
- Use only the 7 phase keys: `plan`, `gather`, `analyze`, `review`, `validate`, `decide`, `compact`.
- Not contain duplicates.

## Domain detection edge cases

- **Mixed signals** (e.g., "review and decide"): default to `code-review` and let the user override.
- **No clear domain**: ask one disambiguating question with the four domain options. Do not guess.
- **High-stakes implementation** (e.g., security-sensitive, irreversible migration): suggest upgrading to `code-review` shape via the override prompt.

## Why these shapes

- **Research skips Decide** because research output is descriptive, not prescriptive. Forcing a recommendation creates false certainty.
- **Code review keeps everything** because the cost of a missed bug exceeds the cost of one extra phase.
- **Implementation skips Review/Validate by default** because design docs are inherently exploratory; quality gates inhibit early exploration. Production-bound designs should opt in.
- **Decision keeps Validate, drops Review** because the Decide phase's options comparison already serves the adversarial role; Validate is preserved to fact-check the recommended option's claims.
