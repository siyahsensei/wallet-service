package account

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	BankAccount       AccountType = "BANK_ACCOUNT"
	SavingsAccount    AccountType = "SAVINGS_ACCOUNT"
	CheckingAccount   AccountType = "CHECKING_ACCOUNT"
	CreditCard        AccountType = "CREDIT_CARD"
	InvestmentAccount AccountType = "INVESTMENT_ACCOUNT"
	CryptoWallet      AccountType = "CRYPTO_WALLET"
	CryptoExchange    AccountType = "CRYPTO_EXCHANGE"
	Broker            AccountType = "BROKER"
	Pension           AccountType = "PENSION"
	Insurance         AccountType = "INSURANCE"
	Home              AccountType = "HOME"
	Safe              AccountType = "SAFE"
	Other             AccountType = "OTHER"
)

type Account struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	UserID       uuid.UUID   `json:"userId" db:"user_id"`
	Name         string      `json:"name" db:"name"`
	AccountType  AccountType `json:"accountType" db:"account_type"`
	Balance      float64     `json:"balance" db:"balance"`
	CurrencyCode string      `json:"currencyCode" db:"currency_code"`
	CreatedAt    time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time   `json:"updatedAt" db:"updated_at"`
}

func NewAccount(command CreateAccountCommand) *Account {
	now := time.Now()
	return &Account{
		ID:           uuid.New(),
		UserID:       uuid.MustParse(command.UserID),
		Name:         command.Name,
		AccountType:  command.AccountType,
		Balance:      command.Balance,
		CurrencyCode: command.CurrencyCode,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
