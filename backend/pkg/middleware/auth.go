package middleware

import (
	"strings"

	apperrors "erp-system/pkg/errors"
	"erp-system/pkg/jwt"
	"erp-system/pkg/response"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	jwtManager *jwt.Manager
	logger     *zap.Logger
}

func NewAuthMiddleware(jwtManager *jwt.Manager, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Authenticate validates the Bearer token and injects claims into context
func (m *AuthMiddleware) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				err := apperrors.NewCustomError("missing authorization header").
					WithErrorCode(apperrors.ErrCodeUnauthorized).
					WithMessageID("error_unauthorized").
					WithHTTPCode(401)
				m.logger.Warn("auth failed: missing header",
					zap.String("path", c.Request().URL.Path),
					zap.String("method", c.Request().Method),
				)
				return response.Error(c, err)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				err := apperrors.NewCustomError("invalid authorization format").
					WithErrorCode(apperrors.ErrCodeUnauthorized).
					WithMessageID("error_unauthorized").
					WithHTTPCode(401)
				m.logger.Warn("auth failed: invalid format",
					zap.String("path", c.Request().URL.Path),
				)
				return response.Error(c, err)
			}

			claims, jwtErr := m.jwtManager.Validate(parts[1])
			if jwtErr != nil {
				m.logger.Warn("auth failed: invalid token",
					zap.String("path", c.Request().URL.Path),
					zap.Error(jwtErr),
				)
				return response.Error(c, jwtErr)
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("roles", claims.Roles)
			c.Set("permissions", claims.Perms)

			return next(c)
		}
	}
}

// RequirePermission checks if user has a specific permission
func (m *AuthMiddleware) RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			perms, ok := c.Get("permissions").([]string)
			if !ok {
				err := apperrors.NewCustomError("no permissions found").
					WithErrorCode(apperrors.ErrCodeForbidden).
					WithMessageID("error_forbidden").
					WithHTTPCode(403)
				m.logger.Warn("authz failed: no permissions",
					zap.String("path", c.Request().URL.Path),
				)
				return response.Error(c, err)
			}

			for _, p := range perms {
				if p == permission || p == "*" {
					return next(c)
				}
			}

			err := apperrors.NewCustomError("insufficient permission: " + permission).
				WithErrorCode(apperrors.ErrCodeInsufficientPermission).
				WithMessageID("error_forbidden").
				WithHTTPCode(403)
			m.logger.Warn("authz failed: insufficient permission",
				zap.String("path", c.Request().URL.Path),
				zap.String("required_permission", permission),
			)
			return response.Error(c, err)
		}
	}
}

// RequireRole checks if user has at least one of the specified roles
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRoles, ok := c.Get("roles").([]string)
			if !ok || len(userRoles) == 0 {
				err := apperrors.NewCustomError("no role found").
					WithErrorCode(apperrors.ErrCodeForbidden).
					WithMessageID("error_forbidden").
					WithHTTPCode(403)
				m.logger.Warn("authz failed: no role",
					zap.String("path", c.Request().URL.Path),
				)
				return response.Error(c, err)
			}

			for _, userRole := range userRoles {
				for _, allowed := range allowedRoles {
					if userRole == allowed {
						return next(c)
					}
				}
			}

			err := apperrors.NewCustomError("role not allowed").
				WithErrorCode(apperrors.ErrCodeForbidden).
				WithMessageID("error_forbidden").
				WithHTTPCode(403)
			m.logger.Warn("authz failed: role not allowed",
				zap.String("path", c.Request().URL.Path),
				zap.Strings("user_roles", userRoles),
				zap.Strings("allowed_roles", allowedRoles),
			)
			return response.Error(c, err)
		}
	}
}
