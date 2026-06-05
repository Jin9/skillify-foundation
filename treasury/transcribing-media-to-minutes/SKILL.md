---
name: transcribing-media-to-minutes
description: >
  Read a meeting recording (audio or video — Thai/English bilingual, usually
  mostly Thai with English banking and tech terms mixed in) and extract it as
  faithful Minutes of Meeting (MOM): a Thai-primary, plain-language minutes file
  plus a timestamped, speaker-labeled transcript sidecar. Use when the user asks
  to "summarize this meeting recording", "turn this audio/video into minutes",
  "extract MOM from this recording", "transcribe and minute this call", or
  "ทำรายงานการประชุมจากไฟล์เสียง/วิดีโอ". Routes by recording length and
  complexity to a fast model variant for transcription and a frontier variant
  for synthesis (Gemini Flash and Pro shown as dated examples; portable across
  hosts, with an external-transcript fallback when the host model cannot ingest
  media). Do NOT use for live real-time transcription, for non-meeting media, or
  for translating a document that has no audio or video.
compatibility: gemini, antigravity, claude-code, codex
---

# Transcribing Media to Minutes

## Purpose

Turn a recorded meeting (audio or video) into faithful **Minutes of Meeting (MOM)** without losing
the original context. Recordings are usually **mostly Thai** with English banking and tech terms
code-switched in. The skill produces a Thai-primary, plain-language minutes file plus a verbatim
timestamped transcript, with every decision and action item traceable back to a moment in the
recording. It does not do live/streaming transcription, and it does not translate documents that
have no media.

## When to use this skill

- Use when: the user gives an audio or video meeting recording and asks to "minute this", "summarize
  this recording", "turn this call into minutes", or "ทำรายงานการประชุม".
- Use when: the user wants both a clean MOM and a speaker-labeled transcript from a recording.
- Do NOT use when: the user wants live/real-time transcription of an ongoing call — this skill works
  on a finished recording or file.
- Do NOT use when: the media is not a meeting (a lecture, a song, a screen-capture demo) and no
  minutes are wanted.
- Do NOT use when: the user only wants a document translated and there is no audio or video — use a
  translation skill instead.

## Input

- Required: a path or URL to one audio or video file (e.g. `.mp3`, `.m4a`, `.wav`, `.mp4`, `.mov`).
- Optional context: attendee names and roles, the agenda, a project glossary, the target output
  directory, and a language override (default is Thai-primary bilingual; see step 5).
- Default output location: the same directory as the source media. Use a user-supplied directory when
  one is given.

## Workflow

1. **Intake and privacy gate.** Confirm the media path, format, and duration. Classify the
   sensitivity of the content *before* sending it anywhere: meeting recordings often contain PII or
   bank-confidential material (PDPA B.E. 2562). Match the data class to an eligible model class using
   `references/model-variant-routing.md`. If no eligible model satisfies the data-residency rule, stop
   and ask the user how to proceed — never upload sensitive audio to an ineligible model.
2. **Capability check and model-variant routing.** Determine whether the host model can natively
   ingest audio/video (a Gemini variant can; some hosts cannot). If it cannot, switch to the
   external-transcript fallback: ask for or generate a transcript first, then resume at step 4. Route
   the work by recording length, speaker count, and synthesis depth — a fast variant for
   straightforward transcription, a frontier variant for long or multi-speaker meetings that need
   heavier reasoning — and choose inline vs File-API delivery by size/duration. See
   `references/model-variant-routing.md`.
3. **Transcribe and diarize.** Produce a timestamped, speaker-labeled transcript following
   `templates/transcript-template.md`. Preserve Thai/English code-switching **verbatim** — keep
   English terms exactly as spoken; do not translate the transcript. Mark inaudible or low-confidence
   spans explicitly rather than guessing.
4. **Extract the MOM structure faithfully.** From the transcript only, derive attendees, agenda /
   topics, a discussion summary, decisions, action items (owner + due date), open questions, and
   risks / next steps. Bind every decision and action item to a transcript timestamp. Never invent an
   owner, a date, a decision, or an attendee that the recording does not support; mark anything
   uncertain. Field list and rules: `references/mom-extraction-and-faithfulness.md`.
5. **Render Thai-primary, plain-language minutes.** Write the MOM mostly in Thai, keeping English
   banking/tech terms where they are normally used (do not over-translate), in simple, plain
   language. Add a one-line plain explanation for any jargon or complex decision so a non-expert can
   follow it — the explanation augments and must not change the original meaning. Fill
   `templates/minutes-of-meeting-template.md`. Rendering rules and a term glossary:
   `references/bilingual-thai-minutes.md`. Honor a language override if the user gave one.
6. **Write outputs and verify.** Write the two files (see Output contract). Run the faithfulness
   check: every decision and action traces to a timestamp; nothing is invented; the original context
   is preserved; the plain-language pass is done. Surface any low-confidence or inaudible spans at the
   top of the MOM so the user can review them.

## Output contract

Two files, written beside the source media (or in a user-supplied directory):

- `<media-basename>-mom.md` — the Minutes of Meeting: Thai-primary bilingual, plain-language, built
  from `templates/minutes-of-meeting-template.md`. Contains a metadata block (title, date, duration,
  attendees), agenda/topics, a discussion summary, a decisions list, an action-item table (owner +
  due + timestamp), open questions, and risks / next steps.
- `<media-basename>-transcript.md` — the verbatim, timestamped, speaker-labeled transcript, built
  from `templates/transcript-template.md`. Thai/English kept exactly as spoken.

Naming rule: `<media-basename>` is the source file name without its extension. Never overwrite an
existing MOM for the same recording without confirming.

## Faithfulness and plain-language rules

- **Keep the same context.** The MOM must mean what the meeting meant. Do not add conclusions, soften,
  or sharpen what was said. Plain-language phrasing simplifies *wording*, never *meaning*.
- **Evidence-bound.** Every decision and action item carries an `[hh:mm:ss]` timestamp into the
  transcript. If it has no timestamp, it does not belong in those sections.
- **No invention.** Unknown owner, date, or attendee → write "ไม่ระบุ (not stated)", not a guess.
- **Bilingual fidelity.** Keep code-switched English terms verbatim; gloss them once in plain Thai.
- **Flag uncertainty.** Inaudible or low-confidence spans are marked, never silently filled.

Deeper guidance lives in the references — do not duplicate it here.

## Constraints and anti-patterns

- Never upload a recording to a model whose data class fails the privacy gate in step 1.
- Never invent decisions, owners, due dates, or attendees; mark uncertainty instead.
- Never translate the transcript — it is verbatim; translation/simplification happens only in the MOM.
- Never hardcode a single vendor as required: Gemini variants are dated examples; other hosts use the
  external-transcript fallback. Keep this body portable.
- Not for live/real-time transcription, non-meeting media, or document-only translation.

## References

| Need | File |
|------|------|
| Model-variant routing, inline vs File API, privacy/data-class gate, non-multimodal fallback | `references/model-variant-routing.md` |
| Thai-primary bilingual rendering, plain-language rules, term glossary | `references/bilingual-thai-minutes.md` |
| MOM field list, evidence-binding, anti-hallucination, context-preservation check | `references/mom-extraction-and-faithfulness.md` |

## Templates and examples

- `templates/minutes-of-meeting-template.md` — the MOM output skeleton.
- `templates/transcript-template.md` — the timestamped, speaker-labeled transcript skeleton.
- `examples/sample-mom.md` — a short before→after (a Thai/English standup recording → filled MOM plus a transcript excerpt).
