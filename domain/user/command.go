package user

type RegisterUserCommand struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type LoginUserCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserCommand struct {
	ID        string `json:"id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type ChangePasswordCommand struct {
	UserID      string `json:"userId" validate:"required"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type DeleteUserCommand struct {
	UserID   string `json:"userId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ValidateUserPasswordCommand struct {
	UserID   string `json:"userId" validate:"required"`
	Password string `json:"password" validate:"required"`
}
