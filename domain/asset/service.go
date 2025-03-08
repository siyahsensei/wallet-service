package asset

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AssetPerformance struct {
	AssetID        uuid.UUID `json:"asset_id"`
	Name           string    `json:"name"`
	Symbol         string    `json:"symbol"`
	Type           AssetType `json:"type"`
	InitialValue   float64   `json:"initial_value"`
	CurrentValue   float64   `json:"current_value"`
	ProfitLoss     float64   `json:"profit_loss"`
	ProfitLossPerc float64   `json:"profit_loss_percentage"`
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

func (s *Service) CreateAsset(
	ctx context.Context,
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
) (*Asset, error) {
	if name == "" {
		return nil, errors.New("asset name cannot be empty")
	}
	if !isValidAssetType(assetType) {
		return nil, errors.New("invalid asset type")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	asset := NewAsset(
		userID,
		accountID,
		name,
		assetType,
		symbol,
		quantity,
		purchasePrice,
		currentPrice,
		currency,
		notes,
		purchaseDate,
	)
	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
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
	if asset.Name == "" {
		return errors.New("asset name cannot be empty")
	}
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

func (s *Service) UpdateAssetPrice(ctx context.Context, id uuid.UUID, price float64) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	asset, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("asset not found")
	}
	asset.UpdatePrice(price)
	return s.repo.Update(ctx, asset)
}

func (s *Service) UpdateAssetQuantity(ctx context.Context, id uuid.UUID, quantity float64, price float64) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	asset, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("asset not found")
	}
	asset.UpdateQuantity(quantity, price)
	return s.repo.Update(ctx, asset)
}

func (s *Service) GetTotalAssetValue(ctx context.Context, userID uuid.UUID, assetTypes []AssetType) (float64, error) {
	return s.repo.GetTotalValue(ctx, userID, assetTypes)
}

func (s *Service) GetAssetPerformance(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*AssetPerformance, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	return s.repo.GetAssetPerformance(ctx, userID, startDate, endDate)
}

func (s *Service) GetAllAssetTypes() []AssetType {
	return []AssetType{
		Cash,
		TermDeposit,
		Stock,
		ETF,
		Fund,
		Bond,
		Option,
		Future,
		Cryptocurrency,
		NFT,
		DeFiToken,
		PreciousMetal,
		RealEstate,
		Debt,
		Receivable,
		Other,
	}
}

func isValidAssetType(t AssetType) bool {
	validTypes := []AssetType{
		Cash,
		TermDeposit,
		Stock,
		ETF,
		Fund,
		Bond,
		Option,
		Future,
		Cryptocurrency,
		NFT,
		DeFiToken,
		PreciousMetal,
		RealEstate,
		Debt,
		Receivable,
		Other,
	}
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}
