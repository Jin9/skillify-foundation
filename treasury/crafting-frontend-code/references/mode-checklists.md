# Mode Checklists

Per-mode checklists for `crafting-frontend-code`. Load and fill the checklist for the mode classified in the Safety workflow (Step 1). Include the filled checklist in the output when it improves reviewability; for direct implementation or small fixes, apply it internally and summarize only the important result.

## `design` — full UI architecture pass

```
- [ ] L1 User Intent + success state stated
- [ ] L2 component tree + state ownership + rendering model named
- [ ] L3 fetching, caching, error/loading, a11y contract chosen
- [ ] L4 components / hooks / types / tests sketched
- [ ] Trade-offs stated on all 4 axes
```

## `optimize` — performance / bundle / DX trade-off

```
- [ ] Current baseline measured (Lighthouse / Web Vitals / bundle stats)
- [ ] Bottleneck identified, not guessed (TTFB / LCP / INP / CLS / TBT / bundle / re-render)
- [ ] Options enumerated with pros/cons (code-split, memo, defer, server-render, virtualize)
- [ ] Recommendation + measurement plan
```

## `fix` — minimal, direct solution

```
- [ ] State "L1-L3 skipped: isolated fix" when useful
- [ ] Root cause vs symptom called out (re-render? stale closure? key? hydration mismatch? effect dependency?)
- [ ] Behavior-preserving unless flagged; visual regression risk noted
- [ ] Existing tests or scripts checked before adding new tools
```

## `analyze` — deep breakdown

```
- [ ] Strengths
- [ ] Gaps (a11y / perf / type-safety / state)
- [ ] Contradictions (e.g., client state mirroring server state)
- [ ] Recommendations, prioritized
```

## `review` — code / component / architecture review

```
- [ ] Risks flagged by severity (P1/P2/P3)
- [ ] A11y, performance, type-safety, security (XSS / CSRF / CSP) checked
- [ ] Concrete fix suggested for each finding
- [ ] Findings are evidence-backed with file/line references when local code is available
```

## `plan` — produce a plan, do not execute

```
- [ ] Priorities set
- [ ] Open decisions listed with owners (design, product, infra)
- [ ] Dependencies + sequencing stated (API ready? feature flag? token bump?)
```
