package account

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	// Cash account types
	BankAccount AccountType = "BANK_ACCOUNT"
	CashWallet  AccountType = "CASH_WALLET"

	// Investment account types
	StockAccount      AccountType = "STOCK_ACCOUNT"
	InvestmentFund    AccountType = "INVESTMENT_FUND"
	BondAccount       AccountType = "BOND_ACCOUNT"
	DerivativeAccount AccountType = "DERIVATIVE_ACCOUNT"

	// Crypto account types
	CryptoExchange AccountType = "CRYPTO_EXCHANGE"
	CryptoWallet   AccountType = "CRYPTO_WALLET"
	DeFiProtocol   AccountType = "DEFI_PROTOCOL"
	NFTCollection  AccountType = "NFT_COLLECTION"

	// Other account types
	PreciousMetals AccountType = "PRECIOUS_METALS"
	RealEstate     AccountType = "REAL_ESTATE"
	DebtAccount    AccountType = "DEBT_ACCOUNT"
	Receivable     AccountType = "RECEIVABLE"
)

type Account struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	UserID      uuid.UUID   `json:"userId" db:"user_id"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Type        AccountType `json:"type" db:"type"`
	Institution string      `json:"institution" db:"institution"`
	Currency    string      `json:"currency" db:"currency"`
	Balance     float64     `json:"balance" db:"balance"`
	IsManual    bool        `json:"isManual" db:"is_manual"`
	Icon        string      `json:"icon" db:"icon"`
	Color       string      `json:"color" db:"color"`
	// API connection details (:müptezel_smile:)
	APIKey      string     `json:"-" db:"api_key"`
	APISecret   string     `json:"-" db:"api_secret"`
	IsConnected bool       `json:"isConnected" db:"is_connected"`
	LastSync    *time.Time `json:"lastSync" db:"last_sync"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

func NewAccount(
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
) *Account {
	return &Account{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: description,
		Type:        accountType,
		Institution: institution,
		Currency:    currency,
		Balance:     balance,
		IsManual:    isManual,
		Icon:        icon,
		Color:       color,
		IsConnected: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (a *Account) SetAPICredentials(apiKey, apiSecret string) {
	a.APIKey = apiKey
	a.APISecret = apiSecret
	a.IsConnected = true
	a.UpdatedAt = time.Now()
}

func (a *Account) UpdateBalance(balance float64) {
	a.Balance = balance
	a.UpdatedAt = time.Now()
	if !a.IsManual {
		now := time.Now()
		a.LastSync = &now
	}
}

func (a *Account) IsCashAccount() bool {
	return a.Type == BankAccount || a.Type == CashWallet
}

func (a *Account) IsInvestmentAccount() bool {
	return a.Type == StockAccount || a.Type == InvestmentFund ||
		a.Type == BondAccount || a.Type == DerivativeAccount
}

func (a *Account) IsCryptoAccount() bool {
	return a.Type == CryptoExchange || a.Type == CryptoWallet ||
		a.Type == DeFiProtocol || a.Type == NFTCollection
}

func (a *Account) IsAssetAccount() bool {
	return a.Type == PreciousMetals || a.Type == RealEstate || a.Type == Receivable
}

func (a *Account) IsLiabilityAccount() bool {
	return a.Type == DebtAccount
}
