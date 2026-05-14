---
artifact_type: qa-test-data
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Test Data Design

## FIX-IDENTITY-customer_id

- PII field refs: customer_id
- Synthesis rules: Surrogate customer_id format CUST-{8-hex-digits-from-seed}. Used in audit events as the only PII-adjacent identifier (per PII inventory rule that customer_id is the audit surrogate).

## FIX-IDENTITY-customer_id-2

- PII field refs: customer_id
- Synthesis rules: Second surrogate customer_id for cross-customer authz scenarios. Distinct deterministic 8-hex suffix from the same frozen seed.

## FIX-IDENTITY-email

- PII field refs: email
- Synthesis rules: Synthetic email pattern: customer+{seq}@shoppilot.test where {seq} is a zero-padded 6-digit counter from a frozen seed. Domain shoppilot.test guarantees non-routable. Never reuse real customer emails (AP-Q7). Bilingual signup glossary preserved for any display-side rendering.
- Time anchors: fixture-seed-2026-05-12T00:00:00Z

## FIX-IDENTITY-email-unknown

- PII field refs: email
- Synthesis rules: Synthetic non-existent email pattern: unknown+{seq}@shoppilot.test guaranteed to not exist in customer table. Used to probe unknown-email branch of login without leaking real PII.

## FIX-IDENTITY-name

- PII field refs: name
- Synthesis rules: Bilingual name fixtures drawn from BA glossary: Thai surface form 'นิรันดร์' / Latin transliteration 'Niran' preserved together. Counter suffix appended for uniqueness ('Niran-{seq}'). Both surface forms travel through display, audit-scrub, and DB layers.

## FIX-IDENTITY-order_number

- PII field refs: order_number
- Synthesis rules: Order-number format ORD-YYYYMMDD-NNNNNN with date pinned to the fixture time anchor (2026-05-12) and NNNNNN drawn from the frozen sequence. Suffix predictability concern parked under OQ-05 (handled by EPIC-CART-CHECKOUT). Used here only to wire the address-referenced-by-active-order scenario.

## FIX-IDENTITY-password

- PII field refs: password
- Synthesis rules: Password fixtures stored as pre-computed one-way hash references (e.g., FIX-IDENTITY-pw-hash-A, FIX-IDENTITY-pw-hash-B); plaintext is generated ephemerally inside the test process only and is NEVER persisted to fixture files, logs, or screenshots. Plaintext is also never echoed in assertion messages. (AP-Q7, AP-Q12.)

## FIX-IDENTITY-password-weak

- PII field refs: password
- Synthesis rules: Pre-canned weak-password literals ('123', 'aaaa', '') held only in the security-test fixture, never in shared fixtures. Asserts the strength predicate rejects them. Exact predicate TBD pending OQ-03.

## FIX-IDENTITY-phone

- PII field refs: phone
- Synthesis rules: Thai mobile format 0812345{NNN} where {NNN} is a 3-digit counter from the frozen seed. Stays inside the reserved test prefix; no real Thai subscriber number reused (AP-Q7).

## FIX-IDENTITY-review_text

- PII field refs: review_text
- Synthesis rules: Synthetic review strings drawn from a fixed pool of bilingual user-generated content templates; intentionally contain NO real PII. Used to validate that review-text PII-leakage scenarios surface a moderation flag rather than leaking real data.

## FIX-IDENTITY-session_token

- PII field refs: session_token
- Synthesis rules: Opaque base64 token rendered from frozen seed: base64(sha256('FIX-IDENTITY-session_token::' || seq)) truncated to 32 bytes. Token value never appears in assertion text, never logged, never asserted-equal in plaintext — assertions check 'present' or 'cleared' only.

## FIX-IDENTITY-shipping_address

- PII field refs: shipping_address
- Synthesis rules: Thai postal address structure: street='Test Soi {seq}', subdistrict='Khlong Toei Nuea', district='Watthana', province='Bangkok', postcode='10110'. All values from a frozen non-routable test catalog of Bangkok postal codes; never a real customer's address (AP-Q7). Scrubbed from authz_denied audit payloads per PII inventory rule.

