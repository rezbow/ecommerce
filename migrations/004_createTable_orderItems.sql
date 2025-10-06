-- +goose Up

CREATE TABLE order_items (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	order_id UUID NOT NULL REFERENCES orders(id),
	product_id UUID NOT NULL REFERENCES products(id),
	quantity INTEGER NOT NULL,
	unit_price BIGINT NOT NULL,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

-- +goose Down

DROP TABLE IF EXISTS order_items;
