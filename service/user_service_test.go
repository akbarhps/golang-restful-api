package service

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-api/domain"
	"go-api/model"
	"go-api/repository"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func generatePassword() string {
	encrypt, _ := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	return string(encrypt)
}

var responseDomain = domain.User{
	Id:        uuid.New(),
	FullName:  "test controller",
	Username:  "testctrl",
	Email:     "testctrl@test.com",
	Password:  generatePassword(),
	CreatedAt: time.Now(),
}

var rawDomain = domain.User{
	Id:        uuid.New(),
	FullName:  "test controller",
	Username:  "testctrl",
	Email:     "testctrl@test.com",
	Password:  "testctrl",
	CreatedAt: time.Now(),
}

func setup() (UserService, *repository.UserRepositoryMock) {
	userRepository := repository.UserRepositoryMock{Mock: mock.Mock{}}
	userService := NewUserService(validator.New(), &userRepository)

	return userService, &userRepository
}

func TestUserServiceImpl_Register(t *testing.T) {
	t.Run("register using valid input should not get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{})

		userRepository.Mock.On("Create", context.Background(), mock.MatchedBy(func(user *domain.User) bool {
			assert.Equal(t, rawDomain.Username, user.Username)
			assert.Equal(t, rawDomain.Email, user.Email)
			assert.Equal(t, rawDomain.FullName, user.FullName)
			return true
		}))

		userRepository.Mock.On("Create", mock.Anything)

		assert.NotPanics(t, func() {
			resp := userService.Register(context.Background(), &model.UserRegister{
				FullName: rawDomain.FullName,
				Username: rawDomain.Username,
				Email:    rawDomain.Email,
				Password: rawDomain.Password,
			})
			assert.NotNil(t, resp)
		})
	})

	t.Run("register using bad input should get panic", func(t *testing.T) {
		userService, _ := setup()

		assert.Panics(t, func() {
			res := userService.Register(context.Background(), &model.UserRegister{
				FullName: "",
				Username: "",
				Email:    "",
				Password: "",
			})
			assert.Nil(t, res)
		})
	})

	t.Run("register using exist email or username should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{rawDomain})

		assert.PanicsWithError(t, "Username or Email already taken", func() {
			res := userService.Register(context.Background(), &model.UserRegister{
				FullName: rawDomain.FullName,
				Username: rawDomain.Username,
				Email:    rawDomain.Email,
				Password: rawDomain.Password,
			})
			assert.Nil(t, res)
		})
	})
}

func TestUserServiceImpl_Login(t *testing.T) {
	t.Run("login with valid input should not get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Username: responseDomain.Username,
		}).Return([]domain.User{responseDomain})

		assert.NotPanics(t, func() {
			res := userService.Login(context.Background(), &model.UserLogin{
				Username: rawDomain.Username,
				Password: rawDomain.Password,
			})
			assert.NotNil(t, res)
		})
	})

	t.Run("login with bad input should get panic", func(t *testing.T) {
		userService, _ := setup()

		assert.Panics(t, func() {
			res := userService.Login(context.Background(), &model.UserLogin{
				Username: "",
				Password: "",
			})
			assert.Nil(t, res)
		})
	})

	t.Run("login with non-exist account should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Username: rawDomain.Username,
		}).Return([]domain.User{})

		assert.PanicsWithError(t, "Record not found", func() {
			res := userService.Login(context.Background(), &model.UserLogin{
				Username: rawDomain.Username,
				Password: rawDomain.Password,
			})
			assert.Nil(t, res)
		})
	})

	t.Run("login with invalid password should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Username: rawDomain.Username,
		}).Return([]domain.User{rawDomain})

		assert.PanicsWithError(t, "Invalid password", func() {
			res := userService.Login(context.Background(), &model.UserLogin{
				Username: rawDomain.Username,
				Password: rawDomain.Password,
			})
			assert.Nil(t, res)
		})
	})
}

