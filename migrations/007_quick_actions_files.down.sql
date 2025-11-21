-- Drop triggers
DROP TRIGGER IF EXISTS update_files_updated_at ON files;
DROP TRIGGER IF EXISTS update_bimbingans_updated_at ON bimbingans;
DROP TRIGGER IF EXISTS update_krss_updated_at ON krss;
DROP TRIGGER IF EXISTS update_transcripts_updated_at ON transcripts;

-- Drop tables
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS bimbingans;
DROP TABLE IF EXISTS krss;
DROP TABLE IF EXISTS transcripts;

