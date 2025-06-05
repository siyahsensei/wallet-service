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
			id, user_id, account_id, definition_id, type, quantity, notes, purchase_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID,
		a.UserID,
		a.AccountID,
		a.DefinitionID,
		a.Type,
		a.Quantity,
		a.Notes,
		a.PurchaseDate,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*asset.Asset, error) {
	query := `
		SELECT id, user_id, account_id, definition_id, type, quantity, notes, purchase_date, created_at, updated_at
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
		SELECT id, user_id, account_id, definition_id, type, quantity, notes, purchase_date, created_at, updated_at
		FROM assets
		WHERE user_id = $1
		ORDER BY created_at DESC
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
		SELECT id, user_id, account_id, definition_id, type, quantity, notes, purchase_date, created_at, updated_at
		FROM assets
		WHERE account_id = $1
		ORDER BY created_at DESC
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
		SELECT id, user_id, account_id, definition_id, type, quantity, notes, purchase_date, created_at, updated_at
		FROM assets
		WHERE user_id = $1 AND type = $2
		ORDER BY created_at DESC
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
		SET account_id = $1, definition_id = $2, type = $3, quantity = $4, notes = $5, purchase_date = $6, updated_at = $7
		WHERE id = $8
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		a.AccountID,
		a.DefinitionID,
		a.Type,
		a.Quantity,
		a.Notes,
		a.PurchaseDate,
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
		// Get total value for all asset types
		query = `
			SELECT COALESCE(SUM(quantity), 0) as total_value
			FROM assets
			WHERE user_id = $1
		`
		args = []interface{}{userID}
	} else {
		// Get total value for specific asset types
		query = `
			SELECT COALESCE(SUM(quantity), 0) as total_value
			FROM assets
			WHERE user_id = $1 AND type = ANY($2)
		`
		// Convert AssetType slice to string slice for PostgreSQL array
		typeStrings := make([]string, len(assetTypes))
		for i, t := range assetTypes {
			typeStrings[i] = string(t)
		}
		args = []interface{}{userID, typeStrings}
	}

	var totalValue float64
	err := r.db.GetContext(ctx, &totalValue, query, args...)
	if err != nil {
		return 0, err
	}
	return totalValue, nil
}

func (r *PostgresRepository) GetAssetPerformance(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*asset.AssetPerformance, error) {
	// This is a simplified implementation. In a real-world scenario, you would need
	// additional tables for price history, market data, etc.
	query := `
		SELECT 
			a.id as asset_id,
			d.name,
			d.abbreviation as symbol,
			a.type,
			a.quantity as initial_value,
			a.quantity as current_value,
			0 as profit_loss,
			0 as profit_loss_perc,
			d.abbreviation as currency
		FROM assets a
		JOIN definitions d ON a.definition_id = d.id
		WHERE a.user_id = $1 
		AND a.purchase_date BETWEEN $2 AND $3
		ORDER BY a.purchase_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*asset.AssetPerformance
	for rows.Next() {
		var perf asset.AssetPerformance
		err := rows.Scan(
			&perf.AssetID,
			&perf.Name,
			&perf.Symbol,
			&perf.Type,
			&perf.InitialValue,
			&perf.CurrentValue,
			&perf.ProfitLoss,
			&perf.ProfitLossPerc,
			&perf.Currency,
		)
		if err != nil {
			return nil, err
		}
		performances = append(performances, &perf)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return performances, nil
}
