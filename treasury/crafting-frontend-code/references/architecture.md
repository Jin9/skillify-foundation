# Frontend Architecture

## Contents
- Reference component model (layering)
- Layer responsibilities
- State ownership map
- Data flow rules
- Composition patterns
- Rendering model decision (CSR / SSR / SSG / RSC)
- Risky patterns

## Low-risk use

Treat this as a decision aid, not a mandate. Preserve existing app layering unless the user asked for refactoring or the current boundary causes a clear bug, performance issue, accessibility issue, or maintenance risk.

## Reference component model

```
Route / Page
    ↓ (fetches, error/loading boundary, layout)
Feature Component
    ↓ (interaction logic, local state)
Primitive / Design-System Component
    ↓ (presentation only)
```

## Layer responsibilities

- **Page / Route**: route-level data fetching, error + loading boundaries, layout shell, SEO meta.
- **Feature**: interaction state, form orchestration, business UI logic. Talks to hooks.
- **Primitive (Design System)**: presentation + a11y. No fetching. No business rules. Pure props in, JSX out.
- **Hooks**: preferred place for side-effects (network, storage, subscriptions, timers), unless the framework requires a route/server boundary.

## State ownership map

| State type | Lives where | Tool |
|---|---|---|
| Server data | Cache | TanStack Query |
| URL / route | URL | Next router / search params |
| Cross-page client | Store | Zustand (or Redux for complex) |
| Local component | `useState` | React |
| Form draft | Form lib | React Hook Form |
| Derived | Compute on render | `useMemo` only when proven needed |

Rule: avoid duplicating server data into client state. The cache should remain the source of truth unless there is an explicit offline/draft requirement.

## Data flow rules

- Top-down props for pure presentation.
- Hooks for side-effects and data dependencies.
- Events bubble up via callbacks. Avoid prop-drilling > 2 levels — use context, composition, or a store.
- Context for transient cross-cutting (theme, auth identity). Avoid it for high-frequency mutable app state.

## Composition patterns

- Prefer composition (children, slots) over configuration props.
- Compound components (`<Tabs><Tabs.List/><Tabs.Panel/></Tabs>`) for related parts.
- Headless primitives (Radix, React Aria) + your styling — do not re-invent menu / dialog / popover a11y.
- Render props or hooks over HOCs.

## Rendering model decision

| Need | Pick |
|---|---|
| SEO matters, content varies per request | SSR |
| SEO matters, content stable across requests | SSG / ISR |
| Auth-gated app, no SEO | CSR (or RSC with client islands) |
| Server-only data, no interactivity | RSC (Server Components) |
| Highly interactive island within an SSR/SSG page | Client Component |

Rule: pick the cheapest model that meets the SEO, freshness, and interactivity requirement. Do not introduce SSR/RSC to an existing CSR surface without a migration reason.

## Risky Patterns

- Business logic in primitive / design-system components.
- Fetching inside leaf components.
- Mirroring server state into a client store.
- Context for high-frequency updates (causes whole-tree re-renders).
- "God hook" that does fetch + form + transform + side-effects.
- Prop-drilling through multiple unused layers when composition, context, or a store would reduce coupling.
- Page components with large, repeated JSX blocks that hide feature boundaries.
