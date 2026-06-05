package repository

import (
	"context"
	"net/http"

	"erp-system/internal/supplier/domain"
	"erp-system/pkg/cache"
	apperrors "erp-system/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CachedSupplierRepository wraps SupplierRepository with caching layer.
// Implements cache-aside pattern for read operations.
type CachedSupplierRepository struct {
	repo         domain.SupplierRepository
	cacheManager *cache.CacheManager
	logger       *zap.Logger
}

// NewCachedSupplierRepository creates a cached repository decorator.
func NewCachedSupplierRepository(
	repo domain.SupplierRepository,
	cacheManager *cache.CacheManager,
	logger *zap.Logger,
) domain.SupplierRepository {
	return &CachedSupplierRepository{
		repo:         repo,
		cacheManager: cacheManager,
		logger:       logger,
	}
}

// ═══════════════════════════════════════════════════════════════
// READ OPERATIONS (with caching)
// ═══════════════════════════════════════════════════════════════

func (r *CachedSupplierRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	cacheKey, _ := cache.BuildSupplierDetailCacheKey(id.String())

	var supplier *domain.Supplier
	err := r.cacheManager.GetFromCacheOrDB(
		ctx,
		cacheKey,
		cache.SupplierDetailCacheTTL,
		&supplier,
		func() error {
			var dbErr error
			supplier, dbErr = r.repo.FindByID(ctx, id)
			return dbErr
		},
	)
	if err != nil {
		return nil, wrapCacheError(err, "FindByID", id.String())
	}

	return supplier, nil
}

func (r *CachedSupplierRepository) FindByCode(ctx context.Context, code string) (*domain.Supplier, error) {
	// Code lookup is less frequent, no caching for now
	return r.repo.FindByCode(ctx, code)
}

func (r *CachedSupplierRepository) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Supplier, int64, error) {
	cacheKey, _ := cache.BuildSupplierListCacheKey(
		filter.Page,
		filter.PerPage,
		string(filter.Status),
		filter.Search,
	)

	type listResult struct {
		Suppliers []*domain.Supplier `json:"suppliers"`
		Total     int64              `json:"total"`
	}

	var result listResult
	err := r.cacheManager.GetFromCacheOrDB(
		ctx,
		cacheKey,
		cache.SupplierListCacheTTL,
		&result,
		func() error {
			suppliers, total, dbErr := r.repo.List(ctx, filter)
			if dbErr == nil {
				result.Suppliers = suppliers
				result.Total = total
			}
			return dbErr
		},
	)
	if err != nil {
		return nil, 0, err
	}

	return result.Suppliers, result.Total, nil
}

func (r *CachedSupplierRepository) GetStats(ctx context.Context) (*domain.SupplierStats, error) {
	cacheKey, _ := cache.BuildSupplierStatsCacheKey()

	var stats *domain.SupplierStats
	err := r.cacheManager.GetFromCacheOrDB(
		ctx,
		cacheKey,
		cache.SupplierStatsCacheTTL,
		&stats,
		func() error {
			var dbErr error
			stats, dbErr = r.repo.GetStats(ctx)
			return dbErr
		},
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *CachedSupplierRepository) GetPerformanceRatings(ctx context.Context, supplierID uuid.UUID) ([]*domain.PerformanceRating, error) {
	cacheKey, _ := cache.BuildSupplierRatingsCacheKey(supplierID.String(), 1, 100)

	var ratings []*domain.PerformanceRating
	err := r.cacheManager.GetFromCacheOrDB(
		ctx,
		cacheKey,
		cache.SupplierRatingsCacheTTL,
		&ratings,
		func() error {
			var dbErr error
			ratings, dbErr = r.repo.GetPerformanceRatings(ctx, supplierID)
			return dbErr
		},
	)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (r *CachedSupplierRepository) GetStageHistory(ctx context.Context, supplierID uuid.UUID) ([]*domain.SupplierStageHistory, error) {
	return r.repo.GetStageHistory(ctx, supplierID)
}

func (r *CachedSupplierRepository) GetAddresses(ctx context.Context, supplierID uuid.UUID) ([]domain.SupplierAddress, error) {
	return r.repo.GetAddresses(ctx, supplierID)
}

func (r *CachedSupplierRepository) AddAddress(ctx context.Context, address *domain.SupplierAddress) error {
	return r.repo.AddAddress(ctx, address)
}

func (r *CachedSupplierRepository) UpdateAddress(ctx context.Context, address *domain.SupplierAddress) error {
	return r.repo.UpdateAddress(ctx, address)
}

func (r *CachedSupplierRepository) DeleteAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	return r.repo.DeleteAddress(ctx, supplierID, addressID)
}

func (r *CachedSupplierRepository) SetMainAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	return r.repo.SetMainAddress(ctx, supplierID, addressID)
}

func (r *CachedSupplierRepository) GetContacts(ctx context.Context, supplierID uuid.UUID) ([]domain.SupplierContact, error) {
	return r.repo.GetContacts(ctx, supplierID)
}

func (r *CachedSupplierRepository) AddContact(ctx context.Context, contact *domain.SupplierContact) error {
	return r.repo.AddContact(ctx, contact)
}

func (r *CachedSupplierRepository) UpdateContact(ctx context.Context, contact *domain.SupplierContact) error {
	return r.repo.UpdateContact(ctx, contact)
}

