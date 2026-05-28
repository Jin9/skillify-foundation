# NFR catalog — the "-ilities" and how to set concrete targets

A non-functional requirement is not done until it carries a **number, a unit, and a condition**. Vague adjectives
("fast", "scalable", "secure") are not requirements — they are wishes. For each NFR below, capture: the metric, the
target, the load/condition it holds under, and how it will be measured. If the realistic target depends on a technical
judgment the BA cannot make, record a draft and route it to the tech lead/architect.

## How to turn an adjective into a target
1. Name the **metric** (what you measure: latency, uptime, error rate…).
2. Set the **target** with a unit (300 ms, 99.9%, < 0.1%).
3. State the **condition** it must hold under (at 50 req/s, at peak, per calendar month).
4. Define **measurement** (where/how: p95 from the API gateway, monthly availability from the SLO dashboard).

Bad: "The API should be fast."
Good: "p95 latency < 300 ms and p99 < 800 ms for the search endpoint at 50 req/s sustained, measured at the gateway."

## The catalog

### Performance & responsiveness
- Latency: target percentiles, not averages — e.g. p95 < 300 ms, p99 < 800 ms.
- Throughput: requests/sec or transactions/sec at stated concurrency.
- Resource budget: CPU/memory ceiling per instance; max payload size.

### Availability & reliability
- Uptime SLO: e.g. 99.9% monthly (≈ 43 min downtime/month); pick a tier the business actually needs.
- Recovery: RTO (time to restore) and RPO (max acceptable data loss), e.g. RTO 15 min, RPO 5 min.
- Error budget / max error rate: e.g. < 0.1% 5xx over a rolling hour.
- Durability: e.g. no acknowledged write is lost; backups retained N days.

### Scalability & capacity
- Current vs projected load (12–24 months): users, peak concurrency, data volume.
- Data growth: e.g. +2 M rows/month; plan retention and archival.
- Scaling model: target headroom (e.g. handle 3× current peak without redesign).

### Security & compliance
- AuthN/AuthZ model; least-privilege; session/token lifetime.
- Data classification: what is PII/PCI/PHI; encryption at rest and in transit.
- Regulatory scope: GDPR, PCI-DSS, HIPAA, SOC 2 — name the specific controls in scope.
- Audit: what events must be logged, and retention period.

### Usability & accessibility
- Task success target / max steps for a key task; error-recovery expectations.
- Accessibility standard (e.g. WCAG 2.1 AA) as a pass/fail target.
- Localization: supported locales, formats, right-to-left needs.

### Maintainability & operability
- Observability: required metrics, logs, traces; alert thresholds tied to the SLOs above.
- Deployability: target release frequency; rollback time.
- Test coverage / quality gate the change must meet.

### Portability & interoperability
- Supported platforms, browsers, devices, API versions.
- Required standards/formats for data exchange.

## Anti-patterns
- Targets with no unit or condition ("handle lots of users").
- Averages where tails matter (use p95/p99 for latency).
- Copy-pasting "five nines" availability when the business does not need or fund it.
- Inventing NFR numbers the BA cannot justify — draft them, then have the technical owner confirm.
