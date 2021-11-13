package repository

import (
	"context"
	"database/sql"
	"go-api/app"
	"go-api/domain"
	"go-api/helper"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var db = app.NewDatabase("test")

func setup() (UserRepository, *sql.Tx) {
	repository := NewUserRepository()

	tx, _ := db.Begin()
	repository.DeleteAll(context.Background(), tx)

	return repository, tx
}

func createUser() *domain.User {
	uid, err := uuid.NewUUID()
	helper.PanicIfError(err)

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
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), tx, user)
		})
	})

	t.Run("create user with registered email or username should panic", func(t *testing.T) {
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), tx, user)
		})
		assert.Panics(t, func() {
			repository.Create(context.Background(), tx, user)
		})
	})
}

func TestUserRepositoryImpl_Find(t *testing.T) {
	t.Run("find success should return users", func(t *testing.T) {
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), tx, user)
			users := repository.Find(context.Background(), tx, user)
			assert.NotEmpty(t, users)
		})
	})

	t.Run("find user not found should return empty slice", func(t *testing.T) {
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)

		assert.NotPanics(t, func() {
			users := repository.Find(context.Background(), tx, &domain.User{})
			assert.Empty(t, users)
		})
	})
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	t.Run("update success should not panic", func(t *testing.T) {
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), tx, user)
			user.FullName = "test repo update"
			repository.Update(context.Background(), tx, user)
		})
	})

	t.Run("update non-exist user should not panic", func(t *testing.T) {
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Update(context.Background(), tx, user)
		})
	})
}

func TestUserRepositoryImpl_Delete(t *testing.T) {
	t.Run("delete user should return empty slice when try to find it", func(t *testing.T) {
		repository, tx := setup()
		defer helper.TXCommitOrRollback(tx)
		user := createUser()

		assert.NotPanics(t, func() {
			repository.Create(context.Background(), tx, user)
			repository.Delete(context.Background(), tx, user)
			users := repository.Find(context.Background(), tx, user)
			assert.Empty(t, users)
		})
	})
}
