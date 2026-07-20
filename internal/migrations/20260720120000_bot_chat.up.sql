CREATE TABLE IF NOT EXISTS bot_chat_messages (
  id          VARCHAR(20)  PRIMARY KEY,
  chat_id     VARCHAR(255) NOT NULL,
  -- message_id is the Zalo webhook's message id, used as an idempotency key so a redelivered
  -- webhook never records the same inbound message (or its auto-reply) twice. Empty string =
  -- no id (outbound bot/admin rows), excluded from the unique index via the partial predicate.
  message_id  VARCHAR(255) NOT NULL DEFAULT '',
  user_id     VARCHAR(20),
  sender_name VARCHAR(255),
  direction   VARCHAR(10)  NOT NULL,
  sender      VARCHAR(10)  NOT NULL,
  text        TEXT         NOT NULL,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bot_chat_messages_chat_id ON bot_chat_messages(chat_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bot_chat_messages_user_id ON bot_chat_messages(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_bot_chat_messages_message_id
  ON bot_chat_messages (message_id) WHERE message_id <> '';
