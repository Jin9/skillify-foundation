# Gates and safety: gate matrix and command-safety policy

These cross-stage tables are shared by all five child skills. They are the authoritative copy — the
`SKILL.md` body points here. The "Never-do guardrails" prose stays in the body.

## Gate matrix
Gates sit where an action is hard to reverse or has high blast radius — not uniformly.

| Gate | After | Action gated | Reversibility | Blast | Type | Approver |
|------|-------|--------------|---------------|-------|------|----------|
| G1 | S1 | confirm scope | reversible draft | anchors all downstream | sync named | BA / PM |
| G2 | S2 | accept Story Set | reversible | low | async review | BA |
| G3 | S3 | resolve a governance blocker | hard to reverse | high (legal/privacy) | sync named | BA + SME (Legal/DPO/Compliance) |
| G4 | S4 | feasibility verdict | sets scope & cost | medium–high | sync named | Tech Lead |
| G5 | S5 | accept handoff bundle | engineering builds on it | high | sync named | Tech Lead (owner of record) |

An agent never passes an irreversible or control-plane gate on its own confidence; confidence only tunes whether a *reversible* step gets async review.

## Command-safety policy (enforced at the tool layer, not by the prompt)
| Tier | Actions |
|------|---------|
| ALLOW | read the raw requirement; classify/parse it; draft any stage contract; render a diagram; write inside `output/` |
| CONFIRM | publish the Handoff Bundle; mark a blocker resolved; set `state: ready-for-tl`; write to a shared backlog; override an upstream contract's scope |
| DENY | echo real PII; auto-resolve a governance blocker; write outside `output/`; call a non-allowlisted tool; modify another node's contract or its own permissions |
