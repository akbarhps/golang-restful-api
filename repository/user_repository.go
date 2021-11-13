package repository

import (
	"context"
	"database/sql"
	"go-api/domain"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *domain.User)
	Find(ctx context.Context, tx *sql.Tx, user *domain.User) []domain.User
	Update(ctx context.Context, tx *sql.Tx, user *domain.User)
	Delete(ctx context.Context, tx *sql.Tx, user *domain.User)
	DeleteAll(ctx context.Context, tx *sql.Tx)
}
