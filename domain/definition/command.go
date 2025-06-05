package definition

type CreateDefinitionCommand struct {
	Name         string `json:"name" validate:"required"`
	Abbreviation string `json:"abbreviation" validate:"required"`
	Suffix       string `json:"suffix"`
}

type UpdateDefinitionCommand struct {
	ID           string `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Abbreviation string `json:"abbreviation" validate:"required"`
	Suffix       string `json:"suffix"`
}

type DeleteDefinitionCommand struct {
	ID string `json:"id" validate:"required"`
}
