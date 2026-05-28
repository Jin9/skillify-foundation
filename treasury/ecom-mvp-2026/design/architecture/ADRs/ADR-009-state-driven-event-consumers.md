# ADR-009 — State-driven event consumers (consumers ignore `event.fromStatus`)

- **Status:** Accepted (LOCKED — closes the dry-run #1 PR-003 race)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

Cross-service event flows in this system traverse three Kafka topics:
- `ecom.payment.events` (payment.completed/failed/expired)
- `ecom.order.events` (order.cancelled)
- `ecom.inventory.events` (reservation.expired)

Per-orderId ordering is preserved WITHIN each topic by the per-orderId partition key + single-threaded outbox relay. Cross-topic ordering is NOT preserved by Kafka; it cannot be.

The PR-003 race (named after the dry-run #1 plan-review finding): a customer customer-cancels their order at the same moment a payment-success callback arrives. Order publishes `events.order.cancelled` (cancelActor=CUSTOMER, fromStatus=PENDING_PAYMENT) on `ecom.order.events`. Payment publishes `events.payment.completed` on `ecom.payment.events`. Inventory consumes both. Depending on which event lands first at Inventory, the consumer must do completely different things — and `event.fromStatus` ('what the producer thought the previous state was') is no help because by the time the second event lands, the producer's view is stale.

## Decision

ALL consumers of cross-topic events are **state-driven**: they read the current state of the affected aggregate inside the consumer transaction (SELECT ... FOR UPDATE) and act on the actual current state. The event's `fromStatus` field is retained in the payload purely for audit/debug; consumers MUST NOT branch on it.

Concretely:

**Inventory consumer of `events.payment.completed`:**
```
BEGIN;
  SELECT * FROM reservations WHERE order_id = $orderId FOR UPDATE;
  -- DO NOT trust event.fromStatus
  if reservation.status == 'RESERVED':
      reserved_qty -= reservation.qty
      sold_qty     += reservation.qty
      reservation.status = 'COMMITTED'
  elif reservation.status == 'COMMITTED':
      log "idempotent no-op"
  elif reservation.status == 'RELEASED':
      log "payment.completed arrived after reservation released — see PR-003"
      metric inventory_reservation_conflict_total += 1
      -- NO mutation; the order will also have ignored the late payment.completed via Order's state-driven consumer
COMMIT;
```

**Inventory consumer of `events.order.cancelled` (the PR-003 race resolution):**
```
BEGIN;
  SELECT * FROM reservations WHERE order_id = $orderId FOR UPDATE;
  -- DO NOT trust event.fromStatus (could say PENDING_PAYMENT but reality is COMMITTED if payment landed first)
  if reservation.status == 'RESERVED':
      reserved_qty -= reservation.qty
      available_qty += reservation.qty
      reservation.status = 'RELEASED'
  elif reservation.status == 'COMMITTED':
      sold_qty -= reservation.qty
      available_qty += reservation.qty
      reservation.status = 'RELEASED'   -- BA PR-006 MVP restock policy
  elif reservation.status == 'RELEASED':
      log "idempotent no-op"
COMMIT;
```

**Order consumer of `events.payment.completed`:**
```
BEGIN;
  SELECT * FROM orders WHERE id = $orderId FOR UPDATE;
  if status == 'PENDING_PAYMENT':
      status = 'PAID'
      INSERT order_status_history (...)
  else:
      log "late payment.completed for non-PENDING_PAYMENT order"
      metric order_consumer_late_event_total{event_type='payment.completed', current_status=$status} += 1
      -- NO state change
COMMIT;
```

The same pattern applies to: Order's consumer of `events.payment.failed` / `events.payment.expired` / `events.reservation.expired`; Payment's consumer of `events.order.cancelled`.

## Consequences

- **Positive:** the PR-003 race resolves correctly under any topic interleaving — there is no "wrong" interleaving because every consumer reads current state.
- **Positive:** the same code-shape covers the redelivery-after-success case (Kafka rebalance redeliver) — the consumer reads the already-mutated state and idempotently no-ops.
- **Positive:** cross-topic ordering can be ignored as a design constraint — consumers tolerate any interleaving.
- **Negative:** every cross-service event handler must be SELECT FOR UPDATE'd, which adds a row lock to every consumed event. At MVP scale this is fine; at high qps it could become a hotspot for very popular orderIds (mitigated because per-orderId qps is naturally low).
- **Negative:** late-event metrics (`order_consumer_late_event_total`) need to be watched in dashboards or the silent no-ops hide bugs.

## Cross-reference

- `events.payment.completed.failure_modes` and `events.order.cancelled.failure_modes` in `contracts.json` document the per-state branches.
- `infra-topology.md` calls out "state-driven; reads current state via FOR UPDATE" on every consumer arrow.
- The `order_consumer_late_event_total` metric (per-event-type, per-current-status) is in `observability-spec.md` and surfaces in the State Machine dashboard.
