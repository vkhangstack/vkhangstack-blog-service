CREATE TABLE IF NOT EXISTS timetable_entries (
  id          VARCHAR(20)  PRIMARY KEY,
  author_id   VARCHAR(255) NOT NULL,
  subject     VARCHAR(255) NOT NULL,
  day_of_week SMALLINT     NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
  start_time  TIME         NOT NULL,
  end_time    TIME         NOT NULL CHECK (end_time > start_time),
  room        VARCHAR(255),
  teacher     VARCHAR(255),
  color       VARCHAR(20)  NOT NULL DEFAULT 'blue',
  note        TEXT,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ,
  deleted_at  TIMESTAMPTZ,
  created_by  VARCHAR(255) NOT NULL,
  updated_by  VARCHAR(255) NULL,
  deleted_by  VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_timetable_entries_day_of_week ON timetable_entries(day_of_week);
CREATE INDEX IF NOT EXISTS idx_timetable_entries_author_id ON timetable_entries(author_id);