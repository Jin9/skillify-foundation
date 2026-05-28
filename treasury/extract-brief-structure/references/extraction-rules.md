# Extraction Rules

## Source Detection

- `jira`: Jira key pattern, project/type/priority headers, description, acceptance criteria, comments, linked issues.
- `slack`: channel banner, timestamped speakers, reactions, linked attachments, chronological replies.
- `meeting-notes`: meeting/date/note-taker/attendees metadata, agenda sections, action items, apologies.
- `email`: From/To/Subject/Date quartet, quoted reply prefixes, signatures.
- `mixed`: no durable single-source markers or multiple merged sources.

Use the user's `source_type` only as a hint. If markers contradict the hint, report both the hint and detected type in `extraction_notes`.

## Preprocessing

Strip ground-truth or audit annotation blocks before extraction. Treat headings like "Intentional Issues", "Hidden from BA Workflow", "Ground Truth", or "Audit Annotation" as strip boundaries. If more than one boundary or partial boundary is detected, block the stage rather than using contaminated input.

## Scope Classification

- `single-story`: one bounded user-value slice.
- `multi-story`: multiple related slices under one epic.
- `single-epic`: one workstream with enough signal for an epic and child stories.
- `multi-epic`: at least three distinct workstreams, each with at least two acceptance criteria, business rules, or decision signals.
- `story_within_epic`: input clearly asks for one story inside a named existing epic.
- `ambiguous`: source asks whether to split, or evidence supports multiple valid scope shapes.

## Story Splitting

Allowed split axes:

- workflow step or state transition
- business-rule variation
- data class variation, such as PII vs non-PII or biometric vs document
- role boundary, such as customer UI vs agent UI
- happy path vs error path
- CRUD operation
- spike or research task

Forbidden split axis: UI/API/database layer alone. If the source forces a technical split, capture it as `extraction_notes[]` and recommend TL review.

## Required Evidence

Each epic and story must include at least one `source_refs[]` entry pointing to a section, speaker, line, or source fragment. Do not invent stakeholder names, dates, regulators, or numeric values absent from the input.
