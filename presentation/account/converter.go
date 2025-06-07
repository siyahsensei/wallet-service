package presentation

import "siyahsensei/wallet-service/domain/account"

func ToAccountResponse(a *account.Account) AccountResponse {
	return AccountResponse{
		ID:          a.ID.String(),
		UserID:      a.UserID.String(),
		Name:        a.Name,
		AccountType: string(a.AccountType),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func ToAccountWithAssetsResponse(a *account.AccountWithAssets) AccountWithAssetsResponse {
	var assets []AssetInfoResponse
	for _, asset := range a.Assets {
		assets = append(assets, AssetInfoResponse{
			ID:           asset.ID.String(),
			DefinitionID: asset.DefinitionID.String(),
			Type:         asset.Type,
			Quantity:     asset.Quantity,
			Symbol:       asset.Symbol,
			Name:         asset.Name,
			Currency:     asset.Currency,
			UpdatedAt:    asset.UpdatedAt,
		})
	}

	return AccountWithAssetsResponse{
		ID:          a.Account.ID.String(),
		UserID:      a.Account.UserID.String(),
		Name:        a.Account.Name,
		AccountType: string(a.Account.AccountType),
		CreatedAt:   a.Account.CreatedAt,
		UpdatedAt:   a.Account.UpdatedAt,
		Assets:      assets,
		AssetCounts: a.AssetCounts,
		LastUpdated: a.LastUpdated,
	}
}

func ToAccountSummaryResponse(s *account.AccountSummary) AccountSummaryResponse {
	return AccountSummaryResponse{
		TotalAccounts: s.TotalAccounts,
		ByType:        s.ByType,
		ByCurrency:    s.ByCurrency,
	}
}
