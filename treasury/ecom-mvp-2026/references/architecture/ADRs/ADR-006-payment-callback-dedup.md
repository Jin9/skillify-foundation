# ADR-006 — Payment callback dedup on (paymentIntentId, providerStatus)

- **Status:** Accepted (LOCKED)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

PAY-005 mandates that a duplicate payment-success callback for the same payment intent does not double-apply state or stock changes. In MVP the callback comes from the in-process mock; in v2 it would come from Stripe/Omise/2C2P, all of which redeliver on transient failures.

Two questions: (1) what is the dedup key? (2) what does the second-arrival do?

## Decision

- **Key:** `(paymentIntentId, providerStatus)` — composite. Persisted as `payment_callback_dedup(paymentIntentId, providerStatus)` UNIQUE.
- **Behavior:** the handler runs in a single Postgres transaction:
  ```
  BEGIN;
    INSERT INTO payment_callback_dedup (payment_intent_id, provider_status, request_hash, response_envelope)
    VALUES ($1, $2, $3, $4);   -- raises unique violation on replay
    UPDATE payment_records SET mock_payment_ref = $5, provider_status = $2, paid_at = $6 WHERE payment_intent_id = $1;
    INSERT INTO outbox_events (id, aggregate_id, event_type, payload_json) VALUES ($id, $1, 'payment.completed', ...);
  COMMIT;
  ```
- **Replay:** the handler catches the unique-constraint violation, looks up the cached `response_envelope`, and returns it verbatim. NO new outbox event. NO new state mutation.
- **Different providerStatus on the same intent:** the FIRST status persisted wins. A later callback with a CONTRADICTING status (e.g. SUCCEEDED then FAILED) is rejected with `409 PAYMENT_INTENT_TERMINAL` carrying the original status. This is what makes the dedup key composite — if it were just `paymentIntentId`, a late FAILED callback would silently win or lose ambiguously.
- **Amount-mismatch:** PAY-006 — if the callback's `amount` != `intent.amount`, the handler rejects with `409 PAYMENT_AMOUNT_MISMATCH` BEFORE inserting the dedup row, so the intent stays REQUIRES_PAYMENT and the reservation is unmolested.
- **Retention:** dedup rows live for >= 30 days (the redelivery horizon for typical providers). Cleanup is a future-iteration concern.

## Consequences

- **Positive:** redelivery is naturally absorbed by the unique-constraint violation; no application-level lock or Redis dedup needed.
- **Positive:** auditable — the dedup table records every callback that ever resolved an intent, including the cached envelope returned on replay.
- **Negative:** if the cached `response_envelope` is poisoned by a bad early callback (bug or attack), every replay returns the bad cached envelope; mitigation is the amount-mismatch guard runs BEFORE the dedup insert.
- **Negative:** the table grows monotonically until the future cleanup job lands; sized at ~1 row per resolved intent so even at 1M intents/year it's a tiny table.
