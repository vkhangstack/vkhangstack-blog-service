CREATE TABLE IF NOT EXISTS quizzes (
  id               VARCHAR(20)   PRIMARY KEY,
  title            VARCHAR(255)  NOT NULL,
  description      TEXT,
  exam_type        VARCHAR(20)   NOT NULL DEFAULT 'general',
  skill            VARCHAR(20),
  status           VARCHAR(20)   NOT NULL DEFAULT 'draft',
  duration_minutes INTEGER,
  created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  deleted_at       TIMESTAMPTZ,
  created_by       VARCHAR(20),
  updated_by       VARCHAR(20),
  deleted_by       VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS quiz_questions (
  id            VARCHAR(20) PRIMARY KEY,
  quiz_id       VARCHAR(20) NOT NULL,
  prompt        TEXT        NOT NULL,
  options       JSONB       NOT NULL,
  correct_index INTEGER     NOT NULL,
  explanation   TEXT,
  position      INTEGER     NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quiz_attempts (
  id               VARCHAR(20) PRIMARY KEY,
  quiz_id          VARCHAR(20) NOT NULL,
  answers          JSONB       NOT NULL,
  score            INTEGER     NOT NULL,
  total            INTEGER     NOT NULL,
  duration_seconds INTEGER,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by       VARCHAR(20)
);

CREATE INDEX IF NOT EXISTS idx_quizzes_status ON quizzes(status);
CREATE INDEX IF NOT EXISTS idx_quizzes_exam_type ON quizzes(exam_type);
CREATE INDEX IF NOT EXISTS idx_quizzes_created_at ON quizzes(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_quiz_id ON quiz_questions(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_quiz_id ON quiz_attempts(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_created_at ON quiz_attempts(created_at DESC);
