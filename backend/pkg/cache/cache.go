package cache

import (
	"context"
	"time"
)

// Cache is a generic interface for caching operations.
// Implementation-agnostic (Redis, Memcached, in-memory, etc.)
type Cache interface {
	// Get retrieves a value by key. Returns nil if not found.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value with TTL. TTL=0 means no expiration.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes a key.
	Delete(ctx context.Context, key string) error

	// DeletePattern removes all keys matching a pattern (e.g., "supplier:*")
	DeletePattern(ctx context.Context, pattern string) error

	// Exists checks if a key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// Ping checks connection health.
	Ping(ctx context.Context) error

	// Close closes the cache connection.
	Close() error
}

// TTL constants for different cache tiers
const (
	TTLShort      = 2 * time.Minute   // Stats, dashboards
	TTLMedium     = 15 * time.Minute  // List queries, search
	TTLLong       = 1 * time.Hour     // Detail views
	TTLVeryLong   = 24 * time.Hour    // Permissions, config
	TTLPersistent = 0                 // No expiration
)
