# Schema flexibility: conform vs infer

The skill produces a contract in one of two modes. The presence of `target_schema` decides which.

## Decision

```
target_schema provided?
├─ yes → CONFORM: map source → the given schema's fields
└─ no  → INFER:   design a minimal contract shape from the source
```

Never guess a schema the caller did not give in order to "be strict." Absence of a schema is a deliberate signal to infer.

## Conform mode

You are given the downstream contract. Your job is faithful mapping, not redesign.

- Populate every **required** field from the source's must-preserve content.
- If the source genuinely does not supply a required field, leave it absent and record `{ item: "<field>", reason: "not-in-source" }` in `_meta.dropped`. Surface this in `_meta.notes` too if the consumer must handle it.
- **Never fabricate** a value to satisfy a required field. A hallucinated value is worse than an honest gap — it corrupts every downstream stage silently. This is the single most important conform-mode rule.
- Respect the schema's types and enums. If the source's value does not fit the declared type/enum, record the mismatch in `dropped` rather than coercing it into a lie.
- Do not add fields outside the schema. Extra signal the consumer did not ask for belongs in `_meta.notes`, not in `contract`.
- Optional fields: fill them when the source supports them; omit silently when it does not (no `dropped` entry needed for optionals).

## Infer mode

No downstream contract exists yet, so you design a lean one.

- Build the **minimal** field set that carries the must-preserve content — nothing speculative.
- Favor **flat, stable, predictable** names (`status`, `amount`, `effective_date`) over deep nesting or clever structure. A downstream consumer must be able to anticipate the shape; unpredictable inferred shapes break composition.
- Use plain, conventional types: strings, numbers, booleans, arrays of those. Avoid bespoke nested objects unless the source's structure genuinely demands them.
- If a downstream contract already exists (the user mentions a consuming stage, or a schema is "somewhere"), prefer conform mode: say so in `_meta.notes` and recommend the caller pass `target_schema` next time. Inferred contracts are a bootstrap, not a substitute for an agreed contract.
- Keep the inferred shape **idempotent**: the same source should yield the same field names on a re-run. Do not let phrasing of the source drive field naming.

## When to fail instead of extract

Stop and report rather than emit a misleading contract when:
- `source` is empty or absent — ask for it.
- In conform mode, `target_schema` is malformed or unparseable — report it; do not silently fall back to infer.
- The source is the wrong artifact for the requested target (e.g., a `target_schema` for invoices but the source is a meeting transcript) — report the mismatch; partial garbage is worse than a clear "this source does not populate this contract."

## Chaining (pre/post-conditions)

- **Pre (conform):** the consuming stage's `target_schema` must be available to the caller.
- **Post (always):** the envelope must pass `scripts/validate_extraction.py` — envelope integrity is a deterministic gate, not an LLM judgment.
- **Post (downstream):** consumers act on `contract`; they pull `source_ref` only when `_meta.dropped`/`coverage` tells them something they need was set aside.
