package routes

import (
	"strconv"
	"time"

	"siyahsensei/wallet-service/domain/asset"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService *asset.Handler
}

func NewAssetHandler(assetService *asset.Handler) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

type AssetResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	AccountID    string    `json:"accountId"`
	DefinitionID string    `json:"definitionId"`
	Type         string    `json:"type"`
	Quantity     float64   `json:"quantity"`
	Notes        string    `json:"notes"`
	PurchaseDate time.Time `json:"purchaseDate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AssetsListResponse struct {
	Assets []AssetResponse `json:"assets"`
	Total  int             `json:"total"`
}

type AssetPerformanceResponse struct {
	AssetID        string  `json:"assetId"`
	Name           string  `json:"name"`
	Symbol         string  `json:"symbol"`
	Type           string  `json:"type"`
	InitialValue   float64 `json:"initialValue"`
	CurrentValue   float64 `json:"currentValue"`
	ProfitLoss     float64 `json:"profitLoss"`
	ProfitLossPerc float64 `json:"profitLossPercentage"`
	Currency       string  `json:"currency"`
}

type TotalValueResponse struct {
	TotalValue float64 `json:"totalValue"`
}

func toAssetResponse(a *asset.Asset) AssetResponse {
	return AssetResponse{
		ID:           a.ID.String(),
		UserID:       a.UserID.String(),
		AccountID:    a.AccountID.String(),
		DefinitionID: a.DefinitionID.String(),
		Type:         string(a.Type),
		Quantity:     a.Quantity,
		Notes:        a.Notes,
		PurchaseDate: a.PurchaseDate,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func toAssetPerformanceResponse(ap *asset.AssetPerformance) AssetPerformanceResponse {
	return AssetPerformanceResponse{
		AssetID:        ap.AssetID.String(),
		Name:           ap.Name,
		Symbol:         ap.Symbol,
		Type:           string(ap.Type),
		InitialValue:   ap.InitialValue,
		CurrentValue:   ap.CurrentValue,
		ProfitLoss:     ap.ProfitLoss,
		ProfitLossPerc: ap.ProfitLossPerc,
		Currency:       ap.Currency,
	}
}

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var command asset.CreateAssetCommand
	if err := c.BodyParser(&command); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command.UserID = userIDValue.String()

	createdAsset, err := h.assetService.HandleCreateAssetCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"asset": toAssetResponse(createdAsset),
	})
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
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

	var command asset.UpdateAssetCommand
	if err := c.BodyParser(&command); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command.ID = assetID
	command.UserID = userIDValue.String()

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
		"asset": toAssetResponse(updatedAsset),
	})
}

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
		"asset": toAssetResponse(foundAsset),
	})
}

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

	var assetResponses []AssetResponse
	for _, a := range assets {
		assetResponses = append(assetResponses, toAssetResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AssetsListResponse{
		Assets: assetResponses,
		Total:  len(assetResponses),
	})
}

func (h *AssetHandler) GetAccountAssets(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID := c.Params("accountId")
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
		})
	}

	query := asset.GetAccountAssetsQuery{
		AccountID: accountID,
		UserID:    userIDValue.String(),
	}

	assets, err := h.assetService.HandleGetAccountAssetsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var assetResponses []AssetResponse
	for _, a := range assets {
		assetResponses = append(assetResponses, toAssetResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AssetsListResponse{
		Assets: assetResponses,
		Total:  len(assetResponses),
	})
}

func (h *AssetHandler) GetAssetsByType(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetType := c.Params("type")
	if assetType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset type is required",
		})
	}

	query := asset.GetAssetsByTypeQuery{
		UserID:    userIDValue.String(),
		AssetType: asset.AssetType(assetType),
	}

	assets, err := h.assetService.HandleGetAssetsByTypeQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var assetResponses []AssetResponse
	for _, a := range assets {
		assetResponses = append(assetResponses, toAssetResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AssetsListResponse{
		Assets: assetResponses,
		Total:  len(assetResponses),
	})
}

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

	var assetResponses []AssetResponse
	for _, a := range assets {
		assetResponses = append(assetResponses, toAssetResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AssetsListResponse{
		Assets: assetResponses,
		Total:  len(assetResponses),
	})
}

func (h *AssetHandler) GetAssetPerformance(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Start date and end date are required",
		})
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start date format. Use RFC3339 format",
		})
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end date format. Use RFC3339 format",
		})
	}

	query := asset.GetAssetPerformanceQuery{
		UserID:    userIDValue.String(),
		StartDate: startDate,
		EndDate:   endDate,
	}

	performances, err := h.assetService.HandleGetAssetPerformanceQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var performanceResponses []AssetPerformanceResponse
	for _, p := range performances {
		performanceResponses = append(performanceResponses, toAssetPerformanceResponse(p))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"performances": performanceResponses,
		"total":        len(performanceResponses),
	})
}

func (h *AssetHandler) GetTotalValue(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := asset.GetTotalValueQuery{
		UserID: userIDValue.String(),
	}

	if assetTypesStr := c.Query("assetTypes"); assetTypesStr != "" {
		query.AssetTypes = []asset.AssetType{asset.AssetType(assetTypesStr)}
	}

	totalValue, err := h.assetService.HandleGetTotalValueQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(TotalValueResponse{
		TotalValue: totalValue,
	})
}

func (h *AssetHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	assetGroup := router.Group("/assets", authMiddleware)

	assetGroup.Post("/", h.CreateAsset)
	assetGroup.Put("/:id", h.UpdateAsset)
	assetGroup.Delete("/:id", h.DeleteAsset)
	assetGroup.Get("/", h.GetUserAssets)
	assetGroup.Get("/filter", h.FilterAssets)
	assetGroup.Get("/performance", h.GetAssetPerformance)
	assetGroup.Get("/total-value", h.GetTotalValue)
	assetGroup.Get("/account/:accountId", h.GetAccountAssets)
	assetGroup.Get("/type/:type", h.GetAssetsByType)
	assetGroup.Get("/:id", h.GetAssetByID)
}
