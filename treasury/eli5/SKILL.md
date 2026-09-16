---
name: eli5
description: >
  Explains one real thing to the user in simple words and short sentences, directly and by its real name,
  with no analogies, metaphors, or stand-in scenarios: a task, plan, or ticket; code, a diff, or an error;
  a concept, term, or technology; a document or spec. First pass about 120 words (one-line answer, 3 to 5
  short steps, what it means for the user); saying "more" adds one layer without repeating. Replies in the
  user's language, English or Thai. Use when the user asks to "explain this simply", "explain in simple
  words", "explain like I'm 5", "ELI5 this", "explain this task to me", "what am I actually doing here?",
  or says "I don't get it" or "what does this mean?" about something concrete (Thai: "อธิบายง่ายๆ",
  "ไม่เข้าใจ"). Chat-only: it explains, it never does the task or makes the decision. Do NOT use for
  trade-off advice inside a live decision (principal-advisor), for plain-Thai rendering of a whole
  document (thai-translator), or for a technical deep dive into code (understand-explain).
---

# ELI5

## Purpose

Explain one real thing to the user in plain words, by its real name, with no analogies, so they understand it and can act on it. Chat markdown only: this skill never writes files, never does the task, and never makes the decision.

## When to use this skill

- Use when the user asks to "explain this simply", "explain in simple words", "explain like I'm 5", or "ELI5 this".
- Use when the user asks "explain this task to me" or "what am I actually doing here?" about the current work item, plan, or ticket.
- Use when the user says "I don't get it" or "what does this mean?" (Thai: "อธิบายง่ายๆ", "ไม่เข้าใจ") about something concrete: a task, a ticket, code, a diff, an error, a term, or a document.
- Do NOT use when the message asks for a change, a fix, or a choice ("fix it", "which should I pick?"): that is the task itself, or the principal-advisor skill for a live decision.
- Do NOT use to render a whole document in plain Thai while keeping every fact (thai-translator), or for a technical deep dive into code (understand-explain).

## Model & cost

Runs on a **mid** model, low to medium effort: the hard part is reading the real target correctly, not the wording. A **small** model is enough for one term or a short error with nothing to read. Escalate to a **frontier** model only for a large diff or a multi-file task whose current state must be reconstructed before it can be explained.

## Plain kernel

Every sentence of the reply follows these ten rules. Detail, word swaps, and the allowed moves live in `references/plain-kernel.md`.

1. Talk to the user as "you". Name the real thing by its real name (file, function, ticket ID, term), verbatim, code in backticks.
2. No analogies, metaphors, comparisons, or stand-in scenarios ("it's like…", "imagine…", "think of it as…"). Say what it is.
3. One idea per sentence. Aim for 15 words, never more than 20. Thai: one clause per sentence.
4. Everyday words. When a technical term is the real name, keep it and define it once under Words; in Thai, keep it in English as typed at work (`database`, not a coined Thai word).
5. Concrete: each step names a real file, value, command, or ticket, never "the system" when the name is known.
6. True first: read the thing before explaining it, and say what you are inferring.
7. Simplify wording, never meaning: numbers, conditions, negations, and exceptions survive.
8. Cap: 120 English words or 500 Thai characters per pass, tail excluded. Never repeat a sentence already sent.
9. No filler: no preamble, no "great question", no restating the question, no closing summary.
10. Explain, do not act: no fixing, no doing, no choosing. At most, name what the user must decide.

## Core workflow

