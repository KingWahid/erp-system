package usecase

import (
	"context"

	"erp-system/internal/supplier/domain"

	"github.com/google/uuid"
)

// SupplierUsecase defines the business logic contract
type SupplierUsecase interface {
	// List & Stats
	ListSuppliers(ctx context.Context, req ListRequest) (*ListResponse, error)
	GetStats(ctx context.Context) (*domain.SupplierStats, error)

	// CRUD
	CreateSupplier(ctx context.Context, req CreateSupplierRequest) (*SupplierDetailResponse, error)
	GetSupplierByID(ctx context.Context, id uuid.UUID) (*SupplierDetailResponse, error)
	UpdateSupplier(ctx context.Context, id uuid.UUID, req UpdateSupplierRequest) (*SupplierDetailResponse, error)
	DeleteSupplier(ctx context.Context, id uuid.UUID) error

	// Block / Unblock
	BlockSupplier(ctx context.Context, id uuid.UUID, reason string) error
	UnblockSupplier(ctx context.Context, id uuid.UUID) error

	// Stage workflow
	AdvanceStage(ctx context.Context, id uuid.UUID, notes string, changedBy string) (*SupplierDetailResponse, error)

	// Materials
	UpdateMaterials(ctx context.Context, id uuid.UUID, req []MaterialRequest) error

	// Performance Rating
	AddPerformanceRating(ctx context.Context, id uuid.UUID, req RatingRequest) error
	GetPerformanceRatings(ctx context.Context, id uuid.UUID) ([]*RatingResponse, error)

	// Addresses
	GetSupplierAddresses(ctx context.Context, supplierID uuid.UUID) ([]AddressResponse, error)
	AddSupplierAddress(ctx context.Context, supplierID uuid.UUID, req AddressRequest) (*AddressResponse, error)
	UpdateSupplierAddress(ctx context.Context, supplierID, addressID uuid.UUID, req AddressRequest) (*AddressResponse, error)
	DeleteSupplierAddress(ctx context.Context, supplierID, addressID uuid.UUID) error
	SetMainAddress(ctx context.Context, supplierID, addressID uuid.UUID) error

	// Contacts
	GetSupplierContacts(ctx context.Context, supplierID uuid.UUID) ([]ContactResponse, error)
	AddSupplierContact(ctx context.Context, supplierID uuid.UUID, req ContactRequest) (*ContactResponse, error)
	UpdateSupplierContact(ctx context.Context, supplierID, contactID uuid.UUID, req ContactRequest) (*ContactResponse, error)
	DeleteSupplierContact(ctx context.Context, supplierID, contactID uuid.UUID) error
	SetPrimaryContact(ctx context.Context, supplierID, contactID uuid.UUID) error

	// Groups
	GetSupplierGroups(ctx context.Context, supplierID uuid.UUID) ([]GroupResponse, error)
	AddSupplierGroup(ctx context.Context, supplierID uuid.UUID, req GroupRequest) (*GroupResponse, error)
	UpdateSupplierGroup(ctx context.Context, supplierID, groupID uuid.UUID, req GroupRequest) (*GroupResponse, error)
	DeleteSupplierGroup(ctx context.Context, supplierID, groupID uuid.UUID) error

	// Stage History
	GetStageHistories(ctx context.Context, supplierID uuid.UUID) ([]StageHistoryResponse, error)

	// Invoices / Outstandings
	GetOutstandingInvoices(ctx context.Context, supplierID uuid.UUID, page, perPage int) ([]InvoiceOutstandingResponse, int64, error)
}
