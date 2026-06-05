package usecase

import (
	"time"

	"github.com/google/uuid"
)

// REQUEST DTOs

type ListRequest struct {
	Search  string `query:"search"`
	Status  string `query:"status"`
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
}

type CreateSupplierRequest struct {
	Name    string `json:"name"    validate:"required"`
	Code    string `json:"code"    validate:"required,max=10"`
	Alias   string `json:"alias"` // Nick Name in UI
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
	Email   string `json:"email"   validate:"omitempty,email"`
	Website string `json:"website"`
	Notes   string `json:"notes"`
}

type UpdateSupplierRequest struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Notes   string `json:"notes"`
}

// AddressRequest — Address tab in modal
type AddressRequest struct {
	Name       string `json:"name"    validate:"required"` // e.g. "Head Office"
	Address    string `json:"address" validate:"required"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	IsMain     bool   `json:"is_main"`
}

// ContactRequest — Contacts tab in modal
type ContactRequest struct {
	Name      string `json:"name"     validate:"required"`
	Position  string `json:"position"` // Job Position
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"` // NEW: mobile phone
	Email     string `json:"email"`
	IsPrimary bool   `json:"is_primary"`
}

// GroupRequest — Groups tab in modal
type GroupRequest struct {
	GroupName string `json:"group_name" validate:"required"` // e.g. "Industry"
	Value     string `json:"value"      validate:"required"` // e.g. "Manufacture"
	IsActive  bool   `json:"is_active"`
}

type MaterialRequest struct {
	MaterialGroup string `json:"material_group" validate:"required"`
	MaterialID    string `json:"material_id"    validate:"required"`
	IsActive      bool   `json:"is_active"`
}

type RatingRequest struct {
	PriceRating    int    `json:"price_rating"    validate:"required,min=1,max=5"`
	DeliveryRating int    `json:"delivery_rating" validate:"required,min=1,max=5"`
	Notes          string `json:"notes"`
	ReviewedBy     string `json:"reviewed_by"`
}

// RESPONSE DTOs

type ListResponse struct {
	Items   []*SupplierListItem `json:"items"`
	Total   int64               `json:"total"`
	Page    int                 `json:"page"`
	PerPage int                 `json:"per_page"`
}

type SupplierListItem struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	SupplierNo string    `json:"supplier_no"`
	Name       string    `json:"name"`
	Alias      string    `json:"alias"`
	Address    string    `json:"address"` // main address city, country
	Contact    string    `json:"contact"` // primary contact name
	Status     string    `json:"status"`
}

type SupplierDetailResponse struct {
	ID             uuid.UUID              `json:"id"`
	Code           string                 `json:"code"`
	SupplierNo     string                 `json:"supplier_no"`
	Name           string                 `json:"name"`
	Alias          string                 `json:"alias"`
	LogoURL        string                 `json:"logo_url"`
	Address        string                 `json:"address"`
	City           string                 `json:"city"`
	Country        string                 `json:"country"`
	Phone          string                 `json:"phone"`
	Email          string                 `json:"email"`
	Website        string                 `json:"website"`
	Status         string                 `json:"status"`
	Stage          string                 `json:"stage"`
	SLAHours       int                    `json:"sla_hours"`
	IsBlocked      bool                   `json:"is_blocked"`
	BlockReason    string                 `json:"block_reason,omitempty"`
	Notes          string                 `json:"notes"`
	Addresses      []AddressResponse      `json:"addresses"`
	Contacts       []ContactResponse      `json:"contacts"`
	Groups         []GroupResponse        `json:"groups"`
	Materials      []MaterialResponse     `json:"materials"`
	StageHistories []StageHistoryResponse `json:"stage_histories"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AddressResponse — Address tab
type AddressResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"` // e.g. "Head Office"
	Address    string    `json:"address"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	IsMain     bool      `json:"is_main"`
}

// ContactResponse — Contacts tab
type ContactResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"` // Job Position
	Phone     string    `json:"phone"`
	Mobile    string    `json:"mobile"` // Mobile phone
	Email     string    `json:"email"`
	IsPrimary bool      `json:"is_primary"`
}

// GroupResponse — Groups tab
type GroupResponse struct {
	ID        uuid.UUID `json:"id"`
	GroupName string    `json:"group_name"`
	Value     string    `json:"value"`
	IsActive  bool      `json:"is_active"`
}

type MaterialResponse struct {
	ID            uuid.UUID `json:"id"`
	MaterialGroup string    `json:"material_group"`
	MaterialID    string    `json:"material_id"`
	IsActive      bool      `json:"is_active"`
}

type RatingResponse struct {
	ID             uuid.UUID `json:"id"`
	PriceRating    int       `json:"price_rating"`
	DeliveryRating int       `json:"delivery_rating"`
	Notes          string    `json:"notes"`
	ReviewedBy     string    `json:"reviewed_by"`
	ReviewedAt     time.Time `json:"reviewed_at"`
}

type StageHistoryResponse struct {
	ID        uuid.UUID `json:"id"`
	FromStage string    `json:"from_stage"`
	ToStage   string    `json:"to_stage"`
	Notes     string    `json:"notes"`
	ChangedBy string    `json:"changed_by"`
	ElapsedMs int64     `json:"elapsed_ms"`
	CreatedAt time.Time `json:"created_at"`
}

// InvoiceOutstandingResponse — Outstandings section in Overview tab
type InvoiceOutstandingResponse struct {
	ID            uuid.UUID `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	ProjectName   string    `json:"project_name"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	InvoiceDate   time.Time `json:"invoice_date"`
	DueDate       time.Time `json:"due_date"`
	AgingDays     int       `json:"aging_days"`
	Status        string    `json:"status"`
	PaidAmount    float64   `json:"paid_amount"`
}
