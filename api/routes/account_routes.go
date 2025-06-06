package routes

import (
	"strconv"
	"time"

	"siyahsensei/wallet-service/domain/account"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountService *account.Handler
}

func NewAccountHandler(accountService *account.Handler) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// Request models (without UserID for security)
type CreateAccountRequest struct {
	Name         string              `json:"name" validate:"required"`
	AccountType  account.AccountType `json:"accountType" validate:"required"`
	Balance      float64             `json:"balance"`
	CurrencyCode string              `json:"currencyCode" validate:"required"`
}

type UpdateAccountRequest struct {
	Name         string              `json:"name" validate:"required"`
	AccountType  account.AccountType `json:"accountType" validate:"required"`
	Balance      float64             `json:"balance"`
	CurrencyCode string              `json:"currencyCode" validate:"required"`
}

type UpdateBalanceRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

// Response models
type AccountResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Name         string    `json:"name"`
	AccountType  string    `json:"accountType"`
	Balance      float64   `json:"balance"`
	CurrencyCode string    `json:"currencyCode"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AccountsListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int               `json:"total"`
}

type AccountSummaryResponse struct {
	TotalAccounts int                         `json:"totalAccounts"`
	TotalBalance  float64                     `json:"totalBalance"`
	ByType        map[account.AccountType]int `json:"byType"`
	ByCurrency    map[string]float64          `json:"byCurrency"`
}

func toAccountResponse(a *account.Account) AccountResponse {
	return AccountResponse{
		ID:           a.ID.String(),
		UserID:       a.UserID.String(),
		Name:         a.Name,
		AccountType:  string(a.AccountType),
		Balance:      a.Balance,
		CurrencyCode: a.CurrencyCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func toAccountSummaryResponse(s *account.AccountSummary) AccountSummaryResponse {
	return AccountSummaryResponse{
		TotalAccounts: s.TotalAccounts,
		TotalBalance:  s.TotalBalance,
		ByType:        s.ByType,
		ByCurrency:    s.ByCurrency,
	}
}

func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
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

	// Map request to command with UserID from JWT token
	command := account.CreateAccountCommand{
		UserID:       userIDValue.String(),
		Name:         req.Name,
		AccountType:  req.AccountType,
		Balance:      req.Balance,
		CurrencyCode: req.CurrencyCode,
	}

	createdAccount, err := h.accountService.HandleCreateAccountCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"account": toAccountResponse(createdAccount),
	})
}

func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID := c.Params("id")
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
		})
	}

	var req UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command := account.UpdateAccountCommand{
		ID:           accountID,
		UserID:       userIDValue.String(),
		Name:         req.Name,
		AccountType:  req.AccountType,
		Balance:      req.Balance,
		CurrencyCode: req.CurrencyCode,
	}

	updatedAccount, err := h.accountService.HandleUpdateAccountCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account": toAccountResponse(updatedAccount),
	})
}

func (h *AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID := c.Params("id")
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
		})
	}

	command := account.DeleteAccountCommand{
		ID:     accountID,
		UserID: userIDValue.String(),
	}

	err := h.accountService.HandleDeleteAccountCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}

func (h *AccountHandler) GetAccountByID(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID := c.Params("id")
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
		})
	}

	query := account.GetAccountByIDQuery{
		ID:     accountID,
		UserID: userIDValue.String(),
	}

	acc, err := h.accountService.HandleGetAccountByIDQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account": toAccountResponse(acc),
	})
}

func (h *AccountHandler) GetUserAccounts(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := account.GetUserAccountsQuery{
		UserID: userIDValue.String(),
	}

	accounts, err := h.accountService.HandleGetUserAccountsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var accountResponses []AccountResponse
	for _, a := range accounts {
		accountResponses = append(accountResponses, toAccountResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AccountsListResponse{
		Accounts: accountResponses,
		Total:    len(accountResponses),
	})
}

func (h *AccountHandler) GetAccountsByType(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountType := c.Params("type")
	if accountType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account type is required",
		})
	}

	query := account.GetAccountsByTypeQuery{
		UserID:      userIDValue.String(),
		AccountType: account.AccountType(accountType),
	}

	accounts, err := h.accountService.HandleGetAccountsByTypeQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var accountResponses []AccountResponse
	for _, a := range accounts {
		accountResponses = append(accountResponses, toAccountResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AccountsListResponse{
		Accounts: accountResponses,
		Total:    len(accountResponses),
	})
}

func (h *AccountHandler) GetAccountsByCurrency(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	currencyCode := c.Params("currency")
	if currencyCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Currency code is required",
		})
	}

	query := account.GetAccountsByCurrencyQuery{
		UserID:       userIDValue.String(),
		CurrencyCode: currencyCode,
	}

	accounts, err := h.accountService.HandleGetAccountsByCurrencyQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var accountResponses []AccountResponse
	for _, a := range accounts {
		accountResponses = append(accountResponses, toAccountResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AccountsListResponse{
		Accounts: accountResponses,
		Total:    len(accountResponses),
	})
}

func (h *AccountHandler) FilterAccounts(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := account.FilterAccountsQuery{
		UserID: userIDValue.String(),
	}

	if accountType := c.Query("accountType"); accountType != "" {
		at := account.AccountType(accountType)
		query.AccountType = &at
	}

	if currencyCode := c.Query("currencyCode"); currencyCode != "" {
		query.CurrencyCode = &currencyCode
	}

	if minBalance := c.Query("minBalance"); minBalance != "" {
		if val, err := strconv.ParseFloat(minBalance, 64); err == nil {
			query.MinBalance = &val
		}
	}

	if maxBalance := c.Query("maxBalance"); maxBalance != "" {
		if val, err := strconv.ParseFloat(maxBalance, 64); err == nil {
			query.MaxBalance = &val
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

	accounts, err := h.accountService.HandleFilterAccountsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var accountResponses []AccountResponse
	for _, a := range accounts {
		accountResponses = append(accountResponses, toAccountResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AccountsListResponse{
		Accounts: accountResponses,
		Total:    len(accountResponses),
	})
}

func (h *AccountHandler) UpdateAccountBalance(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID := c.Params("id")
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
		})
	}

	var req UpdateBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command := account.UpdateAccountBalanceCommand{
		ID:     accountID,
		UserID: userIDValue.String(),
		Amount: req.Amount,
	}

	updatedAccount, err := h.accountService.HandleUpdateAccountBalanceCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account": toAccountResponse(updatedAccount),
	})
}

func (h *AccountHandler) GetAccountSummary(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := account.GetAccountSummaryQuery{
		UserID: userIDValue.String(),
	}

	summary, err := h.accountService.HandleGetAccountSummaryQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"summary": toAccountSummaryResponse(summary),
	})
}

func (h *AccountHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	accountGroup := router.Group("/accounts", authMiddleware)

	accountGroup.Post("/", h.CreateAccount)
	accountGroup.Put("/:id", h.UpdateAccount)
	accountGroup.Delete("/:id", h.DeleteAccount)
	accountGroup.Get("/", h.GetUserAccounts)
	accountGroup.Get("/filter", h.FilterAccounts)
	accountGroup.Get("/summary", h.GetAccountSummary)
	accountGroup.Get("/type/:type", h.GetAccountsByType)
	accountGroup.Get("/currency/:currency", h.GetAccountsByCurrency)
	accountGroup.Put("/:id/balance", h.UpdateAccountBalance)
	accountGroup.Get("/:id", h.GetAccountByID)
}
