# Route spec — `/orders` (orders tab)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_ORDERS`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_ORDERS.md)
- **Auth:** required
- **Type:** Server Component + Client filter chips
- **Tab route:** yes

## Purpose

Customer's order history hub. Lists orders sorted descending by `createdAt` with Thai status badges. Filter chips (`ทั้งหมด` / `กำลังดำเนินการ` / `เสร็จสิ้น`) bucket orders by lifecycle state. Tapping a row navigates to `/orders/[orderId]`.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `order.list-mine` (page, limit, status?) | Order list paginated |

Routed through `withAuth`. The layout is wrapped in `ProtectedShell` — unauthenticated visit 307-redirects to `/login?next=/orders`.

## Key components

- `ShopHeader`
- `Chip` group (filter all/active/done — single-select; updates URL `?status=...`)
- `OrderSummaryCard` x N — `orderNumber`, `StatusBadge`, total, `createdAt`
- `StatusBadge` — Thai-labeled colored pill
- `EmptyOrders` — Thai placeholder with CTA → `/`
- `Pagination`
- `TabBar`

## Filter buckets

| Chip | Status filter | Statuses included |
|---|---|---|
| ทั้งหมด | none | all |
| กำลังดำเนินการ | `active` | PENDING_PAYMENT, PAID, PACKING, SHIPPED |
| เสร็จสิ้น | `done` | DELIVERED, CANCELLED, PAYMENT_FAILED, PAYMENT_EXPIRED |

The `active`/`done` aliases are translated to a `status[]` filter in the request body.

## Mobile layout

- One card per order with status pill, total in ฿ prefix, and `createdAt` formatted via `Intl.DateTimeFormat('th-TH-u-ca-buddhist')`
- Filter chips sticky at top under `ShopHeader`

## Acceptance criteria covered

- **AC1 (filter + Thai badges + tap → detail):** Three chips toggle filter; each row shows a Thai status badge; tap row → `/orders/[orderId]`.
- **AC2 (empty state):** No orders → `EmptyOrders` with CTA back to `/`.

## Edge cases

| Case | Expected |
|---|---|
| `PAYMENT_EXPIRED` in list | Distinct rose-colored badge + Thai label "หมดเวลาชำระเงิน" |
| Pagination | Next page replaces visible list (simpler for MVP); existing rows do not re-render — page param drives SC re-fetch |
| Auth expired during navigation | Layout's `withAuth` refreshes cookies transparently OR redirects to login |

## Notes

- The `/orders` Server Component reads `searchParams.status` and `searchParams.page` directly. Client `Chip` group uses `router.replace(...)` on tap.
