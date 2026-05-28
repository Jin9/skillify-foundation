# Route spec — `/addresses` (address book + picker mode)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_ADDRESSES`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_ADDRESSES.md)
- **Auth:** required
- **Type:** Client Component (CRUD + picker)
- **Tab route:** no

## Purpose

Two-mode route:

1. **Manage mode** (`/addresses`) — list, add, edit, delete, set-default.
2. **Picker mode** (`/addresses?mode=pick&returnTo=/checkout`) — same list with a "เลือกที่อยู่นี้" affordance per card; on tap, `router.push(returnTo + '?addressId=' + id)`.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `identity.address.list` | List |
| `identity.address.create` | Add new (Sheet form) |
| `identity.address.update` | Edit existing (Sheet form) |
| `identity.address.delete` | Delete |
| `identity.address.set-default` | Set as default |

## Key components

- `ShopHeader` (back)
- `AddressCard` — receiver, phone, address line, district, province, postalCode + default badge + edit/delete icons
- `BtnAddAddress` — primary brand; on tap opens `Sheet` with `AddressForm`
- `AddressForm` (RHF + zod):
  ```ts
  z.object({
    receiverName: z.string().min(1).max(120),
    phone: z.string().regex(/^[+]?[0-9]{9,15}$/),
    addressLine: z.string().min(1).max(200),
    district: z.string().min(1).max(80),
    province: z.string().min(1).max(80),
    postalCode: z.string().regex(/^\d{5}$/),
    isDefault: z.boolean().optional(),
  })
  ```
- `PickerSelectButton` — visible only in picker mode

## Mobile layout

- One card per address with default badge
- Sticky `BtnAddAddress` at the bottom (above iOS safe area)
- Bottom Sheet for add/edit form

## Acceptance criteria covered

- **AC1 (list + default + add CTA):** Each card shows the required atoms; default badge marks current default; `BtnAddAddress` is brand emerald.
- **AC2 (picker returns to checkout):** Picker tap → `router.push('/checkout?addressId=' + id)`; `/checkout` re-runs `checkout.preview` with the new `shippingAddressId`, totals recalculate.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Form submit with missing/invalid field | Field-level Thai error per zod resolver; envelope `VALIDATION_ERROR` if backend disagrees | RHF + zod resolver |
| Picker mode with zero saved addresses | Open `Sheet` with `AddressForm` immediately; on save, set as default and return to `/checkout?addressId=...` | Conditional render based on `addresses.length === 0 && mode === 'pick'` |
| Edit form `isDefault` toggle | Per TD ambiguity #4: address-edit form does NOT include `isDefault`; users must use "ตั้งเป็นค่าเริ่มต้น" affordance which calls `identity.address.set-default` | Form schema omits `isDefault` |

## Notes

- Picker mode uses URL `?mode=pick&returnTo=...`; back button respects the `returnTo` so customers don't get stranded.
- Delete asks for confirmation (Thai dialog: "ลบที่อยู่นี้?") to prevent accidental loss.
