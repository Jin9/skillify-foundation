# RACI & role ownership

RACI assigns, for each checkpoint, who is **R**esponsible (does the work), **A**ccountable (owns the outcome,
final say), **C**onsulted (two-way input before the decision), and **I**nformed (told after). Use it so the BA
knows what to produce versus what to route to a human.

## The one-Accountable rule
**Exactly one A per row.** Two A's means no one is accountable. R can be shared; C and I can be many; A is singular.
If you cannot name the single A, that is an open question to resolve before the row is "done".

## Roles in a software BA flow
- **BA (Business Analyst)** — owns problem framing, as-is/to-be analysis, requirements, and the BA artifacts (BRD/SRS, user stories, RACI).
- **Tech Lead / Architect** — owns technical feasibility, effort, solution/architecture shape, NFR realism.
- **PM (Product/Project Manager)** — owns priority, scope trade-offs, timeline, resourcing, and sponsor sign-off.
- **QA** — owns testability of acceptance criteria, test strategy, and coverage view.
- **SME / Domain expert** — validates that the as-is and business rules are correct.
- **Sponsor / Stakeholder** — funds the effort and approves the final requirements.

## Typical ownership by checkpoint
| Checkpoint | R | A | C | I |
|---|---|---|---|---|
| Problem statement & goals | BA | PM | SME, Tech Lead | Stakeholders |
| As-is process map | BA | BA | SME, Tech Lead | PM |
| Functional requirements / user stories | BA | BA | SME, QA | Tech Lead, PM |
| Acceptance criteria (testable) | BA | QA | BA | Tech Lead |
| Non-functional requirements (targets) | BA | Tech Lead | BA, PM | QA |
| To-be process | BA | BA | Tech Lead, SME | PM |
| Solution / architecture shape | Tech Lead | Tech Lead | BA | PM |
| Feasibility & effort | Tech Lead | Tech Lead | BA | PM |
| Priority / MVP scope | BA | PM | Tech Lead, SME | QA |
| Validation & traceability | BA | BA | QA, SME | PM |
| Final sign-off | PM | Sponsor | BA, Tech Lead | Team |

Adapt names to the actual org; keep the one-A rule.

## How the two modes use this map
- **Autonomous BA:** for rows where the BA is A, produce the item directly. For rows where the BA is *not* A
  (feasibility, effort, architecture, sign-off), insert a labeled `[NEEDS INPUT — <role>]` placeholder and ask the
  user to route it to that owner. Never fill another role's A with a guess.
- **Co-pilot:** draft BA-owned rows and propose them to the human BA; for non-BA rows, note the typical owner as a
  suggestion and defer. Final sign-off always defers to the human BA.

## Anti-patterns
- Two A's on one row (or zero).
- The BA marking itself A for feasibility, effort, or architecture.
- Treating "Informed" people as if they approved — I is one-way, after the fact.
- Filling a non-BA owner's verdict with a fabricated value instead of a placeholder.
