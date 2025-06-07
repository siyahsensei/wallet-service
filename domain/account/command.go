package account

type CreateAccountCommand struct {
	UserID      string      `json:"userId" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	AccountType AccountType `json:"accountType" validate:"required"`
}

type UpdateAccountCommand struct {
	ID          string      `json:"id" validate:"required"`
	UserID      string      `json:"userId" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	AccountType AccountType `json:"accountType" validate:"required"`
}

type DeleteAccountCommand struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}
