---
name: implement-frontend-feature
version: 1.0.0
description: >
  Generate production-grade React/TypeScript code for one frontend feature from
  an approved UI design, with banking-grade discipline: WCAG 2.1 AA a11y, no
  any outside parsers, no localStorage auth tokens, no unsanitized
  dangerouslySetInnerHTML, design-token-only styling, generated API types, and
  PII field-level handling. Use when implementing a React component from an
  approved UI design. Use when generating a Next.js page with server state and
  forms from a spec. Use when implementing a feature module with TypeScript
  strict mode and a11y guarantees. Do NOT use for backend code (use
  implement-backend-feature). Do NOT use for pure CSS or styling tweaks
  without component logic. Do NOT use for greenfield architecture decisions
  (defer to design-frontend-feature). Do NOT use for infrastructure or build
  tooling. Do NOT use for fixing existing UI bugs (use generate-frontend-fix).
stage_type: generate
status: ready-for-phase-6
input_schema: schemas/input.json
output_schema: schemas/output.json
banking_grade: {idempotent: true, reversible: soft, audit_level: detailed}
expected_duration_p95_seconds: 180
max_retries_recommended: 2
compatibility: claude-code, codex, opencode
---

# Implement Frontend Feature

## Purpose

Take a single approved frontend design and emit production-grade React /
TypeScript code plus companion tests for one Generate stage in the frontend
Dev workflow (`COGNITIVE_OS.md` Phase 6+). Owns code synthesis only —
analysis, design, review, and validation live in sibling atomic skills.
Banking-grade non-negotiables apply on every output: WCAG 2.1 AA a11y,
TypeScript strict, no auth token in `localStorage`, no unsanitized HTML
injection, design-token-only styling, generated API types, PII field-level
treatment, analytics events at user-significant actions.

## When to use this skill

- Use when: implementing a React component, page, or feature module from an
  approved UI design.
- Use when: a workflow stage of `type: generate` selects this skill for a
  TypeScript / React target package.
- Use when: a design specifies a11y, PII, and token-storage requirements and
  the next stage is code emission.
- Do NOT use when: the target is backend, infra, or pure CSS.
- Do NOT use when: greenfield architecture decisions are still open — defer
  to `design-frontend-feature`.
- Do NOT use when: fixing an existing UI bug — defer to `generate-frontend-fix`.
- Do NOT use when: optimizing performance of existing code — defer to
  `analyze-frontend-performance`.
- Do NOT use when: reviewing already-emitted code — defer to
  `review-frontend-code`.

## Input

Input MUST validate against `schemas/input.json`. Required fields:

| Field | Type | Notes |
|-------|------|-------|
| `design_document` | string (markdown) OR object | Approved design. Must declare L1 intent, L2 component tree + state ownership + rendering model, L3 a11y contract + error/loading/empty states + token-storage strategy, L4 file sketch. |
| `target_feature_path` | string | Existing path under `src/` or `app/` (e.g., `app/loan/application/[id]`). |
| `idempotency_key` | string (UUID v4) | Same key → same output bytes. |
| `a11y_requirements` | object `{wcag_level, keyboard, screen_reader}` | Minimum WCAG level (default `AA`); blocking — no override below AA. |
| `pii_field_classification` | array of `{field, treatment}` | Per-field rule: `mask` / `redact` / `audit-on-view` / `none`. Empty array only if zero PII rendered. |

Optional: `convention_overrides` (object), `bundle_budget_kb` (number).

Reject the input with `loop_back` if any required field is missing, malformed,
or if the design has open a11y / security / PII questions.

## Procedure

Run all 9 steps in order. Do not skip steps for "small" features.

1. **Pre-flight — design completeness.** Verify the design declares L1 / L2 /
   L3 / L4 AND has no `TBD` on a11y, security, PII, or token storage. Any
   `TBD` on those four = `loop_back` to design.
2. **Pre-flight — target path exists.** Read `target_feature_path`. If
   missing, `loop_back`. If present, discover conventions: component layer
   structure, hook patterns, TypeScript settings, state libraries, test
   stack, token/theme system. Discovery outranks the greenfield defaults in
   `references/react-typescript-conventions.md`.
