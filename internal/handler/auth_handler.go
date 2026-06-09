package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rwndy/bookmark-api/internal/domain"
	"github.com/rwndy/bookmark-api/internal/handler/dto"
	"github.com/rwndy/bookmark-api/internal/service"
	"github.com/rwndy/bookmark-api/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, 400, "invalid request body")
	}

	if err := ValidateStruct(&req); err != nil {
		return response.Fail(c, 400, err.Error())
	}

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		return failFromError(c, err)
	}

	return response.Created(c, "user registered", user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, 400, "invalid request body")
	}

	if err := ValidateStruct(&req); err != nil {
		return response.Fail(c, 400, err.Error())
	}

	tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return failFromError(c, err)
	}

	return response.OK(c, "login successful", tokens)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, 400, "invalid request body")
	}

	if err := ValidateStruct(&req); err != nil {
		return response.Fail(c, 400, err.Error())
	}

	tokens, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		return failFromError(c, err)
	}

	return response.OK(c, "token refreshed", tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, 400, "invalid request body")
	}

	if err := ValidateStruct(&req); err != nil {
		return response.Fail(c, 400, err.Error())
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		return failFromError(c, err)
	}

	return response.OK(c, "logout successful", nil)
}

func failFromError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*domain.AppError); ok {
		return response.Fail(c, appErr.Code, appErr.Message)
	}
	return response.Fail(c, 500, "internal server error")
}
