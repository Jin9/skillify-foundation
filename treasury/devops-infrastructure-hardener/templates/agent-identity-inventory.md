# Agent Identity Inventory — <scope>

Every agent is registered as a unique non-human identity with a scope, a
named owner, and a lifecycle-event rotation trigger. No identity is shared
across agents or with a human engineer. This inventory is the basis for
prompt revocation and scoped audit.

## Inventory

| agent | non-human identity | scope / least-privilege | owner | rotation trigger |
|-------|--------------------|--------------------------|-------|------------------|
| <agent name / role> | <unique service account or workload identity id> | <minimum permission set: read/write, resource ids, time-bound> | <named human> | <deploy change \| config change \| anomaly alert \| scope change> |
| <agent name / role> | <unique id — never shared> | <least-privilege grant> | <named human> | <lifecycle event> |

## Per-agent trust-tier map (repeat block per agent)

### <agent name>

- **Identity:** <unique non-human identity; not shared>
- **Data-plane actions and tiers:** <action → Tier 1–4>
- **Control-plane actions:** <action → Tier 1 or 2 only, human approver: name>
- **Highest tier any action holds:** <1–4> (control-plane MUST be <= 2)
- **Owner:** <named human accountable for lifecycle>
- **Rotation trigger:** <lifecycle event, never calendar>

## Invariants

- [ ] Every agent has its own unique identity (no sharing).
- [ ] Every identity has a named human owner.
- [ ] Every identity has a least-privilege scope, not a broad grant.
- [ ] Every rotation trigger is a lifecycle event, not a calendar interval.
- [ ] No control-plane action for any agent exceeds Tier 2.
