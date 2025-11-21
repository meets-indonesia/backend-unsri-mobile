-- Migration for broadcast and notification tables

-- Create broadcasts table
CREATE TABLE IF NOT EXISTS broadcasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('general', 'class', 'campus')),
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP,
    scheduled_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_broadcasts_type ON broadcasts(type);
CREATE INDEX idx_broadcasts_created_by ON broadcasts(created_by);
CREATE INDEX idx_broadcasts_class_id ON broadcasts(class_id);
CREATE INDEX idx_broadcasts_is_published ON broadcasts(is_published);
CREATE INDEX idx_broadcasts_scheduled_at ON broadcasts(scheduled_at);
CREATE INDEX idx_broadcasts_deleted_at ON broadcasts(deleted_at);

-- Create broadcast_audiences table
CREATE TABLE IF NOT EXISTS broadcast_audiences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    broadcast_id UUID NOT NULL REFERENCES broadcasts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20),
    prodi VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_broadcast_audiences_broadcast_id ON broadcast_audiences(broadcast_id);
CREATE INDEX idx_broadcast_audiences_user_id ON broadcast_audiences(user_id);
CREATE INDEX idx_broadcast_audiences_role ON broadcast_audiences(role);
CREATE INDEX idx_broadcast_audiences_prodi ON broadcast_audiences(prodi);

-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('info', 'warning', 'error', 'success')),
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_deleted_at ON notifications(deleted_at);

-- Create device_tokens table
CREATE TABLE IF NOT EXISTS device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    platform VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(token)
);

CREATE INDEX idx_device_tokens_user_id ON device_tokens(user_id);
CREATE INDEX idx_device_tokens_is_active ON device_tokens(is_active);
CREATE INDEX idx_device_tokens_deleted_at ON device_tokens(deleted_at);

-- Create triggers for updated_at
CREATE TRIGGER update_broadcasts_updated_at BEFORE UPDATE ON broadcasts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_device_tokens_updated_at BEFORE UPDATE ON device_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

