CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE orders (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        user_id UUID NOT NULL,
                        total_price DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),
                        delivery_price DECIMAL(10, 2) NOT NULL CHECK (delivery_price >= 0),
                        currency TEXT NOT NULL,
                        status TEXT NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             product_id UUID NOT NULL,
                             quantity INT NOT NULL CHECK (quantity > 0),
                             unit_price DECIMAL(10, 2) NOT NULL CHECK (unit_price >= 0)
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
