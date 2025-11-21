-- Drop triggers
DROP TRIGGER IF EXISTS update_academic_periods_updated_at ON academic_periods;

DROP TRIGGER IF EXISTS update_study_programs_updated_at ON study_programs;

DROP TRIGGER IF EXISTS update_rooms_updated_at ON rooms;

DROP TRIGGER IF EXISTS update_enrollments_updated_at ON enrollments;

DROP TRIGGER IF EXISTS update_classes_updated_at ON classes;

DROP TRIGGER IF EXISTS update_courses_updated_at ON courses;

-- Drop columns added to existing tables
ALTER TABLE dosen
DROP COLUMN IF EXISTS study_program_id;

ALTER TABLE mahasiswa
DROP COLUMN IF EXISTS study_program_id;

ALTER TABLE schedules
DROP COLUMN IF EXISTS class_id;

-- Drop tables in reverse order
DROP TABLE IF EXISTS academic_periods;

DROP TABLE IF EXISTS study_programs;

DROP TABLE IF EXISTS rooms;

DROP TABLE IF EXISTS enrollments;

DROP TABLE IF EXISTS classes;

DROP TABLE IF EXISTS courses;