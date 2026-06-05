package jwt

import (
	"net/http"
	"time"

	"erp-system/pkg/config"
	apperrors "erp-system/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims is the JWT payload
type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	Perms  []string `json:"permissions"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret          []byte
	expirationHours int
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		secret:          []byte(cfg.JWT.Secret),
		expirationHours: cfg.JWT.ExpirationHours,
	}
}

// Generate creates a signed JWT token
func (m *Manager) Generate(userID, email string, roles []string, permissions []string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Perms:  permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.expirationHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Validate parses and validates the token, returning claims
func (m *Manager) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.NewCustomError("unexpected signing method").
				WithErrorCode(apperrors.ErrCodeInvalidToken).
				WithMessageID("error_invalid_token").
				WithHTTPCode(http.StatusUnauthorized)
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, apperrors.NewCustomError("invalid or expired token").
			WithErrorCode(apperrors.ErrCodeInvalidToken).
			WithMessageID("error_invalid_token").
			WithHTTPCode(http.StatusUnauthorized).
			WithCause(err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperrors.NewCustomError("invalid token").
			WithErrorCode(apperrors.ErrCodeInvalidToken).
			WithMessageID("error_invalid_token").
			WithHTTPCode(http.StatusUnauthorized)
	}

	return claims, nil
}
