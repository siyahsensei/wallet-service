package account

type GetAccountByIDQuery struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

type GetUserAccountsQuery struct {
	UserID string `json:"userId" validate:"required"`
}

type GetAccountsByTypeQuery struct {
	UserID      string      `json:"userId" validate:"required"`
	AccountType AccountType `json:"accountType" validate:"required"`
}

type GetAccountsByCurrencyQuery struct {
	UserID       string `json:"userId" validate:"required"`
	CurrencyCode string `json:"currencyCode" validate:"required"`
}

type FilterAccountsQuery struct {
	UserID       string       `json:"userId" validate:"required"`
	AccountType  *AccountType `json:"accountType,omitempty"`
	CurrencyCode *string      `json:"currencyCode,omitempty"`
	MinBalance   *float64     `json:"minBalance,omitempty"`
	MaxBalance   *float64     `json:"maxBalance,omitempty"`
	Limit        int          `json:"limit,omitempty"`
	Offset       int          `json:"offset,omitempty"`
}

type GetAccountSummaryQuery struct {
	UserID string `json:"userId" validate:"required"`
}