3. **Inspect — repo-first.** Read at least 2 sibling components / hooks /
   types in the same area. Mirror their boundary, naming, and import
   layout. Conflicts with design → `uncertainty_flag`, prefer repo.
4. **Plan — pillars.** Classify every emitted file into one pillar:
   `Page` / `Feature` / `Primitive` / `Hook` / `Type` / `Util` per
   `references/react-typescript-conventions.md`. A Primitive may not own
   business logic. A Page owns fetching / boundaries / layout shell.
5. **Plan — state ownership.** Map every state piece to one owner:
   `server` (TanStack Query or repo equivalent), `client` (Zustand /
   repo equivalent), `URL` (search params / router), `form` (RHF / repo
   equivalent), `local` (`useState`), or `derived` (compute on render).
   Never mirror `server` into `client`.
6. **Generate — code.** Emit files. Every emitted file MUST satisfy
   `references/implementation-rules.md`, `react-typescript-conventions.md`,
   `a11y-checklist.md`, `security-checklist.md`, and
   `state-management-rules.md`. Specifically:
   - Type safety: no `any` outside parser / boundary code; use `unknown` and
     narrow with type predicates / `satisfies` / `as` only at parsers.
   - State: no `useEffect` to sync server state into client state.
   - A11y: semantic HTML before ARIA; every input labeled; tab order matches
     visual order; focus moved correctly on route change / modal open / close;
     contrast 4.5:1 normal / 3:1 large.
   - Security: no `dangerouslySetInnerHTML` without an explicit `DOMPurify`
     call and a `// SAFE:` comment naming the threat model; no auth token
     write to `localStorage` / `sessionStorage`; URL props validated against
     scheme allowlist; `target="_blank"` always with `rel="noopener noreferrer"`.
   - PII: every field in `pii_field_classification` rendered through its
     declared treatment helper; never logged to console or analytics.
   - Tokens: every color / spacing / radius / shadow / type / motion read
     from the repo's design tokens — no hex literals, no arbitrary `[437px]`
     values. Missing token → `uncertainty_flag` of kind `token_gap`.
   - API types: imported from `codegen/` (or repo equivalent) — never
     hand-rolled `Response` types beyond a parser-local adapter.
   - Analytics: every user-significant action (submit, navigate, toggle that
     changes persisted state) emits an event whose `event_type` appears in
     `audit_events_emitted` output field.
7. **Generate — tests.** Companion tests per `react-typescript-conventions.md`
   § Tests: Vitest + RTL for components/hooks, MSW for network mocks. Tests
   query by role / label, not by class / test-id. Coverage `>= test_coverage_target`
   (default 0.80) per file. Failing coverage = `loop_back` (over-scoped).
8. **Self-review.** Apply every YES/NO in `references/self-review-checklist.md`.
   Any NO halts emission and routes per the Failure Modes table.
9. **Emit.** Produce a payload conforming to `schemas/output.json`. Populate
   nested `a11y_compliance` and `security_review` objects in full — they are
   required, not best-effort.

## Output Contract

Output MUST validate against `schemas/output.json`. Structured fields:

| Field | Type | Purpose |
|-------|------|---------|
| `files_generated` | array of `{path, content_hash, lines_added, component_pillar}` | Files written. `component_pillar` ∈ `Page | Feature | Primitive | Hook | Type | Util`. |
| `tests_generated` | array of `{path, content_hash, test_type, coverage_pct}` | `test_type` ∈ `unit | component | integration | e2e | visual | a11y`. |
| `a11y_compliance` | nested object — REQUIRED `{wcag_level, keyboard_navigable, screen_reader_tested, color_contrast_verified, focus_management_implemented, axe_clean}` | Every sub-field present, no defaults. |
| `security_review` | nested object — REQUIRED `{xss_surfaces, csrf_protected, csp_compliant, token_storage_strategy, pii_fields_handled}` | `xss_surfaces` MUST be empty or each entry has a `mitigation`. `token_storage_strategy` ∈ `httpOnly-cookie | in-memory | n/a`. |
| `state_ownership` | map `state_piece → server | client | URL | form | local | derived` | Every state piece named in the design must appear. |
| `bundle_impact_estimate_kb` | number | Best-effort estimate. Over `bundle_budget_kb` triggers `uncertainty_flag`. |
| `compensating_actions` | array | Required for any mutation path (optimistic UI rollback, undo affordance, idempotent retry). Empty for read-only. |
| `audit_events_emitted` | array of event-type strings | Every user-significant action contributes one entry. Empty only for display-only components. |
| `uncertainty_flags` | array of `{kind, location, note}` | Non-empty triggers downstream `loop_back`. |
| `decision_metadata` | object `{pillar_choices, state_library_choices, repo_conventions_followed}` | For audit. |

