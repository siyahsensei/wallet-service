package definition

type GetDefinitionByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

type GetAllDefinitionsQuery struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetDefinitionByAbbreviationQuery struct {
	Abbreviation string `json:"abbreviation" validate:"required"`
}

type SearchDefinitionsQuery struct {
	SearchTerm string `json:"searchTerm" validate:"required"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}
