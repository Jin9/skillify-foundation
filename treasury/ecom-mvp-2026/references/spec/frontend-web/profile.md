# Route spec — `/profile` (profile tab)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_PROFILE`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_PROFILE.md)
- **Auth:** required
- **Type:** Server Component + Client logout button
- **Tab route:** yes

## Purpose

Customer's session-control surface — shows name, email, default address, and a "ออกจากระบบ" (logout) action. Logout calls `identity.logout` and always clears local cookies.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `identity.profile.read` | Name + email |
| `identity.address.list` (filter `isDefault=true` client-side) | Default address card |
| `POST /api/auth/logout` → `identity.logout` | Logout action |
| `identity.profile.update` | Optional inline edit (name/phone) — Frontend Spec §3 doesn't gate it; we expose it under a `แก้ไข` action |

## Key components

- `ShopHeader`
- `ProfileHeader` — name + email
- `DefaultAddressCard` — single card from address list
- `BtnLogout` (`'use client'`) — primary brand button
- `Link` → `/addresses` for full address book
- `TabBar`

## Mobile layout

- Top: avatar placeholder + name + email
- Middle: default-address card (link to `/addresses` to manage)
- Bottom: large "ออกจากระบบ" button above `TabBar`

## Acceptance criteria covered

- **AC1 (renders + logout):**
  - Renders name, email, default address (or empty state "ยังไม่มีที่อยู่จัดส่ง")
  - Tap "ออกจากระบบ" → POST `/api/auth/logout` → handler clears both cookies (Max-Age=0) → client clears `SessionContext` → `router.push('/')`

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Logout network failure | Local cookies still cleared client-side via `SessionContext` reset; effective logout; next request with stale token will be rejected when connectivity returns | Route handler always emits `Set-Cookie name=; Max-Age=0` regardless of upstream outcome |
| Customer has no default address | `DefaultAddressCard` shows empty state with CTA → `/addresses` | Filter `isDefault==true`; conditional render |

## Notes

- `BtnLogout` is the only Client Component on the page.
- The Server Component shell fetches profile + addresses in parallel (`Promise.all`).
