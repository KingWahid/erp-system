package server

import (
	authhandler "erp-system/internal/auth/handler"
	"erp-system/internal/generated"
	supplierhandler "erp-system/internal/supplier/handler"
	"erp-system/pkg/middleware"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// combinedServer composes both adapters to satisfy generated.ServerInterface.
type combinedServer struct {
	*authhandler.AuthAdapter
	*supplierhandler.SupplierAdapter
}

// Router wires everything together.
type Router struct {
	echo            *echo.Echo
	authAdapter     *authhandler.AuthAdapter
	supplierAdapter *supplierhandler.SupplierAdapter
	authMiddleware  *middleware.AuthMiddleware
	logger          *zap.Logger
}

func NewRouter(
	e *echo.Echo,
	authAdapter *authhandler.AuthAdapter,
	supplierAdapter *supplierhandler.SupplierAdapter,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) *Router {
	return &Router{
		echo:            e,
		authAdapter:     authAdapter,
		supplierAdapter: supplierAdapter,
		authMiddleware:  authMiddleware,
		logger:          logger,
	}
}

func (r *Router) Register() {
	e := r.echo

	// Global middleware
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(middleware.ZapLogger(r.logger))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"service": "erp-supplier-management",
		})
	})

	// Compose single ServerInterface
	srv := &combinedServer{
		AuthAdapter:     r.authAdapter,
		SupplierAdapter: r.supplierAdapter,
	}

	// Build per-operationId middleware map
	//
	// Only authentication middleware is applied here.
	// Permission checking is delegated to the database via RBAC,
	// so no hardcoded permission checks are needed in this router.
	//
	// Public routes: LoginUser, RegisterUser (no auth)
	// Protected routes: everything else (requires valid JWT)
	//
	auth := r.authMiddleware.Authenticate()

	opMW := map[string][]echo.MiddlewareFunc{
		// Auth (public)
		// LoginUser and RegisterUser have no middleware

		// Auth (protected)
		"GetProfile": {auth},

		// Supplier Management (all protected)
		"GetSupplierStats":        {auth},
		"ListSuppliers":           {auth},
		"CreateSupplier":          {auth},
		"GetSupplier":             {auth},
		"UpdateSupplier":          {auth},
		"DeleteSupplier":          {auth},
		"BlockSupplier":           {auth},
		"AdvanceSupplierStage":    {auth},
		"GetSupplierMaterials":    {auth},
		"UpdateSupplierMaterials": {auth},
		"GetSupplierRatings":      {auth},
		"AddSupplierRating":       {auth},

		// Addresses
		"GetSupplierAddresses":   {auth},
		"AddSupplierAddress":     {auth},
		"UpdateSupplierAddress":  {auth},
		"DeleteSupplierAddress":  {auth},
		"SetMainSupplierAddress": {auth},

		// Contacts
		"GetSupplierContacts":       {auth},
		"AddSupplierContact":        {auth},
		"UpdateSupplierContact":     {auth},
		"DeleteSupplierContact":     {auth},
		"SetPrimarySupplierContact": {auth},

		// Groups
		"GetSupplierGroups":   {auth},
		"AddSupplierGroup":    {auth},
		"UpdateSupplierGroup": {auth},
		"DeleteSupplierGroup": {auth},

		// Stage History
		"GetStageHistory": {auth},

		// Outstandings
		"GetSupplierOutstandings": {auth},
	}

	// Single registration point — all routes from OpenAPI
	generated.RegisterHandlersWithOptions(e, srv, generated.RegisterHandlersOptions{
		OperationMiddlewares: opMW,
	})
}
