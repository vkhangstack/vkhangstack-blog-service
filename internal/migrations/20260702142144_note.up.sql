CREATE TABLE IF NOT EXISTS notes (
  id          VARCHAR(20)   PRIMARY KEY,
  title       VARCHAR(255)  NOT NULL,
  source_url  VARCHAR(2048)[],
  status      VARCHAR(50)   NOT NULL DEFAULT 'draft',
  html        TEXT,
  lexical     TEXT,
  description TEXT,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ,
  created_by  VARCHAR(20),
  updated_by  VARCHAR(20),
  deleted_by  VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS note_tags (
  note_id  VARCHAR(20) NOT NULL,
  tag_id   VARCHAR(20) NOT NULL,
  PRIMARY KEY (note_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_notes_status ON notes(status);
CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_note_tags_note_id ON note_tags(note_id);
CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id ON note_tags(tag_id);