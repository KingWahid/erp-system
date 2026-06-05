package repository

import (
	"context"
	"errors"
	"net/http"

	"erp-system/internal/supplier/domain"
	apperrors "erp-system/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type supplierRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSupplierRepository(db *gorm.DB, logger *zap.Logger) domain.SupplierRepository {
	return &supplierRepository{db: db, logger: logger}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	if err := r.db.WithContext(ctx).Create(supplier).Error; err != nil {
		if apperrors.IsDuplicateKeyError(err) {
			customErr := apperrors.NewCustomError("supplier code already exists").
				WithErrorCode(apperrors.ErrCodeSupplierCodeExists).
				WithMessageID("error_supplier_code_exists").
				WithHTTPCode(http.StatusConflict).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "CreateSupplier"),
				zap.String("code", supplier.Code),
			)
			return customErr
		}
		customErr := apperrors.NewCustomError("failed to create supplier").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "CreateSupplier"),
			zap.String("code", supplier.Code),
		)
		return customErr
	}
	r.logger.Info("supplier created", zap.String("supplier_id", supplier.ID.String()), zap.String("code", supplier.Code))
	return nil
}

func (r *supplierRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	var supplier domain.Supplier
	err := r.db.WithContext(ctx).
		Preload("Contacts").
		Preload("Materials").
		Preload("PerformanceRatings").
		Preload("StageHistories").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&supplier).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customErr := apperrors.NewCustomError("supplier not found").
				WithErrorCode(apperrors.ErrCodeSupplierNotFound).
				WithMessageID("error_supplier_not_found").
				WithHTTPCode(http.StatusNotFound).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "FindSupplierByID"),
				zap.String("supplier_id", id.String()),
			)
			return nil, customErr
		}
		customErr := apperrors.NewCustomError("failed to find supplier").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "FindSupplierByID"),
			zap.String("supplier_id", id.String()),
		)
		return nil, customErr
	}
	return &supplier, nil
}

func (r *supplierRepository) FindByCode(ctx context.Context, code string) (*domain.Supplier, error) {
	var supplier domain.Supplier
	err := r.db.WithContext(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&supplier).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customErr := apperrors.NewCustomError("supplier not found").
				WithErrorCode(apperrors.ErrCodeSupplierNotFound).
				WithMessageID("error_supplier_not_found").
				WithHTTPCode(http.StatusNotFound).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "FindSupplierByCode"),
				zap.String("code", code),
			)
			return nil, customErr
		}
		customErr := apperrors.NewCustomError("failed to find supplier by code").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "FindSupplierByCode"),
			zap.String("code", code),
		)
		return nil, customErr
	}
	return &supplier, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	if err := r.db.WithContext(ctx).Save(supplier).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to update supplier").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpdateSupplier"),
			zap.String("supplier_id", supplier.ID.String()),
		)
		return customErr
	}
	r.logger.Info("supplier updated", zap.String("supplier_id", supplier.ID.String()))
	return nil
}

func (r *supplierRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&domain.Supplier{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to delete supplier").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "SoftDeleteSupplier"),
			zap.String("supplier_id", id.String()),
		)
		return customErr
	}
	r.logger.Info("supplier soft-deleted", zap.String("supplier_id", id.String()))
	return nil
}

func (r *supplierRepository) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Supplier, int64, error) {
	var suppliers []*domain.Supplier
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Supplier{}).
		Preload("Contacts", "is_primary = true").
		Where("deleted_at IS NULL")

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR supplier_no ILIKE ?", like, like, like)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to count suppliers").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger, zap.String("operation", "ListSuppliers"))
		return nil, 0, customErr
	}

	offset := (filter.Page - 1) * filter.PerPage
	if err := query.
		Order("created_at DESC").
		Limit(filter.PerPage).
		Offset(offset).
		Find(&suppliers).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to list suppliers").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger, zap.String("operation", "ListSuppliers"))
		return nil, 0, customErr
	}

	return suppliers, total, nil
}

func (r *supplierRepository) GetStats(ctx context.Context) (*domain.SupplierStats, error) {
	var stats domain.SupplierStats

	r.db.WithContext(ctx).Model(&domain.Supplier{}).
		Where("deleted_at IS NULL").
		Count(&stats.TotalSupplier)

	r.db.WithContext(ctx).Model(&domain.Supplier{}).
		Where("deleted_at IS NULL AND created_at >= NOW() - INTERVAL '30 days'").
		Count(&stats.NewSupplier)

	r.db.WithContext(ctx).Model(&domain.Supplier{}).
		Where("deleted_at IS NULL AND is_blocked = true").
		Count(&stats.BlockedSupplier)

	return &stats, nil
}

