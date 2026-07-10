CREATE TABLE IF NOT EXISTS menus (
  id          VARCHAR(20)  PRIMARY KEY,
  group_title VARCHAR(255) NOT NULL,
  parent_id   VARCHAR(20)  REFERENCES menus(id) ON DELETE CASCADE,
  title       VARCHAR(255) NOT NULL,
  url         VARCHAR(500),
  icon        VARCHAR(100),
  badge       VARCHAR(100),
  resource    VARCHAR(255),
  sort_order  INT          NOT NULL DEFAULT 0,
  is_active   BOOLEAN      NOT NULL DEFAULT true,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ,
  deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_menus_group_title ON menus(group_title);
CREATE INDEX IF NOT EXISTS idx_menus_parent_id   ON menus(parent_id);
