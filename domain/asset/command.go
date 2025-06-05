package asset

type CreateAssetCommand struct {
	UserID       string    `json:"userId" validate:"required"`
	AccountID    string    `json:"accountId" validate:"required"`
	DefinitionID string    `json:"definitionId" validate:"required"`
	Type         AssetType `json:"type" validate:"required"`
	Quantity     float64   `json:"quantity" validate:"required"`
	Notes        string    `json:"notes"`
	PurchaseDate int64     `json:"purchaseDate" validate:"required"`
}

type UpdateAssetCommand struct {
	ID           string    `json:"id" validate:"required"`
	UserID       string    `json:"userId" validate:"required"`
	AccountID    string    `json:"accountId" validate:"required"`
	DefinitionID string    `json:"definitionId" validate:"required"`
	Type         AssetType `json:"type" validate:"required"`
	Quantity     float64   `json:"quantity" validate:"required"`
	Notes        string    `json:"notes"`
	PurchaseDate int64     `json:"purchaseDate" validate:"required"`
}

type DeleteAssetCommand struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}
