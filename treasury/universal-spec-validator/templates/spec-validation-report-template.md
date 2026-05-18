# Spec validation report

> Rendered by `scripts/validate_spec.py`. This is the shape of the human
> `spec-validation.md`; the canonical machine record is `spec-validation.json`.

- specs: {N}  block: {B}  warn: {W}  info: {I}
- exit_code: {0|1|2}  (0 pass / 1 blocked / 2 error)
- by axis: {portability: .., command-safety: .., schema-evolution: ..}

## CRITICAL
- [block] **command-safety/C1** `path/to/spec` @ `$.tools[i].command` — <message> (command-safety-rules.md#C1)

## HIGH
- [block] **schema-evolution/E2** `path/to/spec` @ `$.field` — <message> (schema-evolution-rules.md)
- [block] **portability/P1** `path/to/spec` @ `tool-call dialect` — <message> (portability-rules.md#P1)

## MEDIUM
- [warn] **portability/P2** `path/to/spec` @ `$.input` — <message> (portability-rules.md#P2)

## INFO
- [info] **schema-evolution/E0** `path/to/spec` @ `-` — checks skipped (no baseline)

<!-- One bullet per finding, grouped by severity (critical→info). Every bullet
carries: gate, axis/rule-id, spec path, locator, message, and the rule_ref into
references/. No remediation is auto-applied — this skill is a gate. -->
