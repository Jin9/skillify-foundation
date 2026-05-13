# Profile selection

The scaffold ships three profiles in `profiles/`. Pick one before
`just workflow`. Custom profiles live alongside; see
`authoring-scaffold-profile` skill.

## Built-in profiles

| Profile | Stages | Default models | Use for |
|---|---|---|---|
| `code` (default) | research → plan → critique → 🔒 implement → review → test | gemini / codex / claude / codex / claude / local | Refactors, feature work, bug fixes that touch source code. |
| `literature` | research → plan → critique → review | gemini-heavy | Paper synthesis, literature review, market scans — no code change. |
| `dataset` | research → plan → critique → 🔒 implement → review | codex-heavy | EDA, analysis notebooks, data-wrangling tasks. |

Set with `WORKFLOW_PROFILE=<name> just workflow ...`.

## Decision rules

1. **Will this change source files?** No → `literature`. Yes → continue.
2. **Is it primarily analysis on data, with notebooks as the artifact?** Yes →
   `dataset`. No → `code`.
3. **Is the goal scoped tightly enough that one of the three fits?** No →
   recommend authoring a custom profile via `authoring-scaffold-profile`.

## Cap-tier interaction

A profile that skips `implement`/`test` should run with a cap from the
$0.50–$2.00 tier. Reserve $5+ for full pipelines that include `implement`.

## Sandbox interaction

The sandbox flag (`IMPLEMENT_SANDBOXED=1`) only matters when the chosen
profile includes `implement`. `literature` skips implement, so the flag is
inert; do not advertise it.

## Surface to the user

When confirming the launch, print:

```
Profile: <name>
Stages:  <stage1> → <stage2> → ... → done
Cap:     $<n>.<nn>
Sandbox: <on|off|n/a>
```

Let the user override before kicking off.
