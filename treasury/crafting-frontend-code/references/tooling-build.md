# Tooling & Build

## Contents
- Framework decision: Next.js vs Vite + React
- Bundler & transpiler defaults
- Monorepo
- Linting / formatting
- CI checks
- Dependency hygiene
- Risky patterns

## Low-risk use

Use the existing framework, package manager, test runner, and CI shape unless the user asked for tooling changes. Treat new tools, package-manager switches, and CI gates as migration work requiring an explicit plan.

## Framework decision

| Need | Pick |
|---|---|
| SEO, marketing pages, mixed SSR / SSG | **Next.js (App Router)** |
| Server Components, edge, streaming | **Next.js** |
| Pure SPA, internal tool, dashboard behind auth | **Vite + React Router** |
| Embedded widget / library | **Vite (lib mode)** |
| Existing CRA repo | Migrate to **Vite** (CRA is unmaintained) |

Greenfield heuristic: Next.js for product-facing, Vite for internal. Existing apps should not switch frameworks without a migration plan.

## Bundler & transpiler defaults

- Vite (Rollup) or Next.js (Turbopack / Webpack) — don't roll your own.
- ESBuild for transforms; SWC where supported.
- TS via the bundler (no separate `tsc` build step except for type-check / library publishing).
- Source maps in production (uploaded to Sentry, hidden from public).

## Monorepo

When to adopt: 2+ deployable apps sharing components or types.

- **Tool**: Turborepo or Nx when the repo already uses or needs monorepo orchestration.
- **Package manager**: keep the repo's current package manager unless switching is the task.
- **Layout**:
  ```
  apps/
    web/
    admin/
  packages/
    ui/        # design-system primitives
    config/    # eslint, tsconfig, tailwind presets
    api/       # generated client + types
  ```
- **Boundaries**: enforce with `eslint-plugin-boundaries` — `apps` may import `packages`, never reverse.

## Linting / formatting

- ESLint flat config. Plugins: `@typescript-eslint`, `react`, `react-hooks`, `jsx-a11y`, `tailwindcss`, `import`.
- Prettier + `prettier-plugin-tailwindcss`.
- Run on pre-commit via `lint-staged` + `husky`.
- CI runs the full lint + type-check; pre-commit only the staged subset.

## CI Checks

```
typecheck → lint → unit/integration tests → build → e2e (subset) → bundle-size guard → lighthouse-ci
```

Preferred mature pipeline. If checks are missing, propose them separately instead of adding gates inside an unrelated feature fix.

## Dependency hygiene

- Use the repo package manager with a frozen/immutable lockfile in CI.
- `renovate` or `dependabot` for weekly bumps.
- Audit on a schedule (`pnpm audit --prod`); investigate any HIGH / CRITICAL.
- Resist trendy deps — measure value before adopting.

## Risky Patterns

- Hand-managed `package.json` versions across many packages when a workspace tool is available.
- Disabling lint rules inline without a comment + tracked issue.
- Committing `node_modules` or `.next` / `dist`.
- Two bundlers in one app.
- "Works locally" — if CI doesn't run it, it doesn't ship.
- Build scripts that swallow errors silently.
