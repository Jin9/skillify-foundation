---
name: implementing-go-template-requirements
description: >-
  Apply a single requirement (spec line, ticket, bug report, user story) to a
  Go service that follows the go-template scaffold by editing ONLY business
  logic under `app/[domain]/` plus narrow `register*` wiring in `router/`.
  Also accepts requirement-driven Fowler refactors (God Handler, Long
  Parameter List, Too Many Returns, Feature Envy) scoped to handlers,
  consumers, and access-layer signatures. Locks `main.go`, `config/`,
  `router/deps.go`, `Makefile`, `Dockerfile`, CI, and tooling from
  modification. Use when the user says "implement this requirement in the
  go-template service", "add this endpoint per spec without touching infra",
  "fix this bug, business logic only", "wire up this Kafka event from the
  requirement", "implement spec.md item N", or "refactor this handler to
  address [named smell]". Do NOT use to create a new service from scratch,
  change middleware, modify `main.go`/`config/`/`router/deps.go`/`Dockerfile`/
  `Makefile`/CI, do frontend or non-Go work, or perform refactors with no
  named smell.
argument-hint: Describe the requirement (spec section, ticket title, bug summary). The skill picks the right files, keeps scaffold untouched, and runs the verification gate.
---

# Implementing Go-Template Requirements

Use this skill to translate **one requirement** into code inside a Go service that follows the `go-template` scaffold. The job is **business logic, not scaffolding**. The scaffold (bootstrap, config, middleware, CI, infrastructure wiring) is **locked**: do not edit it unless the requirement explicitly forces a scaffold delta, and only after STOP-and-ask confirmation.

---

## 1. When to use this skill

Trigger on a pasted requirement/spec/ticket/bug asking for an implementation in this repo: "implement this requirement in the go-template service", "add this endpoint per spec without touching infra", "fix this bug, business logic only", "wire up this Kafka event", "implement spec.md item N", "extend `<domain>` to satisfy `<requirement>`", or "refactor `<handler|consumer|service>` to address a named smell from `references/fowler-patterns.md` §4 (God Handler, Long Parameter List, Too Many Returns, Feature Envy)".

Do NOT use for: creating a new service from scratch; changing middleware, security headers, CORS, timeouts, or Gin engine setup; modifying `main.go`, `config/`, `router/deps.go`, `Makefile`, `Dockerfile`, CI, `.scripts/`, `.golangci.yaml`, `.mockery.yaml`, `.env.template`; frontend/Terraform/Kubernetes/non-Go work; or "clean this up" requests with no named §4 smell (ask which smell applies; if none, decline).

---

## 2. The scaffold-lock policy (the defining rule)

Three zones. Every edit must be classified before it is made.

| Zone | Files | What you may do |
|---|---|---|
| ✅ ALLOWED | `app/<domain>/**` (every business-logic file, including `access/` and generated `mocks/mocks.go`) | Edit freely within the conventions in §4 and `references/naming-conventions.md`. |
| ⚠️ NARROW | `router/router.go` (only the `register<Domain>Routes(r, d)` block for your domain), `router/subscriber.go` (only the entry inside `registerEventRoutes`), `spec.md` (only to document new endpoint/event) | Add or change ONLY the named block. Do not touch middleware chain, `New()` skeleton, Kafka consumer-group lifecycle, security headers, CORS, timeout, or logger. |
| 🛑 FORBIDDEN | `main.go`, `config/`, `router/deps.go`, `Makefile`, `Dockerfile`, `docker-compose.yml`, `gitlabci.yml`, `.github/`, `.scripts/`, `.golangci.yaml`, `.mockery.yaml`, `.env.template`, `VERSION`, `CHANGELOG.md`, `go.mod`, `go.sum` | STOP. Do not edit. Ask the user. |

> The ALLOWED file prefixes realise Fowler's Repository (`storage_*`), Cache (`cache_*`), Gateway (`client_*`), Service Layer (`service_*`), and Domain Model patterns. `router/deps.go` is the Composition Root (FORBIDDEN). See `references/fowler-patterns.md` for the cross-walk and the small-safe-steps refactoring discipline that applies when extending handler/consumer code.

