# Pattern rules

The filter file uses tinyproxy's POSIX extended regex matching.
`FilterDefaultDeny Yes` and `FilterExtended On` are set in
`tinyproxy.conf`, so:

- Default action is **deny**.
- Patterns must match against the full hostname.
- Patterns are extended regex (so `(a|b)` and `\.` work).

## Required pattern shape

- Anchor both ends: `^...$`.
- Escape every literal dot: `\.`.
- Either one hostname or an explicit OR-set joined with `|`.

## Good examples

```
^api\.anthropic\.com$
^api\.openai\.com$
^(api|cdn)\.example-bank\.com$
^proxy\.golang\.org$
```

## Bad examples and why

```
api.openai.com           # not anchored; matches "api.openai.com.evil.tld"
^api.openai.com$         # unescaped dot; "api1openai1com" matches
^.*\.openai\.com$        # multi-label wildcard; "foo.bar.openai.com" passes
^api\..*$                # multi-label wildcard; "api.evil.tld" passes
^[0-9.]+$                # IP literal; bypasses hostname attribution
```

## Refusal table

| Pattern signal | Reason | Safer alternative |
|---|---|---|
| `.*` in subdomain position | subdomain takeover bypass | enumerate exact subdomains: `^(api|cdn|files)\.example\.com$` |
| `^.*\.<tld>$` | matches every host on the TLD | enumerate the actual hosts |
| Bare IP literal | bypasses DNS attribution | use a hostname; if no DNS exists, escalate to team-lead |
| Metadata IP / host (`169.254.169.254`, `metadata.google.internal`) | cloud-credential exfil vector | refuse; if the runner truly needs cloud creds, redesign with IRSA / workload identity |
| `*.internal`, `*.corp`, `*.local`, `*.cluster.local` | runner has no business hitting internal services | refuse |
| Hosts that resolve to `127.0.0.0/8`, `10.0.0.0/8`, `172.16/12`, `192.168/16` | bypass via Docker networking | refuse |
| Comment-only entry without a hostname | doesn't allow anything; clutter | drop |

## Section convention

Group entries by purpose. The default `filter` ships with these
sections — keep new entries in the matching one or add a new section
header for a new purpose:

```
# Model APIs
# Local LiteLLM proxy
# Package managers
# Public source hosts
# Partner banks      (example new section)
```

## Multi-host CONNECT considerations

The proxy permits CONNECT only to ports 443 and 563 (`tinyproxy.conf`).
If a target host needs another port, that is a `tinyproxy.conf` change
— **not** a filter change — and is out of scope for this skill.
