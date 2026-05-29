<!-- Exemplar referenced by SKILL.md: a full `user-story-template` instance at the
     richness a story reaches WHEN the Scope Sheet supplies conditional logic —
     numbered Business-Logic rules with a./b. sub-items, a `### <name>` decision
     table, and a `> ⚠️ Note`. Synthetic values only; never real PII. When the scope
     gives no conditional logic, omit the table — do not invent one to look rich. -->

<!-- ba-pipeline: story-set REQ-jaipa-tier-rate · ST-04 · state: drafted
     scope_ref: ../01-scope-sheet.md · open_questions: [OQ-2] -->

`< Back`  ·  `<none>`

# [Feature] Offer tier-rate interest on the Jai Pa personal loan

## Description

*As a* lending business unit

*I want* to offer tiered interest rates on the Jai Pa personal-loan product to customers whose verified monthly income clears the risk-engine threshold

*So that* the offered rate tracks the customer's salary base and stays competitive without manual underwriting.

## Product

Jai Pa — personal loan (new rate offering)

## Background

1. Today the product offers a single flat rate regardless of income band.
2. The risk engine already returns an income-tier decision; the loan-origination flow does not yet consume it.
3. Marketing wants a promotional first-period rate for higher-income, payroll-linked customers.

## Business Logic

1. Tier eligibility is driven by the **risk engine**, not computed in the loan flow:
   a. inputs are income source (`SA` salaried / `SE` self-employed), payroll-with-bank flag (`Yes`/`No`), and verified monthly income.
   b. the loan flow consumes the returned tier verbatim; it never re-derives or overrides a tier.
2. The system must store and apply a **multi-tier** rate schedule: the requirement covers 2 tiers; the platform must support at least 5.
3. Below the income threshold the customer is offered a **single-tier** (flat) rate; at or above it, the multi-tier schedule in the table below applies.
4. The promotional first-period rate applies only to `SA` + payroll `Yes`; every other eligible profile takes the standard multi-tier rate.

### Income → interest-rate tiers

| Income source | Payroll with bank | Verified monthly income | Offered rate |
|---|---|---|---|
| SA | Yes | &lt; 30,000 THB | single-tier 18.00% p.a. |
| SA | Yes | ≥ 30,000 THB (multi-tier) | months 1–3: 7.99% p.a.; month 4+: 20.00% p.a. |
| SA | No  | ≥ 30,000 THB (multi-tier) | months 1–3: 12.99% p.a.; month 4+: 22.00% p.a. |
| SE | —   | ≥ 30,000 THB (multi-tier) | months 1–6: 15.99% p.a.; month 7+: 24.00% p.a. |

> ⚠️ **Note:** rates follow whatever the risk engine returns for the customer's profile; the four driving inputs are income, income source, payroll flag, and the offer's effective date. A schedule already disbursed is never re-tiered retroactively.

## Acceptance Criteria

1. Given an `SA` customer with payroll `Yes` and verified income 45,000 THB, when the offer is generated, then a multi-tier schedule is stored with months 1–3 at 7.99% p.a. and month 4+ at 20.00% p.a.
2. Given an `SA` customer with payroll `Yes` and verified income 24,000 THB, when the offer is generated, then a single-tier flat rate of 18.00% p.a. is stored (no promotional tier).
3. Given an `SE` customer with verified income 50,000 THB, when the offer is generated, then a multi-tier schedule is stored with months 1–6 at 15.99% p.a. and month 7+ at 24.00% p.a.
4. Given any eligible customer, when more than 5 rate tiers are requested, then the system rejects the schedule with "max 5 tiers supported" and stores no offer.

## Out of Scope

1. Renaming the "interest rate" label on the customer-facing offer screen (owned by the Channels team).
2. Retroactive re-tiering of loans already disbursed (Phase 2 — carried as `OQ-2`).
