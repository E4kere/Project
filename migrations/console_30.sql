CREATE TABLE guns (
                      id SERIAL PRIMARY KEY,
                      name VARCHAR(50) NOT NULL,
                      price DECIMAL(10, 2) NOT NULL,
                      damage INT NOT NULL,
                      created_at TIMESTAMPTZ DEFAULT NOW(),
                      updated_at TIMESTAMPTZ DEFAULT NOW()
);
