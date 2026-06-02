# Plain-language Thai rewriting

Depth guidance for Step 3 of `SKILL.md`. Goal: rewrite each section's *wording* into Thai a general reader can follow, without changing the source's meaning, structure, or coverage.

## Target reader

A general Thai adult with no background in the source's domain. Assume they are smart but unfamiliar with the field's jargon, acronyms, and bureaucratic phrasing.

## Sentences

- One idea per sentence. Split long compound or nested sentences into several short ones.
- Prefer active voice and a direct subject–verb order over passive constructions.
- Cut filler and throat-clearing; keep the substance.

## Vocabulary

- Use everyday Thai words over formal, legal, or academic ones (e.g. say "ต้องทำ" rather than a heavy nominalized phrase).
- Avoid legalese and academic nominalizations; turn them into plain verbs and nouns.
- Expand an abbreviation or acronym in plain Thai the first time it appears, then reuse the short form.

## Jargon and technical terms

- Translate the term into Thai, then put the original English term in parentheses once: `การเข้ารหัส (encryption)`.
- After the first mention, reuse the Thai term alone.
- Do **not** add a separate definition or gloss — keep the bilingual term only (this skill does not explain terms beyond the parenthetical English).

## Must NOT be simplified away

Simplify the *wording*, never the *substance*. These must survive verbatim in meaning:

- Numbers, dates, amounts, percentages, thresholds, and units.
- Conditions and logic: "if / unless / only when / except / provided that".
- Exceptions, caveats, and warnings.
- Negations — never flip or drop a "not".
- Named entities: people, products, organizations, systems, file/feature names.

## Not summarization

Every piece of information in the source must still be present in the output. You may split, reorder words within a sentence, or rephrase — but never drop a fact, merge away a distinction, or shorten by omission. If the output is shorter only because wording got tighter, that is fine; if it is shorter because content is missing, that is wrong.

## Tables and lists

Keep tables as tables and lists as lists, with the same rows, columns, and items in the same order. Simplify only the wording inside each cell or item.

## Untranslated elements

Leave these strictly in English: code blocks, inline code, filenames, URLs, command names, and variable or identifier names.
