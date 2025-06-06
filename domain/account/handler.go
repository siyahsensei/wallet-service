package account

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleCreateAccountCommand(ctx context.Context, command CreateAccountCommand) (*Account, error) {
	if command.Name == "" {
		return nil, errors.New("account name is required")
	}
	if !isValidAccountType(command.AccountType) {
		return nil, errors.New("invalid account type")
	}

	account := NewAccount(command)
	if err := h.repo.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (h *Handler) HandleUpdateAccountCommand(ctx context.Context, command UpdateAccountCommand) (*Account, error) {
	if command.Name == "" {
		return nil, errors.New("account name is required")
	}
	if !isValidAccountType(command.AccountType) {
		return nil, errors.New("invalid account type")
	}

	accountID, err := uuid.Parse(command.ID)
	if err != nil {
		return nil, errors.New("invalid account ID")
	}

	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	existingAccount, err := h.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	if existingAccount.UserID != userID {
		return nil, errors.New("unauthorized: account does not belong to user")
	}

	existingAccount.Name = command.Name
	existingAccount.AccountType = command.AccountType
	existingAccount.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, existingAccount); err != nil {
		return nil, err
	}

	return existingAccount, nil
}

func (h *Handler) HandleDeleteAccountCommand(ctx context.Context, command DeleteAccountCommand) error {
	accountID, err := uuid.Parse(command.ID)
	if err != nil {
		return errors.New("invalid account ID")
	}

	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	existingAccount, err := h.repo.GetByID(ctx, accountID)
	if err != nil {
		return errors.New("account not found")
	}

	if existingAccount.UserID != userID {
		return errors.New("unauthorized: account does not belong to user")
	}

	return h.repo.Delete(ctx, accountID)
}

func (h *Handler) HandleGetAccountByIDQuery(ctx context.Context, query GetAccountByIDQuery) (*Account, error) {
	accountID, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, errors.New("invalid account ID")
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	account, err := h.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, errors.New("unauthorized: account does not belong to user")
	}

	return account, nil
}

func (h *Handler) HandleGetAccountByIDWithAssetsQuery(ctx context.Context, query GetAccountByIDQuery) (*AccountWithAssets, error) {
	accountID, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, errors.New("invalid account ID")
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return h.repo.GetByIDWithAssets(ctx, accountID, userID)
}

func (h *Handler) HandleGetUserAccountsQuery(ctx context.Context, query GetUserAccountsQuery) ([]*Account, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return h.repo.GetByUserID(ctx, userID)
}

func (h *Handler) HandleGetUserAccountsWithAssetsQuery(ctx context.Context, query GetUserAccountsQuery) ([]*AccountWithAssets, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return h.repo.GetByUserIDWithAssets(ctx, userID)
}

func (h *Handler) HandleGetAccountsByTypeQuery(ctx context.Context, query GetAccountsByTypeQuery) ([]*Account, error) {
	if !isValidAccountType(query.AccountType) {
		return nil, errors.New("invalid account type")
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return h.repo.GetByType(ctx, userID, query.AccountType)
}

func (h *Handler) HandleFilterAccountsQuery(ctx context.Context, query FilterAccountsQuery) ([]*Account, error) {
	_, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if query.AccountType != nil && !isValidAccountType(*query.AccountType) {
		return nil, errors.New("invalid account type")
	}

	return h.repo.Filter(ctx, query)
}

func (h *Handler) HandleGetAccountSummaryQuery(ctx context.Context, query GetAccountSummaryQuery) (*AccountSummary, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return h.repo.GetAccountSummary(ctx, userID)
}

func isValidAccountType(t AccountType) bool {
	switch t {
	case BankAccount, SavingsAccount, CheckingAccount, CreditCard, InvestmentAccount,
		CryptoWallet, CryptoExchange, Broker, Pension, Insurance, Home, Safe, Other:
		return true
	default:
		return false
	}
}
