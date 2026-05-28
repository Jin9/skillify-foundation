# PLAN_NOTES.md — Plan-Stage Advisories

Low-severity Plan-Reviewer findings (`plan_polish`) accumulate here per `docs/orchestrator.md` § Severity Policy. They never block; they're a backlog of polish items the producers can address opportunistically.

Source: `.squad-run/plan-review.json` (verdict: `block`; 14 findings, 5 advisory).

---

## PR-010 — payment.callback handler ordering as prose, not checklist

**Locus:** `contracts.json#/contracts[name=payment.callback].idempotency_rules and ordering_guarantees`

The three callback shapes (replay on `(intentId, providerStatus)`, terminal with different status, fresh callback) are functionally specified across `idempotency_rules` and notes, but never as a single ordered 1-2-3 checklist. Tech-Designer would benefit; not a correctness issue.

## PR-011 — ORD-008 order-number format in_scope but no AC

**Locus:** `ba.json#/in_scope (ORD-008)`; no AC pins the `ORD-YYYYMMDD-NNNNNN` regex.

Advisory for MVP. Matters only if external systems parse `orderNumber`.

## PR-012 — `events.product.created` is a TL addition beyond BA

**Locus:** `contracts.json#/contracts[name=events.product.created]`

TL added this to let Inventory seed `stock_levels(sku)` on new products. Structurally necessary (otherwise `inStock` filtering breaks for new SKUs); assumption is defensible. Captured here for traceability — not a finding to repair.

## PR-013 — `cart.add-item` same-SKU merge mechanism unspecified

**Locus:** `contracts.json#/contracts[name=cart.add-item].idempotency_rules`

`MERGE into one line` is the rule; the mechanism (`UNIQUE(cart_id, product_id)` + `INSERT…ON CONFLICT UPDATE qty=qty+excluded.qty` in one tx, or row-lock-and-update) isn't pinned. Tech-Designer can resolve; one-line lock prevents divergent implementations.

## PR-015 — BA cites PAY-004 in abandoned-checkout edge case (round-2 nit)

**Locus:** `ba.json#/edge_cases` new abandoned-checkout entry.

The new edge case the abandoned-checkout transition is driven by Inventory's reservation sweeper (PR-001), not by a customer-simulated payment timeout (PAY-004). Citation typo only — observable acceptance signal is correct and unambiguous. Polish for a future BA pass.

## PR-014 — `catalog.product.detail` "historical-order paths" wording is misleading

**Locus:** `contracts.json#/contracts[name=catalog.product.detail].failure_modes`

`order_items` persist `productNameSnapshot`, `productImageUrlSnapshot`, `priceSnapshot` per `cross-cutting.persistence`, so Order-detail does NOT call Catalog for DELETED products. The phrase "still resolvable for admin and historical-order paths" should be tightened to "admin-context only" to avoid steering Tech-Designer toward an unnecessary Order→Catalog cross-call.
