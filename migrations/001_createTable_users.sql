
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	is_admin BOOLEAN,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
