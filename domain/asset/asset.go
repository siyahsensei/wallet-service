package asset

import (
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	// Cash assets
	Cash        AssetType = "CASH"
	TermDeposit AssetType = "TERM_DEPOSIT"

	// Investment assets
	Stock  AssetType = "STOCK"
	ETF    AssetType = "ETF"
	Fund   AssetType = "FUND"
	Bond   AssetType = "BOND"
	Option AssetType = "OPTION"
	Future AssetType = "FUTURE"

	// Crypto assets
	Cryptocurrency AssetType = "CRYPTOCURRENCY"
	NFT            AssetType = "NFT"
	DeFiToken      AssetType = "DEFI_TOKEN"

	// Other assets
	PreciousMetal AssetType = "PRECIOUS_METAL"
	RealEstate    AssetType = "REAL_ESTATE"
	Debt          AssetType = "DEBT"
	Receivable    AssetType = "RECEIVABLE"
	Other         AssetType = "OTHER"
)

type Asset struct {
	ID            uuid.UUID `json:"id" db:"id"`
	AccountID     uuid.UUID `json:"account_id" db:"account_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Name          string    `json:"name" db:"name"`
	Type          AssetType `json:"type" db:"type"`
	Symbol        string    `json:"symbol" db:"symbol"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	PurchasePrice float64   `json:"purchase_price" db:"purchase_price"`
	CurrentPrice  float64   `json:"current_price" db:"current_price"`
	Currency      string    `json:"currency" db:"currency"`
	Notes         string    `json:"notes" db:"notes"`
	PurchaseDate  time.Time `json:"purchase_date" db:"purchase_date"`
	LastUpdated   time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func NewAsset(
	userID uuid.UUID,
	accountID uuid.UUID,
	name string,
	assetType AssetType,
	symbol string,
	quantity float64,
	purchasePrice float64,
	currentPrice float64,
	currency string,
	notes string,
	purchaseDate time.Time,
) *Asset {
	now := time.Now()
	return &Asset{
		ID:            uuid.New(),
		UserID:        userID,
		AccountID:     accountID,
		Name:          name,
		Type:          assetType,
		Symbol:        symbol,
		Quantity:      quantity,
		PurchasePrice: purchasePrice,
		CurrentPrice:  currentPrice,
		Currency:      currency,
		Notes:         notes,
		PurchaseDate:  purchaseDate,
		LastUpdated:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (a *Asset) CurrentValue() float64 {
	return a.Quantity * a.CurrentPrice
}

func (a *Asset) PurchaseValue() float64 {
	return a.Quantity * a.PurchasePrice
}

func (a *Asset) ProfitLoss() float64 {
	return a.CurrentValue() - a.PurchaseValue()
}

func (a *Asset) ProfitLossPercentage() float64 {
	if a.PurchaseValue() == 0 {
		return 0
	}
	return (a.ProfitLoss() / a.PurchaseValue()) * 100
}

func (a *Asset) UpdatePrice(price float64) {
	a.CurrentPrice = price
	a.LastUpdated = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *Asset) UpdateQuantity(quantity float64, price float64) {
	if quantity > a.Quantity && price > 0 {
		totalOldValue := a.Quantity * a.PurchasePrice
		additionalValue := (quantity - a.Quantity) * price
		a.PurchasePrice = (totalOldValue + additionalValue) / quantity
	}
	a.Quantity = quantity
	a.UpdatedAt = time.Now()
}
