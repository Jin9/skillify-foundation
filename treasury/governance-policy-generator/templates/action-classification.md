# Action classification (SKELETON)

Complete this table BEFORE writing any allow rule in policy.rego. Two axes:
reversibility (undoable by a Git revert or equivalent) times blast radius
(resources/users affected if the action is wrong). Auto-approve eligibility
is granted ONLY to high-reversibility / low-blast-radius actions; everything
else is human review until a tested, audited rollback story exists.

If reversibility is unknown, mark Irreversible.

| Action | Reversibility (Reversible / Irreversible) | Blast radius (Low / Medium / High) | Rollback story (for irreversible) | Auto-approve eligibility (Eligible / Human review / Human review only) |
|---|---|---|---|---|
| `<read tracked file>` | Reversible | Low | n/a | Eligible |
| `<edit tracked file>` | Reversible | Low | git revert | Eligible |
| `<git commit (non-protected branch)>` | Reversible | Low | git revert | Eligible |
| `<git branch create (agent/ namespace)>` | Reversible | Low | branch delete | Eligible |
| `<dependency patch/minor, CI green>` | Reversible | Low | revert lockfile | Eligible (HUMAN REVIEW gate) |
| `<package install>` | Irreversible | Medium | uninstall + registry check | Human review (registry-allowlist gated) |
| `<git push to protected branch>` | Irreversible | High | `<define>` | Human review only |
| `<deploy>` | Irreversible | High | `<define rollback target>` | Human review only |
| `<external API call>` | Irreversible | `<assess>` | `<none — cannot undo>` | Human review only |
| `<DB migration>` | Irreversible | High | `<down migration tested?>` | Human review only |
| `<ADD EVERY OTHER INVENTORIED ACTION>` | | | | |

## Notes

- Three-tier permission separation: observe (always) / recommend (no
  immediate effect) / execute (only inside the auto-approve envelope).
- Promotion to widen the envelope is telemetry-gated (reversal rate below a
  pre-defined threshold), never intuition-gated.
- Every "Eligible" row must map to exactly one conditioned allow rule in
  policy.rego; no row maps to an unconditioned allow.
