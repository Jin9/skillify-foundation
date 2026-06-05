# Model-variant routing for media → minutes

Deep guidance for step 1 (privacy gate) and step 2 (capability check + routing) of `SKILL.md`.
The skill body stays vendor-neutral; this file names Gemini variants only as **dated examples
(as of 2026-06)** — re-verify against the host's current docs before relying on a number.

## 1. Privacy / data-class gate (run BEFORE any upload)

Classify the recording's most sensitive content, then eliminate any model class that fails the rule
for that class. Choose the model only from what survives.

| Data class | Examples | Eligible model class |
|------------|----------|----------------------|
| Public / non-sensitive | public webinar, town hall | any hosted model |
| Internal | team standup, sprint planning | hosted-standard OK |
| Confidential / PII | customer names, account data | ZDR + BAA, or single-tenant VPC |
| Restricted (PDPA B.E. 2562 / bank-confidential / no-egress) | KYC calls, credit decisions | open-weight self-host, or no upload |

If nothing survives the gate, **stop and ask the user** — do not downgrade the requirement silently.
Krungthai DGL / banking recordings default to Confidential or Restricted unless told otherwise.

## 2. Native-multimodal capability check

- A multimodal host (e.g. a Gemini variant) can take the audio/video file directly — no separate ASR.
- A text-only host cannot. Switch to the **external-transcript fallback**: obtain a transcript from a
  dedicated ASR tool (or ask the user to supply one), then resume the skill at step 4 using that
  transcript as the source of truth.
- Detect, do not assume: infer from the host/model in use, or ask once if unclear.

## 3. Variant routing — match the recording to a capability tier

Two named tiers, vendor-neutral. Pick per stage; a long meeting can use the fast tier for raw
transcription and the frontier tier only for synthesis.

| Tier | When | Dated example (2026-06) |
|------|------|--------------------------|
| Fast / low-cost | short (under ~30 min), one/few speakers, clear audio, straightforward ASR | Gemini 2.x Flash |
| Frontier / high-reasoning | long (~30 min+), many speakers, heavy code-switching, noisy audio, or dense decisions needing careful synthesis | Gemini 2.x Pro |

Routing signals: duration, speaker count, audio quality, density of decisions, and how much
Thai/English code-switching must be disambiguated. When in doubt for a banking meeting, prefer the
frontier tier for the synthesis stage — faithfulness matters more than cost here.

## 4. Delivery: inline vs File API

Dated examples (2026-06) — verify against current host limits:

- **Inline** for small/short clips (roughly under ~20 MB / a few minutes), sent directly in the request.
- **File API upload** for larger or longer files; reference the uploaded handle in the request.
- Very long audio (around an hour or more) may need chunking. Keep chunk boundaries on natural pauses
  and carry a 5–10 s overlap so no sentence is split across chunks; renumber timestamps continuously.

## 5. Re-evaluation

Model names, tiers, and limits move fast. Treat every model name here as a **dated example**, keep the
decision behind the fast/frontier abstraction, and re-confirm the mapping and limits roughly quarterly
or whenever the host's model lineup changes.
