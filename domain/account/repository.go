package account

import (
	"context"

	"github.com/google/uuid"
)

type AccountSummary struct {
	TotalAccounts int                 `json:"totalAccounts"`
	TotalBalance  float64             `json:"totalBalance"`
	ByType        map[AccountType]int `json:"byType"`
	ByCurrency    map[string]float64  `json:"byCurrency"`
}

type Repository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error)
	GetByType(ctx context.Context, userID uuid.UUID, accountType AccountType) ([]*Account, error)
	GetByCurrency(ctx context.Context, userID uuid.UUID, currencyCode string) ([]*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error
	GetAccountSummary(ctx context.Context, userID uuid.UUID) (*AccountSummary, error)
	Filter(ctx context.Context, query FilterAccountsQuery) ([]*Account, error)
}
