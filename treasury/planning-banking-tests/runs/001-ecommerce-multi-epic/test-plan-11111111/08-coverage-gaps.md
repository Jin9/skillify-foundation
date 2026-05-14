---
artifact_type: qa-coverage-gaps
gap_count: 37
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Coverage Gaps

## [P1] missing-compliance
- Description: Five PDPA-TH-2019 compliance tests (COMP-PDPA-001..005) are scoped but compliance_scope_blocked=true because the BA brief has not produced a PII inventory, a retention schedule, a DSAR workflow, or a legal sign-off owner. Cannot execute these tests until governance gaps close.
- Required action: Resolve pii_inventory_missing, legal_absent_on_regulatory, retention_policy_unstated. Then unblock COMP-PDPA-001..005.
- Escalate to: Legal

## [P1] missing-compliance
- Story: EPIC-CART-CHECKOUT-3
- Description: Thai Revenue Code (TRC) and Consumer Protection Act (CPA) regulators have citation_status=pending in the BA brief. Compliance tests COMP-TRC-001 and COMP-CPA-001 are scope-blocked.
- Required action: Legal/Compliance to finalize regulator citations (TRC VAT receipt format, CPA price-transparency clause) before COMP tests can be designated mandatory.
- Escalate to: Legal

## [P1] missing-compliance
- Story: EPIC-ORDER-FULFILL-2
- Description: BA citation_status: pending for TRC (5-year accounting retention) and CPA (refund-timeline disclosure). Compliance tests are scope-blocked until citations land.
- Required action: BA / Compliance to confirm regulator citations (TRC section reference, CPA section reference) before COMP-TRC-001 and COMP-CPA-001 can be unblocked.
- Escalate to: Compliance

## [P1] missing-nfr-target
- Story: EPIC-CATALOG-1
- Description: Listing read:write ratio, expected scale, and pagination page-size are unspecified; capacity model and index strategy depend on these inputs.
- Required action: Resolve OQ-16 with PM and TL before TL design freezes the catalog query layer.
- Escalate to: BA

## [P1] missing-tl-design
- Story: EPIC-CART-CHECKOUT-3
- Description: TL has not specified idempotency_key_derivation_rule: which request fields (customer_id, cart_snapshot_hash, client-supplied nonce?) compose the key and how it is hashed/normalized. Without this, replay-safety tests TC-EPIC-CART-CHECKOUT-3-003 and webhook-idempotency TC-EPIC-CART-CHECKOUT-4-004 cannot deterministically assert key matching across replays.
- Required action: TL to publish idempotency-key derivation rule (input fields, normalization, hash algorithm, server-side storage TTL) before T-7 of execution start.
- Escalate to: TL

## [P1] missing-tl-design
- Story: EPIC-CART-CHECKOUT-3
- Description: TL has not specified concurrency_primitive for stock reservation: lock granularity (per-SKU row lock vs reservation-table advisory lock vs optimistic CAS), isolation level, and timeout. Race-resolution tests TC-EPIC-CART-CHECKOUT-3-002 and TC-EPIC-CART-CHECKOUT-3-008 need this to wire deterministic harness probes.
- Required action: TL to publish reservation-concurrency design (primitive, granularity, isolation, deadlock policy) before T-7.
- Escalate to: TL

## [P1] missing-tl-design
- Story: EPIC-ORDER-FULFILL-1
- Description: State-machine model not yet committed by TL: allowed transition matrix (awaiting_payment, paid, packing, shipped, delivered, cancelled), terminal states, and compensation-hook spec required to finalise invalid-transition test matrix.
- Required action: TL to publish state-machine spec (transition table + terminal-set + compensation hooks) and link to story EPIC-ORDER-FULFILL-1.
- Escalate to: TL

## [P1] missing-tl-design
- Story: EPIC-ORDER-FULFILL-2
- Description: Compensation-hook specification (stock-restore semantics, OOB notification contract, manual-refund operator endpoint shape) not yet committed by TL.
- Required action: TL to publish compensation-hook spec covering stock-restore atomicity and OOB-refund notification payload contract.
- Escalate to: TL

## [P1] qa_blocking_governance_gap
- Description: legal_absent_on_regulatory: Legal stakeholder absent on every epic while Thai-market regulators (PDPA, CCA, CPA, TRC) are implicit. Blocks compliance test scope for this epic (all 7 compliance_tests carry compliance_scope_blocked=true).
- Required action: Engage Legal counsel and DPO to confirm PDPA controller, lawful basis, retention schedule, DSAR workflow, breach-notification process, CCA log-retention duration, and CPA mandatory-disclosure scope before QA execution.
- Escalate to: Legal

