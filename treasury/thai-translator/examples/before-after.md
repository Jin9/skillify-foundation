# Example: source → plain-Thai

A short worked example of the skill's output. The source has a heading, a bullet
list, a technical term, and a numeric condition. The plain-Thai version keeps the
same structure, carries the bilingual term, and preserves the number and the
condition exactly — only the wording is simplified.

## Source (input)

```markdown
## Account Lockout Policy

The authentication subsystem enforces a lockout after repeated failures:

- After 5 consecutive failed login attempts, the account is temporarily locked.
- The lockout persists for 30 minutes unless an administrator clears it.
- Encryption of credentials at rest is mandatory.
```

## Output (`*-th.md`)

```markdown
## นโยบายการล็อกบัญชี (Account Lockout Policy)

ระบบยืนยันตัวตน (authentication) จะล็อกบัญชีเมื่อใส่รหัสผิดซ้ำหลายครั้ง:

- ถ้าใส่รหัสผ่านผิดติดต่อกัน 5 ครั้ง บัญชีจะถูกล็อกชั่วคราว
- บัญชีจะถูกล็อกนาน 30 นาที เว้นแต่ผู้ดูแลระบบจะปลดล็อกให้ก่อน
- ต้องเข้ารหัส (encryption) รหัสผ่านที่เก็บไว้เสมอ
```

## What to notice

- Same heading and same three bullets, in the same order (structural parity).
- `authentication` and `encryption` kept as Thai + English in parentheses.
- The number `5`, the duration `30 นาที`, the condition `เว้นแต่` (unless), and the
  requirement `ต้อง...เสมอ` (mandatory) are all preserved — nothing was summarized away.
- Sentences are short and use everyday words.
