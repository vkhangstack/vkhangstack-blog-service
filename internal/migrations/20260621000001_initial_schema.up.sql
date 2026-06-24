CREATE TABLE IF NOT EXISTS messages (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL,
    body       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id         BIGINT PRIMARY KEY,
    username   TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    email      TEXT UNIQUE,
    full_name  TEXT NOT NULL,
    role       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS stores (
    id          BIGINT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    address     TEXT NOT NULL,
    image_url   TEXT,
    phone       TEXT NOT NULL,
    email       TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS customers (
    id          BIGINT PRIMARY KEY,
    email       TEXT NOT NULL,
    external_id TEXT,
    password    TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT FALSE,
    balance     DOUBLE PRECISION DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS categories (
    id             BIGINT PRIMARY KEY,
    name           TEXT NOT NULL,
    description    TEXT NOT NULL,
    image_url      TEXT NOT NULL,
    is_active      BOOLEAN NOT NULL DEFAULT FALSE,
    subcategories  JSONB,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS products (
    id          BIGINT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    cost_price  DOUBLE PRECISION NOT NULL,
    sale_price  DOUBLE PRECISION,
    is_active   BOOLEAN NOT NULL DEFAULT FALSE,
    image_url   TEXT,
    category_id BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS staff (
    id                BIGINT PRIMARY KEY,
    full_name         TEXT NOT NULL,
    store_id          BIGINT NOT NULL,
    role              TEXT NOT NULL,
    email             TEXT,
    authentication_id TEXT NOT NULL,
    is_active         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS orders (
    id          BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    store_id    BIGINT NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    status      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS order_items (
    id         BIGINT PRIMARY KEY,
    order_id   BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity   INTEGER NOT NULL,
    unit_price DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS inventories (
    id            BIGINT PRIMARY KEY,
    product_id    BIGINT NOT NULL,
    store_id      BIGINT NOT NULL,
    quantity      INTEGER NOT NULL,
    reorder_level INTEGER NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS payments (
    id          BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    checkout_id TEXT NOT NULL,
    order_id    BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS order_infos (
    id             BIGINT PRIMARY KEY,
    order_id       BIGINT,
    seller_account TEXT,
    amount         TEXT,
    currency       TEXT,
    status         TEXT
);

CREATE TABLE IF NOT EXISTS jti_sessions (
    id      BIGINT PRIMARY KEY,
    user_id TEXT NOT NULL,
    jti     TEXT NOT NULL
);
