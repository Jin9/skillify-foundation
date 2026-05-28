# Route spec — `/search` (search tab)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_SEARCH`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_SEARCH.md)
- **Auth:** public
- **Type:** Server Component shell + Client filter pane
- **Tab route:** yes

## Purpose

Discovery surface beyond home — combines a query input, category filter chips, optional price range, in-stock toggle, and a paginated 2-column product grid. URL search params drive the SC fetch so links are shareable.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `catalog.product.list` (page, limit, sort, categoryId, minPrice, maxPrice, inStock, search) | Main grid |
| `catalog.category.list` (activeOnly=true) | Filter chips |

## Key components

- `ShopHeader`
- `SearchBar` (`'use client'`, debounced 300ms input → `router.replace('?...search=...')`)
- `Chip` x N (category pills, single-select)
- `PriceRangeSheet` (`'use client'`, bottom Sheet with min/max sliders → URL params)
- `Toggle` (`inStock` checkbox)
- `ProductGrid` + `ProductCard`
- `Pagination` (page links computed from total/limit)
- `TabBar`

## Mobile layout

- Search input fixed at top below `ShopHeader`
- Filter chips horizontal scroll under search
- 2-col product grid; 12px gap
- Pagination footer above bottom `TabBar`

## Acceptance criteria covered

- **AC1 (filters drive results):** typing query + selecting category + price range + inStock re-renders the grid against `catalog.product.list` with matching pagination + filter contract; each card has name + ฿ price + category; tap → `/detail/[productId]`.
- **AC2 (empty state):** zero results → Thai empty-state placeholder; clearing filters restores full result set.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Pagination beyond first page | Page=2 fetch grows the visible list (or replaces, per UX choice — MVP replaces, simpler) | URL `?page=2` triggers SC re-fetch |
| `inStock=true` toggle | Only products with `availableQty > 0` render | URL `?inStock=true` → `catalog.product.list` filter |
| Request returns INTERNAL_ERROR | Toast with traceId; chips and grid keep last good state | Standard `ApiError` handling in fetch wrapper |

## Notes

- The Server Component reads `searchParams` from the route props; the Client filter pane uses `useRouter().replace(...)` to update them without losing scroll.
- No `'use client'` on the page itself — only on `SearchBar`, `Chip` group, `PriceRangeSheet`, and `Toggle`.
