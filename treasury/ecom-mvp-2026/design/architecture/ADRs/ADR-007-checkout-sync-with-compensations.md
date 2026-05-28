# ADR-007 — Checkout uses synchronous orchestration with compensations (not a saga choreography)

- **Status:** Accepted (LOCKED)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

Checkout commits cart-to-order across five downstream services (cart, identity, catalog, inventory, order, payment). Two architectural shapes are common:

1. **Synchronous orchestration with explicit compensations.** Checkout calls services in order; on a partial failure it issues compensating calls (release reservation, cancel order). Uses HTTP under the hood; a single durable `idempotency_keys` row in checkout_db survives the orchestrator's own crash.
2. **Choreographed saga (event-driven only).** Checkout writes a saga-start event; each downstream service consumes events and emits its outcome; Checkout consumes outcomes and decides the next step.

For an MVP with strict latency budget (PERF-002: p95 checkout < 500ms) and exactly one orchestration shape, (1) is materially simpler.

## Decision

Synchronous orchestration with compensations. The 10-step flow is documented in `contracts.json#/contracts/checkout.commit.internal_orchestration` and reproduced in `connectivity.md`. The compensation matrix is:

| Failure point | Compensations (in order) |
|---|---|
| Step 6 (`inventory.reservation.create`) | none — no state mutated |
| Step 7 (`order.create-from-checkout`) | `inventory.reservation.release(reason=ORDER_CANCELLED)`; on retry-exhaustion log CRITICAL + saga DEGRADED |
| Step 8 (`payment.intent.create`) | `order.cancel-on-checkout-failure(reason=...)` THEN `inventory.reservation.release(reason=ORDER_CANCELLED)`; on retry-exhaustion log CRITICAL + saga DEGRADED |
| Step 9 (`cart.clear-on-checkout`) | none — log WARN; do NOT roll back order/payment (cart inconsistency is cosmetic) |

Hard requirements that follow from this shape:

- **Every internal call is idempotent on `orderId`** (cross-cutting.idempotency.internal_orderId_rule). Compensation calls and retries cannot create duplicate state.
- **`order.cancel-on-checkout-failure` MUST be a real handler, not a stub** — REV-L2-004 dry-run #1 left this as a 501 stub and a compound dead-end resulted. ADR-003 says Order is the sole writer of order.status; this handler is the canonical sole-writer entry-point for the compensation path.
- **Reservation TTL (15min, see cross-cutting.pricing.constants.RESERVATION_TTL_MINUTES) is the safety net** — if a compensation call itself fails AND its retry budget exhausts, the inventory.reservation.sweep-expired sweeper releases the reserved stock and emits `events.reservation.expired`, which Order's state-driven consumer translates into PAYMENT_EXPIRED. Both halves (sweeper + Order consumer) must ship.
- **Saga reconciler:** Checkout runs a 60s job that scans `saga_log WHERE status='DEGRADED' AND created_at < NOW() - INTERVAL '1 hour'` and emits a metric `checkout_saga_degraded_unreconciled_total` so on-call sees stuck rows. REV-L2-008 followup.

## Consequences

- **Positive:** simple to reason about; one place to read the orchestration shape (`checkout.commit.internal_orchestration`).
- **Positive:** latency is bounded by the longest synchronous chain (about 5 hops); meets PERF-002.
- **Positive:** exactly-once semantics under the orderId-keyed idempotency rule for steps 6/7/8.
- **Negative:** Checkout becomes a single point of failure for the orchestration. If Checkout crashes mid-flight, the in-flight tx is aborted and the customer's POST returns a network error; the customer retries with the same Idempotency-Key which causes the persisted INFLIGHT row to short-circuit. The reservation TTL safety-net handles the case where Checkout never recovers.
- **Negative:** the compensation matrix is the most-tested code surface — every cross-component reviewer pass should re-validate it.
