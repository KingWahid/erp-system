package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ZapLogger returns an Echo middleware that logs requests using zap
func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			req := c.Request()
			res := c.Response()

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("ip", c.RealIP()),
				zap.String("request_id", res.Header().Get(echo.HeaderXRequestID)),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				log.Error("request error", fields...)
			} else if res.Status >= 500 {
				log.Error("server error", fields...)
			} else if res.Status >= 400 {
				log.Warn("client error", fields...)
			} else {
				log.Info("request", fields...)
			}

			return err
		}
	}
}
