package presentation

import "siyahsensei/wallet-service/domain/user"

func ToPublicUser(u *user.User) *UserPublic {
	return &UserPublic{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}