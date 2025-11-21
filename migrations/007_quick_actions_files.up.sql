-- Migration for quick actions and file storage

-- Create transcripts table
CREATE TABLE IF NOT EXISTS transcripts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    semester VARCHAR(20) NOT NULL,
    ipk DECIMAL(3,2),
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transcripts_student_id ON transcripts(student_id);
CREATE INDEX idx_transcripts_semester ON transcripts(semester);

-- Create krss table
CREATE TABLE IF NOT EXISTS krss (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    semester VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'approved', 'rejected')),
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_krss_student_id ON krss(student_id);
CREATE INDEX idx_krss_semester ON krss(semester);
CREATE INDEX idx_krss_status ON krss(status);

-- Create bimbingans table
CREATE TABLE IF NOT EXISTS bimbingans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    dosen_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    topic VARCHAR(255),
    notes TEXT,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bimbingans_student_id ON bimbingans(student_id);
CREATE INDEX idx_bimbingans_dosen_id ON bimbingans(dosen_id);
CREATE INDEX idx_bimbingans_date ON bimbingans(date);

-- Create files table
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50),
    mime_type VARCHAR(100),
    size BIGINT NOT NULL,
    path VARCHAR(500) NOT NULL,
    url VARCHAR(500),
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_file_type ON files(file_type);
CREATE INDEX idx_files_deleted_at ON files(deleted_at);

-- Create triggers for updated_at
CREATE TRIGGER update_transcripts_updated_at BEFORE UPDATE ON transcripts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_krss_updated_at BEFORE UPDATE ON krss
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bimbingans_updated_at BEFORE UPDATE ON bimbingans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_files_updated_at BEFORE UPDATE ON files
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

