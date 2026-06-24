CREATE TABLE blog_tags (
    id         BIGINT PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE blog_categories (
    id          BIGINT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    parent_id   BIGINT REFERENCES blog_categories(id) ON DELETE SET NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE blog_posts (
    id              BIGINT PRIMARY KEY,
    title           VARCHAR(500) NOT NULL,
    slug            VARCHAR(500) NOT NULL UNIQUE,
    excerpt         TEXT,
    content         TEXT NOT NULL,
    cover_image_url VARCHAR(1000),
    category_id     BIGINT REFERENCES blog_categories(id) ON DELETE SET NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'draft',
    published_at    TIMESTAMPTZ,
    scheduled_at    TIMESTAMPTZ,
    view_count      BIGINT NOT NULL DEFAULT 0,
    author_id       BIGINT NOT NULL REFERENCES accounts(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE blog_post_tags (
    post_id BIGINT NOT NULL REFERENCES blog_posts(id) ON DELETE CASCADE,
    tag_id  BIGINT NOT NULL REFERENCES blog_tags(id)       ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_blog_posts_status    ON blog_posts(status)       WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_posts_scheduled ON blog_posts(scheduled_at) WHERE status = 'scheduled';
CREATE INDEX idx_blog_categories_slug ON blog_categories(slug)    WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_post_tags_tag   ON blog_post_tags(tag_id);
