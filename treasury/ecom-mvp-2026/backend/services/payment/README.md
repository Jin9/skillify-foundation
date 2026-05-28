# Payment Service

Mock payment service for the B2C E-Commerce Platform. Owns the payment-intent state machine, HMAC-verified webhook callbacks, customer-facing simulation, and transactional outbox.

---

## Endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/api/v1/payment/intent/create` | `X-Internal-Secret` header | Create payment intent (called by checkout). Idempotent on `order_id`. |
| POST | `/api/v1/payment/intent/simulate` | JWT CUSTOMER | Simulate a payment outcome in-process (calls `process_callback` internally). |
| POST | `/api/v1/payment/intent/callback` | HMAC-SHA256 signature | Webhook from mock provider. Signature verified before any DB access. |

---

## Canonical `process_callback` — 7-Step Ordering

All state transitions for payment intents flow through a single internal function (`ProcessCallback` in `app/payment/process_callback.go`). Both the HTTP webhook handler and the simulate handler call it. The 7 steps, in order:

1. **Lock intent row.** `BEGIN TX` + `SELECT … FOR UPDATE` on `payment_intents.intent_id`. Returns `404 PAYMENT_INTENT_NOT_FOUND` if not found.

2. **Amount-mismatch check (PAY-006).** Compare `amountFromCaller` to `intent.amount_minor` (server truth). On mismatch: `ROLLBACK`, return `409 PAYMENT_AMOUNT_MISMATCH`. No dedup row written. No state transition.

3. **Derive dedup key.** `sha256(intent_id || '|' || provider_status)` — server-validated fields only. No user-controlled body content participates in the key.

4. **Dedup lookup.** `SELECT FROM payment_callback_dedup WHERE dedup_key = $1`. If found: `COMMIT` (no writes), return cached envelope. **Outbox is never re-fired on a replay hit.**

5. **Terminal-state guard.** If `intent.status IN (SUCCEEDED, FAILED, EXPIRED)` (same-status replay would have hit Step 4): return `409 PAYMENT_INTENT_TERMINAL`. Persist dedup row with the 409 envelope so retries return the same response without re-entering the guard.

6. **Apply transition.** `UPDATE payment_intents SET status=…, mock_provider_ref=…, paid_at=… WHERE intent_id=$1 AND version=$2`. Then `INSERT INTO payment.outbox` for `payment.completed` (SUCCEEDED) or `payment.failed` (FAILED) — in the **same transaction**. `EXPIRED` status does not emit an outbox event (see Sweeper Scope below).

7. **Persist success dedup row.** `INSERT INTO payment_callback_dedup … ON CONFLICT DO NOTHING`. `COMMIT`. Return success envelope.

### Ordering invariant: amount check BEFORE terminal check

Amount mismatch (Step 2) is checked before deriving the dedup key (Step 3) and before the terminal-state check (Step 5). This means a mismatched-amount replay of a terminal intent returns `409 PAYMENT_AMOUNT_MISMATCH`, not `409 PAYMENT_INTENT_TERMINAL`. The TD specifies this deliberately: amount errors are not persisted in the dedup table.

---

## Dedup Key Composition

```
dedup_key = hex(SHA-256(intent_id::text || '|' || provider_status))
```

**Fields used:**
- `intent_id` — from the looked-up `payment_intents` row (server-recorded UUID).
- `provider_status` — validated enum from the request (`SUCCEEDED | FAILED | EXPIRED`).

**Fields NOT used:**
- `amount` — never mixed into the key (amount errors are rejected before dedup).
- `mockPaymentRef` — user-supplied string; excluded per security posture.
- `providerTimestamp` — excluded; same (intent, outcome) pair must always dedup regardless of when it arrives.

---

## Signature Verification

Webhook authentication uses **HMAC-SHA256** over the **raw request body bytes** (not re-serialised JSON):

```
X-Mock-Provider-Signature: hex(HMAC-SHA256(rawBody, CALLBACK_HMAC_SECRET))
```

