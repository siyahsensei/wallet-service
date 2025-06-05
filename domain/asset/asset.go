package asset

import (
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	Cash        AssetType = "CASH"
	TermDeposit AssetType = "TERM_DEPOSIT"

	Stock AssetType = "STOCK"
	ETF   AssetType = "ETF"
	Fund  AssetType = "FUND"
	Bond  AssetType = "BOND"

	Cryptocurrency AssetType = "CRYPTOCURRENCY"
	NFT            AssetType = "NFT"
	DeFiToken      AssetType = "DEFI_TOKEN"

	PreciousMetal AssetType = "PRECIOUS_METAL"
	RealEstate    AssetType = "REAL_ESTATE"
	Debt          AssetType = "DEBT"
	Receivable    AssetType = "RECEIVABLE"
	Salary        AssetType = "SALARY"
	Other         AssetType = "OTHER"
)

type Asset struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"userId" db:"user_id"`
	AccountID    uuid.UUID `json:"accountId" db:"account_id"`
	DefinitionID uuid.UUID `json:"definitionId" db:"definition_id"`
	Type         AssetType `json:"type" db:"type"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	Notes        string    `json:"notes" db:"notes"`
	PurchaseDate time.Time `json:"purchaseDate" db:"purchase_date"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

func NewAsset(command CreateAssetCommand) *Asset {
	now := time.Now()
	purchaseDate := time.Unix(command.PurchaseDate, 0)
	return &Asset{
		ID:           uuid.New(),
		UserID:       uuid.MustParse(command.UserID),
		AccountID:    uuid.MustParse(command.AccountID),
		DefinitionID: uuid.MustParse(command.DefinitionID),
		Type:         command.Type,
		Quantity:     command.Quantity,
		Notes:        command.Notes,
		PurchaseDate: purchaseDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