**STOP triggers** — when the requirement crosses the scaffold boundary (new infra/SDK client, new env var, logging/signal/shutdown change, anything touching `go.mod`, `Dockerfile` base image, middleware, or CI), surface the conflict, name the forbidden files, and ask the user to confirm the scaffold delta before continuing. Do not silently edit forbidden files. The full STOP-trigger list, per-file rationale, and the two-path STOP handling are authoritative in `references/scaffold-lock-policy.md`.

---

## 3. Workflow

Run these in order. Every step has an entry and exit condition.

**Step 1 — Parse the requirement.**
- Entry: user message containing a requirement, spec excerpt, ticket, or bug description.
- Extract: kind (new endpoint / new consumer / new storage method / bugfix), target domain, name of action or event, inputs, outputs, side effects (DB writes, Kafka publish), dependencies (storage, cache, external client).
- Exit: a one-paragraph summary of *what* the change is and *which files* it implies. Restate to the user before editing if any extraction is ambiguous.

**Step 2 — Locate the domain.**
- Entry: target domain name from Step 1.
- List `app/<domain>/` to confirm it exists. Read `app/<domain>/handler.go` to see existing `HandlerConfig` dependencies.
- If `app/<domain>/` does not exist: this is a new-domain creation, still allowed (the whole `app/` tree is in the allowlist). Announce the new-domain decision and continue.
- Exit: confirmed target domain path.

**Step 3 — Audit existing access interfaces.**
- Entry: target domain confirmed.
- Read every `access/storage_*.go`, `access/cache_*.go`, `access/client_*.go` in the domain. Match the requirement's I/O needs against existing methods.
- Reuse existing methods where possible. Extend an interface only when the requirement genuinely needs a new query.
- Exit: a list of access-layer additions (new methods, new interfaces) or "no access changes needed".

**Step 4 — Apply the scaffold-lock test.**
- Entry: a draft list of files to create or modify.
- For each file, classify into ALLOWED / NARROW / FORBIDDEN per §2.
- If any file is FORBIDDEN: stop, surface the conflict, ask the user.
- If any file is NARROW: confirm you will edit only the named block.
- Exit: a clean edit plan, all files classified.

**Step 5 — Implement file deltas in dependency order.**
- Entry: clean edit plan from Step 4. Match the requirement to one recipe (add-endpoint / add-consumer / add-storage-method / fix-bug / new-domain / refactor) and follow its exact ordered file set, templates, and narrow-wiring step in `references/implementation-recipe.md`.
- Exit: every file in the recipe's set is written and follows naming + style conventions.

**Step 6 — Regenerate mocks.**
- Entry: access-layer interfaces changed. Run `mockery` (existing `.mockery.yaml`, do not edit it); never hand-edit `app/<domain>/access/mocks/mocks.go`.
- Exit: mocks are current.

**Step 7 — Write tests.**
- Entry: handler/consumer/service code complete.
- One test file per code file (`handler_<action>_test.go`, `consumer_<action>_test.go`). Use the `mockArgs`/`args`/`want`/`prepare` table-driven shape and the per-branch coverage rules (every `if err != nil`, every required-field case, model-getter parse failures, consumer invalid-JSON + validation cases) authoritative in `references/testing-pattern.md`. Templates: `templates/handler_test.go.tmpl`, `templates/consumer_test.go.tmpl`.
- Exit: `go test -race -cover ./app/<domain>/...` reports 100% statement coverage in `app/<domain>/` (excluding generated `mocks/`).

**Step 8 — Run the verification gate.**
- Entry: code + tests complete.
- Run `make precommit` (or equivalent: `go fmt ./...`, `go vet ./...`, `golangci-lint run --fix`, `go test -race -cover ./...`).
- All must pass. Coverage must remain at the existing threshold.
- Exit: green build.

**Step 9 — Diff allowlist check.**
- Entry: green build.
- Run `git diff --name-only` against the merge base of your working branch.
- Confirm every file is in ALLOWED (or is a NARROW edit limited to its named block).
- If any FORBIDDEN file appears: revert it. Surface to the user that the gate caught a scaffold leak.
- Exit: a clean diff that touches only business logic + narrow wiring + spec doc.

---

## 4. Naming and style — load-bearing conventions

