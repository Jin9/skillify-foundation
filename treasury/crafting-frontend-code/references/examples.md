# Output Examples

Concrete output shapes for the most common modes. Adapt depth to the task; small fixes should be shorter than these examples.

## Contents
- `design` example — small feature
- `review` example — PR-scope finding set
- `fix` example — isolated component bug
- `optimize` example — bundle regression
- `analyze` example — feature audit

## `design` example

**User prompt:** "Design the saved-filters dropdown for the loan-search page."

**Response shape:**

```
Mode: design

- [x] L1 User Intent: user saves a search and re-applies it later without re-typing.
      Success state: user opens dropdown → selects saved filter → URL updates → list refreshes in <300ms perceived.
- [x] L2 Component tree:
      LoanSearchPage (Next.js page, SSR)
        └─ SavedFiltersDropdown (client feature)
             ├─ Trigger (button + count badge)
             └─ Menu (Radix Popover)
                 └─ FilterRow × N (delete + apply)
      State ownership:
      - list of filters → TanStack Query (server cache)
      - current selection → URL search params (?filter=<id>)
      - open/close → local component state (useState)
- [x] L3 Strategy:
      - Query key: ["savedFilters", userId], staleTime 30s
      - Apply: setSearchParams({ filter: id }) — back-button restores prior filter
      - Delete: optimistic update, rollback on error, invalidate on settle
      - Empty state: "No saved filters yet — save your current search."
      - A11y: Radix Popover handles focus trap + Esc; arrow-key nav between rows; delete confirmed via destructive-button pattern
- [x] L4 Sketch:
      hooks/useSavedFilters.ts → list query + delete mutation
      components/SavedFiltersDropdown.tsx → composition root
      components/FilterRow.tsx → primitive
      types/savedFilter.ts → from openapi-typescript codegen
      tests:
      - vitest+RTL: keyboard nav, optimistic delete + rollback
      - Playwright: save → dropdown shows it → apply changes URL
- [x] Trade-offs:
      complexity ↔ maintainability: one hook owns query+mutations (single invalidation point) > splitting into two
      coupling ↔ flexibility: tied to URL search-param shape — back-button works, but renaming the key is a breaking change for shared links
```

## `review` example

**User prompt:** "Review this LoanForm component."

**Response shape:**

```
Mode: review

P1 — must fix
- src/features/loan/LoanForm.tsx:42 — useState mirrors server data from useLoan(). Causes double-render on every refetch and breaks optimistic-update guarantees.
  Fix: drop the local state, derive directly from query.data.
- LoanForm.tsx:88 — onClick on a <div> with no role/tabIndex/keyboard handler. A11y blocker.
  Fix: <button type="button"> with explicit aria-label if icon-only.

P2 — should fix
- types/loan.ts:1 — hand-written `LoanResponse` type. Will drift from backend schema.
  Fix: replace with output from openapi-typescript (already wired in package.json).
- LoanForm.tsx:120 — useMemo around a 3-property object literal — no expensive compute, no memoized child consumer. Net negative.
  Fix: drop the useMemo.

P3 — nit
- Inline tailwind colors (`bg-zinc-100`) bypass design tokens.
  Fix: use `bg-surface` / `bg-surface-muted`.
- Missing test for the error-submission path (server returns 422).
```

## `fix` example

**User prompt:** "Why does the dropdown re-mount and lose focus when I type in the search box?"

**Response shape:**

````
L1-L3 skipped: isolated fix.

Root cause: `Dropdown` is defined inside the parent component body. Every keystroke re-renders the parent, which creates a new `Dropdown` function reference, and React unmounts/mounts the subtree.

Symptom you saw: focus lost, internal state lost, animations restart.

Fix: hoist `Dropdown` to module scope.

```tsx
// before
function SearchPage() {
  function Dropdown(props) { /* ... */ }   // new identity each render
  return <SearchInput><Dropdown/></SearchInput>
}

// after
function Dropdown(props) { /* ... */ }     // stable identity
function SearchPage() {
  return <SearchInput><Dropdown/></SearchInput>
}
```

Behavior-preserving. No visual regression risk. Verify with Profiler: `Dropdown` should now show "Did not render" between keystrokes when its props are stable.
````

## `optimize` example

**User prompt:** "LCP regressed from 1.9s to 3.4s after the marketing page redesign. Why?"

**Response shape:**

```
Mode: optimize

- [x] Baseline: LCP p75 was 1.9s pre-redesign. Now 3.4s. Goal: ≤2.5s ("good").
- [x] Bottleneck identified (Lighthouse trace + WebPageTest):
      - LCP element changed from a small <h1> to a 1920x1080 hero <img>.
      - Image is 480KB JPEG (no AVIF/WebP), no width/height, lazy-loaded.
      - Loading via <img src="..." loading="lazy"> — browser deprioritizes the LCP element.
      Re-render is NOT the issue. INP is unchanged.
- [x] Options:
      A. Use next/image with priority, AVIF/WebP, explicit dimensions.
         + Likely brings LCP back to ~2.0s, fixes CLS.
         + Standard, no infra change.
         − One PR, designer needs to confirm crops at responsive sizes.
      B. Inline a low-res placeholder + load full image after FCP.
         + Best perceived perf.
         − More code, manual maintenance per hero.
      C. Move the hero below the fold (smaller LCP element).
         − Defeats the design intent.
- [x] Recommendation: Option A.
      Measurement plan:
      - Build the change behind a feature flag for 10% traffic.
      - Compare LCP p75 + CLS p75 vs control over 48h.
      - Promote if LCP ≤ 2.3s p75 AND CLS ≤ 0.05 p75.
```

## `analyze` example

**User prompt:** "Take a look at our checkout flow code and tell me what's working and what's not."

**Response shape:**

```
Mode: analyze

Strengths
- Clean separation: CheckoutPage (route) → CheckoutFlow (feature) → step components (primitives).
- Server state lives in TanStack Query consistently. No useEffect-fetching seen.
- Forms use RHF + Zod with shared schemas — same schema validates client + server.

Gaps
- No error boundary around the payment step. A throw in PaymentProvider crashes the whole page.
- Client state (selected shipping method) is in useState at the top of CheckoutFlow but read 4 levels deep — passes through 3 components that don't use it.
- No analytics on step transitions. Funnel cannot be debugged.
- A11y: radio buttons for shipping method are <div>s with onClick. Not keyboard-operable.
- Bundle: stripe-js is in the main bundle (~120KB) — only needed on the payment step.

Contradictions
- README claims "feature-flag every checkout change" but PaymentProvider was added without a flag in commit 4b1a9c2.
- TanStack Query staleTime is 0 globally, but the cart data is fetched on every step transition — this is the staleTime bug, not a design choice.

Recommendations (prioritized)
1. P1 — Wrap payment step in ErrorBoundary with a fallback that preserves cart state.
2. P1 — Fix shipping-method radio a11y (<input type="radio"> in a <fieldset>).
3. P2 — Lazy-load stripe-js on entering the payment step (dynamic import).
4. P2 — Move shippingMethod selection into Zustand or URL state to remove prop-drilling.
5. P2 — Add `checkout.step.entered` analytics event with stepId + cart total.
6. P3 — Set staleTime per-query for cart (60s) instead of global.
```
