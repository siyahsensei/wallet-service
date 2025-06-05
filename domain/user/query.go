package user

type GetUserByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

type GetUserByEmailQuery struct {
	Email string `json:"email" validate:"required,email"`
}

type ListUsersQuery struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}
