package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	repo        Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func NewHandler(repo Repository, jwtSecret string, tokenExpiry time.Duration) *Handler {
	return &Handler{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

func (s *Handler) HandleRegisterUserCommand(ctx context.Context, command RegisterUserCommand) (*User, error) {
	existingUser, err := s.repo.GetByEmail(ctx, command.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user, err := NewUser(command.Email, command.Password, command.FirstName, command.LastName)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Handler) HandleLoginUserCommand(ctx context.Context, command LoginUserCommand) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, command.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := user.ComparePassword(command.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Handler) HandleUpdateUserCommand(ctx context.Context, command UpdateUserCommand) (*User, error) {
	userID, err := uuid.Parse(command.ID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Email = command.Email
	user.FirstName = command.FirstName
	user.LastName = command.LastName
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Handler) HandleChangePasswordCommand(ctx context.Context, command ChangePasswordCommand) error {
	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := user.ComparePassword(command.OldPassword); err != nil {
		return errors.New("incorrect current password")
	}

	if err := user.UpdatePassword(command.NewPassword); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

func (s *Handler) HandleDeleteUserCommand(ctx context.Context, command DeleteUserCommand) error {
	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := user.ComparePassword(command.Password); err != nil {
		return errors.New("invalid password")
	}

	return s.repo.Delete(ctx, userID)
}

func (s *Handler) HandleValidateUserPasswordCommand(ctx context.Context, command ValidateUserPasswordCommand) error {
	userID, err := uuid.Parse(command.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	return user.ComparePassword(command.Password)
}

func (s *Handler) HandleGetUserByIDQuery(ctx context.Context, query GetUserByIDQuery) (*User, error) {
	userID, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByID(ctx, userID)
}

func (s *Handler) HandleGetUserByEmailQuery(ctx context.Context, query GetUserByEmailQuery) (*User, error) {
	return s.repo.GetByEmail(ctx, query.Email)
}

func (s *Handler) HandleListUsersQuery(ctx context.Context, query ListUsersQuery) ([]*User, error) {
	return s.repo.List(ctx, query.Offset, query.Limit)
}

func (s *Handler) GetTokenExpiry() time.Duration {
	return s.tokenExpiry
}