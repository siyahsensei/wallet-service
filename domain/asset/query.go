package asset

import (
	"time"
)

type GetAssetByIDQuery struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

type GetUserAssetsQuery struct {
	UserID string `json:"userId" validate:"required"`
}

type GetAccountAssetsQuery struct {
	AccountID string `json:"accountId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
}

type GetAssetsByTypeQuery struct {
	UserID    string    `json:"userId" validate:"required"`
	AssetType AssetType `json:"assetType" validate:"required"`
}

type FilterAssetsQuery struct {
	UserID      string     `json:"userId" validate:"required"`
	AccountID   *string    `json:"accountId,omitempty"`
	AssetType   *AssetType `json:"assetType,omitempty"`
	MinQuantity *float64   `json:"minQuantity,omitempty"`
	MaxQuantity *float64   `json:"maxQuantity,omitempty"`
	CreatedFrom *time.Time `json:"createdFrom,omitempty"`
	CreatedTo   *time.Time `json:"createdTo,omitempty"`
	Limit       int        `json:"limit,omitempty"`
	Offset      int        `json:"offset,omitempty"`
}

type GetAssetPerformanceQuery struct {
	UserID    string    `json:"userId" validate:"required"`
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required"`
}

type GetTotalValueQuery struct {
	UserID     string      `json:"userId" validate:"required"`
	AssetTypes []AssetType `json:"assetTypes,omitempty"`
}
