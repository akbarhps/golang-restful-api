package repository

import (
	"context"
	"go-api/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User)
	Find(ctx context.Context, user *domain.User) []domain.User
	Update(ctx context.Context, user *domain.User)
	Delete(ctx context.Context, user *domain.User)
	DeleteAll(ctx context.Context)
}