The prefix-name rule: the **type prefix comes first** in the filename — `handler_<action>.go`, `consumer_<action>.go`, `service_<action>.go`, `storage_<dep>.go`, `cache_<dep>.go`, `client_<dep>.go`. Never `<action>_handler.go` or `<dep>_storage.go`. `access/` files **co-locate** interface + unexported impl + domain model + sentinel errors + constants in a single file — no separate `model.go`/`errors.go`/`constants.go`.

All other conventions (prefix↔Fowler-pattern cross-walk, package/import rules, constructor-returns-interface, context-first, error wrapping, the `(ctx, ≤3 params) (T, error)` signature shape and its Long-Parameter-List / Too-Many-Returns triggers, `wrapper`/`kafka.BindMessage`/event-name/`slog` rules) are authoritative in `references/naming-conventions.md`. Refactor mechanics: `references/fowler-patterns.md` §4–§5.

---

## 5. Pre-commit checklist

The diff-allowlist gate is the one box that must be ticked **in the body** before claiming done:

- [ ] **Diff-allowlist gate (Step 9):** `git diff --name-only` against the merge base shows only ALLOWED `app/<domain>/**` paths and NARROW edits confined to their named block (`router/router.go` `register<Domain>Routes`, `router/subscriber.go` `registerEventRoutes`, `spec.md` new entry). Any FORBIDDEN file means the gate failed — revert it.

The full 18-item verification checklist (naming, code conventions, tests, coverage, build gates, docs, commit hygiene) with rationale and per-failure remediation is authoritative in `references/verification-checklist.md`. Tick every box there before declaring the requirement done.

---

## 6. Constraints (hard rules)

- **DO NOT** edit any file in the FORBIDDEN zone (§2) without explicit user confirmation. The scaffold lock is the core safety property of this skill.
- **DO NOT** skip the diff allowlist check (Step 9) — it is the final gate that catches scaffold leakage.
- **DO NOT** mix refactor commits with feature commits. Refactor first, commit, run `make precommit`, then add the feature, commit. One refactoring action per commit (Fowler small-safe-steps).
- **DO NOT** accept refactor requests without a named smell from `references/fowler-patterns.md` §4. "Clean this up" is not a requirement — ask which smell applies; if none, decline.

The remaining hard rules are authoritative in their references: no cross-domain imports and no `model.go`/`errors.go`/`constants.go` inside `access/` (`references/naming-conventions.md`); always wrap raw SDK handles in `access/` and never `json.Unmarshal` in consumers (`references/scaffold-lock-policy.md`, `references/common-module-quickref.md`); no preemptive abstractions and never hand-edit generated mocks (`references/fowler-patterns.md`, `references/implementation-recipe.md`).

---

## 7. References

| Need | File |
|---|---|
| Full scaffold-lock allowlist + per-file rationale | `references/scaffold-lock-policy.md` |
| File naming conventions and package rules | `references/naming-conventions.md` |
| Per-task recipes: add-endpoint, add-consumer, add-storage-method, fix-bug | `references/implementation-recipe.md` |
| Table-driven test pattern with `mockArgs`/`args`/`want`/`prepare` | `references/testing-pattern.md` |
| Quick reference for common-module helpers (`wrapper`, `serror`, `kafka`, `app`) | `references/common-module-quickref.md` |
| Pre-commit checklist and diff allowlist check | `references/verification-checklist.md` |
| Fowler pattern cross-walk + code smells + small-safe-steps refactoring discipline (handlers/consumers scope) | `references/fowler-patterns.md` |
| Worked end-to-end example: input shapes, parsed requirement, file plan | `examples/usage.md` |

## 8. Templates

| Need | File |
|---|---|
| `handler_<action>.go` skeleton | `templates/handler.go.tmpl` |
| `handler_<action>_test.go` skeleton | `templates/handler_test.go.tmpl` |
| `consumer_<action>.go` skeleton | `templates/consumer.go.tmpl` |
| `consumer_<action>_test.go` skeleton | `templates/consumer_test.go.tmpl` |
| `access/storage_<dep>.go` skeleton with interface + impl + model | `templates/storage.go.tmpl` |

---

## 9. Troubleshooting

Signal → action table (STOP triggers, coverage gaps, consumer/payload issues, mockery surprises, and each named-smell refactor path) is in `references/implementation-recipe.md` under "Troubleshooting".
