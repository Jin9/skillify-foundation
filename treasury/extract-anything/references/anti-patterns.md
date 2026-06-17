# Anti-patterns

Extraction mistakes that defeat the purpose of a chainable contract. Each has a symptom and the fix.

## 1. Lossy prose summary masquerading as a contract
- **Symptom:** `contract` is a paragraph, or one big `summary` string, instead of structured fields.
- **Why it's bad:** downstream stages (often small/mid-tier models) cannot key on prose; edge cases dissolve into fluent text.
- **Fix:** structure the must-preserve content into typed fields. The summary, if any, is one field among many — never the whole contract.

## 2. Silent dropping
- **Symptom:** content was left out but `_meta.dropped` does not mention it.
- **Why it's bad:** the consumer believes the contract is complete; the loss is invisible and unrecoverable as a known gap.
- **Fix:** every intentional omission gets a `{ item, reason }` entry. If you cannot describe what you dropped, you do not understand the source well enough to extract it.

## 3. Hallucinated fields (conform mode)
- **Symptom:** a required target field is filled with a plausible-but-unsupported value.
- **Why it's bad:** corrupts every downstream stage silently; the worst failure mode of the skill.
- **Fix:** leave the field absent and record `reason: not-in-source`. An honest gap beats a confident lie.

## 4. Unbounded / inlined blobs
- **Symptom:** a contract field holds a multi-kilobyte verbatim span.
- **Why it's bad:** blows the consumer's context budget; the contract stops being lean.
- **Fix:** store a short form plus a `source_ref` anchor; record the trim as `reason: size-bounded`.

## 5. Deep nesting
- **Symptom:** the contract is four-plus levels deep, mirroring the source's structure.
- **Why it's bad:** unpredictable and expensive for downstream stages to traverse; breaks context economy.
- **Fix:** flatten. Mirror the consumer's needs, not the source's layout.

## 6. Re-reading the source downstream
- **Symptom:** the contract is so thin that the next stage must open the source to act.
- **Why it's bad:** defeats the entire point of a handoff contract; couples stages to the source.
- **Fix:** raise `coverage`. Ensure every must-preserve unit is in `contract`, not deferred to `source_ref`. `source_ref` is for recovering *dropped* context, not the *primary* content.

## 7. Schema invention (infer mode drifting into conform)
- **Symptom:** you invented a strict schema and validated against it, even though the caller gave none.
- **Why it's bad:** imposes a shape the consumer never agreed to; brittle.
- **Fix:** in infer mode keep the shape minimal and flat, and recommend the caller pass a `target_schema` next time (in `_meta.notes`).

## 8. Multi-source concatenation
- **Symptom:** several sources merged into one envelope.
- **Why it's bad:** muddies provenance and coverage; the loss ledger no longer maps to one origin.
- **Fix:** one source → one envelope. Extract each separately; let the orchestrator combine contracts.

## 9. Editorializing
- **Symptom:** the contract contains the extractor's opinions, recommendations, or judgments about whether the source's claims are true.
- **Why it's bad:** extraction is faithful transcription of *what the source says*, not adjudication of *whether it is right* (that is a different skill's job).
- **Fix:** carry the source's content; put any extraction caveat in `_meta.notes`, never disguised as source content.
