# System Performance

## Contents
- SLO-driven design
- Caching strategy
- Database performance
- Concurrency & resource management
- Network & latency
- Profiling approach
- Capacity planning

## SLO-driven design

- Define performance targets BEFORE optimizing: p50, p95, p99 latency, throughput (RPS), error budget.
- Every service MUST have SLOs. No SLO means no way to know if it is "fast enough".
- Measure from the caller's perspective, not internal processing time.

## Caching strategy

| Layer | Tool | Invalidation | Use when |
|---|---|---|---|
| Application | In-process (LRU, `sync.Map`) | TTL or event-driven | Hot path, low-cardinality, read-heavy |
| Distributed | Redis, Memcached | TTL + explicit invalidation on write | Cross-instance, session, computed results |
| CDN / Edge | CloudFront, Fastly | TTL + purge API | Static assets, public API responses |

Rules:
- NEVER cache without an invalidation strategy — stale data is a silent bug.
- Cache-aside pattern as default; write-through only when consistency is critical.
- Tag cache entries with a version/hash for safe rollouts.

## Database performance

- Index strategy: cover queries, not tables. Analyze EXPLAIN before adding indexes.
- Connection pooling: size pool to `(cores × 2) + spindle_count`, not arbitrary large numbers.
- Read replicas for read-heavy paths. MUST handle replication lag explicitly.
- Batch writes over individual inserts — chunked upserts (align with staging→merge pattern).
- NEVER `SELECT *` in production code — explicit column lists only.

## Concurrency & resource management

- Bounded concurrency: worker pools / semaphores. NEVER unbounded goroutines or threads.
- Backpressure: reject or queue at the system boundary when capacity is exceeded. Do not degrade silently.
- Circuit breaker on external dependencies — fail fast over cascading timeouts.
- Timeout budget: allocate per hop; total request timeout = sum of critical-path hops + buffer.

## Network & latency

- Minimize serialization overhead: prefer binary (protobuf, msgpack) for service-to-service; JSON at the system boundary only.
- Connection reuse (HTTP/2, gRPC persistent connections). NEVER create connections per request.
- Colocate services that chat frequently — cross-AZ latency adds up.
- Batch external API calls where possible. N+1 at network level is worse than N+1 at DB level.

## Profiling approach

- Profile in production-like environments. Local profiling catches CPU issues but misses infrastructure bottlenecks.
- Go: `pprof` (CPU, heap, goroutine), continuous profiling via Pyroscope or Datadog.
- Flame graphs first, line-level optimization second.
- MUST identify the bottleneck before optimizing. Measure, do not guess.

## Capacity planning

- Load test against 2× expected peak — headroom for spikes.
- Autoscaling: scale on the metric that saturates first (CPU, memory, connections, queue depth).
- Right-size instances: profile actual resource usage, not theoretical maximums.
- Cost-aware: performance optimization MUST justify its infrastructure cost. Scaling horizontally is usually cheaper than micro-optimizing.
