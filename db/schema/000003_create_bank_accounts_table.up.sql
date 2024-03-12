CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY,
    bank_name VARCHAR NOT NULL,
    bank_account_name VARCHAR NOT NULL,
    bank_account_number VARCHAR NOT NULL,
    user_id UUID REFERENCES users(id)
);