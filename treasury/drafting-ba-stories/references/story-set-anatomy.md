# Story Set anatomy — per-story section sourcing + worked example

Deep guidance for filling each per-story `user-story-template` instance and a full worked ShopPilot example.
The cross-stage contract (output directory/envelope shape, the count invariant, deterministic naming/idempotency)
stays in `SKILL.md`; this file expands only *how to fill a story* and *what a filled story looks like*.

## Per-story file = a full `user-story-template` instance

Each `ST-NN-<slug>.md` follows the vendored template exactly. Where each section's content comes from:

| Template section | Required? | Source in the Scope Sheet (S1) |
|---|---|---|
| breadcrumb `< Back · <ISSUE-KEY> · <RELATED-KEY>` | optional | tracker key(s) if any, else drop the line |
| `# [<TYPE>] <title>` | **required** | one capability from `in_scope`; tag `[Feature]` / `[BAU]` / `[Spike]` |
| `## Description` (As a / I want / So that) | **required** | the capability + its value (`business_goal`) |
| `## Product` | optional | product / feature area named in scope |
| `## Background` | optional | `business_goal` + the trigger for this capability |
| `## Business Logic` (numbered rules, a./b. sub-items) | optional | rules the scope says the system must enforce |
| `### <Decision matrix>` table | optional | **only when the scope gives conditional logic** (rate / limit / state tables) — never invent one |
| `> ⚠️ Note:` callout | optional | a cross-cutting rule / effective date / ownership |
| `## Acceptance Criteria` | **kept** (pipeline override) | ≥1 testable Gherkin Given/When/Then with concrete values — the G2 gate requires it |
| `## Out of Scope` | optional | this story's slice of `out_of_scope`; name the owning team |
| `## Mock-up references` (AS-IS / TO-BE, Spec API) | optional | only if the scope references screens / APIs; else delete the section |
| METADATA (Attachments / Subtasks / Linked / Confluence / Activity) | optional | `<none>` placeholders unless mirroring a real ticket export |

Per the template, **only Title and Description are required; keep the sections you need and delete the rest** — except
Acceptance Criteria, which this pipeline keeps (the G2 gate is about testable criteria). Banking / compliance rigor is
**not** force-filled here: a story carries rigor through Business-Logic rules + decision tables + Acceptance Criteria;
standalone PII / compliance flags live in `checking-ba-governance` (S3). When the scope supports it a story reaches the
richness of `example_TierRateStory.md` (numbered Business-Logic rules, a `### <name>` decision table, a
`> ⚠️ Note`). When the scope gives no conditional logic, omit the table — do not invent one.

## Language

Story **content** follows the source language (the template states content may be written in any language). Pick the
scaffold to match: English content → `templates/user-story-template.md`; Thai content →
`templates/user-story-template-th.md`. The exemplar is Thai content on the English scaffold — that is allowed; match
the source. Do **not** emit parallel EN+TH copies of a story.

## Example (ShopPilot MVP)

`INDEX.md` carries the envelope + epic, then a story table whose first rows are:

```
| ST-ID | issue-key | title                                 | summary (1 line)                     | state   | file                                       |
|-------|-----------|---------------------------------------|--------------------------------------|---------|--------------------------------------------|
| ST-03 | <none>    | [Feature] Checkout with server-side totals | Pay exactly the server-computed total | drafted | ST-03-checkout-with-server-side-totals.md  |
| ST-05 | <none>    | [Feature] Stock reservation & last-item race | One winner on the last unit, at confirm | drafted | ST-05-stock-reservation-last-item-race.md |
```

`ST-03-checkout-with-server-side-totals.md` (abbreviated — a full `user-story-template` instance):

```
<!-- ba-pipeline: story-set REQ-shoppilot-mvp · ST-03 · state: drafted
     scope_ref: ../01-scope-sheet.md · open_questions: [] -->

`< Back`  ·  `<none>`

# [Feature] Checkout with server-side totals

## Description
*As a* logged-in customer
*I want* to review cart, address, coupon, shipping and net total before paying
*So that* I pay exactly what the system computed, not a tampered client value.

## Business Logic
1. Totals are computed **server-side** (§4.6); the client value is never trusted.
2. Shipping is free at ≥ 1,500 THB, else 60 THB (§6.3).
3. One coupon per order; the total is floored at 0 and never negative (§6.4).
4. The coupon is **re-validated at confirm** (§4.5).

## Acceptance Criteria
1. Given a 1,200 THB cart with coupon `WELCOME100`, when I confirm, then the total is 1,160 THB (1,200 − 100 + 60 shipping) and an order is created in `awaiting-payment`.
2. Given a 1,600 THB cart, when I confirm, then shipping is 0 (free ≥ 1,500).
3. Given a coupon that expired while I browsed, when I confirm, then the order is rejected with "coupon no longer valid" and no order is created.

## Out of Scope
1. Coupon stacking (Phase 2).
```

`ST-05` carries `open_questions: [OQ-6]` (the reservation-TTL question) in its locator stub — carried verbatim from
the Scope Sheet, not invented here.
