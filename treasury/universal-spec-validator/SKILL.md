---
name: universal-spec-validator
description: >
  Validate an agent spec (tool/function JSON Schema, MCP tool manifest,
  OpenAPI tool def, SKILL.md frontmatter, or data contract) as a CI /
  pre-commit ENFORCEMENT GATE that fails the build before malformed tool calls
  or unsafe execution ship. Three axes: cross-model portability drift, unsafe
  command surface, breaking schema evolution vs a baseline. Use when the user
  asks to "add a spec-validation gate to CI", "fail the build if this tool
  schema breaks existing callers", "check these MCP tool / skill specs for
  cross-model drift before merge", "pre-commit validate our agent
  command/schema specs", or "is this schema change backward-compatible for the
  tool registry". Produces spec-validation.json + spec-validation.md and a
  non-zero exit on any blocking finding. Do NOT use for general code/app
  security review, runtime sandboxing or executing the commands,
  authoring/generating the schema or skill, or generic config/lint.
---

# Universal Spec Validator

## Purpose

Gate agent specs in CI/pre-commit: deterministically detect cross-model
portability drift, unsafe command surface, and breaking schema evolution, then
fail the pipeline before a broken or dangerous contract reaches an agent. This
skill owns the *gate*, not the fix — it reports and blocks; it never rewrites
the spec or executes anything it inspects.

## When to use this skill

- Use when: "add a spec-validation gate to CI" / "pre-commit validate our agent command/schema specs".
- Use when: "fail the build if this tool schema breaks existing callers" / "is this schema change backward-compatible for the tool registry".
- Use when: "check these MCP tool / skill specs for cross-model drift before merge".
- Do NOT use when: the request is general code/app security review, runtime sandboxing or *running* the commands, authoring/generating the schema or skill, or generic YAML/JSON lint. Hand those to the appropriate dedicated skill.

## Inputs

| Input | How supplied | Notes |
|-------|--------------|-------|
| Spec target(s) | file path(s) or glob | JSON Schema, MCP tool manifest, OpenAPI tool def, `SKILL.md` (frontmatter), or JSON/YAML data contract |
| Baseline (optional) | path to the previously-shipped spec/dir | required for schema-evolution checks; absent → those checks report `info: skipped` |
| Config (optional) | `.spec-validator.yaml` | severity-to-gate tiers, ignore rules, target model families, compatibility mode |

## Workflow

1. **Resolve inputs.** Identify the spec target(s) and, if present, the baseline and `.spec-validator.yaml`. If no spec target is given, ask once; do not guess.
2. **Run the gate script — do not hand-evaluate.** Pass/fail logic is deterministic and lives in `scripts/`. Run:
   `python3 scripts/validate_spec.py <spec-or-glob> [--baseline PATH] [--config PATH] [--out DIR] [--format both]`
   The script applies the three rule sets (see References) and writes the reports.
3. **Do not modify the spec.** This skill is read-only on every file it inspects. It never executes a command surface it is validating.
4. **Surface the result.** Report the exit code, the blocking-finding count per axis, and the path to `spec-validation.md`. The exit code IS the gate signal — wire it directly into CI (`&&` / job step) or a pre-commit hook (see `scripts/ci-gate.sh`).
5. **Tier the response (do not blanket-fail on noise).** Severity-to-gate mapping comes from `.spec-validator.yaml` or the default in `references/severity-and-gate-policy.md`: blocking findings exit non-zero; `warn`/`info` print but exit 0. Recommend pinning + risk-tiered automation, never an unconditional kill.

## Output contract

The skill (via `scripts/validate_spec.py`) produces, in `--out` (default `./`):

- `spec-validation.json` — machine-readable: `{summary, exit_code, findings:[{axis, id, severity, gate, spec, locator, message, rule_ref}]}`. Axes: `portability`, `command-safety`, `schema-evolution`.
- `spec-validation.md` — human report rendered from the same findings (shape: `templates/spec-validation-report-template.md`).
- Process exit code: `0` = no blocking findings; `1` = at least one blocking finding; `2` = usage/parse error. CI must treat non-zero as a failed gate.

No other files are written. No network calls, no command execution, no LLM calls — the gate is deterministic and CI-safe.

## Constraints

- DO NOT edit, normalize, or "auto-fix" the spec — gate only; remediation is the caller's job.
- DO NOT execute, sandbox, or fetch anything described by the spec (that is runtime safety, out of scope).
- DO NOT add findings the rule sets do not define; keep the gate reproducible across runs and machines.
- DO NOT auto-disable on every signal — honor the risk-tiered severity-to-gate policy; a false alarm must degrade to `warn`, not a blanket build kill.
- DO NOT duplicate the rule taxonomies here; they live one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals the folder, kebab-case, no XML, description under 1024 chars with triggers + negatives.
- [ ] `python3 scripts/validate_spec.py --help` runs; gate exits non-zero on a known-bad fixture, zero on a clean one.
- [ ] Every finding carries an `axis`, `severity`, `gate`, and `rule_ref` into a `references/` rule.

## References

- Cross-model portability drift taxonomy and checks: `references/portability-rules.md`
- Safe-vs-breaking schema-evolution rulebook: `references/schema-evolution-rules.md`
- Unsafe command-surface checks: `references/command-safety-rules.md`
- Severity tiers, exit-code and risk-tiered gate policy: `references/severity-and-gate-policy.md`
- Report shape: `templates/spec-validation-report-template.md`; config: `templates/spec-validator.example.yaml`
