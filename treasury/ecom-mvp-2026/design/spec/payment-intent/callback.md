# POST /api/v1/payment/intent/callback

## Summary

Webhook entry point that the (mock) payment provider invokes to report the outcome of a payment intent. Authenticated via HMAC-SHA256 signature over the raw body using `MOCK_PROVIDER_SECRET`, with timestamp-based replay protection (+/-5min skew). After validating the signature and timestamp, delegates to the canonical `process_callback()` 7-step handler — the single function shared with `payment.simulate` so both paths get identical idempotency, dedup, and outbox semantics.

The callback's two non-negotiable security properties are: (1) the dedup key is composed of server-recorded fields ONLY (`sha256(intent_id || '|' || provider_status)` per ADR-006), and (2) the amount-mismatch check runs BEFORE the terminal-state guard so an attacker cannot bypass amount validation by replaying a (paymentIntentId, providerStatus) pair against a terminal intent with a tweaked amount.

## Story refs

- `STORY_PAYMENT_SIMULATE_SUCCESS` (idempotency AC: PAY-005)
- `STORY_PAYMENT_AMOUNT_VALIDATION` (PAY-006)
- `STORY_PAYMENT_SIMULATE_FAILED_OR_TIMEOUT` (PAY-003 + PAY-004)

## Contract ref

