package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo        Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
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

// Command Handlers
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

func (s *Handler) HandleLoginUserCommand(ctx context.Context, command LoginUserCommand) (*LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, command.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := user.ComparePassword(command.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: signedToken,
		User:  user,
	}, nil
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

// Query Handlers
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

// Legacy methods for backward compatibility (will be deprecated)
func (s *Handler) RegisterUser(ctx context.Context, email, password, firstName, lastName string) (*User, error) {
	command := RegisterUserCommand{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}
	return s.HandleRegisterUserCommand(ctx, command)
}

func (s *Handler) LoginUser(ctx context.Context, email, password string) (string, error) {
	command := LoginUserCommand{
		Email:    email,
		Password: password,
	}
	response, err := s.HandleLoginUserCommand(ctx, command)
	if err != nil {
		return "", err
	}
	return response.Token, nil
}

func (s *Handler) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := GetUserByIDQuery{
		ID: id.String(),
	}
	return s.HandleGetUserByIDQuery(ctx, query)
}

func (s *Handler) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := GetUserByEmailQuery{
		Email: email,
	}
	return s.HandleGetUserByEmailQuery(ctx, query)
}

func (s *Handler) UpdateUser(ctx context.Context, user *User) error {
	command := UpdateUserCommand{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	_, err := s.HandleUpdateUserCommand(ctx, command)
	return err
}

func (s *Handler) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	command := ChangePasswordCommand{
		UserID:      userID.String(),
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
	return s.HandleChangePasswordCommand(ctx, command)
}

func (s *Handler) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Note: This method doesn't validate password, which might be a security issue
	// For proper deletion, use HandleDeleteUserCommand instead
	return s.repo.Delete(ctx, userID)
}

func (s *Handler) ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	command := ValidateUserPasswordCommand{
		UserID:   userID.String(),
		Password: password,
	}
	return s.HandleValidateUserPasswordCommand(ctx, command)
}

func (s *Handler) GetTokenExpiry() time.Duration {
	return s.tokenExpiry
}
