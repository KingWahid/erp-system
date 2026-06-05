package cache

import (
	"bytes"
	"text/template"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// CACHE KEY TEMPLATES
// ═══════════════════════════════════════════════════════════════

var (
	// Supplier keys
	supplierDetailTmpl       = template.Must(template.New("supplierDetail").Parse(`supplier:detail:{{.supplierID}}`))
	supplierListTmpl         = template.Must(template.New("supplierList").Parse(`supplier:list:p{{.page}}:l{{.limit}}:s{{.status}}:q{{.search}}`))
	supplierListPatternTmpl  = template.Must(template.New("supplierListPattern").Parse(`supplier:list:*`))
	supplierStatsTmpl        = template.Must(template.New("supplierStats").Parse(`supplier:stats`))
	supplierMaterialsTmpl    = template.Must(template.New("supplierMaterials").Parse(`supplier:{{.supplierID}}:materials`))
	supplierRatingsTmpl      = template.Must(template.New("supplierRatings").Parse(`supplier:{{.supplierID}}:ratings:p{{.page}}:l{{.limit}}`))
	supplierRatingsPatternTmpl = template.Must(template.New("supplierRatingsPattern").Parse(`supplier:{{.supplierID}}:ratings:*`))

	// User/Auth keys
	userPermissionsTmpl = template.Must(template.New("userPermissions").Parse(`user:{{.userID}}:permissions`))
	userDetailTmpl      = template.Must(template.New("userDetail").Parse(`user:detail:{{.userID}}`))
)

// ═══════════════════════════════════════════════════════════════
// CACHE TTL CONSTANTS
// ═══════════════════════════════════════════════════════════════

const (
	SupplierDetailCacheTTL     = 15 * time.Minute
	SupplierListCacheTTL       = 2 * time.Minute
	SupplierStatsCacheTTL      = 2 * time.Minute
	SupplierMaterialsCacheTTL  = 10 * time.Minute
	SupplierRatingsCacheTTL    = 5 * time.Minute
	UserPermissionsCacheTTL    = 1 * time.Hour
	UserDetailCacheTTL         = 10 * time.Minute
)

// ═══════════════════════════════════════════════════════════════
// SUPPLIER CACHE KEY BUILDERS
// ═══════════════════════════════════════════════════════════════

// BuildSupplierDetailCacheKey builds cache key for supplier detail.
// Pattern: "supplier:detail:{supplierID}".
func BuildSupplierDetailCacheKey(supplierID string) (string, error) {
	return executeTemplate(supplierDetailTmpl, map[string]any{
		"supplierID": supplierID,
	})
}

// BuildSupplierListCacheKey builds cache key for paginated supplier list.
// Pattern: "supplier:list:p{page}:l{limit}:s{status}:q{search}".
func BuildSupplierListCacheKey(page, limit int, status, search string) (string, error) {
	return executeTemplate(supplierListTmpl, map[string]any{
		"page":   page,
		"limit":  limit,
		"status": status,
		"search": search,
	})
}

// BuildSupplierListCachePattern builds invalidation pattern for all supplier lists.
// Pattern: "supplier:list:*".
func BuildSupplierListCachePattern() (string, error) {
	return executeTemplate(supplierListPatternTmpl, nil)
}

// BuildSupplierStatsCacheKey builds cache key for supplier stats.
// Pattern: "supplier:stats".
func BuildSupplierStatsCacheKey() (string, error) {
	return executeTemplate(supplierStatsTmpl, nil)
}

// BuildSupplierMaterialsCacheKey builds cache key for supplier materials.
// Pattern: "supplier:{supplierID}:materials".
func BuildSupplierMaterialsCacheKey(supplierID string) (string, error) {
	return executeTemplate(supplierMaterialsTmpl, map[string]any{
		"supplierID": supplierID,
	})
}

// BuildSupplierRatingsCacheKey builds cache key for supplier ratings list.
// Pattern: "supplier:{supplierID}:ratings:p{page}:l{limit}".
func BuildSupplierRatingsCacheKey(supplierID string, page, limit int) (string, error) {
	return executeTemplate(supplierRatingsTmpl, map[string]any{
		"supplierID": supplierID,
		"page":       page,
		"limit":      limit,
	})
}

// BuildSupplierRatingsCachePattern builds invalidation pattern for supplier ratings.
// Pattern: "supplier:{supplierID}:ratings:*".
func BuildSupplierRatingsCachePattern(supplierID string) (string, error) {
	return executeTemplate(supplierRatingsPatternTmpl, map[string]any{
		"supplierID": supplierID,
	})
}

// ═══════════════════════════════════════════════════════════════
// USER/AUTH CACHE KEY BUILDERS
// ═══════════════════════════════════════════════════════════════

// BuildUserPermissionsCacheKey builds cache key for user permissions.
// Pattern: "user:{userID}:permissions".
func BuildUserPermissionsCacheKey(userID string) (string, error) {
	return executeTemplate(userPermissionsTmpl, map[string]any{
		"userID": userID,
	})
}

// BuildUserDetailCacheKey builds cache key for user detail.
// Pattern: "user:detail:{userID}".
func BuildUserDetailCacheKey(userID string) (string, error) {
	return executeTemplate(userDetailTmpl, map[string]any{
		"userID": userID,
	})
}

// ═══════════════════════════════════════════════════════════════
// HELPER
// ═══════════════════════════════════════════════════════════════

// executeTemplate executes a template with given data and returns the result.
func executeTemplate(tmpl *template.Template, data map[string]any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
