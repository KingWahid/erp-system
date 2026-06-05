package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	// ErrCacheMiss indicates the key was not found in cache
	ErrCacheMiss = errors.New("cache miss")
	// ErrCacheDisabled indicates caching is disabled
	ErrCacheDisabled = errors.New("cache disabled")
)

// ═══════════════════════════════════════════════════════════════
// CACHE-ASIDE PATTERN IMPLEMENTATION
// ═══════════════════════════════════════════════════════════════

// CacheManager wraps Cache with JSON marshaling helpers
type CacheManager struct {
	cache   Cache
	logger  *zap.Logger
	enabled bool
}

// NewCacheManager creates a new cache manager
func NewCacheManager(cache Cache, logger *zap.Logger) *CacheManager {
	return &CacheManager{
		cache:   cache,
		logger:  logger,
		enabled: cache != nil,
	}
}

// GetJSON retrieves and unmarshals JSON data from cache.
// Returns ErrCacheMiss if key not found, ErrCacheDisabled if caching disabled.
func (cm *CacheManager) GetJSON(ctx context.Context, key string, result interface{}) error {
	if !cm.enabled {
		return ErrCacheDisabled
	}

	data, err := cm.cache.Get(ctx, key)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrCacheMiss
	}

	return json.Unmarshal(data, result)
}

// SetJSON marshals and stores JSON data in cache.
func (cm *CacheManager) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !cm.enabled {
		return ErrCacheDisabled
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cm.cache.Set(ctx, key, data, ttl)
}

// Delete removes a key from cache.
func (cm *CacheManager) Delete(ctx context.Context, key string) error {
	if !cm.enabled {
		return ErrCacheDisabled
	}
	return cm.cache.Delete(ctx, key)
}

// DeletePattern removes all keys matching pattern.
func (cm *CacheManager) DeletePattern(ctx context.Context, pattern string) error {
	if !cm.enabled {
		return ErrCacheDisabled
	}
	return cm.cache.DeletePattern(ctx, pattern)
}

// ═══════════════════════════════════════════════════════════════
// CACHE-ASIDE PATTERN: GetFromCacheOrDB
// ═══════════════════════════════════════════════════════════════

// GetFromCacheOrDB implements the cache-aside pattern:
// 1. Try to get from cache first
// 2. On cache miss, fetch from DB
// 3. Update cache asynchronously (non-blocking)
//
// This ensures the request is never blocked by cache writes.
func (cm *CacheManager) GetFromCacheOrDB(
	ctx context.Context,
	cacheKey string,
	ttl time.Duration,
	result interface{},
	dbFetch func() error,
) error {
	// If caching is disabled, go directly to DB
	if !cm.enabled {
		return dbFetch()
	}

	// Try to get from cache first
	err := cm.GetJSON(ctx, cacheKey, result)
	if err == nil {
		// Cache hit
		cm.logger.Debug("cache hit", zap.String("key", cacheKey))
		return nil
	}

	// Check if it's a cache miss vs cache error
	if !errors.Is(err, ErrCacheMiss) && !errors.Is(err, ErrCacheDisabled) {
		// Cache error - log warning but continue to DB
		cm.logger.Warn("cache error", zap.String("key", cacheKey), zap.Error(err))
	}

	// Cache miss or error - fetch from DB
	cm.logger.Debug("cache miss, fetching from DB", zap.String("key", cacheKey))
	err = dbFetch()
	if err != nil {
		return err
	}

	// Marshal data synchronously to avoid race conditions with callers
	// who may mutate the result after this function returns.
	// Only the Redis write operation runs in the background.
	data, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		cm.logger.Warn("failed to marshal data for cache",
			zap.String("key", cacheKey),
			zap.Error(marshalErr))
		// Don't fail the request, just skip caching
		return nil
	}

	// Update cache with pre-marshaled data (non-blocking)
	go func(jsonData []byte) {
		// Use background context to avoid cancellation
		bgCtx := context.Background()
		err := cm.cache.Set(bgCtx, cacheKey, jsonData, ttl)
		if err != nil {
			cm.logger.Warn("failed to update cache",
				zap.String("key", cacheKey),
				zap.Error(err))
		} else {
			cm.logger.Debug("updated cache", zap.String("key", cacheKey))
		}
	}(data)

	return nil
}

// ═══════════════════════════════════════════════════════════════
// CACHE INVALIDATION HELPERS
// ═══════════════════════════════════════════════════════════════

// InvalidateSupplierCache invalidates all cache keys related to a supplier.
// This is called after Create/Update/Delete operations.
func (cm *CacheManager) InvalidateSupplierCache(ctx context.Context, supplierID string) error {
	if !cm.enabled {
		return nil
	}

	// Invalidate supplier detail
	detailKey, _ := BuildSupplierDetailCacheKey(supplierID)
	_ = cm.Delete(ctx, detailKey)

	// Invalidate supplier materials
	materialsKey, _ := BuildSupplierMaterialsCacheKey(supplierID)
	_ = cm.Delete(ctx, materialsKey)

	// Invalidate all rating pages for this supplier
	ratingsPattern, _ := BuildSupplierRatingsCachePattern(supplierID)
	_ = cm.DeletePattern(ctx, ratingsPattern)

	// Invalidate all supplier list pages
	listPattern, _ := BuildSupplierListCachePattern()
	_ = cm.DeletePattern(ctx, listPattern)

	// Invalidate stats
	statsKey, _ := BuildSupplierStatsCacheKey()
	_ = cm.Delete(ctx, statsKey)

	cm.logger.Debug("invalidated supplier cache", zap.String("supplierID", supplierID))
	return nil
}

// InvalidateUserCache invalidates all cache keys related to a user.
func (cm *CacheManager) InvalidateUserCache(ctx context.Context, userID string) error {
	if !cm.enabled {
		return nil
	}

	// Invalidate user detail
	detailKey, _ := BuildUserDetailCacheKey(userID)
	_ = cm.Delete(ctx, detailKey)

	// Invalidate user permissions
	permsKey, _ := BuildUserPermissionsCacheKey(userID)
	_ = cm.Delete(ctx, permsKey)

	cm.logger.Debug("invalidated user cache", zap.String("userID", userID))
	return nil
}
