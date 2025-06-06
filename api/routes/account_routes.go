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

// Request models (without Balance and CurrencyCode)
type CreateAccountRequest struct {
	Name        string              `json:"name" validate:"required"`
	AccountType account.AccountType `json:"accountType" validate:"required"`
}

type UpdateAccountRequest struct {
	Name        string              `json:"name" validate:"required"`
	AccountType account.AccountType `json:"accountType" validate:"required"`
}

// Response models
type AccountResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	AccountType string    `json:"accountType"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AccountWithAssetsResponse struct {
	ID            string              `json:"id"`
	UserID        string              `json:"userId"`
	Name          string              `json:"name"`
	AccountType   string              `json:"accountType"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	Assets        []AssetInfoResponse `json:"assets"`
	TotalBalances map[string]float64  `json:"totalBalances"`
	AssetCounts   map[string]int      `json:"assetCounts"`
	LastUpdated   *time.Time          `json:"lastUpdated"`
}

type AssetInfoResponse struct {
	ID           string    `json:"id"`
	DefinitionID string    `json:"definitionId"`
	Type         string    `json:"type"`
	Quantity     float64   `json:"quantity"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	Currency     string    `json:"currency"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AccountsListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int               `json:"total"`
}

type AccountsWithAssetsListResponse struct {
	Accounts []AccountWithAssetsResponse `json:"accounts"`
	Total    int                         `json:"total"`
}

type AccountSummaryResponse struct {
	TotalAccounts int                         `json:"totalAccounts"`
	ByType        map[account.AccountType]int `json:"byType"`
	ByCurrency    map[string]float64          `json:"byCurrency"`
}

func toAccountResponse(a *account.Account) AccountResponse {
	return AccountResponse{
		ID:          a.ID.String(),
		UserID:      a.UserID.String(),
		Name:        a.Name,
		AccountType: string(a.AccountType),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func toAccountWithAssetsResponse(a *account.AccountWithAssets) AccountWithAssetsResponse {
	var assets []AssetInfoResponse
	for _, asset := range a.Assets {
		assets = append(assets, AssetInfoResponse{
			ID:           asset.ID.String(),
			DefinitionID: asset.DefinitionID.String(),
			Type:         asset.Type,
			Quantity:     asset.Quantity,
			Symbol:       asset.Symbol,
			Name:         asset.Name,
			Currency:     asset.Currency,
			UpdatedAt:    asset.UpdatedAt,
		})
	}

	return AccountWithAssetsResponse{
		ID:            a.Account.ID.String(),
		UserID:        a.Account.UserID.String(),
		Name:          a.Account.Name,
		AccountType:   string(a.Account.AccountType),
		CreatedAt:     a.Account.CreatedAt,
		UpdatedAt:     a.Account.UpdatedAt,
		Assets:        assets,
		TotalBalances: a.TotalBalances,
		AssetCounts:   a.AssetCounts,
		LastUpdated:   a.LastUpdated,
	}
}

func toAccountSummaryResponse(s *account.AccountSummary) AccountSummaryResponse {
	return AccountSummaryResponse{
		TotalAccounts: s.TotalAccounts,
		ByType:        s.ByType,
		ByCurrency:    s.ByCurrency,
	}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account body CreateAccountRequest true "Account creation data"
// @Success 201 {object} map[string]AccountResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /accounts [post]
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
		UserID:      userIDValue.String(),
		Name:        req.Name,
		AccountType: req.AccountType,
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

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an existing account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param account body UpdateAccountRequest true "Account update data"
// @Success 200 {object} map[string]AccountResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /accounts/{id} [put]
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
		ID:          accountID,
		UserID:      userIDValue.String(),
		Name:        req.Name,
		AccountType: req.AccountType,
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

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an existing account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /accounts/{id} [delete]
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

// GetAccountByID godoc
// @Summary Get account by ID
// @Description Get a specific account by ID for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]AccountResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /accounts/{id} [get]
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

	foundAccount, err := h.accountService.HandleGetAccountByIDQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account": toAccountResponse(foundAccount),
	})
}

// GetAccountByIDWithAssets godoc
// @Summary Get account by ID with assets
// @Description Get a specific account by ID with its assets for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]AccountWithAssetsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /accounts/{id}/with-assets [get]
func (h *AccountHandler) GetAccountByIDWithAssets(c *fiber.Ctx) error {
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

	foundAccount, err := h.accountService.HandleGetAccountByIDWithAssetsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"account": toAccountWithAssetsResponse(foundAccount),
	})
}

// GetUserAccounts godoc
// @Summary Get all user accounts
// @Description Get all accounts for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AccountsListResponse
// @Failure 401 {object} map[string]string
// @Router /accounts [get]
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

// GetUserAccountsWithAssets godoc
// @Summary Get all user accounts with assets
// @Description Get all accounts with their assets for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AccountsWithAssetsListResponse
// @Failure 401 {object} map[string]string
// @Router /accounts/with-assets [get]
func (h *AccountHandler) GetUserAccountsWithAssets(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := account.GetUserAccountsQuery{
		UserID: userIDValue.String(),
	}

	accounts, err := h.accountService.HandleGetUserAccountsWithAssetsQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var accountResponses []AccountWithAssetsResponse
	for _, a := range accounts {
		accountResponses = append(accountResponses, toAccountWithAssetsResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(AccountsWithAssetsListResponse{
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

// GetAccountSummary godoc
// @Summary Get account summary
// @Description Get summary statistics for all user accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]AccountSummaryResponse
// @Failure 401 {object} map[string]string
// @Router /accounts/summary [get]
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
	accountGroup.Get("/with-assets", h.GetUserAccountsWithAssets)
	accountGroup.Get("/filter", h.FilterAccounts)
	accountGroup.Get("/summary", h.GetAccountSummary)
	accountGroup.Get("/type/:type", h.GetAccountsByType)
	accountGroup.Get("/:id", h.GetAccountByID)
	accountGroup.Get("/:id/with-assets", h.GetAccountByIDWithAssets)
}
