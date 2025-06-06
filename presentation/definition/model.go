package presentation

// Request models
type CreateDefinitionRequest struct {
	Name         string `json:"name" validate:"required"`
	Abbreviation string `json:"abbreviation" validate:"required"`
	Suffix       string `json:"suffix"`
}

type UpdateDefinitionRequest struct {
	Name         string `json:"name" validate:"required"`
	Abbreviation string `json:"abbreviation" validate:"required"`
	Suffix       string `json:"suffix"`
}

// Response models
type DefinitionResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Suffix       string `json:"suffix"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type DefinitionsListResponse struct {
	Definitions []DefinitionResponse `json:"definitions"`
	Total       int                  `json:"total"`
}