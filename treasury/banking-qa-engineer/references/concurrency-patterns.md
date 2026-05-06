# Concurrency patterns and anti-patterns

Audit every state-mutating flow. Money movement, balance updates, KYC
state transitions, and audit-log writes must each pass these checks.

## 1. Transaction boundary

- [ ] The smallest write unit that must be atomic is wrapped in a single
      DB transaction.
- [ ] No external API call inside a DB transaction — extract before
      `BEGIN` or after `COMMIT` to avoid long-held locks.
- [ ] Transaction scope is named (function-level, request-level, saga
      step) and matches the business invariant.

**Anti-pattern:** transaction wraps the whole HTTP handler including
`http.Post` to a partner bank. A 30-second partner timeout holds DB
locks for 30 seconds.

## 2. Isolation level

- [ ] Default isolation level chosen explicitly (not silently
      `READ COMMITTED`).
- [ ] Money-movement flows use `REPEATABLE READ` or higher, or use
      `SELECT ... FOR UPDATE` to prevent phantom reads on balance.
- [ ] Optimistic concurrency (version column / `WHERE version = N`)
      paired with a retry budget.

**Anti-pattern:** `READ COMMITTED` + `SELECT balance; ... ; UPDATE
balance = balance - X` allows two parallel disbursements both reading
the same starting balance.

## 3. Idempotency

- [ ] Every state-mutating endpoint accepts an idempotency key.
- [ ] Idempotency key is persisted with the result hash; replay with the
      same key returns the same result, not a fresh side-effect.
- [ ] Dedup window matches the operating model (commonly 24h–7d).
- [ ] Kafka consumers use a dedup key per message and persist the
      processed-set.

**Anti-pattern:** retry-safe at the infra layer (Kafka redelivery) but
the consumer's `INSERT INTO ledger` is non-idempotent → duplicate ledger
entries on broker restart.

## 4. Lock ordering

- [ ] If multiple rows are locked, the order is fixed across all code
      paths (e.g., always `SELECT FOR UPDATE` by `account_id ASC`).
- [ ] No nested transactions across services with circular dependencies.
- [ ] Distributed locks (Redis / DB advisory) have explicit TTL and
      fence tokens.

**Anti-pattern:** transfer A→B locks A then B; reversal B→A locks B
then A. Concurrent calls deadlock.

## 5. Saga / process manager

- [ ] Long-running flows (loan origination, multi-step payment) modeled
      as a saga with explicit compensating actions.
- [ ] Each saga step is idempotent.
- [ ] Saga state is persisted; replay from any step is safe.
- [ ] No business logic in the orchestrator — orchestrator owns flow,
      domain processor owns invariants.

## 6. Read-your-writes

- [ ] After a write, the same actor's next read sees the write
      (read-your-writes guarantee).
- [ ] Read replicas not used for read-after-write paths without
      session-stickiness or explicit primary read.

## 7. Dual-control

- [ ] Underwriter overrides, manual disbursement, refunds, and
      write-offs require two distinct actor approvals.
- [ ] The two approvals are recorded with actor + timestamp; same-actor
      double-approve is rejected at the app layer (not just the UI).

## 8. Race-on-creation (uniqueness)

- [ ] Unique constraints on natural keys (account number, loan id,
      idempotency key) enforced at the DB.
- [ ] Application code handles unique-violation as a normal flow, not
      an error.

**Anti-pattern:** `if not exists: insert` without a unique constraint;
two parallel calls both observe "not exists" and both insert.

## 9. Time correctness

- [ ] Use monotonic clock (`time.Since`) for durations.
- [ ] Use UTC for absolute timestamps.
- [ ] Assume DST changes do not roll back ledger timestamps.

## 10. Failure-path correctness

- [ ] Every failure path either rolls back the transaction or compensates
      in a saga.
- [ ] No partial writes survive a failure ("ledger entry written, balance
      not decremented").
- [ ] Background-job retry policy bounded; poison-pill DLQ in place.

## Auto-rejections

The QA skill auto-rejects on any of:

- Missing transaction boundary on a state-mutating flow.
- `READ COMMITTED` + read-modify-write on monetary state.
- Non-idempotent endpoint that mutates money or KYC state.
- Idempotency key absent on disbursement / settlement / refund.
- Single-actor approval on flows requiring dual-control.
- Unbounded retry on a poison message.
