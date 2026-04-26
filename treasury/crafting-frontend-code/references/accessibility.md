# Accessibility (a11y)

## Contents
- Baseline contract
- Semantic HTML before ARIA
- Keyboard
- Forms
- Focus management
- Color & contrast
- Motion & preferences
- Testing
- Risky patterns

## Low-risk use

Accessibility fixes should be behavior-preserving unless the current UI is inaccessible. Prefer semantic HTML and small targeted changes before introducing new component libraries.

## Baseline contract

Target WCAG 2.2 Level AA. Shipped UI should:
- Be operable with keyboard alone.
- Have a visible focus state.
- Pass `axe` automated checks (no critical violations).
- Work with a screen reader (NVDA / VoiceOver smoke test for critical flows).

## Semantic HTML before ARIA

- Use `<button>` for actions, `<a href>` for navigation. Do not put `onClick` on `<div>`.
- Use `<label htmlFor>` (or wrap inputs) — never use placeholder as label.
- Use `<nav>`, `<main>`, `<header>`, `<footer>`, `<section>` for landmarks.
- Headings in order (h1 → h2 → h3). One `<h1>` per page.
- Lists for list-shaped content (`<ul>`, `<ol>`).

ARIA is a fallback when semantics don't exist. Wrong ARIA is worse than no ARIA.

## Keyboard

- Tab order matches visual order.
- Custom widgets follow [WAI-ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/) keyboard contracts (menu, listbox, dialog, tabs).
- `Esc` closes overlays. `Enter` / `Space` activates buttons.
- No keyboard traps — except modals, which should trap correctly while open.

## Forms

- Every input has an associated label.
- Required fields marked in label text, not just color or `*`.
- Validation errors are programmatically associated (`aria-describedby`) and announced.
- Group related fields with `<fieldset>` + `<legend>`.

## Focus management

- After route change, move focus to the page heading or main content.
- Modal opens → focus first interactive element. Modal closes → focus returns to trigger.
- Use `inert` (or `aria-hidden`) on background when a modal is open.

## Color & contrast

- Text contrast: 4.5:1 (normal), 3:1 (large or UI).
- Never convey state by color alone — pair with icon, text, or pattern.
- Test with simulated color-blindness (Chrome DevTools).

## Motion & preferences

- Respect `prefers-reduced-motion`: disable / shorten transforms and parallax.
- Respect `prefers-color-scheme` for theming.

## Testing

- `@axe-core/react` in dev to catch issues at render time.
- `@axe-core/playwright` for CI or release checks on key pages.
- Manual: keyboard-only walkthrough + 1 screen-reader smoke test per critical flow per release.

## Risky Patterns

- `<div onClick>` for clickable things.
- Removing focus outline without replacing it.
- Placeholder as the only label.
- `aria-label` on a `<div>` that pretends to be a button.
- Auto-focusing a deep field on page load (steals focus from screen readers).
- `tabindex` > 0 (hijacks tab order).
- Toast-only error messaging (screen readers may miss it — pair with inline `role="alert"`).
