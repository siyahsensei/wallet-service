package routes

import (
	"strconv"

	"siyahsensei/wallet-service/domain/definition"
	presentation "siyahsensei/wallet-service/presentation/definition"

	"github.com/gofiber/fiber/v2"
)

type DefinitionRoute struct {
	definitionService *definition.Handler
}

func NewDefinitionRoute(definitionService *definition.Handler) *DefinitionRoute {
	return &DefinitionRoute{
		definitionService: definitionService,
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
func (r *DefinitionRoute) CreateDefinition(c *fiber.Ctx) error {
	var req presentation.CreateDefinitionRequest
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
	createdDefinition, err := r.definitionService.HandleCreateDefinitionCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"definition": presentation.ToDefinitionResponse(createdDefinition),
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
func (r *DefinitionRoute) UpdateDefinition(c *fiber.Ctx) error {
	definitionID := c.Params("id")
	if definitionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Definition ID is required",
		})
	}

	var req presentation.UpdateDefinitionRequest
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
	updatedDefinition, err := r.definitionService.HandleUpdateDefinitionCommand(c.Context(), command)
	if err != nil {
		if err.Error() == "definition not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"definition": presentation.ToDefinitionResponse(updatedDefinition),
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
func (r *DefinitionRoute) DeleteDefinition(c *fiber.Ctx) error {
	definitionID := c.Params("id")
	if definitionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Definition ID is required",
		})
	}

	command := definition.DeleteDefinitionCommand{
		ID: definitionID,
	}
	err := r.definitionService.HandleDeleteDefinitionCommand(c.Context(), command)
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
func (h *DefinitionRoute) GetDefinitionByID(c *fiber.Ctx) error {
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
		"definition": presentation.ToDefinitionResponse(foundDefinition),
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
// @Param type query string false "Definition type"
// @Success 200 {object} DefinitionsListResponse
// @Failure 500 {object} map[string]string
// @Router /definitions [get]
func (h *DefinitionRoute) GetAllDefinitions(c *fiber.Ctx) error {
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

	var definitionResponses []presentation.DefinitionResponse
	for _, d := range definitions {
		definitionResponses = append(definitionResponses, presentation.ToDefinitionResponse(d))
	}

	return c.Status(fiber.StatusOK).JSON(presentation.DefinitionsListResponse{
		Definitions: definitionResponses,
		Total:       len(definitionResponses),
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
func (h *DefinitionRoute) SearchDefinitions(c *fiber.Ctx) error {
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

	if definitionType := c.Query("type"); definitionType != "" {
		query.DefinitionType = definitionType
	}

	definitions, err := h.definitionService.HandleSearchDefinitionsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var definitionResponses []presentation.DefinitionResponse
	for _, d := range definitions {
		definitionResponses = append(definitionResponses, presentation.ToDefinitionResponse(d))
	}

	return c.Status(fiber.StatusOK).JSON(presentation.DefinitionsListResponse{
		Definitions: definitionResponses,
		Total:       len(definitionResponses),
	})
}

func (h *DefinitionRoute) RegisterRoutes(router fiber.Router) {
	definitionGroup := router.Group("/definitions")

	definitionGroup.Post("/", h.CreateDefinition)
	definitionGroup.Put("/:id", h.UpdateDefinition)
	definitionGroup.Delete("/:id", h.DeleteDefinition)
	definitionGroup.Get("/:id", h.GetDefinitionByID)
	definitionGroup.Get("/", h.GetAllDefinitions)
	definitionGroup.Get("/search", h.SearchDefinitions)
}
