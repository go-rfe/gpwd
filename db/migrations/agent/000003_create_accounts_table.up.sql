CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR PRIMARY KEY,
    server TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password BLOB NOT NULL,
    registered BOOLEAN DEFAULT false NOT NULL,
    UNIQUE (server, username)
);