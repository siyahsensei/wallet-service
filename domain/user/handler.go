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

func NewHandler(repo Repository, jwtSecret string, tokenExpiry time.Duration) *Handler {
	return &Handler{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

func (s *Handler) RegisterUser(ctx context.Context, email, password, firstName, lastName string) (*User, error) {
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	user, err := NewUser(email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Handler) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := user.ComparePassword(password); err != nil {
		return "", errors.New("invalid credentials")
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
		return "", err
	}
	return signedToken, nil
}

func (s *Handler) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Handler) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Handler) UpdateUser(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *Handler) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := user.ComparePassword(oldPassword); err != nil {
		return errors.New("incorrect current password")
	}
	if err := user.UpdatePassword(newPassword); err != nil {
		return err
	}
	return s.repo.Update(ctx, user)
}

func (s *Handler) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Delete(ctx, userID)
}

func (s *Handler) ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	return user.ComparePassword(password)
}

func (s *Handler) GetTokenExpiry() time.Duration {
	return s.tokenExpiry
}