## [P1] qa_blocking_governance_gap
- Description: pii_inventory_missing: pii_inventory exists at field level but retention/residency for several fields is 'TBD per PDPA review'. Blocks PII test-data design because masking-rule edges (e.g. partial-mask pattern length, residency boundary) are not finalized.
- Required action: Privacy-Lead + Legal to finalize per-field retention, residency, masking pattern, and access-audit rules. QA cannot ship redaction-correctness NFR target without finalized masking rules.
- Escalate to: Compliance

## [P1] qa_blocking_governance_gap
- Description: regulatory_citation_unresolved: All four regulators (PDPA-TH-2019, TRC-VAT-7PCT, CCA-TH-LOG-90D, CPA-TH-DISCLOSURE) carry citation_status='pending'. Blocks compliance test scope — exact citation required to bind test assertions to statute clauses.
- Required action: Legal + Compliance Officer to produce exact statute citations and confirm scope of application to ShopPilot MVP per regulator. Due 2026-05-26.
- Escalate to: Legal

## [P1] qa_blocking_governance_gap
- Description: retention_policy_unstated: Audit-log retention duration absent (5.9), customer-account retention not stated, order retention not stated. Blocks retention-test thresholds (cannot test 90-day enforcement without a stated threshold; AP-Q9 forbids inventing one).
- Required action: Compliance Officer + Legal to confirm retention schedule per data class and per audit-event class before QA execution of COMP-CCA-001 and COMP-CCA-004.
- Escalate to: Compliance

## [P1] qa_blocking_governance_gap
- Description: BA P1 governance gap: legal_absent_on_regulatory — blocks downstream QA execution
- Required action: Resolve BA P1 governance gap 'legal_absent_on_regulatory' in source brief; regenerate via ba-elicit-from-raw and re-run this skill.
- Escalate to: Legal

## [P1] qa_blocking_governance_gap
- Description: BA P1 governance gap: pii_inventory_missing — blocks downstream QA execution
- Required action: Resolve BA P1 governance gap 'pii_inventory_missing' in source brief; regenerate via ba-elicit-from-raw and re-run this skill.
- Escalate to: Compliance

## [P1] qa_blocking_governance_gap
- Description: BA P1 governance gap: regulatory_citation_unresolved — blocks downstream QA execution
- Required action: Resolve BA P1 governance gap 'regulatory_citation_unresolved' in source brief; regenerate via ba-elicit-from-raw and re-run this skill.
- Escalate to: Legal

## [P1] qa_blocking_governance_gap
- Description: BA P1 governance gap: retention_policy_unstated — blocks downstream QA execution
- Required action: Resolve BA P1 governance gap 'retention_policy_unstated' in source brief; regenerate via ba-elicit-from-raw and re-run this skill.
- Escalate to: Compliance

## [P1] qa_blocking_governance_gap
- Story: EPIC-IDENTITY-1
- Description: BA brief carries P1 governance gaps (legal_absent_on_regulatory, pii_inventory_missing, regulatory_citation_unresolved, retention_policy_unstated). IDENTITY is the primary PDPA surface — signup cannot ship to production until DPO + Legal confirm controller identity, lawful basis, retention schedule, and DSAR workflow.
- Required action: Engage Legal + DPO + Security Reviewer to close the four P1 governance gaps. Cannot execute COMP-PDPA-001..005 until closed.
- Escalate to: Legal

## [P1] qa_blocking_governance_gap
- Story: EPIC-IDENTITY-2
- Description: Login flow inherits the same four P1 governance gaps. Audit retention for login_success / login_failure / authz_denied events is undefined.
- Required action: Compliance Officer + Legal to confirm audit-event retention per class before sign-off.
- Escalate to: Compliance

## [P1] qa_blocking_governance_gap
- Story: EPIC-IDENTITY-3
- Description: Address book persists Thai postal PII (recipient, phone, full street) and feeds order snapshots. Residency rule + retention schedule for order-attached snapshots vs un-attached entries unstated.
- Required action: Legal + Compliance Officer to declare residency and per-row retention policy (live address vs order-snapshot copy).
- Escalate to: Legal

## [P2] ambiguous-requirement
- Story: EPIC-CATALOG-3
- Description: Behaviour for admin soft-deleting a product currently held in a customer's cart is unspecified on the admin side of the race.
- Required action: Resolve OQ-34 with PM and TL to define admin-side outcome.
- Escalate to: BA

## [P2] ambiguous-requirement
- Story: EPIC-CATALOG-3
- Description: Optimistic-locking policy for two admins concurrently editing the same product is unspecified; risk of silent lost update on price.
- Required action: Resolve OQ-35 with TL to define lock policy and conflict-response shape.
- Escalate to: TL

## [P2] ambiguous-requirement
- Story: EPIC-IDENTITY-1
- Description: Post-signup routing (auto-login vs route-to-login) unresolved by OQ-09. TC-EPIC-IDENTITY-1-001 last assertion is conditional until UX rule lands.
- Required action: PM + UX to publish the post-signup routing rule.
- Escalate to: BA

