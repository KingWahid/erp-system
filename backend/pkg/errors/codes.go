package errors

// ErrorCode represents application-specific error codes
type ErrorCode string

// Error codes for different error scenarios
const (
	ErrCodeDatabaseConnection  ErrorCode = "ERR_DATABASE_CONNECTION"
	ErrCodeDatabaseQuery       ErrorCode = "ERR_DATABASE_QUERY"
	ErrCodeDatabaseTransaction ErrorCode = "ERR_DATABASE_TRANSACTION"
	ErrCodeRecordNotFound      ErrorCode = "ERR_RECORD_NOT_FOUND"
	ErrCodeDuplicateKey        ErrorCode = "ERR_DUPLICATE_KEY"

	ErrCodeRedisConnection ErrorCode = "ERR_REDIS_CONNECTION"
	ErrCodeCacheOperation  ErrorCode = "ERR_CACHE_OPERATION"
	ErrCodeCacheMiss       ErrorCode = "ERR_CACHE_MISS"

	ErrCodeInvalidInput     ErrorCode = "ERR_INVALID_INPUT"
	ErrCodeValidationFailed ErrorCode = "ERR_VALIDATION_FAILED"
	ErrCodeInvalidUUID      ErrorCode = "ERR_INVALID_UUID"
	ErrCodeInvalidEmail     ErrorCode = "ERR_INVALID_EMAIL"

	ErrCodeUnauthorized       ErrorCode = "ERR_UNAUTHORIZED"
	ErrCodeInvalidCredentials ErrorCode = "ERR_INVALID_CREDENTIALS"
	ErrCodeInvalidToken       ErrorCode = "ERR_INVALID_TOKEN"
	ErrCodeExpiredToken       ErrorCode = "ERR_EXPIRED_TOKEN"
	ErrCodeTokenGeneration    ErrorCode = "ERR_TOKEN_GENERATION"

	ErrCodeForbidden              ErrorCode = "ERR_FORBIDDEN"
	ErrCodeInsufficientPermission ErrorCode = "ERR_INSUFFICIENT_PERMISSION"
	ErrCodeInactiveAccount        ErrorCode = "ERR_INACTIVE_ACCOUNT"

	ErrCodeResourceNotFound   ErrorCode = "ERR_RESOURCE_NOT_FOUND"
	ErrCodeResourceConflict   ErrorCode = "ERR_RESOURCE_CONFLICT"
	ErrCodeBadRequest         ErrorCode = "ERR_BAD_REQUEST"
	ErrCodeInternalServer     ErrorCode = "ERR_INTERNAL_SERVER"
	ErrCodeServiceUnavailable ErrorCode = "ERR_SERVICE_UNAVAILABLE"

	ErrCodeSupplierNotFound       ErrorCode = "ERR_SUPPLIER_NOT_FOUND"
	ErrCodeSupplierCodeExists     ErrorCode = "ERR_SUPPLIER_CODE_EXISTS"
	ErrCodeSupplierInvalidStage   ErrorCode = "ERR_SUPPLIER_INVALID_STAGE"
	ErrCodeSupplierAlreadyBlocked ErrorCode = "ERR_SUPPLIER_ALREADY_BLOCKED"

	ErrCodeUserNotFound    ErrorCode = "ERR_USER_NOT_FOUND"
	ErrCodeEmailExists     ErrorCode = "ERR_EMAIL_EXISTS"
	ErrCodePasswordWeak    ErrorCode = "ERR_PASSWORD_WEAK"
	ErrCodePasswordHashing ErrorCode = "ERR_PASSWORD_HASHING"

	ErrCodeRoleNotFound       ErrorCode = "ERR_ROLE_NOT_FOUND"
	ErrCodePermissionNotFound ErrorCode = "ERR_PERMISSION_NOT_FOUND"
	ErrCodeRoleAssignFailed   ErrorCode = "ERR_ROLE_ASSIGN_FAILED"
)
