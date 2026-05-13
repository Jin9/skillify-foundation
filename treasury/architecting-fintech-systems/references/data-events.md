# Data & Event Patterns

## Contents
- Ingestion: staging → merge (default)
- When to pick something else
- Kafka / event strategy

## Ingestion: staging → merge (default)

Standard flow for batch/bulk data:

1. Ingest into a staging table (chunked, e.g. 10k rows).
2. Validate and clean in staging.
3. Merge into master via upsert:

```sql
INSERT INTO master SELECT ... FROM staging
ON DUPLICATE KEY UPDATE ...
```

Use this pattern by default when any of the following apply: multi-source input, high volume, partial failure recovery matters, or lock contention is a concern.

## When to pick something else

| Pattern | Use when |
|---|---|
| Truncate / rename (swap table) | Full reload, no partial-update need, near-zero stale-read tolerance |
| Upsert in place (no staging) | Low volume, simple schema, single source of truth |
| CDC (Debezium / DMS) | Real-time replication, cross-system sync, event-sourcing feed |

## Kafka / event strategy

- Track consumer lag and offset control.
- Large payloads: split, store externally (S3), reference in event, reassemble on consume.
- NEVER produce unbounded payloads directly to Kafka topics.
- Dead-letter queue for poison messages, with alerting.
