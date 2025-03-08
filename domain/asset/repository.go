package asset

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Create adds a new asset to the database
	Create(ctx context.Context, asset *Asset) error

	// GetByID retrieves an asset by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Asset, error)

	// GetByUserID retrieves all assets for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Asset, error)

	// GetByAccountID retrieves all assets for a specific account
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*Asset, error)

	// GetByType retrieves all assets of a specific type for a user
	GetByType(ctx context.Context, userID uuid.UUID, assetType AssetType) ([]*Asset, error)

	// Update updates an asset's information
	Update(ctx context.Context, asset *Asset) error

	// Delete removes an asset from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// GetTotalValue calculates the total current value of assets for a user
	// or for specific asset types if provided
	GetTotalValue(ctx context.Context, userID uuid.UUID, assetTypes []AssetType) (float64, error)

	// GetAssetPerformance gets the performance of assets over a time period
	GetAssetPerformance(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*AssetPerformance, error)
}
