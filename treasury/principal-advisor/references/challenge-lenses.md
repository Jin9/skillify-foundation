# Challenge lenses — steelman-first assumption testing

How to run the `challenge` stance without turning hostile or, worse, useless.

## Steelman first — always

Before hunting flaws, state the strongest honest case for the user's idea in 2–4 sentences. This is not politeness: critics that skip the steelman reliably reject good ideas and pattern-match on superficial flaws. A finding only counts if it still stands after the idea's best case has been argued.

## Attack lenses

Sweep these; raise a finding only where there is evidence in what the user said or showed — never on style or taste.

- **Goal fit** — does the idea actually serve the stated goal, or a more enjoyable adjacent goal?
- **Complexity** — is the machinery bigger than the problem? What is the simplest version that still works?
- **Data & state** — who owns which data, what consistency does the idea silently assume, what happens at migration time?
- **Failure & operations** — what breaks at 3 a.m., who notices, and how is it rolled back?
- **Cost** — build cost, run cost, and the blast radius when the optimistic assumption fails.
- **Security & compliance exposure** — where are the auth boundaries, secrets, and personal data? Flag for a real review; do not deep-model it in chat.
- **Team & capability** — can this team, at its real size and skill mix, build it and then operate it?
- **Reversibility** — which parts are one-way doors, and is the user walking through one casually?

## Etiquette of findings

- Severity-rank; raise at most three per turn — the strongest three, not the first three.
- Attack the idea, never the person.
- End every finding with its dissolve condition: the evidence or answer that would make it go away. A challenge without a dissolve condition is just an opinion.

## Bias check before sending

Re-read each finding and cut or convert it if:

- **Hallucinated** — it references something the user never said or showed. Delete it.
- **Position or verbosity bias** — flagged because it appeared first or its section was long. Re-judge on substance.
- **Unfalsifiable** — nothing available could confirm or deny it. Convert it to an open question, not a finding.
- **Taste** — a style preference wearing a risk costume. Drop it.

## Exit protocol

Track each finding to accepted, rebutted, or parked. Then restate the surviving position in 1–2 sentences — the user should leave knowing what their idea looks like after the fire, not just that there was fire. If the user says stop, stop.
