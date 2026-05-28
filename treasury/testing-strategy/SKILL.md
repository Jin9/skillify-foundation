---
name: testing-strategy
description: >
  Decide what to test and at which level (unit / integration / contract / e2e / regression) and map a requirement's business logic into concrete test scenarios, so coverage guarantees behavior rather than a number. Use when the user asks "what's our test strategy for this", "which test levels do we need", "turn this business logic into test scenarios", or "is our coverage meaningful". Produces a layered test-level plan, a business-logic-to-scenario map, edge-case parameters, contract-test boundaries, and a meaningful-coverage check. Do NOT use for localizing a specific failing test or bug (use progressive-bug-hunter) or for security test cases (use reviewing-software-security).
---

# testing-strategy

## Purpose
Decide **what to test, at which level, and as which scenarios** — turning a requirement's business logic into concrete cases so coverage *guarantees behavior*, not a percentage. A high number with the wrong tests is no coverage.

## When to use
- Triggers: *"what's our test strategy for this"*, *"which test levels do we need"*, *"turn this business logic into test scenarios"*, *"is our coverage meaningful"*.
- **Not this skill:** localizing a specific failing test or bug → `progressive-bug-hunter`; security test cases → `reviewing-software-security`.

## Input
- A requirement + its business logic (rules, branches, edge cases) and the service/API surface it touches.
- (Optional) existing coverage numbers, current test suite, service boundaries.

## Output
A **Test Strategy** artifact + a filled checklist + a verdict. Skeleton:
```
# Test Strategy — <feature>
Test levels (level — what it covers — who owns):
  - unit         — <business rule in isolation>        — <owner>
  - integration  — <handler + infra wiring>            — <owner>
  - contract     — <service / API boundary>            — <owner>   # see contract-testing
  - e2e          — <critical user journey only>        — <owner>
  - regression   — <set replayed on every change>      — <owner>
Business-logic → scenario map:
  - <rule> → <happy path> / <failure path> / <boundary>
Edge-case parameters (probe to reduce risk):
  - <param>: <min / max / null / overflow / duplicate>
Contract tests (which boundaries, consumer-driven):
  - <consumer> ↔ <provider>: <pact / tool-schema>
Coverage target:  business logic = 100%
"Meaningful?" check:  [ ] tests assert behavior, not lines  [ ] edge params covered
Verdict:  STRATEGY READY  |  NEEDS SCENARIOS
```

## Decision rules
1. **Business-logic coverage = 100%**; additionally probe **possible edge-case parameters** to reduce risk.
2. Iron rule — **unit tests must be written to guarantee every API**, and domain logic stays separated from infra so it is testable in isolation.
3. **Testable** ranks among the top good-code properties because it "guarantees the system works correctly" — prioritize it over performance.
4. In review, **unit-test coverage of ALL business logic** is the #2 priority — the strategy must make that coverage reachable.
5. **100% coverage but still bugs → the tests don't test the right things** — treat the number as a smoke alarm, not a guarantee; demand behavior assertions.
6. Put a **contract-test level between unit and e2e**: consumer-driven contracts (Pact) hold service/REST boundaries without running both sides; for agent/tool boundaries, tool-schema governance is the highest-leverage contract point.

## Checklist
- [ ] Test levels chosen with what-each-covers + owner
- [ ] Every business rule mapped to happy / failure / boundary scenarios
- [ ] Edge-case parameters enumerated (null, max, duplicate, overflow)
- [ ] Contract boundaries identified and consumer-driven
- [ ] Regression set defined
- [ ] Coverage target set to 100% of business logic
- [ ] "Meaningful coverage" check passed (asserts behavior, not lines)
- [ ] Verdict recorded (ready / needs-scenarios)

## Anti-patterns (never do)
- Chase a coverage **number** while the tests assert nothing meaningful.
- Leave any **business logic** path untested.
- Test domain logic only through infra/e2e because it wasn't separated.
- Skip the **contract level** and discover boundary drift in production.
- Gold-plate e2e for paths a unit/contract test already guarantees.
- Write production code without tests.

## Example
**Input:** "Wallet top-up via Provider X." **Output (excerpt):** *Levels:* unit (ledger math, idempotency key), contract (wallet ↔ Provider X callback, consumer-driven Pact), e2e (one happy top-up journey). *Scenario map:* idempotency rule → duplicate callback must not double-credit; provider declines → no ledger entry. *Edge params:* amount = 0 / max / negative; callback arrives twice; currency mismatch. *Meaningful check:* assert balance + ledger row, not just "200 OK." → **STRATEGY READY**.

## Human approval gate
**Stop.** The team owns the final scenario set and which boundaries get contract tests; the skill recommends, it does not sign off. The lead signs off on the strategy artifact, and team-affecting decisions are not made alone.
