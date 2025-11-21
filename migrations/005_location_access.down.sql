-- Drop triggers
DROP TRIGGER IF EXISTS update_access_permissions_updated_at ON access_permissions;
DROP TRIGGER IF EXISTS update_geofences_updated_at ON geofences;

-- Drop tables
DROP TABLE IF EXISTS access_permissions;
DROP TABLE IF EXISTS access_logs;
DROP TABLE IF EXISTS location_history;
DROP TABLE IF EXISTS geofences;

