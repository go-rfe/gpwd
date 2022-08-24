CREATE TABLE IF NOT EXISTS secrets (
    id VARCHAR PRIMARY KEY,
    labels TEXT,
    data BLOB
);