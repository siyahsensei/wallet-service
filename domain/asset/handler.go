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

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) HandleCreateAssetCommand(
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

func (s *Service) HandleUpdateAssetCommand(ctx context.Context, command UpdateAssetCommand) (*Asset, error) {
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

func (s *Service) HandleDeleteAssetCommand(ctx context.Context, command DeleteAssetCommand) error {
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

func (s *Service) HandleGetAssetByIDQuery(ctx context.Context, query GetAssetByIDQuery) (*Asset, error) {
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

func (s *Service) HandleGetUserAssetsQuery(ctx context.Context, query GetUserAssetsQuery) ([]*Asset, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) HandleGetAccountAssetsQuery(ctx context.Context, query GetAccountAssetsQuery) ([]*Asset, error) {
	accountID, err := uuid.Parse(query.AccountID)
	if err != nil {
		return nil, errors.New("invalid account ID")
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	assets, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var userAssets []*Asset
	for _, asset := range assets {
		if asset.UserID == userID {
			userAssets = append(userAssets, asset)
		}
	}

	return userAssets, nil
}

func (s *Service) HandleGetAssetsByTypeQuery(ctx context.Context, query GetAssetsByTypeQuery) ([]*Asset, error) {
	if !isValidAssetType(query.AssetType) {
		return nil, errors.New("invalid asset type")
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByType(ctx, userID, query.AssetType)
}

func (s *Service) HandleFilterAssetsQuery(ctx context.Context, query FilterAssetsQuery) ([]*Asset, error) {
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

func (s *Service) HandleGetAssetPerformanceQuery(ctx context.Context, query GetAssetPerformanceQuery) ([]*AssetPerformance, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetAssetPerformance(ctx, userID, query.StartDate, query.EndDate)
}

func (s *Service) HandleGetTotalValueQuery(ctx context.Context, query GetTotalValueQuery) (float64, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}

	return s.repo.GetTotalValue(ctx, userID, query.AssetTypes)
}

func (s *Service) matchesFilter(asset *Asset, query FilterAssetsQuery) bool {
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

func (s *Service) CreateAsset(
	ctx context.Context, command CreateAssetCommand) (*Asset, error) {
	return s.HandleCreateAssetCommand(ctx, command)
}

func (s *Service) GetAssetByID(ctx context.Context, id uuid.UUID) (*Asset, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetUserAssets(ctx context.Context, userID uuid.UUID) ([]*Asset, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) GetAccountAssets(ctx context.Context, accountID uuid.UUID) ([]*Asset, error) {
	return s.repo.GetByAccountID(ctx, accountID)
}

func (s *Service) GetAssetsByType(ctx context.Context, userID uuid.UUID, assetType AssetType) ([]*Asset, error) {
	if !isValidAssetType(assetType) {
		return nil, errors.New("invalid asset type")
	}

	return s.repo.GetByType(ctx, userID, assetType)
}

func (s *Service) UpdateAsset(ctx context.Context, asset *Asset) error {
	if !isValidAssetType(asset.Type) {
		return errors.New("invalid asset type")
	}
	if asset.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	_, err := s.repo.GetByID(ctx, asset.ID)
	if err != nil {
		return errors.New("asset not found")
	}
	return s.repo.Update(ctx, asset)
}

func (s *Service) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("asset not found")
	}
	return s.repo.Delete(ctx, id)
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
