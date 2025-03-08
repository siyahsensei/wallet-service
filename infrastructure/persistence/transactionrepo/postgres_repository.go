package transactionrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"siyahsensei/wallet-service/domain/transaction"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, t *transaction.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
			currency, description, category, date, to_account_id, transaction_hash,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		t.ID,
		t.UserID,
		t.AccountID,
		t.AssetID,
		t.Type,
		t.Amount,
		t.Quantity,
		t.Price,
		t.Fee,
		t.Currency,
		t.Description,
		t.Category,
		t.Date,
		t.ToAccountID,
		t.TransactionHash,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE id = $1
	`
	var t transaction.Transaction
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`
	var transactions []*transaction.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE account_id = $1 OR to_account_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`
	var transactions []*transaction.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE asset_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`
	var transactions []*transaction.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, assetID, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
		LIMIT $4 OFFSET $5
	`
	var transactions []*transaction.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresRepository) GetByType(ctx context.Context, userID uuid.UUID, transactionType transaction.TransactionType, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND type = $2
		ORDER BY date DESC
		LIMIT $3 OFFSET $4
	`
	var transactions []*transaction.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID, transactionType, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresRepository) GetByCategory(ctx context.Context, userID uuid.UUID, category string, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, asset_id, type, amount, quantity, price, fee,
		       currency, description, category, date, to_account_id, transaction_hash,
		       created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND category = $2
		ORDER BY date DESC
		LIMIT $3 OFFSET $4
	`
	var transactions []*transaction.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID, category, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresRepository) Update(ctx context.Context, t *transaction.Transaction) error {
	t.UpdatedAt = time.Now()
	query := `
		UPDATE transactions
		SET account_id = $1, asset_id = $2, type = $3, amount = $4, quantity = $5,
		    price = $6, fee = $7, currency = $8, description = $9, category = $10,
		    date = $11, to_account_id = $12, transaction_hash = $13, updated_at = $14
		WHERE id = $15
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		t.AccountID,
		t.AssetID,
		t.Type,
		t.Amount,
		t.Quantity,
		t.Price,
		t.Fee,
		t.Currency,
		t.Description,
		t.Category,
		t.Date,
		t.ToAccountID,
		t.TransactionHash,
		t.UpdatedAt,
		t.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM transactions
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
		return errors.New("transaction not found")
	}
	return nil
}

func (r *PostgresRepository) GetTotalsByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]float64, error) {
	query := `
		SELECT category, SUM(
			CASE
				WHEN type IN ('WITHDRAWAL', 'BUY', 'TRANSFER', 'FEE', 'EXPENSE', 'TAX', 'BORROWING', 'REPAYMENT', 'LENDING') THEN -amount
				ELSE amount
			END
		) as total
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY category
	`
	rows, err := r.db.QueryxContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	totals := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		totals[category] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return totals, nil
}

func (r *PostgresRepository) GetTotalsByType(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[transaction.TransactionType]float64, error) {
	query := `
		SELECT type, SUM(
			CASE
				WHEN type IN ('WITHDRAWAL', 'BUY', 'TRANSFER', 'FEE', 'EXPENSE', 'TAX', 'BORROWING', 'REPAYMENT', 'LENDING') THEN -amount
				ELSE amount
			END
		) as total
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY type
	`
	rows, err := r.db.QueryxContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	totals := make(map[transaction.TransactionType]float64)
	for rows.Next() {
		var txType transaction.TransactionType
		var total float64
		if err := rows.Scan(&txType, &total); err != nil {
			return nil, err
		}
		totals[txType] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return totals, nil
}

func (r *PostgresRepository) GetMonthlyTotals(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*transaction.MonthlyTotal, error) {
	query := `
		SELECT 
			EXTRACT(YEAR FROM date) as year,
			EXTRACT(MONTH FROM date) as month,
			SUM(CASE WHEN type IN ('DEPOSIT', 'SELL', 'DIVIDEND', 'INTEREST', 'INCOME', 'STAKING', 'MINING', 'AIRDROP', 'BORROWING') THEN amount ELSE 0 END) as total_in,
			SUM(CASE WHEN type IN ('WITHDRAWAL', 'BUY', 'TRANSFER', 'FEE', 'EXPENSE', 'TAX', 'BORROWING', 'REPAYMENT', 'LENDING') THEN amount ELSE 0 END) as total_out,
			SUM(CASE 
				WHEN type IN ('DEPOSIT', 'SELL', 'DIVIDEND', 'INTEREST', 'INCOME', 'STAKING', 'MINING', 'AIRDROP', 'BORROWING') THEN amount 
				WHEN type IN ('WITHDRAWAL', 'BUY', 'TRANSFER', 'FEE', 'EXPENSE', 'TAX', 'BORROWING', 'REPAYMENT', 'LENDING') THEN -amount
				ELSE 0 
			END) as net_amount
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY year, month
		ORDER BY year, month
	`
	var monthlyTotals []*transaction.MonthlyTotal
	err := r.db.SelectContext(ctx, &monthlyTotals, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return monthlyTotals, nil
}
