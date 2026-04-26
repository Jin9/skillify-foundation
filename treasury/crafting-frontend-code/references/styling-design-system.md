# Styling & Design System

## Contents
- Tailwind defaults
- Design tokens
- Theming (light / dark / brand)
- Responsive strategy
- Component variants (cva / tailwind-variants)
- Design system layering
- Risky patterns

## Low-risk use

Use the repo's existing styling system first. Treat Tailwind, tokens, and variant libraries as preferred patterns only when already present or when the user asks for a design-system proposal.

## Tailwind Defaults

- If the repo uses Tailwind, use it for layout, spacing, typography, and color.
- Avoid arbitrary values (`w-[437px]`) when a token exists. They signal a missing token.
- Compose long class lists via `clsx` / `cn` helper. Pull recurring sets into a component, not into a string constant.
- Use the repo's existing class-ordering formatter. If Tailwind formatting is missing, propose it separately.

## Design tokens

Tokens live in `tailwind.config.ts` (or a token JSON consumed by it). Categories:

| Category | Examples |
|---|---|
| Color | `bg-surface`, `text-muted`, `border-subtle` |
| Spacing | `p-2`, `gap-4` (4px scale) |
| Radius | `rounded-md`, `rounded-pill` |
| Shadow | `shadow-sm`, `shadow-overlay` |
| Type | `text-body`, `text-heading-lg` |
| Motion | `duration-fast`, `ease-emphasized` |

Rule: components should consume **semantic** tokens (`bg-surface`) when the design system defines them. Re-skinning becomes a token swap, not a code change.

## Theming

- Use CSS variables for token values; Tailwind reads via `var(--color-surface)`.
- Theme toggle = swap a class on `<html>` (`light` / `dark` / `brand-x`).
- Server-render the initial theme to avoid FOUC (read cookie / header).

## Responsive strategy

- Mobile-first. Add `sm:` / `md:` / `lg:` breakpoints up.
- Container queries (`@container`) for component-driven responsiveness when layout varies by parent, not viewport.
- Avoid fixed pixel widths in layout; use flex / grid + min/max constraints.

## Component variants — `cva` or `tailwind-variants`

```ts
const button = cva('inline-flex items-center font-medium rounded-md', {
  variants: {
    intent: { primary: 'bg-brand text-on-brand', secondary: 'bg-surface text-default' },
    size: { sm: 'h-8 px-3 text-sm', md: 'h-10 px-4' },
  },
  defaultVariants: { intent: 'primary', size: 'md' },
})
```

Type-safe, exhaustive, easy to test.

## Design system layering

```
Tokens (JSON / CSS vars)
    ↓
Primitives (Button, Input, Dialog) — headless behavior + tokens
    ↓
Patterns (FormField, EmptyState, DataTable) — composed primitives
    ↓
Features — app-specific composition
```

- Primitives never import from features.
- Patterns never import from features.
- Features import freely from patterns / primitives / tokens.

Enforce with existing boundary tooling when available; propose `eslint-plugin-boundaries` or `dependency-cruiser` only as a separate tooling change.

## Risky Patterns

- Inline styles for anything reachable by a token.
- Hard-coded hex / rgb in components.
- Global CSS overrides targeting design-system primitives.
- Re-implementing Dialog / Menu / Popover from scratch — use Radix / React Aria.
- "Just one more `!important`" — find the specificity bug instead.
- A component that ships its own color palette parallel to tokens.
