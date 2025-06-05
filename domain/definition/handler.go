package definition

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleCreateDefinitionCommand(ctx context.Context, command CreateDefinitionCommand) (*Definition, error) {
	if command.Name == "" {
		return nil, errors.New("name is required")
	}
	if command.Abbreviation == "" {
		return nil, errors.New("abbreviation is required")
	}

	// Check if abbreviation already exists
	existingDef, _ := h.repo.GetByAbbreviation(ctx, command.Abbreviation)
	if existingDef != nil {
		return nil, errors.New("abbreviation already exists")
	}

	definition := NewDefinition(command)
	if err := h.repo.Create(ctx, definition); err != nil {
		return nil, err
	}
	return definition, nil
}

func (h *Handler) HandleUpdateDefinitionCommand(ctx context.Context, command UpdateDefinitionCommand) (*Definition, error) {
	if command.Name == "" {
		return nil, errors.New("name is required")
	}
	if command.Abbreviation == "" {
		return nil, errors.New("abbreviation is required")
	}

	definitionID, err := uuid.Parse(command.ID)
	if err != nil {
		return nil, errors.New("invalid definition ID")
	}

	existingDefinition, err := h.repo.GetByID(ctx, definitionID)
	if err != nil {
		return nil, errors.New("definition not found")
	}

	// Check if abbreviation already exists for another definition
	existingByAbbr, _ := h.repo.GetByAbbreviation(ctx, command.Abbreviation)
	if existingByAbbr != nil && existingByAbbr.ID != definitionID {
		return nil, errors.New("abbreviation already exists")
	}

	existingDefinition.Name = command.Name
	existingDefinition.Abbreviation = command.Abbreviation
	existingDefinition.Suffix = command.Suffix
	existingDefinition.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, existingDefinition); err != nil {
		return nil, err
	}

	return existingDefinition, nil
}

func (h *Handler) HandleDeleteDefinitionCommand(ctx context.Context, command DeleteDefinitionCommand) error {
	definitionID, err := uuid.Parse(command.ID)
	if err != nil {
		return errors.New("invalid definition ID")
	}

	_, err = h.repo.GetByID(ctx, definitionID)
	if err != nil {
		return errors.New("definition not found")
	}

	return h.repo.Delete(ctx, definitionID)
}

func (h *Handler) HandleGetDefinitionByIDQuery(ctx context.Context, query GetDefinitionByIDQuery) (*Definition, error) {
	definitionID, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, errors.New("invalid definition ID")
	}

	return h.repo.GetByID(ctx, definitionID)
}

func (h *Handler) HandleGetAllDefinitionsQuery(ctx context.Context, query GetAllDefinitionsQuery) ([]*Definition, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	return h.repo.GetAll(ctx, limit, query.Offset)
}

func (h *Handler) HandleGetDefinitionByAbbreviationQuery(ctx context.Context, query GetDefinitionByAbbreviationQuery) (*Definition, error) {
	if query.Abbreviation == "" {
		return nil, errors.New("abbreviation is required")
	}

	return h.repo.GetByAbbreviation(ctx, query.Abbreviation)
}

func (h *Handler) HandleSearchDefinitionsQuery(ctx context.Context, query SearchDefinitionsQuery) ([]*Definition, error) {
	if query.SearchTerm == "" {
		return nil, errors.New("search term is required")
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	return h.repo.Search(ctx, query.SearchTerm, limit, query.Offset)
}
