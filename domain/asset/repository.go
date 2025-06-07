package asset

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, asset *Asset) error
	Update(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Asset, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Asset, error)
	GetByType(ctx context.Context, userID uuid.UUID, assetType AssetType) ([]*Asset, error)
}
