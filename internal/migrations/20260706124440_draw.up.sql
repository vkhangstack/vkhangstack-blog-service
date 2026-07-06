CREATE TABLE IF NOT EXISTS drawings (
    id          VARCHAR(20)   PRIMARY KEY,
    owner_id    VARCHAR(20)   NOT NULL,
    title       VARCHAR(255)  NOT NULL DEFAULT 'Untitled drawing',
    elements    JSONB         NOT NULL DEFAULT '[]',
    app_state   JSONB         NOT NULL DEFAULT '{}',
    files       JSONB         NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    created_by  VARCHAR(20)   NOT NULL,
    updated_by  VARCHAR(20)   NULL,
    deleted_by  VARCHAR(20)   NULL
);

CREATE INDEX IF NOT EXISTS idx_drawings_owner_id ON drawings(owner_id);
CREATE INDEX IF NOT EXISTS idx_drawings_created_at ON drawings(created_at);