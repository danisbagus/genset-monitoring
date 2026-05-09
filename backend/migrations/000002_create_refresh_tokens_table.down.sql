-- ══════════════════════════════════════════════════════════════════
-- Migration: 000002_create_refresh_tokens_table (DOWN)
-- Description: Drops the refresh_tokens table
-- ══════════════════════════════════════════════════════════════════

DROP TABLE IF EXISTS refresh_tokens;
