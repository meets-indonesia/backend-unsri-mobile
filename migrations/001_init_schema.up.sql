-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('mahasiswa', 'dosen', 'staff')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Create mahasiswa table
CREATE TABLE IF NOT EXISTS mahasiswa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nim VARCHAR(50) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    prodi VARCHAR(255),
    angkatan INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_mahasiswa_user_id ON mahasiswa(user_id);
CREATE INDEX idx_mahasiswa_nim ON mahasiswa(nim);
CREATE INDEX idx_mahasiswa_deleted_at ON mahasiswa(deleted_at);

-- Create dosen table
CREATE TABLE IF NOT EXISTS dosen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nip VARCHAR(50) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    prodi VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_dosen_user_id ON dosen(user_id);
CREATE INDEX idx_dosen_nip ON dosen(nip);
CREATE INDEX idx_dosen_deleted_at ON dosen(deleted_at);

-- Create staff table
CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nip VARCHAR(50) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    jabatan VARCHAR(255),
    unit VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_staff_user_id ON staff(user_id);
CREATE INDEX idx_staff_nip ON staff(nip);
CREATE INDEX idx_staff_deleted_at ON staff(deleted_at);

-- Create schedules table (expandable for academic system)
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID, -- For future expansion with course service
    course_code VARCHAR(50),
    course_name VARCHAR(255),
    dosen_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    room VARCHAR(100),
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_schedules_dosen_id ON schedules(dosen_id);
CREATE INDEX idx_schedules_date ON schedules(date);
CREATE INDEX idx_schedules_course_id ON schedules(course_id);
CREATE INDEX idx_schedules_deleted_at ON schedules(deleted_at);

-- Create attendance_sessions table
CREATE TABLE IF NOT EXISTS attendance_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID REFERENCES schedules(id) ON DELETE SET NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('kelas', 'kampus')),
    qr_code TEXT,
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_attendance_sessions_schedule_id ON attendance_sessions(schedule_id);
CREATE INDEX idx_attendance_sessions_created_by ON attendance_sessions(created_by);
CREATE INDEX idx_attendance_sessions_expires_at ON attendance_sessions(expires_at);
CREATE INDEX idx_attendance_sessions_deleted_at ON attendance_sessions(deleted_at);

-- Create attendances table
CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    session_id UUID REFERENCES attendance_sessions(id) ON DELETE SET NULL,
    schedule_id UUID REFERENCES schedules(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('kelas', 'kampus')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('hadir', 'izin', 'sakit', 'alpa', 'terlambat')),
    date DATE NOT NULL,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_attendances_user_id ON attendances(user_id);
CREATE INDEX idx_attendances_session_id ON attendances(session_id);
CREATE INDEX idx_attendances_schedule_id ON attendances(schedule_id);
CREATE INDEX idx_attendances_date ON attendances(date);
CREATE INDEX idx_attendances_type ON attendances(type);
CREATE INDEX idx_attendances_deleted_at ON attendances(deleted_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_mahasiswa_updated_at BEFORE UPDATE ON mahasiswa
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_dosen_updated_at BEFORE UPDATE ON dosen
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_attendance_sessions_updated_at BEFORE UPDATE ON attendance_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_attendances_updated_at BEFORE UPDATE ON attendances
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

