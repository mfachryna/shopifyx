CREATE TABLE payments (
    id UUID PRIMARY KEY,
    payment_proof_image_url VARCHAR NOT NULL,
    quantity INTEGER  NOT NULL,
    product_id UUID REFERENCES products(id),
    user_id UUID REFERENCES users(id),
    bank_account_id UUID REFERENCES bank_accounts(id)
);