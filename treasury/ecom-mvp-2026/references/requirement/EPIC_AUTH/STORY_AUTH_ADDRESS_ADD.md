---
template_version: 0.1.0
story_id: STORY_AUTH_ADDRESS_ADD
epic_id: EPIC_AUTH
title: Customer adds a shipping address and selects default
as_a: Customer
i_want: to add a shipping address and choose it as my default
so_that: checkout can ship orders to me and pre-select my default address
acceptance_criteria:
  - given: An authenticated customer with no addresses on file
    when: The client POSTs identity address-create with receiverName, phone, addressLine, province, district, postalCode, isDefault=true
    then: |
      Response is 200/201 with code=SUCCESS, the address is persisted bound to the
      requesting userId, and the customer's defaultAddressId points to the new
      address (CUST-003, CUST-004).
  - given: An authenticated customer with one default address A and a fresh request adding address B with isDefault=true
    when: The client POSTs identity address-create with isDefault=true for address B
    then: |
      Address B becomes the customer's default; address A persists but is no longer
      default. The customer never has two default addresses simultaneously
      (CUST-004).
priority: Must
size: M
edge_cases:
  - case: POST address-create with a missing required field (e.g. no postal code)
    expected: 400 VALIDATION_FAILED listing the missing field; no address row created
  - case: POST address-create as customer A but include a userId field for customer B in the body
    expected: The body's userId is ignored; the address is bound to the requesting customer's userId from the access token
  - case: GET addresses for a different customer using customer A's token
    expected: Listing returns only customer A's addresses; never customer B's
requirement_refs: [CUST-003, CUST-004]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.2 CUST-003 + CUST-004; userId-from-token enforcement added since checkout (CHK-002) depends on ownership.
---

# STORY_AUTH_ADDRESS_ADD — Customer adds a shipping address and selects default

## User narrative

As a **Customer**, I want **to add a shipping address and choose it as my default** so that **checkout can ship orders to me and pre-select my default address**.

## Why this story exists in EPIC_AUTH

The mobile webview's checkout route depends on the customer having at least one shipping address; the addresses route lists / adds them and exposes a picker mode for checkout (Frontend Spec §3). CHK-002 also requires checkout to verify the chosen addressId belongs to the requesting userId — that ownership check is enforced here at create-time.

## Edge cases

| Case | Expected |
|---|---|
| Missing required field | 400 VALIDATION_FAILED listing the missing field |
| Body's userId differs from token | Body is ignored; token's userId is authoritative |
| GET addresses across customers | Listing scoped to requesting customer's userId only |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.2 CUST-003 + CUST-004. |
