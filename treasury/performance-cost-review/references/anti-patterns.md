# Anti-patterns (never do)

Loaded on demand by `performance-cost-review`. A budget that commits any of these
is not a budget.

- Set a budget with no baseline and no NFR derivation — that is a wish, not a budget.
- Budget the total monthly bill instead of cost-per-unit — it cannot tell efficiency from growth.
- Average a target across the whole service instead of per user journey — it hides the journey that is failing.
- Scale capacity or jump to a bigger instance before exhausting caching, query/index, cardinality, and rightsizing levers.
- Add an index or a high-cardinality metric dimension without budgeting its write / cost amplification.
- Divide a per-outcome cost by attempts instead of *verified* successes — a cheap-but-failing path then looks efficient.
- Report a single current number with no headroom or time-to-wall — leadership cannot plan a scaling decision from it.
- Page on every threshold breach or cosmetic trend — guardrails are tiered (alert -> cap -> breaker); pages are High+ only.