## FIX-IDENTITY-shipping_address-2

- PII field refs: shipping_address
- Synthesis rules: Second synthetic Thai postal address for default-swap scenarios. postcode='10330' (Pathum Wan). Same fabrication rules as FIX-IDENTITY-shipping_address.

## FIX-IDENTITY-tracking_number

- PII field refs: tracking_number
- Synthesis rules: Synthetic tracking-number pattern TRK-{12-alnum-from-seed}; not present on any real carrier. Referenced from IDENTITY fixtures only for cross-customer read-attempt negative tests (no positive use in EPIC-IDENTITY).

## FX-AUDIT-sink

- PII field refs: _(none)_
- Synthesis rules: In-memory append-only audit sink with assert helpers for field-by-field event inspection and allow-list scans for forbidden fields (password, token, PAN, CVV).
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-active-product

- PII field refs: customer_id
- Synthesis rules: Generate P1 with status=active, price=300 THB, stock_available=10; customer C1 with synthetic email animate90+c1@example.test (no real PII).
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-customer-empty

- PII field refs: customer_id
- Synthesis rules: Customer C1 authenticated, cart empty, no prior orders.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-customer-with-line

- PII field refs: customer_id
- Synthesis rules: Customer C1 with existing cart line P1 qty=2 and P1 stock_available=10.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-idem-key-001

- PII field refs: _(none)_
- Synthesis rules: Idempotency key 'ik-cart-001' replayed twice on add-to-cart request.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-low-stock-product

- PII field refs: customer_id
- Synthesis rules: Product P1 with stock_available=5, status=active.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-paused-product-mid-session

- PII field refs: customer_id
- Synthesis rules: Product P1 starts active with cart line; admin action flips P1.status to paused at T0+30s. Frozen clock.
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0+30s

## FX-CART-subtotal-4000

- PII field refs: customer_id
- Synthesis rules: Cart configured so subtotal == 4000 exactly.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-subtotal-600

- PII field refs: customer_id
- Synthesis rules: Cart configured so subtotal == 600 exactly (e.g. P1 price=300 qty=2).
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-subtotal-80

- PII field refs: customer_id
- Synthesis rules: Cart with subtotal == 80 to exercise discount clamp.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CART-two-customers

- PII field refs: customer_id
- Synthesis rules: Customers C1 and C2 with separate carts and addresses; never share session.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-address-C1

- PII field refs: customer_id, shipping_address, phone
- Synthesis rules: Synthetic address owned by C1 with non-real phone and address strings; address_id='addr-c1-001'.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-address-C2

- PII field refs: customer_id, shipping_address
- Synthesis rules: Synthetic address owned by C2; used to attempt cross-owner reference from C1.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-bulk-100-orders

- PII field refs: customer_id
- Synthesis rules: 100 sequential checkout submissions to inspect order_number generation distribution.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-cart-P1q2

- PII field refs: customer_id, shipping_address
- Synthesis rules: Customer C1 with valid cart of P1 qty=2, address A1 owned by C1, optional coupon WELCOME100.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-cart-with-stale-coupon

- PII field refs: customer_id
- Synthesis rules: Cart with EXPIRED50 coupon applied at preview; clock advanced past coupon expiry before submit.
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0+1h

## FX-CHK-customer-C1

- PII field refs: customer_id
- Synthesis rules: Authenticated session for customer C1 with bearer token.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-five-customers

- PII field refs: customer_id
- Synthesis rules: Five distinct synthetic customers C1..C5, each with P1 qty=1 in cart.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-last-unit-P1

- PII field refs: _(none)_
- Synthesis rules: Product P1 with stock_available=1 to drive race-resolution test.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-stock-3

- PII field refs: _(none)_
- Synthesis rules: Product P1 stock_available=3.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CHK-two-customers

- PII field refs: customer_id
- Synthesis rules: Customers C1 and C2 with identical cart state {P1 qty=1}.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CLOCK-control

- PII field refs: _(none)_
- Synthesis rules: Programmatic clock-advance API used by sweep tests; supports advance(seconds) and reset().
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-CONCURRENCY-harness