[`payment.callback`](../../architecture/contracts.json#payment.callback)

## Auth

- Tier: **internal_secret** (webhook signature, not customer JWT, not X-Internal-Secret).
- Header `X-Mock-Provider-Signature: <hex(HMAC-SHA256(rawBody, MOCK_PROVIDER_SECRET))>` REQUIRED.
- Header `X-Mock-Provider-Timestamp: <RFC3339>` REQUIRED; reject if skew > +/-5min from `now()`.
- Validation algorithm:
  1. Read `X-Mock-Provider-Timestamp`; parse RFC3339; reject 401 `AUTH_INVALID` if unparseable or out-of-skew.
  2. Read `rawBody = c.GetRawData()` BEFORE `c.ShouldBindJSON` (the HMAC must hash exactly the bytes the provider signed; re-marshaling the parsed JSON does not produce byte-identical output).
  3. Compute `expected = hmac.New(sha256, MOCK_PROVIDER_SECRET).Sum(rawBody)` -> `hex(expected)`.
  4. Read `X-Mock-Provider-Signature`; hex-decode; reject 401 if malformed.
  5. Compare via `crypto/subtle.ConstantTimeCompare(expectedBytes, suppliedBytes) != 1` -> reject 401. Plain `==` / `!=` is FORBIDDEN.
  6. Algorithm SHA256 is PINNED in code — no caller-controlled `algo` field in the header (defends against CVE-2018-1000613 class downgrade attacks).
- `MOCK_PROVIDER_SECRET` MIN 32 bytes; recipient validates length at startup; refuses to boot if shorter.
- On reject: log structured `webhook_signature_reject` with `remoteAddr`, `path`, `traceId`; NEVER log the supplied signature, `rawBody`, `MOCK_PROVIDER_SECRET`, or any portion of the body.

## Request

Headers:
- `Content-Type: application/json`
- `X-Mock-Provider-Signature: <hex hmac>` (REQUIRED)
- `X-Mock-Provider-Timestamp: <RFC3339>` (REQUIRED; +/-5min skew)

Body:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["paymentIntentId", "providerStatus", "mockPaymentRef", "amount", "providerTimestamp"],
  "properties": {
    "paymentIntentId":   { "type": "string", "format": "uuid" },
    "providerStatus":    { "type": "string", "enum": ["SUCCEEDED", "FAILED", "EXPIRED"] },
    "mockPaymentRef":    { "type": "string", "maxLength": 64 },
    "amount":            { "type": "integer", "minimum": 0 },
    "providerTimestamp": { "type": "string", "format": "date-time" }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "callback applied",
  "data": {
    "paymentIntentId": "01935b9c-2a17-7c3d-9e8f-a1b2c3d4e5f6",
    "applied": true,
    "providerStatus": "SUCCEEDED"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | schema violation (missing field, malformed UUID, bad enum) |
| `AUTH_INVALID` | 401 | missing/wrong `X-Mock-Provider-Signature` (constant-time compare reject) OR `X-Mock-Provider-Timestamp` out-of-skew |
| `PAYMENT_INTENT_NOT_FOUND` | 404 | no `payment_intents` row matches `paymentIntentId` |
| `PAYMENT_AMOUNT_MISMATCH` | 409 | `callback.amount != intent.amount_minor` (PAY-006); intent stays REQUIRES_PAYMENT, NO dedup row written, NO outbox row written |
| `PAYMENT_INTENT_TERMINAL` | 409 | intent already TERMINAL AND incoming `providerStatus` differs from the original (returned by `process_callback` Step 5 terminal-state guard); body carries `{currentStatus, attemptedStatus}` |
| `INTERNAL_ERROR` | 500 | DB failure or panic |

## Business logic steps — the canonical 7-step ordering

This is the load-bearing security contract for the callback path. The function is named `process_callback` and lives in `app/payment/process_callback.go`. It is invoked by both this HTTP handler and by `payment.simulate` so both paths share the SAME idempotency, dedup, and outbox semantics.

```
function process_callback(intent_id, provider_status, mock_provider_ref, amount_from_caller, provider_timestamp):
    BEGIN TRANSACTION (single tx)

    // STEP 1 — LOCK INTENT
    intent := SELECT * FROM payment_intents WHERE intent_id = :intent_id FOR UPDATE
    IF intent IS NULL THEN
        ROLLBACK
        RETURN 404 PAYMENT_INTENT_NOT_FOUND

    // STEP 2 — AMOUNT-CHECK (BEFORE terminal-state guard)
    // SECURITY: This ordering matters. Placing AMOUNT before TERMINAL prevents an
    // attacker from bypassing amount validation by replaying a (intentId, providerStatus)
    // pair on a terminal intent with a tweaked amount. Without this ordering, the terminal-state
    // path would write a dedup row keyed on the contradicting status, masking the amount tampering.
    IF amount_from_caller != intent.amount_minor THEN
        ROLLBACK
        log.warn(intent_id, action='AMOUNT_MISMATCH', expected=intent.amount_minor, got=<REDACTED>)
        RETURN 409 PAYMENT_AMOUNT_MISMATCH (NO dedup row written, NO outbox row written)

    // STEP 3 — DEDUP-KEY DERIVE
    // SECURITY: server-recorded fields ONLY. NO user-controlled body fields participate.
    dedup_key := sha256(intent.intent_id::text || '|' || provider_status)

    // STEP 4 — REPLAY: RETURN CACHED ENVELOPE
    existing := SELECT * FROM payment_callback_dedup WHERE dedup_key = :dedup_key
    IF existing IS NOT NULL THEN
        COMMIT (no writes)
        log.info(dedup_key, intent_id, action='REPLAY_HIT')
        RETURN existing.envelope WITH existing.http_status   // NEVER re-emits outbox

    // STEP 5 — TERMINAL-STATE GUARD
    // Reaching here means: same intent already TERMINAL AND incoming status DIFFERS from original
    // (because Step 4 would have caught a same-status replay).
    IF intent.status IN ('SUCCEEDED', 'FAILED', 'EXPIRED') THEN
        envelope := build_error_envelope(409, 'PAYMENT_INTENT_TERMINAL',
            {currentStatus: intent.status, attemptedStatus: provider_status})
        // Persist the conflict envelope so future replays of THIS conflict get the same answer
        INSERT INTO payment_callback_dedup (dedup_key, intent_id, provider_status, envelope, http_status, expires_at)
            VALUES (:dedup_key, intent_id, provider_status, envelope, 409, now() + 30d)
        COMMIT
        log.info(dedup_key, intent_id, action='TERMINAL', currentStatus=intent.status, attemptedStatus=provider_status)
        RETURN envelope

    // STEP 6 — TRANSITION + OUTBOX in SAME TX (intent.status MUST be REQUIRES_PAYMENT here)
    new_status := MAP provider_status: SUCCEEDED -> SUCCEEDED, FAILED -> FAILED, EXPIRED -> EXPIRED
    paid_at    := (provider_status == SUCCEEDED) ? provider_timestamp : NULL

    UPDATE payment_intents SET
        status            = new_status,
        mock_provider_ref = :mock_provider_ref,
        provider_status   = :provider_status,
        paid_at           = :paid_at,
        updated_at        = now(),
        version           = version + 1
    WHERE intent_id = :intent_id
      AND version   = intent.version    // OPTLOCK; should always succeed under FOR UPDATE

    // ADR-004 outbox: INSERT in same tx; relay drains to Kafka asynchronously
    IF new_status == 'SUCCEEDED' THEN
        INSERT INTO outbox_events (
            id, aggregate_type, aggregate_id, event_id, event_type, topic, partition_key,
            payload, status, created_at)
        VALUES (
            uuidv7(), 'payment_intent', intent.intent_id, uuidv7(), 'payment.completed',
            'ecom.payment.events', intent.order_id::text,
            jsonb_build_object(
                'eventId', :event_id,
                'occurredAt', now(),
                'orderId', intent.order_id,
                'paymentIntentId', intent.intent_id,
                'amount', intent.amount_minor,
                'mockPaymentRef', :mock_provider_ref),
            'PENDING', now())
    ELIF new_status == 'FAILED' THEN
        INSERT INTO outbox_events (... event_type='payment.failed',
            payload={eventId, occurredAt, orderId, paymentIntentId, reason: 'PROVIDER_DECLINED'},
            status='PENDING')
    ELIF new_status == 'EXPIRED' THEN
        // MVP: NO outbox row. Inventory's reservation sweeper drives Order/Inventory transitions
        // via events.reservation.expired (ADR-009 + TL PR-001 LOCKED).
        // Feature flag PAYMENT_EMIT_EXPIRED_EVENT (default false) gates future enablement.
        skip

    // STEP 7 — DEDUP-SUCCESS ROW + COMMIT
    envelope := build_success_envelope(200, 'SUCCESS',
        {paymentIntentId: intent.intent_id, applied: true, providerStatus: provider_status})

    INSERT INTO payment_callback_dedup (dedup_key, intent_id, provider_status, envelope, http_status, expires_at)
        VALUES (:dedup_key, intent.intent_id, provider_status, envelope, 200, now() + 30d)
        ON CONFLICT (dedup_key) DO NOTHING
        // Defense-in-depth: a concurrent racer that bypassed FOR UPDATE (cross-replica split-brain)
        // and beat us past Step 4 will own the canonical envelope; we no-op.

    COMMIT
    log.info(dedup_key, intent_id, action='APPLIED', new_status, durationMs)
    RETURN envelope
```

### Why this exact ordering

1. **`FOR UPDATE` first** — serializes concurrent callbacks for the same intent. Race losers wait at Step 1 and find the world post-mutation when they run.
2. **AMOUNT before TERMINAL** — without this, an attacker on a terminal intent could replay with a tweaked amount and the terminal-state path would write a dedup row that masks the tampering. Putting amount-check first means a tampered amount is rejected with 409 `PAYMENT_AMOUNT_MISMATCH` and NO dedup row is written, leaving the intent's existing dedup envelope authoritative.
3. **`dedup_key` from server-trusted fields ONLY** — `sha256(intent_id || '|' || provider_status)`. If `mockPaymentRef` or `amount` were in the key, an attacker could force dedup misses by changing those fields and replay state mutations.
4. **REPLAY before TERMINAL guard** — same-status replays (the common provider-redelivery case) get the cached SUCCESS envelope, not a misleading 409.
5. **TERMINAL guard** — same-intent-different-status: cache the 409 envelope so the next time the provider replays the SAME conflicting status, it gets the SAME 409 (no oscillation if a tweaked amount changes the body again).
6. **TRANSITION + OUTBOX in SAME tx** — ADR-004 outbox-pattern guarantee: no phantom event without state change, no missing event with state change. EXPIRED skipped per ADR-009 + TL PR-001.
7. **DEDUP-SUCCESS row inside the tx** — a crash between TRANSITION and DEDUP rolls back BOTH; on retry the same callback re-runs cleanly through Step 4. `ON CONFLICT DO NOTHING` is defense for the cross-replica racer.

## Side effects

- On dedup-clear (Step 6 + 7): UPDATE `payment_intents` (status, mock_provider_ref, provider_status, paid_at, version+1), INSERT `outbox_events` row (SUCCEEDED -> `payment.completed`; FAILED -> `payment.failed`; EXPIRED -> NO event), INSERT `payment_callback_dedup` row — all in ONE Postgres tx.
- On replay (Step 4): zero writes; cached envelope returned.
- On terminal-different (Step 5): one INSERT into `payment_callback_dedup` (the conflict envelope); zero changes to `payment_intents`; zero outbox writes.
- On amount-mismatch (Step 2): zero writes — clean rollback.
- The `outbox_events` row is drained to Kafka by the in-process `outbox_relay` goroutine (separate from this handler's tx). Per-`orderId` ordering is preserved by `partition_key = order_id`.

## Idempotency

- Key shape: `dedup_key = sha256(intent_id || '|' || provider_status)`. Composite, server-derived.
- Scope: per-`(paymentIntentId, providerStatus)`. Backed by PRIMARY KEY on `payment_callback_dedup.dedup_key`.
- TTL: 30 days (matches cross-cutting.idempotency for payment.callback; >= provider redelivery horizon).
- Collision behavior:
  - SAME (intent, status) replay -> 200 with cached envelope (Step 4).
  - DIFFERENT status on a TERMINAL intent -> 409 `PAYMENT_INTENT_TERMINAL` with cached envelope (Step 5; cached on first conflict).
  - AMOUNT mismatch -> 409 `PAYMENT_AMOUNT_MISMATCH` with NO cache (Step 2; the next retry with a corrected amount can still apply cleanly).

## Performance

- p95 < 250ms (BA non_functional.latency for sync webhook).
- Single Postgres tx; one `SELECT FOR UPDATE` + one `UPDATE` + two `INSERT`s on the apply path.
- Provider redelivery storm (e.g., 10 redeliveries over 1 hour) is absorbed at Step 4 with zero contention beyond the per-intent row lock.

## Sequence diagram

```mermaid
sequenceDiagram
  autonumber
  participant PR as Mock Provider
  participant PM as Payment (handler)
  participant PC as process_callback
  participant DB as Postgres (payment)
  participant OBR as Outbox Relay (goroutine)
  participant KF as Kafka (ecom.payment.events)
  participant OD as Order consumer
  participant IV as Inventory consumer

  PR->>+PM: POST /api/v1/payment/intent/callback<br/>X-Mock-Provider-Signature, X-Mock-Provider-Timestamp<br/>{paymentIntentId, providerStatus, mockPaymentRef, amount, providerTimestamp}

  PM->>PM: validate Timestamp +/-5min skew
  PM->>PM: read rawBody (BEFORE ShouldBindJSON)
  PM->>PM: HMAC-SHA256(rawBody, MOCK_PROVIDER_SECRET)<br/>subtle.ConstantTimeCompare
  alt signature/timestamp invalid
    PM-->>PR: 401 AUTH_INVALID<br/>(NEVER log body/sig/secret)
  else valid signature
    PM->>+PC: process_callback(intent_id, provider_status, mock_provider_ref, amount, ts)

    Note over PC,DB: STEP 1 — LOCK
    PC->>+DB: BEGIN; SELECT * FROM payment_intents<br/>WHERE intent_id=$1 FOR UPDATE
    DB-->>-PC: intent row (or NULL)

    alt intent not found
      PC->>DB: ROLLBACK
      PC-->>PM: 404 PAYMENT_INTENT_NOT_FOUND
    else intent found
      Note over PC: STEP 2 — AMOUNT-CHECK (before TERMINAL!)
      alt amount_from_caller != intent.amount_minor
        PC->>DB: ROLLBACK
        PC-->>PM: 409 PAYMENT_AMOUNT_MISMATCH<br/>(no dedup row, no outbox)
      else amount matches
        Note over PC: STEP 3 — DEDUP-KEY DERIVE<br/>dedup_key = sha256(intent_id||"|"||provider_status)
        PC->>+DB: SELECT * FROM payment_callback_dedup<br/>WHERE dedup_key = $1
        DB-->>-PC: existing row (or NULL)

        alt STEP 4 — REPLAY HIT
          PC->>DB: COMMIT
          PC-->>PM: cached envelope<br/>(NEVER re-emits outbox)
        else dedup miss
          alt STEP 5 — TERMINAL GUARD<br/>(intent.status TERMINAL & status differs)
            PC->>DB: INSERT payment_callback_dedup<br/>(dedup_key, envelope=409 TERMINAL, ...)
            PC->>DB: COMMIT
            PC-->>PM: 409 PAYMENT_INTENT_TERMINAL
          else apply path (intent.status = REQUIRES_PAYMENT)
            Note over PC: STEP 6 — TRANSITION + OUTBOX (same tx)
            PC->>DB: UPDATE payment_intents SET status=new_status,<br/>mock_provider_ref, provider_status, paid_at, version+1
            alt new_status IN (SUCCEEDED, FAILED)
              PC->>DB: INSERT outbox_events<br/>(event_type, payload, partition_key=order_id, PENDING)
            else new_status = EXPIRED
              Note over PC: NO outbox in MVP<br/>(per ADR-009 + TL PR-001)
            end
            Note over PC: STEP 7 — DEDUP-SUCCESS ROW
            PC->>DB: INSERT payment_callback_dedup<br/>(dedup_key, envelope=200 SUCCESS, ON CONFLICT DO NOTHING)
            PC->>DB: COMMIT
            PC-->>-PM: 200 SUCCESS {applied: true, providerStatus}
          end
        end
      end
    end
    PM-->>-PR: response
  end

  Note over OBR,KF: Asynchronous (separate goroutine)
  OBR->>DB: SELECT outbox_events WHERE published_at IS NULL<br/>ORDER BY id LIMIT 100
  OBR->>+KF: producer.Send(topic=ecom.payment.events,<br/>key=order_id, value=payload)
  KF-->>-OBR: ACK
  OBR->>DB: UPDATE outbox_events SET published_at=now()
  KF-->>OD: events.payment.completed/failed (per-orderId order)
  KF-->>IV: events.payment.completed/failed (per-orderId order)
  Note over OD,IV: Consumers are STATE-DRIVEN per ADR-009<br/>SELECT FOR UPDATE on current state, ignore event.fromStatus
```

## Test cases

- `callback_with_valid_HMAC_and_matching_amount_succeeds`
- `callback_with_missing_signature_header_returns_401_AUTH_INVALID`
- `callback_with_bad_signature_returns_401_AUTH_INVALID_no_DB_writes`
- `callback_with_timestamp_skew_over_5min_returns_401_AUTH_INVALID`
- `callback_with_amount_mismatch_returns_409_PAYMENT_AMOUNT_MISMATCH_no_state_change`
- `callback_with_amount_mismatch_writes_NO_dedup_row`
- `callback_with_amount_mismatch_writes_NO_outbox_row`
- `callback_with_unknown_intent_returns_404_PAYMENT_INTENT_NOT_FOUND`
- `callback_replay_same_status_returns_cached_envelope_no_new_outbox_row` (PAY-005)
- `callback_replay_different_status_after_TERMINAL_returns_409_PAYMENT_INTENT_TERMINAL`
- `callback_amount_check_runs_BEFORE_terminal_state_guard_attacker_cannot_bypass_amount_validation_via_terminal_replay`
- `callback_concurrency_10_parallel_5_SUCCEEDED_5_FAILED_yields_exactly_one_terminal_state_one_outbox_row`
- `callback_dedup_key_uses_only_server_recorded_fields_changing_mockPaymentRef_or_amount_in_replay_does_NOT_force_dedup_miss`
- `callback_SUCCEEDED_emits_payment_completed_outbox_with_partition_key_orderId`
- `callback_FAILED_emits_payment_failed_outbox_with_partition_key_orderId`
- `callback_EXPIRED_emits_NO_outbox_event_in_MVP_feature_flag_PAYMENT_EMIT_EXPIRED_EVENT_default_false`
- `callback_logs_NEVER_contain_body_or_signature_or_secret`
- `callback_logs_DO_contain_dedup_key_intent_id_action`
- `callback_HMAC_uses_raw_request_bytes_not_re_marshaled_JSON`
- `callback_constant_time_compare_SHA256_pinned_no_caller_controlled_algo_field`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Opus 4.7 1M, dry-run #2) | Created. Canonical 7-step ordering pinned with amount-check at Step 2 (BEFORE terminal-state guard). HMAC-SHA256 + timestamp +/-5min replay protection. Dedup key composed of server-recorded fields ONLY. EXPIRED branch in Step 6 explicitly skipped per ADR-009 + TL PR-001. Mandatory mermaid sequence diagram covering signature -> 7-step -> outbox -> Kafka -> consumers. |
