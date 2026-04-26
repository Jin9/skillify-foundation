# Frontend Performance

## Contents
- Performance budget (Core Web Vitals)
- Measurement before optimization
- Render performance
- Bundle size
- Network & data
- Asset / image strategy
- Risky patterns

## Low-risk use

Measure before changing behavior. Do not add instrumentation, CI gates, third-party scripts, or large dependency swaps inside a performance fix unless the user approves the broader change.

## Performance budget (Core Web Vitals targets)

| Metric | Good | Needs improvement |
|---|---|---|
| LCP | ≤ 2.5s | ≤ 4.0s |
| INP | ≤ 200ms | ≤ 500ms |
| CLS | ≤ 0.1 | ≤ 0.25 |
| TTFB | ≤ 800ms | ≤ 1.8s |
| FCP | ≤ 1.8s | ≤ 3.0s |

Set a budget per route. If the project already has budgets, enforce them; otherwise propose budgets separately and report measured deltas.

## Measurement before optimization

- Lighthouse + WebPageTest for synthetic.
- RUM (Vercel Analytics, Sentry, or `web-vitals` lib) for field data.
- React DevTools Profiler for render bottlenecks.
- Chrome Coverage tab for unused JS/CSS.
- `source-map-explorer` or `vite-bundle-visualizer` for bundle composition.

Do not optimize without a number or a reproducible user-visible symptom. "Feels slow" alone is not enough for a broad refactor.

## Render performance

- `React.memo` only after profiler shows wasted re-renders.
- `useMemo` / `useCallback` only when (a) referential stability matters for memoized children or (b) computation is provably expensive.
- Stable keys on lists. Avoid array indexes when the list can reorder, insert, or delete.
- Split context providers — narrower scope reduces blast radius.
- `useTransition` / `useDeferredValue` for non-urgent updates.
- Virtualize lists > ~200 items (TanStack Virtual, react-window).

## Bundle size

- Per-route code splitting via dynamic `import()`.
- Tree-shake: import named symbols, verify with bundle visualizer.
- Replace heavyweight deps:
  - `date-fns` over `moment`
  - native `fetch` over `axios` (unless you actually need interceptors)
  - `clsx` over `classnames`
  - per-icon imports over full icon libraries
- Defer non-critical JS (`<script defer>`, `next/script` strategy).
- Polyfill only for browsers you support; check bundle for legacy transforms.

## Network & data

- Server state: TanStack Query with `staleTime` set to avoid refetch storms.
- Prefetch on hover / intent for likely navigations.
- Parallelize independent requests; avoid waterfalls.
- Use HTTP cache headers; let the browser do work.
- `Suspense` boundaries to start streaming UI early (RSC / Next App Router).
- `<link rel="preload">`, `preconnect`, or framework-supported hints for critical assets when measurement shows a loading bottleneck.

## Asset / image strategy

- `next/image` (or equivalent) — automatic responsive + AVIF/WebP + lazy.
- Explicit `width` / `height` on images to prevent CLS.
- Inline critical CSS; defer the rest.
- Self-host fonts with `font-display: swap`; use variable fonts.
- Use `loading="lazy"` for below-the-fold images; eager for the LCP image.

## Risky Patterns

- Optimizing without measuring.
- Wrapping everything in `useMemo` / `React.memo` "just in case" (slows things down).
- Importing entire utility libraries (`import _ from 'lodash'`).
- Loading hero images via JS instead of `<img>` / `next/image`.
- Hydration mismatches "fixed" by suppressing the warning instead of fixing the cause.
- Lazy-loading the LCP image.
