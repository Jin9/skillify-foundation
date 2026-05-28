# BA Output Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [BA](../roles.md) · **`template_version`:** 1.0.0

Authority: BA owns this artifact. Tech-Lead, Tech-Designer, Reviewer-L1, Reviewer-L2 may not produce or modify requirements; if a downstream role thinks a requirement is missing, the routing tag is `requirements_gap` (back to BA), not gap-fill.

---

## Fields

- `template_version` — string (semver), required
- `workflow_name` — string, required
- `purpose` — single sentence, required (≤ 200 chars)
- `in_scope` — list of testable behaviors, required (≥ 1)
- `out_of_scope` — list of explicit non-goals, required (empty array allowed but suspicious)
- `acceptance_criteria` — list of `{given, when, then}` objects, required (≥ 1, each independently checkable)
- `edge_cases` — list of `{case, expected}` objects, required (empty array is a smell)
- `success_conditions` — string, required (end-to-end "done")
- `non_functional` — object with optional `latency`, `idempotency`, `failure_tolerance`
- `compliance_sensitive` — boolean, required (true triggers dual-pass BA + Reviewer-L2 → Opus Max; see [`../model-routing.md`](../model-routing.md))

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "ba-output/v1.0.0",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "template_version", "workflow_name", "purpose",
    "in_scope", "out_of_scope", "acceptance_criteria",
    "edge_cases", "success_conditions", "compliance_sensitive"
  ],
  "properties": {
    "template_version": { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "workflow_name":    { "type": "string", "minLength": 1 },
    "purpose":          { "type": "string", "minLength": 1, "maxLength": 200 },
    "in_scope": {
      "type": "array", "minItems": 1,
      "items": { "type": "string", "minLength": 1 }
    },
    "out_of_scope": {
      "type": "array",
      "items": { "type": "string", "minLength": 1 }
    },
    "acceptance_criteria": {
      "type": "array", "minItems": 1,
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["given", "when", "then"],
        "properties": {
          "given": { "type": "string", "minLength": 1 },
          "when":  { "type": "string", "minLength": 1 },
          "then":  { "type": "string", "minLength": 1 }
        }
      }
    },
    "edge_cases": {
      "type": "array",
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["case", "expected"],
        "properties": {
          "case":     { "type": "string", "minLength": 1 },
          "expected": { "type": "string", "minLength": 1 }
        }
      }
    },
    "success_conditions": { "type": "string", "minLength": 1 },
    "non_functional": {
      "type": "object", "additionalProperties": false,
      "properties": {
        "latency":           { "type": "string" },
        "idempotency":       { "type": "string" },
        "failure_tolerance": { "type": "string" }
      }
    },
    "compliance_sensitive": { "type": "boolean" }
  }
}
```

## Example (known-good)

Workflow: produce a structured review digest for a GitHub pull request.

```json
{
  "template_version": "1.0.0",
  "workflow_name": "pr-review-digest",
  "purpose": "Produce a structured review digest for a GitHub pull request that captures risk, blockers, and a reviewer-friendly summary.",
  "in_scope": [
    "Read PR title, description, diff, and CI status from the GitHub API",
    "Classify changed files by area (frontend, backend, infra, docs)",
    "Identify risk signals (large diffs, schema migrations, secret-touching files)",
    "Emit a markdown digest grouped by area with a top-line risk summary"
  ],
  "out_of_scope": [
    "Posting the digest as a PR comment (operator does this)",
    "Running tests or static analysis — only consume CI status",
    "Approving or rejecting the PR"
  ],
  "acceptance_criteria": [
    {
      "given": "A PR with a non-empty diff and passing CI",
      "when": "The workflow is invoked with the PR URL",
      "then": "The digest contains an area-grouped summary, a risk score, and a list of files touched per area"
    },
    {
      "given": "A PR with a failing CI status",
      "when": "The workflow is invoked",
      "then": "The digest top-line includes the failing job name and the risk score escalates to 'high' regardless of diff size"
    },
    {
      "given": "A PR that touches schema migrations",
      "when": "The workflow is invoked",
      "then": "The digest flags 'schema-migration' explicitly with the affected migration file paths"
    }
  ],
  "edge_cases": [
    { "case": "Draft PR (not ready for review)",                 "expected": "Workflow declines with a 'PR is draft' message; no digest emitted" },
    { "case": "PR with binary-only changes (no diff visible)",   "expected": "Digest notes 'binary-only change' and lists file paths; risk score 'manual-review-required'" },
    { "case": "GitHub API rate-limit hit mid-fetch",             "expected": "Workflow retries with exponential backoff (3 attempts); on final failure emits a 'fetch-incomplete' digest with what was retrieved" },
    { "case": "PR with zero changed files (empty diff)",         "expected": "Workflow declines with 'empty diff' message" }
  ],
  "success_conditions": "Operator receives a markdown digest within 30 seconds of invocation that captures area, risk, and per-area file lists, and accurately reflects the PR's CI status.",
  "non_functional": {
    "latency": "p95 < 30s end-to-end on a 200-file PR",
    "idempotency": "Same PR URL + same commit SHA must produce the same digest (modulo timestamp)",
    "failure_tolerance": "Single GitHub API failure tolerated via retry; cumulative failures degrade to 'fetch-incomplete' digest"
  },
  "compliance_sensitive": false
}
```

## Negative examples

These instances are valid JSON against the schema. Plan-Reviewer should still reject them.

**Negative #1 — Unfalsifiable criteria, missing edges, no non-functionals**

```json
{
  "template_version": "1.0.0",
  "workflow_name": "pr-review-digest",
  "purpose": "Help reviewers review PRs.",
  "in_scope": ["Review PRs"],
  "out_of_scope": [],
  "acceptance_criteria": [
    { "given": "A PR exists", "when": "Workflow runs", "then": "The digest is good" }
  ],
  "edge_cases": [],
  "success_conditions": "PRs get reviewed.",
  "compliance_sensitive": false
}
```

What Plan-Reviewer should catch:

1. `purpose` paraphrases the workflow name — Tech-Lead has nothing to architect against. Tag: `requirements_gap` (high).
2. `in_scope` is a single entry that restates the workflow name. Not testable behaviors. Tag: `requirements_gap` (high).
3. `acceptance_criteria[0].then` is unfalsifiable ("The digest is good"). Tag: `requirement_ambiguity` (high).
4. `edge_cases: []` — no real workflow has zero edges. Tag: `missing_edge_case` (high).
5. `non_functional` omitted — latency/idempotency/failure expectations unstated. Tag: `requirements_gap` (medium).

Expected routing: Plan-Reviewer → BA with stacked `requirements_gap` + `requirement_ambiguity` + `missing_edge_case`. Per [`../orchestrator.md`](../orchestrator.md) routing table, plan-stage cap is 1 — if BA can't repair, run terminates `HardFail`.

**Negative #2 — BA invented requirements (scope creep, unstated assumptions)**

```json
{
  "template_version": "1.0.0",
  "workflow_name": "pr-review-digest",
  "purpose": "Produce a structured review digest for a GitHub pull request.",
  "in_scope": [
    "Read PR data from GitHub API",
    "Emit a markdown digest",
    "Auto-merge the PR if all CI checks pass and the digest's risk score is low",
    "Post the digest as a PR comment via the GitHub API",
    "Notify the PR author on Slack with the digest summary"
  ],
  "out_of_scope": [],
  "acceptance_criteria": [
    {
      "given": "A PR with passing CI and low risk",
      "when": "The workflow runs",
      "then": "The PR is auto-merged within 5 minutes"
    }
  ],
  "edge_cases": [
    { "case": "Slack API is down", "expected": "Workflow continues; failure logged" }
  ],
  "success_conditions": "PRs are auto-merged when safe.",
  "compliance_sensitive": false
}
```

What Plan-Reviewer should catch:

1. `in_scope` adds three behaviors the user didn't ask for: auto-merge, PR commenting, Slack notification. BA is supposed to *structure* requirements, not invent them. Tag: `unstated_assumption` (high) routed to BA.
2. `out_of_scope: []` — no fence around the inventions; Tech-Lead has no signal on what was added vs requested.
3. `acceptance_criteria` is auto-merge-centric, redirecting the workflow's center of gravity away from the user's "produce a digest" ask. Tag: `unstated_assumption` (high).
4. `edge_cases` covers only the side-quest behaviors (Slack outage); core workflow edges (draft PR, empty diff, rate limit) absent. Tag: `missing_edge_case` (high).

Expected routing: Plan-Reviewer → BA. The hard rule "if it isn't in BA's spec, it isn't in the workflow" is being inverted here — BA is putting things *in* the spec the user didn't say. Same plan-stage cap = 1 applies.
