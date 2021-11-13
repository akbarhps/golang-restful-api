package repository

import (
	"context"
	"go-api/app"
	"go-api/domain"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setup() UserRepository {
	db := app.NewDatabase("test")
	repository := NewUserRepository(db)
	repository.DeleteAll(context.Background())

	return repository
}

func createUser() *domain.User {
	uid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return &domain.User{
		Id:        uid,
		FullName:  "test repo",
		Username:  "testrepo",
		Email:     "testrepo@test.com",
		Password:  "testrepo",
		CreatedAt: time.Now(),
	}
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	t.Run("create user with valid input should not panic", func(t *testing.T) {
		repository := setup()
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), user)
		})
	})

	t.Run("create user with registered email or username should panic", func(t *testing.T) {
		repository := setup()
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), user)
		})
		assert.Panics(t, func() {
			repository.Create(context.Background(), user)
		})
	})
}

func TestUserRepositoryImpl_Find(t *testing.T) {
	t.Run("find success should return users", func(t *testing.T) {
		repository := setup()
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), user)
			users := repository.Find(context.Background(), user)
			assert.NotEmpty(t, users)
		})
	})

	t.Run("find user not found should return empty slice", func(t *testing.T) {
		repository := setup()

		assert.NotPanics(t, func() {
			users := repository.Find(context.Background(), &domain.User{})
			assert.Empty(t, users)
		})
	})
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	t.Run("update success should not panic", func(t *testing.T) {
		repository := setup()
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), user)
			user.FullName = "test repo update"
			repository.Update(context.Background(), user)
		})
	})

	t.Run("update non-exist user should not panic", func(t *testing.T) {
		repository := setup()
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Update(context.Background(), user)
		})
	})
}

func TestUserRepositoryImpl_Delete(t *testing.T) {
	t.Run("delete user should return empty slice when try to find it", func(t *testing.T) {
		repository := setup()
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), user)
			repository.Delete(context.Background(), user)
			users := repository.Find(context.Background(), user)
			assert.Empty(t, users)
		})
	})
}
