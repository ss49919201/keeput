CREATE TABLE entries (
	id              INTEGER PRIMARY KEY AUTOINCREMENT, 
    title           TEXT NOT NULL,
    published_at    INTEGER NOT NULL,
    platform        TEXT CHECK(platform = 'hatena' OR platform = 'zenn') NOT NULL,
    UNIQUE (platform, published_at, title)
);
