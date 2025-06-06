package routes

import (
	"strconv"

	"siyahsensei/wallet-service/domain/definition"

	"github.com/gofiber/fiber/v2"
)

type DefinitionHandler struct {
	definitionService *definition.Handler
}

func NewDefinitionHandler(definitionService *definition.Handler) *DefinitionHandler {
	return &DefinitionHandler{
		definitionService: definitionService,
	}
}

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

func toDefinitionResponse(d *definition.Definition) DefinitionResponse {
	return DefinitionResponse{
		ID:           d.ID.String(),
		Name:         d.Name,
		Abbreviation: d.Abbreviation,
		Suffix:       d.Suffix,
		CreatedAt:    d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateDefinition godoc
// @Summary Create a new definition
// @Description Create a new asset definition
// @Tags definitions
// @Accept json
// @Produce json
// @Param definition body CreateDefinitionRequest true "Definition creation data"
// @Success 201 {object} map[string]DefinitionResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /definitions [post]
func (h *DefinitionHandler) CreateDefinition(c *fiber.Ctx) error {
	var req CreateDefinitionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command := definition.CreateDefinitionCommand{
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
		Suffix:       req.Suffix,
	}

	createdDefinition, err := h.definitionService.HandleCreateDefinitionCommand(c.Context(), command)
	if err != nil {
		if err.Error() == "abbreviation already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"definition": toDefinitionResponse(createdDefinition),
	})
}

// UpdateDefinition godoc
// @Summary Update a definition
// @Description Update an existing asset definition
// @Tags definitions
// @Accept json
// @Produce json
// @Param id path string true "Definition ID"
// @Param definition body UpdateDefinitionRequest true "Definition update data"
// @Success 200 {object} map[string]DefinitionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /definitions/{id} [put]
func (h *DefinitionHandler) UpdateDefinition(c *fiber.Ctx) error {
	definitionID := c.Params("id")
	if definitionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Definition ID is required",
		})
	}

	var req UpdateDefinitionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command := definition.UpdateDefinitionCommand{
		ID:           definitionID,
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
		Suffix:       req.Suffix,
	}

	updatedDefinition, err := h.definitionService.HandleUpdateDefinitionCommand(c.Context(), command)
	if err != nil {
		if err.Error() == "definition not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "abbreviation already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"definition": toDefinitionResponse(updatedDefinition),
	})
}

// DeleteDefinition godoc
// @Summary Delete a definition
// @Description Delete an existing asset definition
// @Tags definitions
// @Accept json
// @Produce json
// @Param id path string true "Definition ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /definitions/{id} [delete]
func (h *DefinitionHandler) DeleteDefinition(c *fiber.Ctx) error {
	definitionID := c.Params("id")
	if definitionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Definition ID is required",
		})
	}

	command := definition.DeleteDefinitionCommand{
		ID: definitionID,
	}

	err := h.definitionService.HandleDeleteDefinitionCommand(c.Context(), command)
	if err != nil {
		if err.Error() == "definition not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetDefinitionByID godoc
// @Summary Get definition by ID
// @Description Get a specific asset definition by ID
// @Tags definitions
// @Accept json
// @Produce json
// @Param id path string true "Definition ID"
// @Success 200 {object} map[string]DefinitionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /definitions/{id} [get]
func (h *DefinitionHandler) GetDefinitionByID(c *fiber.Ctx) error {
	definitionID := c.Params("id")
	if definitionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Definition ID is required",
		})
	}

	query := definition.GetDefinitionByIDQuery{
		ID: definitionID,
	}

	foundDefinition, err := h.definitionService.HandleGetDefinitionByIDQuery(c.Context(), query)
	if err != nil {
		if err.Error() == "definition not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Definition not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"definition": toDefinitionResponse(foundDefinition),
	})
}

// GetAllDefinitions godoc
// @Summary Get all definitions
// @Description Get all asset definitions with optional pagination
// @Tags definitions
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} DefinitionsListResponse
// @Failure 500 {object} map[string]string
// @Router /definitions [get]
func (h *DefinitionHandler) GetAllDefinitions(c *fiber.Ctx) error {
	query := definition.GetAllDefinitionsQuery{}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			query.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			query.Offset = val
		}
	}

	definitions, err := h.definitionService.HandleGetAllDefinitionsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var definitionResponses []DefinitionResponse
	for _, d := range definitions {
		definitionResponses = append(definitionResponses, toDefinitionResponse(d))
	}

	return c.Status(fiber.StatusOK).JSON(DefinitionsListResponse{
		Definitions: definitionResponses,
		Total:       len(definitionResponses),
	})
}

// GetDefinitionByAbbreviation godoc
// @Summary Get definition by abbreviation
// @Description Get a specific asset definition by abbreviation
// @Tags definitions
// @Accept json
// @Produce json
// @Param abbreviation path string true "Definition Abbreviation"
// @Success 200 {object} map[string]DefinitionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /definitions/abbreviation/{abbreviation} [get]
func (h *DefinitionHandler) GetDefinitionByAbbreviation(c *fiber.Ctx) error {
	abbreviation := c.Params("abbreviation")
	if abbreviation == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Abbreviation is required",
		})
	}

	query := definition.GetDefinitionByAbbreviationQuery{
		Abbreviation: abbreviation,
	}

	foundDefinition, err := h.definitionService.HandleGetDefinitionByAbbreviationQuery(c.Context(), query)
	if err != nil {
		if err.Error() == "definition not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Definition not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"definition": toDefinitionResponse(foundDefinition),
	})
}

// SearchDefinitions godoc
// @Summary Search definitions
// @Description Search asset definitions by name or abbreviation
// @Tags definitions
// @Accept json
// @Produce json
// @Param q query string true "Search term"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} DefinitionsListResponse
// @Failure 400 {object} map[string]string
// @Router /definitions/search [get]
func (h *DefinitionHandler) SearchDefinitions(c *fiber.Ctx) error {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search term (q) is required",
		})
	}

	query := definition.SearchDefinitionsQuery{
		SearchTerm: searchTerm,
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			query.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			query.Offset = val
		}
	}

	definitions, err := h.definitionService.HandleSearchDefinitionsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var definitionResponses []DefinitionResponse
	for _, d := range definitions {
		definitionResponses = append(definitionResponses, toDefinitionResponse(d))
	}

	return c.Status(fiber.StatusOK).JSON(DefinitionsListResponse{
		Definitions: definitionResponses,
		Total:       len(definitionResponses),
	})
}

func (h *DefinitionHandler) RegisterRoutes(router fiber.Router) {
	definitionGroup := router.Group("/definitions")

	definitionGroup.Post("/", h.CreateDefinition)
	definitionGroup.Put("/:id", h.UpdateDefinition)
	definitionGroup.Delete("/:id", h.DeleteDefinition)
	definitionGroup.Get("/", h.GetAllDefinitions)
	definitionGroup.Get("/search", h.SearchDefinitions)
	definitionGroup.Get("/abbreviation/:abbreviation", h.GetDefinitionByAbbreviation)
	definitionGroup.Get("/:id", h.GetDefinitionByID)
}
