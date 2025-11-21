-- Migration for location and access control

-- Create geofences table
CREATE TABLE IF NOT EXISTS geofences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius DOUBLE PRECISION NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_geofences_is_active ON geofences(is_active);
CREATE INDEX idx_geofences_deleted_at ON geofences(deleted_at);

-- Create location_history table
CREATE TABLE IF NOT EXISTS location_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('tap_in', 'tap_out')),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    geofence_id UUID REFERENCES geofences(id) ON DELETE SET NULL,
    is_valid BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_location_history_user_id ON location_history(user_id);
CREATE INDEX idx_location_history_geofence_id ON location_history(geofence_id);
CREATE INDEX idx_location_history_created_at ON location_history(created_at);

-- Create access_logs table
CREATE TABLE IF NOT EXISTS access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    gate_id VARCHAR(100),
    access_type VARCHAR(20) NOT NULL CHECK (access_type IN ('entry', 'exit')),
    is_allowed BOOLEAN DEFAULT TRUE,
    reason TEXT,
    qr_code_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_access_logs_user_id ON access_logs(user_id);
CREATE INDEX idx_access_logs_gate_id ON access_logs(gate_id);
CREATE INDEX idx_access_logs_created_at ON access_logs(created_at);

-- Create access_permissions table
CREATE TABLE IF NOT EXISTS access_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    gate_id VARCHAR(100) NOT NULL,
    is_allowed BOOLEAN DEFAULT TRUE,
    valid_from TIMESTAMP,
    valid_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_access_permissions_user_id ON access_permissions(user_id);
CREATE INDEX idx_access_permissions_gate_id ON access_permissions(gate_id);
CREATE INDEX idx_access_permissions_deleted_at ON access_permissions(deleted_at);

-- Create trigger for updated_at
CREATE TRIGGER update_geofences_updated_at BEFORE UPDATE ON geofences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_access_permissions_updated_at BEFORE UPDATE ON access_permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

