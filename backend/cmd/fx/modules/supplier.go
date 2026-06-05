package modules

import (
	supplierdomain "erp-system/internal/supplier/domain"
	supplierhandler "erp-system/internal/supplier/handler"
	supplierrepo "erp-system/internal/supplier/repository"
	supplierusecase "erp-system/internal/supplier/usecase"

	"go.uber.org/fx"
)

// SupplierModule provides all supplier-layer dependencies
var SupplierModule = fx.Module("supplier",
	fx.Provide(
		// Base repository: NewSupplierRepository(db, logger) → domain.SupplierRepository
		// Named "base" so fx can distinguish it from the cached version
		fx.Annotate(
			supplierrepo.NewSupplierRepository,
			fx.ResultTags(`name:"baseSupplierRepo"`),
		),

		// Cached repository decorator:
		// NewCachedSupplierRepository(baseRepo, cacheManager, logger) → domain.SupplierRepository
		fx.Annotate(
			supplierrepo.NewCachedSupplierRepository,
			fx.ParamTags(`name:"baseSupplierRepo"`, ``, ``),
			fx.As(new(supplierdomain.SupplierRepository)),
		),

		// Usecase layer
		fx.Annotate(
			supplierusecase.NewSupplierUsecase,
			fx.As(new(supplierusecase.SupplierUsecase)),
		),

		// Handler/Adapter layer
		supplierhandler.NewSupplierAdapter,
	),
)
