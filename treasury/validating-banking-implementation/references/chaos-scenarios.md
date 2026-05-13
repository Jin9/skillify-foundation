# Chaos test scenarios

Draft (do not execute). Choose scenarios that match the flow being
QA'd. Every scenario specifies: trigger, expected behavior, rollback
verification, and the metric / alert that proves the system handled it.

## Universal expectations

For every chaos scenario:

- The system either degrades gracefully (read-only, queued retry,
  documented user-visible error) or explicitly fails closed (rejects
  the operation).
- No partial monetary state is left behind.
- Audit log entries reflect the failure, not silence.
- The user-visible error message does not leak internal detail.

## Scenario library

### 1. Partner bank API timeout (disbursement)

**Trigger:** the partner bank's `/disburse` endpoint times out at 30s.
**Expected:** disbursement marked `pending`, idempotency key persisted,
saga's compensating step armed but not fired. Customer-facing status
shows "processing".
**Rollback:** when partner returns `success` later (or the saga retries
and gets `success`), no double-disbursement.
**Metric:** `disbursement.partner_timeout_total` increments;
`disbursement.duplicate_attempt_total` stays at 0.

### 2. Partner bank returns ambiguous error

**Trigger:** partner returns 5xx with a body like "Transaction may have
succeeded".
**Expected:** the system pauses and asks for human reconciliation
rather than auto-retrying. Disbursement saga state = `awaiting-recon`.
**Rollback:** human operator confirms or denies via the recon UI; saga
proceeds.
**Metric:** alert fires when `disbursement.awaiting_recon` > 0.

### 3. DB primary unavailable for 60s (planned failover)

**Trigger:** DB primary terminates; replica promoted.
**Expected:** writes return retriable-error during the gap; reads can
serve from the replica with a stale-read warning (where the flow
permits).
**Rollback:** in-flight transactions either commit or abort cleanly
post-failover; no half-committed states.
**Metric:** `db.write_unavailable_seconds` measured; `< 90s` SLO.

### 4. Kafka broker partial outage

**Trigger:** one of three brokers becomes unreachable.
**Expected:** producers continue with reduced ack guarantee; consumers
re-balance; no duplicate processing of in-flight messages.
**Rollback:** broker recovers; consumer-lag drains within SLO.
**Metric:** `kafka.consumer_lag_seconds` returns under SLO within 5
min.

### 5. Sudden 10× traffic spike on KYC verification

**Trigger:** synthetic load on `/kyc/verify`.
**Expected:** rate limiter triggers; queue backpressure ripples to the
callers; no service crash.
**Rollback:** spike subsides; queue drains.
**Metric:** `kyc.rate_limit_triggered_total` increments; no
`http.5xx{route="/kyc/verify"}` increase.

### 6. Cache eviction storm (Redis)

**Trigger:** Redis primary OOM-kills; failover to replica with empty
cache.
**Expected:** services degrade to DB lookups with bounded concurrency;
no thundering herd kills the DB.
**Rollback:** cache warms via background prefetch; latency returns to
SLO within 10 min.
**Metric:** `db.qps` returns to baseline; `cache.miss_ratio` returns
under SLO.

### 7. Clock skew between services

**Trigger:** one node's NTP drifts 5 minutes ahead.
**Expected:** JWT validation tolerance window holds; audit-log
timestamps are flagged but not silently corrupted.
**Rollback:** NTP corrects; subsequent timestamps are normal.
**Metric:** alert on `clock.skew_seconds` > 60.

### 8. Underwriter override during disbursement

**Trigger:** underwriter manually changes credit decision while
disbursement saga is running.
**Expected:** saga detects the override and either pauses for re-
approval or completes only if the override is dual-controlled.
**Rollback:** override flow is auditable, single-actor overrides are
rejected.
**Metric:** `override.single_actor_rejected_total` increments;
`override.applied_total` only on dual-control flows.

### 9. Sandbox / proxy egress failure (sandboxed implement)

**Trigger:** the implement runner attempts to fetch a non-allowlisted
host.
**Expected:** the runner aborts the flow rather than ship code that
ran outside the sandbox.
**Rollback:** workflow halted; no code mutation persisted.
**Metric:** `sandbox.egress_denied_total` increments.

### 10. Replay attack on idempotency key

**Trigger:** legitimate request retried with same idempotency key 24h
later (post-dedup-window).
**Expected:** request rejected as expired-key, not silently re-executed.
**Rollback:** caller refreshes key and retries.
**Metric:** `idempotency.expired_key_total` increments.

## How to use this library

1. Identify the flow under QA.
2. Select 3–5 scenarios from this library that match the flow's failure
   modes.
3. For each, fill the four fields: trigger, expected, rollback, metric.
4. Add to the chaos-test plan section of the QA report.
5. Mark the plan **draft** — execution is not part of this skill.
