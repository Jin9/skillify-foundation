# Worked example: one source, both modes

A short source extracted twice — once in **infer** mode (no schema), once in **conform** mode (caller supplies `target_schema`). Both emit one envelope that passes `scripts/validate_extraction.py`.

## Source (`policy-memo-2026-06-12.md`)

> Following last quarter's complaint volume, we propose raising the per-customer
> daily outbound transfer cap on retail savings accounts from 50,000 to 200,000 THB,
> effective 1 July 2026. The raise applies only to customers at KYC tier 2 or above
> whose account is at least 90 days old. Dormant accounts keep the 50,000 cap until
> reactivated. (See BOT notice 12/2566 for the governing limit framework.)

---

## Run A — infer mode

**Input:** `{ "source": "policy-memo-2026-06-12.md" }` (no `target_schema`).

The agent tags salience: the change, scope, date, conditions, and the dormant-account exception are must-preserve; the complaint-volume background is out-of-scope; the verbatim BOT citation is droppable (size-bounded, reachable via `source_ref`). It infers a flat shape:

```json
{
  "contract": {
    "title": "Daily transfer limit raise",
    "change": "Raise per-customer daily outbound transfer cap from 50000 to 200000 THB",
    "applies_to": "retail savings accounts",
    "effective_date": "2026-07-01",
    "conditions": ["KYC tier >= 2", "account age >= 90 days"],
    "exceptions": ["dormant accounts keep the 50000 cap until reactivated"]
  },
  "_meta": {
    "mode": "infer",
    "source_ref": "policy-memo-2026-06-12.md",
    "coverage": "high",
    "confidence": "high",
    "dropped": [
      { "item": "background on last quarter's complaint volume", "reason": "out-of-scope" },
      { "item": "verbatim BOT notice 12/2566 citation text", "reason": "size-bounded" }
    ],
    "provenance": [
      { "field": "exceptions", "from": "memo sentence 4 (dormant accounts)" }
    ],
    "notes": "Inferred shape. If a downstream limit-change contract exists, pass it as target_schema."
  }
}
```

Summary line: `mode=infer coverage=high dropped=2 confidence=high`

---

## Run B — conform mode

**Input:** the same `source`, plus a downstream `target_schema`:

```json
{
  "type": "object",
  "required": ["limit_name", "old_value", "new_value", "currency", "effective_date", "approver"],
  "properties": {
    "limit_name": { "type": "string" },
    "old_value": { "type": "number" },
    "new_value": { "type": "number" },
    "currency": { "type": "string" },
    "effective_date": { "type": "string" },
    "approver": { "type": "string" }
  }
}
```

The agent maps source → schema. The source supports every field **except `approver`** — the memo names no approver. The rule is: do not fabricate; record the gap.

```json
{
  "contract": {
    "limit_name": "per-customer daily outbound transfer cap (retail savings)",
    "old_value": 50000,
    "new_value": 200000,
    "currency": "THB",
    "effective_date": "2026-07-01"
  },
  "_meta": {
    "mode": "conform",
    "source_ref": "policy-memo-2026-06-12.md",
    "coverage": "partial",
    "confidence": "high",
    "dropped": [
      { "item": "approver", "reason": "not-in-source" },
      { "item": "KYC tier >= 2 and 90-day age conditions; dormant-account exception", "reason": "out-of-scope" }
    ],
    "provenance": [
      { "field": "old_value", "from": "memo: 'from 50,000'" },
      { "field": "new_value", "from": "memo: 'to 200,000 THB'" }
    ],
    "notes": "Required field 'approver' is not present in the source — left absent, not invented. Caller must supply or source it before downstream use."
  }
}
```

Summary line: `mode=conform coverage=partial dropped=2 confidence=high`

**Key contrasts:**
- Run B left a *required* field (`approver`) absent and flagged it `not-in-source` rather than fabricating it.
- The conditions/exception that were first-class fields in the inferred contract are out-of-scope for *this* target schema, so they move to `dropped` (still reason-tagged, still recoverable via `source_ref`).
- Both envelopes pass `validate_extraction.py`; neither drops anything silently.
