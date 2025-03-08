package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create adds a new user to the database
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates a user's information
	Update(ctx context.Context, user *User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]*User, error)
}
