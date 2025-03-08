-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    institution VARCHAR(100),
    currency VARCHAR(10) NOT NULL,
    balance DECIMAL(16, 2) NOT NULL DEFAULT 0,
    is_manual BOOLEAN NOT NULL DEFAULT TRUE,
    icon VARCHAR(50),
    color VARCHAR(20),
    api_key VARCHAR(255),
    api_secret VARCHAR(255),
    is_connected BOOLEAN NOT NULL DEFAULT FALSE,
    last_sync TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create assets table
CREATE TABLE assets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    symbol VARCHAR(50),
    quantity DECIMAL(20, 8) NOT NULL,
    purchase_price DECIMAL(20, 8) NOT NULL,
    current_price DECIMAL(20, 8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    notes TEXT,
    purchase_date TIMESTAMP NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    asset_id UUID REFERENCES assets(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    quantity DECIMAL(20, 8),
    price DECIMAL(20, 8),
    fee DECIMAL(20, 8) NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    date TIMESTAMP NOT NULL,
    to_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    transaction_hash VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create indexes for better performance
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_assets_user_id ON assets(user_id);
CREATE INDEX idx_assets_account_id ON assets(account_id);
CREATE INDEX idx_assets_type ON assets(type);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_asset_id ON transactions(asset_id);
CREATE INDEX idx_transactions_date ON transactions(date);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_category ON transactions(category);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users; 