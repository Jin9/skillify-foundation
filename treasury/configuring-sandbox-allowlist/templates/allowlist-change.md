# PR description: sandbox allowlist change

```markdown
## sandbox: <add|remove|tighten> <host or pattern>

### Why

<one paragraph: which workflow / stage / partner needs this; what
breaks without it; why no narrower alternative was available>

### Diff

```diff
--- a/docker/sandbox-proxy/filter
+++ b/docker/sandbox-proxy/filter
@@ <section>
+^api\.example\.com$
```

### Refusal-rule check

- [ ] Anchored with `^...$`
- [ ] Literal dots escaped (`\.`)
- [ ] No multi-label wildcards (`.*` in subdomain position)
- [ ] Not an IP literal
- [ ] Not a metadata service or `*.internal` / `*.local` host
- [ ] Host resolves (`getent hosts <host>`)
- [ ] Not a duplicate of an existing entry

### Verification (run after merge)

```bash
docker compose -f docker/sandbox-compose.yml build proxy --no-cache
docker compose -f docker/sandbox-compose.yml up -d proxy
just doctor --full
```

`just doctor --full` exercises the egress-test path: the dispatcher
attempts to fetch a non-allowlisted host and aborts implement if the
fetch unexpectedly succeeds.

### Rollback

```bash
git revert <this-commit>
docker compose -f docker/sandbox-compose.yml build proxy --no-cache
docker compose -f docker/sandbox-compose.yml up -d proxy
```

### Reviewer ask

This change touches the security boundary for sandboxed implement.
Please review:

- The refusal-rule checklist above.
- The "why" paragraph — is the reason still valid in 90 days?
- Whether a narrower pattern would meet the same use case.
```
