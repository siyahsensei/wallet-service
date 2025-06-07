package routes

import (
	"strconv"

	"siyahsensei/wallet-service/domain/account"
	presentation "siyahsensei/wallet-service/presentation/account"

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

func (h *AccountHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	accountGroup := router.Group("/accounts", authMiddleware)

	accountGroup.Post("/", h.CreateAccount)
	accountGroup.Put("/:id", h.UpdateAccount)
	accountGroup.Delete("/:id", h.DeleteAccount)
	accountGroup.Get("/", h.GetUserAccounts)
	accountGroup.Get("/filter", h.FilterAccounts)
	accountGroup.Get("/summary", h.GetAccountSummary)
	accountGroup.Get("/:id", h.GetAccountByID)
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account body presentation.CreateAccountRequest true "Account creation data"
// @Success 201 {object} map[string]presentation.AccountResponse
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

	var req presentation.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

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
		"account": presentation.ToAccountResponse(createdAccount),
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
// @Param account body presentation.UpdateAccountRequest true "Account update data"
// @Success 200 {object} map[string]presentation.AccountResponse
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

	var req presentation.UpdateAccountRequest
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
		"account": presentation.ToAccountResponse(updatedAccount),
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
// @Param with-assets query bool false "Include assets in response"
// @Success 200 {object} map[string]interface{}
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

	// Check if assets should be included
	withAssets := c.Query("with-assets") == "true"

	if withAssets {
		foundAccount, err := h.accountService.HandleGetAccountByIDWithAssetsQuery(c.Context(), query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"account": presentation.ToAccountWithAssetsResponse(foundAccount),
		})
	} else {
		foundAccount, err := h.accountService.HandleGetAccountByIDQuery(c.Context(), query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"account": presentation.ToAccountResponse(foundAccount),
		})
	}
}

// GetUserAccounts godoc
// @Summary Get all user accounts
// @Description Get all accounts for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param with-assets query bool false "Include assets in response"
// @Success 200 {object} interface{}
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

	// Check if assets should be included
	withAssets := c.Query("with-assets") == "true"

	if withAssets {
		accounts, err := h.accountService.HandleGetUserAccountsWithAssetsQuery(c.Context(), query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var accountResponses []presentation.AccountWithAssetsResponse
		for _, a := range accounts {
			accountResponses = append(accountResponses, presentation.ToAccountWithAssetsResponse(a))
		}

		return c.Status(fiber.StatusOK).JSON(presentation.AccountsWithAssetsListResponse{
			Accounts: accountResponses,
			Total:    len(accountResponses),
		})
	} else {
		accounts, err := h.accountService.HandleGetUserAccountsQuery(c.Context(), query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var accountResponses []presentation.AccountResponse
		for _, a := range accounts {
			accountResponses = append(accountResponses, presentation.ToAccountResponse(a))
		}

		return c.Status(fiber.StatusOK).JSON(presentation.AccountsListResponse{
			Accounts: accountResponses,
			Total:    len(accountResponses),
		})
	}
}

// FilterAccounts godoc
// @Summary Filter accounts
// @Description Filter accounts with optional parameters for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param accountType query string false "Account Type"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} presentation.AccountsListResponse
// @Failure 401 {object} map[string]string
// @Router /accounts/filter [get]
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

	var accountResponses []presentation.AccountResponse
	for _, a := range accounts {
		accountResponses = append(accountResponses, presentation.ToAccountResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(presentation.AccountsListResponse{
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
// @Success 200 {object} map[string]presentation.AccountSummaryResponse
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
		"summary": presentation.ToAccountSummaryResponse(summary),
	})
}
