# Rust Idioms, Async & Ownership

Generation-side guidance for the language-particular shapes `crafting-rust-code` produces. This is the constructive companion to `review-rust-code`'s adversarial `rust-idioms-and-async.md` — write code this way so the reviewer's R1/R2/R3 scans pass.

## 1. Error taxonomy (thiserror at boundaries, anyhow in bins)

- **Libraries / domain modules**: a `#[derive(thiserror::Error)]` enum, one variant per failure mode a caller may branch on. Preserve the source chain with `#[from]` / `#[source]` / `#[error(transparent)]`. Do not over-engineer — group related failures; if callers always handle 20 variants identically, you have too many.
- **Binaries / top-level glue**: `anyhow::Result` with `.context("doing X")` for human-readable context. `anyhow` is acceptable only where no caller branches on the type.
- **Never** `map_err(|e| MyError::Other(e.to_string()))` — it discards the cause. Keep the error a typed value, not a sniffed string.
- **Error classification** (drives behavior): give each domain error a class at the type level — an `enum` discriminant or `fn class(&self) -> ErrorClass { Client | Server | Dependency }`. The class, not edge guessing, decides the HTTP status, the retry decision, and the `tracing` error attribute.
- **No panics on production paths**: no `.unwrap()` / `.expect()` / `panic!` / `unreachable!` / `todo!` / `unimplemented!` reachable from a handler/worker. Use `?` and return the typed error. The single allowance is `expect("invariant: <why this cannot fail>")` for a genuinely unreachable, documented invariant.

## 2. Async & tokio discipline

The executor is cooperative: anything that blocks a worker thread or holds a lock across an await point starves every other task on that thread. Current production consensus (Oxide RFD 400, Tokio docs, common-mistakes write-ups):

- **No blocking in async**: no `std::thread::sleep`, blocking `std::fs::*` / `std::net::*`, or >~10µs of CPU work inside an `async fn` / `.await` path. Use the async equivalent (`tokio::fs`, `tokio::time::sleep`) or `tokio::task::spawn_blocking` for CPU/blocking work.
- **No `std::sync::Mutex` guard held across `.await`**: the resulting future is `!Send` and the guard pins the lock across a yield. Either (a) drop the guard before the await, or (b) restructure. Hold a lock only for the brief, await-free critical section.
- **Prefer message-passing over shared mutable state**: instead of `Arc<Mutex<State>>` shared across tasks, give one **manager task** sole ownership of the state, have callers send requests over an `mpsc` channel, and reply over a `oneshot`. This sidesteps cancellation-corrupts-shared-state entirely. Reach for `tokio::sync::Mutex` only when a lock genuinely must be held across an await and the message-passing refactor is disproportionate.
- **Cancellation safety**: a future may be dropped at any `.await`. In `tokio::select!`, only use cancel-safe futures in branches — `AsyncReadExt::read_exact`, `AsyncWriteExt::write_all`, `Mutex::lock`, and an in-flight DB transaction are **not** cancel-safe and can lose data when their branch loses the race. Keep partial state out of the future, or move non-cancel-safe work behind a manager task.
- **No detached `tokio::spawn` on a request path**: a detached task outlives the request and drops its failure on the floor. Use `JoinSet` (or scoped tasks) so failures and cancellation propagate. Don't hold a pooled resource (DB connection, `PoolConnection`) across an await you don't need to — it starves the pool.
- **Propagate deadlines**: thread the request timeout / `CancellationToken` into every `sqlx`, `reqwest`, `tonic`, and producer call. A handler with no deadline can hang forever.
- **Send/'static**: make `Send`/`Sync` bounds visible in signatures. A `tokio::spawn` body capturing `Rc`/`RefCell` or a non-`Send` borrow only compiles under `current_thread` and will break a multi-thread runtime.

## 3. Unsafe discipline

- Default to `#![forbid(unsafe_code)]`. In banking/regulated code, the absence of `unsafe` is the norm.
- If `unsafe` is unavoidable: every block carries a `// SAFETY:` comment naming the invariants it relies on (not "this is safe"); every `unsafe fn` carries a `# Safety` rustdoc section stating the precondition; `unsafe impl Send`/`Sync` is justified field-by-field (including future fields); `mem::transmute`, `Pin::new_unchecked`, `Box::leak`, `hint::unreachable_unchecked` each carry a written rationale. Exercise the `unsafe` path with tests (and MIRI when available).

## 4. Ownership, traits & visibility

- **Ownership first**: model who owns each value. Prefer move/borrow over `clone()`; reach for `Arc` only for genuinely shared immutable data, `Arc<Mutex<_>>` only when shared mutation can't be a manager task. Borrow-checker fights are a design signal, not a `clone()` problem.
- **Newtypes for value objects**: wrap primitives that carry an invariant (`struct Money(Decimal)`, `struct AccountId(Uuid)`) so invalid states are unrepresentable and signatures are self-documenting.
- **Trait vs generic**: generics (`<T: Trait>` / `impl Trait`) monomorphize — fast, larger binary, one concrete type per call site; trait objects (`dyn Trait`, must be object-safe) give dynamic dispatch and heterogeneous storage. Pick by whether you store mixed implementors and by call-site/binary-size cost.
- **Tell, Don't Ask**: push behavior with an invariant onto the owning type's method (`account.debit(amount)?`) rather than pulling fields out and deciding in the handler — but only when the mutation has an invariant the type can defend and is reused; don't add premature getters/setters.
- **Keep the API surface narrow**: default `pub(crate)`; widen to `pub` only deliberately. Use `#[non_exhaustive]` on growable public enums/structs. Avoid glob `pub use` that re-exports more than the contract names. Use a sealed-trait pattern when a trait must be public but not externally implementable.
