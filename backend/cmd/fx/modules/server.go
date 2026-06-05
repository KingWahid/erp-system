package modules

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"erp-system/cmd/server"
	authhandler "erp-system/internal/auth/handler"
	supplierhandler "erp-system/internal/supplier/handler"
	"erp-system/pkg/config"
	"erp-system/pkg/jwt"
	"erp-system/pkg/middleware"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerModule provides Echo, Router, AuthMiddleware, and lifecycle hooks
var ServerModule = fx.Module("server",
	fx.Provide(
		newEcho,
		newAuthMiddleware,
		newRouter,
	),
	fx.Invoke(startServer),
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}

func newAuthMiddleware(jwtManager *jwt.Manager, logger *zap.Logger) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(jwtManager, logger)
}

func newRouter(
	e *echo.Echo,
	authAdapter *authhandler.AuthAdapter,
	supplierAdapter *supplierhandler.SupplierAdapter,
	authMW *middleware.AuthMiddleware,
	logger *zap.Logger,
) *server.Router {
	return server.NewRouter(e, authAdapter, supplierAdapter, authMW, logger)
}

// startServer registers routes and manages server lifecycle via fx
func startServer(
	lc fx.Lifecycle,
	e *echo.Echo,
	router *server.Router,
	cfg *config.Config,
	logger *zap.Logger,
) {
	router.Register()

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("server starting",
					zap.String("addr", addr),
					zap.String("env", cfg.App.Env),
				)
				if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
					logger.Fatal("server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("server shutting down...")
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := e.Shutdown(shutdownCtx); err != nil {
				logger.Error("forced shutdown", zap.Error(err))
				return err
			}
			logger.Info("server exited")
			return nil
		},
	})
}
