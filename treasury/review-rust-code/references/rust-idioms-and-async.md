# Rust Idioms, Async, and Unsafe Reference

> Adversarial scan reference for LLM reviewers. Cite each finding as `file:line` and tag with the category enum from the final table. When in doubt, prefer false positives — let the human triage. Every code shape below is a **flagworthy pattern**, not an exemplar.

---

## Section 1: Error Handling

### `?` operator and `From`/`Into` plumbing

The `?` operator extracts `Ok` or returns the `Err`, but the returned error type must satisfy `From<InnerErr> for OuterErr`. Missing `From` impls produce `the trait From<X> is not implemented for Y` errors. See [doc.rust-lang.org/book/ch09-02](https://doc.rust-lang.org/book/ch09-02-recoverable-errors-with-result.html).

```rust
// FLAG: ad-hoc `.map_err(|e| MyError::Io(e.to_string()))?` patterns when a
// proper `#[from]` would carry source chain + backtrace. The map_err is
// lossy — `source()` will return None downstream.
let f = File::open(p).map_err(|e| MyError::Io(e.to_string()))?;
```

### `thiserror` vs `anyhow` — boundary rule

| Crate       | Use in                | Why                                                                 |
|-------------|-----------------------|---------------------------------------------------------------------|
| `thiserror` | libraries, domain     | Stable, named variants; callers can match. Source: [docs.rs/thiserror](https://docs.rs/thiserror/latest/thiserror/) |
| `anyhow`    | binaries, app glue    | Erased `dyn Error`; context strings; not matchable. Source: [docs.rs/anyhow](https://docs.rs/anyhow/latest/anyhow/) |

**Flag**: a library crate exporting `pub fn foo() -> anyhow::Result<T>` in its public API. Callers cannot match on variants. **Flag**: an application top-level using `thiserror` everywhere when an `anyhow::Error` with context would do.

### `#[from]`, `#[source]`, `#[error(transparent)]`

```rust
#[derive(thiserror::Error, Debug)]
pub enum DomainError {
    #[error("io: {0}")]
    Io(#[from] std::io::Error),          // auto From<io::Error>; preserves source()
    #[error("parse")]
    Parse { #[source] inner: ParseError }, // source chain without From
    #[error(transparent)]
    Other(#[from] anyhow::Error),         // forwards Display + source unchanged
}
```

- `#[from]` generates `From<T>`; the variant must have exactly one non-attribute field (plus optional backtrace).
- `#[source]` only marks the source chain — does not generate `From`.
- `#[error(transparent)]` forwards `Display` and `source()` straight through; use only for pass-through wrappers.

### `Context::context` vs `with_context`

```rust
// context: static string, evaluated unconditionally
file.read_to_string(&mut s).context("read config")?;
// with_context: closure, only formatted on Err — use for any allocation or fmt
let bytes = fs::read(&p).with_context(|| format!("read {}", p.display()))?;
```

**Flag**: `.context(format!(...))` on a hot path — the format runs on every success.

### Result type alias

```rust
// Reduces noise and lets one error type evolve in one place.
pub type Result<T, E = DomainError> = std::result::Result<T, E>;
```

### Type-level vs edge-decoded error class

Encode the error category in the **type system**, not in a string field. `match err { DomainError::NotFound => 404, .. }` is reviewable; `if err.to_string().contains("not found")` is not. **Flag** any string-sniffing on errors.

### Panicking APIs — when each is acceptable

| API                 | Acceptable in production?                                  |
|---------------------|------------------------------------------------------------|
| `panic!()`          | Truly unrecoverable; rare. Prefer typed errors.            |
| `unwrap()`          | Tests, prototypes, `const`. **Flag** in library code.      |
| `expect("invariant: <why>")` | OK with a documented invariant the compiler can't see. |
| `unreachable!()`    | Match arms eliminated by guards/types ([std docs](https://doc.rust-lang.org/std/macro.unreachable.html)). |
| `todo!()` / `unimplemented!()` | **Flag** anywhere outside a WIP branch.         |

```rust
// OK
let port: u16 = env::var("PORT").unwrap().parse()
    .expect("invariant: PORT is validated at startup");
// FLAG
let user = repo.find(id).unwrap();   // silent crash in request path
```

### `let ... else { return ... }` pattern

```rust
let Some(account) = repo.get(id).await? else {
    return Err(DomainError::AccountMissing(id));
};
```

**Flag** deeply nested `if let Some(x) = ... { if let Some(y) = ... { ... } }` ladders where `let else` would flatten control flow.

### Adversarial scan questions — Section 1

- Is any `?` site silently losing the error chain (`.map_err(|e| Custom::Msg(e.to_string()))`)?
- Does the public library API leak `anyhow::Error`?
- Does any `expect` lack a documented invariant?
- Does the code branch on `err.to_string()` instead of variants?

---

## Section 2: Async + Tokio

### `async fn` vs `impl Future`, and `async fn` in traits

`async fn foo() -> T` desugars to `fn foo() -> impl Future<Output = T>`. Native `async fn` in traits is stable for static dispatch but is **not `dyn`-compatible** — `dyn Trait` for a trait with `async fn` will not compile. Use `#[async_trait]` for `dyn` dispatch. See [RFC 3185](https://rust-lang.github.io/rfcs/3185-static-async-fn-in-trait.html).

**Flag**: a service trait with native `async fn` that is then used as `Box<dyn Service>`.

### Blocking-in-async violations

| Pattern                              | Fix                                          |
|--------------------------------------|----------------------------------------------|
| `std::thread::sleep(d)` in async     | `tokio::time::sleep(d).await`                |
| `std::fs::read(p)` in async          | `tokio::fs::read(p).await` or `spawn_blocking` |
| `std::net::TcpStream::connect` async | `tokio::net::TcpStream::connect`             |
| CPU loop > ~10 µs without yield      | `tokio::task::spawn_blocking(move || ...)`   |

See [`spawn_blocking` docs](https://docs.rs/tokio/latest/tokio/task/fn.spawn_blocking.html): use for **short** blocking work; long-lived workers belong on `std::thread::spawn`.

```rust
// FLAG: blocking call inside an async fn
async fn read_cfg(p: &Path) -> io::Result<String> {
    std::fs::read_to_string(p)   // blocks the executor thread
}
```

### `std::sync::Mutex` held across `.await`

`std::sync::MutexGuard` is `!Send`. If held across `.await`, the resulting future is `!Send` and cannot be `tokio::spawn`ed on a multi-thread runtime. Even on `current_thread`, holding it deadlocks the executor. See [tokio.rs/tutorial/shared-state](https://tokio.rs/tokio/tutorial/shared-state) and the [async book](https://rust-lang.github.io/async-book/07_workarounds/03_send_approximation.html).

```rust
// FLAG: lock alive across await
let mut g = state.lock().unwrap();
g.value = client.fetch().await?;     // future is !Send

// Fix A — scope the guard
let req = { let g = state.lock().unwrap(); g.build_request() };
let resp = client.send(req).await?;

// Fix B — only when truly needed: tokio::sync::Mutex (more expensive)
let mut g = state.lock().await;
g.value = client.fetch().await?;
```

Prefer **restructuring** over reaching for `tokio::sync::Mutex`.

### `Send + 'static` on `tokio::spawn`

`tokio::spawn` requires the future to be `Send + 'static`. **Flag** spawned closures capturing `Rc<T>`, `RefCell<T>`, raw pointers, or borrowed references with a non-`'static` lifetime. See [tokio.rs/tutorial/spawning](https://tokio.rs/tokio/tutorial/spawning).

```rust
// FLAG: Rc is !Send
let counter = Rc::new(RefCell::new(0));
tokio::spawn(async move { *counter.borrow_mut() += 1; });  // compile error on multi-thread
```

Note: this **compiles** on `tokio::task::spawn_local` under a `current_thread` runtime — that is the surprise to call out.

### Structured concurrency

| Primitive                                 | Use when                                                    |
|-------------------------------------------|-------------------------------------------------------------|
| `tokio::join!(a, b)`                      | Fixed N, all must complete, no early-exit on Err            |
| `tokio::try_join!(a, b)`                  | Fixed N, short-circuit on first Err                         |
| `tokio::select! { ... }`                  | First-ready wins; **all other branches dropped**            |
| [`JoinSet`](https://docs.rs/tokio/latest/tokio/task/struct.JoinSet.html) | Dynamic N spawned tasks, awaited together   |

**Flag** detached `tokio::spawn` on the request path: the parent does not await the child, so panics, errors, and cancellation are lost. Prefer `JoinSet` so drop aborts the children.

### Cancellation safety in `select!`

A future is cancel-safe if **dropping and recreating it is a no-op**. The [`select!` docs](https://docs.rs/tokio/latest/tokio/macro.select.html) list cancel-safe vs not.

| Cancel-safe                                | NOT cancel-safe — flag in `select!`         |
|--------------------------------------------|---------------------------------------------|
| `mpsc::Receiver::recv`                     | `AsyncReadExt::read_exact`                  |
| `TcpListener::accept`                      | `AsyncWriteExt::write_all` (partial writes) |
| `AsyncReadExt::read` (single call)         | `Mutex::lock`, `Semaphore::acquire` (lose queue position) |
| `StreamExt::next`                          | Anything mid-transaction                    |

```rust
// FLAG: write_all in a select branch can leave a partial frame on the wire
tokio::select! {
    _ = sock.write_all(&frame) => {}
    _ = shutdown.cancelled() => {}
}
```

### Timeouts, deadlines, cancellation token

```rust
use tokio::time::{timeout, Duration};
let r = timeout(Duration::from_secs(2), client.call()).await??;
```

For propagated cancellation across many tasks, use [`CancellationToken`](https://docs.rs/tokio-util/latest/tokio_util/sync/struct.CancellationToken.html). **Flag** request-scoped code that has no deadline and no token — it is a leak waiting to happen.

### Async drop pitfalls

There is no `async Drop`. Resources that need `.await` to release (graceful shutdown, flush, COMMIT) must expose an explicit `async fn close(self)` / `finish(self)`. **Flag** types with a `Drop` impl that block, ignore errors, or rely on a background task to "eventually" clean up.

### Adversarial scan questions — Section 2

- Does any `async fn` call `std::thread::sleep`, `std::fs::*`, `std::net::*`, or do >10 µs of CPU work without `spawn_blocking`?
- Is any `std::sync::Mutex` guard alive across an `.await`?
- Does any `tokio::spawn` capture `Rc`, `RefCell`, or a borrow that compiles only under `current_thread`?
- Is `tokio::spawn` used in detach mode on the request path where `JoinSet` would propagate failures?
- Does any `select!` branch contain a non-cancel-safe future (`read_exact`, `write_all`, `Mutex::lock`, in-flight transaction)?
- Does any I/O-issuing function lack a deadline or `CancellationToken`?
- Does any resource needing async cleanup rely on `Drop`?

---

## Section 3: Unsafe Rust

### When `unsafe` is legitimately needed

Per [std `unsafe` docs](https://doc.rust-lang.org/std/keyword.unsafe.html) and the [Nomicon](https://doc.rust-lang.org/nomicon/):

1. FFI (`extern "C"` calls).
2. Manual memory layout (`MaybeUninit`, raw pointers, custom allocators).
3. Performance hotspots with a measured win **and** a safe wrapper.
4. Building a safe abstraction over a primitive (`Vec`, `Mutex`).

Anything else: **flag**.

### `// SAFETY:` discipline

Every `unsafe { }` block must be preceded by a `// SAFETY:` comment naming the invariants the caller has upheld. Every `unsafe fn` must document a `# Safety` section in its rustdoc.

```rust
/// # Safety
/// `ptr` must be non-null, aligned, and point to an initialized `T`.
unsafe fn read_t<T>(ptr: *const T) -> T {
    // SAFETY: caller upheld the doc-comment preconditions; pointer is valid for reads of T.
    unsafe { ptr.read() }
}
```

### Crate-level lints

```rust
#![forbid(unsafe_code)]                  // crate contains zero unsafe
#![deny(unsafe_op_in_unsafe_fn)]         // unsafe fn body needs explicit unsafe { }
```

See [RFC 2585](https://rust-lang.github.io/rfcs/2585-unsafe-block-in-unsafe-fn.html). **Flag** crates that *should* be `#![forbid(unsafe_code)]` (pure business logic) but are not.

### Footguns to flag on sight

| Pattern                              | Why it's flagworthy                          |
|--------------------------------------|----------------------------------------------|
| `mem::transmute::<A, B>(x)`          | Almost always wrong; prefer `as`, `From`, or `bytemuck` |
| `*const T` / `*mut T` deref          | No lifetime in the pointer — needs SAFETY justification |
| `Box::leak(b)` without "why static"  | Memory leak by design — must be intentional  |
| `Pin::new_unchecked(...)`            | Caller must guarantee no-move forever ([Pin docs](https://doc.rust-lang.org/std/pin/index.html)) |
| `unsafe impl Send for T`             | Hand-asserting thread safety — high stakes   |
| `unsafe impl Sync for T`             | Same                                         |
| `std::hint::unreachable_unchecked()` | UB if hit — `unreachable!()` is almost always correct |

### `unsafe trait` / `unsafe impl`

`unsafe trait` means implementors must uphold invariants the compiler cannot check (`Send`, `Sync`, `GlobalAlloc`, `Allocator`). An `unsafe impl` block is a load-bearing assertion — **always** review the SAFETY justification.

### Miri

Recommend the author run `cargo +nightly miri test` for any crate touching raw pointers, transmute, manual `Send`/`Sync`, or `Pin::new_unchecked`. Miri detects out-of-bounds, use-after-free, aliasing violations (Stacked/Tree Borrows), and uninitialized reads. See [github.com/rust-lang/miri](https://github.com/rust-lang/miri).

### Adversarial scan questions — Section 3

- Does every `unsafe { }` block have a `// SAFETY:` comment that actually names invariants (not "this is safe")?
- Does every `unsafe fn` have a `# Safety` rustdoc section?
- Are there `unsafe impl Send`/`Sync` blocks? Is the justification sound for all fields (including future ones)?
- Is `transmute`, `Box::leak`, or `Pin::new_unchecked` present without a written rationale?
- Should this crate be `#![forbid(unsafe_code)]`?

---

## Section 4: Ownership and Lifetime Traps

### `Arc<Mutex<T>>` overuse

Reach for `Arc<Mutex<T>>` only after answering: can ownership be passed by message instead (channel) or split by partition (sharding)? **Flag** large structs wrapped in a single global `Arc<Mutex<_>>` — that is a contention bottleneck and usually a missing design.

### Lifetime hacks

```rust
// FLAG: lifetime laundering via transmute — almost always unsound
let s: &'static str = unsafe { std::mem::transmute::<&'_ str, &'static str>(local) };

// FLAG: Box::leak for &'static — only OK at startup with a written reason
static CFG: OnceLock<&'static Config> = OnceLock::new();
CFG.get_or_init(|| Box::leak(Box::new(load_config())));   // requires justification
```

### Self-referential structs and `Pin`

A `Pin<&mut T>` in code that is **not** a manually written `Future` impl is suspicious. Pair every `Pin::new_unchecked` with a written argument for why the pointee will not move. Prefer crates like `ouroboros` or `pin-project` over hand-rolled `unsafe`.

### `Rc<RefCell<T>>` inside async

`Rc` is `!Send`. Code holding `Rc<RefCell<_>>` across `.await` will **fail to compile** under `tokio::spawn` (multi-thread runtime) but **will compile** under `current_thread`/`spawn_local`. This is a runtime-flavor surprise: a feature works in tests, breaks in prod. **Flag** any `Rc`/`RefCell` reachable from an `async fn`.

### Borrow-checker cargo-cult

| Smell                                | Better                                         |
|--------------------------------------|------------------------------------------------|
| `.clone()` to silence the borrow checker on a hot path | Restructure scope; borrow shorter |
| `.to_owned()` / `.to_string()` on every call | Accept `&str`; intern; `Cow<'_, str>` |
| `let v = lock.lock().unwrap().clone()` to "drop the guard" | Acknowledge the copy cost or redesign |
| `Arc::clone` everywhere in a single-task path | Use `&Arc<T>` borrow                    |

### Adversarial scan questions — Section 4

- Is `Arc<Mutex<T>>` masking a design that wants a channel or partitioning?
- Is any `'static` widened by `transmute` or `Box::leak`?
- Does any `Pin::new_unchecked` lack a written soundness argument?
- Is `Rc<RefCell<_>>` reachable from async code that may move to a multi-thread runtime?
- Is `.clone()` on a hot path concealing a borrow-checker fight rather than a real ownership requirement?

---

## Section 5: Trait and Generic Discipline

### Return position and parameter position

| Form                  | Dispatch  | When                                                            |
|-----------------------|-----------|-----------------------------------------------------------------|
| `&dyn Trait` arg      | dynamic   | Heterogeneous callers; small code size; no monomorphization     |
| `&impl Trait` arg     | static    | One call site, want inlining                                    |
| `Box<dyn Trait>` ret  | dynamic   | Return type varies at runtime; need owned trait object          |
| `impl Trait` ret      | static    | Single concrete type per call site; do not want to name it      |

### Trait object safety (dyn-compatibility)

A trait is `dyn`-compatible only if (per [reference/items/traits](https://doc.rust-lang.org/reference/items/traits.html#dyn-compatibility)):

- All supertraits are dyn-compatible
- `Self: Sized` is not a supertrait
- No associated constants
- No associated types with generic parameters
- Every method either takes a dispatchable receiver (`&self`, `&mut self`, `Box<Self>`, `Rc<Self>`, `Arc<Self>`, `Pin<P<Self>>`), has no type params, returns no `impl Trait`, and is not `async fn` — **or** carries a `where Self: Sized` bound to opt out

**Flag** a trait used as `dyn Trait` that contains `async fn` (use `#[async_trait]`) or `-> impl Future`.

### `where` clauses vs inline bounds

Prefer `where` for anything beyond a single bound — it keeps the signature readable and surfaces complex bounds in one place.

```rust
// Prefer
fn run<T, S>(t: T, s: S) -> Result<()>
where
    T: Service + Send + 'static,
    S: Store<Item = T::Output>,
{ ... }
```

### Marker traits

`Send`, `Sync`, `Sized`, `Unpin` are auto-implemented when all fields qualify. **Flag** every `unsafe impl Send` / `unsafe impl Sync` — it overrides the auto-derive and must justify why a raw pointer / `Cell` / `*mut` field is actually thread-safe.

### Sealed trait pattern

For public traits that the library must remain free to extend (per [api-guidelines/future-proofing](https://rust-lang.github.io/api-guidelines/future-proofing.html)):

```rust
pub trait Storage: private::Sealed {
    fn get(&self, k: &str) -> Option<&str>;
}

mod private {
    pub trait Sealed {}
    impl Sealed for super::MemStorage {}
}
```

**Flag** a public trait that is *meant* to be closed but is not sealed — adding a method is a breaking change.

### Adversarial scan questions — Section 5

- Is any `dyn Trait` used on a trait with `async fn` or `-> impl Future`?
- Is there an `unsafe impl Send`/`Sync` whose justification covers every field, including ones a future maintainer may add?
- Should this public trait be sealed?
- Is `Box<dyn Trait>` allocated on a hot path where `impl Trait` with monomorphization would inline?

---

## Severity and Category Map

| Idiom / Violation                                                   | Severity | Category enum       |
|---------------------------------------------------------------------|----------|---------------------|
| `.unwrap()` / `panic!()` on a request path                          | P1       | `panicking_api`     |
| `todo!()` / `unimplemented!()` in non-WIP code                      | P1       | `panicking_api`     |
| `expect("...")` without a documented invariant                      | P2       | `panicking_api`     |
| `unreachable!()` used as control flow (not a true impossibility)    | P2       | `panicking_api`     |
| Library public API returning `anyhow::Result`                       | P2       | `errors`            |
| `.map_err(|e| Custom::Msg(e.to_string()))` (lossy error chain)      | P2       | `errors`            |
| Branching on `err.to_string()` instead of variants                  | P1       | `errors`            |
| Missing `#[from]`/`#[source]` on a wrapping variant                 | P3       | `errors`            |
| `context(format!(...))` on a hot path                               | P3       | `errors`            |
| Blocking std I/O / `std::thread::sleep` inside `async fn`           | P1       | `async_blocking`    |
| CPU-bound work in async without `spawn_blocking`                    | P1       | `async_blocking`    |
| `std::sync::MutexGuard` held across `.await`                        | P1       | `async_blocking`    |
| `tokio::spawn` capturing `Rc` / `RefCell` (multi-thread runtime)    | P1       | `async_blocking`    |
| Detached `tokio::spawn` on request path (no `JoinSet`/await)        | P2       | `async_blocking`    |
| Non-cancel-safe future in `select!` (`read_exact`, `write_all`, lock) | P1     | `cancellation`      |
| I/O without deadline / `CancellationToken`                          | P2       | `cancellation`      |
| Async resource cleanup hidden in `Drop`                             | P2       | `cancellation`      |
| `unsafe { }` block without `// SAFETY:` comment                     | P1       | `unsafe_rust`       |
| `unsafe fn` without `# Safety` rustdoc section                      | P1       | `unsafe_rust`       |
| `unsafe impl Send` / `unsafe impl Sync` without justification       | P1       | `unsafe_rust`       |
| `mem::transmute`, `Pin::new_unchecked`, `Box::leak` without rationale | P1     | `unsafe_rust`       |
| Crate could be `#![forbid(unsafe_code)]` but is not                 | P3       | `unsafe_rust`       |
| `'static` widened via `transmute` or `Box::leak`                    | P1       | `unsafe_rust`       |
| `Rc<RefCell<_>>` reachable from `async fn`                          | P2       | `async_blocking`    |
| `dyn Trait` over a trait with `async fn` / `-> impl Future`         | P2       | `errors`            |
| Public unsealed trait meant to be closed                            | P3       | `errors`            |

---

### Sources

- [doc.rust-lang.org — Recoverable Errors with Result](https://doc.rust-lang.org/book/ch09-02-recoverable-errors-with-result.html)
- [docs.rs — thiserror](https://docs.rs/thiserror/latest/thiserror/)
- [docs.rs — anyhow](https://docs.rs/anyhow/latest/anyhow/)
- [doc.rust-lang.org — `unreachable!` macro](https://doc.rust-lang.org/std/macro.unreachable.html)
- [tokio.rs — Shared state tutorial](https://tokio.rs/tokio/tutorial/shared-state)
- [tokio.rs — Spawning tutorial](https://tokio.rs/tokio/tutorial/spawning)
- [docs.rs — `tokio::select!`](https://docs.rs/tokio/latest/tokio/macro.select.html)
- [docs.rs — `tokio::task::spawn_blocking`](https://docs.rs/tokio/latest/tokio/task/fn.spawn_blocking.html)
- [docs.rs — `tokio::task::JoinSet`](https://docs.rs/tokio/latest/tokio/task/struct.JoinSet.html)
- [docs.rs — `tokio_util::sync::CancellationToken`](https://docs.rs/tokio-util/latest/tokio_util/sync/struct.CancellationToken.html)
- [async book — Send approximation](https://rust-lang.github.io/async-book/07_workarounds/03_send_approximation.html)
- [RFC 3185 — `async fn` in trait](https://rust-lang.github.io/rfcs/3185-static-async-fn-in-trait.html)
- [doc.rust-lang.org — `unsafe` keyword](https://doc.rust-lang.org/std/keyword.unsafe.html)
- [doc.rust-lang.org — Nomicon](https://doc.rust-lang.org/nomicon/)
- [doc.rust-lang.org — `std::pin`](https://doc.rust-lang.org/std/pin/index.html)
- [RFC 2585 — unsafe blocks in `unsafe fn`](https://rust-lang.github.io/rfcs/2585-unsafe-block-in-unsafe-fn.html)
- [github.com/rust-lang/miri](https://github.com/rust-lang/miri)
- [doc.rust-lang.org — Trait dyn-compatibility](https://doc.rust-lang.org/reference/items/traits.html#dyn-compatibility)
- [Rust API Guidelines — Sealed traits](https://rust-lang.github.io/api-guidelines/future-proofing.html#sealed-traits-protect-against-downstream-implementations-c-sealed)
