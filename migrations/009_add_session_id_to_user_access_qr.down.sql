-- Rollback migration: Remove session_id and expires_at from user_access_qrs table

-- Drop composite index
DROP INDEX IF EXISTS idx_user_access_qrs_session_active_expires;

-- Drop expires_at index
DROP INDEX IF EXISTS idx_user_access_qrs_expires_at;

-- Drop session_id unique index
DROP INDEX IF EXISTS idx_user_access_qrs_session_id;

-- Drop expires_at column
ALTER TABLE user_access_qrs 
DROP COLUMN IF EXISTS expires_at;

-- Drop session_id column
ALTER TABLE user_access_qrs 
DROP COLUMN IF EXISTS session_id;
