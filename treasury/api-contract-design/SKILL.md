---
name: api-contract-design
description: >
  Design a stable, clear API contract (action, request/response, error cases, validation, idempotency, versioning) for an API you expose, so consumers can rely on it and it stays backward-compatible. Use when the user asks "design this API/endpoint", "what's the request/response contract", "do we need to version this", or "how should this API handle errors". Produces an API Contract artifact plus a checklist and a human approval gate, with unit-test scenarios derived from the business logic. Do NOT use for validating an existing schema for breaking changes (use universal-spec-validator).
---

# api-contract-design

## Purpose
Design a stable, clear API contract — the action, request/response, error cases, validation, idempotency, and
versioning — for an API **you expose** (not a third party's), so consumers can rely on it and it stays
backward-compatible.

## When to use
- Triggers: *"design this API/endpoint"*, *"what's the request/response contract"*, *"do we need to version this"*, *"how should this API handle errors"*.
- **Not this skill:** validating an existing schema for breaking changes → `universal-spec-validator`.

## Input
- A domain model or use-case.

## Output
An **API Contract** artifact + checklist. Skeleton:

```
# API — <domain>/<aggregate>/<action>
Method + path:   POST /api/v1/<domain>/<aggregate>/<action>
Request:   <fields, types, required>
Response:  <success shape + status>
Errors:    <code → meaning → client action>
Validation:    <rules>
Idempotency:   <key? scope?>
Versioning:    <needed? strategy>  Backward-compat:  <impact>
Change-log:    <date — change>
Unit-test scenarios:  <from the business logic>
```

## Decision rules
1. **Start from the command:** what action is the customer sending?
2. Derive **unit-test scenarios from the business logic** as you design — the contract isn't done until they're listed.
3. **Naming baseline:** `/api/v1/<domain>/<aggregate>/<action>`.
4. **Check backward compatibility**, then decide if versioning is needed — don't break consumers silently.
5. The API spec is **owner-driven**: the team writes it, the squad-lead + team review it, and the **lead signs off**; critical issues the lead discusses directly, minor ones the team decides.
6. For money/duplicate-prone operations, require an **idempotency key**.

## Checklist
- [ ] Action/command is explicit
- [ ] Request + response schemas defined
- [ ] Error cases enumerated with client action
- [ ] Validation rules stated
- [ ] Idempotency decided (key + scope where needed)
- [ ] Versioning + backward-compat assessed
- [ ] Change-log entry added
- [ ] Unit-test scenarios derived from business logic

## Anti-patterns (never do)
- Break backward compatibility silently.
- Ship an endpoint without its error cases or test scenarios.
- Omit idempotency on money-moving operations.
- Promise an endpoint/feature without prior design.
- Expose secrets or credentials in plaintext (in examples, headers, or error bodies).

## Example
**Input:** wallet top-up. **Output (excerpt):** `POST /api/v1/wallet/wallet/top-up`, request `{amount, currency,
idempotency_key}`, response `201 {wallet_id, balance}`, errors `409 duplicate → treat as success`, `422 invalid
amount`. Idempotency key required. v1 — no break. Tests: happy path, duplicate key, negative amount.

## Human approval gate
**Stop.** Squad-lead + team review the spec and the lead signs off before it's published as the contract. The
skill drafts the contract; humans ratify it.
