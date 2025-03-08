package account

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateAccount(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	description string,
	accountType AccountType,
	institution string,
	currency string,
	balance float64,
	isManual bool,
	icon string,
	color string,
) (*Account, error) {
	if name == "" {
		return nil, errors.New("account name cannot be empty")
	}
	if !isValidAccountType(accountType) {
		return nil, errors.New("invalid account type")
	}
	account := NewAccount(
		userID,
		name,
		description,
		accountType,
		institution,
		currency,
		balance,
		isManual,
		icon,
		color,
	)
	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) GetAccountsByType(ctx context.Context, userID uuid.UUID, accountType AccountType) ([]*Account, error) {
	if !isValidAccountType(accountType) {
		return nil, errors.New("invalid account type")
	}

	return s.repo.GetByType(ctx, userID, accountType)
}

func (s *Service) UpdateAccount(ctx context.Context, account *Account) error {
	if account.Name == "" {
		return errors.New("account name cannot be empty")
	}
	if !isValidAccountType(account.Type) {
		return errors.New("invalid account type")
	}
	_, err := s.repo.GetByID(ctx, account.ID)
	if err != nil {
		return errors.New("account not found")
	}
	return s.repo.Update(ctx, account)
}

func (s *Service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("account not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) SetAPICredentials(ctx context.Context, id uuid.UUID, apiKey, apiSecret string) error {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("account not found")
	}
	account.SetAPICredentials(apiKey, apiSecret)
	return s.repo.Update(ctx, account)
}

func (s *Service) UpdateAccountBalance(ctx context.Context, id uuid.UUID, balance float64) error {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("account not found")
	}
	account.UpdateBalance(balance)
	return s.repo.Update(ctx, account)
}

func (s *Service) GetTotalBalance(ctx context.Context, userID uuid.UUID, accountTypes []AccountType) (float64, error) {
	return s.repo.GetTotalBalance(ctx, userID, accountTypes)
}

func (s *Service) GetAllAccountTypes() []AccountType {
	return []AccountType{
		BankAccount,
		CashWallet,
		StockAccount,
		InvestmentFund,
		BondAccount,
		DerivativeAccount,
		CryptoExchange,
		CryptoWallet,
		DeFiProtocol,
		NFTCollection,
		PreciousMetals,
		RealEstate,
		DebtAccount,
		Receivable,
	}
}

func isValidAccountType(t AccountType) bool {
	validTypes := []AccountType{
		BankAccount,
		CashWallet,
		StockAccount,
		InvestmentFund,
		BondAccount,
		DerivativeAccount,
		CryptoExchange,
		CryptoWallet,
		DeFiProtocol,
		NFTCollection,
		PreciousMetals,
		RealEstate,
		DebtAccount,
		Receivable,
	}
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}
