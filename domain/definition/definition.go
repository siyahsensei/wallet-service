package definition

import (
	"time"

	"github.com/google/uuid"
)

type Definition struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Abbreviation string    `json:"abbreviation" db:"abbreviation"`
	Suffix       string    `json:"suffix" db:"suffix"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

func NewDefinition(command CreateDefinitionCommand) *Definition {
	now := time.Now()
	return &Definition{
		ID:           uuid.New(),
		Name:         command.Name,
		Abbreviation: command.Abbreviation,
		Suffix:       command.Suffix,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
