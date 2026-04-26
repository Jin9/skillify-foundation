# Frontend Security

## Contents
- Threat model
- XSS
- CSRF
- Auth-token storage
- Content Security Policy (CSP)
- Dependency / supply chain
- Sensitive data in the client
- High-risk patterns

## Low-risk use

Security-sensitive changes need explicit context. Do not change auth/session storage, CSP, CSRF strategy, analytics scrubbing, or third-party script loading without understanding the existing backend and deployment model.

## Threat model

The browser is hostile territory. Assume:
- Any value you didn't render server-side may be tampered with.
- The user may install a malicious extension.
- Third-party scripts you load can read your DOM and storage.

## XSS

- Default safe: React escapes by default. The danger is `dangerouslySetInnerHTML`.
- If you must render HTML, sanitize with `DOMPurify` first. Never trust server HTML "because we wrote it."
- URL props (`href`, `src`): validate against an allowlist of schemes (`http`, `https`, `mailto`). Reject `javascript:`.
- Avoid `eval`, `new Function()`, dynamic `import()` of user-controlled paths.

## CSRF

- For cookie-based auth: enforce `SameSite=Lax` (or `Strict`) on session cookies.
- Mutating endpoints require a CSRF token OR fetch with `credentials: 'include'` + same-site policy + custom header (`X-Requested-With`).
- For token-in-header auth (Bearer), CSRF is not applicable — but XSS becomes the bigger risk.

## Auth-token storage

| Storage | Risk |
|---|---|
| `localStorage` | Readable by any script — XSS = full takeover |
| `sessionStorage` | Same as above, scoped to tab |
| Cookie (`HttpOnly`, `Secure`, `SameSite`) | NOT readable by JS — preferred for session tokens |
| In-memory (React state / module var) | Lost on reload; safe from storage exfil |

Preferred default for new designs: HttpOnly cookie for session, in-memory for short-lived access tokens. Do not store new auth tokens in `localStorage` without an accepted threat model.

## Content Security Policy (CSP)

- Start from `default-src 'self'` and loosen only for required origins.
- `script-src` with nonces or hashes; avoid `'unsafe-inline'` and `'unsafe-eval'`.
- `frame-ancestors 'none'` (or specific origins) to prevent clickjacking.
- Report violations to a collector (`report-uri` / `report-to`) and triage them.

Next.js: configure via `next.config.js` headers. Vite: via your reverse proxy / CDN.

## Dependency / supply chain

- Pin versions; lockfile in CI.
- Use the repo's package manager audit regularly. Investigate HIGH / CRITICAL findings.
- Subresource Integrity (`integrity=`) for any third-party CDN script.
- Avoid `postinstall` scripts you haven't audited.
- Treat copy-pasted snippets from blog posts / Stack Overflow as untrusted code.

## Sensitive data in the client

- Don't render secrets the user shouldn't see (API keys, internal IDs) into the DOM, even hidden.
- Don't ship server-only env vars to the client. In Next.js: only `NEXT_PUBLIC_*` is safe.
- PII in the URL = PII in browser history, referrer headers, and analytics. Avoid.
- Logs sent from the client: scrub fields before send (PII, tokens).

## High-Risk Patterns

- `dangerouslySetInnerHTML` on user input without sanitization.
- Storing JWTs in `localStorage` without a documented threat model.
- `<a href={userInput}>` without scheme validation.
- Disabling React's escaping via string concatenation in JSX.
- `target="_blank"` without `rel="noopener noreferrer"`.
- Logging full requests / responses to a third-party tool without scrubbing.
- Hardcoding API keys in client bundles.
