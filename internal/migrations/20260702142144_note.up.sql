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
  created_by  VARCHAR(36),
  updated_by  VARCHAR(36),
  deleted_by  VARCHAR(36)
);

CREATE TABLE IF NOT EXISTS note_tags (
  note_id  VARCHAR(20) NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
  tag_id   VARCHAR(20) NOT NULL REFERENCES blog_tags(id) ON DELETE CASCADE,
  PRIMARY KEY (note_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_notes_status ON notes(status);
CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);