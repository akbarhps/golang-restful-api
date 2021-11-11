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
	repository := NewUserRepositoryImpl(db)

	err := repository.DeleteAll(context.Background())
	if err != nil {
		panic(err)
	}

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

func TestRepositoryImpl_CreateSuccess(t *testing.T) {
	repository := setup()
	user := createUser()

	err := repository.Create(context.Background(), user)
	assert.NoError(t, err)
}

func TestRepositoryImpl_CreateDuplicate(t *testing.T) {
	repository := setup()
	user := createUser()

	err := repository.Create(context.Background(), user)
	assert.NoError(t, err)

	err = repository.Create(context.Background(), user)
	assert.Error(t, err)
}

func TestRepositoryImpl_FindByUsername(t *testing.T) {
	repository := setup()
	user := createUser()

	err := repository.Create(context.Background(), user)
	assert.NoError(t, err)

	findUser, err := repository.Find(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEmpty(t, findUser)
	assert.Equal(t, user.Id, findUser[0].Id)
}

func TestRepositoryImpl_Update(t *testing.T) {
	repository := setup()
	user := createUser()

	err := repository.Create(context.Background(), user)
	assert.NoError(t, err)

	user.Username = "testupdate"

	err = repository.Update(context.Background(), user)
	assert.NoError(t, err)
}

func TestRepositoryImpl_Delete(t *testing.T) {
	repository := setup()
	user := createUser()

	err := repository.Create(context.Background(), user)
	assert.NoError(t, err)

	err = repository.Delete(context.Background(), user)
	assert.NoError(t, err)

	findUser, err := repository.Find(context.Background(), user)
	assert.Empty(t, findUser)
	assert.NoError(t, err)
}