func (r *supplierRepository) UpsertMaterials(ctx context.Context, supplierID uuid.UUID, materials []domain.SupplierMaterial) error {
	if len(materials) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&materials).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to upsert materials").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpsertMaterials"),
			zap.String("supplier_id", supplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) DeleteMaterials(ctx context.Context, supplierID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Delete(&domain.SupplierMaterial{}).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to delete materials").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "DeleteMaterials"),
			zap.String("supplier_id", supplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) UpsertContacts(ctx context.Context, supplierID uuid.UUID, contacts []domain.SupplierContact) error {
	if len(contacts) == 0 {
		return nil
	}
	r.db.WithContext(ctx).Where("supplier_id = ?", supplierID).Delete(&domain.SupplierContact{})
	if err := r.db.WithContext(ctx).Create(&contacts).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to upsert contacts").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpsertContacts"),
			zap.String("supplier_id", supplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) AddPerformanceRating(ctx context.Context, rating *domain.PerformanceRating) error {
	if err := r.db.WithContext(ctx).Create(rating).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to add performance rating").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "AddPerformanceRating"),
			zap.String("supplier_id", rating.SupplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) GetPerformanceRatings(ctx context.Context, supplierID uuid.UUID) ([]*domain.PerformanceRating, error) {
	var ratings []*domain.PerformanceRating
	if err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("created_at DESC").
		Find(&ratings).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to get performance ratings").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetPerformanceRatings"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, customErr
	}
	return ratings, nil
}

func (r *supplierRepository) AddStageHistory(ctx context.Context, history *domain.SupplierStageHistory) error {
	if err := r.db.WithContext(ctx).Create(history).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to add stage history").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "AddStageHistory"),
			zap.String("supplier_id", history.SupplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) GetStageHistory(ctx context.Context, supplierID uuid.UUID) ([]*domain.SupplierStageHistory, error) {
	var histories []*domain.SupplierStageHistory
	if err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("created_at ASC").
		Find(&histories).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to get stage history").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetStageHistory"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, customErr
	}
	return histories, nil
}

func (r *supplierRepository) GetAddresses(ctx context.Context, supplierID uuid.UUID) ([]domain.SupplierAddress, error) {
	var addresses []domain.SupplierAddress
	if err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("created_at ASC").
		Find(&addresses).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to get supplier addresses").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetAddresses"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, customErr
	}
	return addresses, nil
}

func (r *supplierRepository) AddAddress(ctx context.Context, address *domain.SupplierAddress) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to add supplier address").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "AddAddress"),
			zap.String("supplier_id", address.SupplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) UpdateAddress(ctx context.Context, address *domain.SupplierAddress) error {
	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to update supplier address").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpdateAddress"),
			zap.String("address_id", address.ID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) DeleteAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND supplier_id = ?", addressID, supplierID).
		Delete(&domain.SupplierAddress{}).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to delete supplier address").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "DeleteAddress"),
			zap.String("supplier_id", supplierID.String()),
			zap.String("address_id", addressID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) SetMainAddress(ctx context.Context, supplierID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.SupplierAddress{}).
			Where("supplier_id = ?", supplierID).
			Update("is_main", false).Error; err != nil {
			customErr := apperrors.NewCustomError("failed to unset main addresses").
				WithErrorCode(apperrors.ErrCodeDatabaseTransaction).
				WithMessageID("error_database_query").
				WithHTTPCode(http.StatusInternalServerError).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "SetMainAddress"),
				zap.String("supplier_id", supplierID.String()),
			)
			return customErr
		}
		if err := tx.Model(&domain.SupplierAddress{}).
			Where("id = ? AND supplier_id = ?", addressID, supplierID).
			Update("is_main", true).Error; err != nil {
			customErr := apperrors.NewCustomError("failed to set main address").
				WithErrorCode(apperrors.ErrCodeDatabaseTransaction).
				WithMessageID("error_database_query").
				WithHTTPCode(http.StatusInternalServerError).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "SetMainAddress"),
				zap.String("address_id", addressID.String()),
			)
			return customErr
		}
		return nil
	})
}

func (r *supplierRepository) GetContacts(ctx context.Context, supplierID uuid.UUID) ([]domain.SupplierContact, error) {
	var contacts []domain.SupplierContact
	if err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("created_at ASC").
		Find(&contacts).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to get supplier contacts").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetContacts"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, customErr
	}
	return contacts, nil
}

