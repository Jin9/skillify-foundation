# Secure Fix and Validation Rules

Apply these rules at Steps 8 and 9 of the review pipeline. A fix block without a validation block is incomplete.

## Secure fix rules

1. **Compilable in target stack.** Go fixes compile against current `go.mod`; YAML passes `kubectl --dry-run=server` or `kong config parse`; SQL passes MySQL syntax check.
2. **Minimal diff.** Smallest patch that closes the gap. No drive-by refactors.
3. **Stack-idiomatic.** Gin middleware not handwritten `http.Handler` glue; `database/sql` placeholders not `fmt.Sprintf`; declarative Kong plugins not custom Lua unless declarative cannot express the control.
4. **Defense in depth.** Two layers where feasible (Kong rate-limit + Gin rate-limit; NetworkPolicy + service-mesh authz).
5. **Specific maintained library + version + last-checked CVE date.** Examples: `github.com/go-jose/go-jose/v3`, `github.com/casbin/casbin/v2`.
6. **Environment-aware.** Call out SIT/UAT/PRD differences explicitly.
7. **Reversibility.** State backward-compat; include migration plan if wire/schema changes.
8. **No fix-by-disable.** Never recommend lowering a security policy to make code work.
9. **Cost note when material** (Vault, mesh, KMS) so the architect can weigh trade-off.
10. **Cite the standard the fix satisfies.**

## Validation rules

Every finding includes a validation block.

1. **Unit / integration test** asserting the bad behavior is now rejected (e.g., `TestBOLAOnLoanGet_ReturnsForbidden`). Provide the test signature.
2. **Negative test** asserting the legitimate path still works.
3. **Static check** with named rule (`gosec G201`, `govulncheck`, `staticcheck`, `kube-linter no-read-only-root-fs`, `checkov CKV_K8S_*`, `trivy config`, `kong config parse`, `apisix-cli check`).
4. **Manual verification** for things that resist automation (`kubectl auth can-i ...` as service account, expect `no`).
5. **Regression guard** preferring a CI-failing policy (OPA/Kyverno/Conftest) over memory-based discipline.
6. **Observability check** — what log line, metric, or alert confirms the control is active in PRD.
7. **Rollout gate** for high-blast-radius fixes (NetworkPolicy, RBAC tightening): staged SIT → UAT → canary PRD with explicit rollback condition.

## Cross-cutting summary block

End every multi-finding review with this shape:

```markdown
## Cross-cutting summary

**Top findings (by impact):**
1. [SEV-#] <title> — <one-line why>
2. ...
3. ...

**Systemic patterns:**
- <pattern name> — observed in <files/keys> — root cause: <one line>

**Coverage caveats:**
- <what was not reviewed and why> (e.g., upstream Kong config not provided)
```

The summary surfaces systemic root causes ("auth middleware missing on 4 of 7 routes — root cause: Gin router-group misuse") and tells the reader what is *not* covered.
