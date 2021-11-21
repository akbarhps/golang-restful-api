package user

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(tx *gorm.DB, user *Entity) {
	r.Called(user)
}

func (r *RepositoryMock) Update(tx *gorm.DB, user *Entity) {
	r.Called(user)
}

func (r *RepositoryMock) Delete(tx *gorm.DB, user *Entity) {
	r.Called(user)
}

func (r *RepositoryMock) FindById(tx *gorm.DB, id string) *Entity {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*Entity)
	}
	return nil
}

func (r *RepositoryMock) FindLike(tx *gorm.DB, keyword string) []*Entity {
	args := r.Called(keyword)
	if args.Get(0) != nil {
		return args.Get(0).([]*Entity)
	}
	return nil
}

func (r *RepositoryMock) FindByEmail(tx *gorm.DB, email string) *Entity {
	args := r.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*Entity)
	}
	return nil
}

func (r *RepositoryMock) FindByUsername(tx *gorm.DB, username string) *Entity {
	args := r.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*Entity)
	}
	return nil
}

func (r *RepositoryMock) FindByEmailOrUsername(tx *gorm.DB, handler string) *Entity {
	args := r.Called(handler)
	if args.Get(0) != nil {
		return args.Get(0).(*Entity)
	}
	return nil
}
