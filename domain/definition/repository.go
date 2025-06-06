package definition

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, definition *Definition) error
	Update(ctx context.Context, definition *Definition) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Definition, error)
	GetAll(ctx context.Context, limit, offset int) ([]*Definition, error)
	Search(ctx context.Context, searchTerm string, limit, offset int, definitionType string) ([]*Definition, error)
}
