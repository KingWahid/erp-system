package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../mocks/repository_mock.go

// UserRepository defines the contract for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}
