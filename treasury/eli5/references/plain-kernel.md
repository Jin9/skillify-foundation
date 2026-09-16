# Plain kernel: detail

Depth for kernel rules 2, 4, 9, and 10 of `SKILL.md`, plus what to do when the user is still confused. Read only the section you need; the ten rules themselves stay in `SKILL.md`.

## Word swaps

Prefer the everyday word when it is still correct. The linter flags most of the left column, and the Thai textbook forms, as rule P5; its message names the swap.

| Instead of | Say |
|---|---|
| utilize, leverage | use |
| facilitate | help |
| in order to | to |
| prior to | before |
| subsequently | later |
| it should be noted, it is important to note | (delete; state the point) |
| aforementioned | (name the thing again) |
| thus, hence | so |
| whereby | where |
| via | through |
| e.g., i.e. | for example, that is |
| additionally, furthermore, moreover | also |
| in the event that | if |
| with regard to | about |
| commence, terminate | start, stop |
| instantiate | create |
| invoke | call |
| persist (a record) | save |

Thai: the textbook forms and their everyday replacements live in the linter's Thai P5 list, and its message names the swap; how to write the sentence itself is under Thai notes.

## Banned devices and their fixes

| Device | Looks like | Fix |
|---|---|---|
| Analogy or simile | "it works like a queue", "similar to a post office" | State what it does: "it holds requests and hands them out in arrival order". |
| Metaphor | "the gateway is the front door" | Name the function: "the gateway is the one service that accepts outside requests". |
| Made-up scenario | "imagine a restaurant", "suppose you have a shop" | Use the user's real case: their real request, file, or ticket. |
| Personification | "the server wants", "the compiler complains" | Say the mechanism: "the server returns 401", "the compiler stops with error E0308". |
| Rhetorical question | "So what went wrong?" | Delete it; the next sentence already answers it. |
| Story or characters | "Alice sends Bob a message" | Use the real parties: "your app sends the payment service a request". |
| Humor, drama | "the dreaded null pointer" | Delete the adjective. |
| Comparison to another technology | "it is Redis for Postgres" | Banned unless the user raised that technology or it is already in their task; then answer in one sentence with the one difference. |

Negative example and its fix. A research report describes a request router as "an air traffic controller for prompts". Plain version: "A router is a small program that reads each request and picks which model gets it."

## Allowed moves when the user is still confused

Use exactly one per turn, applied to the part the user named, inside the same skeleton and cap.

1. Split the step: break the unclear step into two or three smaller steps.
2. Real example: show the actual input and output for the user's own case, with the real values.
3. Define: define the one word the step leaned on, in one sentence.
4. Say what it is not: name the nearest real thing it is confused with and the one difference (a real neighbor concept, not a metaphor).
5. Show the line: quote the exact line, command, or field in backticks and say what it does.
6. Ask: if you cannot tell which part is unclear, end with an imperative: "Tell me which step is unclear, 2 or 3."

## Thai notes

- Pronoun คุณ for the user; no ท่าน, and no ครับ/ค่ะ tails, since the agent has no gender to pick one.
- Register: type it the way a Thai developer types to a colleague in chat. Thai sentence glue with everyday connectors (ถ้า, แต่, คือ, พอ…ก็, ก็เลย) and ๆ where speech doubles the word (ง่ายๆ, จริงๆ); one นะ per pass where a colleague would type it. Technical nouns stay in English, in Latin script, exactly as typed at work (database, DB, server, deploy, commit, subquery, CTE): no Thai coinage, no parenthetical gloss. Verbs go the way people type them, Thai or English (รัน, เช็ค, เซฟ, merge, deploy). Identifiers, error text, commands, product names, and code stay in English inside backticks.
- Write each sentence from the idea, not word for word from the English skeleton: "what you are actually doing" is not "สิ่งที่คุณทำจริง" but "จริงๆ แล้วงานนี้คือ…". If it would sound odd read aloud to a colleague, retype it.
- Arabic numerals, not Thai numerals.
- Sentence rule: one clause per sentence, with a space at every clause end. The linter flags a run of more than 90 Thai characters without a space (rule P3).
- Cap: 500 Thai characters above the tail, English words and Words entries included (rule P4).

## Running the linter

Invocation: the quoted heredoc in workflow step 5 of `SKILL.md`. `--cap N` overrides the cap; `--lang en|th` overrides detection (the unit of P3 and P4, and the Thai P5 list). Violations print as `file:line: RULE message` on stderr; a clean run prints `ok: 1 input(s) clean (en, 97 words)` on stdout, so the count tells you how much budget is left.

The manual pass owns what the linter cannot see:

- personification and stand-in actors ("the server wants", "a waiter")
- "suppose you have a…" and "say you have a…" scenarios, and a bare "as if"
- soft Thai markers used figuratively: เหมือน, เหมือนกับ, คล้ายกับ, สมมติว่า, เปรียบเทียบ
- textbook coinages or essay connectors the Thai P5 list does not know (add a recurring one to the list), and `user ต้อง…` addressing the reader, which P6 only knows as ผู้ใช้ต้อง
- one idea per Thai sentence
- at most three Words terms, each defined where it first appears
- meaning drift: every number, negation, condition, and exception still present

## Why these rules

Background only; no action needed. "Define or gloss every term at first use" and "start with plain definitions" are applied acceptance criteria in the research vault's review notes (hybrid execution boundary routing, 2026-05-17; OpenTelemetry agent-trace report, 2026-05-17). "Short sentences, no filler" follows the Caveman skill write-up and the "one to three sentences, no filler" rules-file pattern in the skill corpus; a 2026-05-18 writing-quality report found that models preferred concise replies over padded ones. The no-analogy rule has no corpus grounding and is the user's house rule: the corpus even lists "create analogy" as a valid explaining move, and the air-traffic-controller line above comes from the vault itself. The Thai register is the user's house rule too (2026-09): the earlier Thai-first gloss rule, borrowed from thai-translator whose reader is a general adult, produced textbook words no developer types; this skill's reader is a developer, so it deliberately differs.
