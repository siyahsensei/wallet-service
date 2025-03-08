package transaction

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Create adds a new transaction to the database
	Create(ctx context.Context, transaction *Transaction) error

	// GetByID retrieves a transaction by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)

	// GetByUserID retrieves all transactions for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Transaction, error)

	// GetByAccountID retrieves all transactions for a specific account
	GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*Transaction, error)

	// GetByAssetID retrieves all transactions for a specific asset
	GetByAssetID(ctx context.Context, assetID uuid.UUID, limit, offset int) ([]*Transaction, error)

	// GetByDateRange retrieves all transactions within a date range for a user
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*Transaction, error)

	// GetByType retrieves all transactions of a specific type for a user
	GetByType(ctx context.Context, userID uuid.UUID, transactionType TransactionType, limit, offset int) ([]*Transaction, error)

	// GetByCategory retrieves all transactions of a specific category for a user
	GetByCategory(ctx context.Context, userID uuid.UUID, category string, limit, offset int) ([]*Transaction, error)

	// Update updates a transaction's information
	Update(ctx context.Context, transaction *Transaction) error

	// Delete removes a transaction from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// GetTotalsByCategory gets the total amount of transactions grouped by category within a date range
	GetTotalsByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]float64, error)

	// GetTotalsByType gets the total amount of transactions grouped by type within a date range
	GetTotalsByType(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[TransactionType]float64, error)

	// GetMonthlyTotals gets the total amount of credits and debits for each month in a date range
	GetMonthlyTotals(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*MonthlyTotal, error)
}