func (r *CachedSupplierRepository) DeleteContact(ctx context.Context, supplierID, contactID uuid.UUID) error {
	return r.repo.DeleteContact(ctx, supplierID, contactID)
}

func (r *CachedSupplierRepository) SetPrimaryContact(ctx context.Context, supplierID, contactID uuid.UUID) error {
	return r.repo.SetPrimaryContact(ctx, supplierID, contactID)
}

func (r *CachedSupplierRepository) GetGroups(ctx context.Context, supplierID uuid.UUID) ([]domain.SupplierGroup, error) {
	return r.repo.GetGroups(ctx, supplierID)
}

func (r *CachedSupplierRepository) AddGroup(ctx context.Context, group *domain.SupplierGroup) error {
	return r.repo.AddGroup(ctx, group)
}

func (r *CachedSupplierRepository) UpdateGroup(ctx context.Context, group *domain.SupplierGroup) error {
	return r.repo.UpdateGroup(ctx, group)
}

func (r *CachedSupplierRepository) DeleteGroup(ctx context.Context, supplierID, groupID uuid.UUID) error {
	return r.repo.DeleteGroup(ctx, supplierID, groupID)
}

func (r *CachedSupplierRepository) GetStageHistories(ctx context.Context, supplierID uuid.UUID) ([]*domain.SupplierStageHistory, error) {
	return r.repo.GetStageHistories(ctx, supplierID)
}

func (r *CachedSupplierRepository) GetOutstandingInvoices(ctx context.Context, supplierID uuid.UUID, page, perPage int) ([]*domain.SupplierInvoice, int64, error) {
	return r.repo.GetOutstandingInvoices(ctx, supplierID, page, perPage)
}

// ═══════════════════════════════════════════════════════════════
// WRITE OPERATIONS (invalidate cache)
// ═══════════════════════════════════════════════════════════════

func (r *CachedSupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	err := r.repo.Create(ctx, supplier)
	if err == nil {
		_ = r.cacheManager.InvalidateSupplierCache(ctx, supplier.ID.String())
	}
	return err
}

func (r *CachedSupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	err := r.repo.Update(ctx, supplier)
	if err == nil {
		_ = r.cacheManager.InvalidateSupplierCache(ctx, supplier.ID.String())
	}
	return err
}

func (r *CachedSupplierRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.repo.SoftDelete(ctx, id)
	if err == nil {
		_ = r.cacheManager.InvalidateSupplierCache(ctx, id.String())
	}
	return err
}

func (r *CachedSupplierRepository) UpsertMaterials(ctx context.Context, supplierID uuid.UUID, materials []domain.SupplierMaterial) error {
	err := r.repo.UpsertMaterials(ctx, supplierID, materials)
	if err == nil {
		materialsKey, _ := cache.BuildSupplierMaterialsCacheKey(supplierID.String())
		_ = r.cacheManager.Delete(ctx, materialsKey)
		detailKey, _ := cache.BuildSupplierDetailCacheKey(supplierID.String())
		_ = r.cacheManager.Delete(ctx, detailKey)
	}
	return err
}

func (r *CachedSupplierRepository) DeleteMaterials(ctx context.Context, supplierID uuid.UUID) error {
	err := r.repo.DeleteMaterials(ctx, supplierID)
	if err == nil {
		materialsKey, _ := cache.BuildSupplierMaterialsCacheKey(supplierID.String())
		_ = r.cacheManager.Delete(ctx, materialsKey)
	}
	return err
}

func (r *CachedSupplierRepository) UpsertContacts(ctx context.Context, supplierID uuid.UUID, contacts []domain.SupplierContact) error {
	err := r.repo.UpsertContacts(ctx, supplierID, contacts)
	if err == nil {
		detailKey, _ := cache.BuildSupplierDetailCacheKey(supplierID.String())
		_ = r.cacheManager.Delete(ctx, detailKey)
	}
	return err
}

func (r *CachedSupplierRepository) AddPerformanceRating(ctx context.Context, rating *domain.PerformanceRating) error {
	err := r.repo.AddPerformanceRating(ctx, rating)
	if err == nil {
		ratingsPattern, _ := cache.BuildSupplierRatingsCachePattern(rating.SupplierID.String())
		_ = r.cacheManager.DeletePattern(ctx, ratingsPattern)
	}
	return err
}

func (r *CachedSupplierRepository) AddStageHistory(ctx context.Context, history *domain.SupplierStageHistory) error {
	err := r.repo.AddStageHistory(ctx, history)
	if err == nil {
		detailKey, _ := cache.BuildSupplierDetailCacheKey(history.SupplierID.String())
		_ = r.cacheManager.Delete(ctx, detailKey)
	}
	return err
}

// wrapCacheError wraps a non-CustomError from cache layer operations.
// If err is already a CustomError, it is returned as-is.
func wrapCacheError(err error, operation, entityID string) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*apperrors.CustomError); ok {
		return err
	}
	return apperrors.NewCustomError("cache operation failed: " + operation).
		WithErrorCode(apperrors.ErrCodeCacheOperation).
		WithMessageID("error_cache_operation").
		WithHTTPCode(http.StatusInternalServerError).
		WithCause(err)
}
