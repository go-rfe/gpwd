CREATE TABLE IF NOT EXISTS secrets (
    id VARCHAR PRIMARY KEY,
    username VARCHAR REFERENCES accounts(username),
    labels VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN DEFAULT false,
    data bytea
);