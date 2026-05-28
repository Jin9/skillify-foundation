# ADR-010 — Mobile-first Thai customer-only frontend (no admin UI in MVP)

- **Status:** Accepted (locked by BA in_scope + Frontend Spec scope)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

The BA `in_scope` lists 14 customer-facing frontend lines and the `out_of_scope` block explicitly rejects an admin UI ("Frontend Spec is customer-facing only; Admin remains backend-only via Postman/curl in MVP"). The Frontend Spec specifies a 390x844 mobile viewport, Thai (`lang="th"`), IBM Plex Sans Thai typography, emerald `#1a8060` brand colour, 5 bottom-tab routes, and 8 full-screen routes.

This is a deliberate scope reduction: shipping ONE polished customer surface in MVP beats shipping a half-finished admin SPA alongside a half-finished customer SPA.

## Decision

`frontend-web` is a Next.js 14 (App Router) + TypeScript + Tailwind app. It serves ONLY the customer journey. Admin endpoints exist and are tested via curl/Postman.

Key constraints:

- **Single viewport target:** 390x844 mobile-first; the layout reflows down to ~360 (Android Chrome) without horizontal scroll. There is no separate desktop layout in MVP.
- **Thai-only:** `lang="th"` on every page; copy in Thai; Buddhist-era date formatting where date is rendered to the customer; THB currency with `฿` prefix (no fractional THB in MVP). The data model does NOT preclude later i18n (no Thai-only enums in API contracts; only in display copy).
- **5 bottom-tab routes:** `home`, `search`, `cart`, `orders`, `profile`. The bottom tab bar appears ONLY on these 5 routes.
- **8 full-screen routes:** `detail`, `login`, `checkout`, `payment`, `orderSuccess`, `orderResult`, `orderDetail`, `addresses` — no bottom tab bar.
- **Auth gating:** `cart`, `orders`, `profile`, `checkout`, `payment`, `orderDetail`, `addresses` are auth-required; an unauthenticated visit redirects to `login` and after successful login returns to the originally requested route. `home`, `search`, `detail` are public.
- **Cookie-only token storage:** the Next.js route handlers act as the gateway and own the 6-step `withAuth` flow per `cross-cutting.auth.token_acquisition_flow`. The frontend SPA NEVER sees the JWT bytes.
- **No internal-secret access:** `frontend/.env.example` explicitly omits `INTERNAL_SHARED_SECRET`; the proxy forwards `Authorization` but NEVER `X-Internal-Secret`. A regression test asserts the absence of `X-Internal-Secret` on customer-originated requests (REV-L2-013 verified-correct posture).
- **Idempotency-Key on checkout commit:** generated client-side as a UUID v4 and persisted in `localStorage` keyed per pending checkout so a refresh-and-retry does not create a second order.
- **Server-driven totals on checkout:** the checkout route renders `subtotal`, `shippingFee`, `total` exactly as returned by `checkout.preview` / `checkout.commit`. The client NEVER recomputes the §11.2 free-shipping rule.
- **Status badges in Thai:** order status enums (`PENDING_PAYMENT`, `PAID`, `PACKING`, `SHIPPED`, `DELIVERED`, `CANCELLED`, `PAYMENT_FAILED`, `PAYMENT_EXPIRED`) map to Thai labels in the UI; the underlying API contract remains English enums.

## Consequences

- **Positive:** the team ships ONE polished customer experience instead of two half-built ones.
- **Positive:** admin interactions are fully testable via curl/Postman against the documented contracts; no admin UI bugs to triage.
- **Negative:** an admin UI WILL be needed post-MVP; deferring it does not delete the work. The contracts.json admin endpoints are designed to be UI-friendly (no quirky binary fields, all enums documented) so the post-MVP UI build is straightforward.
- **Negative:** Thai-only copy means a future i18n pass needs a copy-extraction sweep across all routes. The data model is i18n-clean (no Thai-only enum values reach the API) so this stays a frontend-only concern.
- **Negative:** mobile-first means the customer cannot complete the journey on a desktop without acceptable degradation; the constraint is that the layout reflows to standard mobile widths, not that desktop is unusable. Frontend Spec §1 covers 360-wide reflow; wider viewports get the same single-column layout centred.
