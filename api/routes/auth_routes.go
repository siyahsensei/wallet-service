package routes

import (
	"siyahsensei/wallet-service/domain/user"
	"siyahsensei/wallet-service/infrastructure/configuration/auth"
	presentation "siyahsensei/wallet-service/presentation/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthRoute struct {
	userService *user.Handler
	jwtAuth     *auth.JWTMiddleware
}

func NewAuthRoute(userService *user.Handler, jwtAuth *auth.JWTMiddleware) *AuthRoute {
	return &AuthRoute{
		userService: userService,
		jwtAuth:     jwtAuth,
	}
}

func (h *AuthRoute) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Get("/me", authMiddleware, h.Me)
	router.Put("/me", authMiddleware, h.UpdateUser)
	router.Put("/change-password", authMiddleware, h.ChangePassword)
	router.Delete("/me", authMiddleware, h.DeleteUser)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.RegisterUserCommand true "User registration data"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthRoute) Register(c *fiber.Ctx) error {
	var command user.RegisterUserCommand
	if err := c.BodyParser(&command); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	newUser, err := h.userService.HandleRegisterUserCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := h.jwtAuth.GenerateToken(newUser.ID, newUser.Email, h.userService.GetTokenExpiry())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(presentation.TokenResponse{
		Token: token,
		User:  presentation.ToPublicUser(newUser),
	})
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body user.LoginUserCommand true "User login credentials"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthRoute) Login(c *fiber.Ctx) error {
	var command user.LoginUserCommand
	if err := c.BodyParser(&command); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userInfo, err := h.userService.HandleLoginUserCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := h.jwtAuth.GenerateToken(userInfo.ID, userInfo.Email, h.userService.GetTokenExpiry())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(presentation.TokenResponse{
		Token: token,
		User:  presentation.ToPublicUser(userInfo),
	})
}

// Me godoc
// @Summary Get current user information
// @Description Get current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]UserPublic
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthRoute) Me(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := user.GetUserByIDQuery{
		ID: userIDValue.String(),
	}

	userInfo, err := h.userService.HandleGetUserByIDQuery(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user information",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": presentation.ToPublicUser(userInfo),
	})
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]UserPublic
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/me [put]
func (h *AuthRoute) UpdateUser(c *fiber.Ctx) error {
	var req presentation.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	command := user.UpdateUserCommand{
		ID:        userIDValue.String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	updatedUser, err := h.userService.HandleUpdateUserCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user information",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": presentation.ToPublicUser(updatedUser),
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current authenticated user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param passwords body ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/change-password [put]
func (h *AuthRoute) ChangePassword(c *fiber.Ctx) error {
	var req presentation.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	command := user.ChangePasswordCommand{
		UserID:      userIDValue.String(),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	err := h.userService.HandleChangePasswordCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete current authenticated user account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body DeleteUserRequest true "Password confirmation"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/me [delete]
func (h *AuthRoute) DeleteUser(c *fiber.Ctx) error {
	var req presentation.DeleteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	command := user.DeleteUserCommand{
		UserID:   userIDValue.String(),
		Password: req.Password,
	}

	err := h.userService.HandleDeleteUserCommand(c.Context(), command)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}