# Worked example: heat-pump-retrofit

Example invocation of the `search-sources` skill on a policy-flavoured technical topic at `standard` depth. Demonstrates: query generation (2 per sub-question on standard), the source-record schema, and how `coverage[]` is used to surface thin sub-questions rather than padding the source list with low-quality results.

## Input (abbreviated)

```yaml
research_plan:
  topic: "Are heat pumps cost-effective for retrofit in cold climates?"
  sub_questions:
    - id: q1
      text: "What is the COP of modern cold-climate heat pumps at -15 C?"
    - id: q2
      text: "What does a residential retrofit cost in the US northeast in 2025?"
    - id: q3
      text: "What incentives or rebates are available in MA/NY/VT in 2025?"
depth: standard
```

## Queries the stage would issue

Standard depth ⇒ 2 queries per sub-question:

- q1: `cold climate heat pump COP at -15C 2025` ;
  `coefficient of performance residential heat pump subzero performance`
- q2: `residential heat pump retrofit cost northeast US 2025` ;
  `total installed price air-source heat pump replacing furnace homeowner`
- q3: `Massachusetts New York Vermont heat pump rebate 2025` ;
  `state incentives residential electrification cold climate 2025`

## Output (abbreviated, illustrative — concrete URLs would come from the configured backend)

```json
{
  "sources": [
    {
      "source_id": "s1",
      "url": "https://www.nrel.gov/.../cchp-field-performance-2025.pdf",
      "title": "Cold-Climate Heat Pump Field Performance, 2025 Update",
      "snippet": "Field-measured COP for cold-climate heat pumps averaged 2.1 at -15 C across 84 sites in the Northeast US.",
      "retrieved_at": "2026-05-14T14:02:11Z",
      "query": "cold climate heat pump COP at -15C 2025",
      "sub_question_ids": ["q1"],
      "source_type": "paper",
      "published_date": "2025-09-01",
      "publisher": "NREL",
      "relevance_score": 0.92
    },
    {
      "source_id": "s2",
      "url": "https://www.energystar.gov/.../cold-climate-air-source-heat-pumps",
      "title": "Cold Climate Air-Source Heat Pumps",
      "snippet": "ENERGY STAR Cold Climate criteria require a minimum HSPF of 10 and rated capacity retention at 5 F.",
      "retrieved_at": "2026-05-14T14:02:13Z",
      "query": "coefficient of performance residential heat pump subzero performance",
      "sub_question_ids": ["q1"],
      "source_type": "docs",
      "publisher": "ENERGY STAR",
      "relevance_score": 0.81
    },
    {
      "source_id": "s3",
      "url": "https://emp.lbl.gov/.../residential-hp-retrofit-cost-2025",
      "title": "Residential Heat Pump Retrofit Costs, 2025",
      "snippet": "Median installed cost for a whole-home cold-climate heat pump retrofit in the US Northeast was $18,400 in 2025, before incentives.",
      "retrieved_at": "2026-05-14T14:02:17Z",
      "query": "residential heat pump retrofit cost northeast US 2025",
      "sub_question_ids": ["q2"],
      "source_type": "paper",
      "published_date": "2025-06-15",
      "publisher": "LBNL",
      "relevance_score": 0.89
    },
    {
      "source_id": "s4",
      "url": "https://www.mass.gov/info-details/masssave-heat-pump-rebates",
      "title": "Mass Save Heat Pump Rebates",
      "snippet": "Mass Save offers up to $10,000 in rebates for whole-home cold-climate heat pump installations in 2025.",
      "retrieved_at": "2026-05-14T14:02:21Z",
      "query": "Massachusetts New York Vermont heat pump rebate 2025",
      "sub_question_ids": ["q3"],
      "source_type": "docs",
      "publisher": "Mass.gov",
      "relevance_score": 0.86
    }
  ],
  "coverage": [
    {"sub_question_id": "q1", "source_count": 2, "gap_note": null},
    {"sub_question_id": "q2", "source_count": 1, "gap_note": "below standard floor (target 3+); consider a deep re-run"},
    {"sub_question_id": "q3", "source_count": 1, "gap_note": "only Massachusetts represented; NY and VT not yet covered"}
  ]
}
```

The example shows: real backend-supplied URLs only, one record per surviving result, every source tagged with its sub-question, two `coverage` entries flagging gaps so `review-report` can challenge thin sub-questions, and a diversity mix (paper / docs across two government and two research sites).
