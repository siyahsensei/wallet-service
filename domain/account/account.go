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
	ID          uuid.UUID   `json:"id" db:"id"`
	UserID      uuid.UUID   `json:"userId" db:"user_id"`
	Name        string      `json:"name" db:"name"`
	AccountType AccountType `json:"accountType" db:"account_type"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at"`
}

type AccountWithAssets struct {
	Account
	Assets      []AssetInfo    `json:"assets"`
	AssetCounts map[string]int `json:"assetCounts"`
	LastUpdated *time.Time     `json:"lastUpdated"`
}

type AssetInfo struct {
	ID           uuid.UUID `json:"id"`
	DefinitionID uuid.UUID `json:"definitionId"`
	Type         string    `json:"type"`
	Quantity     float64   `json:"quantity"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	Currency     string    `json:"currency"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewAccount(command CreateAccountCommand) *Account {
	now := time.Now()
	return &Account{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(command.UserID),
		Name:        command.Name,
		AccountType: command.AccountType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
