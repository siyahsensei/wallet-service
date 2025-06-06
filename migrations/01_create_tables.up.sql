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

-- Create definitions table (asset details like currency codes, crypto symbols, stock tickers)
CREATE TABLE definitions (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    abbreviation VARCHAR(50) UNIQUE NOT NULL,
    suffix VARCHAR(10),
    definition_type VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create accounts table (where assets are stored)
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

-- Create assets table (what assets user owns)
CREATE TABLE assets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    definition_id UUID NOT NULL REFERENCES definitions(id) ON DELETE RESTRICT,
    asset_type VARCHAR(50) NOT NULL,
    quantity DECIMAL(20,8) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create indexes for better performance
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_assets_user_id ON assets(user_id);
CREATE INDEX idx_assets_account_id ON assets(account_id);
CREATE INDEX idx_assets_definition_id ON assets(definition_id);
CREATE INDEX idx_assets_asset_type ON assets(asset_type);
CREATE INDEX idx_definitions_type ON definitions(definition_type);
CREATE INDEX idx_definitions_name ON definitions(name);
CREATE INDEX idx_definitions_abbreviation ON definitions(abbreviation);

INSERT INTO definitions (id, name, abbreviation, suffix, definition_type, description, created_at, updated_at) VALUES
-- Currencies (CURRENCY) - Top 10 Most Traded
(gen_random_uuid(), 'US Dollar', 'USD', '$', 'CURRENCY', 'United States Dollar', NOW(), NOW()),
(gen_random_uuid(), 'Euro', 'EUR', '€', 'CURRENCY', 'Currency of the European Union', NOW(), NOW()),
(gen_random_uuid(), 'Japanese Yen', 'JPY', '¥', 'CURRENCY', 'National currency of Japan', NOW(), NOW()),
(gen_random_uuid(), 'British Pound', 'GBP', '£', 'CURRENCY', 'British Pound Sterling', NOW(), NOW()),
(gen_random_uuid(), 'Chinese Yuan', 'CNY', '¥', 'CURRENCY', 'National currency of China', NOW(), NOW()),
(gen_random_uuid(), 'Australian Dollar', 'AUD', 'A$', 'CURRENCY', 'National currency of Australia', NOW(), NOW()),
(gen_random_uuid(), 'Canadian Dollar', 'CAD', 'C$', 'CURRENCY', 'National currency of Canada', NOW(), NOW()),
(gen_random_uuid(), 'Swiss Franc', 'CHF', 'Fr', 'CURRENCY', 'National currency of Switzerland', NOW(), NOW()),
(gen_random_uuid(), 'Hong Kong Dollar', 'HKD', 'HK$', 'CURRENCY', 'Currency of Hong Kong', NOW(), NOW()),
(gen_random_uuid(), 'Turkish Lira', 'TRY', '₺', 'CURRENCY', 'National currency of Turkey', NOW(), NOW()),

-- Cryptocurrencies (CRYPTOCURRENCY) - Top 20 by Market Cap
(gen_random_uuid(), 'Bitcoin', 'BTC', '₿', 'CRYPTOCURRENCY', 'The first and largest cryptocurrency', NOW(), NOW()),
(gen_random_uuid(), 'Ethereum', 'ETH', 'Ξ', 'CRYPTOCURRENCY', 'Smart contract platform cryptocurrency', NOW(), NOW()),
(gen_random_uuid(), 'Tether USDt', 'USDT', '₮', 'CRYPTOCURRENCY', 'A stablecoin pegged to the US Dollar', NOW(), NOW()),
(gen_random_uuid(), 'BNB', 'BNB', 'BNB', 'CRYPTOCURRENCY', 'The cryptocurrency of the Binance exchange', NOW(), NOW()),
(gen_random_uuid(), 'Solana', 'SOL', 'SOL', 'CRYPTOCURRENCY', 'A high-performance blockchain platform', NOW(), NOW()),
(gen_random_uuid(), 'USDC', 'USDC', '$', 'CRYPTOCURRENCY', 'A stablecoin pegged to the US Dollar', NOW(), NOW()),
(gen_random_uuid(), 'XRP', 'XRP', 'XRP', 'CRYPTOCURRENCY', 'A digital payment protocol developed by Ripple', NOW(), NOW()),
(gen_random_uuid(), 'Dogecoin', 'DOGE', 'Ð', 'CRYPTOCURRENCY', 'A meme-based cryptocurrency derived from Litecoin', NOW(), NOW()),
(gen_random_uuid(), 'Toncoin', 'TON', 'TON', 'CRYPTOCURRENCY', 'A blockchain originally developed by Telegram', NOW(), NOW()),
(gen_random_uuid(), 'Cardano', 'ADA', '₳', 'CRYPTOCURRENCY', 'A proof-of-stake blockchain platform', NOW(), NOW()),
(gen_random_uuid(), 'Shiba Inu', 'SHIB', 'SHIB', 'CRYPTOCURRENCY', 'A token created as an alternative to Dogecoin', NOW(), NOW()),
(gen_random_uuid(), 'Avalanche', 'AVAX', 'AVAX', 'CRYPTOCURRENCY', 'A platform for decentralized applications and blockchains', NOW(), NOW()),
(gen_random_uuid(), 'Polkadot', 'DOT', 'DOT', 'CRYPTOCURRENCY', 'A protocol that enables different blockchains to interoperate', NOW(), NOW()),
(gen_random_uuid(), 'Chainlink', 'LINK', 'LINK', 'CRYPTOCURRENCY', 'A decentralized oracle network', NOW(), NOW()),
(gen_random_uuid(), 'TRON', 'TRX', 'TRX', 'CRYPTOCURRENCY', 'A decentralized content sharing platform', NOW(), NOW()),
(gen_random_uuid(), 'Polygon', 'MATIC', 'MATIC', 'CRYPTOCURRENCY', 'A scaling solution for Ethereum', NOW(), NOW()),
(gen_random_uuid(), 'Litecoin', 'LTC', 'Ł', 'CRYPTOCURRENCY', 'One of the first cryptocurrencies created as an alternative to Bitcoin', NOW(), NOW()),
(gen_random_uuid(), 'Internet Computer', 'ICP', 'ICP', 'CRYPTOCURRENCY', 'A decentralized computing platform that runs at web speed', NOW(), NOW()),
(gen_random_uuid(), 'Kaspa', 'KAS', 'KAS', 'CRYPTOCURRENCY', 'A Proof-of-Work cryptocurrency using the BlockDAG structure', NOW(), NOW()),
(gen_random_uuid(), 'Ethereum Classic', 'ETC', 'ETC', 'CRYPTOCURRENCY', 'The continuation of the original Ethereum blockchain', NOW(), NOW()),

-- Commodities (COMMODITY) - Top 25
-- Precious Metals
(gen_random_uuid(), 'Gold', 'XAU', 'oz', 'COMMODITY', 'Precious metal gold (per ounce)', NOW(), NOW()),
(gen_random_uuid(), 'Silver', 'XAG', 'oz', 'COMMODITY', 'Precious metal silver (per ounce)', NOW(), NOW()),
(gen_random_uuid(), 'Platinum', 'XPT', 'oz', 'COMMODITY', 'Precious metal platinum (per ounce)', NOW(), NOW()),
(gen_random_uuid(), 'Palladium', 'XPD', 'oz', 'COMMODITY', 'Precious metal palladium (per ounce)', NOW(), NOW()),
-- Industrial Metals
(gen_random_uuid(), 'Copper', 'HG', 'lb', 'COMMODITY', 'Industrial metal copper (per pound)', NOW(), NOW()),
(gen_random_uuid(), 'Aluminum', 'ALI', 't', 'COMMODITY', 'Industrial metal aluminum (per ton)', NOW(), NOW()),
(gen_random_uuid(), 'Iron Ore', 'TIO', 't', 'COMMODITY', 'Industrial metal iron ore (per ton)', NOW(), NOW()),
(gen_random_uuid(), 'Nickel', 'NI', 't', 'COMMODITY', 'Industrial metal nickel (per ton)', NOW(), NOW()),
(gen_random_uuid(), 'Zinc', 'ZN', 't', 'COMMODITY', 'Industrial metal zinc (per ton)', NOW(), NOW()),
(gen_random_uuid(), 'Lead', 'PB', 't', 'COMMODITY', 'Industrial metal lead (per ton)', NOW(), NOW()),
-- Energy
(gen_random_uuid(), 'Crude Oil (WTI)', 'WTI', 'bbl', 'COMMODITY', 'West Texas Intermediate crude oil (per barrel)', NOW(), NOW()),
(gen_random_uuid(), 'Brent Crude Oil', 'BRENT', 'bbl', 'COMMODITY', 'Brent crude oil (per barrel)', NOW(), NOW()),
(gen_random_uuid(), 'Natural Gas', 'NG', 'MMBtu', 'COMMODITY', 'Natural Gas (per MMBtu)', NOW(), NOW()),
(gen_random_uuid(), 'Gasoline', 'RB', 'gal', 'COMMODITY', 'Gasoline (per gallon)', NOW(), NOW()),
(gen_random_uuid(), 'Heating Oil', 'HO', 'gal', 'COMMODITY', 'Heating Oil (per gallon)', NOW(), NOW()),
-- Agricultural Products
(gen_random_uuid(), 'Corn', 'ZC', 'bu', 'COMMODITY', 'Corn (per bushel)', NOW(), NOW()),
(gen_random_uuid(), 'Soybeans', 'ZS', 'bu', 'COMMODITY', 'Soybeans (per bushel)', NOW(), NOW()),
(gen_random_uuid(), 'Wheat', 'ZW', 'bu', 'COMMODITY', 'Wheat (per bushel)', NOW(), NOW()),
(gen_random_uuid(), 'Coffee', 'KC', 'lb', 'COMMODITY', 'Coffee (per pound)', NOW(), NOW()),
(gen_random_uuid(), 'Sugar', 'SB', 'lb', 'COMMODITY', 'Sugar (per pound)', NOW(), NOW()),
(gen_random_uuid(), 'Cotton', 'CT', 'lb', 'COMMODITY', 'Cotton (per pound)', NOW(), NOW()),
(gen_random_uuid(), 'Cocoa', 'CC', 't', 'COMMODITY', 'Cocoa (per ton)', NOW(), NOW()),
(gen_random_uuid(), 'Lumber', 'LBS', '1kbf', 'COMMODITY', 'Lumber (per 1000 board feet)', NOW(), NOW()),
(gen_random_uuid(), 'Live Cattle', 'LE', 'lb', 'COMMODITY', 'Live Cattle (per pound)', NOW(), NOW()),
(gen_random_uuid(), 'Lean Hogs', 'HE', 'lb', 'COMMODITY', 'Lean Hogs (per pound)', NOW(), NOW());