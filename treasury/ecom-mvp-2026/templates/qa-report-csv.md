# QA Report CSV Template (L1 + L2)

**Parent:** [`../templates.md`](../templates.md)
**Owner roles:** [QA-L1](../roles.md), [QA-L2](../roles.md) · **`template_version`:** 0.1.0

The CSV sidecar to QA's JSON verdict. The user explicitly asked for a per-case report Dev can act on without ambiguity. Structured so a wrong test result has direct trace-back to story → AC → expected → actual → evidence.

**File locations:**
- QA-L1: `<workflow_root>/.squad-run/components/<svc>/qa-l1.csv`
- QA-L2: `<workflow_root>/.squad-run/layer-2/qa-l2.csv`

The JSON verdict (`qa-l1.json` / `qa-l2.json`) remains the orchestrator's routing input. The CSV is human-readable + Dev-actionable.

---

## CSV header (REQUIRED, exact column order)

```csv
case_id,scenario,verdict,expected,actual,ac_id,story_id,evidence
```

## Column semantics

| Column | Type | Required | Description |
|---|---|---|---|
| `case_id` | string | Yes | Stable identifier for the test case. Format `<svc>-<NNN>` for L1 (e.g., `identity-001`); `l2-<NNN>` for L2. |
| `scenario` | string | Yes | One sentence describing what the test exercises. Must be parseable in isolation (don't reference "the previous case"). |
| `verdict` | enum | Yes | `pass` \| `fail` \| `skipped` \| `error` |
| `expected` | string | Yes | What should happen, in concrete terms (HTTP status, code, value, etc.). |
| `actual` | string | Yes | What actually happened. For pass: same as expected (or "matched"). For fail: the divergence in concrete terms. For skipped/error: reason. |
| `ac_id` | string | Yes for pass/fail | The AC reference, e.g., `STORY_AUTH_LOGIN/AC-1` or `STORY_AUTH_LOGIN/AC-2`. Empty for cross-cutting checks not tied to a single AC. |
| `story_id` | string | Yes for pass/fail | The owning `STORY_<SLUG>` ID. Empty for cross-cutting. |
| `evidence` | string | Yes | Pointer to evidence: log line ref, screenshot path, request/response body file, test framework run ID. NEVER inline a giant blob — link to a file. |

## CSV escaping rules

- All cells use `,` as separator and `"` as quote character.
- Cells containing `,`, `"`, or newlines MUST be quoted.
- Inner `"` doubled (`"" `).
- Newlines in cells preserved as `\n` literal characters within a quoted cell (avoid where possible — keep cells single-line).

---

## Sidecar JSON (the orchestrator's routing input)

```json
{
  "template_version": "0.1.0",
  "review_kind": "qa-l1" | "qa-l2",
  "component_name": "<svc>",
  "skills_invoked": ["simplify"],
  "verdict": "pass" | "code_mismatch" | "spec_incomplete" | "spec_wrong"          // L1
            | "pass" | "contract_violation" | "criteria_unmeetable" | "criteria_unmet",  // L2
  "summary": "1-paragraph",
  "csv_path": ".squad-run/components/<svc>/qa-l1.csv",
  "findings": [
    {
      "id": "QA-IDENT-001",
      "tag": "code_mismatch | spec_incomplete | spec_wrong",   // L1
      "locus": "...",
      "description": "...",
      "recommendation": "..."
    }
  ]
}
```

The JSON references the CSV via `csv_path`. The orchestrator routes on the JSON's `findings[]`; the CSV is for human consumption.

---

## Example CSV (L1, identity service)

```csv
case_id,scenario,verdict,expected,actual,ac_id,story_id,evidence
identity-001,Register a new user with valid email + password,pass,"HTTP 200, code=SUCCESS, data.userId is UUID","matched",STORY_AUTH_REGISTER/AC-1,STORY_AUTH_REGISTER,tests/output/register_001.log
identity-002,Register with duplicate email,pass,"HTTP 409, code=EMAIL_ALREADY_REGISTERED","matched",STORY_AUTH_REGISTER/AC-2,STORY_AUTH_REGISTER,tests/output/register_002.log
identity-003,Login with correct credentials,pass,"HTTP 200, accessToken is ES256 JWT, refreshToken issued","matched",STORY_AUTH_LOGIN/AC-1,STORY_AUTH_LOGIN,tests/output/login_001.log
identity-004,Login with wrong password,fail,"HTTP 401, code=AUTH_INVALID, message=""Invalid email or password.""","HTTP 401, code=AUTH_INVALID, message=""Wrong password""",STORY_AUTH_LOGIN/AC-2,STORY_AUTH_LOGIN,tests/output/login_002.log
identity-005,Login with unknown email returns same message as wrong password,fail,"both calls return identical {code, message}","unknown-email returns 404 NOT_FOUND",STORY_AUTH_LOGIN/AC-3,STORY_AUTH_LOGIN,tests/output/login_003.log
identity-006,Refresh rotates jti and returns new pair,skipped,—,"endpoint stubbed (501 NOT_IMPLEMENTED_MVP)",STORY_AUTH_REFRESH/AC-1,STORY_AUTH_REFRESH,
identity-007,Bcrypt hash uses cost 12,pass,"hash starts with $2a$12$","matched",—,STORY_AUTH_REGISTER,tests/output/register_001.log
identity-008,Login latency constant-time vs unknown-email,error,"|t_wrong_pw - t_unknown_email| < 50ms","measurement framework not wired (manual: 0ms vs 30ms; informal pass)",STORY_AUTH_LOGIN/AC-2,STORY_AUTH_LOGIN,
```

Reading this CSV, Dev sees:
- `identity-004` failed — the message text doesn't match AC-2 (Dev returned `"Wrong password"` instead of the locked `"Invalid email or password."`).
- `identity-005` failed — Dev returned 404 on unknown email, leaking enumeration. AC-3 explicitly forbids this.

These two findings route back to Dev with precise file/handler context.

---

## Sidecar JSON example (paired with CSV above)

```json
{
  "template_version": "0.1.0",
  "review_kind": "qa-l1",
  "component_name": "identity",
  "skills_invoked": ["simplify"],
  "verdict": "code_mismatch",
  "summary": "Login wrong-password and unknown-email paths violate AUTH-007. Wrong-password returns a different message than spec; unknown-email returns 404 instead of the same generic AUTH_INVALID 401.",
  "csv_path": ".squad-run/components/identity/qa-l1.csv",
  "findings": [
    {
      "id": "QA-IDENT-001",
      "tag": "code_mismatch",
      "locus": "app/identity/handler_login.go:60-65 (wrong-password branch)",
      "description": "Returns message=\"Wrong password\" instead of locked \"Invalid email or password.\" (AUTH-007).",
      "recommendation": "Change the literal message string to match TD spec; ensure the same string for both wrong-password and unknown-email branches."
    },
    {
      "id": "QA-IDENT-002",
      "tag": "code_mismatch",
      "locus": "app/identity/handler_login.go:42-50 (unknown-email branch)",
      "description": "Returns HTTP 404 NOT_FOUND for unknown email; spec requires HTTP 401 AUTH_INVALID with the same message as wrong-password.",
      "recommendation": "Collapse the branches: any auth failure returns 401 AUTH_INVALID with the locked message."
    }
  ]
}
```

---

## Negative examples

### Negative #1 — `expected` and `actual` are identical on a fail

```csv
case_id,scenario,verdict,expected,actual,ac_id,story_id,evidence
identity-099,Some test,fail,"works correctly","works correctly",STORY_AUTH_LOGIN/AC-1,STORY_AUTH_LOGIN,
```

What's wrong:

1. If `expected == actual`, the verdict can't be `fail` — schema/sanity violation.
2. `expected` is unfalsifiable ("works correctly") — Dev can't act. QA-L1 must use concrete language (HTTP code, body field, value).

### Negative #2 — Missing `evidence` on a fail

```csv
case_id,scenario,verdict,expected,actual,ac_id,story_id,evidence
identity-100,Login fails,fail,"HTTP 200","HTTP 401",STORY_AUTH_LOGIN/AC-1,STORY_AUTH_LOGIN,
```

What's wrong:

1. `evidence` empty on a `fail` — Dev can't reproduce or examine. QA-L1 must point to a log file, request/response capture, or test run.
2. Also: `expected: HTTP 200` for a "Login fails" scenario is contradictory — the scenario is mislabeled.