- PII field refs: _(none)_
- Synthesis rules: Deterministic concurrency harness with controllable thread/coroutine scheduling and lock instrumentation; supports barrier-synced parallel submits.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-COUPON-EXPIRED50

- PII field refs: _(none)_
- Synthesis rules: Coupon 'EXPIRED50' with valid_until=T0-1h (already expired at submit).
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0-1h

## FX-COUPON-SAVE10

- PII field refs: _(none)_
- Synthesis rules: Coupon code 'SAVE10', percent=10, cap=300, min-order=0.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-COUPON-WELCOME100

- PII field refs: _(none)_
- Synthesis rules: Coupon code 'WELCOME100', fixed-THB=100, min-order=500, quota_remaining=1000, valid_from=T0-1d, valid_until=T0+30d.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-COUPON-WELCOME100-lowered-min

- PII field refs: _(none)_
- Synthesis rules: Variant of WELCOME100 with min-order=50 to test clamp on subtotal=80.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-IDEM-key-chk-001

- PII field refs: _(none)_
- Synthesis rules: Idempotency key 'ik-chk-001' derived from sha256(customer_id+cart_snapshot_hash+frozen_seed='SEED-2026-05-12'); reused across two submit calls.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-IDEM-key-coup-001

- PII field refs: _(none)_
- Synthesis rules: Idempotency key 'ik-coup-001' derived from sha256(customer_id+coupon_code+frozen_seed='SEED-2026-05-12').
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-ORDER-awaiting-payment-560

- PII field refs: customer_id
- Synthesis rules: Order O1 with grand_total=560, status='awaiting-payment'.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-ORDER-awaiting-payment-O1

- PII field refs: customer_id
- Synthesis rules: Order O1 in status='awaiting-payment' with P1 stock_reserved=2.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-ORDER-cancelled-pre-payment

- PII field refs: customer_id
- Synthesis rules: Order with WELCOME100 coupon applied, then cancelled before payment to drive quota-reversal test.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-ORDER-paid-O1

- PII field refs: customer_id
- Synthesis rules: Order O1 already transitioned to 'paid' via prior webhook for duplicate-callback test.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-PSP-evt-001

- PII field refs: _(none)_
- Synthesis rules: Replay of provider_event_id='evt-001' for webhook idempotency test.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-PSP-mock-all-outcomes

- PII field refs: _(none)_
- Synthesis rules: Three sibling payloads covering success/failure/timeout for parameterized audit-event test.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-PSP-mock-failure

- PII field refs: _(none)_
- Synthesis rules: Mock PSP webhook with outcome='failure', reason_code='declined'.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-PSP-mock-no-signal

- PII field refs: _(none)_
- Synthesis rules: Mock PSP that never delivers callback; used with FX-CLOCK-control to drive timeout sweep.
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0+10min

## FX-PSP-mock-success

- PII field refs: _(none)_
- Synthesis rules: Mock PSP webhook payload with outcome='success', amount matches order, provider_event_id='evt-001'. Synthetic credit-card payload (card_holder='TEST USER', card_last4='4242', no PAN/CVV captured by app per PCI-out-of-scope rule).
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-PSP-mock-success-wrong-amount

- PII field refs: _(none)_
- Synthesis rules: Mock PSP webhook with outcome='success' but amount=500 vs order grand_total=560.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-PSP-rogue-signature

- PII field refs: _(none)_
- Synthesis rules: Webhook call from a non-configured origin or with invalid signature header.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-TIME-T0

- PII field refs: _(none)_
- Synthesis rules: Static T0 anchor for non-time-dependent tests.
- Time anchors: T0=2026-05-12T09:00:00+07:00

## FX-TIME-T0-T+10min

- PII field refs: _(none)_
- Synthesis rules: Two-anchor frozen clock; T+10min crosses reservation TTL (assumed; TBD pending OQ-CART-reservation-ttl).
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0+10min

## FX-TIME-T0-T+5min

- PII field refs: _(none)_
- Synthesis rules: Two-anchor frozen clock for short-window failure tests.
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0+5min

## FX-TIME-anchor-T0-T+1h

