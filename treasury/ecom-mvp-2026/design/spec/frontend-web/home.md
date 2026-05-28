# Route spec — `/` (home tab)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_HOME`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_HOME.md)
- **Auth:** public (no login required)
- **Type:** Server Component
- **Tab route:** yes (5-tab bottom bar visible)

## Purpose

Landing screen for the customer mobile webview. Establishes the visual contract for every other route (Thai language, 390-wide single column, IBM Plex Sans Thai, emerald `#1a8060` brand, ฿ Thai-Baht prefix, 5-tab bottom navigation). Surfaces three sections (featured, categories, new arrivals) so the visitor can tap into product detail or browse by category.

## Backend dependencies

| Endpoint | Why we call it |
|---|---|
| `catalog.product.list` (sort=newest, page=1, limit=12) | "New arrivals" + "Featured" sections |
| `catalog.category.list` (activeOnly=true) | "Categories" section |

Both are public — no auth header forwarded; the generic `/api/proxy/[...path]/route.ts` is used for both calls.

## Key components

- `ShopHeader` (top bar, no back button)
- `Section` x3 (Featured / Categories / New arrivals)
- `ProductCard` (image placeholder + name + ฿ price)
- `Chip` (category pill in horizontal scroll)
- `TabBar` (bottom — `home` highlighted)

## Mobile layout

- Single 390-wide column, page padding 16px
- Card gap 12px, section spacing 8px
- Horizontal scroll for category chips and featured carousel — vertical scroll for new arrivals grid (2-col)
- Reflows down to 360-wide without horizontal scroll

## Acceptance criteria covered

- **AC1 (visual contract):** lang='th'; 390 column with no horizontal scroll; bottom tab bar visible with 5 icons; primary CTA computed-style emerald `#1a8060`; ฿ prefix on every price.
- **AC2 (data sections render):** Featured + Categories + New-arrivals each render at least one item from `catalog.product.list` and `catalog.category.list`. Tapping a product card navigates to `/detail/[productId]`.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Catalog returns empty list | Empty-state placeholder in Thai per section | Server Component renders `<EmptyState message='ยังไม่มีสินค้า' />` when result array is empty |
| Catalog request fails (envelope code != SUCCESS) | Error state with traceId in mono footer | `ErrorBoundary` (segment `error.tsx`) renders fallback |
| Visitor offline | Browser-native error → segment `error.tsx` shows Thai message; recovers on reconnect | No special handling beyond standard ErrorBoundary |

## Notes

- No interactivity on this route → no `'use client'` boundary needed. `TabBar` is a Client wrapper for `usePathname()` highlighting only.
- This route is the visual oracle for every subsequent page. Any token regression (color, typography, spacing) will reproduce here.
