package repository

import (
	"context"
	"errors"
	"net/http"

	authdomain "erp-system/internal/auth/domain"
	apperrors "erp-system/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RoleRepository handles role and permission queries
type RoleRepository interface {
	FindRoleByName(ctx context.Context, name string) (*authdomain.Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error
}

type roleRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRoleRepository(db *gorm.DB, logger *zap.Logger) RoleRepository {
	return &roleRepository{db: db, logger: logger}
}

func (r *roleRepository) FindRoleByName(ctx context.Context, name string) (*authdomain.Role, error) {
	var role authdomain.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name = ? AND is_active = TRUE", name).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customErr := apperrors.NewCustomError("role not found: " + name).
				WithErrorCode(apperrors.ErrCodeRoleNotFound).
				WithMessageID("error_role_not_found").
				WithHTTPCode(http.StatusNotFound).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "FindRoleByName"),
				zap.String("role_name", name),
			)
			return nil, customErr
		}
		customErr := apperrors.NewCustomError("failed to find role").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "FindRoleByName"),
			zap.String("role_name", name),
		)
		return nil, customErr
	}
	return &role, nil
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	userRole := authdomain.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	if err := r.db.WithContext(ctx).Create(&userRole).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to assign role to user").
			WithErrorCode(apperrors.ErrCodeRoleAssignFailed).
			WithMessageID("error_role_assign_failed").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "AssignRoleToUser"),
			zap.String("user_id", userID.String()),
			zap.String("role_id", roleID.String()),
		)
		return customErr
	}
	return nil
}
