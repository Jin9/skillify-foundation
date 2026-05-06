# implement gate review block

```markdown
## implement gate review for <workflow_id>

**Inputs**
- Plan:     .agent/stages/plan.md (<n> lines)
  Sections: <h1/h2 list>
- Critique: .agent/stages/critique.md
  P1: <count>, P2: <count>, P3: <count>
- Cap:      $<spent> / $<cap>  (remaining $<rem>)
- Sandbox:  <on|off>  (required by signals: <yes|no>)

### P1/P2 cross-walk

| Sev | Issue | Tag | Mitigation |
|-----|-------|-----|------------|
| P1  | <title> | <addressed|partial|unaddressed|accepted> | <plan.md:N "quote"> |
| P2  | <title> | ... | ... |
| ... |       |     |     |

### Five-question check

| # | Question | Answer |
|---|----------|--------|
| 1 | Plan read end-to-end? | <user-confirms> |
| 2 | Critique blockers addressed? | <yes|partial|no> |
| 3 | Cap tight enough? | <yes|loose|too-tight> |
| 4 | Sandbox set when needed? | <yes|n/a|missing> |
| 5 | Ready to babysit / accept? | <user-confirms> |

### Recommendation

**<APPROVE|REJECT>** — <one-sentence reason>

### Command (you type this)

```bash
# Approve:
just approve implement

# Reject (with reason):
just reject implement "<concise reason>"
```
```

## Filling rules

- Every P1 from the critique must appear in the cross-walk table. Drop
  none.
- The recommendation is the *worst* answer across questions:
  - Any `unaddressed P1` → REJECT.
  - Sandbox required but missing → REJECT.
  - Q3 = loose AND Q5 = "not ready to babysit" → REJECT.
  - All else → APPROVE.
- Print exactly one approve and one reject command. Do not run either.
