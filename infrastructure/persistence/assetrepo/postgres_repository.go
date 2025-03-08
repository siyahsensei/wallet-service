package assetrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"siyahsensei/wallet-service/domain/asset"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, a *asset.Asset) error {
	query := `
		INSERT INTO assets (
			id, user_id, account_id, name, type, symbol, quantity, purchase_price,
			current_price, currency, notes, purchase_date, last_updated, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID,
		a.UserID,
		a.AccountID,
		a.Name,
		a.Type,
		a.Symbol,
		a.Quantity,
		a.PurchasePrice,
		a.CurrentPrice,
		a.Currency,
		a.Notes,
		a.PurchaseDate,
		a.LastUpdated,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*asset.Asset, error) {
	query := `
		SELECT id, user_id, account_id, name, type, symbol, quantity, purchase_price,
		       current_price, currency, notes, purchase_date, last_updated, created_at, updated_at
		FROM assets
		WHERE id = $1
	`
	var a asset.Asset
	err := r.db.GetContext(ctx, &a, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*asset.Asset, error) {
	query := `
		SELECT id, user_id, account_id, name, type, symbol, quantity, purchase_price,
		       current_price, currency, notes, purchase_date, last_updated, created_at, updated_at
		FROM assets
		WHERE user_id = $1
		ORDER BY name
	`
	var assets []*asset.Asset
	err := r.db.SelectContext(ctx, &assets, query, userID)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *PostgresRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*asset.Asset, error) {
	query := `
		SELECT id, user_id, account_id, name, type, symbol, quantity, purchase_price,
		       current_price, currency, notes, purchase_date, last_updated, created_at, updated_at
		FROM assets
		WHERE account_id = $1
		ORDER BY name
	`
	var assets []*asset.Asset
	err := r.db.SelectContext(ctx, &assets, query, accountID)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *PostgresRepository) GetByType(ctx context.Context, userID uuid.UUID, assetType asset.AssetType) ([]*asset.Asset, error) {
	query := `
		SELECT id, user_id, account_id, name, type, symbol, quantity, purchase_price,
		       current_price, currency, notes, purchase_date, last_updated, created_at, updated_at
		FROM assets
		WHERE user_id = $1 AND type = $2
		ORDER BY name
	`
	var assets []*asset.Asset
	err := r.db.SelectContext(ctx, &assets, query, userID, assetType)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *PostgresRepository) Update(ctx context.Context, a *asset.Asset) error {
	a.UpdatedAt = time.Now()
	query := `
		UPDATE assets
		SET name = $1, type = $2, symbol = $3, quantity = $4, purchase_price = $5,
		    current_price = $6, currency = $7, notes = $8, purchase_date = $9,
		    last_updated = $10, updated_at = $11
		WHERE id = $12
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		a.Name,
		a.Type,
		a.Symbol,
		a.Quantity,
		a.PurchasePrice,
		a.CurrentPrice,
		a.Currency,
		a.Notes,
		a.PurchaseDate,
		a.LastUpdated,
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
		return errors.New("asset not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM assets
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
		return errors.New("asset not found")
	}
	return nil
}

func (r *PostgresRepository) GetTotalValue(ctx context.Context, userID uuid.UUID, assetTypes []asset.AssetType) (float64, error) {
	var query string
	var args []interface{}

	if len(assetTypes) == 0 {
		query = `
			SELECT COALESCE(SUM(quantity * current_price), 0) as total_value
			FROM assets
			WHERE user_id = $1
		`
		args = []interface{}{userID}
	} else {
		query = `
			SELECT COALESCE(SUM(quantity * current_price), 0) as total_value
			FROM assets
			WHERE user_id = $1 AND type = ANY($2)
		`
		strTypes := make([]string, len(assetTypes))
		for i, t := range assetTypes {
			strTypes[i] = string(t)
		}
		args = []interface{}{userID, strTypes}
	}

	var totalValue float64
	err := r.db.GetContext(ctx, &totalValue, query, args...)
	if err != nil {
		return 0, err
	}

	return totalValue, nil
}

func (r *PostgresRepository) GetAssetPerformance(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*asset.AssetPerformance, error) {
	query := `
		SELECT 
			id, name, symbol, type, 
			(quantity * purchase_price) as initial_value,
			(quantity * current_price) as current_value,
			(quantity * current_price - quantity * purchase_price) as profit_loss,
			CASE 
				WHEN purchase_price = 0 THEN 0
				ELSE ((current_price - purchase_price) / purchase_price) * 100 
			END as profit_loss_percentage,
			currency
		FROM assets
		WHERE user_id = $1
		ORDER BY name
	`
	var performances []*asset.AssetPerformance
	err := r.db.SelectContext(ctx, &performances, query, userID)
	if err != nil {
		return nil, err
	}
	return performances, nil
}
