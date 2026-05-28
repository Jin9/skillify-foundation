# Ambiguity and Hidden-Requirement Rules

## Eight Linguistic Detectors

1. `lexical`: vague words such as urgent, reasonable, appropriate, fast, simple, normal.
2. `syntactic`: clauses or parentheses that change the meaning of the main sentence.
3. `pragmatic`: context-dependent references such as this, that, usual, standard, current flow.
4. `pronominal`: pronouns crossing more than one sentence or possible antecedent.
5. `quantifier`: many, several, most, some, a lot, or percentages without denominators.
6. `modal`: may, might, could, should in prescriptive sections.
7. `commitment_conditionality`: yes, ok, will, agreed, pending, subject to, once, if.
8. `phase_boundary_drift`: phase 1, phase 2, deferred, later, not now, pressure to pull forward.

Placeholder tokens such as `TBD`, `?`, `(?),` `$Xk`, `<owner>`, or anonymous numeric values also produce findings.

## Severity Floors

- `P1`: Legal absent on regulatory content, regulator citation unresolved on T1, tipping-off violation, missing PII inventory, irreversible action without compensation.
- `P2`: value needed before sprint planning, unnamed policy owner, conflicting AC, modal hedge in rule text, missing SLA, ambiguous state machine.
- `P3`: documented assumption can proceed with review, such as vague label without regulatory or operational impact.

## Ten Hidden-Requirement Frames

1. Scale and capacity: always apply.
2. Time and timing: always apply.
3. Money and economics: apply when pricing, settlement, fees, chargeback, revenue, cost, or limits appear.
4. Regulatory and legal: apply when PII, payment, named jurisdiction, regulator, consumer-facing status, sanctions, AML, KYC, EDD, SAR, biometric, or audit logging appears.
5. Operational and organizational: always apply.
6. Failure and edge cases: always apply.
7. Integration and dependencies: apply when external systems, vendors, APIs, files, events, or batch jobs appear.
8. Localization and culture: apply when named market, language, currency, timezone, holiday, or jurisdiction appears.
9. Lifecycle: always apply.
10. Customer experience: apply when any human-facing surface or support workflow appears.

## Coverage Rules

For non-failure outputs, `frames_applied` plus `frames_skipped` must cover all frame numbers 1 through 10. Every skipped frame needs a non-empty reason. If an applied frame yields no question or assumption, record a P2 coverage question unless the source is unusually complete and that reason is explicit.
