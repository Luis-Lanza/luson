-- Drop indexes first
DROP INDEX IF EXISTS idx_suppliers_active;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_hash;

-- Drop tables
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS suppliers;
