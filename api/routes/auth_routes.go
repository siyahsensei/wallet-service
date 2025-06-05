package routes

import (
	"siyahsensei/wallet-service/domain/user"
	"siyahsensei/wallet-service/infrastructure/configuration/auth"

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

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token string      `json:"token"`
	User  *UserPublic `json:"user"`
}

type UserPublic struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type DeleteUserRequest struct {
	Password string `json:"password" validate:"required"`
}

func toPublicUser(u *user.User) *UserPublic {
	return &UserPublic{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (h *AuthRoute) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	newUser, err := h.userService.RegisterUser(
		c.Context(),
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)
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

	return c.Status(fiber.StatusCreated).JSON(TokenResponse{
		Token: token,
		User:  toPublicUser(newUser),
	})
}

func (h *AuthRoute) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, err := h.userService.LoginUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	userInfo, err := h.userService.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user information",
		})
	}

	return c.Status(fiber.StatusOK).JSON(TokenResponse{
		Token: token,
		User:  toPublicUser(userInfo),
	})
}

func (h *AuthRoute) Me(c *fiber.Ctx) error {
	userIDValue, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userInfo, err := h.userService.GetUserByID(c.Context(), userIDValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user information",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": toPublicUser(userInfo),
	})
}

func (h *AuthRoute) UpdateUser(c *fiber.Ctx) error {
	var req UpdateUserRequest
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

	userInfo, err := h.userService.GetUserByID(c.Context(), userIDValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user information",
		})
	}

	userInfo.FirstName = req.FirstName
	userInfo.LastName = req.LastName
	userInfo.Email = req.Email

	err = h.userService.UpdateUser(c.Context(), userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user information",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": toPublicUser(userInfo),
	})
}

func (h *AuthRoute) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
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

	err := h.userService.ChangePassword(c.Context(), userIDValue, req.OldPassword, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

func (h *AuthRoute) DeleteUser(c *fiber.Ctx) error {
	var req DeleteUserRequest
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

	err := h.userService.ValidateUserPassword(c.Context(), userIDValue, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	err = h.userService.DeleteUser(c.Context(), userIDValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *AuthRoute) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Get("/me", authMiddleware, h.Me)
	router.Put("/me", authMiddleware, h.UpdateUser)
	router.Put("/change-password", authMiddleware, h.ChangePassword)
	router.Delete("/me", authMiddleware, h.DeleteUser)
}
