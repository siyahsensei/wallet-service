package presentation

import (
	"siyahsensei/wallet-service/domain/account"
	"time"
)

type CreateAccountRequest struct {
	Name        string              `json:"name" validate:"required"`
	AccountType account.AccountType `json:"accountType" validate:"required"`
}

type UpdateAccountRequest struct {
	Name        string              `json:"name" validate:"required"`
	AccountType account.AccountType `json:"accountType" validate:"required"`
}

type AccountResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	AccountType string    `json:"accountType"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AccountWithAssetsResponse struct {
	ID            string              `json:"id"`
	UserID        string              `json:"userId"`
	Name          string              `json:"name"`
	AccountType   string              `json:"accountType"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	Assets        []AssetInfoResponse `json:"assets"`
	TotalBalances map[string]float64  `json:"totalBalances"`
	AssetCounts   map[string]int      `json:"assetCounts"`
	LastUpdated   *time.Time          `json:"lastUpdated"`
}

type AssetInfoResponse struct {
	ID           string    `json:"id"`
	DefinitionID string    `json:"definitionId"`
	Type         string    `json:"type"`
	Quantity     float64   `json:"quantity"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	Currency     string    `json:"currency"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AccountsListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int               `json:"total"`
}

type AccountsWithAssetsListResponse struct {
	Accounts []AccountWithAssetsResponse `json:"accounts"`
	Total    int                         `json:"total"`
}

type AccountSummaryResponse struct {
	TotalAccounts int                         `json:"totalAccounts"`
	ByType        map[account.AccountType]int `json:"byType"`
	ByCurrency    map[string]float64          `json:"byCurrency"`
}
