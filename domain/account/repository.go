package account

import (
	"context"

	"github.com/google/uuid"
)

type AccountSummary struct {
	TotalAccounts int                 `json:"totalAccounts"`
	ByType        map[AccountType]int `json:"byType"`
	ByCurrency    map[string]float64  `json:"byCurrency"` // calculated from assets
}

type Repository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByIDWithAssets(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*AccountWithAssets, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error)
	GetByUserIDWithAssets(ctx context.Context, userID uuid.UUID) ([]*AccountWithAssets, error)
	GetByType(ctx context.Context, userID uuid.UUID, accountType AccountType) ([]*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAccountSummary(ctx context.Context, userID uuid.UUID) (*AccountSummary, error)
	Filter(ctx context.Context, query FilterAccountsQuery) ([]*Account, error)
}
