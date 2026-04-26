# Testing

## Contents
- Test pyramid
- Vitest + RTL conventions
- MSW for network
- Playwright (E2E)
- Visual regression
- Accessibility tests
- Performance regressions
- Risky patterns

## Low-risk use

Prefer the repo's existing test tools and scripts. Add only the smallest test that proves the change unless the task is explicitly to expand testing strategy.

## Test pyramid

| Layer | Tool | What it tests |
|---|---|---|
| Unit | Vitest | Pure functions, hooks, reducers |
| Component | Vitest + RTL | One component, mocked deps |
| Integration | Vitest + RTL + MSW | Feature, real network mocks |
| E2E | Playwright | Real browser, real backend or seeded |
| Visual | Playwright + Percy / Chromatic / Argos | Pixel diffs |

Shape: lots of unit + component, fewer integration, few E2E. Slow E2E for critical paths only.

## Vitest + RTL conventions

- Test behavior, not implementation. Query by role / label, not by class / test-id (use `data-testid` only when no other option exists).
- One assertion per behavior — multiple `expect`s per `it` is fine if they test one behavior together.
- Use `userEvent` (not `fireEvent`) for interactions — it simulates real user flow.
- Don't snapshot whole components — snapshot small derived structures (e.g., serialized state).

```ts
test('submits valid form', async () => {
  const user = userEvent.setup()
  render(<SignupForm onSubmit={onSubmit}/>)
  await user.type(screen.getByLabelText(/email/i), 'a@b.com')
  await user.click(screen.getByRole('button', { name: /sign up/i }))
  expect(onSubmit).toHaveBeenCalledWith({ email: 'a@b.com' })
})
```

## MSW for network

- Define handlers per feature in `mocks/handlers/`.
- Use the same handlers in tests AND in dev (`worker.start()`).
- Override per-test for error / loading variants — don't mock `fetch` directly.

## Playwright (E2E)

- One project per browser you actually ship to (Chromium minimum).
- Page Object Model for stable locators (`page.getByRole`, `page.getByLabel`).
- Network: prefer real backend in test env over mocking. Mock only third parties you don't own.
- For mature apps, run E2E in CI on PR + main. Smoke subset on every PR; full suite on main / nightly.

## Visual regression

- Snapshot only stable-looking states (loading, empty, success, error).
- Mask dynamic regions (timestamps, animations).
- Review diffs as part of PR. Auto-approve only on intentional design system bumps.

## Accessibility tests

- `@axe-core/react` in dev (logs violations to console).
- `@axe-core/playwright` in E2E for key routes; gate CI once the team agrees on thresholds.

## Performance regressions

- Lighthouse CI on key pages; fail only on agreed budget breaches.
- Bundle-size guard: `size-limit` or `@next/bundle-analyzer` diff on PR.

## Risky Patterns

- Testing implementation (`expect(component.state).toBe(...)`).
- `data-testid` everywhere — masks bad DOM semantics. Use roles / labels first.
- Mocking `useState` or other React internals.
- Sleeping / timing-based assertions — use `findBy*` / `waitFor`.
- E2E tests that depend on each other's state.
- Skipped tests left in `main` without a tracked issue.
- Coverage thresholds without conversation about what they actually catch.
