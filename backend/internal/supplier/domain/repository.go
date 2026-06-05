package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../mocks/repository_mock.go

// ListFilter defines query parameters for supplier list
type ListFilter struct {
	Search  string
	Status  SupplierStatus
	Page    int
	PerPage int
}

// SupplierRepository defines the data access contract
type SupplierRepository interface {
	// Core CRUD
	Create(ctx context.Context, supplier *Supplier) error
	FindByID(ctx context.Context, id uuid.UUID) (*Supplier, error)
	FindByCode(ctx context.Context, code string) (*Supplier, error)
	Update(ctx context.Context, supplier *Supplier) error
	SoftDelete(ctx context.Context, id uuid.UUID) error

	// List & Stats
	List(ctx context.Context, filter ListFilter) ([]*Supplier, int64, error)
	GetStats(ctx context.Context) (*SupplierStats, error)

	// Materials
	UpsertMaterials(ctx context.Context, supplierID uuid.UUID, materials []SupplierMaterial) error
	DeleteMaterials(ctx context.Context, supplierID uuid.UUID) error

	// Contacts (legacy bulk upsert)
	UpsertContacts(ctx context.Context, supplierID uuid.UUID, contacts []SupplierContact) error

	// Addresses (new granular CRUD)
	GetAddresses(ctx context.Context, supplierID uuid.UUID) ([]SupplierAddress, error)
	AddAddress(ctx context.Context, address *SupplierAddress) error
	UpdateAddress(ctx context.Context, address *SupplierAddress) error
	DeleteAddress(ctx context.Context, supplierID, addressID uuid.UUID) error
	SetMainAddress(ctx context.Context, supplierID, addressID uuid.UUID) error

	// Contacts (new granular CRUD)
	GetContacts(ctx context.Context, supplierID uuid.UUID) ([]SupplierContact, error)
	AddContact(ctx context.Context, contact *SupplierContact) error
	UpdateContact(ctx context.Context, contact *SupplierContact) error
	DeleteContact(ctx context.Context, supplierID, contactID uuid.UUID) error
	SetPrimaryContact(ctx context.Context, supplierID, contactID uuid.UUID) error

	// Groups
	GetGroups(ctx context.Context, supplierID uuid.UUID) ([]SupplierGroup, error)
	AddGroup(ctx context.Context, group *SupplierGroup) error
	UpdateGroup(ctx context.Context, group *SupplierGroup) error
	DeleteGroup(ctx context.Context, supplierID, groupID uuid.UUID) error

	// Performance
	AddPerformanceRating(ctx context.Context, rating *PerformanceRating) error
	GetPerformanceRatings(ctx context.Context, supplierID uuid.UUID) ([]*PerformanceRating, error)

	// Stage
	AddStageHistory(ctx context.Context, history *SupplierStageHistory) error
	GetStageHistory(ctx context.Context, supplierID uuid.UUID) ([]*SupplierStageHistory, error)
	GetStageHistories(ctx context.Context, supplierID uuid.UUID) ([]*SupplierStageHistory, error)

	// Invoices / Outstandings
	GetOutstandingInvoices(ctx context.Context, supplierID uuid.UUID, page, perPage int) ([]*SupplierInvoice, int64, error)
}
