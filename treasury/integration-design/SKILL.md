---
name: integration-design
description: >
  Design resilient integration with an external/3rd-party system (timeouts, retries, fallback, error mapping, idempotency, and a clear ownership/support model) so their failures do not become our outages. Use when the user asks "integrate with this external/3rd-party system", "call this vendor API", "how do we handle their downtime/timeouts", or "who owns this integration". Produces an Integration Design artifact, a checklist, and a human approval gate before a live external dependency is enabled. Do NOT use for security review of the integration such as auth, secrets, or PII (use reviewing-software-security).
---

# integration-design

## Purpose
Design resilient integration with an external/3rd-party system — timeouts, retries, fallback, error mapping, and a
clear ownership/support model — so their failures don't become our outages. This designs how **we consume** an
external/third-party system (their API) — not the API we expose.

## When to use
- Triggers: *"integrate with this external/3rd-party system"*, *"call this vendor API"*, *"how do we handle their downtime/timeouts"*, *"who owns this integration"*.
- **Not this skill:** security review of the integration (auth, secrets, PII) → `reviewing-software-security`.

## Input
- An integration requirement + the external system's characteristics (protocol, SLA, known failure modes, owner).

## Output
An **Integration Design** artifact + checklist. Skeleton:

```
# Integration — <external system>
Owner / support model:  <who, escalation path>
Protocol + SLA:         <rest/grpc/webhook · their stated SLA>
Timeout:                <per call>
Retry:                  <count + backoff + which errors are retryable>
Idempotency:            <key for retried/duplicate calls>
Fallback:               <degrade path when they're down>
Error mapping:          <their error → our domain error → customer message>
Monitoring/alerts:      <what we watch>
```

## Decision rules
1. Before committing, analyze the external stack's **maturity, vendor/community maintenance, and prevention methods** — a clean non-prod test is not proof it survives production.
2. **Assume it will fail:** every external call gets a timeout + bounded retry (with backoff) + a fallback path.
3. **Map their errors to our domain errors** — never leak raw vendor errors to customers.
4. Make **retried calls idempotent** so retries can't double-charge or double-create.
5. Clarify **ownership + support model** up front, and weigh the integration's **cost + maintenance** like any stack choice.

## Checklist
- [ ] Protocol + their SLA recorded
- [ ] Timeout per call
- [ ] Retry policy (count, backoff, retryable errors)
- [ ] Idempotency on retried/duplicate calls
- [ ] Fallback/degrade path defined
- [ ] Error mapping (their error → ours → customer)
- [ ] Monitoring + alerts
- [ ] Ownership + escalation path

## Anti-patterns (never do)
- Call an external system with no timeout or fallback.
- Treat a non-prod integration test as production proof.
- Leak raw vendor errors to customers.
- Ship with no clear owner for when it breaks.
- Expose secrets or credentials in plaintext.

## Example
**Input:** payment-provider top-up. **Output (excerpt):** 3s timeout, 2 retries w/ exponential backoff using the same
idempotency key, fallback = enqueue + reconcile, map provider `5xx → "top-up pending"`, alert on error-rate >1%.
Owner: Payments squad; escalate to vendor TAM after 15 min.

## Human approval gate
**Stop.** A human (and the external system's owner) approves the contract and failure handling before go-live. The
skill designs the integration; it does not enable a live external dependency on its own.
