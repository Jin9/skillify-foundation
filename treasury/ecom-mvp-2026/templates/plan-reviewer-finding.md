# Plan-Reviewer Finding Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [Plan-Reviewer](../roles.md) · **`template_version`:** 0.1.0

The output Plan-Reviewer emits at the plan-stage gate, BEFORE fan-out. Promoted from inline v0.1 used in dry-run #1.

**File location:** `<workflow_root>/.squad-run/plan-review.json` (and `plan-review-r2.json` after one repair cycle).

---

## Top-level shape

```json
{
  "template_version": "0.1.0",
  "verdict": "pass" | "advisory_only" | "block",
  "summary": "1-paragraph: state of the plan",
  "findings": [/* finding objects, see below */]
}
```

`verdict`:
- `pass` — no high or medium findings; fan-out proceeds immediately.
- `advisory_only` — only `low` (`plan_polish`) findings; lands in `PLAN_NOTES.md`; fan-out proceeds.
- `block` — at least one high or medium; routes back per per-finding `routes_to`. Plan-stage cap is 1; second `block` ⇒ `HardFail`.

## Finding object fields

- `id` — string, required, `PR-NNN`
- `severity` — enum: `high` | `medium` | `low`, required
- `diagnosis_tag` — enum (per `../orchestrator.md` routing table), required:
  - `requirements_gap` (high) → BA
  - `requirement_ambiguity` (high) → BA
  - `missing_edge_case` (high|medium) → BA
  - `architecture_risk` (high) → Tech-Lead
  - `contract_ambiguity` (high) → Tech-Lead
  - `unstated_assumption` (high|medium) → BA OR Tech-Lead (specify)
  - `complexity_misjudged` (medium) → Tech-Lead
  - `plan_polish` (low) → no route, advisory
- `routes_to` — enum: `BA` | `Tech-Lead` | `null` (low only), required
- `locus` — string, required: artifact + field/section path (e.g., `ba.json#/acceptance_criteria[12]`, `contracts.json#/contracts[name=payment.callback]/idempotency_rules`)
- `description` — string, required, 1-3 sentences (problem only, not solution)
- `recommendation` — string, required, one line (concrete fix; producer decides if/how to act)

---

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "plan-review/v0.1.0",
  "type": "object",
  "additionalProperties": false,
  "required": ["template_version", "verdict", "summary", "findings"],
  "properties": {
    "template_version": { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "verdict":          { "enum": ["pass", "advisory_only", "block"] },
    "summary":          { "type": "string", "minLength": 1 },
    "findings": {
      "type": "array",
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["id", "severity", "diagnosis_tag", "routes_to", "locus", "description", "recommendation"],
        "properties": {
          "id":            { "type": "string", "pattern": "^PR-\\d{3,}$" },
          "severity":      { "enum": ["high", "medium", "low"] },
          "diagnosis_tag": { "enum": ["requirements_gap", "requirement_ambiguity", "missing_edge_case", "architecture_risk", "contract_ambiguity", "unstated_assumption", "complexity_misjudged", "plan_polish"] },
          "routes_to":     { "anyOf": [{"enum": ["BA", "Tech-Lead"]}, {"type": "null"}] },
          "locus":         { "type": "string", "minLength": 1 },
          "description":   { "type": "string", "minLength": 1 },
          "recommendation":{ "type": "string", "minLength": 1 }
        }
      }
    }
  }
}
```

---

## Example (real, from dry-run #1)

```json
{
  "id": "PR-002",
  "severity": "high",
  "diagnosis_tag": "requirements_gap",
  "routes_to": "BA",
  "locus": "ba.json#/in_scope (checkout pricing entries) and ba.json#/non_functional vs requirement §11.2 'MVP Shipping Fee Rule'",
  "description": "Requirement §11.2 specifies a tiered shipping fee — 60 THB when subtotal < 1500, 0 THB otherwise — as part of Pricing Rules, separate from §8.11 Shipping Mock (which BA correctly deferred). BA's checkout in_scope only states 'recomputes subtotal and total server-side' and does not capture the fee rule. Tech-Lead, with no BA signal, hardcoded shippingFee:'number (0 in MVP)' in checkout.preview. The resulting grandTotal will be wrong for any subtotal under 1500 — a directly observable functional regression against the requirement.",
  "recommendation": "Add an in_scope entry and a dedicated AC for §11.2: subtotal<1500 ⇒ shippingFee=60, subtotal≥1500 ⇒ shippingFee=0, grandTotal=subtotal+shippingFee in this run (no coupon)."
}
```

---

## Negative examples

### Negative #1 — Spurious `high` over a defensible call

```json
{
  "id": "PR-099",
  "severity": "high",
  "diagnosis_tag": "architecture_risk",
  "routes_to": "Tech-Lead",
  "locus": "contracts.json#/contracts[name=cart.add-item]",
  "description": "Cart uses pgx instead of an ORM.",
  "recommendation": "Switch to an ORM."
}
```

What's wrong:

1. `pgx vs ORM` is a defensible choice, not an architecture risk; firing `high` here burns the plan-stage cycle for nothing. Plan-Reviewer's calibration rule applies: "could a careful human producer reasonably have written this and been right?" Yes → downgrade to `low` `plan_polish` or skip.
2. `recommendation` is prescriptive on tech-stack rather than diagnostic.

### Negative #2 — Missing `routes_to` on an `unstated_assumption`

```json
{
  "id": "PR-100",
  "severity": "high",
  "diagnosis_tag": "unstated_assumption",
  "locus": "ba.json#/in_scope",
  "description": "BA added auto-merge behavior the user didn't ask for.",
  "recommendation": "Remove auto-merge from in_scope."
}
```

What's wrong:

1. `routes_to` is required for `unstated_assumption` (per orchestrator routing table — could route to BA OR Tech-Lead). Schema fail; orchestrator can't route. Tag this finding's structure as broken (Plan-Reviewer self-correction needed).
