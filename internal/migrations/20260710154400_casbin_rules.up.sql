CREATE TABLE IF NOT EXISTS casbin_rules (
  id     BIGSERIAL    PRIMARY KEY,
  ptype  VARCHAR(10)  NOT NULL,          -- "p" (policy) or "g" (grouping/role)
  v0     VARCHAR(255) NOT NULL DEFAULT '',
  v1     VARCHAR(255) NOT NULL DEFAULT '',
  v2     VARCHAR(255) NOT NULL DEFAULT '',
  v3     VARCHAR(255) NOT NULL DEFAULT '',
  v4     VARCHAR(255) NOT NULL DEFAULT '',
  v5     VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_casbin_rules_unique
  ON casbin_rules(ptype, v0, v1, v2, v3, v4, v5);
