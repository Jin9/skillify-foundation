# Thai-primary bilingual minutes — rendering and plain language

Deep guidance for step 5 of `SKILL.md`. Goal: minutes that a Thai reader understands at a glance,
that keep the English terms people actually use, and that never change the meaning of what was said.

## Language policy

- **Thai-primary.** Write the minutes mostly in Thai. This is the default when the spoken meeting is
  mostly Thai.
- **Keep English terms verbatim** when they are the words people use on the job — product names,
  banking/tech jargon, acronyms (KYC, AML, API, PR, sprint, deploy). Do **not** force Thai
  translations that no one says.
- **Gloss once, simply.** The first time a jargon term or a complex decision appears, add a short
  plain-Thai explanation in parentheses or on the next line. Explain the *idea*; do not restate it.
- **Language override.** If the user asked for English-only or fully mirrored bilingual output, follow
  that instead; otherwise use Thai-primary.

## Plain-language rules (simple explanation, same meaning)

1. Short sentences, everyday Thai. Prefer common words over bureaucratic ones.
2. One idea per bullet. Split run-on discussion into discrete points.
3. Active voice with a clear actor: "ทีม X จะ…", not "จะมีการดำเนินการ…".
4. Spell out an acronym once, then reuse it.
5. The plain explanation is an *aid*. It must not add a conclusion, soften, or sharpen the original.

## Worked phrasing

- Spoken: "เดี๋ยว dev จะ deploy ตัว hotfix ขึ้น prod คืนนี้นะ"
  → Minutes: "ทีม Dev จะ deploy (ปล่อยขึ้นใช้งานจริง) hotfix ขึ้น production คืนนี้ — ผู้รับผิดชอบ: Dev"
- Spoken: "อันนี้ติด compliance ต้องรอ legal เคลียร์ก่อน"
  → Minutes: "ประเด็นนี้ติดเงื่อนไข compliance (กฎ/ระเบียบ) — รอทีม Legal อนุมัติก่อน → เป็นประเด็นค้าง"

## Mini glossary (Thai ↔ English: meeting + banking)

| English | Thai (plain) |
|---------|--------------|
| Minutes of Meeting (MOM) | รายงานการประชุม / บันทึกการประชุม |
| Action item | สิ่งที่ต้องทำ / งานที่มอบหมาย |
| Owner | ผู้รับผิดชอบ |
| Due date | กำหนดเสร็จ |
| Decision | มติ / ข้อสรุป |
| Open question | ประเด็นค้าง / คำถามที่ยังไม่ตอบ |
| Risk | ความเสี่ยง |
| Blocker | สิ่งที่ติดขัด / ตัวขวาง |
| Compliance | การปฏิบัติตามกฎเกณฑ์ |
| KYC | การพิสูจน์ตัวตนลูกค้า |
| Deploy / release | ปล่อยขึ้นใช้งานจริง |

Extend the glossary from the meeting's own vocabulary; always keep the source term beside the gloss.

## Tone

Neutral, professional minutes register. Attribute decisions to roles/teams, not to gossip. Do not
record side chatter unless it produced a decision or an action.
