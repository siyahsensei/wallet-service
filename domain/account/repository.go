package account

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create adds a new account to the database
	Create(ctx context.Context, account *Account) error

	// GetByID retrieves an account by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)

	// GetByUserID retrieves all accounts for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error)

	// GetByType retrieves all accounts of a specific type for a user
	GetByType(ctx context.Context, userID uuid.UUID, accountType AccountType) ([]*Account, error)

	// Update updates an account's information
	Update(ctx context.Context, account *Account) error

	// Delete removes an account from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// GetTotalBalance calculates the total balance for a user across all accounts
	// or for specific account types if provided
	GetTotalBalance(ctx context.Context, userID uuid.UUID, accountTypes []AccountType) (float64, error)
}
