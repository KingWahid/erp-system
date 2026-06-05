package modules

import (
	authdomain "erp-system/internal/auth/domain"
	authhandler "erp-system/internal/auth/handler"
	authrepo "erp-system/internal/auth/repository"
	authusecase "erp-system/internal/auth/usecase"

	"go.uber.org/fx"
)

// AuthModule provides all auth-layer dependencies
var AuthModule = fx.Module("auth",
	fx.Provide(
		// Repository layer
		// NewUserRepository(db, logger) → authdomain.UserRepository
		fx.Annotate(
			authrepo.NewUserRepository,
			fx.As(new(authdomain.UserRepository)),
		),
		// NewRoleRepository(db, logger) → authrepo.RoleRepository
		fx.Annotate(
			authrepo.NewRoleRepository,
			fx.As(new(authrepo.RoleRepository)),
		),

		// Usecase layer
		// NewAuthUsecase(UserRepository, RoleRepository, *jwt.Manager, *zap.Logger) → AuthUsecase
		fx.Annotate(
			authusecase.NewAuthUsecase,
			fx.As(new(authusecase.AuthUsecase)),
		),

		// Handler/Adapter layer
		authhandler.NewAuthAdapter,
	),
)
