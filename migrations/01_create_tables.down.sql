-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Drop indexes first
DROP INDEX IF EXISTS idx_assets_asset_type;
DROP INDEX IF EXISTS idx_definitions_type;
DROP INDEX IF EXISTS idx_assets_definition_id;
DROP INDEX IF EXISTS idx_assets_account_id;
DROP INDEX IF EXISTS idx_assets_user_id;
DROP INDEX IF EXISTS idx_accounts_user_id;
DROP INDEX IF EXISTS idx_definitions_abbreviation;
DROP INDEX IF EXISTS idx_definitions_name;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS definitions;
DROP TABLE IF EXISTS users; 