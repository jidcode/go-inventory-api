package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

func (handler *AuthHandler) Register(ctx echo.Context) error {
	var input struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Validation failed")
	}

	user, err := handler.service.RegisterUser(input.Username, input.Email, input.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Registration failed")
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (handler *AuthHandler) Login(ctx echo.Context) error {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Validation failed")
	}

	token, err := handler.service.LoginUser(input.Email, input.Password)
	if err != nil {
		fmt.Printf("User login failed for email %s: %s\n", input.Email, err)
		return ctx.JSON(http.StatusUnauthorized, "Invalid email or password. Try again.")
	}

	return ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (handler *AuthHandler) GetProfile(ctx echo.Context) error {
	user := ctx.Get("user").(*models.User)

	return ctx.JSON(http.StatusOK, user)
}

func (handler *AuthHandler) CheckTokenExpiration(ctx echo.Context) error {
	tokenString := ctx.Request().Header.Get("Authorization")
	if tokenString == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "No token provided"})
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	isExpired := handler.service.IsTokenExpired(tokenString)

	return ctx.JSON(http.StatusOK, map[string]bool{
		"expired": isExpired,
	})
}
