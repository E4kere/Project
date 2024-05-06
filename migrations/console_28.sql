-- tokens table
CREATE TABLE tokens (
                        hash BYTEA PRIMARY KEY,
                        user_id INT REFERENCES users(id) ON DELETE CASCADE,
                        expiry TIMESTAMPTZ NOT NULL,
                        scope VARCHAR(20) NOT NULL
);
