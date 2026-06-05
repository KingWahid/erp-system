package handler

import (
	authusecase "erp-system/internal/auth/usecase"
	"erp-system/internal/generated"
	apperrors "erp-system/pkg/errors"
	"erp-system/pkg/response"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AuthAdapter implements the auth portion of generated.ServerInterface.
type AuthAdapter struct {
	usecase authusecase.AuthUsecase
	logger  *zap.Logger
}

func NewAuthAdapter(uc authusecase.AuthUsecase, logger *zap.Logger) *AuthAdapter {
	return &AuthAdapter{usecase: uc, logger: logger}
}

// RegisterUser implements generated.ServerInterface
func (a *AuthAdapter) RegisterUser(ctx echo.Context) error {
	var body generated.RegisterUserJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	role := ""
	if body.Role != nil {
		role = string(*body.Role)
	}

	res, err := a.usecase.Register(ctx.Request().Context(), authusecase.RegisterRequest{
		Name:     body.Name,
		Email:    string(body.Email),
		Password: body.Password,
		Role:     role,
	})
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Created(ctx, res)
}

// LoginUser implements generated.ServerInterface
func (a *AuthAdapter) LoginUser(ctx echo.Context) error {
	var body generated.LoginUserJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return response.Error(ctx, apperrors.BadRequest("invalid request body").
			WithErrorCode(apperrors.ErrCodeInvalidInput))
	}

	res, err := a.usecase.Login(ctx.Request().Context(), authusecase.LoginRequest{
		Email:    string(body.Email),
		Password: body.Password,
	})
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}

// GetProfile implements generated.ServerInterface
func (a *AuthAdapter) GetProfile(ctx echo.Context, params generated.GetProfileParams) error {
	userID, _ := ctx.Get("user_id").(string)
	res, err := a.usecase.GetProfile(ctx.Request().Context(), userID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, res)
}
