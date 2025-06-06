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
			id, user_id, name, account_type, balance, currency_code, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID,
		a.UserID,
		a.Name,
		a.AccountType,
		a.Balance,
		a.CurrencyCode,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, balance, currency_code, created_at, updated_at
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

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, balance, currency_code, created_at, updated_at
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

func (r *PostgresRepository) GetByType(ctx context.Context, userID uuid.UUID, accountType account.AccountType) ([]*account.Account, error) {
	query := `
		SELECT id, user_id, name, account_type, balance, currency_code, created_at, updated_at
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
		SET name = $1, account_type = $2, balance = $3, currency_code = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		a.Name,
		a.AccountType,
		a.Balance,
		a.CurrencyCode,
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
	// Get total accounts and balance
	totalQuery := `
		SELECT COUNT(*) as total_accounts, COALESCE(SUM(balance), 0) as total_balance
		FROM accounts
		WHERE user_id = $1
	`
	var totalAccounts int
	var totalBalance float64
	err := r.db.QueryRowContext(ctx, totalQuery, userID).Scan(&totalAccounts, &totalBalance)
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
	typeRows, err := r.db.QueryContext(ctx, typeQuery, userID)
	if err != nil {
		return nil, err
	}
	defer typeRows.Close()

	byType := make(map[account.AccountType]int)
	for typeRows.Next() {
		var accountType account.AccountType
		var count int
		if err := typeRows.Scan(&accountType, &count); err != nil {
			return nil, err
		}
		byType[accountType] = count
	}

	// Get balance by currency
	currencyQuery := `
		SELECT currency_code, SUM(balance) as total_balance
		FROM accounts
		WHERE user_id = $1
		GROUP BY currency_code
	`
	currencyRows, err := r.db.QueryContext(ctx, currencyQuery, userID)
	if err != nil {
		return nil, err
	}
	defer currencyRows.Close()

	byCurrency := make(map[string]float64)
	for currencyRows.Next() {
		var currencyCode string
		var balance float64
		if err := currencyRows.Scan(&currencyCode, &balance); err != nil {
			return nil, err
		}
		byCurrency[currencyCode] = balance
	}

	return &account.AccountSummary{
		TotalAccounts: totalAccounts,
		TotalBalance:  totalBalance,
		ByType:        byType,
		ByCurrency:    byCurrency,
	}, nil
}

func (r *PostgresRepository) Filter(ctx context.Context, query account.FilterAccountsQuery) ([]*account.Account, error) {
	baseQuery := `
		SELECT id, user_id, name, account_type, balance, currency_code, created_at, updated_at
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

	if query.CurrencyCode != nil {
		conditions = append(conditions, fmt.Sprintf("currency_code = $%d", argIndex))
		args = append(args, *query.CurrencyCode)
		argIndex++
	}

	if query.MinBalance != nil {
		conditions = append(conditions, fmt.Sprintf("balance >= $%d", argIndex))
		args = append(args, *query.MinBalance)
		argIndex++
	}

	if query.MaxBalance != nil {
		conditions = append(conditions, fmt.Sprintf("balance <= $%d", argIndex))
		args = append(args, *query.MaxBalance)
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
