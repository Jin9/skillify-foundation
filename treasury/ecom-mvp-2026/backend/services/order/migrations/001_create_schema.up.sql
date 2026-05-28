-- Migration 001: Create the "order" schema.
-- All tables in this service live under the "order" schema (schema-per-service pattern).

CREATE SCHEMA IF NOT EXISTS "order";

-- Enable uuid-ossp for uuid_generate_v7() if needed (pgx handles UUIDs in application layer).
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
