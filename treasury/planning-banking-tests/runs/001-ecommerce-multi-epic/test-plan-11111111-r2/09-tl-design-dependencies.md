---
artifact_type: qa-tl-design-dependencies
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# TL Design Dependencies

## Compensation-hook specification (stock-restore semantics, OOB notification contr
- Description: Compensation-hook specification (stock-restore semantics, OOB notification contract, manual-refund operator endpoint shape) not yet committed by TL.
- Needed by: 

## State-machine model not yet committed by TL: allowed transition matrix (awaiting
- Description: State-machine model not yet committed by TL: allowed transition matrix (awaiting_payment, paid, packing, shipped, delivered, cancelled), terminal states, and compensation-hook spec required to finalise invalid-transition test matrix.
- Needed by: 

## TL has not specified concurrency_primitive for stock reservation: lock granulari
- Description: TL has not specified concurrency_primitive for stock reservation: lock granularity (per-SKU row lock vs reservation-table advisory lock vs optimistic CAS), isolation level, and timeout. Race-resolution tests TC-EPIC-CART-CHECKOUT-3-002 and TC-EPIC-CART-CHECKOUT-3-008 need this to wire deterministic harness probes.
- Needed by: 

## TL has not specified error_envelope shape for payment/checkout failures: error c
- Description: TL has not specified error_envelope shape for payment/checkout failures: error code namespace, reason field, message-vs-code split, structured details, localization key. Error tests TC-EPIC-CART-CHECKOUT-3-002, TC-EPIC-CART-CHECKOUT-3-009, TC-EPIC-CART-CHECKOUT-4-005 currently assert on prose strings ('out_of_stock_at_confirm', 'coupon_expired') which are brittle.
- Needed by: 

## TL has not specified idempotency_key_derivation_rule: which request fields (cust
- Description: TL has not specified idempotency_key_derivation_rule: which request fields (customer_id, cart_snapshot_hash, client-supplied nonce?) compose the key and how it is hashed/normalized. Without this, replay-safety tests TC-EPIC-CART-CHECKOUT-3-003 and webhook-idempotency TC-EPIC-CART-CHECKOUT-4-004 cannot deterministically assert key matching across replays.
- Needed by: 

