# Worked example: four-day-week

Example invocation of the `extract-findings` skill on a small input (4 sources, 3 sub-questions). Demonstrates: evidence-first extraction, corroboration via `supporting_source_ids`, contradiction surfacing via `disputed_by` (separate findings emitted, not collapsed), the qualitative confidence rubric in use, and `partial` coverage with a `gap_note` for the under-covered sub-question.

## Tiny input

```yaml
research_plan:
  topic: "Are four-day work weeks viable for knowledge work?"
  sub_questions:
    - id: sq-1
      question: "What productivity effects have trials reported?"
    - id: sq-2
      question: "What employee well-being effects have trials reported?"
    - id: sq-3
      question: "Which industries have shown poor fit?"

sources:
  - id: s-1
    title: "Iceland 4-day-week trial results (Autonomy/Alda 2021)"
    url: "https://example.org/iceland-2021"
    text: "Trials covering 2,500 workers (1% of working population)
           ran 2015-2019. Productivity 'remained the same or improved
           in the majority of workplaces.' Workers reported reduced
           stress and burnout."
  - id: s-2
    title: "UK pilot summary (4 Day Week Foundation 2022)"
    url: "https://example.org/uk-2022"
    text: "61 companies, 2,900 employees. 92% continued with 4-day
           week after the trial. Average revenue change +1.4% vs.
           same period prior year. Burnout dropped 71%."
  - id: s-3
    title: "Manufacturing op-ed (Industry Week 2023)"
    url: "https://example.org/iw-2023"
    text: "Author argues continuous-process plants cannot compress
           into four days without idle capacity or extra shifts.
           Acknowledges office support functions may differ."
  - id: s-4
    title: "Meta-analysis of compressed schedules (J. Org. Behav. 2019)"
    url: "https://example.org/jobeh-2019"
    text: "Effects on objective productivity were 'mixed and
           context-dependent.' Self-reported productivity rose in 7
           of 11 included studies; objective output rose in 3 of 8."
```

## Expected output

```json
{
  "findings": [
    {
      "id": "f-1",
      "sub_question_id": "sq-1",
      "claim": "Iceland's 2015-2019 trials reported productivity remaining the same or improving in the majority of workplaces.",
      "evidence": "\"Productivity remained the same or improved in the majority of workplaces.\" Trials covered 2,500 workers (1% of working population).",
      "source_id": "s-1",
      "supporting_source_ids": [],
      "disputed_by": ["s-4"],
      "confidence": "medium",
      "note": "Self-reported productivity; meta-analysis (s-4) finds objective measures are more mixed."
    },
    {
      "id": "f-2",
      "sub_question_id": "sq-1",
      "claim": "The UK 2022 pilot reported average revenue +1.4% vs. the same period prior year across 61 companies.",
      "evidence": "\"Average revenue change +1.4% vs. same period prior year.\" 61 companies, 2,900 employees.",
      "source_id": "s-2",
      "supporting_source_ids": [],
      "disputed_by": [],
      "confidence": "high",
      "note": null
    },
    {
      "id": "f-3",
      "sub_question_id": "sq-1",
      "claim": "A 2019 meta-analysis found objective productivity effects of compressed schedules were mixed and context-dependent.",
      "evidence": "Effects on objective productivity were \"mixed and context-dependent.\" Self-reported productivity rose in 7 of 11 included studies; objective output rose in 3 of 8.",
      "source_id": "s-4",
      "supporting_source_ids": [],
      "disputed_by": ["s-1"],
      "confidence": "high",
      "note": "Counterpoint to single-trial positive reports."
    },
    {
      "id": "f-4",
      "sub_question_id": "sq-2",
      "claim": "Iceland trial participants reported reduced stress and burnout.",
      "evidence": "\"Workers reported reduced stress and burnout.\"",
      "source_id": "s-1",
      "supporting_source_ids": ["s-2"],
      "disputed_by": [],
      "confidence": "high",
      "note": null
    },
    {
      "id": "f-5",
      "sub_question_id": "sq-2",
      "claim": "The UK 2022 pilot reported a 71% drop in burnout among participating employees.",
      "evidence": "\"Burnout dropped 71%.\"",
      "source_id": "s-2",
      "supporting_source_ids": [],
      "disputed_by": [],
      "confidence": "high",
      "note": null
    },
    {
      "id": "f-6",
      "sub_question_id": "sq-3",
      "claim": "Continuous-process manufacturing may not compress into four days without idle capacity or extra shifts.",
      "evidence": "\"Continuous-process plants cannot compress into four days without idle capacity or extra shifts. Acknowledges office support functions may differ.\"",
      "source_id": "s-3",
      "supporting_source_ids": [],
      "disputed_by": [],
      "confidence": "low",
      "note": "Single op-ed; author acknowledges office functions may differ. No empirical study cited."
    }
  ],
  "coverage": [
    { "sub_question_id": "sq-1", "finding_ids": ["f-1", "f-2", "f-3"], "status": "covered", "gap_note": null },
    { "sub_question_id": "sq-2", "finding_ids": ["f-4", "f-5"], "status": "covered", "gap_note": null },
    { "sub_question_id": "sq-3", "finding_ids": ["f-6"], "status": "partial", "gap_note": "Only one op-ed; no empirical sector-level data on poor-fit industries." }
  ],
  "unsupported_sources": []
}
```

## Notes on the example

- f-1 and f-3 are emitted as **separate findings** with mutual `disputed_by` pointers — the trial reports positive productivity and the meta-analysis finds mixed effects. Collapsing them would hide the contested-question signal.
- f-4 lists `supporting_source_ids: ["s-2"]` — both Iceland and UK pilots reported reduced burnout independently. Corroboration upgrades this finding to `high` confidence even though each individual quote is short.
- f-6 is `low` confidence: single op-ed, no empirical study, and the author's own caveat. Coverage for sq-3 is `partial` and the `gap_note` tells `synthesize-report` not to overstate this conclusion.
- `unsupported_sources` is empty — every input source produced at least one finding.