## Failure Modes

| Failure | Detection | Recovery |
|---------|-----------|----------|
| Design has unresolved a11y / security / PII / token-storage question | Step 1 | `loop_back` to design with `uncertainty_flags` populated |
| Target path missing or unreadable | Step 2 | `loop_back` to design (component path undecided) |
| Repo convention conflicts with design pattern | Step 3 | Prefer repo, emit `uncertainty_flag` of kind `convention_conflict` |
| Cannot meet WCAG `AA` (e.g., design requires gradient text under 4.5:1) | Step 6 a11y | `loop_back` to design — a11y is a BLOCKER, not a warning |
| Cannot generate component test for an interaction | Step 7 | `loop_back` to design (component is over-scoped or untestable as drawn) |
| Generated code introduces an XSS / CSRF / token-storage / PII-leak surface without mitigation | Step 8 self-review | `human-queue` — security NO is never `loop_back` |
| Bundle impact exceeds `bundle_budget_kb` | Step 9 | `uncertainty_flag` of kind `bundle_overrun`, verdict deferred to Review stage |
| Self-review finds NO on type-safety / state / a11y / security | Step 8 | First NO on security / PII: `human-queue`; everything else: `loop_back` |

## Anti-Patterns

- DO NOT use `any` type outside parser / boundary code — use `unknown` and narrow.
- DO NOT mirror server state into client state (`useEffect` syncing TanStack Query data into Zustand is the canonical mistake).
- DO NOT write an auth token to `localStorage` or `sessionStorage` — both are XSS-readable. Use HttpOnly cookie or in-memory.
- DO NOT use `dangerouslySetInnerHTML` without a `DOMPurify` call AND a `// SAFE:` comment naming the threat model.
- DO NOT skip a11y attributes (`aria-*`, `role`, semantic tags); do NOT remove focus outline without replacement; do NOT use placeholder as label.
- DO NOT introduce new npm dependencies that are not in the design — emit `dependency_addition` uncertainty flag instead.
- DO NOT generate a component without a companion test in the same payload.
- DO NOT log PII to console, error reporters, or analytics — render through the field's declared treatment helper only.
- DO NOT bypass the repo's design token system — no hex literals, no `[437px]` arbitrary values, no inline `style={{color: ...}}` for reachable tokens.
- DO NOT hand-roll API client types — import from `codegen/` (or repo equivalent); the only exception is a parser-local adapter.
- DO NOT put business logic in a `Primitive` pillar component; do NOT fetch in a leaf.
- DO NOT skip self-review steps for "obviously small" features — banking-grade applies on every Generate.

## References

| Need | File |
|------|------|
| Banking-grade frontend non-negotiables + 9 v2 augmentations (analytics events, optimistic rollback, error classes, token discipline, bundle guard, etc.) | `references/implementation-rules.md` |
| Component pillars (Page / Feature / Primitive / Hook / Type / Util), TypeScript strict patterns, test conventions | `references/react-typescript-conventions.md` |
| YES/NO a11y scan applied at step 6 + step 8 | `references/a11y-checklist.md` |
| YES/NO security scan applied at step 6 + step 8 | `references/security-checklist.md` |
| State ownership decision tree (server / client / URL / form / local / derived) | `references/state-management-rules.md` |
| YES/NO self-review applied at step 8 before emit | `references/self-review-checklist.md` |
| Extraction lineage, augmentations, drops, deviations | `RATIONALE.md` (human audit only — do not load into LLM context) |
