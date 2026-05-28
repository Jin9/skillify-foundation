# Audience profiles

`synthesize-report` accepts a free-form `audience` string and maps it to one
of four named profiles. Audience does NOT change which findings get selected
(that's settled in Procedure Step 3); it changes how those findings are
framed when drafted.

| Profile | Tone | Vocabulary | Depth bias | Examples / framing |
|---|---|---|---|---|
| `general` | conversational, accessible | plain English; expand acronyms first use; minimize jargon | breadth over depth; one good example per concept | familiar analogies; "what this means for you" framings |
| `executive` | declarative, "so what" first, decisive | business / strategic vocab; shared industry jargon OK | shallow — recommendations up top, evidence in tail; cut mechanism detail | dollar impact, risk framing, market comparators, action implications |
| `expert` | precise, neutral, cite-or-it-didn't-happen | full domain terminology assumed | depth over breadth; nuance and edge cases included | mechanism, parameter values, comparative method, references to prior work |
| `academic` | formal, hedged ("the evidence suggests"); passive voice OK | full technical + methodology vocab | full depth; methodology section preserved | citations to peer-reviewed work; replication context; epistemic caveats |

## Mapping free-form strings

When `audience` is a free-form string, pick the closest profile:

- "general", "public", "lay reader", "non-specialist" → `general`
- "executive", "CXO", "decision maker", "leadership", "investor" → `executive`
- "engineer", "developer", "practitioner", "specialist", "technical" → `expert`
- "researcher", "scientist", "academic", "PhD", "peer review" → `academic`

When the string is unrecognizable, default to `general` but lift the bar to
"informed adult reader" rather than absolute novice.

## Validation hook

Validation-gate item 9 (`Audience tone check`) enforces:

- No `general` report contains undefined jargon in body prose.
- No `executive` report defers the recommendation past the executive summary.
- No `expert` report omits mechanism / parameter detail when the findings
  supply it.