1. **Identify the target and the confusion point.** *Entry:* a trigger phrase about something concrete. Resolve "this": the last message, pasted text, a file, a ticket, an error. Classify its kind: task/plan/ticket, code/diff, error, concept/term, or document/spec. Note the language the user wrote in. Take the most recent thing shown; ask one question only when two candidates are equally recent. *Exit when* the target, its kind, and the reply language are named.
2. **Read the real thing first.** *Entry:* target named. Open the actual file, diff, error location, ticket, or document when it is reachable in the session, read-only. Never explain from memory of it. If it is unreachable, explain from what the user gave and say so in the reply. Note the one to three facts the user needs: the goal, what failed, the definition. *Exit when* you can say in one sentence what it is and what the user needs from it.
3. **Pick the shape.** *Entry:* kind known. Take the per-kind filler for the bold line, the steps, and the For you line from `references/shapes-by-target.md`. Choose at most three terms for Words. *Exit when* the bold line is drafted and the steps are planned.
4. **Draft under the cap.** *Entry:* shape chosen. Fill the skeleton in Output format, applying the kernel. Use the word swaps and the banned-device list in `references/plain-kernel.md`. *Exit when* the draft fits the skeleton and the cap.
5. **Check.** *Entry:* a draft exists. Run the linter from this skill's folder with the draft on stdin, using a quoted heredoc so backticks survive:

   ```sh
   python3 scripts/check_plain_style.py - <<'EOF'
   [draft]
   EOF
   ```

   Do not save the draft as a file. Fix every reported violation and rerun until it exits 0. A clean run is necessary, not sufficient: do one visual pass for analogies in disguise (stand-in actors, "suppose you have", personification) and for meaning drift (numbers, negations, conditions). If the host cannot run scripts, apply the kernel list by hand. *Exit when* the exit code is 0, or the manual pass is done, and the visual pass is clean.
6. **Reply.** *Entry:* a clean draft. Send it as chat markdown and nothing else: no preamble, no file, no action. The plain style ends with this reply; the next non-explanation turn uses the normal style. *Exit when* the reply is sent.
7. **Follow up.** *Entry:* the user answers. "more" (Thai "เพิ่ม" or "ต่อ"): send the next layer per the layer contract below, linted as its own input with the default cap. "still don't get X": re-explain X only, with one allowed move from `references/plain-kernel.md`, same skeleton and cap. A pivot ("ok do it", "fix it", "which one?"): leave this skill; the normal workflow and skills resume, and nothing is done inside this skill. *Exit when* the user says they get it, or pivots.

## Output format

Chat markdown only. No files, no code changes, no commands run on the user's behalf. Every pass, first or later, uses this skeleton:

```
**[One line: what it is / what happened / what you are doing.]**

1. [Short step. Real names in backticks.]
2. [Short step.]
3. [Short step.]                      (2 to 5 numbered lines)

**For you:** [what this means for you right now, 1 or 2 short sentences.]

**Words:**                            (optional, at most 3 terms)
- **term**: plain meaning.

Say **more** for the next layer, or tell me which step is unclear.
```

- Cap: everything above the tail is at most 120 English words or 500 Thai characters. The tail is fixed text and does not count.
- Thai labels: `**สำหรับคุณ:**` and `**คำศัพท์:**`; Thai tail: `พิมพ์ **เพิ่ม** ถ้าอยากรู้อีกชั้น หรือบอกว่าข้อไหนยังไม่เข้าใจ`. A Thai Words entry keeps the English headword: `- **subquery**: ความหมายภาษาไทย`.
- Layer contract: layer 2 is why and mechanism (what causes what underneath); layer 3 is details, edge cases, exceptions, and numbers. A layer never repeats a sentence already sent. After layer 3, a further "more" narrows to one named part and restarts the skeleton on that part.
- A question to the user inside a reply is phrased as an imperative ("Tell me which step is unclear, 2 or 3."), because the linter rejects question marks in prose.
- For you may name options and what the user must decide, never which option to choose.

## Constraints & anti-patterns

- DO NOT use analogies, metaphors, similes, "imagine…" scenarios, stories, characters, personification, or rhetorical questions.
- DO NOT compare to a technology the user did not raise. Comparing to something already in the user's own task (their other function, their last ticket) is allowed.
- DO NOT do the task, fix the code, run the command, or choose for the user.
- DO NOT write files or change anything; the reply is chat markdown only.
- DO NOT drop or soften numbers, conditions, negations, or exceptions. Simplify wording, never meaning.
- DO NOT state what you did not read; say what was unreachable.
- DO NOT translate or alter identifiers, error text, commands, or product names.
- DO NOT carry the plain style into the next non-explanation turn.
- DO NOT repeat an earlier sentence in a "more" layer.
- MUST ALWAYS reply in the language the user wrote in; Thai replies read like a developer's chat message, technical nouns in English as typed at work (Thai notes in `references/plain-kernel.md`).
- MUST ALWAYS read the real thing first when it is reachable.
- MUST ALWAYS end every pass with the fixed tail.

