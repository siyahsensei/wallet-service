package accountrepo

import (
	"context"
	"database/sql"
	"errors"
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
			id, user_id, name, description, type, institution, currency, balance, 
			is_manual, icon, color, api_key, api_secret, is_connected, last_sync, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID,
		a.UserID,
		a.Name,
		a.Description,
		a.Type,
		a.Institution,
		a.Currency,
		a.Balance,
		a.IsManual,
		a.Icon,
		a.Color,
		a.APIKey,
		a.APISecret,
		a.IsConnected,
		a.LastSync,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	query := `
		SELECT id, user_id, name, description, type, institution, currency, balance, 
		       is_manual, icon, color, api_key, api_secret, is_connected, last_sync, 
		       created_at, updated_at
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
		SELECT id, user_id, name, description, type, institution, currency, balance, 
		       is_manual, icon, color, api_key, api_secret, is_connected, last_sync, 
		       created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY name
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
		SELECT id, user_id, name, description, type, institution, currency, balance, 
		       is_manual, icon, color, api_key, api_secret, is_connected, last_sync, 
		       created_at, updated_at
		FROM accounts
		WHERE user_id = $1 AND type = $2
		ORDER BY name
	`
	var accounts []*account.Account
	err := r.db.SelectContext(ctx, &accounts, query, userID, accountType)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *PostgresRepository) Update(ctx context.Context, a *account.Account) error {
	a.UpdatedAt = time.Now()
	query := `
		UPDATE accounts
		SET name = $1, description = $2, type = $3, institution = $4, currency = $5,
		    balance = $6, is_manual = $7, icon = $8, color = $9, api_key = $10,
		    api_secret = $11, is_connected = $12, last_sync = $13, updated_at = $14
		WHERE id = $15
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		a.Name,
		a.Description,
		a.Type,
		a.Institution,
		a.Currency,
		a.Balance,
		a.IsManual,
		a.Icon,
		a.Color,
		a.APIKey,
		a.APISecret,
		a.IsConnected,
		a.LastSync,
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

func (r *PostgresRepository) GetTotalBalance(ctx context.Context, userID uuid.UUID, accountTypes []account.AccountType) (float64, error) {
	var query string
	var args []interface{}

	if len(accountTypes) == 0 {
		query = `
			SELECT COALESCE(SUM(balance), 0) as total_balance
			FROM accounts
			WHERE user_id = $1
		`
		args = []interface{}{userID}
	} else {
		query = `
			SELECT COALESCE(SUM(balance), 0) as total_balance
			FROM accounts
			WHERE user_id = $1 AND type = ANY($2)
		`
		strTypes := make([]string, len(accountTypes))
		for i, t := range accountTypes {
			strTypes[i] = string(t)
		}
		args = []interface{}{userID, strTypes}
	}
	var totalBalance float64
	err := r.db.GetContext(ctx, &totalBalance, query, args...)
	if err != nil {
		return 0, err
	}
	return totalBalance, nil
}
