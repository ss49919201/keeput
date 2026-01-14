CREATE TABLE entries (
    id              SERIAL,
    title           VARCHAR(100) NOT NULL,
    published_at    DATETIME NOT NULL,
    platform        ENUM('zenn', 'hatena'),
    fetched_at      DATETIME NOT NULL,
    UNIQUE (platform, published_at, title)
);
