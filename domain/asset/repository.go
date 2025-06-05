package asset

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, asset *Asset) error
	GetByID(ctx context.Context, id uuid.UUID) (*Asset, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Asset, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*Asset, error)
	GetByType(ctx context.Context, userID uuid.UUID, assetType AssetType) ([]*Asset, error)
	Update(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalValue(ctx context.Context, userID uuid.UUID, assetTypes []AssetType) (float64, error)
	GetAssetPerformance(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*AssetPerformance, error)
}
