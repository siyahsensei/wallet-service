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

-- Create definitions table
CREATE TABLE definitions (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    abbreviation VARCHAR(50) UNIQUE NOT NULL,
    suffix VARCHAR(10),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create index for faster searches
CREATE INDEX idx_definitions_name ON definitions(name);
CREATE INDEX idx_definitions_abbreviation ON definitions(abbreviation);

-- Create accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0.00,
    currency_code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create assets table
CREATE TABLE assets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    definition_id UUID NOT NULL REFERENCES definitions(id) ON DELETE RESTRICT,
    type VARCHAR(50) NOT NULL,
    quantity DECIMAL(20,8) NOT NULL,
    notes TEXT,
    purchase_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create indexes for better performance
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_assets_user_id ON assets(user_id);
CREATE INDEX idx_assets_account_id ON assets(account_id);
CREATE INDEX idx_assets_definition_id ON assets(definition_id);
CREATE INDEX idx_assets_type ON assets(type);
CREATE INDEX idx_assets_purchase_date ON assets(purchase_date);

-- Insert some common currency definitions
INSERT INTO definitions (id, name, abbreviation, suffix, created_at, updated_at) VALUES
    (gen_random_uuid(), 'Turkish Lira', 'TL', '₺', NOW(), NOW()),
    (gen_random_uuid(), 'US Dollar', 'USD', '$', NOW(), NOW()),
    (gen_random_uuid(), 'Euro', 'EUR', '€', NOW(), NOW()),
    (gen_random_uuid(), 'Bitcoin', 'BTC', '₿', NOW(), NOW()),
    (gen_random_uuid(), 'Ethereum', 'ETH', 'Ξ', NOW(), NOW()),
    (gen_random_uuid(), 'Tether USD', 'USDT', '$', NOW(), NOW()),
    (gen_random_uuid(), 'British Pound', 'GBP', '£', NOW(), NOW()),
    (gen_random_uuid(), 'Japanese Yen', 'JPY', '¥', NOW(), NOW()),
    (gen_random_uuid(), 'Swiss Franc', 'CHF', 'Fr', NOW(), NOW()),
    (gen_random_uuid(), 'Canadian Dollar', 'CAD', 'C$', NOW(), NOW()); 