## [P2] ambiguous-requirement
- Story: EPIC-ORDER-FULFILL-3
- Description: BA notes customer self-reads 'typically logged but not as audit (decision noted as P2 OQ)'. Renderer treats as info-log not audit; requires BA confirmation.
- Required action: BA to confirm whether customer self-read on own order requires audit row.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-CART-CHECKOUT-3
- Description: Checkout-submit p95 latency target unresolved; AP-Q9 forbids inventing the number. NFR-EPIC-CART-CHECKOUT-perf-001 currently marked TBD.
- Required action: BA to resolve OQ-CART-checkout-p95-target with Payment-Lead and Ops.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-CART-CHECKOUT-3
- Description: Reservation TTL duration unresolved; FX-TIME-T0-T+10min assumes 10 min but actual TTL is TBD.
- Required action: BA to resolve OQ-CART-reservation-ttl with Payment-Lead.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-CART-CHECKOUT-3
- Description: Concurrent-submit load level N for race test (NFR-EPIC-CART-CHECKOUT-perf-002) unresolved.
- Required action: BA to resolve OQ-CART-concurrent-load-target.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-CATALOG-1
- Description: Listing-page first-paint p95 latency target unresolved; BA section 8.3 says 'within a fraction of a second' but gives no numeric SLO.
- Required action: Resolve OQ-12 to obtain numeric p50/p95 SLO before enabling perf gate.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-CATALOG-1
- Description: WCAG accessibility conformance level for customer-facing listing and PDP is unspecified; section 8.6 says only 'basic accessibility'.
- Required action: Resolve OQ-49 to confirm WCAG A vs AA vs AAA target.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-CATALOG-1
- Description: Bot / scraping / credential-stuffing defence threshold for the public catalog endpoints is unspecified beyond a generic rate-limit mention in section 8.1.
- Required action: Resolve OQ-37 with Security Reviewer to set concrete defence thresholds.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-IDENTITY-1
- Description: Password strength predicate is lexically vague ('reasonable strength') in BA section 4.2. Test TC-EPIC-IDENTITY-1-003 cannot bind a concrete assertion until OQ-03 resolves.
- Required action: Security Reviewer to publish the exact length / character-class / blacklist predicate.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-IDENTITY-1
- Description: Performance SLO for signup p95 absent (OQ-12 unresolved). NFR-EPIC-IDENTITY-perf-001 carries TBD target.
- Required action: PM + Performance Engineer to publish numeric p95 SLO for signup and login.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-IDENTITY-2
- Description: Session idle-timeout numeric undefined (OQ-02, OQ-11). Brute-force throttle numeric undefined (OQ-04).
- Required action: Security Reviewer + PM to publish numeric SLOs for session timeout and brute-force throttle.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-ORDER-FULFILL-1
- Description: State-machine reliability target (acceptable retry-storm error rate, skipped-state tolerance) not specified.
- Required action: BA / TL to specify retry-storm tolerance and zero-skipped-state acceptance criterion.
- Escalate to: BA

## [P2] missing-nfr-target
- Story: EPIC-ORDER-FULFILL-3
- Description: Order-tracking p95 latency target not specified (AP-Q9: no invented thresholds).
- Required action: BA to specify p95 latency target for order-detail GET; default placeholder 'TBD pending OQ-NFR-perf-fulfil'.
- Escalate to: BA

## [P2] missing-tl-design
- Story: EPIC-CART-CHECKOUT-3
- Description: TL has not specified error_envelope shape for payment/checkout failures: error code namespace, reason field, message-vs-code split, structured details, localization key. Error tests TC-EPIC-CART-CHECKOUT-3-002, TC-EPIC-CART-CHECKOUT-3-009, TC-EPIC-CART-CHECKOUT-4-005 currently assert on prose strings ('out_of_stock_at_confirm', 'coupon_expired') which are brittle.
- Required action: TL to publish standard error-envelope schema (JSON Schema or proto) and stable error-code registry before T-7.
- Escalate to: TL

## [P2] qa_blocking_governance_gap
- Story: EPIC-CART-CHECKOUT-3
- Description: Audit-event SLA (write latency, durability guarantee) unresolved (NFR-EPIC-CART-CHECKOUT-obs-001). TC-EPIC-CART-CHECKOUT-3-004 asserts write-before-respond but exact SLA is TBD.
- Required action: BA + TL to resolve OQ-CART-audit-sla.
- Escalate to: BA

## [P3] missing-nfr-target
- Story: EPIC-CART-CHECKOUT-4
- Description: Payment retry policy (max retries, backoff) for NFR-EPIC-CART-CHECKOUT-reliab-001 unresolved.
- Required action: BA to resolve OQ-CART-payment-retry-policy.
- Escalate to: BA