## Examples

### Example 1: an error, in English

**User says**: pastes `! [rejected] main -> main (non-fast-forward)` and writes "I don't get it, what does this mean?"

**Action**:
1. Kind = error; language = English. Read the full `git push` output in the session.
2. Shape = error from `references/shapes-by-target.md`: what failed, where it stopped, what you can do.
3. Draft, lint on stdin until exit 0, visual pass, reply.

**Result** (78 words above the tail; the same text is the fixture `examples/error-first-pass.md`):

> **Your push was rejected because the remote branch has commits your local branch does not have.**
>
> 1. Someone pushed to `main` after your last `git pull`.
> 2. Git refuses to overwrite their commits, so it stops with `! [rejected] main -> main (non-fast-forward)`.
> 3. Nothing is lost. Your commits are still on your machine.
>
> **For you:** run `git pull --rebase origin main`, fix any conflicts, then `git push` again.
>
> **Words:**
> - **non-fast-forward**: your branch is missing commits that the remote already has.
> - **rebase**: replay your commits on top of the newest remote commits.
>
> Say **more** for the next layer, or tell me which step is unclear.

### Example 2: a concept, in Thai, then "more"

**User says**: "CTE คืออะไร ไม่เข้าใจ" while editing a query that has a subquery in `FROM`.

**Action**:
1. Kind = concept/term; language = Thai. Read the query the user is editing.
2. Shape = concept: a one-sentence definition, what it does, what it is not, why it matters in this query.
3. Draft in Thai with Arabic numerals, keeping `subquery` and `temp table` in English and defining them under คำศัพท์; lint (detected as `th`, 425 characters); reply with the Thai tail.
4. On "เพิ่ม", send layer 2 (how this database runs a CTE, and when it materializes it) with only new sentences.

**Result**: the Thai reply is the fixture `examples/thai-concept-first-pass.md`; its textbook-register twin, which the linter rejects (P5 on six lines, P6 once), is `examples/thai-textbook-before.md`. The shape of a first pass followed by a "more" layer, in English, is `examples/concept-with-more.md`.

## Validation checklist

Before sending each pass, verify:
- [ ] The reply talks to "you" and names the real thing by its real name.
- [ ] `scripts/check_plain_style.py` exits 0 on the exact text, and the visual pass found no analogy, stand-in scenario, personification, or rhetorical question.
- [ ] Above the tail: at most 120 English words or 500 Thai characters, 2 to 5 numbered steps, a For you line, at most 3 Words terms.
- [ ] Every number, condition, negation, and exception from the source survives.
- [ ] Identifiers, error text, and commands are verbatim; Thai keeps technical nouns in English.
- [ ] Nothing was done, written, or decided on the user's behalf.
- [ ] A "more" layer adds only new sentences.

## References

| Need | Reference |
|---|---|
| Kernel detail: word swaps (EN and TH), banned devices with fixes, allowed moves when the user is still confused, Thai notes, manual-pass list | `references/plain-kernel.md` |
| Per-kind fillers for the skeleton and the layer plan | `references/shapes-by-target.md` |
| Linter usage, rule IDs P1 to P9, known limits | `scripts/check_plain_style.py` |
| Worked replies that are also linter fixtures | `examples/error-first-pass.md`, `examples/concept-with-more.md`, `examples/thai-concept-first-pass.md`; the analogy-laden `examples/bad-analogy-before.md` and its plain rewrite `examples/bad-analogy-after.md`; the textbook-Thai `examples/thai-textbook-before.md` that the Thai P5 list rejects |
