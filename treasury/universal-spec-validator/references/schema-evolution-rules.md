# Safe-vs-breaking schema-evolution rulebook

Distilled from research report *Strict Schema Evolution for Tool Registries*.
A tool registry is the model's sole source of truth; when the underlying spec
drifts the model keeps emitting calls shaped to the stale contract. Measured:
~8% of responses unparseable, ~5% wrong field types; at 5%/step a 12-step run
has ~46% chance of at least one failure. JSON Schema has **no built-in
compatibility enforcement** — the enforcement must live in this gate.

Evolution checks require a `--baseline` (the previously-shipped spec). With no
baseline, emit one `info` finding `schema-evolution: skipped (no baseline)` and
exit-neutral on this axis.

## Compatibility model
- **Backward** (default): a caller built on the NEW schema still handles OLD-shaped data.
- **Forward**: a caller on the OLD schema still handles NEW-shaped data.
- **Full** = backward + forward. **Transitive** = holds across all prior versions.
- Config `compatibility_mode: backward|forward|full` (default `backward`); `transitive: true|false`.

## SAFE changes (never blocking)
- Add an **optional** field **with a default**.
- Add a new enum value when the content model is **open**.
- Widen a type (int → number), relax a constraint (lower `minLength`, raise `maxItems`).
- Add a new tool/operation (backward-safe; *not* forward-safe — flag only under forward/full).

## BREAKING changes (blocking under the active compatibility mode)
- Remove a field, or remove/rename a `required` field.
- Add a **required** field **without a default**.
- Change a field's type (string → integer, scalar → object, etc.).
- Narrow an enum (remove a value) or tighten a constraint (raise `minLength`, add `pattern`, add `required`).
- Make a previously-optional field required.
- Change a stable identifier: Protobuf field number reuse, renamed `$id`/operationId, reused deprecated field.

## Detection & response
- Diff baseline vs candidate structurally (oasdiff-style: 450+ change categories collapse to the SAFE/BREAKING split above).
- A breaking change under the active mode → `severity: high`, `gate: block`.
- Risk-tiered: config may map a specific rule to `warn` (pause+notify) instead of `block` for medium-risk evolutions — but type-change and required-without-default stay `block` by default.
- Recommend, in the report (not auto-applied): regenerate the tool def from the new contract OR pause/disable the tool + alert owner; enforce exact version pins; govern the registry's own schema with these same rules.

Boundary: this axis guarantees *shape* compatibility only. A well-formed call
to the *wrong* operation is semantic correctness — out of scope. Every finding
cites `schema-evolution-rules.md`.
