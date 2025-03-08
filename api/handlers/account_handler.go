package handlers

import (
	"siyahsensei/wallet-service/domain/account"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountService *account.Service
}

func NewAccountHandler(accountService *account.Service) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

type CreateAccountRequest struct {
	Name        string              `json:"name" validate:"required"`
	Description string              `json:"description"`
	Type        account.AccountType `json:"type" validate:"required"`
	Institution string              `json:"institution"`
	Currency    string              `json:"currency" validate:"required"`
	Balance     float64             `json:"balance"`
	IsManual    bool                `json:"isManual"`
	Icon        string              `json:"icon"`
	Color       string              `json:"color"`
	APIKey      string              `json:"apiKey,omitempty"`
	APISecret   string              `json:"apiSecret,omitempty"`
}

type UpdateAccountRequest struct {
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Type        account.AccountType `json:"type,omitempty"`
	Institution string              `json:"institution,omitempty"`
	Currency    string              `json:"currency,omitempty"`
	Balance     *float64            `json:"balance,omitempty"`
	IsManual    *bool               `json:"isManual,omitempty"`
	Icon        string              `json:"icon,omitempty"`
	Color       string              `json:"color,omitempty"`
}

type APICredentialsRequest struct {
	APIKey    string `json:"apiKey" validate:"required"`
	APISecret string `json:"apiSecret" validate:"required"`
}

type AccountResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        account.AccountType `json:"type"`
	Institution string              `json:"institution"`
	Currency    string              `json:"currency"`
	Balance     float64             `json:"balance"`
	IsManual    bool                `json:"isManual"`
	Icon        string              `json:"icon"`
	Color       string              `json:"color"`
	IsConnected bool                `json:"isConnected"`
	LastSync    string              `json:"lastSync,omitempty"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

func toAccountResponse(a *account.Account) *AccountResponse {
	resp := &AccountResponse{
		ID:          a.ID.String(),
		Name:        a.Name,
		Description: a.Description,
		Type:        a.Type,
		Institution: a.Institution,
		Currency:    a.Currency,
		Balance:     a.Balance,
		IsManual:    a.IsManual,
		Icon:        a.Icon,
		Color:       a.Color,
		IsConnected: a.IsConnected,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   a.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if a.LastSync != nil {
		resp.LastSync = a.LastSync.Format("2006-01-02T15:04:05Z")
	}

	return resp
}

func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	newAccount, err := h.accountService.CreateAccount(
		c.Context(),
		userIDStr,
		req.Name,
		req.Description,
		req.Type,
		req.Institution,
		req.Currency,
		req.Balance,
		req.IsManual,
		req.Icon,
		req.Color,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if req.APIKey != "" && req.APISecret != "" {
		if err := h.accountService.SetAPICredentials(c.Context(), newAccount.ID, req.APIKey, req.APISecret); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to set API credentials",
			})
		}
		newAccount, _ = h.accountService.GetAccountByID(c.Context(), newAccount.ID)
	}

	return c.Status(fiber.StatusCreated).JSON(toAccountResponse(newAccount))
}

func (h *AccountHandler) GetUserAccounts(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountType := c.Query("type")
	var accounts []*account.Account
	var err error
	if accountType != "" {
		accounts, err = h.accountService.GetAccountsByType(c.Context(), userIDStr, account.AccountType(accountType))
	} else {
		accounts, err = h.accountService.GetUserAccounts(c.Context(), userIDStr)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve accounts",
		})
	}

	response := make([]*AccountResponse, len(accounts))
	for i, acc := range accounts {
		response[i] = toAccountResponse(acc)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accounts": response,
	})
}

func (h *AccountHandler) GetAccount(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountIDStr := c.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	acc, err := h.accountService.GetAccountByID(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	if acc.UserID != userIDStr {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to access this account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(toAccountResponse(acc))
}

func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountIDStr := c.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	acc, err := h.accountService.GetAccountByID(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	if acc.UserID != userIDStr {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this account",
		})
	}

	var req UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name != "" {
		acc.Name = req.Name
	}
	if req.Description != "" {
		acc.Description = req.Description
	}
	if req.Type != "" {
		acc.Type = req.Type
	}
	if req.Institution != "" {
		acc.Institution = req.Institution
	}
	if req.Currency != "" {
		acc.Currency = req.Currency
	}
	if req.Balance != nil {
		acc.UpdateBalance(*req.Balance)
	}
	if req.IsManual != nil {
		acc.IsManual = *req.IsManual
	}
	if req.Icon != "" {
		acc.Icon = req.Icon
	}
	if req.Color != "" {
		acc.Color = req.Color
	}

	if err := h.accountService.UpdateAccount(c.Context(), acc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(toAccountResponse(acc))
}

func (h *AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountIDStr := c.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	acc, err := h.accountService.GetAccountByID(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	if acc.UserID != userIDStr {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this account",
		})
	}

	if err := h.accountService.DeleteAccount(c.Context(), accountID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}

func (h *AccountHandler) SetAPICredentials(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountIDStr := c.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	acc, err := h.accountService.GetAccountByID(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	if acc.UserID != userIDStr {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this account",
		})
	}

	var req APICredentialsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.accountService.SetAPICredentials(c.Context(), accountID, req.APIKey, req.APISecret); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set API credentials",
		})
	}

	acc, _ = h.accountService.GetAccountByID(c.Context(), accountID)

	return c.Status(fiber.StatusOK).JSON(toAccountResponse(acc))
}

func (h *AccountHandler) GetAccountTypes(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account_types": h.accountService.GetAllAccountTypes(),
	})
}

func (h *AccountHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("/", authMiddleware, h.CreateAccount)
	router.Get("/", authMiddleware, h.GetUserAccounts)
	router.Get("/types", h.GetAccountTypes)
	router.Get("/:id", authMiddleware, h.GetAccount)
	router.Put("/:id", authMiddleware, h.UpdateAccount)
	router.Delete("/:id", authMiddleware, h.DeleteAccount)
	router.Post("/:id/credentials", authMiddleware, h.SetAPICredentials)
}