func TestUserServiceImpl_UpdateProfile(t *testing.T) {
	t.Run("update profile with valid input should not get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{responseDomain})
		userRepository.Mock.On("Update", context.Background(), &responseDomain)

		assert.NotPanics(t, func() {
			res := userService.UpdateProfile(context.Background(), &model.UserUpdateProfile{
				Id:       rawDomain.Id,
				Email:    rawDomain.Email,
				FullName: rawDomain.FullName,
				Username: rawDomain.Username,
			})
			assert.NotNil(t, res)
		})
	})

	t.Run("update profile with bad input should get panic", func(t *testing.T) {
		userService, _ := setup()

		assert.Panics(t, func() {
			res := userService.UpdateProfile(context.Background(), &model.UserUpdateProfile{
				Id:       uuid.New(),
				FullName: "",
				Username: "",
				Email:    "",
			})
			assert.Nil(t, res)
		})
	})

	t.Run("update profile with non-exist user should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{})

		assert.PanicsWithError(t, "Record not found", func() {
			res := userService.UpdateProfile(context.Background(), &model.UserUpdateProfile{
				Id:       rawDomain.Id,
				FullName: rawDomain.FullName,
				Username: rawDomain.Username,
				Email:    rawDomain.Email,
			})
			assert.Nil(t, res)
		})
	})
}

func TestUserServiceImpl_UpdatePassword(t *testing.T) {
	t.Run("update password with valid input should not get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{responseDomain})

		userRepository.Mock.On("Update", context.Background(), mock.MatchedBy(func(u *domain.User) bool {
			assert.Equal(t, responseDomain.Id, u.Id)
			assert.Equal(t, responseDomain.Email, u.Email)
			assert.Equal(t, responseDomain.Username, u.Username)
			return true
		}))

		userRepository.Mock.On("Update", mock.Anything)

		assert.NotPanics(t, func() {
			userService.UpdatePassword(context.Background(), &model.UserUpdatePassword{
				Id:          rawDomain.Id,
				Email:       rawDomain.Email,
				Username:    rawDomain.Username,
				OldPassword: rawDomain.Password,
				NewPassword: rawDomain.Password,
			})
		})
	})

	t.Run("update password with bad input should get panic", func(t *testing.T) {
		assert.Panics(t, func() {
			userService, _ := setup()

			userService.UpdatePassword(context.Background(), &model.UserUpdatePassword{
				Id:          uuid.New(),
				OldPassword: "",
				NewPassword: "",
			})
		})
	})

	t.Run("update password non-exist should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{})

		assert.PanicsWithError(t, "Record not found", func() {
			userService.UpdatePassword(context.Background(), &model.UserUpdatePassword{
				Id:          rawDomain.Id,
				Email:       rawDomain.Email,
				Username:    rawDomain.Username,
				OldPassword: rawDomain.Password,
				NewPassword: rawDomain.Password,
			})
		})
	})

	t.Run("update password with wrong old password should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{rawDomain})

		assert.PanicsWithError(t, "Invalid old password", func() {
			userService.UpdatePassword(context.Background(), &model.UserUpdatePassword{
				Id:          rawDomain.Id,
				Username:    rawDomain.Username,
				Email:       rawDomain.Email,
				OldPassword: "wrongpswd",
				NewPassword: rawDomain.Password,
			})
		})
	})
}

func TestUserServiceImpl_Delete(t *testing.T) {
	t.Run("delete with valid input should not get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{responseDomain})
		userRepository.Mock.On("Delete", context.Background(), &responseDomain)

		assert.NotPanics(t, func() {
			userService.Delete(context.Background(), &model.UserDelete{
				Id:       rawDomain.Id,
				Email:    rawDomain.Email,
				Password: rawDomain.Password,
				Username: rawDomain.Username,
			})
		})
	})

	t.Run("delete with bad input should get panic", func(t *testing.T) {
		assert.Panics(t, func() {
			userService, _ := setup()

			userService.Delete(context.Background(), &model.UserDelete{
				Id: uuid.New(),
			})
		})
	})

	t.Run("delete with non-exist user should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{})

		assert.PanicsWithError(t, "Record not found", func() {
			userService.Delete(context.Background(), &model.UserDelete{
				Id:       rawDomain.Id,
				Username: rawDomain.Username,
				Email:    rawDomain.Email,
				Password: rawDomain.Password,
			})
		})
	})

	t.Run("delete with wrong old password should get panic", func(t *testing.T) {
		userService, userRepository := setup()

		userRepository.Mock.On("Find", context.Background(), &domain.User{
			Id:       rawDomain.Id,
			Email:    rawDomain.Email,
			Username: rawDomain.Username,
		}).Return([]domain.User{rawDomain})

		assert.PanicsWithError(t, "Invalid password", func() {
			userService.Delete(context.Background(), &model.UserDelete{
				Id:       rawDomain.Id,
				Username: rawDomain.Username,
				Email:    rawDomain.Email,
				Password: "wrongpswd",
			})
		})
	})
}
