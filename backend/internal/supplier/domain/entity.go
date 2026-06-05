package domain

import (
	"time"

	"github.com/google/uuid"
)

// SupplierStatus represents supplier lifecycle status
type SupplierStatus string

const (
	StatusActive     SupplierStatus = "active"
	StatusInProgress SupplierStatus = "in_progress"
	StatusBlocked    SupplierStatus = "blocked"
	StatusDraft      SupplierStatus = "draft"
	StatusInactive   SupplierStatus = "inactive"
)

// SupplierStage represents the onboarding workflow stage
type SupplierStage string

const (
	StageDraft      SupplierStage = "draft"
	StageInReview   SupplierStage = "in_review"
	StageAssessment SupplierStage = "in_assessment"
	StageActive     SupplierStage = "active"
)

// Supplier is the core domain entity
type Supplier struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code        string    `gorm:"uniqueIndex;not null"`
	SupplierNo  string    `gorm:"uniqueIndex;not null"`
	Name        string    `gorm:"not null"`
	Alias       string    // Nick Name in UI
	LogoURL     string
	Address     string // main address (legacy single field)
	City        string
	Province    string
	Country     string `gorm:"default:'Indonesia'"`
	PostalCode  string
	Phone       string
	Email       string
	Website     string
	Status      SupplierStatus `gorm:"type:varchar(50);default:'draft'"`
	Stage       SupplierStage  `gorm:"type:varchar(50);default:'draft'"`
	SLAHours    int            `gorm:"default:72"`
	IsBlocked   bool           `gorm:"default:false"`
	BlockReason string
	Notes       string
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`

	// Relations
	Addresses          []SupplierAddress      `gorm:"foreignKey:SupplierID"`
	Contacts           []SupplierContact      `gorm:"foreignKey:SupplierID"`
	Groups             []SupplierGroup        `gorm:"foreignKey:SupplierID"`
	Materials          []SupplierMaterial     `gorm:"foreignKey:SupplierID"`
	PerformanceRatings []PerformanceRating    `gorm:"foreignKey:SupplierID"`
	StageHistories     []SupplierStageHistory `gorm:"foreignKey:SupplierID"`
}

func (Supplier) TableName() string { return "suppliers" }

// SupplierAddress stores multiple office addresses per supplier.
// Screenshot: Address tab — Name (Head Office/Branch), Address, Main radio
type SupplierAddress struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name       string    `gorm:"not null"` // e.g. "Head Office", "Branch Office"
	Address    string    `gorm:"not null"`
	City       string
	Province   string
	Country    string `gorm:"default:'Indonesia'"`
	PostalCode string
	IsMain     bool `gorm:"default:false"` // radio "Main"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (SupplierAddress) TableName() string { return "supplier_addresses" }

// SupplierContact stores contact persons for a supplier.
// Screenshot: Contacts tab — Name, Job Position, Email, Phone, Mobile, Main
type SupplierContact struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name       string    `gorm:"not null"`
	Position   string    // Job Position
	Phone      string    // Office phone
	Mobile     string    // Mobile phone
	Email      string
	IsPrimary  bool `gorm:"default:false"` // radio "Main"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (SupplierContact) TableName() string { return "supplier_contacts" }

// SupplierGroup stores group/classification tags for a supplier.
// Screenshot: Groups tab — Group Name, Value, Active checkbox
type SupplierGroup struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID uuid.UUID `gorm:"type:uuid;not null;index"`
	GroupName  string    `gorm:"not null"` // e.g. "Industry", "Telkom Group"
	Value      string    `gorm:"not null"` // e.g. "Manufacture", "Non Telkom Group"
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (SupplierGroup) TableName() string { return "supplier_groups" }

// SupplierMaterial stores materials provided by a supplier.
// Screenshot: Material List tab — Material Group, Material ID, Active
type SupplierMaterial struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID    uuid.UUID `gorm:"type:uuid;not null;index"`
	MaterialGroup string    `gorm:"not null"` // e.g. "IT - Device"
	MaterialID    string    `gorm:"not null"` // e.g. "Computer / Notebook"
	IsActive      bool      `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (SupplierMaterial) TableName() string { return "supplier_materials" }

// PerformanceRating stores supplier performance reviews.
// Screenshot: Performance Ratings — Price, Delivery Time, Notes, Reviewed by
type PerformanceRating struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID     uuid.UUID `gorm:"type:uuid;not null;index"`
	PriceRating    int       `gorm:"check:price_rating >= 1 AND price_rating <= 5"`
	DeliveryRating int       `gorm:"check:delivery_rating >= 1 AND delivery_rating <= 5"`
	Notes          string
	ReviewedBy     string
	ReviewedAt     time.Time
	CreatedAt      time.Time
}

func (PerformanceRating) TableName() string { return "supplier_performance_ratings" }

// SupplierStageHistory tracks stage transitions.
// Screenshot: Stage widget — Draft → In Review → In Assessment → Active
type SupplierStageHistory struct {
	ID         uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID uuid.UUID     `gorm:"type:uuid;not null;index"`
	FromStage  SupplierStage `gorm:"type:varchar(50)"`
	ToStage    SupplierStage `gorm:"type:varchar(50);not null"`
	Notes      string
	ChangedBy  string
	ElapsedMs  int64
	CreatedAt  time.Time
}

func (SupplierStageHistory) TableName() string { return "supplier_stage_histories" }

// SupplierInvoice stores invoice records for a supplier.
// Status: unpaid | partial | paid | overdue
type SupplierInvoice struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SupplierID    uuid.UUID `gorm:"type:uuid;not null;index"`
	InvoiceNumber string    `gorm:"not null"`
	ProjectName   string
	Amount        float64 `gorm:"type:numeric(18,2)"`
	Currency      string  `gorm:"default:'IDR'"`
	InvoiceDate   time.Time
	DueDate       time.Time
	PaidDate      *time.Time // nil = not yet paid
	Status        string     `gorm:"default:'unpaid'"` // unpaid|partial|paid|overdue
	PaidAmount    float64    `gorm:"type:numeric(18,2)"`
	Notes         string
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (SupplierInvoice) TableName() string { return "supplier_invoices" }

// SupplierStats is a read-model for the list page stats cards.
// Screenshot: Total Supplier, New Supplier, Avg Cost per Supplier, Blocked Supplier
type SupplierStats struct {
	TotalSupplier   int64   `json:"total_supplier"`
	NewSupplier     int64   `json:"new_supplier"`
	AvgCostSupplier float64 `json:"avg_cost_supplier"`
	BlockedSupplier int64   `json:"blocked_supplier"`
}
