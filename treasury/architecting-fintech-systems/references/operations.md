# Operations & Observability

## Contents
- Structured logging
- Metrics (RED / USE)
- SLOs & error budgets
- Graduation criteria
- Retry & incident handling
- Alerting

## Structured logging

JSON by default. Every log entry MUST include: `status`, `path`, `business_id` (e.g. `loan_app_id`), `latency_ms`, `error_class`, `tenant_id`.

Prefer log-based metrics — cost-efficient relative to custom metrics. Extract only the necessary fields so queries do not parse raw text.

## Metrics (RED / USE)

- **RED** (request-scoped): Rate, Error rate, Duration — for every service endpoint.
- **USE** (resource-scoped): Utilization, Saturation, Errors — for infrastructure (CPU, DB pool, queue depth).
- Instrument both. RED drives SLO tracking; USE drives capacity planning.

## SLOs & error budgets

- Define SLOs for every user-facing path (e.g. p99 < 200ms, availability > 99.9%).
- Error budget = 1 − SLO target, measured over a rolling 30-day window.
- Burn alerts: fast burn (5% in 1h) and slow burn (10% in 6h).
- When the budget is exhausted, freeze features and prioritize reliability.

## Graduation criteria

- <1M events/day: log-based metrics + basic dashboards are sufficient.
- 1–10M events/day: add pre-aggregated counters; consider a time-series DB.
- >10M events/day: MUST move to pre-aggregated metrics (StatsD/Prometheus histograms). Raw log queries become cost-prohibitive.

## Retry & incident handling

- Prefer DB-driven retry workers over blind Kafka retries.
- Support delayed retry, scheduled retry, and manual recovery.
- Every retry MUST track: attempt count, last error, next retry time, and state.
- Dead-letter with alerting for exhausted retries.

## Alerting

- Alert on symptoms (error rate, latency), not causes.
- Business-level alerts alongside infrastructure alerts (e.g. "0 loan applications processed in 30min").
- Route alerts by severity: page (P1), ticket (P2), dashboard (P3+).
