package handler

import (
	"net/http"

	authusecase "erp-system/internal/auth/usecase"
	"erp-system/pkg/response"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	usecase authusecase.AuthUsecase
}

func NewAuthHandler(uc authusecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

// Register godoc
// @Summary      Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  authusecase.RegisterRequest  true  "Register payload"
// @Success      201   {object}  authusecase.AuthResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req authusecase.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	res, err := h.usecase.Register(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, res)
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  authusecase.LoginRequest  true  "Login payload"
// @Success      200   {object}  authusecase.AuthResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req authusecase.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	res, err := h.usecase.Login(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, res)
}

// GetProfile godoc
// @Summary      Get current user profile
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  authusecase.ProfileResponse
// @Router       /auth/me [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	res, err := h.usecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, res)
}
