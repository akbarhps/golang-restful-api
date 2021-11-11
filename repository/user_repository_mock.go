package repository

import (
	"context"
	"go-api/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	args := repository.Mock.Called(ctx, user)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return nil
}

func (repository *UserRepositoryMock) Find(ctx context.Context, user *domain.User) ([]domain.User, error) {
	args := repository.Mock.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), nil
}

func (repository *UserRepositoryMock) Update(ctx context.Context, user *domain.User) error {
	args := repository.Mock.Called(ctx, user)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return nil
}

func (repository *UserRepositoryMock) Delete(ctx context.Context, user *domain.User) error {
	args := repository.Mock.Called(ctx, user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (repository *UserRepositoryMock) DeleteAll(ctx context.Context) error {
	args := repository.Mock.Called(ctx)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Get(1).(error)
}
