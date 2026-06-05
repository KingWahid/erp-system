package modules

import (
	"erp-system/pkg/cache"
	"erp-system/pkg/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// CacheModule provides Redis cache and CacheManager
var CacheModule = fx.Module("cache",
	fx.Provide(
		newRedisCache,
		newCacheManager,
	),
)

// newRedisCache attempts Redis connection.
// Returns nil (not an error) if Redis is unavailable — cache degrades gracefully.
func newRedisCache(cfg *config.Config, logger *zap.Logger) (*cache.RedisCache, error) {
	rc, err := cache.NewRedis(cfg)
	if err != nil {
		logger.Warn("redis unavailable, caching disabled", zap.Error(err))
		return nil, nil // graceful degradation — NOT a fatal error
	}
	logger.Info("redis connected",
		zap.String("host", cfg.Redis.Host),
		zap.String("port", cfg.Redis.Port),
	)
	return rc, nil
}

// newCacheManager wraps RedisCache (may be nil) into CacheManager
func newCacheManager(rc *cache.RedisCache, logger *zap.Logger) *cache.CacheManager {
	if rc == nil {
		// Pass nil Cache interface — CacheManager handles this gracefully
		return cache.NewCacheManager(nil, logger)
	}
	return cache.NewCacheManager(rc, logger)
}
