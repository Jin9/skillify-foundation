# ADR-003 — Order service is sole writer of `orders.status`

- **Status:** Accepted (LOCKED; carried over from dry-run #1; verified by grep in REV-L2)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

The order's lifecycle is the consistency anchor of the system. Five different forces want to mutate it: customer-cancel from PENDING_PAYMENT (`order.cancel-mine`), admin-cancel from PENDING_PAYMENT or PAID (`order.update-status-admin`), payment.completed (PENDING_PAYMENT -> PAID), payment.failed/expired (PENDING_PAYMENT -> PAYMENT_FAILED/PAYMENT_EXPIRED), reservation.expired (PENDING_PAYMENT -> PAYMENT_EXPIRED), and admin fulfillment moves (PAID -> PACKING -> SHIPPED -> DELIVERED). If any service other than Order writes the column, race resolution becomes intractable across schema-per-service boundaries.

Dry-run #1 demonstrated this with the PR-003 race (cancel + payment-success arriving on different Kafka topics). The fix was to (a) keep Order as the sole writer and (b) make all consumers state-driven (ADR-009).

## Decision

`UPDATE orders SET status = ...` exists in exactly ONE place in the codebase: `backend/services/order/app/order/access/storage_order.go::UpdateStatus`. Reviewer-L2 verified this with grep across all 7 backends in dry-run #1 and the rule remains in force.

- **Payment** never writes `orders` rows. It publishes `payment.completed | payment.failed | payment.expired` events; Order consumes and transitions.
- **Inventory** never writes `orders` rows. It publishes `reservation.expired`; Order consumes and transitions to PAYMENT_EXPIRED.
- **Checkout** never writes `orders` rows directly — it calls `order.create-from-checkout` (a NEW formally-authored internal contract in dry-run #2) and on compensation calls `order.cancel-on-checkout-failure`. Both of these are Order's own handlers; Checkout does not bypass.

Inside Order, all status mutations go through one transactional pattern:
```
BEGIN;
  SELECT * FROM orders WHERE id = $1 FOR UPDATE;
  -- validate (currentStatus, requestedToStatus) against the §9.2 transitions table
  UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1;
  INSERT INTO order_status_history (order_id, from_status, to_status, actor_user_id, actor_role, reason, at) VALUES (...);
  -- if cancellation: INSERT INTO outbox_events events.order.cancelled
COMMIT;
```

## Consequences

- **Positive:** race resolution is centralized — the Order tx is the serialization point for every status mutation.
- **Positive:** auditability is trivial — `order_status_history` has every transition with from/to/actor/reason/at.
- **Negative:** Checkout cannot create an Order in PENDING_PAYMENT inside its own transaction; it must HTTP-call `order.create-from-checkout`. This is the cost of the schema-per-service constraint (ADR-001) combined with sole-writer.
- **Negative:** every cross-service event consumer must be state-driven (ADR-009) because `event.fromStatus` cannot be trusted across topics.

## Reviewer-L2 follow-up

REV-L2 (dry-run #1) confirmed the rule held in shipped code. Dry-run #2 must add `order.create-from-checkout` and `order.cancel-on-checkout-failure` as fully-implemented handlers (not stubs) so the rule is preserved on the compensation path too — see ADR-007.
