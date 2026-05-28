-- Rollback migration 001.
DROP TABLE IF EXISTS checkout.idempotency_keys;
DROP SCHEMA IF EXISTS checkout;
