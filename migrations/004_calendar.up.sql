-- Migration for academic calendar

-- Create academic_events table
CREATE TABLE IF NOT EXISTS academic_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    is_all_day BOOLEAN DEFAULT FALSE,
    location VARCHAR(255),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_academic_events_start_date ON academic_events(start_date);
CREATE INDEX idx_academic_events_end_date ON academic_events(end_date);
CREATE INDEX idx_academic_events_event_type ON academic_events(event_type);
CREATE INDEX idx_academic_events_created_by ON academic_events(created_by);
CREATE INDEX idx_academic_events_deleted_at ON academic_events(deleted_at);

-- Create trigger for updated_at
CREATE TRIGGER update_academic_events_updated_at BEFORE UPDATE ON academic_events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

