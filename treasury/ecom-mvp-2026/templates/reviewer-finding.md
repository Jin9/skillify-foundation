# Reviewer Finding Template (L1 + L2)

**Parent:** [`../templates.md`](../templates.md)
**Owner roles:** [Reviewer-L1](../roles.md), [Reviewer-L2](../roles.md) · **`template_version`:** 0.1.0

Same template for L1 and L2; the diagnosis tag distinguishes scope. Promoted from inline v0.1 used in dry-run #1.

**File locations:**
- Reviewer-L1: `<workflow_root>/.squad-run/components/<svc>/rev-l1.json`
- Reviewer-L2: `<workflow_root>/.squad-run/layer-2/rev-l2.json`

---

## Top-level shape

```json
{
  "template_version": "0.1.0",
  "review_kind": "reviewer-l1" | "reviewer-l2",
  "component_name": "<svc>",        // L1 only; L2 omits
  "skills_invoked": ["expert-software-security-reviewer", "simplify"],
  "verdict": "pass" | "advisory_only" | "block",
  "summary": "1-paragraph",
  "findings": [/* finding objects */]
}
```

`verdict` is for human-readability; the orchestrator routes mechanically on per-finding `severity` + `diagnosis_tag`.

## Finding object fields

- `id` — string, required, format `REV-<COMPONENT>-NNN` for L1, `REV-L2-NNN` for L2
- `severity` — enum: `high` | `medium` | `low`, required
- `diagnosis_tag` — enum, required:
  - **L1:** `code_security_issue` | `latent_bug` | `code_quality_issue` | `style_nit` | `spec_induced_security_issue` | `spec_induced_latent_bug`
  - **L2:** `cross_component_security_issue` | `cross_component_data_flow_issue` | `cross_component_error_propagation_issue` | `cross_component_maintainability_issue`
- `routes_to` — enum: `Dev` | `Tech-Designer` | `Tech-Lead` | `null` (low advisory only), required
- `locus` — string, required:
  - **L1:** `file:line` (e.g., `app/identity/service.go:80`)
  - **L2:** `component-pair (a ↔ b) + file:line if applicable`
- `description` — string, required, 1-3 sentences
- `recommendation` — string, required, one line
- `cwe_refs` — array of strings, optional: CWE / OWASP ASVS refs (e.g., `CWE-208`, `ASVS-V2.4.7`)

---

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "reviewer-finding/v0.1.0",
  "type": "object",
  "additionalProperties": false,
  "required": ["template_version", "review_kind", "verdict", "summary", "findings"],
  "properties": {
    "template_version": { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "review_kind":      { "enum": ["reviewer-l1", "reviewer-l2"] },
    "component_name":   { "type": "string" },
    "skills_invoked":   { "type": "array", "items": { "type": "string" } },
    "verdict":          { "enum": ["pass", "advisory_only", "block"] },
    "summary":          { "type": "string", "minLength": 1 },
    "findings": {
      "type": "array",
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["id", "severity", "diagnosis_tag", "routes_to", "locus", "description", "recommendation"],
        "properties": {
          "id":            { "type": "string", "minLength": 1 },
          "severity":      { "enum": ["high", "medium", "low"] },
          "diagnosis_tag": { "type": "string" },
          "routes_to":     { "anyOf": [{"enum": ["Dev", "Tech-Designer", "Tech-Lead"]}, {"type": "null"}] },
          "locus":         { "type": "string", "minLength": 1 },
          "description":   { "type": "string", "minLength": 1 },
          "recommendation":{ "type": "string", "minLength": 1 },
          "cwe_refs":      { "type": "array", "items": { "type": "string" } }
        }
      }
    }
  }
}
```

---

## Examples

### Example #1 — L1 high `code_security_issue`

```json
{
  "id": "REV-IDENT-001",
  "severity": "high",
  "diagnosis_tag": "code_security_issue",
  "routes_to": "Dev",
  "locus": "app/identity/service.go:80",
  "description": "Calls HashSha256EncodePepper for password hashing — unsalted base64(sha256(pw+pepper)) with no work factor; trivially crackable offline.",
  "recommendation": "Switch to bcrypt cost 12 via golang.org/x/crypto/bcrypt; migrate via lazy-rehash on next login.",
  "cwe_refs": ["CWE-916", "CWE-759", "ASVS-V2.4.1"]
}
```

### Example #2 — L2 high `cross_component_data_flow_issue`

```json
{
  "id": "REV-L2-001",
  "severity": "high",
  "diagnosis_tag": "cross_component_data_flow_issue",
  "routes_to": "Tech-Lead",
  "locus": "checkout ↔ order (handler_commit.go:200 vs handler_create_from_checkout.go:45)",
  "description": "Checkout's orderCreateRequest omits buyerEmail; Order's CreateFromCheckoutRequest declares buyerEmail required+email. Every checkout commit 400s at step 7.",
  "recommendation": "Either Checkout calls identity.profile.read to fetch email and forwards it, OR Identity adds email to JWT claims and Checkout extracts it from claims, OR Order makes buyerEmail optional and looks it up itself.",
  "cwe_refs": []
}
```

---

## Negative examples

### Negative #1 — Severity inflated for a style nit

```json
{
  "id": "REV-IDENT-099",
  "severity": "high",
  "diagnosis_tag": "code_security_issue",
  "routes_to": "Dev",
  "locus": "app/identity/handler_register.go:42",
  "description": "Variable name `req` is too short; should be `registerRequest`.",
  "recommendation": "Rename for clarity."
}
```

What's wrong:

1. Variable naming is a `style_nit` (low, advisory). Inflating to `code_security_issue` (high) burns a Dev cycle on a nit. Calibration violation.
2. Same as Plan-Reviewer's spurious-high calibration rule: ask "could a careful human have written this and been right?"

### Negative #2 — Missing `routes_to` on an L1 finding

```json
{
  "id": "REV-INV-005",
  "severity": "high",
  "diagnosis_tag": "latent_bug",
  "locus": "app/inventory/handler_reservation_create.go:363",
  "description": "Lock-order inversion vs sweeper — deadlock under concurrent load."
}
```

What's wrong:

1. `routes_to` missing — schema fail; orchestrator can't route. Latent bug high routes to Dev (cap 2).
2. Missing `recommendation`.
