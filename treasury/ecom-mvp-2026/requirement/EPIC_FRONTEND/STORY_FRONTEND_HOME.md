---
template_version: 0.1.0
story_id: STORY_FRONTEND_HOME
epic_id: EPIC_FRONTEND
title: home tab renders featured + categories + new arrivals (Thai labels, ฿ prefix, emerald CTAs, 5-tab bar)
as_a: Guest
i_want: a Thai-language mobile home tab with featured products, categories, and new arrivals
so_that: I can browse the storefront from my phone and tap into a product detail
acceptance_criteria:
  - given: An unauthenticated visitor opens the customer mobile webview
    when: The home tab route renders
    then: |
      The page declares lang='th', the layout fits a 390-wide viewport with a single
      column without horizontal scroll, the bottom tab bar shows the 5 tabs
      (home/search/cart/orders/profile), the primary CTA colour matches the emerald
      brand token, prices render with the ฿ Thai-Baht prefix, and product cards open
      the detail route on tap (Frontend Spec §1, §2, §3).
  - given: The home tab is open and the catalog has at least one product per featured / new-arrivals section
    when: The home tab loads its data from the catalog service
    then: |
      The featured section, the categories section, and the new-arrivals section
      each render at least one item; tapping a product card navigates to the detail
      route for that product (Frontend Spec §3 home row + CAT-001).
priority: Must
size: M
edge_cases:
  - case: Catalog has zero products
    expected: Each section renders an empty-state placeholder in Thai (no crash, no broken layout)
  - case: Visitor opens home tab while offline
    expected: The page surfaces a Thai error message via the standard envelope and recovers when connectivity returns
requirement_refs: [CAT-001, "Frontend Spec §1", "Frontend Spec §2", "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 home row.
---

# STORY_FRONTEND_HOME — home tab renders featured + categories + new arrivals (Thai labels, ฿ prefix, emerald CTAs, 5-tab bar)

## User narrative

As a **Guest**, I want **a Thai-language mobile home tab with featured products, categories, and new arrivals** so that **I can browse the storefront from my phone and tap into a product detail**.

## Why this story exists in EPIC_FRONTEND

The home tab is the first impression of the entire customer mobile webview. The visual contract (Thai, 390×844, emerald, ฿ prefix, 5-tab bar) is asserted here as a representative case for every subsequent route.

## Edge cases

| Case | Expected |
|---|---|
| Empty catalog | Empty-state placeholders in Thai |
| Offline | Standard envelope error in Thai; recovery on reconnect |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
