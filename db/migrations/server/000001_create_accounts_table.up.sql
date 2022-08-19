CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    username VARCHAR (50) UNIQUE,
    password VARCHAR (255)
);