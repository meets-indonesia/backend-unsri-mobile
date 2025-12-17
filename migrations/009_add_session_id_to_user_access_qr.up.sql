-- Migration: Add session_id and expires_at to user_access_qrs table
-- This migration adds support for tap-in/tap-out functionality for gate access QR codes
--
-- Changes:
-- 1. Add session_id column (UUID, unique) - unique session ID per QR generation
-- 2. Add expires_at column (TIMESTAMP, nullable) - set when tap-out occurs
-- 3. Create indexes for performance

-- Add session_id column (UUID, unique, not null)
ALTER TABLE user_access_qrs 
ADD COLUMN IF NOT EXISTS session_id UUID;

-- Update existing records: generate session_id from id (for existing data)
UPDATE user_access_qrs 
SET session_id = gen_random_uuid() 
WHERE session_id IS NULL;

-- Make session_id NOT NULL after populating existing records
ALTER TABLE user_access_qrs 
ALTER COLUMN session_id SET NOT NULL;

-- Add unique constraint on session_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_access_qrs_session_id ON user_access_qrs(session_id);

-- Add expires_at column (TIMESTAMP, nullable)
ALTER TABLE user_access_qrs 
ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP;

-- Create index on expires_at for query performance (filtering expired sessions)
CREATE INDEX IF NOT EXISTS idx_user_access_qrs_expires_at ON user_access_qrs(expires_at);

-- Create composite index for common query pattern: (session_id, is_active, expires_at)
CREATE INDEX IF NOT EXISTS idx_user_access_qrs_session_active_expires 
ON user_access_qrs(session_id, is_active, expires_at) 
WHERE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP;
