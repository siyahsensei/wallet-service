package presentation

import (
	"siyahsensei/wallet-service/domain/asset"
	"time"
)

type CreateAssetRequest struct {
	AccountID    string          `json:"accountId" validate:"required"`
	DefinitionID string          `json:"definitionId" validate:"required"`
	Type         asset.AssetType `json:"type" validate:"required"`
	Quantity     float64         `json:"quantity" validate:"required"`
	Notes        string          `json:"notes"`
	PurchaseDate int64           `json:"purchaseDate" validate:"required"`
}

type UpdateAssetRequest struct {
	AccountID    string          `json:"accountId" validate:"required"`
	DefinitionID string          `json:"definitionId" validate:"required"`
	Type         asset.AssetType `json:"type" validate:"required"`
	Quantity     float64         `json:"quantity" validate:"required"`
	Notes        string          `json:"notes"`
	PurchaseDate int64           `json:"purchaseDate" validate:"required"`
}

type AssetResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	AccountID    string    `json:"accountId"`
	DefinitionID string    `json:"definitionId"`
	Type         string    `json:"type"`
	Quantity     float64   `json:"quantity"`
	Notes        string    `json:"notes"`
	PurchaseDate time.Time `json:"purchaseDate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AssetsListResponse struct {
	Assets []AssetResponse `json:"assets"`
	Total  int             `json:"total"`
}

type AssetPerformanceResponse struct {
	AssetID        string  `json:"assetId"`
	Name           string  `json:"name"`
	Symbol         string  `json:"symbol"`
	Type           string  `json:"type"`
	InitialValue   float64 `json:"initialValue"`
	CurrentValue   float64 `json:"currentValue"`
	ProfitLoss     float64 `json:"profitLoss"`
	ProfitLossPerc float64 `json:"profitLossPercentage"`
	Currency       string  `json:"currency"`
}

type TotalValueResponse struct {
	TotalValue float64 `json:"totalValue"`
}