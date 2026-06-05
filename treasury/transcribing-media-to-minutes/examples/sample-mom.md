# Example — Thai/English standup recording → MOM + transcript

A worked, shortened example of the skill's two outputs for a ~6-minute bilingual team standup
(`sprint-standup-2026-06-05.m4a`). Trimmed for illustration.

## Input (what the agent receives)

- Media: `sprint-standup-2026-06-05.m4a` (audio, ~6 min, mostly Thai with English dev terms)
- Optional context given: attendees = Ploy (PM), Nattapong "Dev" (Backend), May (QA)

## Output 1 — `sprint-standup-2026-06-05-transcript.md` (excerpt, verbatim)

    [00:00:12] Ploy: เริ่มเลยนะคะ เมื่อวาน dev ทำอะไรไปบ้าง
    [00:00:21] Dev: เมื่อวานปิด API ตัว payment เสร็จ วันนี้จะ deploy ขึ้น staging ก่อนเที่ยง
    [00:01:03] May: ขอ test case ของ payment ด้วยนะ เดี๋ยว QA รันให้ภายในวันศุกร์
    [00:02:30] Dev: อันนี้ติด compliance นิดนึง ต้องรอ legal เคลียร์เรื่อง KYC ก่อน
    [00:03:10] [inaudible]

## Output 2 — `sprint-standup-2026-06-05-mom.md` (excerpt)

    # รายงานการประชุม (Minutes of Meeting)

    ## ข้อมูลการประชุม (Meeting info)
    - หัวข้อ (Title): Sprint Standup
    - วันที่ (Date): 2026-06-05
    - ความยาว (Duration): 00:06:14
    - ไฟล์ต้นทาง (Source media): sprint-standup-2026-06-05.m4a
    - ผู้เข้าร่วม (Attendees): Ploy — PM, Nattapong (Dev) — Backend, May — QA

    ## ⚠️ จุดที่ควรตรวจ (Low-confidence / inaudible — review first)
    - [00:03:10] เสียงไม่ชัด (inaudible) — มีประเด็นที่ฟังไม่ออก โปรดตรวจไฟล์ต้นทาง

    ## มติ / ข้อสรุป (Decisions)
    - [00:00:21] Dev จะ deploy (ปล่อยขึ้นใช้งานจริง) API payment ขึ้น staging ก่อนเที่ยงวันนี้

    ## สิ่งที่ต้องทำ (Action items)
    | งาน (Task) | ผู้รับผิดชอบ (Owner) | กำหนดเสร็จ (Due) | อ้างอิงเวลา (Timestamp) |
    |------------|----------------------|------------------|--------------------------|
    | Deploy API payment ขึ้น staging | Dev (Nattapong) | ก่อนเที่ยงวันนี้ | [00:00:21] |
    | รัน test case ของ payment | May (QA) | ภายในวันศุกร์ | [00:01:03] |

    ## ประเด็นค้าง (Open questions)
    - การทำ KYC ของ payment ติดเงื่อนไข compliance (กฎ/ระเบียบ) — รอ: ทีม Legal เคลียร์ก่อน [00:02:30]

## Why this is faithful

- Every decision/action carries a timestamp into the transcript.
- The inaudible span at [00:03:10] is surfaced for review, not invented.
- English dev terms (deploy, staging, API, KYC, compliance) are kept verbatim and glossed once in
  plain Thai — the meaning is unchanged.
- Nothing was added: the standup gave no exact due date, so due dates copy what was said ("ภายในวัน
  ศุกร์"), and the KYC point is an Open question, not an Action, because no owner committed to it.
