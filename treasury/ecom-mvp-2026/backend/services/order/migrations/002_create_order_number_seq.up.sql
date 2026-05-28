-- Migration 002: Global order_number sequence (ORD-008).
--
-- A single global sequence is used (not per-day) per td.json §order_number_strategy:
-- "A per-day sequence requires either DDL at midnight or a sequence-per-day naming scheme
-- (operationally noisy in MVP). A single global sequence... keeps allocation race-free."
-- The date prefix is injected at allocation time in Go: fmt.Sprintf("ORD-%s-%06d", date, seq).
-- NNNNNN does not reset to 000001 at midnight — ORD-008 does not require that.

CREATE SEQUENCE IF NOT EXISTS "order".order_number_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
