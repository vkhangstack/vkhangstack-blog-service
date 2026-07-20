CREATE TABLE IF NOT EXISTS flashcard_decks (
  id          VARCHAR(20)  PRIMARY KEY,
  title       VARCHAR(255) NOT NULL,
  description TEXT,
  language    VARCHAR(10)  NOT NULL DEFAULT 'en',
  status      VARCHAR(20)  NOT NULL DEFAULT 'draft',
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ,
  created_by  VARCHAR(20),
  updated_by  VARCHAR(20),
  deleted_by  VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS flashcard_cards (
  id         VARCHAR(20) PRIMARY KEY,
  deck_id    VARCHAR(20) NOT NULL REFERENCES flashcard_decks(id) ON DELETE CASCADE,
  front      TEXT        NOT NULL,
  back       TEXT        NOT NULL,
  example    TEXT,
  phonetic   VARCHAR(255),
  position   INTEGER     NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS flashcard_reviews (
  id                VARCHAR(20)   PRIMARY KEY,
  user_id           VARCHAR(20)   NOT NULL,
  card_id           VARCHAR(20)   NOT NULL REFERENCES flashcard_cards(id) ON DELETE CASCADE,
  deck_id           VARCHAR(20)   NOT NULL REFERENCES flashcard_decks(id) ON DELETE CASCADE,
  ease_factor       NUMERIC(4,2)  NOT NULL DEFAULT 2.5,
  interval_days     INTEGER       NOT NULL DEFAULT 0,
  repetitions       INTEGER       NOT NULL DEFAULT 0,
  due_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  last_reviewed_at  TIMESTAMPTZ,
  last_rating       VARCHAR(20),
  created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, card_id)
);

CREATE INDEX IF NOT EXISTS idx_flashcard_decks_status ON flashcard_decks(status);
CREATE INDEX IF NOT EXISTS idx_flashcard_decks_created_at ON flashcard_decks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_flashcard_cards_deck_id ON flashcard_cards(deck_id);
CREATE INDEX IF NOT EXISTS idx_flashcard_reviews_card_id ON flashcard_reviews(card_id);
CREATE INDEX IF NOT EXISTS idx_flashcard_reviews_deck_id ON flashcard_reviews(deck_id);
CREATE INDEX IF NOT EXISTS idx_flashcard_reviews_user_due ON flashcard_reviews(user_id, due_at);