Implementation in `app/payment/signature.go`:
- Algorithm is **pinned to SHA-256** — no algorithm negotiation.
- Comparison uses `crypto/subtle.ConstantTimeCompare` — timing-safe; prevents timing-oracle attacks.
- Raw body is read **before** any JSON parsing so the MAC covers exactly what the provider signed.
- Timestamp skew check (`X-Mock-Provider-Timestamp`) rejects requests more than ±5 minutes old.
- On mismatch: `401 AUTH_INVALID`. Signature header and body are **never** logged.

---

## Internal Auth (`payment.intent.create`)

MVP uses a shared secret via `X-Internal-Secret` header (configured as `INTERNAL_SHARED_SECRET` env var). The checkout service includes this header on every call. Future: SERVICE-role JWT.

---

## Sweeper Scope — No Event Emission

The `ExpirySweeper` (`app/payment/expiry_sweeper.go`) runs every 60 seconds:

```sql
SELECT intent_id FROM payment.payment_intents
WHERE status = 'REQUIRES_PAYMENT' AND expires_at < now()
ORDER BY expires_at
LIMIT $batch_size
FOR UPDATE SKIP LOCKED
```

For each expired row, it calls `MarkExpired` in its own short transaction.

**The sweeper does NOT emit `events.payment.expired` to the outbox.**

Rationale (TL PR-001 LOCKED): The reservation sweeper (inventory service) publishes `events.reservation.expired`, which Order consumes to transition `PENDING_PAYMENT → PAYMENT_EXPIRED` and Inventory consumes to release stock. Emitting `events.payment.expired` in addition would create duplicate-cancel cascades. Payment's sweeper only makes Payment's own state honest so that subsequent late callbacks for an EXPIRED intent cleanly return `409 PAYMENT_INTENT_TERMINAL`.

Future enablement: set env `PAYMENT_EMIT_EXPIRED_EVENT=true` (default `false`).

---

## Configuration

| Env Var | Default | Purpose |
|---------|---------|---------|
| `PORT` | — | HTTP listen port |
| `DB_URL` | — | Postgres connection URL |
| `JWT_PUBLIC_KEY` | — | EC public key (PEM) for verifying customer JWTs |
| `JWT_ISSUER` | `shoppilot-identity` | Expected JWT issuer |
| `JWT_AUDIENCE` | `shoppilot-api` | Expected JWT audience |
| `CALLBACK_HMAC_SECRET` | — | Shared secret for HMAC-SHA256 webhook verification |
| `INTERNAL_SHARED_SECRET` | — | Shared secret for checkout → payment internal calls |
| `INTENT_TTL_MINUTES` | `15` | Payment intent lifetime (must match reservation TTL) |
| `SWEEPER_CADENCE_SECONDS` | `60` | How often the expiry sweeper ticks |
| `SWEEPER_BATCH_SIZE` | `200` | Max intents expired per sweeper tick |
| `KAFKA_ENABLED` | `false` | Enable Kafka outbox relay |
| `KAFKA_BROKERS` | — | Kafka broker addresses |
| `PAYMENT_EMIT_EXPIRED_EVENT` | `false` | Future: emit events.payment.expired (see Sweeper Scope) |

---

## Logging Redaction

Per TD security posture:
- **Always logged:** `dedup_key`, `intent_id`, `order_id`, `status`, `action` (REPLAY_HIT / APPLIED / AMOUNT_MISMATCH / TERMINAL), `traceId`.
- **Never logged:** full callback request body, `mockPaymentRef`, signature header value, `CALLBACK_HMAC_SECRET`, or any field that could carry provider PII.

---

## Outbox Events

| Event | Trigger | Topic | Partition Key |
|-------|---------|-------|---------------|
| `payment.completed` | `process_callback(SUCCEEDED)` | `ecom.payment.events` | `order_id` |
| `payment.failed` | `process_callback(FAILED)` | `ecom.payment.events` | `order_id` |
| `payment.expired` | NOT emitted in MVP | — | — |

Partition key = `order_id` ensures strict per-order ordering within `ecom.payment.events`.
