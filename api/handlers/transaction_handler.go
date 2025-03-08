package handlers

import (
	"time"

	"siyahsensei/wallet-service/domain/transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *transaction.Service
}

func NewTransactionHandler(transactionService *transaction.Service) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

type CreateTransactionRequest struct {
	AccountID       string                      `json:"accountId" validate:"required"`
	AssetID         string                      `json:"assetId,omitempty"`
	Type            transaction.TransactionType `json:"type" validate:"required"`
	Amount          float64                     `json:"amount" validate:"required"`
	Quantity        float64                     `json:"quantity,omitempty"`
	Price           float64                     `json:"price,omitempty"`
	Fee             float64                     `json:"fee,omitempty"`
	Currency        string                      `json:"currency" validate:"required"`
	Description     string                      `json:"description,omitempty"`
	Category        string                      `json:"category,omitempty"`
	Date            string                      `json:"date" validate:"required"`
	ToAccountID     string                      `json:"toAccountId,omitempty"`
	TransactionHash string                      `json:"transactionHash,omitempty"`
}

type UpdateTransactionRequest struct {
	AccountID       string                      `json:"accountId,omitempty"`
	AssetID         string                      `json:"assetId,omitempty"`
	Type            transaction.TransactionType `json:"type,omitempty"`
	Amount          *float64                    `json:"amount,omitempty"`
	Quantity        *float64                    `json:"quantity,omitempty"`
	Price           *float64                    `json:"price,omitempty"`
	Fee             *float64                    `json:"fee,omitempty"`
	Currency        string                      `json:"currency,omitempty"`
	Description     string                      `json:"description,omitempty"`
	Category        string                      `json:"category,omitempty"`
	Date            string                      `json:"date,omitempty"`
	ToAccountID     string                      `json:"toAccountId,omitempty"`
	TransactionHash string                      `json:"transactionHash,omitempty"`
}

type TransactionResponse struct {
	ID              string                      `json:"id"`
	AccountID       string                      `json:"accountId"`
	AssetID         string                      `json:"assetId,omitempty"`
	Type            transaction.TransactionType `json:"type"`
	Amount          float64                     `json:"amount"`
	TotalAmount     float64                     `json:"totalAmount"`
	Quantity        float64                     `json:"quantity,omitempty"`
	Price           float64                     `json:"price,omitempty"`
	Fee             float64                     `json:"fee"`
	Currency        string                      `json:"currency"`
	Description     string                      `json:"description,omitempty"`
	Category        string                      `json:"category,omitempty"`
	Date            string                      `json:"date"`
	ToAccountID     string                      `json:"toAccountId,omitempty"`
	TransactionHash string                      `json:"transactionHash,omitempty"`
	IsCredit        bool                        `json:"isCredit"`
	IsDebit         bool                        `json:"isDebit"`
	CreatedAt       string                      `json:"createdAt"`
	UpdatedAt       string                      `json:"updatedAt"`
}

func toTransactionResponse(t *transaction.Transaction) *TransactionResponse {
	response := &TransactionResponse{
		ID:              t.ID.String(),
		AccountID:       t.AccountID.String(),
		Type:            t.Type,
		Amount:          t.Amount,
		TotalAmount:     t.TotalAmount(),
		Quantity:        t.Quantity,
		Price:           t.Price,
		Fee:             t.Fee,
		Currency:        t.Currency,
		Description:     t.Description,
		Category:        t.Category,
		Date:            t.Date.Format("2006-01-02"),
		TransactionHash: t.TransactionHash,
		IsCredit:        t.IsCredit(),
		IsDebit:         t.IsDebit(),
		CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if t.AssetID != nil {
		response.AssetID = t.AssetID.String()
	}

	if t.ToAccountID != nil {
		response.ToAccountID = t.ToAccountID.String()
	}

	return response
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	var req CreateTransactionRequest
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
	var assetID *uuid.UUID
	if req.AssetID != "" {
		parsedAssetID, err := uuid.Parse(req.AssetID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid asset ID",
			})
		}
		assetID = &parsedAssetID
	}

	var toAccountID *uuid.UUID
	if req.ToAccountID != "" {
		parsedToAccountID, err := uuid.Parse(req.ToAccountID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid destination account ID",
			})
		}
		toAccountID = &parsedToAccountID
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format, expected YYYY-MM-DD",
		})
	}

	newTransaction, err := h.transactionService.CreateTransaction(
		c.Context(),
		userIDValue,
		accountID,
		assetID,
		req.Type,
		req.Amount,
		req.Quantity,
		req.Price,
		req.Fee,
		req.Currency,
		req.Description,
		req.Category,
		date,
		toAccountID,
		req.TransactionHash,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(toTransactionResponse(newTransaction))
}

