-- +goose Up
CREATE TABLE products (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	price BIGINT NOT NULL,
	stock_quantity INTEGER NOT NULL,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS products;
