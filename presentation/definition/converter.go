package presentation

import "siyahsensei/wallet-service/domain/definition"

func ToDefinitionResponse(d *definition.Definition) DefinitionResponse {
	return DefinitionResponse{
		ID:           d.ID.String(),
		Name:         d.Name,
		Abbreviation: d.Abbreviation,
		Suffix:       d.Suffix,
		CreatedAt:    d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}