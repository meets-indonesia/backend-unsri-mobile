-- Migration for user access QR (unique QR per user for gate access)

-- Create user_access_qrs table
CREATE TABLE IF NOT EXISTS user_access_qrs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    qr_token VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id),
    UNIQUE(qr_token)
);

CREATE INDEX idx_user_access_qrs_user_id ON user_access_qrs(user_id);
CREATE INDEX idx_user_access_qrs_qr_token ON user_access_qrs(qr_token);
CREATE INDEX idx_user_access_qrs_is_active ON user_access_qrs(is_active);
CREATE INDEX idx_user_access_qrs_deleted_at ON user_access_qrs(deleted_at);

-- Create trigger for updated_at
CREATE TRIGGER update_user_access_qrs_updated_at BEFORE UPDATE ON user_access_qrs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

