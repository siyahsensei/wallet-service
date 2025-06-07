package routes

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"siyahsensei/wallet-service/domain/asset"
	presentation "siyahsensei/wallet-service/presentation/asset"
)

type AssetHandler struct {
	assetService *asset.Handler
}

func NewAssetHandler(assetService *asset.Handler) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

func (h *AssetHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	assetGroup := router.Group("/assets", authMiddleware)

	assetGroup.Post("/", h.CreateAsset)
	assetGroup.Put("/:id", h.UpdateAsset)
	assetGroup.Delete("/:id", h.DeleteAsset)
	assetGroup.Get("/", h.GetUserAssets)
	assetGroup.Get("/filter", h.FilterAssets)
	assetGroup.Get("/:id", h.GetAssetByID)
}

// CreateAsset godoc
// @Summary Create a new asset
// @Description Create a new asset for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param asset body presentation.CreateAssetRequest true "Asset creation data"
// @Success 201 {object} map[string]presentation.AssetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /assets [post]
func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req presentation.CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Map request to command with UserID from JWT token
	command := asset.CreateAssetCommand{
		UserID:       userIDValue,
		AccountID:    req.AccountID,
		DefinitionID: req.DefinitionID,
		Type:         req.Type,
		Quantity:     req.Quantity,
		Notes:        req.Notes,
		PurchaseDate: req.PurchaseDate,
	}

	createdAsset, err := h.assetService.HandleCreateAssetCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"asset": presentation.ToAssetResponse(createdAsset),
	})
}

// UpdateAsset godoc
// @Summary Update an asset
// @Description Update an existing asset for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Param asset body presentation.UpdateAssetRequest true "Asset update data"
// @Success 200 {object} map[string]presentation.AssetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [put]
func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetID := c.Params("id")
	if assetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset ID is required",
		})
	}

	var req presentation.UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Map request to command with UserID from JWT token and ID from URL params
	command := asset.UpdateAssetCommand{
		ID:           assetID,
		UserID:       userIDValue,
		AccountID:    req.AccountID,
		DefinitionID: req.DefinitionID,
		Type:         req.Type,
		Quantity:     req.Quantity,
		Notes:        req.Notes,
		PurchaseDate: req.PurchaseDate,
	}

	updatedAsset, err := h.assetService.HandleUpdateAssetCommand(c.Context(), command)
	if err != nil {
		if err.Error() == "asset not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "unauthorized: asset does not belong to user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"asset": presentation.ToAssetResponse(updatedAsset),
	})
}

// DeleteAsset godoc
// @Summary Delete an asset
// @Description Delete an existing asset for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetID := c.Params("id")
	if assetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset ID is required",
		})
	}

	command := asset.DeleteAssetCommand{
		ID:     assetID,
		UserID: userIDValue.String(),
	}

	err := h.assetService.HandleDeleteAssetCommand(c.Context(), command)
	if err != nil {
		if err.Error() == "asset not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "unauthorized: asset does not belong to user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetAssetByID godoc
// @Summary Get asset by ID
// @Description Get a specific asset by ID for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Success 200 {object} map[string]presentation.AssetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [get]
func (h *AssetHandler) GetAssetByID(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetID := c.Params("id")
	if assetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset ID is required",
		})
	}

	query := asset.GetAssetByIDQuery{
		ID:     assetID,
		UserID: userIDValue.String(),
	}

	foundAsset, err := h.assetService.HandleGetAssetByIDQuery(c.Context(), query)
	if err != nil {
		if err.Error() == "asset not found" || err.Error() == "unauthorized: asset does not belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"asset": presentation.ToAssetResponse(foundAsset),
	})
}

// GetUserAssets godoc
// @Summary Get all user assets
// @Description Get all assets for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} presentation.AssetsListResponse
// @Failure 401 {object} map[string]string
// @Router /assets [get]
func (h *AssetHandler) GetUserAssets(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := asset.GetUserAssetsQuery{
		UserID: userIDValue.String(),
	}

	assets, err := h.assetService.HandleGetUserAssetsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var assetResponses []presentation.AssetResponse
	for _, a := range assets {
		assetResponses = append(assetResponses, presentation.ToAssetResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(presentation.AssetsListResponse{
		Assets: assetResponses,
		Total:  len(assetResponses),
	})
}

// FilterAssets godoc
// @Summary Filter assets
// @Description Filter assets with optional parameters for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param accountId query string false "Account ID"
// @Param assetType query string false "Asset Type"
// @Param minQuantity query number false "Minimum Quantity"
// @Param maxQuantity query number false "Maximum Quantity"
// @Param createdFrom query string false "Created From Date (RFC3339)"
// @Param createdTo query string false "Created To Date (RFC3339)"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} presentation.AssetsListResponse
// @Failure 401 {object} map[string]string
// @Router /assets/filter [get]
func (h *AssetHandler) FilterAssets(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := asset.FilterAssetsQuery{
		UserID: userIDValue.String(),
	}

	if accountID := c.Query("accountId"); accountID != "" {
		query.AccountID = &accountID
	}

	if assetType := c.Query("assetType"); assetType != "" {
		at := asset.AssetType(assetType)
		query.AssetType = &at
	}

	if minQuantity := c.Query("minQuantity"); minQuantity != "" {
		if val, err := strconv.ParseFloat(minQuantity, 64); err == nil {
			query.MinQuantity = &val
		}
	}

	if maxQuantity := c.Query("maxQuantity"); maxQuantity != "" {
		if val, err := strconv.ParseFloat(maxQuantity, 64); err == nil {
			query.MaxQuantity = &val
		}
	}

	if createdFrom := c.Query("createdFrom"); createdFrom != "" {
		if val, err := time.Parse(time.RFC3339, createdFrom); err == nil {
			query.CreatedFrom = &val
		}
	}

	if createdTo := c.Query("createdTo"); createdTo != "" {
		if val, err := time.Parse(time.RFC3339, createdTo); err == nil {
			query.CreatedTo = &val
		}
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

	assets, err := h.assetService.HandleFilterAssetsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var assetResponses []presentation.AssetResponse
	for _, a := range assets {
		assetResponses = append(assetResponses, presentation.ToAssetResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(presentation.AssetsListResponse{
		Assets: assetResponses,
		Total:  len(assetResponses),
	})
}
