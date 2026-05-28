# RACI matrix — <project / feature name>

> **R** = Responsible (does the work) · **A** = Accountable (owns the outcome, final say) ·
> **C** = Consulted (two-way input before) · **I** = Informed (told after).
> **Exactly one A per row.** R may be shared; C and I may be many. See references/raci-ownership.md.

| Checkpoint | BA | Tech Lead / Architect | PM | QA | SME | Sponsor |
|------------|----|-----------------------|----|----|-----|---------|
| Problem statement & goals | R | C | A | — | C | I |
| As-is process map | A/R | C | I | — | C | — |
| Functional requirements / stories | A/R | C | I | C | C | — |
| Acceptance criteria (testable) | R | C | — | A | — | — |
| Non-functional targets | R | A | C | C | — | — |
| To-be process | A/R | C | I | — | C | — |
| Solution / architecture shape | C | A/R | I | — | — | — |
| Feasibility & effort | C | A/R | I | — | — | — |
| Priority / MVP scope | R | C | A | I | C | I |
| Validation & traceability | A/R | C | I | C | C | — |
| Final sign-off | C | C | R | — | — | A |

Adapt role names to the real org and add rows as needed. Replace any cell the BA cannot decide with
`[NEEDS INPUT — <role>]` rather than guessing the owner.

## Validity check
- [ ] Every row has **exactly one A**
- [ ] No checkpoint left without an R
- [ ] BA is not marked A for feasibility, effort, or architecture
- [ ] Sign-off A is a human decision-maker
