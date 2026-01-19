CREATE TABLE entries (
	id              INTEGER PRIMARY KEY AUTOINCREMENT, 
    title           TEXT NOT NULL,
    published_at    INTEGER NOT NULL,
    platform        TEXT NOT NULL,
    UNIQUE (platform, published_at, title)
);
