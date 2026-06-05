# MOM extraction and faithfulness

Deep guidance for step 4 (and the verify part of step 6) of `SKILL.md`. The source of truth is the
**transcript**. If a claim is not supported by the transcript, it does not go in the minutes.

## Fields to extract

| Section | What goes in it | Rule |
|---------|-----------------|------|
| Metadata | title, date, duration, attendees (name + role) | attendees only if named or clearly identified; else "ไม่ระบุ (not stated)" |
| Agenda / topics | the topics actually discussed | derive from content if no agenda was stated |
| Discussion summary | the key points per topic, plain language | summarize, do not transcribe; keep meaning |
| Decisions | what was agreed | each with an `[hh:mm:ss]` timestamp |
| Action items | task, owner, due date | each with owner + timestamp; due date only if stated |
| Open questions | unresolved points, things "to confirm" | note who must answer, if stated |
| Risks / next steps | risks raised, agreed next steps | a next step not yet owned goes here, not in Action items |

## Evidence-binding

- Every **decision** and **action item** must carry an `[hh:mm:ss]` timestamp pointing into the
  transcript. No timestamp → it is not a decision/action; move it to the discussion summary or drop it.
- When several moments support one item, cite the moment the decision was actually made.

## Anti-hallucination rules

1. Never invent an owner, due date, decision, or attendee. Unknown → "ไม่ระบุ (not stated)".
2. Do not infer agreement from silence. "No objection recorded" is not "approved".
3. Do not merge two different people's points into one false consensus.
4. Do not promote a casual "maybe we should…" into a Decision; it is at most an Open question.
5. Distinguish proposal vs decision vs action. A proposal that was not agreed is not a decision.

## Context-preservation / information-loss check (run in step 6)

Summarizing loses detail on purpose — bound the loss, do not let it distort:

- [ ] Every decision in the recording appears in Decisions (none dropped).
- [ ] Every assigned task appears in Action items with an owner.
- [ ] No item appears that the transcript does not support (nothing added).
- [ ] Numbers, dates, amounts, and names are copied exactly, not approximated.
- [ ] Each item still means what it meant in context (no softening/sharpening).
- [ ] Low-confidence or inaudible spans are listed at the top of the MOM for review.

If any box fails, fix the MOM before delivering. Report residual low-confidence items rather than
hiding them.
