package errors

import (
	stderrors "errors"
	"strings"

	"gorm.io/gorm"
)

// IsDuplicateKeyError checks if error is a PostgreSQL unique constraint violation
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if stderrors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "23505")
}
