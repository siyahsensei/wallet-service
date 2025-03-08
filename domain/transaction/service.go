package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type MonthlyTotal struct {
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	TotalIn   float64 `json:"totalIn"`
	TotalOut  float64 `json:"totalOut"`
	NetAmount float64 `json:"netAmount"`
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTransaction(
	ctx context.Context,
	userID uuid.UUID,
	accountID uuid.UUID,
	assetID *uuid.UUID,
	transactionType TransactionType,
	amount float64,
	quantity float64,
	price float64,
	fee float64,
	currency string,
	description string,
	category string,
	date time.Time,
	toAccountID *uuid.UUID,
	transactionHash string,
) (*Transaction, error) {
	if !isValidTransactionType(transactionType) {
		return nil, errors.New("invalid transaction type")
	}
	if amount <= 0 && transactionType != Withdrawal && transactionType != Expense &&
		transactionType != Fee && transactionType != Tax && transactionType != Repayment {
		return nil, errors.New("amount must be greater than zero for non-expense transactions")
	}
	if transactionType == Transfer && toAccountID == nil {
		return nil, errors.New("destination account must be specified for transfers")
	}
	if (transactionType == Buy || transactionType == Sell) && assetID == nil {
		return nil, errors.New("asset must be specified for buy/sell transactions")
	}
	transaction := NewTransaction(
		userID,
		accountID,
		assetID,
		transactionType,
		amount,
		quantity,
		price,
		fee,
		currency,
		description,
		category,
		date,
		toAccountID,
		transactionHash,
	)
	if err := s.repo.Create(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

func (s *Service) GetAccountTransactions(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	return s.repo.GetByAccountID(ctx, accountID, limit, offset)
}

func (s *Service) GetAssetTransactions(ctx context.Context, assetID uuid.UUID, limit, offset int) ([]*Transaction, error) {
	return s.repo.GetByAssetID(ctx, assetID, limit, offset)
}

func (s *Service) GetTransactionsByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*Transaction, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	return s.repo.GetByDateRange(ctx, userID, startDate, endDate, limit, offset)
}

func (s *Service) GetTransactionsByType(ctx context.Context, userID uuid.UUID, transactionType TransactionType, limit, offset int) ([]*Transaction, error) {
	if !isValidTransactionType(transactionType) {
		return nil, errors.New("invalid transaction type")
	}
	return s.repo.GetByType(ctx, userID, transactionType, limit, offset)
}

func (s *Service) GetTransactionsByCategory(ctx context.Context, userID uuid.UUID, category string, limit, offset int) ([]*Transaction, error) {
	return s.repo.GetByCategory(ctx, userID, category, limit, offset)
}

func (s *Service) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	if !isValidTransactionType(transaction.Type) {
		return errors.New("invalid transaction type")
	}
	if transaction.Amount <= 0 && transaction.Type != Withdrawal && transaction.Type != Expense &&
		transaction.Type != Fee && transaction.Type != Tax && transaction.Type != Repayment {
		return errors.New("amount must be greater than zero for non-expense transactions")
	}
	if transaction.Type == Transfer && transaction.ToAccountID == nil {
		return errors.New("destination account must be specified for transfers")
	}
	if (transaction.Type == Buy || transaction.Type == Sell) && transaction.AssetID == nil {
		return errors.New("asset must be specified for buy/sell transactions")
	}
	_, err := s.repo.GetByID(ctx, transaction.ID)
	if err != nil {
		return errors.New("transaction not found")
	}
	return s.repo.Update(ctx, transaction)
}

func (s *Service) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("transaction not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetTotalsByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]float64, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	return s.repo.GetTotalsByCategory(ctx, userID, startDate, endDate)
}

func (s *Service) GetTotalsByType(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[TransactionType]float64, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	return s.repo.GetTotalsByType(ctx, userID, startDate, endDate)
}

func (s *Service) GetMonthlyTotals(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*MonthlyTotal, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	return s.repo.GetMonthlyTotals(ctx, userID, startDate, endDate)
}

func (s *Service) GetAllTransactionTypes() []TransactionType {
	return []TransactionType{
		Deposit,
		Withdrawal,
		Transfer,
		Buy,
		Sell,
		Dividend,
		Interest,
		Fee,
		Income,
		Expense,
		Tax,
		Rebalance,
		Split,
		Merger,
		Staking,
		Mining,
		Airdrop,
		Lending,
		Borrowing,
		Repayment,
	}
}

func isValidTransactionType(t TransactionType) bool {
	validTypes := []TransactionType{
		Deposit,
		Withdrawal,
		Transfer,
		Buy,
		Sell,
		Dividend,
		Interest,
		Fee,
		Income,
		Expense,
		Tax,
		Rebalance,
		Split,
		Merger,
		Staking,
		Mining,
		Airdrop,
		Lending,
		Borrowing,
		Repayment,
	}
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}
