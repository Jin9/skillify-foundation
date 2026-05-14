---
artifact_type: qa-nfr-tests
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# NFR Tests

## accessibility

### NFR-EPIC-CATALOG-a11y-001 — WCAG conformance level on listing and PDP customer-facing surfaces
- Epic: EPIC-CATALOG
- Target: TBD pending OQ-49
- Blocking OQs: OQ-49

### NFR-EPIC-IDENTITY-a11y-001 — Signup form WCAG 2.1 AA conformance (labels, focus order, color contrast, keyboard nav, error announcement)
- Epic: EPIC-IDENTITY
- Target: 0 critical/serious axe-core violations; SR-narrated error messages for email/password/phone fields

### NFR-EPIC-IDENTITY-a11y-002 — Login form WCAG 2.1 AA conformance
- Epic: EPIC-IDENTITY
- Target: 0 critical/serious axe-core violations; generic auth error is announced via aria-live

### NFR-EPIC-IDENTITY-a11y-003 — Address book form WCAG 2.1 AA conformance (Thai postal fields)
- Epic: EPIC-IDENTITY
- Target: 0 critical/serious axe-core violations; subdistrict/district/province/postcode each have programmatic labels

### NFR-EPIC-ORDER-FULFILL-a11y-001 — order tracking page WCAG 2.1 AA conformance (customer-facing)
- Epic: EPIC-ORDER-FULFILL
- Target: WCAG 2.1 AA - zero Critical, zero Serious axe-core violations; status timeline screen-reader announces transitions in order

## data-integrity

### NFR-EPIC-CATALOG-data-001 — customer complaints citing 'product visible but cannot be bought' per week
- Epic: EPIC-CATALOG
- Target: 0 in the first 30 days post-launch

### NFR-EPIC-GOVERNANCE-OBS-data-001 — pii_redaction_correctness_pct_on_sample_log_corpus
- Epic: EPIC-GOVERNANCE-OBS
- Target: TBD pending pii_inventory_missing gap resolution (masking_rule per field finalized) and OQ-resolved sample size; AP-Q9 forbids invented thresholds

## latency

### NFR-EPIC-CATALOG-perf-001 — p95 listing-page first-paint latency under normal load
- Epic: EPIC-CATALOG
- Target: TBD pending OQ-12
- Blocking OQs: OQ-12

## observability

### NFR-EPIC-CART-CHECKOUT-obs-001 — audit_event_emission_rate_per_payment_outcome
- Epic: EPIC-CART-CHECKOUT
- Target: TBD pending OQ-CART-audit-sla
- Blocking OQs: OQ-CART-audit-sla

## performance

### NFR-EPIC-CART-CHECKOUT-perf-001 — checkout_submit_p95_latency_ms
- Epic: EPIC-CART-CHECKOUT
- Target: TBD pending OQ-CART-checkout-p95-target
- Blocking OQs: OQ-CART-checkout-p95-target

### NFR-EPIC-CART-CHECKOUT-perf-002 — concurrent_checkout_race_resolution_no_oversell_under_N_concurrent
- Epic: EPIC-CART-CHECKOUT
- Target: TBD pending OQ-CART-concurrent-load-target
- Blocking OQs: OQ-CART-concurrent-load-target

### NFR-EPIC-IDENTITY-perf-001 — Signup p95 response time under normal load
- Epic: EPIC-IDENTITY
- Target: TBD pending OQ-12 (no numeric SLO declared in BA brief 8.3/3.3)
- Blocking OQs: OQ-12

### NFR-EPIC-IDENTITY-perf-002 — Login p95 response time under normal load
- Epic: EPIC-IDENTITY
- Target: TBD pending OQ-12
- Blocking OQs: OQ-12

### NFR-EPIC-ORDER-FULFILL-perf-001 — order tracking / order-detail page p95 latency (customer-facing GET)
- Epic: EPIC-ORDER-FULFILL
- Target: TBD pending OQ-NFR-perf-fulfil
- Blocking OQs: OQ-NFR-perf-fulfil

## resilience

### NFR-EPIC-CART-CHECKOUT-reliab-001 — payment_retry_under_timeout_idempotency_preservation
- Epic: EPIC-CART-CHECKOUT
- Target: TBD pending OQ-CART-payment-retry-policy
- Blocking OQs: OQ-CART-payment-retry-policy

### NFR-EPIC-GOVERNANCE-OBS-reliab-001 — audit_event_durability_loss_count_under_induced_crash
- Epic: EPIC-GOVERNANCE-OBS
- Target: TBD pending OQ-20 (durability target not stated; AP-Q9)
- Blocking OQs: OQ-20

### NFR-EPIC-ORDER-FULFILL-reliab-001 — state-machine consistency under client retry / network partition; no skipped intermediate states observed across 1000 simulated retries
- Epic: EPIC-ORDER-FULFILL
- Target: TBD pending OQ-NFR-reliab-fulfil (acceptance: zero skipped-state events; zero out-of-order audit rows)
- Blocking OQs: OQ-NFR-reliab-fulfil

## scalability

### NFR-EPIC-CATALOG-perf-002 — listing pagination page-size and read:write ratio capacity
- Epic: EPIC-CATALOG
- Target: TBD pending OQ-16
- Blocking OQs: OQ-16

## security

### NFR-EPIC-CART-CHECKOUT-sec-001 — pan_scope_isolation_no_card_data_in_app_logs_or_audit
- Epic: EPIC-CART-CHECKOUT
- Target: 0 occurrences of PAN/CVV across audit-sink and app-log scans (PCI out-of-scope invariant)

### NFR-EPIC-CATALOG-sec-001 — bot / scraping / credential-stuffing defence threshold on public catalog endpoints
- Epic: EPIC-CATALOG
- Target: TBD pending OQ-37
- Blocking OQs: OQ-37

### NFR-EPIC-IDENTITY-sec-001 — Cross-customer authz denial rate on address edit
- Epic: EPIC-IDENTITY
- Target: 100% denial of cross-owner address edits (DoD §5: 0 leakage incidents)

### NFR-EPIC-IDENTITY-sec-002 — Login anti-enumeration parity (response body + timing)
- Epic: EPIC-IDENTITY
- Target: Identical message text AND p95 timing delta < 25ms between unknown-email and wrong-password paths

### NFR-EPIC-IDENTITY-sec-003 — Session token entropy
- Epic: EPIC-IDENTITY
- Target: Token derived from CSPRNG with >= 128 bits of entropy; base64-opaque; never echoed in logs or audit payload

### NFR-EPIC-IDENTITY-sec-004 — Brute-force throttle on /login
- Epic: EPIC-IDENTITY
- Target: TBD pending OQ-04
- Blocking OQs: OQ-04

### NFR-EPIC-IDENTITY-sec-005 — Password strength predicate (length + character classes + blacklist)
- Epic: EPIC-IDENTITY
- Target: TBD pending OQ-03
- Blocking OQs: OQ-03

### NFR-EPIC-IDENTITY-sec-006 — Session idle timeout
- Epic: EPIC-IDENTITY
- Target: TBD pending OQ-02 / OQ-11
- Blocking OQs: OQ-02, OQ-11

## throughput

### NFR-EPIC-GOVERNANCE-OBS-obs-001 — audit_log_write_throughput_events_per_second_sustained_p95_latency
- Epic: EPIC-GOVERNANCE-OBS
- Target: TBD pending OQ-20 (audit-log throughput / retention SLO not stated; AP-Q9 forbids invented thresholds)
- Blocking OQs: OQ-20

