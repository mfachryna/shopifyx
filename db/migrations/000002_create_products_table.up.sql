CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    price INTEGER  NOT NULL,
    image_url VARCHAR  NOT NULL,
    stock INTEGER  NOT NULL,
    condition VARCHAR(50)  NOT NULL,
    is_purchasable BOOLEAN  NOT NULL,
    purchase_count INTEGER  DEFAULT 0,
    tags VARCHAR[],
    user_id UUID REFERENCES users(id)
);