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

type FilterAccountsQuery struct {
	UserID      string       `json:"userId" validate:"required"`
	AccountType *AccountType `json:"accountType,omitempty"`
	Limit       int          `json:"limit,omitempty"`
	Offset      int          `json:"offset,omitempty"`
}

type GetAccountSummaryQuery struct {
	UserID string `json:"userId" validate:"required"`
}