- PII field refs: _(none)_
- Synthesis rules: Frozen clock harness; T0 fixed; advance to T0+1h for expiry tests.
- Time anchors: T0=2026-05-12T09:00:00+07:00, T0+1h

## catalog-admin-fixture-set

- PII field refs: _(none)_
- Synthesis rules: Admin actor records (role=admin, username='aoy' and a second admin for concurrency tests) used as audit-event actors only — no customer PII. Pre-existing product with SKU 'SKU-EL-001', price=1290, status=active, category=Electronics, stock=10 to support duplicate-SKU rejection, price-change audit, SKU-immutability, and idempotent-replay scenarios. Pre-issued idempotency_key 'ik-prod-001' bound to one prior successful edit. Reason-for-change free-text seed values for audit payload assertions.
- Time anchors: T0, T+1m

## catalog-fixture-set

- PII field refs: _(none)_
- Synthesis rules: 20 seed products: 12 with status=active and stock_available>0 (10 in-stock + 2 will be toggled out-of-stock for the in-stock-only filter test), 8 with status in {draft, paused, deleted}. Distribute the 12 active products across 5 categories with Electronics having at least 3 SKUs. Include bilingual product names (Thai + English) sourced from the BA glossary surface forms. Each product carries: SKU (immutable), slug, name, description, image, price (>0), compare-at price (nullable), category_id, stock_available, status, public/hidden flag. Includes a soft-deleted product P3 reachable by deep link, an active product P1 (price=1290, stock=5, rating=4.2, review_count=11), and a paused product P2.
- Time anchors: T0

## fx-catalog-mutated

- PII field refs: _(none)_
- Synthesis rules: Catalog state where P1.price was 1290 at order time then mutated to 999. Used to verify snapshot immutability per 6.8.
- Time anchors: frozen-order-ts=2025-11-01T09:05:00+07:00, frozen-mutation-ts=2025-11-03T10:00:00+07:00

## fx-gov-admin-aoy

- PII field refs: name, email
- Synthesis rules: Synthetic admin identity 'admin:aoy'; email pattern aoy+admin@example.test; no real PII (AP-Q7). Role=admin.
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-audit-corpus-mixed

