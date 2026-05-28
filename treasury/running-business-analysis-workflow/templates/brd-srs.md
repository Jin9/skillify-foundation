# BRD / SRS — <project / feature name>

> Mode: <autonomous | co-pilot>   ·   Stage reached: <1–6>   ·   Author/BA: <name>   ·   Date: <date>
> Status: <draft | in-review | validated | signed-off>
> Mark anything the BA cannot determine as `[NEEDS INPUT — <role>]`. Never fabricate.

## 1. Problem & goals  (Stage 1)
- **Problem statement (1 paragraph):**
- **Measurable goals / success metrics:** (metric + target + timeframe)
  - <metric> → <target> by <when>
- **Stakeholders:** (and sponsor)
- **In scope:**
- **Out of scope:**
- **Constraints & assumptions:**

## 2. Current state — as-is  (Stage 2)
- **Process today (actors → steps → systems):**
- **Pain points / workarounds:**
- **Current data & integrations:**
- **Baseline metrics (volumes today):**
- **Must-not-break:**

## 3. Requirements  (Stage 3)
### 3.1 Functional requirements
| ID | User story (As a… I want… so that…) | Priority (MoSCoW) | Acceptance criteria ref | Owner |
|----|--------------------------------------|-------------------|-------------------------|-------|
| FR-1 | | | see user-story.md | |

### 3.2 Non-functional requirements  (every row needs a number — see references/nfr-catalog.md)
| ID | Category | Target (metric + unit) | Condition | Measured how | Owner |
|----|----------|------------------------|-----------|--------------|-------|
| NFR-1 | Performance | p95 < 300 ms | at 50 req/s | gateway | Tech Lead |
| NFR-2 | Availability | 99.9% / month | — | SLO dashboard | Tech Lead |

## 4. Future state — to-be  (Stage 4)
- **To-be process:**
- **Gap (as-is → to-be):**
- **New data / integrations / interfaces:**
- **MVP slice (core-problem cut test):**
- **Deferred to later phases:**

## 5. Feasibility & risk  (Stage 5 — verdicts from human owners only)
| Requirement ID | Feasibility (feasible / needs-spike / infeasible) | Rough effort | Risks / dependencies | Verdict owner |
|----------------|---------------------------------------------------|--------------|----------------------|---------------|
| FR-1 | `[NEEDS INPUT — tech lead]` | | | Tech Lead |

## 6. Validation & sign-off  (Stage 6)
- **Validation results / conflicts resolved:**
- **Traceability:** every requirement links goal → requirement → AC (see RTM / raci-matrix.md).
- **Open items:**
- **RACI:** see templates/raci-matrix.md
- **Sign-off (human decision):**
  - Sponsor: __________  Date: ____    PM: __________  Date: ____    Tech Lead (technical sections): __________  Date: ____
