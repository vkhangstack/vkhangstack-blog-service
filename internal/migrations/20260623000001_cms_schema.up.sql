CREATE TABLE IF NOT EXISTS blog_tags (
    id         VARCHAR(20) PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS blog_categories (
    id          VARCHAR(20) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    parent_id   VARCHAR(20) NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS blog_posts (
    id              VARCHAR(20) PRIMARY KEY,
    title           VARCHAR(500) NOT NULL,
    slug            VARCHAR(500) NOT NULL UNIQUE,
    excerpt         TEXT,
    content         TEXT NOT NULL,
    cover_image_url VARCHAR(1000),
    category_id     VARCHAR(20) NULL,
    `status`          VARCHAR(50) NOT NULL DEFAULT 'draft',
    published_at    TIMESTAMPTZ,
    scheduled_at    TIMESTAMPTZ,
    view_count      BIGINT NOT NULL DEFAULT 0,
    lexical_state   TEXT NULL,
    author_id       VARCHAR(20) NULL,
    `type`            VARCHAR(50) NOT NULL DEFAULT 'post',
    visibility      VARCHAR(50) NOT NULL DEFAULT 'public',
    locale          VARCHAR(10) NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS blog_post_tags (
    post_id VARCHAR(20) NOT NULL,
    tag_id  VARCHAR(20) NOT NULL,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_blog_posts_status    ON blog_posts(`status`)       WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_posts_scheduled ON blog_posts(scheduled_at) WHERE `status` = 'scheduled';
CREATE INDEX idx_blog_posts_published ON blog_posts(published_at) WHERE `status` = 'published';
CREATE INDEX idx_blog_posts_category  ON blog_posts(category_id)   WHERE deleted_at IS NULL
CREATE INDEX idx_blog_categories_slug ON blog_categories(slug)    WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_post_tags_tag   ON blog_post_tags(tag_id);
