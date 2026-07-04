DROP INDEX IF EXISTS idx_notes_status;
DROP INDEX IF EXISTS idx_notes_created_at;
DROP INDEX IF EXISTS idx_note_tags_note_id;
DROP INDEX IF EXISTS idx_note_tags_tag_id;
DROP TABLE IF EXISTS note_tags;
DROP TABLE IF EXISTS notes;