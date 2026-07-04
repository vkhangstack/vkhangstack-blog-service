CREATE TABLE IF NOT EXISTS tags (
    id         VARCHAR(20) PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    type       VARCHAR(50) NOT NULL,
    author_id  VARCHAR(20) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_tags_slug ON tags(slug);

