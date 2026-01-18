CREATE TYPE platform AS ENUM ('zenn', 'hatena');

CREATE TABLE entries (
    id              SERIAL PRIMARY KEY,
    title           CHARACTER VARYING(100) NOT NULL,
    published_at    TIMESTAMPTZ NOT NULL,
    platform        platform NOT NULL,
    UNIQUE (platform, published_at, title)
);
