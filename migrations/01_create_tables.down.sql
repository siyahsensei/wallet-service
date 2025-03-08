-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users; 