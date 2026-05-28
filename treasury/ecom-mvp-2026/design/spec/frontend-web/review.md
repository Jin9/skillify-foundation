# Route spec — `/orders/[orderId]/review/[itemId]` (review writer — DEFERRED)

- **Component:** `frontend-web`
- **Story:** N/A in this run (review module is deferred per [EPIC_FRONTEND.out_of_scope[1]](../../../requirement/EPIC_FRONTEND/EPIC_FRONTEND.md))
- **Auth:** required (would be, when wired)
- **Type:** placeholder Server Component
- **Tab route:** no

## Purpose

Frontend Spec §3 lists a `review` route for writing a star rating + comment for an order item. EPIC_FRONTEND explicitly defers the review module (`out_of_scope[1]`), so this route ships as a placeholder that returns 200 with a Thai "เร็ว ๆ นี้" message — preserving graceful behavior on deep links from external sources.

## Backend dependencies

None (no review contract exists in `contracts.json` for this MVP).

## Key components

- `ShopHeader` (back)
- `EmptyState` — Thai message: "เร็ว ๆ นี้ — ฟีเจอร์การรีวิวจะเปิดให้ใช้งานในเวอร์ชันถัดไป"
- `Btn` (light) "กลับ" → `/orders/[orderId]`

## Acceptance criteria covered

None binding in MVP. The route exists so:

- Deep links from external sources don't 404 (graceful 200 with placeholder copy).
- The route surface is documented for the post-MVP review module work, which will lift this from placeholder to a star-rating + textarea form bound to a future `review.create` contract.

## Notes

- This is the 13th route in the Frontend Spec §3 inventory minus the `orderResult` row. The TD's pages array enumerates 13 distinct routes — `review` is included only as a placeholder file (no AC backing); not counted as a "live" route in the test plan or `ba_acceptance_mapping`.
- If Tech-Lead prefers we omit the placeholder file entirely (return 404 instead), it's a one-line removal.
