package usecase

import (
	"context"
	"net/http"
	"time"

	authdomain "erp-system/internal/auth/domain"
	authrepo "erp-system/internal/auth/repository"
	apperrors "erp-system/pkg/errors"
	"erp-system/pkg/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo   authdomain.UserRepository
	roleRepo   authrepo.RoleRepository
	jwtManager *jwt.Manager
	logger     *zap.Logger
}

func NewAuthUsecase(
	userRepo authdomain.UserRepository,
	roleRepo authrepo.RoleRepository,
	jwtManager *jwt.Manager,
	logger *zap.Logger,
) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.NewCustomError("email already registered").
			WithErrorCode(apperrors.ErrCodeEmailExists).
			WithMessageID("error_email_exists").
			WithHTTPCode(http.StatusConflict)
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		customErr := apperrors.NewCustomError("failed to hash password").
			WithErrorCode(apperrors.ErrCodePasswordHashing).
			WithMessageID("error_password_hashing").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(uc.logger, zap.String("operation", "Register"))
		return nil, customErr
	}

	user := &authdomain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPw),
		IsActive: true,
	}

	roleName := req.Role
	if roleName == "" {
		roleName = "viewer"
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	role, err := uc.roleRepo.FindRoleByName(ctx, roleName)
	if err != nil {
		return nil, apperrors.NewCustomError("invalid role: " + roleName).
			WithErrorCode(apperrors.ErrCodeRoleNotFound).
			WithMessageID("error_role_not_found").
			WithHTTPCode(http.StatusBadRequest).
			WithCause(err)
	}

	if err := uc.roleRepo.AssignRoleToUser(ctx, user.ID, role.ID); err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	user, err = uc.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	permissions := user.GetPermissions()
	roleNames := user.GetRoleNames()

	token, err := uc.jwtManager.Generate(user.ID.String(), user.Email, roleNames, permissions)
	if err != nil {
		customErr := apperrors.NewCustomError("failed to generate token").
			WithErrorCode(apperrors.ErrCodeTokenGeneration).
			WithMessageID("error_token_generation").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(uc.logger, zap.String("operation", "Register"))
		return nil, customErr
	}

	uc.logger.Info("user registered successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
		zap.String("role", roleName),
	)

	return &AuthResponse{
		Token: token,
		User: UserPayload{
			ID:          user.ID.String(),
			Name:        user.Name,
			Email:       user.Email,
			Role:        roleNames,
			Permissions: permissions,
		},
	}, nil
}

func (uc *authUsecase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewCustomError("invalid email or password").
			WithErrorCode(apperrors.ErrCodeInvalidCredentials).
			WithMessageID("error_invalid_credentials").
			WithHTTPCode(http.StatusUnauthorized)
	}

	if !user.IsActive {
		return nil, apperrors.NewCustomError("account is inactive").
			WithErrorCode(apperrors.ErrCodeInactiveAccount).
			WithMessageID("error_inactive_account").
			WithHTTPCode(http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.NewCustomError("invalid email or password").
			WithErrorCode(apperrors.ErrCodeInvalidCredentials).
			WithMessageID("error_invalid_credentials").
			WithHTTPCode(http.StatusUnauthorized)
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = uc.userRepo.Update(ctx, user)

	permissions := user.GetPermissions()
	roleNames := user.GetRoleNames()

	token, err := uc.jwtManager.Generate(user.ID.String(), user.Email, roleNames, permissions)
	if err != nil {
		customErr := apperrors.NewCustomError("failed to generate token").
			WithErrorCode(apperrors.ErrCodeTokenGeneration).
			WithMessageID("error_token_generation").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(uc.logger, zap.String("operation", "Login"))
		return nil, customErr
	}

	uc.logger.Info("user logged in successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return &AuthResponse{
		Token: token,
		User: UserPayload{
			ID:          user.ID.String(),
			Name:        user.Name,
			Email:       user.Email,
			Role:        roleNames,
			Permissions: permissions,
		},
	}, nil
}

func (uc *authUsecase) GetProfile(ctx context.Context, userID string) (*ProfileResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.NewCustomError("invalid user id").
			WithErrorCode(apperrors.ErrCodeInvalidUUID).
			WithMessageID("error_invalid_uuid").
			WithHTTPCode(http.StatusBadRequest).
			WithCause(err)
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewCustomErrorFromError(err)
	}

	uc.logger.Info("profile fetched",
		zap.String("user_id", user.ID.String()),
	)

	return &ProfileResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.GetRoleNames(),
	}, nil
}
