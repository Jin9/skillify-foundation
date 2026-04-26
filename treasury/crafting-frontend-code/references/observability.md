# Observability

## Contents
- What to instrument
- Error tracking
- RUM (Real User Monitoring)
- Logging
- Feature flags
- Analytics
- Privacy
- Risky patterns

## Low-risk use

Instrumentation changes can alter privacy, cost, and production behavior. Do not add analytics, replay, third-party scripts, or PII-bearing fields without confirming consent, sampling, and scrubbing expectations.

## What to instrument

| Signal | Tool category | Purpose |
|---|---|---|
| Uncaught errors | Sentry / Rollbar | Find production bugs |
| Web Vitals | `web-vitals` lib + RUM | Real-user perf |
| Console errors / warnings | Sentry breadcrumbs | Pre-error context |
| Route changes | RUM / analytics | User flow |
| Feature-flag exposure | Flag SDK (LaunchDarkly / GrowthBook) | A/B + rollout safety |

## Error tracking

- Sentry (or equivalent) initialized once at app entry.
- Capture: source maps uploaded on build; release tagged with git SHA; user ID set after auth (anonymized if needed).
- Wrap React tree with `<ErrorBoundary>` per route + a global one. Show fallback UI; report the error.
- Add breadcrumbs for key actions (route change, mutation, auth event).
- Filter noise: known third-party errors, browser extension noise.

```ts
Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  release: process.env.NEXT_PUBLIC_GIT_SHA,
  tracesSampleRate: 0.1,
  replaysSessionSampleRate: 0.0,
  replaysOnErrorSampleRate: 1.0,
  beforeSend(event) { return scrubPII(event) },
})
```

## RUM (Real User Monitoring)

- Send Core Web Vitals (LCP, INP, CLS, TTFB, FCP) per page load.
- Tag samples with route, device class, connection type.
- Compare p50 vs p75 vs p95 — averages lie.
- Alert on regressions vs the previous release window.

## Logging

- Production client warnings/errors should go to an approved collector when one exists, not only the browser console.
- Levels: `debug` (dev only), `info`, `warn`, `error`. Ship `warn`+ in prod.
- Structured (JSON) — log objects, not concatenated strings.
- Always scrub PII / tokens before send.

## Feature flags

- Wrap risky launches in a flag; default off in prod.
- Track exposure in analytics so you can compute lift.
- Server-evaluate flags when possible to avoid client flash + tampering.
- Schedule flag cleanup — track removal in the same ticket as launch.

## Analytics

- Define events with a schema (name, required props). Treat events as a contract.
- Centralize through one wrapper — easy to swap providers / add scrubbing.
- Use semantic event names (`checkout.started`), not "Button Click 17."
- Rate-limit / debounce noisy events (scroll, mouse-move).

## Privacy

- Cookie consent banner before any non-essential tracking fires.
- Honor DNT and regional opt-outs (GDPR, CCPA).
- Anonymize IPs at collector when possible.
- Document what you collect in a privacy notice. Match reality.

## Risky Patterns

- Shipping `console.log` debugging to prod.
- Sending unscrubbed request bodies to Sentry.
- Tracking everything "in case we need it" — privacy + cost time bomb.
- Flags that have been "on for everyone" for >3 months — clean them up.
- Using analytics events as a substitute for a backend audit log.
- Sampling errors aggressively then wondering why bug reports don't match dashboards.
