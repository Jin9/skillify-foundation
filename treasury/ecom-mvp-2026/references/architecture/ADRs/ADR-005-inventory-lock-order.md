# ADR-005 — Inventory lock-order pin (per-SKU, ascending sku string)

- **Status:** Accepted
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

`inventory.reservation.create` may reserve N items in a single call (e.g. a checkout cart with 3 different SKUs). Each item requires a `SELECT stock_levels FOR UPDATE WHERE sku = $1` row lock to enforce the no-negative invariant under concurrent reservations. If two concurrent reservations from two different checkouts touch the same two SKUs in opposite orders, Postgres can deadlock.

Without a documented lock-order rule, deadlock-frequency scales with concurrent checkout traffic, even at low qps if the catalog has hot SKUs (e.g. a flash-sale item).

## Decision

For any transaction that locks more than one `stock_levels` row, locks MUST be acquired in ASCENDING SKU STRING ORDER.

Implementation requirement on `inventory.reservation.create`:
```go
sort.Slice(req.Items, func(i, j int) bool {
    return req.Items[i].Sku < req.Items[j].Sku
})
for _, item := range req.Items {
    // SELECT stock_levels FOR UPDATE WHERE sku = item.Sku
    // ... validate available_qty >= item.Qty
    // ... UPDATE stock_levels SET available_qty -= item.Qty, reserved_qty += item.Qty
}
INSERT INTO reservations ...
COMMIT
```

The same rule applies anywhere a future caller bulk-locks multiple `stock_levels` rows. Single-row paths (`reservation.commit`, `reservation.release`, `reservation.sweep-expired` per-row, `stock.adjust`) cannot deadlock and don't need ordering.

## Consequences

- **Positive:** deadlock-free under arbitrary concurrent reservation patterns; Postgres' lock manager grants in a deterministic order.
- **Positive:** the rule is testable — a unit test can construct two concurrent reservation requests with SKUs in opposite order and assert both succeed sequentially without ERR-DEADLOCK.
- **Negative:** the sort step adds O(N log N) per request, which is irrelevant at MVP scale (N <= 50 items per cart).
- **Negative:** if a future caller bulk-mutates SKUs from a different ordering source (e.g. category-wide stock adjustment), it must adopt the same rule or risk deadlocking with reservations.
