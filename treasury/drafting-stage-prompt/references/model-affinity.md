# Model affinity per stage

The scaffold's default profile binds each stage to a model in
`AGENTS.md` § "Model assignment". A library entry should declare which
model it was written for so prompts that hard-code model quirks aren't
silently reused on a different model.

| Stage | Default model | Why this model | What to write *for* |
|---|---|---|---|
| `research` | Gemini | Long context; reads whole repos and docs. | Prompts that reference whole-repo scans, multi-file synthesis, long-tail document fetch. |
| `plan` | Codex | Strong at structured planning, terminal-native. | Prompts that demand numbered steps, file/line targets, acceptance criteria. |
| `critique` | Claude | Best for architectural/security review. | Prompts that ask for STRIDE walks, P1/P2/P3 tagging, cross-aggregate reasoning. |
| `implement` | Codex | Terminal-native, sandbox-aware. | Prompts that demand minimal-diff edits, compileable code, test-name-first thinking. |
| `review` | Claude | Catches what implementation missed. | Prompts that diff intent vs. result and surface drift. |
| `test` | local | Just runs the test suite; no model. | Not a model-affinity stage. Skip the affinity field. |

## Frontmatter shape

```yaml
---
stage: critique
topic: ddd-aggregate-boundaries
model_affinity: claude
last_validated: 2026-05-06
outcome: win
---
```

## Cross-model migration

When a library entry written for one model is reused on another, the user
must update `model_affinity` and add a `last_validated` row. This skill
does not migrate automatically — the failure modes are too model-specific.

## Anti-affinity

Some prompts win on one model and lose on another. When the squad
confirms an anti-affinity, capture it in `_anti-patterns.md`:

```
- 2026-05-08 · plan/jwt-validation — codex wins; gemini drifts to discussing alternatives instead of writing the plan
```
