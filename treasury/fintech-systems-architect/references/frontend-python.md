# Frontend & Python Code Style

## Contents
- TypeScript / React frontend — stack decisions
- Frontend structure
- Frontend principles
- Python structure
- Python principles

## TypeScript / React frontend — stack decisions

- **Framework**: Next.js (SSR/SSG needed, SEO) or Vite + React (pure SPA, internal tools).
- **Data fetching**: TanStack Query for server state (caching, retry, stale-while-revalidate).
- **Client state**: Zustand (simple, no boilerplate). Redux only when state is truly complex.
- **Types from API**: OpenAPI codegen (`openapi-typescript` or `orval`). NEVER hand-write API types.
- **Styling**: Tailwind CSS by default; CSS Modules if the team prefers scoped classes.
- **Forms**: React Hook Form + Zod for the validation schema.
- **Testing**: Vitest (unit), Playwright (E2E), MSW for API mocking.

## Frontend structure

Feature-folder components. Colocate component + hook + types + test per feature.

```
/features
  /loan
    LoanForm.tsx
    useLoan.ts
    loan.types.ts
    loan.api.ts
    LoanForm.test.tsx
```

## Frontend principles

- Strict TypeScript (`strict: true`). No `any` unless truly unavoidable.
- Prefer `type` over `interface` for data shapes; use `interface` for contracts.
- Explicit return types on public functions.
- Use `Result<T, E>` pattern or discriminated unions over thrown exceptions for domain errors.
- Dependency injection via constructor (backend) or React context (frontend).

## Python structure

Package-per-domain, same layout philosophy as Go.

```
/domain
  /loan
    handler.py
    service.py
    repository.py
    model.py
```

## Python principles

- Type hints everywhere (`mypy --strict` compatible).
- Pydantic for data validation / serialization at boundaries.
- Prefer dataclasses for internal domain models.
- Explicit dependency injection. No module-level singletons.
- Use `__init__.py` to control the public API surface per package.
- Async: use `asyncio` only when IO-bound concurrency is needed. Do not force it.
