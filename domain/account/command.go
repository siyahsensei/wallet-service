package account

type CreateAccountCommand struct {
	UserID       string      `json:"userId" validate:"required"`
	Name         string      `json:"name" validate:"required"`
	AccountType  AccountType `json:"accountType" validate:"required"`
	Balance      float64     `json:"balance"`
	CurrencyCode string      `json:"currencyCode" validate:"required"`
}

type UpdateAccountCommand struct {
	ID           string      `json:"id" validate:"required"`
	UserID       string      `json:"userId" validate:"required"`
	Name         string      `json:"name" validate:"required"`
	AccountType  AccountType `json:"accountType" validate:"required"`
	Balance      float64     `json:"balance"`
	CurrencyCode string      `json:"currencyCode" validate:"required"`
}

type DeleteAccountCommand struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

type UpdateAccountBalanceCommand struct {
	ID     string  `json:"id" validate:"required"`
	UserID string  `json:"userId" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}
