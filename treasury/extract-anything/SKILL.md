---
name: extract-anything
description: >
  Extract ANY source — a document, a prior workflow stage's output, a spec,
  code, or notes — into a single LEAN, chainable JSON contract that preserves
  the source's salient context with bounded, inline-recorded information loss,
  so decoupled agentic-workflow stages hand off without re-reading the source.
  Conforms to a caller-provided target_schema when given, otherwise infers a
  minimal shape. Output is one JSON envelope: the contract plus an inline
  _meta loss record. Use when the user asks to "extract this into a JSON
  contract", "pull everything from this source into a contract for the next
  stage", "turn this source into a lean chainable contract / handoff", or
  "extract X to JSON preserving context". Do NOT use for research-plan
  evidence (extract-findings), BA requirements (extract-brief-structure),
  code-to-rules-spec extraction (business-logic-extractor), or validating an
  existing contract (universal-spec-validator).
compatibility: claude-code, codex, copilot, gemini, antigravity
version: 0.1.0
---

# Skill: Extract Anything

## Purpose

Turn one arbitrary source into one lean JSON contract that preserves the source's salient context with bounded, explicitly recorded loss — so the next stage in an agentic workflow can act on the contract alone, without re-reading the source.

## Input

- `source` (required): the material to extract from — inline text, or a path/ref to a file. Any kind: prose, notes, a spec, code, structured data, or a prior stage's output.
- `target_schema` (optional): the downstream contract shape (a JSON Schema or a field list). Its presence selects **conform** mode; its absence selects **infer** mode.
- `focus` (optional): what the next stage cares about (a goal, a question, the consuming stage's name). Steers salience; never invents content.
- `tier_hint` (optional): `small` | `mid` | `frontier` — the model tier that will consume the contract. Lower tiers get flatter, smaller contracts.
- `idempotency_key` (optional): caller-supplied id echoed into `_meta` for trace correlation.

If `source` is missing or empty, stop and ask for it. Do not invent a source.

## Modes

| Mode | Selected when | Behavior |
|------|---------------|----------|
| `conform` | `target_schema` is provided | Map the source onto the target schema's fields. Populate every required field or mark it missing-with-reason in `_meta.dropped`. NEVER fabricate a value the source does not support. |
| `infer` | no `target_schema` | Design a minimal, flat contract whose fields capture the source's salient content. Prefer stable, predictable field names. Recommend conform mode in `_meta.notes` if a downstream contract already exists. |

See `references/schema-flexibility.md` for the conform-vs-infer decision rules and lean-contract heuristics.

## Workflow

1. **Resolve inputs and mode.**
   - Entry: a non-empty `source` is available.
   - Load the source (read the file if a path/ref was given).
   - Set `mode = conform` if `target_schema` is present, else `mode = infer`.
   - Exit: source content in hand and `mode` decided.

2. **Salience-first pass.**
   - Entry: source content loaded.
   - Identify what the consuming stage needs (use `focus` if given). Tag each unit of source content as **must-preserve**, **droppable**, or **out-of-scope**.
   - Apply the preserve-vs-drop rules in `references/loss-bounding.md`. Bound size: keep the contract flat and lean; push bulky spans behind a `source_ref` pointer rather than inlining them.
   - Exit: a salience map of must-preserve vs droppable content.

3. **Build the contract.**
   - Entry: salience map ready.
   - `conform`: place each must-preserve unit into its `target_schema` field. Leave a required field unfilled only when the source genuinely lacks it — then record it in `_meta.dropped` with reason `not-in-source`.
   - `infer`: design the minimal field set that carries the must-preserve content. Favor flat objects and stable names over deep nesting (deep nesting breaks context economy for downstream stages).
   - Carry no field the source does not support. See `references/anti-patterns.md`.
   - Exit: the `contract` object is populated.

4. **Bound the loss inline.**
   - Entry: `contract` populated.
   - Fill `_meta`: `mode`, `coverage` (`high|partial|low` — how much must-preserve content was captured), `dropped[]` (each `{ item, reason }`), `confidence` (`high|medium|low`), `source_ref` (how to reach the original — path/ref/anchor), and `provenance[]` (per non-obvious field, where in the source it came from).
   - Every intentional omission MUST appear in `dropped[]`. Silent loss is forbidden.
   - Exit: `_meta` complete; the same context is recoverable via `dropped` + `source_ref`.

5. **Self-validate.**
   - Entry: full envelope assembled.
   - Check: envelope has `contract` and a complete `_meta`; in conform mode every required target field is populated OR present in `dropped[]` with a reason; no field contains content unsupported by the source.
   - Run `scripts/validate_extraction.py <output.json>` when available; if it exits non-zero, fix the envelope and re-run. Do not rely on self-report for envelope integrity.
   - Exit: validation passes.

6. **Emit.**
   - Write the single JSON envelope (see Output Format) to `extraction.json` (or the caller's path).
   - Print a one-line summary: `mode=<m> coverage=<c> dropped=<n> confidence=<lvl>`.
   - Exit: envelope delivered.

## Output Format

Exactly one JSON object — the envelope. No prose around it.

```json
{
  "contract": { },
  "_meta": {
    "mode": "conform | infer",
    "source_ref": "path / ref / anchor to the original source",
    "coverage": "high | partial | low",
    "confidence": "high | medium | low",
    "dropped": [ { "item": "what was left out", "reason": "why" } ],
    "provenance": [ { "field": "contract field path", "from": "where in source" } ],
    "idempotency_key": "echoed if provided",
    "notes": "optional: e.g. recommend conform mode / flags for the consumer"
  }
}
```

- `contract` holds the extracted, lean, downstream-facing payload (shape per mode).
- `_meta` is the loss-bounding record; it travels with the contract so loss is auditable and the source is recoverable. The envelope schema lives in `schemas/output.json`; the input contract in `schemas/input.json`.

## Constraints

- DO NOT fabricate fields, values, or detail the source does not support — in either mode.
- DO NOT emit a lossy prose summary; the output is structured JSON with an explicit loss record, not a paraphrase.
- DO NOT inline bulky content; reference it via `source_ref` to protect context economy.
- DO NOT drop content silently; every omission goes in `_meta.dropped`.
- DO NOT validate, test, or design contracts — this skill only extracts. For validating/testing an existing contract use `universal-spec-validator` or `contract-testing-pact`; for designing a new one use `api-contract-design` / `befe-contract-design` / `data-modeling`.
- DO NOT take over a sibling extractor's job: research-plan evidence → `extract-findings`; BA requirements → `extract-brief-structure`; code → business-rules spec with file:line → `business-logic-extractor`.
- DO NOT duplicate the reference material inside this file.

## Validation

- [ ] Output is one JSON envelope with `contract` and a complete `_meta`.
- [ ] `mode` reflects whether a `target_schema` was supplied.
- [ ] In conform mode, every required target field is populated or listed in `dropped[]` with a reason.
- [ ] No fabricated content; every non-obvious field has a `provenance` entry.
- [ ] Every intentional omission is recorded in `dropped[]`.
- [ ] `scripts/validate_extraction.py` exits 0 on the emitted envelope.

## References

- Preserve-vs-drop rules and size discipline: `references/loss-bounding.md`
- Conform-vs-infer decision and lean-contract heuristics: `references/schema-flexibility.md`
- Extraction mistakes to avoid: `references/anti-patterns.md`
- Envelope skeleton: `templates/extraction.example.json`
- Worked end-to-end run (both modes): `examples/doc-to-contract.example.md`
- Deterministic envelope check: `scripts/validate_extraction.py`
- I/O contracts: `schemas/input.json`, `schemas/output.json`