func (r *supplierRepository) AddContact(ctx context.Context, contact *domain.SupplierContact) error {
	if err := r.db.WithContext(ctx).Create(contact).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to add supplier contact").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "AddContact"),
			zap.String("supplier_id", contact.SupplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) UpdateContact(ctx context.Context, contact *domain.SupplierContact) error {
	if err := r.db.WithContext(ctx).Save(contact).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to update supplier contact").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpdateContact"),
			zap.String("contact_id", contact.ID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) DeleteContact(ctx context.Context, supplierID, contactID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND supplier_id = ?", contactID, supplierID).
		Delete(&domain.SupplierContact{}).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to delete supplier contact").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "DeleteContact"),
			zap.String("supplier_id", supplierID.String()),
			zap.String("contact_id", contactID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) SetPrimaryContact(ctx context.Context, supplierID, contactID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.SupplierContact{}).
			Where("supplier_id = ?", supplierID).
			Update("is_primary", false).Error; err != nil {
			customErr := apperrors.NewCustomError("failed to unset primary contacts").
				WithErrorCode(apperrors.ErrCodeDatabaseTransaction).
				WithMessageID("error_database_query").
				WithHTTPCode(http.StatusInternalServerError).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "SetPrimaryContact"),
				zap.String("supplier_id", supplierID.String()),
			)
			return customErr
		}
		if err := tx.Model(&domain.SupplierContact{}).
			Where("id = ? AND supplier_id = ?", contactID, supplierID).
			Update("is_primary", true).Error; err != nil {
			customErr := apperrors.NewCustomError("failed to set primary contact").
				WithErrorCode(apperrors.ErrCodeDatabaseTransaction).
				WithMessageID("error_database_query").
				WithHTTPCode(http.StatusInternalServerError).
				WithCause(err)
			customErr.LogError(r.logger,
				zap.String("operation", "SetPrimaryContact"),
				zap.String("contact_id", contactID.String()),
			)
			return customErr
		}
		return nil
	})
}

func (r *supplierRepository) GetGroups(ctx context.Context, supplierID uuid.UUID) ([]domain.SupplierGroup, error) {
	var groups []domain.SupplierGroup
	if err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("created_at ASC").
		Find(&groups).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to get supplier groups").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetGroups"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, customErr
	}
	return groups, nil
}

func (r *supplierRepository) AddGroup(ctx context.Context, group *domain.SupplierGroup) error {
	if err := r.db.WithContext(ctx).Create(group).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to add supplier group").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "AddGroup"),
			zap.String("supplier_id", group.SupplierID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) UpdateGroup(ctx context.Context, group *domain.SupplierGroup) error {
	if err := r.db.WithContext(ctx).Save(group).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to update supplier group").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "UpdateGroup"),
			zap.String("group_id", group.ID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) DeleteGroup(ctx context.Context, supplierID, groupID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND supplier_id = ?", groupID, supplierID).
		Delete(&domain.SupplierGroup{}).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to delete supplier group").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "DeleteGroup"),
			zap.String("supplier_id", supplierID.String()),
			zap.String("group_id", groupID.String()),
		)
		return customErr
	}
	return nil
}

func (r *supplierRepository) GetStageHistories(ctx context.Context, supplierID uuid.UUID) ([]*domain.SupplierStageHistory, error) {
	return r.GetStageHistory(ctx, supplierID)
}

func (r *supplierRepository) GetOutstandingInvoices(ctx context.Context, supplierID uuid.UUID, page, perPage int) ([]*domain.SupplierInvoice, int64, error) {
	var invoices []*domain.SupplierInvoice
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.SupplierInvoice{}).
		Where("supplier_id = ? AND status IN ('unpaid', 'partial', 'overdue')", supplierID)

	if err := query.Count(&total).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to count outstanding invoices").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetOutstandingInvoices"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, 0, customErr
	}

	offset := (page - 1) * perPage
	if err := query.
		Order("due_date ASC").
		Limit(perPage).
		Offset(offset).
		Find(&invoices).Error; err != nil {
		customErr := apperrors.NewCustomError("failed to get outstanding invoices").
			WithErrorCode(apperrors.ErrCodeDatabaseQuery).
			WithMessageID("error_database_query").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
		customErr.LogError(r.logger,
			zap.String("operation", "GetOutstandingInvoices"),
			zap.String("supplier_id", supplierID.String()),
		)
		return nil, 0, customErr
	}

	return invoices, total, nil
}
