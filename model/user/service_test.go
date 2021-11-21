package user_test

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-api/app"
	"go-api/model/user"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func setupServiceTest() (*user.RepositoryMock, user.Service) {
	app.TestDBInit()
	repository := &user.RepositoryMock{mock.Mock{}}
	service := user.NewService(validator.New(), repository)
	return repository, service
}

func TestServiceImpl_Register(t *testing.T) {
	var (
		registerValid = &user.RegisterRequest{
			Email:       "testservice@test.com",
			Username:    "testservice",
			DisplayName: "test service",
			Password:    "test service",
		}

		registerEmpty = &user.RegisterRequest{
			Email:       "",
			Username:    "",
			DisplayName: "",
			Password:    "",
		}
	)

	t.Run("success should return auth response with no error", func(t *testing.T) {
		repository, service := setupServiceTest()
		repository.On("FindByUsername", registerValid.Username).Return(&user.User{})
		repository.On("FindByEmail", registerValid.Email).Return(&user.User{})
		repository.On("Create", mock.MatchedBy(func(u *user.User) bool {
			assert.Equal(t, registerValid.Email, u.Email)
			assert.Equal(t, registerValid.Username, u.Username)
			assert.Equal(t, registerValid.DisplayName, u.DisplayName)
			return true
		})).Return(nil)

		assert.NotPanics(t, func() {
			res := service.Register(context.Background(), registerValid)
			t.Log(res)
			assert.NotEmpty(t, res)
		})
	})

	t.Run("empty input should panic", func(t *testing.T) {
		_, service := setupServiceTest()

		assert.Panics(t, func() {
			res := service.Register(context.Background(), registerEmpty)
			assert.Empty(t, res)
			t.Log(res)
		})
	})

	t.Run("taken username or email should panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		repository.On("FindByUserID", registerValid.Username).Return(&user.User{}, nil)
		repository.On("FindByEmail", registerValid.Email).Return(&user.User{}, nil)

		assert.Panics(t, func() {
			res := service.Register(context.Background(), registerEmpty)
			assert.Empty(t, res)
			t.Log(res)
		})
	})
}

func TestServiceImpl_Login(t *testing.T) {
	requestValid := &user.LoginRequest{
		Handler:  "testservice",
		Password: "testservice",
	}
	requestEmpty := &user.LoginRequest{
		Handler:  "",
		Password: "",
	}

	t.Run("success should return auth response object", func(t *testing.T) {
		repository, service := setupServiceTest()
		enc, _ := bcrypt.GenerateFromPassword([]byte(requestValid.Password), bcrypt.DefaultCost)
		repository.On("FindByEmailOrUsername", requestValid.Handler).Return(&user.User{
			ID:          uuid.NewV4().String(),
			Username:    requestValid.Handler,
			DisplayName: "test service",
			Password:    string(enc),
		}, nil)

		assert.NotPanics(t, func() {
			res := service.Login(context.Background(), requestValid)
			assert.NotEmpty(t, res)
			t.Log(res)
		})
	})

	t.Run("empty input should panic", func(t *testing.T) {
		_, service := setupServiceTest()
		assert.Panics(t, func() {
			res := service.Login(context.Background(), requestEmpty)
			assert.NotEmpty(t, res)
			t.Log(res)
		})
	})

	t.Run("wrong password should panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		enc, _ := bcrypt.GenerateFromPassword([]byte(requestValid.Password), bcrypt.DefaultCost)
		repository.On("FindByEmailOrUsername", requestValid.Handler).Return(&user.User{
			ID:          uuid.NewV4().String(),
			Username:    requestValid.Handler,
			DisplayName: "test service",
			Password:    string(enc),
		}, nil)

		assert.Panics(t, func() {
			res := service.Login(context.Background(), &user.LoginRequest{
				Handler:  requestValid.Handler,
				Password: "wrong password",
			})
			assert.Empty(t, res)
			t.Log(res)
		})
	})
}

func TestServiceImpl_UpdateProfile(t *testing.T) {
	requestValid := &user.UpdateProfileRequest{
		UserID:      uuid.NewV4().String(),
		Email:       "testservice@test.com",
		Username:    "testservice",
		DisplayName: "test service",
		Biography:   "test service",
	}
	requestEmpty := &user.UpdateProfileRequest{
		UserID:      "",
		Email:       "",
		Username:    "",
		DisplayName: "",
		Biography:   "",
	}

	t.Run("success should not panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		repository.On("FindById", requestValid.UserID).Return(&user.User{
			ID:          requestValid.UserID,
			Username:    requestValid.Username,
			DisplayName: requestValid.DisplayName,
			Email:       requestValid.Email,
			Biography:   requestValid.Biography,
		})
		repository.On("FindByUsername", requestValid.Username).Return(&user.User{})
		repository.On("FindByEmail", requestValid.Email).Return(&user.User{})
		repository.On("Update", mock.MatchedBy(func(u *user.User) bool {
			assert.Equal(t, requestValid.Email, u.Email)
			assert.Equal(t, requestValid.Username, u.Username)
			assert.Equal(t, requestValid.DisplayName, u.DisplayName)
			assert.Equal(t, requestValid.Biography, u.Biography)
			return true
		})).Return(nil)

		assert.NotPanics(t, func() {
			service.UpdateProfile(context.Background(), requestValid)
		})
	})

	t.Run("empty input should panic", func(t *testing.T) {
		_, service := setupServiceTest()

		assert.Panics(t, func() {
			service.UpdateProfile(context.Background(), requestEmpty)
		})
	})
}

