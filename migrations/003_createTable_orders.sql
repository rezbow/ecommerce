-- +goose Up

CREATE TABLE orders (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users(id), 
	status VARCHAR(50) NOT NULL,
	total_amount BIGINT NOT NULL,
	shipping_address TEXT NOT NULL,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS orders;
