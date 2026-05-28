---
template_version: 0.1.0
epic_id: EPIC_CATALOG
title: Product browse, search, filter, and detail
summary: |
  Guests and customers must browse the product catalog, filter by category and price,
  check stock availability, and open a product-detail surface that the cart and checkout
  routes consume. Admin-side product create / edit / soft-delete and category management
  are exposed only as backend endpoints in MVP (no admin UI).
business_value: The catalog is the customer's entry point. Without it the home and search tabs of the customer mobile webview cannot render and the funnel into cart and checkout never starts.
in_scope_stories:
  - STORY_CATALOG_LIST
  - STORY_CATALOG_DETAIL
  - STORY_CATALOG_ADMIN_PRODUCT_CRUD
  - STORY_CATALOG_ADMIN_CATEGORY
out_of_scope:
  - Product Review module (REV-001..005) deferred per §4.2; review summary on detail returns placeholder data
  - AI/personalised recommendations (§4.2 #8 and #13) deferred
  - Admin product / category management UI deferred — backend endpoints only this run
  - Multi-currency listing and full multi-language UI deferred per §4.2 #11 and #12
  - Multi-vendor catalog (§4.2 #3) deferred
success_metrics:
  - Listing endpoint p95 < 300 ms on a 1,000-SKU dataset (PERF-001)
  - Listing returns ONLY products with status=ACTIVE and visible=true (CAT-001)
  - Soft-deleted products remain resolvable from historical orders (CAT-009)
  - Duplicate-SKU on admin create returns CONFLICT (CAT-007, CAT-010)
dependencies:
  - EPIC_AUTH
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Refactored from runs/01; explicit out-of-scope on review summary and admin UI surfaces (per Frontend Spec scope rule).
---

# EPIC_CATALOG — Product browse, search, filter, and detail

## Why

The mobile webview's home and search tabs (Frontend Spec §3) are the funnel mouth for the entire customer journey. Without listing, filter, and product detail working, no item ever reaches a cart and no order is ever created.

The original requirement (§8.3, §8.4) defines the full catalog surface: listing with pagination/filter/sort (CAT-001..005), product detail (CAT-006), admin product create / edit / soft-delete (CAT-007..009), the four-status enum (CAT-010), and category create / edit / inactive (CATE-001..003).

In MVP we ship the **customer-visible** half of this on the frontend (listing, filter, detail) and keep the admin operations as **backend-only** endpoints exercised via Postman / curl, not through a built admin UI. This is the deliberate scope cut from the Frontend Spec ("customer-facing only — admin is out of scope").

The product status enum is exactly four values (DRAFT, ACTIVE, INACTIVE, DELETED) and only ACTIVE+visible products are exposed in listing — every other status is admin-internal.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_CATALOG_LIST` | Guest/Customer lists products with pagination, filter, sort | Must |
| `STORY_CATALOG_DETAIL` | Guest/Customer reads product detail | Must |
| `STORY_CATALOG_ADMIN_PRODUCT_CRUD` | Admin creates / edits / soft-deletes a product | Must |
| `STORY_CATALOG_ADMIN_CATEGORY` | Admin creates / edits / deactivates a category | Must |

## Out-of-scope rationale

- **Review summary on detail** — REV-001..005 deferred (no review module in MVP); the detail endpoint returns an empty/placeholder review block that the frontend can render as a "no reviews yet" state.
- **AI/personalised recommendations** — explicitly deferred per §4.2; we do not invent a "you may also like" surface that the user did not ask for.
- **Admin UI** — Frontend Spec is customer-facing; admin product / category endpoints exist but no admin web UI is built this run.

## Risks

- **Soft-delete vs historical order resolution** — CAT-009 requires that a soft-deleted product still resolves from a prior order's snapshot. Implementations that hard-cascade product deletes will break order detail; the order-item snapshot must own enough of the product fields.
- **Listing performance ceiling** — PERF-001 caps listing at p95 < 300 ms on 1,000 SKUs; combined filter + sort queries must be index-backed.
- **Category-name uniqueness** — CATE-001 requires uniqueness *within a level*, not globally; a parent's "Headphones" must not collide with another parent's "Headphones".

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Refactored from runs/01; explicit out-of-scope on review summary and admin UI surfaces. |
