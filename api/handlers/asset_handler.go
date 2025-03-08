package handlers

import (
	"time"

	"siyahsensei/wallet-service/domain/asset"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService *asset.Service
}

func NewAssetHandler(assetService *asset.Service) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

type CreateAssetRequest struct {
	AccountID     string          `json:"accountId" validate:"required"`
	Name          string          `json:"name" validate:"required"`
	Type          asset.AssetType `json:"type" validate:"required"`
	Symbol        string          `json:"symbol"`
	Quantity      float64         `json:"quantity" validate:"required"`
	PurchasePrice float64         `json:"purchasePrice" validate:"required"`
	CurrentPrice  float64         `json:"currentPrice" validate:"required"`
	Currency      string          `json:"currency" validate:"required"`
	Notes         string          `json:"notes"`
	PurchaseDate  string          `json:"purchaseDate" validate:"required"`
}

type UpdateAssetRequest struct {
	Name          string          `json:"name,omitempty"`
	Type          asset.AssetType `json:"type,omitempty"`
	Symbol        string          `json:"symbol,omitempty"`
	Quantity      *float64        `json:"quantity,omitempty"`
	PurchasePrice *float64        `json:"purchasePrice,omitempty"`
	CurrentPrice  *float64        `json:"currentPrice,omitempty"`
	Currency      string          `json:"currency,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	PurchaseDate  string          `json:"purchaseDate,omitempty"`
}

type AssetResponse struct {
	ID                   string          `json:"id"`
	AccountID            string          `json:"accountId"`
	Name                 string          `json:"name"`
	Type                 asset.AssetType `json:"type"`
	Symbol               string          `json:"symbol"`
	Quantity             float64         `json:"quantity"`
	PurchasePrice        float64         `json:"purchasePrice"`
	CurrentPrice         float64         `json:"currentPrice"`
	Currency             string          `json:"currency"`
	Notes                string          `json:"notes"`
	PurchaseDate         string          `json:"purchaseDate"`
	CurrentValue         float64         `json:"currentValue"`
	ProfitLoss           float64         `json:"profitLoss"`
	ProfitLossPercentage float64         `json:"profitLossPercentage"`
	LastUpdated          string          `json:"lastUpdated"`
	CreatedAt            string          `json:"createdAt"`
	UpdatedAt            string          `json:"updatedAt"`
}

func toAssetResponse(a *asset.Asset) *AssetResponse {
	return &AssetResponse{
		ID:                   a.ID.String(),
		AccountID:            a.AccountID.String(),
		Name:                 a.Name,
		Type:                 a.Type,
		Symbol:               a.Symbol,
		Quantity:             a.Quantity,
		PurchasePrice:        a.PurchasePrice,
		CurrentPrice:         a.CurrentPrice,
		Currency:             a.Currency,
		Notes:                a.Notes,
		PurchaseDate:         a.PurchaseDate.Format("2006-01-02"),
		CurrentValue:         a.CurrentValue(),
		ProfitLoss:           a.ProfitLoss(),
		ProfitLossPercentage: a.ProfitLossPercentage(),
		LastUpdated:          a.LastUpdated.Format("2006-01-02T15:04:05Z"),
		CreatedAt:            a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:            a.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid purchase date format, expected YYYY-MM-DD",
		})
	}

	newAsset, err := h.assetService.CreateAsset(
		c.Context(),
		userIDValue,
		accountID,
		req.Name,
		req.Type,
		req.Symbol,
		req.Quantity,
		req.PurchasePrice,
		req.CurrentPrice,
		req.Currency,
		req.Notes,
		purchaseDate,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(toAssetResponse(newAsset))
}

func (h *AssetHandler) GetUserAssets(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetType := c.Query("type")
	var assets []*asset.Asset
	var err error
	if assetType != "" {
		assets, err = h.assetService.GetAssetsByType(c.Context(), userIDValue, asset.AssetType(assetType))
	} else {
		assets, err = h.assetService.GetUserAssets(c.Context(), userIDValue)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve assets",
		})
	}

	response := make([]*AssetResponse, len(assets))
	for i, a := range assets {
		response[i] = toAssetResponse(a)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets": response,
	})
}

func (h *AssetHandler) GetAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetIDStr := c.Params("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	a, err := h.assetService.GetAssetByID(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	if a.UserID != userIDValue {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to access this asset",
		})
	}

	return c.Status(fiber.StatusOK).JSON(toAssetResponse(a))
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetIDStr := c.Params("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	a, err := h.assetService.GetAssetByID(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	if a.UserID != userIDValue {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this asset",
		})
	}

	var req UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name != "" {
		a.Name = req.Name
	}
	if req.Type != "" {
		a.Type = req.Type
	}
	if req.Symbol != "" {
		a.Symbol = req.Symbol
	}
	if req.Quantity != nil {
		if req.PurchasePrice != nil {
			a.UpdateQuantity(*req.Quantity, *req.PurchasePrice)
		} else {
			a.UpdateQuantity(*req.Quantity, 0)
		}
	}
	if req.CurrentPrice != nil {
		a.UpdatePrice(*req.CurrentPrice)
	}
	if req.Currency != "" {
		a.Currency = req.Currency
	}
	if req.Notes != "" {
		a.Notes = req.Notes
	}
	if req.PurchaseDate != "" {
		purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid purchase date format, expected YYYY-MM-DD",
			})
		}
		a.PurchaseDate = purchaseDate
	}

	if err := h.assetService.UpdateAsset(c.Context(), a); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(toAssetResponse(a))
}

func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetIDStr := c.Params("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	a, err := h.assetService.GetAssetByID(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	if a.UserID != userIDValue {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this asset",
		})
	}

	if err := h.assetService.DeleteAsset(c.Context(), assetID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete asset",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Asset deleted successfully",
	})
}

func (h *AssetHandler) GetAssetTypes(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"asset_types": h.assetService.GetAllAssetTypes(),
	})
}

func (h *AssetHandler) GetAccountAssets(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountIDStr := c.Params("accountId")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	assets, err := h.assetService.GetAccountAssets(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve assets",
		})
	}

	var userAssets []*asset.Asset
	for _, a := range assets {
		if a.UserID == userIDValue {
			userAssets = append(userAssets, a)
		}
	}

	response := make([]*AssetResponse, len(userAssets))
	for i, a := range userAssets {
		response[i] = toAssetResponse(a)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets": response,
	})
}

func (h *AssetHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("/", authMiddleware, h.CreateAsset)
	router.Get("/", authMiddleware, h.GetUserAssets)
	router.Get("/types", h.GetAssetTypes)
	router.Get("/:id", authMiddleware, h.GetAsset)
	router.Put("/:id", authMiddleware, h.UpdateAsset)
	router.Delete("/:id", authMiddleware, h.DeleteAsset)
	router.Get("/account/:accountId", authMiddleware, h.GetAccountAssets)
}
