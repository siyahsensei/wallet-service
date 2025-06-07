package asset

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AssetPerformance struct {
	AssetID        uuid.UUID `json:"assetId"`
	Name           string    `json:"name"`
	Symbol         string    `json:"symbol"`
	Type           AssetType `json:"type"`
	InitialValue   float64   `json:"initialValue"`
	CurrentValue   float64   `json:"currentValue"`
	ProfitLoss     float64   `json:"profitLoss"`
	ProfitLossPerc float64   `json:"profitLossPercentage"`
	Currency       string    `json:"currency"`
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (s *Handler) HandleCreateAssetCommand(
	ctx context.Context, command CreateAssetCommand) (*Asset, error) {
	if !isValidAssetType(command.Type) {
		return nil, errors.New("invalid asset type")
	}
	if command.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	asset := NewAsset(command)
	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *Handler) HandleUpdateAssetCommand(ctx context.Context, command UpdateAssetCommand) (*Asset, error) {
	if !isValidAssetType(command.Type) {
		return nil, errors.New("invalid asset type")
	}
	if command.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	assetID, err := uuid.Parse(command.ID)
	if err != nil {
		return nil, errors.New("invalid asset ID")
	}

	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	existingAsset, err := s.repo.GetByID(ctx, assetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	if existingAsset.UserID != userID {
		return nil, errors.New("unauthorized: asset does not belong to user")
	}
	existingAsset.AccountID = uuid.MustParse(command.AccountID)
	existingAsset.DefinitionID = uuid.MustParse(command.DefinitionID)
	existingAsset.Type = command.Type
	existingAsset.Quantity = command.Quantity
	existingAsset.Notes = command.Notes
	existingAsset.PurchaseDate = time.Unix(command.PurchaseDate, 0)
	existingAsset.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existingAsset); err != nil {
		return nil, err
	}

	return existingAsset, nil
}

func (s *Handler) HandleDeleteAssetCommand(ctx context.Context, command DeleteAssetCommand) error {
	assetID, err := uuid.Parse(command.ID)
	if err != nil {
		return errors.New("invalid asset ID")
	}

	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	existingAsset, err := s.repo.GetByID(ctx, assetID)
	if err != nil {
		return errors.New("asset not found")
	}

	if existingAsset.UserID != userID {
		return errors.New("unauthorized: asset does not belong to user")
	}

	return s.repo.Delete(ctx, assetID)
}

func (s *Handler) HandleGetAssetByIDQuery(ctx context.Context, query GetAssetByIDQuery) (*Asset, error) {
	assetID, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, errors.New("invalid asset ID")
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	asset, err := s.repo.GetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if asset.UserID != userID {
		return nil, errors.New("unauthorized: asset does not belong to user")
	}

	return asset, nil
}

func (s *Handler) HandleGetUserAssetsQuery(ctx context.Context, query GetUserAssetsQuery) ([]*Asset, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByUserID(ctx, userID)
}

func (s *Handler) HandleFilterAssetsQuery(ctx context.Context, query FilterAssetsQuery) ([]*Asset, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	assets, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*Asset
	for _, asset := range assets {
		if s.matchesFilter(asset, query) {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	if query.Offset > 0 && query.Offset < len(filteredAssets) {
		filteredAssets = filteredAssets[query.Offset:]
	}

	if query.Limit > 0 && query.Limit < len(filteredAssets) {
		filteredAssets = filteredAssets[:query.Limit]
	}

	return filteredAssets, nil
}

func (s *Handler) matchesFilter(asset *Asset, query FilterAssetsQuery) bool {
	if query.AccountID != nil {
		accountID, err := uuid.Parse(*query.AccountID)
		if err != nil || asset.AccountID != accountID {
			return false
		}
	}

	if query.AssetType != nil && asset.Type != *query.AssetType {
		return false
	}

	if query.MinQuantity != nil && asset.Quantity < *query.MinQuantity {
		return false
	}

	if query.MaxQuantity != nil && asset.Quantity > *query.MaxQuantity {
		return false
	}

	if query.CreatedFrom != nil && asset.CreatedAt.Before(*query.CreatedFrom) {
		return false
	}

	if query.CreatedTo != nil && asset.CreatedAt.After(*query.CreatedTo) {
		return false
	}

	return true
}

func isValidAssetType(t AssetType) bool {
	validTypes := []AssetType{
		Cash,
		TermDeposit,
		Stock,
		ETF,
		Fund,
		Bond,
		Cryptocurrency,
		NFT,
		DeFiToken,
		PreciousMetal,
		RealEstate,
		Debt,
		Receivable,
		Salary,
		Other,
	}
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}
