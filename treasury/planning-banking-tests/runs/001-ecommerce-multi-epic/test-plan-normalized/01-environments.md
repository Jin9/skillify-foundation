---
artifact_type: qa-environments
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Environments

## ENV-CART-integration

- Purpose: Integration env for cart CRUD: real DB, real cart service; auth stubbed to mint session for synthetic customers.
- Required mocks: auth-token-mint, product-catalog-readonly-mirror
- Real services: cart-service, product-stock-read-api, postgres

## ENV-CHK-e2e

- Purpose: End-to-end env for checkout submit and payment callback flows: real order service, real stock reservation, real audit sink, mock PSP.
- Required mocks: mock-PSP-webhook-server, notification-log-sink, audit-sink
- Real services: order-service, cart-service, stock-reservation-service, postgres

## ENV-CHK-race-deterministic

- Purpose: Deterministic concurrency env for race-condition and TTL-sweep tests; provides clock-control API, lock instrumentation, and barrier-synced parallel submit harness. Required for TC-3-002, TC-3-007, TC-3-008, TC-4-003.
- Required mocks: mock-PSP-webhook-server, frozen-clock, lock-instrumentation-probe, concurrency-harness
- Real services: order-service, stock-reservation-service, postgres-with-isolation-level-instrumented

## ENV-IDENTITY-compliance

- Purpose: Compliance scenarios (DSAR, breach-notification simulation, consent flow). Scope blocked until P1 governance gaps resolved.
- Required mocks: pdpc-notification-stub, dsar-workflow-engine-stub
- Real services: customer-db (ephemeral), audit-log-sink (in-memory verifier)

## ENV-IDENTITY-e2e

- Purpose: End-to-end browser-driven scenarios for signup, login, address-book happy paths.
- Required mocks: email-verification-stub, rate-limit-stub-disabled-by-default
- Real services: customer-db (ephemeral), auth-service (ephemeral), session-store (ephemeral)

## ENV-IDENTITY-integration

- Purpose: Service-level integration scenarios: uniqueness, anti-enumeration, throttle, authz denial, audit emission.
- Required mocks: clock-control (time-anchor injection), rate-limit-stub-controllable
- Real services: customer-db (ephemeral), auth-service (ephemeral), audit-log-sink (in-memory verifier)

## env-ci

- Purpose: CI integration tests for public catalog read paths (listing, PDP) with seeded relational DB and stubbed downstream services
- Required mocks: postgres, auth-stub-guest
- Real services: _(none)_

## env-e2e-fulfil

- Purpose: End-to-end browser tests for customer order-detail page, status timeline render, OOB refund operator workflow.
- Required mocks: payment-gateway-mock
- Real services: order-service, stock-service, audit-log-store, rbac-service, web-frontend, notification-service

## env-gov-integration

- Purpose: Integration env for governance/observability epic: real app, real audit-log store, real notification-log store, mock-shipping and mock-payment for cross-border transfer test.
- Required mocks: mock-shipping, mock-payment, mock-clock
- Real services: app, audit-log-store, notification-log-store, pdp-renderer

## env-integration

- Purpose: Integration tests for admin product CRUD with audit emission, idempotency store, and authz enforcement
- Required mocks: postgres, audit-event-sink, idempotency-store, auth-stub-admin
- Real services: _(none)_

## env-integration-fulfil

- Purpose: Integration tests for order state machine, audit emission, authz, idempotency - in-process service with mocked payment gateway and real audit-log DB.
- Required mocks: payment-gateway-mock, notification-channel-mock
- Real services: order-service, stock-service, audit-log-store, rbac-service

