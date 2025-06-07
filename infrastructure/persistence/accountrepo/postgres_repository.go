package accountrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"siyahsensei/wallet-service/domain/account"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, a *account.Account) error {
	query := `
		INSERT INTO accounts (
			id, user_id, name, account_type, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID,
		a.UserID,
		a.Name,
		a.AccountType,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`
	var a account.Account
	err := r.db.GetContext(ctx, &a, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *PostgresRepository) GetByIDWithAssets(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*account.AccountWithAssets, error) {
	// First get the account
	acc, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if account belongs to user
	if acc.UserID != userID {
		return nil, errors.New("unauthorized: account does not belong to user")
	}

	// Get assets for this account with definition details
	assetsQuery := `
		SELECT 
			a.id, a.definition_id, a.asset_type, a.quantity, a.updated_at,
			d.name, d.abbreviation, d.suffix
		FROM assets a
		JOIN definitions d ON a.definition_id = d.id
		WHERE a.account_id = $1
		ORDER BY a.updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, assetsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []account.AssetInfo
	assetCounts := make(map[string]int)
	var lastUpdated *time.Time

	for rows.Next() {
		var asset account.AssetInfo
		var suffix sql.NullString

		err := rows.Scan(
			&asset.ID, &asset.DefinitionID, &asset.Type, &asset.Quantity, &asset.UpdatedAt,
			&asset.Name, &asset.Symbol, &suffix,
		)
		if err != nil {
			return nil, err
		}

		if suffix.Valid {
			asset.Currency = suffix.String
		} else {
			asset.Currency = asset.Symbol // fallback to symbol if no suffix
		}

		assets = append(assets, asset)

		// Count assets by type
		assetCounts[asset.Type]++

		if lastUpdated == nil || asset.UpdatedAt.After(*lastUpdated) {
			lastUpdated = &asset.UpdatedAt
		}
	}

	return &account.AccountWithAssets{
		Account:     *acc,
		Assets:      assets,
		AssetCounts: assetCounts,
		LastUpdated: lastUpdated,
	}, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	var accounts []*account.Account
	err := r.db.SelectContext(ctx, &accounts, query, userID)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *PostgresRepository) GetByUserIDWithAssets(ctx context.Context, userID uuid.UUID) ([]*account.AccountWithAssets, error) {
	accounts, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var accountsWithAssets []*account.AccountWithAssets
	for _, acc := range accounts {
		accWithAssets, err := r.GetByIDWithAssets(ctx, acc.ID, userID)
		if err != nil {
			return nil, err
		}
		accountsWithAssets = append(accountsWithAssets, accWithAssets)
	}

	return accountsWithAssets, nil
}

func (r *PostgresRepository) GetByType(ctx context.Context, userID uuid.UUID, accountType account.AccountType) ([]*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, created_at, updated_at
		FROM accounts
		WHERE user_id = $1 AND account_type = $2
		ORDER BY created_at DESC
	`
	var accounts []*account.Account
	err := r.db.SelectContext(ctx, &accounts, query, userID, accountType)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *PostgresRepository) GetByCurrency(ctx context.Context, userID uuid.UUID, currencyCode string) ([]*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, balance, currency_code, created_at, updated_at
		FROM accounts
		WHERE user_id = $1 AND currency_code = $2
		ORDER BY created_at DESC
	`
	var accounts []*account.Account
	err := r.db.SelectContext(ctx, &accounts, query, userID, currencyCode)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *PostgresRepository) Update(ctx context.Context, a *account.Account) error {
	a.UpdatedAt = time.Now()
	query := `
		UPDATE accounts
		SET name = $1, account_type = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		a.Name,
		a.AccountType,
		a.UpdatedAt,
		a.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *PostgresRepository) UpdateBalance(ctx context.Context, id uuid.UUID, balance float64) error {
	query := `
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := r.db.ExecContext(ctx, query, balance, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *PostgresRepository) GetAccountSummary(ctx context.Context, userID uuid.UUID) (*account.AccountSummary, error) {
	// Get total accounts
	totalQuery := `
		SELECT COUNT(*) as total_accounts
		FROM accounts
		WHERE user_id = $1
	`
	var totalAccounts int
	err := r.db.QueryRowContext(ctx, totalQuery, userID).Scan(&totalAccounts)
	if err != nil {
		return nil, err
	}

	// Get accounts by type
	typeQuery := `
		SELECT account_type, COUNT(*) as count
		FROM accounts
		WHERE user_id = $1
		GROUP BY account_type
	`
	rows, err := r.db.QueryContext(ctx, typeQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byType := make(map[account.AccountType]int)
	for rows.Next() {
		var accountType account.AccountType
		var count int
		err := rows.Scan(&accountType, &count)
		if err != nil {
			return nil, err
		}
		byType[accountType] = count
	}

	// Get currency totals from assets
	currencyQuery := `
		SELECT d.suffix, SUM(a.quantity) as total
		FROM assets a
		JOIN definitions d ON a.definition_id = d.id
		JOIN accounts acc ON a.account_id = acc.id
		WHERE acc.user_id = $1 AND d.suffix IS NOT NULL
		GROUP BY d.suffix
	`
	currencyRows, err := r.db.QueryContext(ctx, currencyQuery, userID)
	if err != nil {
		return nil, err
	}
	defer currencyRows.Close()

	byCurrency := make(map[string]float64)
	for currencyRows.Next() {
		var currency string
		var total float64
		err := currencyRows.Scan(&currency, &total)
		if err != nil {
			return nil, err
		}
		byCurrency[currency] = total
	}

	return &account.AccountSummary{
		TotalAccounts: totalAccounts,
		ByType:        byType,
		ByCurrency:    byCurrency,
	}, nil
}

func (r *PostgresRepository) Filter(ctx context.Context, query account.FilterAccountsQuery) ([]*account.Account, error) {
	baseQuery := `
		SELECT id, user_id, name, account_type, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
	`

	var conditions []string
	var args []interface{}
	args = append(args, query.UserID)
	argIndex := 2

	if query.AccountType != nil {
		conditions = append(conditions, fmt.Sprintf("account_type = $%d", argIndex))
		args = append(args, *query.AccountType)
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY created_at DESC"

	if query.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, query.Limit)
		argIndex++
	}

	if query.Offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, query.Offset)
	}

	var accounts []*account.Account
	err := r.db.SelectContext(ctx, &accounts, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
