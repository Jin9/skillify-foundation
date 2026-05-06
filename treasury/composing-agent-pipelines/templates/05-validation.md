# 05 — Validation (claim verdicts)

**Validated**: [ISO timestamp]
**Inputs**: `03-analysis.md`, `04-review.md` (if exists)

Each P1 and P2 claim from the analysis gets one section. Verdict is one of: `confirmed`, `refuted`, `unverifiable`.

---

## Claim C1

**Claim text** (verbatim from analysis): "[claim]"

**Cited sources**: `02-evidence/q1.md#F1`

**Verdict**: confirmed | refuted | unverifiable

**Source quote** (re-fetched and verbatim):
```
[quoted lines from the source, with file:line]
```

**Reasoning**: [1–3 sentences explaining how the source supports / contradicts / fails to address the claim]

---

## Claim C2

**Claim text**: "[claim]"

**Cited sources**: `02-evidence/q3.md#F2`

**Verdict**: confirmed | refuted | unverifiable

**Source quote**:
```
[quoted lines]
```

**Reasoning**: [explanation]

---

## Summary

| Verdict       | Count |
|---------------|-------|
| confirmed     | [n]   |
| refuted       | [n]   |
| unverifiable  | [n]   |

## Refuted P1 claims

If any P1 claim was refuted, list it here. The pipeline must NOT proceed to Decide based on these claims.

- C1: [claim summary] → refuted by `path/to/source:line`.
