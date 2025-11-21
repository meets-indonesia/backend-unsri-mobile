-- Drop triggers
DROP TRIGGER IF EXISTS update_attendances_updated_at ON attendances;
DROP TRIGGER IF EXISTS update_attendance_sessions_updated_at ON attendance_sessions;
DROP TRIGGER IF EXISTS update_schedules_updated_at ON schedules;
DROP TRIGGER IF EXISTS update_staff_updated_at ON staff;
DROP TRIGGER IF EXISTS update_dosen_updated_at ON dosen;
DROP TRIGGER IF EXISTS update_mahasiswa_updated_at ON mahasiswa;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS attendances;
DROP TABLE IF EXISTS attendance_sessions;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS dosen;
DROP TABLE IF EXISTS mahasiswa;
DROP TABLE IF EXISTS users;