func TestServiceImpl_UpdatePassword(t *testing.T) {
	var (
		updatePasswordValid = &user.UpdatePasswordRequest{
			UserID:      uuid.NewV4().String(),
			OldPassword: "testservice",
			NewPassword: "testservice123",
		}

		updatePasswordEmpty = &user.UpdatePasswordRequest{
			UserID:      "",
			OldPassword: "",
			NewPassword: "",
		}

		updatePasswordWrongPassword = &user.UpdatePasswordRequest{
			UserID:      updatePasswordValid.UserID,
			OldPassword: "wrongpassword",
			NewPassword: "testservice123",
		}
	)
	t.Run("success should not panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		enc, _ := bcrypt.GenerateFromPassword([]byte(updatePasswordValid.OldPassword), bcrypt.DefaultCost)
		mUser := &user.User{
			ID:          updatePasswordValid.UserID,
			Email:       "testservice@test.com",
			Username:    "testservice",
			DisplayName: "test service",
			Password:    string(enc),
		}

		repository.On("FindById", updatePasswordValid.UserID).Return(mUser)
		repository.On("Update", mock.MatchedBy(func(u *user.User) bool {
			assert.Equal(t, mUser.ID, u.ID)
			return true
		}))

		assert.NotPanics(t, func() {
			service.UpdatePassword(context.Background(), updatePasswordValid)
		})
	})

	t.Run("empty input should panic", func(t *testing.T) {
		_, service := setupServiceTest()
		assert.Panics(t, func() {
			service.UpdatePassword(context.Background(), updatePasswordEmpty)
		})
	})

	t.Run("not found user should panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		repository.On("FindById", updatePasswordValid.UserID).Return(nil)
		assert.Panics(t, func() {
			service.UpdatePassword(context.Background(), updatePasswordValid)
		})
	})

	t.Run("wrong old password should panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		enc, _ := bcrypt.GenerateFromPassword([]byte(updatePasswordValid.OldPassword), bcrypt.DefaultCost)
		mUser := &user.User{
			ID:          updatePasswordValid.UserID,
			Email:       "testservice@test.com",
			Username:    "testservice",
			DisplayName: "test service",
			Password:    string(enc),
		}

		repository.On("FindById", updatePasswordValid.UserID).Return(mUser)
		assert.Panics(t, func() {
			service.UpdatePassword(context.Background(), updatePasswordWrongPassword)
		})
	})
}

func TestServiceImpl_FindByUsername(t *testing.T) {
	t.Run("success should return user detail response", func(t *testing.T) {
		repository, service := setupServiceTest()
		userID := "testservice"

		repository.On("FindByUsername", userID).Return(&user.User{
			ID:          userID,
			Email:       "testservice@test.com",
			Username:    "testservice",
			DisplayName: "test service",
		})

		assert.NotPanics(t, func() {
			res := service.FindByUsername(context.Background(), userID)
			assert.NotNil(t, res)
		})
	})

	t.Run("not found user should panic", func(t *testing.T) {
		repository, service := setupServiceTest()
		userID := "testservice"

		repository.On("FindByUsername", userID).Return(nil)

		assert.Panics(t, func() {
			res := service.FindByUsername(context.Background(), userID)
			assert.Nil(t, res)
		})
	})
}

func TestServiceImpl_SearchLike(t *testing.T) {
	t.Run("success should return slice of user", func(t *testing.T) {
		repository, service := setupServiceTest()
		keyword := "testservice"

		repository.On("FindLike", keyword).Return([]*user.User{
			{
				Username:    "user1",
				DisplayName: "user1",
			},
			{
				Username:    "user2",
				DisplayName: "user2",
			},
			{
				Username:    "user3",
				DisplayName: "user3",
			},
		})

		assert.NotPanics(t, func() {
			res := service.SearchLike(context.Background(), keyword)
			assert.NotEmpty(t, res)
		})
	})

	t.Run("not found user should return empty slice", func(t *testing.T) {
		repository, service := setupServiceTest()
		keyword := "notfounduser"

		repository.On("FindLike", keyword).Return([]*user.User{})

		assert.NotPanics(t, func() {
			res := service.SearchLike(context.Background(), keyword)
			assert.Empty(t, res)
		})
	})
}
