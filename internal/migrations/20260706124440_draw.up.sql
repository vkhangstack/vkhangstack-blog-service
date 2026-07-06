CREATE TABLE IF NOT EXISTS drawing (
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

CREATE INDEX IF NOT EXISTS idx_drawing_owner_id ON drawing(owner_id);
CREATE INDEX IF NOT EXISTS idx_drawing_created_at ON drawing(created_at);