- PII field refs: customer_id, order_number, shipping_address
- Synthesis rules: Synth corpus of mixed audit rows covering catalog/order/identity/payment/fulfillment events. Marker tokens injected for scrubbing assertions; never any real PII (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00, T0-30d, T0-90d, T0-91d

## fx-gov-audit-row-price-changed

- PII field refs: customer_id
- Synthesis rules: Seeded audit row 'product.price_changed' with before={price:1290}, after={price:999}. No customer PII.

## fx-gov-cross-epic-event-corpus

- PII field refs: customer_id, order_number
- Synthesis rules: Synth corpus exercising one instance of each enumerated state-changing event class across all epics (catalog.*, order.*, inventory.adjusted, review.created, review.hidden, identity.*, payment.*, fulfillment.*).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-customer-c1

- PII field refs: email, phone, name, shipping_address, customer_id
- Synthesis rules: Synthetic customer C1 with synth email c1@example.test, synth phone +66-XX-XXX-XXX0 (non-routable test range), name and Thai-format address generated from a closed test gazetteer; never use real customer data (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-customer-c2-not-owner

- PII field refs: email, customer_id
- Synthesis rules: Second synthetic customer used for authz negative tests; same synthesis rules as C1 (AP-Q7).

## fx-gov-notification-corpus

- PII field refs: email, phone, customer_id, order_number
- Synthesis rules: Synth corpus of customer-facing notifications: order-status updates, signup confirm, password reset, review-status. Synth recipients only (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-order-o1-already-reviewed

- PII field refs: order_number, customer_id, review_text
- Synthesis rules: Synthetic order O1 delivered with one pre-existing review row for the P1 line (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-order-o1-delivered

- PII field refs: shipping_address, order_number, tracking_number, customer_id
- Synthesis rules: Synthetic order O1 in status='delivered' for C1, one line for P1. Address and order_number synthesized; tracking_number from synth carrier ID space (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00, delivered_at=T0-7d

## fx-gov-order-o1-paid-not-delivered

- PII field refs: order_number, customer_id
- Synthesis rules: Synthetic order O1 in status='paid' (not delivered) for C1 (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-review-r-flagged

- PII field refs: review_text, customer_id
- Synthesis rules: Synthetic review whose content is configured to trigger a specific moderation rule_id for tipping-off-copy tests (AP-Q7, AP-Q11).

## fx-gov-review-r1

- PII field refs: review_text, customer_id
- Synthesis rules: Review row R1 with synth text in Thai; no incidental PII; tied to C1+O1+P1 (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-gov-sku-p1

- PII field refs: _(none)_
- Synthesis rules: Catalog SKU P1 seeded with stock_available=5, stock_reserved=0, stock_sold=0. No PII.

## fx-gov-sku-p1-low-stock

- PII field refs: _(none)_
- Synthesis rules: SKU P1 seeded with stock_available=1, stock_reserved=0. No PII.

## fx-gov-sku-p1-with-reservations

- PII field refs: _(none)_
- Synthesis rules: SKU P1 seeded with stock_available=2, stock_reserved=3 and 3 live reservations. No PII.

## fx-gov-synth-pii-corpus

- PII field refs: email, phone, name, shipping_address, password, session_token, customer_id
- Synthesis rules: Marker-tokenized synthetic PII corpus driven through the system to produce realistic log volume. Markers are well-known sentinel strings ('SENTINEL_PWD', 'SENTINEL_TOKEN', etc.) so post-hoc log scans can deterministically count leaks. No real PII (AP-Q7).
- Time anchors: T0=2026-05-12T00:00:00+07:00

## fx-idem-keys

- PII field refs: _(none)_
- Synthesis rules: Pre-seeded idempotency keys: ik-fulfil-001 (paid->packing), ik-cancel-001 (admin-cancel). Replay tests reuse exact keys to verify dedupe.

## fx-orders-state-frozen

- PII field refs: email, phone, name, shipping_address
- Synthesis rules: Frozen-state order fixtures at each state checkpoint: O-awaiting (awaiting_payment), O-paid (paid), O-packing (packing), O-shipped (shipped, tracking=TH123456789), O-delivered (delivered, terminal), O-cancelled (cancelled from paid). PII is synthetic-fake: deterministic seed=ecommerce-v5-fulfil; tipping-off scan clean. Order numbers conform to ORD-YYYYMMDD-NNNNNN with frozen date 2025-11-01.
- Time anchors: frozen-now=2025-11-01T09:00:00+07:00, frozen-paid-ts=2025-11-01T09:05:00+07:00, frozen-packing-ts=2025-11-01T10:00:00+07:00, frozen-shipped-ts=2025-11-02T14:00:00+07:00, frozen-delivered-ts=2025-11-04T11:30:00+07:00

## fx-payment-ledger

- PII field refs: _(none)_
- Synthesis rules: Payment-captured marker rows for O-paid and O-cancelled orders, with no auto-refund record. Used by OOB-refund tests to assert compensating-action workflow.
- Time anchors: frozen-paid-ts=2025-11-01T09:05:00+07:00

## fx-stock-ledger

- PII field refs: _(none)_
- Synthesis rules: Stock ledger snapshot: P1 stock_available=8, stock_sold=2 prior to cancel; expected post-cancel: stock_available=10, stock_sold=0. Idempotent fixture; resets per test run.
- Time anchors: frozen-now=2025-11-01T09:00:00+07:00

## fx-tracking-numbers

- PII field refs: tracking_number
- Synthesis rules: Deterministic synthetic Thailand-Post-format tracking numbers (TH + 9 digits). Seeded so TC-EPIC-ORDER-FULFILL-1-003/-005 always observe 'TH123456789'. Indirect identifier - flagged for tipping-off scan.
- Time anchors: frozen-now=2025-11-01T09:00:00+07:00

## fx-users-roles

- PII field refs: email
- Synthesis rules: Role fixtures: admin:aoy (role=admin, tenant=T1), customer:C1 (role=customer, owns O1, tenant=T1), customer:C2 (role=customer, tenant=T1, no orders). Used for authz tests.

