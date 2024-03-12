CREATE DATABASE gunstore;

CREATE TABLE guns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    damage INT NOT NULL
);
