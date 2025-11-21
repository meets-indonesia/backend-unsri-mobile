-- Drop triggers
DROP TRIGGER IF EXISTS update_device_tokens_updated_at ON device_tokens;
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
DROP TRIGGER IF EXISTS update_broadcasts_updated_at ON broadcasts;

-- Drop tables
DROP TABLE IF EXISTS device_tokens;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS broadcast_audiences;
DROP TABLE IF EXISTS broadcasts;

