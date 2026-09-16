**Idempotent means doing the same request twice has the same effect as doing it once.**

1. The client sends a request, the network drops the reply, and the client sends it again.
2. Without idempotency, the server runs the action twice, so you get two payments.
3. With idempotency, the server sees the repeat and does not run the action again.
4. The usual way is an idempotency key: a unique id the client sends with each request.

**For you:** add an `Idempotency-Key` header to `POST /payments` and store each key with its result.

**Words:**
- **idempotency key**: a unique id that marks two requests as the same request.

Say **more** for the next layer, or tell me which step is unclear.

---

**Store the key before you run the action, or a fast retry can still run it twice.**

1. Insert the key with status `pending` in the same database transaction as the action.
2. If a request arrives with a key that already exists, return the stored result and stop.
3. Keep keys for at least 24 hours, because clients retry long after the first attempt.
4. Use a key for one request body only. A new body needs a new key.

**For you:** add a unique index on the key column, so two requests cannot both insert it.

Say **more** for the next layer, or tell me which step is unclear.
