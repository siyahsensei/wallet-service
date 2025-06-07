package presentation

import "siyahsensei/wallet-service/domain/asset"

func ToAssetResponse(a *asset.Asset) AssetResponse {
	return AssetResponse{
		ID:           a.ID.String(),
		UserID:       a.UserID.String(),
		AccountID:    a.AccountID.String(),
		DefinitionID: a.DefinitionID.String(),
		Type:         string(a.Type),
		Quantity:     a.Quantity,
		Notes:        a.Notes,
		PurchaseDate: a.PurchaseDate,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func ToAssetPerformanceResponse(ap *asset.AssetPerformance) AssetPerformanceResponse {
	return AssetPerformanceResponse{
		AssetID:        ap.AssetID.String(),
		Name:           ap.Name,
		Symbol:         ap.Symbol,
		Type:           string(ap.Type),
		InitialValue:   ap.InitialValue,
		CurrentValue:   ap.CurrentValue,
		ProfitLoss:     ap.ProfitLoss,
		ProfitLossPerc: ap.ProfitLossPerc,
		Currency:       ap.Currency,
	}
}