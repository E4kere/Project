CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       name VARCHAR(500) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash BYTEA NOT NULL,
                       activated BOOLEAN NOT NULL DEFAULT FALSE,
                       version INT NOT NULL DEFAULT 1
);
