# EPIC Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [BA](../roles.md) · **`template_version`:** 0.1.0

An EPIC is a coarse business outcome that decomposes into stories. The BA emits one EPIC MD per coarse outcome named in the source requirement. Tech-Lead's component decomposition typically maps several stories of one or more EPICs to each microservice.

**File location:** `<workflow_root>/requirement/<EPIC_SLUG>/<EPIC_SLUG>.md`
**ID convention:** `EPIC_<DOMAIN>` (uppercase, underscore-separated). Examples: `EPIC_AUTH`, `EPIC_CATALOG`, `EPIC_CHECKOUT`. Same-domain epics may add a discriminator: `EPIC_AUTH_CUSTOMER` vs `EPIC_AUTH_ADMIN`. Domain slug is derived from the requirement's bounded contexts, not invented.

---

## Fields

- `epic_id` — string, required, format `^EPIC_[A-Z][A-Z0-9_]*$`
- `title` — single-line, required
- `summary` — 2-3 sentences, required, captures the business outcome (the "why")
- `business_value` — string, required, ties to a measurable goal from the requirement (revenue, retention, compliance, latency, cost)
- `in_scope_stories` — list of `STORY_<SLUG>` IDs, required (≥1)
- `out_of_scope` — list of explicit non-goals, required (empty array allowed but suspicious — same smell as BA's `out_of_scope`)
- `success_metrics` — list of measurable signals, required (≥1)
- `dependencies` — list of `EPIC_<SLUG>` IDs this epic blocks on, optional
- `compliance_sensitive` — boolean, required (mirrors BA's flag; affects downstream fan-out tier)
- `change_log` — table, required (≥1 entry on creation)

---

## File shape

```markdown
---
epic_id: EPIC_AUTH
title: Customer authentication and session management
summary: |
  Customers must register, log in, refresh sessions, and log out securely. This is foundational
  to every other epic (catalog browsing is the only customer-facing capability allowed without auth).
business_value: Unblocks every authenticated path in the platform. Without auth, no checkout, no orders, no admin.
in_scope_stories:
  - STORY_AUTH_REGISTER
  - STORY_AUTH_LOGIN
  - STORY_AUTH_REFRESH
  - STORY_AUTH_LOGOUT
out_of_scope:
  - Social login (Google/Facebook/Apple) — deferred to v2
  - Multi-factor authentication — deferred to v2
  - Password reset via email — deferred to v2 (no email infrastructure in MVP)
success_metrics:
  - Login success rate ≥ 99% on valid credentials
  - p95 login latency < 300ms
  - Refresh token rotation has zero replay-windows ≥ 5s in production logs
dependencies: []
compliance_sensitive: false
---

# EPIC_AUTH — Customer authentication and session management

## Why

…2-3 paragraphs of plain-language context. What customer behavior does this enable?
What's the business motivation? Reference the requirement section that's the source of truth.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_AUTH_REGISTER` | Register with email + password | Must |
| `STORY_AUTH_LOGIN` | Log in with email + password | Must |
| `STORY_AUTH_REFRESH` | Refresh expired access token | Must |
| `STORY_AUTH_LOGOUT` | Revoke refresh token | Should |

## Out-of-scope rationale

For each item in `out_of_scope` above, one sentence on why it's out (cost, dependency, scope reduction).

## Risks

- Password hashing must use a memory-hard or work-factor algorithm (bcrypt cost 12 minimum).
- Refresh token rotation must be atomic — race in the rotation tx allows replay.
- Generic login error must not leak whether the email exists (AUTH-007).

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-07 | BA (Claude Opus 4.7 xHigh) | Created from requirement §8.1, §8.2 |
```

---

## Schema (validated by orchestrator before fan-out)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "epic-output/v0.1.0",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "template_version", "epic_id", "title", "summary", "business_value",
    "in_scope_stories", "out_of_scope", "success_metrics", "dependencies",
    "compliance_sensitive", "change_log"
  ],
  "properties": {
    "template_version": { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "epic_id":          { "type": "string", "pattern": "^EPIC_[A-Z][A-Z0-9_]*$" },
    "title":            { "type": "string", "minLength": 1, "maxLength": 120 },
    "summary":          { "type": "string", "minLength": 1, "maxLength": 600 },
    "business_value":   { "type": "string", "minLength": 1 },
    "in_scope_stories": {
      "type": "array", "minItems": 1,
      "items": { "type": "string", "pattern": "^STORY_[A-Z][A-Z0-9_]*$" }
    },
    "out_of_scope":     { "type": "array", "items": { "type": "string", "minLength": 1 } },
    "success_metrics":  { "type": "array", "minItems": 1, "items": { "type": "string", "minLength": 1 } },
    "dependencies":     { "type": "array", "items": { "type": "string", "pattern": "^EPIC_[A-Z][A-Z0-9_]*$" } },
    "compliance_sensitive": { "type": "boolean" },
    "change_log": {
      "type": "array", "minItems": 1,
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["date", "author", "change"],
        "properties": {
          "date":   { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
          "author": { "type": "string", "minLength": 1 },
          "change": { "type": "string", "minLength": 1 }
        }
      }
    }
  }
}
```

The frontmatter is YAML for readability; the orchestrator parses it as JSON-equivalent and validates against the schema above.

---

## Negative examples

These pass YAML/JSON syntax. Plan-Reviewer should reject them.

### Negative #1 — Vague EPIC, missing measurability

```yaml
epic_id: EPIC_CHECKOUT
title: Make checkout work
summary: Customer should be able to check out.
business_value: Important for the business.
in_scope_stories: [STORY_CHK_DO_IT]
out_of_scope: []
success_metrics:
  - Checkout works
dependencies: []
compliance_sensitive: false
change_log:
  - {date: 2026-05-07, author: BA, change: Created}
```

What Plan-Reviewer should catch:

1. `summary` paraphrases the title — Tech-Lead has nothing to architect against. Tag: `requirements_gap` (high).
2. `business_value` is content-free ("important for the business"). Tag: `requirement_ambiguity` (medium).
3. `in_scope_stories` has one entry that isn't a real story slug (`STORY_CHK_DO_IT`). Tag: `requirements_gap` (high).
4. `success_metrics` is unfalsifiable ("Checkout works"). Tag: `requirement_ambiguity` (high).
5. `out_of_scope: []` — for a checkout epic, no scope fence. Tag: `unstated_assumption` (medium).

### Negative #2 — EPIC inventing stories outside the requirement

```yaml
epic_id: EPIC_CATALOG
title: Catalog browsing and management
summary: Customers browse products; admins manage inventory and AI-assisted product recommendations.
business_value: Drive conversion via personalized recommendations.
in_scope_stories:
  - STORY_CATALOG_LIST
  - STORY_CATALOG_DETAIL
  - STORY_CATALOG_AI_RECOMMENDATIONS
  - STORY_CATALOG_PERSONALIZED_HOME
  - STORY_CATALOG_LOYALTY_BADGE
out_of_scope: []
success_metrics:
  - Recommendation CTR > 8%
dependencies: []
compliance_sensitive: false
change_log:
  - {date: 2026-05-07, author: BA, change: Created}
```

What Plan-Reviewer should catch:

1. `STORY_CATALOG_AI_RECOMMENDATIONS` and `STORY_CATALOG_PERSONALIZED_HOME` reference behaviors the requirement explicitly defers (§4.2 lists "AI product recommendation" as out of scope). Tag: `unstated_assumption` (high) → BA.
2. `STORY_CATALOG_LOYALTY_BADGE` references the loyalty system also explicitly deferred. Tag: `unstated_assumption` (high).
3. `out_of_scope: []` — exactly the scope-creep smell `ba.md`'s negative example #2 warns about, surfaced at the EPIC level. Tag: `missing_edge_case` (medium).
4. `success_metrics` measures a metric (CTR) on a behavior that isn't shipped — unmeasurable in this run. Tag: `requirement_ambiguity` (high).

Expected routing in both cases: Plan-Reviewer → BA, plan-stage cap 1, then HardFail.
