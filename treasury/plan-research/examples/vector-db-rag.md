# Worked example: vector-db-rag

Example invocation of the `plan-research` skill on a mid-complexity technical topic at `standard` depth for a `technical` audience. Demonstrates: a question-shaped thesis, 5 MECE pillars + sub-questions, a non-empty `out_of_scope` list, 5 testable `success_criteria`, and audience-appropriate jargon (HNSW, IVF-PQ, pgvector) preserved without inline expansion.

## Input

```yaml
topic: "Vector databases for retrieval-augmented generation"
depth: standard
audience: technical
```

## Output (`research_plan`)

```json
{
  "topic": "Vector databases for retrieval-augmented generation",
  "depth": "standard",
  "audience": "technical",
  "thesis_question": "Which vector databases are best suited to retrieval-augmented generation workloads, and on what dimensions do they meaningfully differ?",
  "coverage_dimensions": [
    "landscape_or_options",
    "mechanisms_or_how_it_works",
    "evidence_and_outcomes",
    "trade_offs_or_constraints",
    "counterpoints_or_controversies"
  ],
  "sub_questions": [
    {
      "id": "q1",
      "question": "What are the leading vector database options used in RAG systems as of 2026, and how do they segment by deployment model (managed, self-hosted, embedded)?",
      "rationale": "Anchors the report in a current landscape so later trade-off discussion has named referents.",
      "coverage_dimension": "landscape_or_options",
      "expected_source_types": ["industry", "vendor", "community"]
    },
    {
      "id": "q2",
      "question": "How do mainstream vector indexes (HNSW, IVF-PQ, DiskANN) trade off recall, latency, and memory for RAG-typical workloads (~1M–100M chunks)?",
      "rationale": "RAG performance is dominated by index choice; engineers need the mechanism-level picture, not vendor marketing.",
      "coverage_dimension": "mechanisms_or_how_it_works",
      "expected_source_types": ["academic", "primary", "docs"]
    },
    {
      "id": "q3",
      "question": "What published benchmarks measure vector-DB performance on RAG-relevant tasks, and what do they show for top systems?",
      "rationale": "Grounds claims in evidence rather than vendor copy.",
      "coverage_dimension": "evidence_and_outcomes",
      "expected_source_types": ["academic", "industry"]
    },
    {
      "id": "q4",
      "question": "What are the operational trade-offs of running a vector database in production for RAG (cost per million vectors, scaling story, hybrid search support, ecosystem)?",
      "rationale": "Production choice rarely follows from latency alone; trade-offs determine the real decision.",
      "coverage_dimension": "trade_offs_or_constraints",
      "expected_source_types": ["industry", "docs", "vendor"]
    },
    {
      "id": "q5",
      "question": "Where do practitioners argue that a vector database is the wrong primitive for RAG (e.g., Postgres pgvector, BM25-only, in-memory FAISS), and on what grounds?",
      "rationale": "A serious technical report must surface the credible 'don't use one' position.",
      "coverage_dimension": "counterpoints_or_controversies",
      "expected_source_types": ["community", "industry", "primary"]
    }
  ],
  "out_of_scope": [
    "Training or fine-tuning embedding models from scratch.",
    "Non-RAG uses of vector search (recommendation, image retrieval, anomaly detection).",
    "LLM hosting / inference infrastructure.",
    "Pre-2023 vector-DB landscape (the field has materially shifted)."
  ],
  "success_criteria": [
    "Each of the five coverage_dimensions is addressed by at least one cited source in the final report.",
    "The landscape section names ≥ 5 specific vector databases with their deployment model.",
    "The mechanisms section explains HNSW, IVF-PQ, and at least one disk-based index at a level a senior engineer can act on.",
    "The trade-offs section cites at least one concrete cost or latency number with its source and date.",
    "The counterpoints section includes at least one credible 'use Postgres / BM25 instead' argument with a named source."
  ],
  "notes_for_downstream": "Field changed materially in 2024–2026; prefer sources <18 months old for landscape and benchmarks. pgvector is a recurring controversial pick — search for both pro and anti arguments. Watch for vendor-funded benchmarks; cross-check with at least one independent source per claim.",
  "plan_skipped": false
}
```

This plan: 5 sub-questions (standard depth), each MECE-mapped to a distinct pillar, audience-appropriate jargon (HNSW, IVF-PQ, pgvector), explicit out-of-scope exclusions, and concrete success criteria the `review-report` stage can check against.
