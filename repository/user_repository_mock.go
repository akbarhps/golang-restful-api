package repository

import (
	"context"
	"database/sql"
	"go-api/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) Create(ctx context.Context, tx *sql.Tx, user *domain.User) {
	repository.Mock.Called(ctx, user)
}

func (repository *UserRepositoryMock) Find(ctx context.Context, tx *sql.Tx, user *domain.User) []domain.User {
	args := repository.Mock.Called(ctx, user)
	return args.Get(0).([]domain.User)
}

func (repository *UserRepositoryMock) Update(ctx context.Context, tx *sql.Tx, user *domain.User) {
	repository.Mock.Called(ctx, user)
}

func (repository *UserRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, user *domain.User) {
	repository.Mock.Called(ctx, user)
}

func (repository *UserRepositoryMock) DeleteAll(ctx context.Context, tx *sql.Tx) {
	repository.Mock.Called(ctx)
}
