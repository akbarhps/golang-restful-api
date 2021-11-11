package repository

import (
	"context"
	"go-api/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Find(ctx context.Context, user *domain.User) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, user *domain.User) error
	DeleteAll(ctx context.Context) error
}
