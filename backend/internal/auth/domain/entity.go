package domain

import (
	"time"

	"github.com/google/uuid"
)

// ═══════════════════════════════════════════════════════════════
// RBAC Entities
// ═══════════════════════════════════════════════════════════════

// User represents a system user with many-to-many role assignment
type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	Password    string    `gorm:"not null"`
	IsActive    bool      `gorm:"default:true"`
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations
	Roles []Role `gorm:"many2many:user_roles"`
}

func (User) TableName() string { return "users" }

// GetPermissions aggregates all permissions from all assigned roles.
// Deduplicates by permission key (resource:action).
func (u *User) GetPermissions() []string {
	seen := make(map[string]bool)
	var result []string

	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			key := perm.Resource + ":" + perm.Action
			if !seen[key] {
				seen[key] = true
				result = append(result, key)
			}
		}
	}
	return result
}

// GetRoleNames returns all role names as string slice
func (u *User) GetRoleNames() []string {
	names := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		names = append(names, r.Name)
	}
	return names
}

// Role represents a named set of permissions
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Description string
	IsSystem    bool `gorm:"default:false"` // system roles cannot be deleted
	IsActive    bool `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

func (Role) TableName() string { return "roles" }

// Permission represents a granular access control entry
type Permission struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Resource       string    `gorm:"not null;index"` // e.g. "supplier", "material"
	Action         string    `gorm:"not null"`       // e.g. "read", "create", "update"
	EndpointPath   string    `gorm:"not null"`       // e.g. "/suppliers/:id"
	EndpointMethod string    `gorm:"not null"`       // e.g. "GET", "POST"
	Description    string
	Hide           bool `gorm:"default:false"` // hidden from UI permission picker
	CreatedAt      time.Time

	// Composite unique constraint handled in migration
}

func (Permission) TableName() string { return "permissions" }

// Key returns the canonical permission identifier "resource:action"
func (p *Permission) Key() string {
	return p.Resource + ":" + p.Action
}

// ═══════════════════════════════════════════════════════════════
// Junction Tables (explicit for documentation, GORM handles these)
// ═══════════════════════════════════════════════════════════════

// UserRole maps users to roles (many-to-many)
type UserRole struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
}

func (UserRole) TableName() string { return "user_roles" }

// RolePermission maps roles to permissions (many-to-many)
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time
}

func (RolePermission) TableName() string { return "role_permissions" }
