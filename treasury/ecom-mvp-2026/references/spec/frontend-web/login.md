# Route spec — `/login` (login + register toggle)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_AUTH`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_AUTH.md)
- **Auth:** public
- **Type:** Client Component (form)
- **Tab route:** no

## Purpose

Single route for login + register via a `ToggleTabs` switcher. On successful login the route handler sets HttpOnly access+refresh cookies and the SPA navigates back to `?next=...` or `/`.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `POST /api/auth/login` → `identity.login` | Login submission |
| `POST /api/auth/register` → `identity.register` | Register submission |

## Key components

- `ShopHeader` (back button)
- `ToggleTabs` (login | register) — clears form on toggle (edge case)
- `Field` x N (email, password, name)
- `Btn` (brand, full-width)
- `ErrorBanner` (envelope code + Thai message)

## Form schema (zod)

```ts
const loginSchema = z.object({
  email: z.string().email().max(254),
  password: z.string().min(8).max(72),
})
const registerSchema = z.object({
  email: z.string().email().max(254),
  password: z.string().min(8).max(72),
  name: z.string().min(1).max(120),
})
```

(per `identity.register.payload_shape` and `identity.login.payload_shape`)

## Mobile layout

- Centered single column, max-width 480 (capped) but in 390 column it just fills
- Brand-emerald submit at bottom

## Acceptance criteria covered

- **AC1 (login flow):** valid email+password → `POST /api/auth/login` → handler sets access+refresh cookies + returns `{user}` → client `router.push(next || '/')`.
- **AC2 (register flow):** valid email+password+name → `POST /api/auth/register` → on success: switch the toggle to `login` mode (BA preference per AUTH-001 — register does NOT auto-login) and surface a Thai success message; then user submits login.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Wrong password | Generic AUTH_INVALID with locked phrasing 'อีเมลหรือรหัสผ่านไม่ถูกต้อง' (Thai equivalent of 'Invalid email or password.' per AUTH-007 + cross-cutting.error-codes.ba_flagged_resolution) | Form `ErrorBanner` from envelope |
| Duplicate email register | `DUPLICATE_EMAIL` 409 → Thai 'อีเมลนี้ถูกใช้แล้ว' | `error_handling.ui_mapping` |
| Toggle clears form | `react-hook-form` `reset()` on toggle | `ToggleTabs.onChange` |

## Notes

- `?next=` query param is preserved through the toggle and applied on successful login.
- The route handler returns ONLY `{user}` to the client — raw access/refresh tokens never reach JavaScript.
