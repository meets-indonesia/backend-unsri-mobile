-- Migration untuk ekspansi sistem akademik
-- Tabel-tabel ini siap untuk digunakan ketika course service dan academic services lainnya diimplementasikan

-- Create courses/mata_kuliah table
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    credits INTEGER NOT NULL DEFAULT 0,
    semester INTEGER,
    prodi VARCHAR(255),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_courses_code ON courses(code);
CREATE INDEX idx_courses_prodi ON courses(prodi);
CREATE INDEX idx_courses_deleted_at ON courses(deleted_at);

-- Create classes table (kelas untuk setiap mata kuliah)
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    class_code VARCHAR(50) NOT NULL,
    class_name VARCHAR(255),
    semester VARCHAR(20) NOT NULL, -- e.g., "2024/2025-Ganjil"
    academic_year VARCHAR(20),
    capacity INTEGER DEFAULT 0,
    enrolled INTEGER DEFAULT 0,
    dosen_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assistant_dosen_id UUID REFERENCES users(id) ON DELETE SET NULL,
    room VARCHAR(100),
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(course_id, class_code, semester)
);

CREATE INDEX idx_classes_course_id ON classes(course_id);
CREATE INDEX idx_classes_dosen_id ON classes(dosen_id);
CREATE INDEX idx_classes_semester ON classes(semester);
CREATE INDEX idx_classes_deleted_at ON classes(deleted_at);

-- Create enrollments table (KRS - Kartu Rencana Studi)
CREATE TABLE IF NOT EXISTS enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE RESTRICT,
    enrollment_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'dropped', 'failed')),
    grade VARCHAR(5), -- e.g., 'A', 'B', 'C', 'D', 'E'
    score DECIMAL(5,2), -- Numeric score
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(student_id, class_id)
);

CREATE INDEX idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX idx_enrollments_class_id ON enrollments(class_id);
CREATE INDEX idx_enrollments_status ON enrollments(status);
CREATE INDEX idx_enrollments_deleted_at ON enrollments(deleted_at);

-- Update schedules table to reference classes
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS class_id UUID REFERENCES classes(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_schedules_class_id ON schedules(class_id);

-- Create rooms table (untuk manajemen ruangan)
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    building VARCHAR(255),
    floor INTEGER,
    capacity INTEGER,
    room_type VARCHAR(50), -- 'classroom', 'lab', 'auditorium', etc.
    facilities TEXT, -- JSON array of facilities
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_rooms_code ON rooms(code);
CREATE INDEX idx_rooms_building ON rooms(building);
CREATE INDEX idx_rooms_deleted_at ON rooms(deleted_at);

-- Create study_programs table (Program Studi)
CREATE TABLE IF NOT EXISTS study_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    faculty VARCHAR(255),
    degree_level VARCHAR(50), -- 'S1', 'S2', 'S3', 'D3', etc.
    accreditation VARCHAR(10),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_study_programs_code ON study_programs(code);
CREATE INDEX idx_study_programs_faculty ON study_programs(faculty);
CREATE INDEX idx_study_programs_deleted_at ON study_programs(deleted_at);

-- Update mahasiswa table to reference study_programs
ALTER TABLE mahasiswa ADD COLUMN IF NOT EXISTS study_program_id UUID REFERENCES study_programs(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_mahasiswa_study_program_id ON mahasiswa(study_program_id);

-- Update dosen table to reference study_programs
ALTER TABLE dosen ADD COLUMN IF NOT EXISTS study_program_id UUID REFERENCES study_programs(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_dosen_study_program_id ON dosen(study_program_id);

-- Create academic_periods table (Semester)
CREATE TABLE IF NOT EXISTS academic_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL, -- e.g., "2024/2025-Ganjil"
    academic_year VARCHAR(20) NOT NULL,
    semester_type VARCHAR(20) NOT NULL CHECK (semester_type IN ('Ganjil', 'Genap', 'Pendek')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    registration_start DATE,
    registration_end DATE,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_academic_periods_code ON academic_periods(code);
CREATE INDEX idx_academic_periods_academic_year ON academic_periods(academic_year);
CREATE INDEX idx_academic_periods_is_active ON academic_periods(is_active);
CREATE INDEX idx_academic_periods_deleted_at ON academic_periods(deleted_at);

-- Create triggers for updated_at
CREATE TRIGGER update_courses_updated_at BEFORE UPDATE ON courses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_classes_updated_at BEFORE UPDATE ON classes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_enrollments_updated_at BEFORE UPDATE ON enrollments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rooms_updated_at BEFORE UPDATE ON rooms
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_study_programs_updated_at BEFORE UPDATE ON study_programs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_academic_periods_updated_at BEFORE UPDATE ON academic_periods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

