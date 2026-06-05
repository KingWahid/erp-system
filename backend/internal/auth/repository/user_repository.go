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

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) authdomain.UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) Create(ctx context.Context, user *authdomain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if apperrors.IsDuplicateKeyError(err) {
			customErr := apperrors.NewCustomError("email already exists").
				WithErrorCode(apperrors.ErrCodeEmailExists).
				WithMessageID("error_email_exists").
				WithHTTPCode(http.StatusConflict).
				WithCause(err)
			customErr.LogError(r.logger, zap.String("operation", "CreateUser"))
			return customErr
		}
		customErr := apperrors.NewCustomError("failed to create user").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger, zap.String("operation", "CreateUser"))
		return customErr
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	var user authdomain.User
	err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customErr := apperrors.NewCustomError("user not found").
				WithErrorCode(apperrors.ErrCodeUserNotFound).
				WithMessageID("error_user_not_found").
				WithHTTPCode(http.StatusNotFound).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "FindUserByID"),
				zap.String("user_id", id.String()),
			)
			return nil, customErr
		}
		customErr := apperrors.NewCustomError("failed to find user").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "FindUserByID"),
			zap.String("user_id", id.String()),
		)
		return nil, customErr
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	var user authdomain.User
	err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customErr := apperrors.NewCustomError("user not found").
				WithErrorCode(apperrors.ErrCodeUserNotFound).
				WithMessageID("error_user_not_found").
				WithHTTPCode(http.StatusNotFound).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "FindUserByEmail"),
				zap.String("email", email),
			)
			return nil, customErr
		}
		customErr := apperrors.NewCustomError("failed to find user by email").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "FindUserByEmail"),
			zap.String("email", email),
		)
		return nil, customErr
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *authdomain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to update user").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpdateUser"),
			zap.String("user_id", user.ID.String()),
		)
		return customErr
	}
	return nil
}
