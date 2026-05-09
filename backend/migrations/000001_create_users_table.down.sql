-- ══════════════════════════════════════════════════════════════════
-- Migration: 000001_create_users_table (DOWN)
-- Description: Drops the users table and all associated objects
-- ══════════════════════════════════════════════════════════════════

DROP TABLE IF EXISTS users;
