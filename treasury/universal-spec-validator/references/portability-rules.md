# Cross-model semantic-portability drift rules

Distilled from research reports *Cross-model schema normalization & semantic
portability* and *Cross-model skill portability*. Portability is a layered
control problem: drift enters independently at the transport, generation,
alignment, coercion, evolution, and artifact layers. Syntactic portability is
largely solved — 35+ tools share the SKILL.md format and constrained decoding
closes the schema-shape gap; residual risk is **semantic and behavioral** and
concentrates at the cross-vendor boundary. Adherence there is sharply
model-dependent: gpt-4o with Structured Outputs scores ~100% on complex
JSON-schema following while an older model in the same family scores under 40%,
and a schema authored for one vendor's strict mode may be silently altered or
rejected by another's.

The gate checks the structurally-detectable subset below. It cannot prove
semantic equivalence (no benchmark does) — it flags the conditions known to
*cause* drift so they are reviewed before merge.

## P1 — Single-vendor tool-call dialect lock-in (severity: high)
Tool/function specs differ across families: OpenAI wraps as
`{"type":"function","function":{...}}` with `parameters`; Anthropic uses a
flat object with `input_schema`; Gemini uses protobuf-style `types.Schema`.
A spec hard-coded to one dialect fails silently on another and passes invalid
data downstream. **Flag** a spec that uses exactly one vendor shape with no
declared normalization layer / gateway and no `x-portable`/adapter marker.

## P2 — Vendor strict-mode incompatibility (severity: medium)
OpenAI strict mode requires `additionalProperties:false`, every property in
`required`, and optionals expressed as a `null` union. A schema authored for
one strict mode is rejected or silently altered by another. **Flag**:
`additionalProperties` absent/true on an object used as a tool input; optional
fields not modeled as nullable; `required` missing properties when strict.

## P3 — Reasoning-killed-by-format (severity: low/info)
Forcing structure too early degrades reasoning; the mitigation is a
reason-then-format split. **Info** when a tool/skill spec mandates structured
output for a step described as requiring multi-step reasoning with no separate
reasoning channel.

## P4 — Non-portable artifact dependency (severity: medium)
Prompts and embeddings do not transfer reliably across (often even within) model
families; embeddings are translatable but that is also a confidentiality risk.
**Flag** a spec/skill that embeds a model-specific prompt format, a
provider-specific embedding reference, or assumes a context-window size, without
a `requires_capabilities` / `min_context_tokens` annotation enabling
capability-gated routing.

## P5 — Missing capability/safety annotation (severity: medium)
Refusal thresholds and capabilities (code_execution, web_access) differ across
providers; unannotated skills fail at runtime on a stricter model. **Flag** a
skill/tool that implies code execution, file write, or network access but
declares no capability requirement the router can gate on.

## P6 — Unidentified shared keys (severity: info)
JSON-Schema/Protobuf carry no semantic typing; reused keys across record
classes drift in meaning. **Info**: recommend JSON-LD `@context` or an IRI/URI
`$id` for cross-agent shared records.

Gate rule: P1/P4/P5 default `high|medium` → blocking per the severity policy
unless the config downgrades; P2 medium; P3/P6 info (never blocking). Every
emitted finding cites `portability-rules.md#P<n>`.