func (h *TransactionHandler) GetUserTransactions(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit, offset := getPaginationParams(c)

	transactions, err := h.transactionService.GetUserTransactions(c.Context(), userIDValue, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	response := make([]*TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = toTransactionResponse(t)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": response,
	})
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactionIDStr := c.Params("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	t, err := h.transactionService.GetTransactionByID(c.Context(), transactionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if t.UserID != userIDValue {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to access this transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(toTransactionResponse(t))
}

func (h *TransactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactionIDStr := c.Params("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	t, err := h.transactionService.GetTransactionByID(c.Context(), transactionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if t.UserID != userIDValue {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this transaction",
		})
	}

	var req UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.AccountID != "" {
		accountID, err := uuid.Parse(req.AccountID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid account ID",
			})
		}
		t.AccountID = accountID
	}

	if req.AssetID != "" {
		assetID, err := uuid.Parse(req.AssetID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid asset ID",
			})
		}
		t.AssetID = &assetID
	}

	if req.Type != "" {
		t.Type = req.Type
	}

	if req.Amount != nil {
		t.Amount = *req.Amount
	}

	if req.Quantity != nil {
		t.Quantity = *req.Quantity
	}

	if req.Price != nil {
		t.Price = *req.Price
	}

	if req.Fee != nil {
		t.Fee = *req.Fee
	}

	if req.Currency != "" {
		t.Currency = req.Currency
	}

	if req.Description != "" {
		t.Description = req.Description
	}

	if req.Category != "" {
		t.Category = req.Category
	}

	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format, expected YYYY-MM-DD",
			})
		}
		t.Date = date
	}

	if req.ToAccountID != "" {
		toAccountID, err := uuid.Parse(req.ToAccountID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid destination account ID",
			})
		}
		t.ToAccountID = &toAccountID
	}

	if req.TransactionHash != "" {
		t.TransactionHash = req.TransactionHash
	}

	if err := h.transactionService.UpdateTransaction(c.Context(), t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(toTransactionResponse(t))
}

func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactionIDStr := c.Params("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	t, err := h.transactionService.GetTransactionByID(c.Context(), transactionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if t.UserID != userIDValue {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this transaction",
		})
	}

	if err := h.transactionService.DeleteTransaction(c.Context(), transactionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction deleted successfully",
	})
}

func (h *TransactionHandler) GetTransactionsByAccount(c *fiber.Ctx) error {
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

	limit, offset := getPaginationParams(c)

	transactions, err := h.transactionService.GetAccountTransactions(c.Context(), accountID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	var userTransactions []*transaction.Transaction
	for _, t := range transactions {
		if t.UserID == userIDValue {
			userTransactions = append(userTransactions, t)
		}
	}

	response := make([]*TransactionResponse, len(userTransactions))
	for i, t := range userTransactions {
		response[i] = toTransactionResponse(t)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": response,
	})
}

func (h *TransactionHandler) GetTransactionsByAsset(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	assetIDStr := c.Params("assetId")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	limit, offset := getPaginationParams(c)

	transactions, err := h.transactionService.GetAssetTransactions(c.Context(), assetID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	var userTransactions []*transaction.Transaction
	for _, t := range transactions {
		if t.UserID == userIDValue {
			userTransactions = append(userTransactions, t)
		}
	}

	response := make([]*TransactionResponse, len(userTransactions))
	for i, t := range userTransactions {
		response[i] = toTransactionResponse(t)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": response,
	})
}

func (h *TransactionHandler) GetTransactionTypes(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transaction_types": h.transactionService.GetAllTransactionTypes(),
	})
}

func getPaginationParams(c *fiber.Ctx) (limit, offset int) {
	limit = 20
	offset = 0

	if limitParam := c.Query("limit"); limitParam != "" {
		parsedLimit := c.QueryInt("limit", 20)
		limit = parsedLimit
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		parsedOffset := c.QueryInt("offset", 0)
		offset = parsedOffset
	}

	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func (h *TransactionHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("/", authMiddleware, h.CreateTransaction)
	router.Get("/", authMiddleware, h.GetUserTransactions)
	router.Get("/types", h.GetTransactionTypes)
	router.Get("/:id", authMiddleware, h.GetTransaction)
	router.Put("/:id", authMiddleware, h.UpdateTransaction)
	router.Delete("/:id", authMiddleware, h.DeleteTransaction)
	router.Get("/account/:accountId", authMiddleware, h.GetTransactionsByAccount)
	router.Get("/asset/:assetId", authMiddleware, h.GetTransactionsByAsset)
}
