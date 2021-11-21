package user_test

//
//import (
//	"context"
//	"github.com/go-playground/validator"
//	"github.com/satori/go.uuid"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"go-api/app"
//	"go-api/model/user"
//	"golang.org/x/crypto/bcrypt"
//	"testing"
//)
//
//func setupServiceTest() (user.Repository, user.Service) {
//	app.TestDBInit()
//
//	repository := &user.RepositoryMock{mock.Mock{}}
//	service := user.NewService(validator.New(), repository)
//	return repository, service
//}
//
//func TestServiceImpl_Register(t *testing.T) {
//	requestValid := &RegisterRequest{
//		Email:       "testservice@test.com",
//		Username:    "testservice",
//		DisplayName: "test service",
//		Password:    "test service",
//	}
//	requestEmpty := &RegisterRequest{
//		Email:       "",
//		Username:    "",
//		DisplayName: "",
//		Password:    "",
//	}
//
//	t.Run("success should return auth response with no error", func(t *testing.T) {
//		uRepository, uService := setupServiceTest()
//		uRepository.On("FindByUserID", requestValid.Username).Return(nil, nil)
//		uRepository.On("FindByEmail", requestValid.Email).Return(nil, nil)
//		uRepository.On("Create", mock.MatchedBy(func(u *User) bool {
//			assert.Equal(t, requestValid.Email, u.Email)
//			assert.Equal(t, requestValid.Username, u.Username)
//			assert.Equal(t, requestValid.DisplayName, u.DisplayName)
//			return true
//		})).Return(nil)
//
//		assert.NotPanics(t, func() {
//			res := uService.Register(context.Background(), requestValid)
//			t.Log(res)
//			assert.NotEmpty(t, res)
//		})
//	})
//
//	t.Run("empty input should panic", func(t *testing.T) {
//		_, uService := setupServiceTest()
//
//		assert.Panics(t, func() {
//			res := uService.Register(context.Background(), requestEmpty)
//			assert.Empty(t, res)
//			t.Log(res)
//		})
//	})
//
//	t.Run("taken username or email should panic", func(t *testing.T) {
//		uRepository, uService := setupServiceTest()
//		uRepository.On("FindByUserID", requestValid.Username).Return(&User{}, nil)
//		uRepository.On("FindByEmail", requestValid.Email).Return(&User{}, nil)
//
//		assert.Panics(t, func() {
//			res := uService.Register(context.Background(), requestEmpty)
//			assert.Empty(t, res)
//			t.Log(res)
//		})
//	})
//}
//
//func TestServiceImpl_Login(t *testing.T) {
//	requestValid := &LoginRequest{
//		Handler:  "testservice",
//		Password: "testservice",
//	}
//	requestEmpty := &LoginRequest{
//		Handler:  "",
//		Password: "",
//	}
//
//	t.Run("success should return auth response object", func(t *testing.T) {
//		uRepository, uService := setupServiceTest()
//		enc, _ := bcrypt.GenerateFromPassword([]byte(requestValid.Password), bcrypt.DefaultCost)
//		uRepository.On("FindByEmailOrUsername", requestValid.Handler).Return(&User{
//			ID:          uuid.NewV4().String(),
//			Username:    requestValid.Handler,
//			DisplayName: "test service",
//			Password:    string(enc),
//		}, nil)
//
//		assert.NotPanics(t, func() {
//			res := uService.Login(context.Background(), requestValid)
//			assert.NotEmpty(t, res)
//			t.Log(res)
//		})
//	})
//
//	t.Run("empty input should panic", func(t *testing.T) {
//		_, uService := setupServiceTest()
//		assert.Panics(t, func() {
//			res := uService.Login(context.Background(), requestEmpty)
//			assert.NotEmpty(t, res)
//			t.Log(res)
//		})
//	})
//
//	t.Run("wrong password should panic", func(t *testing.T) {
//		uRepository, uService := setupServiceTest()
//		enc, _ := bcrypt.GenerateFromPassword([]byte(requestValid.Password), bcrypt.DefaultCost)
//		uRepository.On("FindByEmailOrUsername", requestValid.Handler).Return(&User{
//			ID:          uuid.NewV4().String(),
//			Username:    requestValid.Handler,
//			DisplayName: "test service",
//			Password:    string(enc),
//		}, nil)
//
//		assert.Panics(t, func() {
//			res := uService.Login(context.Background(), &LoginRequest{
//				Handler:  requestValid.Handler,
//				Password: "wrong password",
//			})
//			assert.Empty(t, res)
//			t.Log(res)
//		})
//	})
//}
//
//func TestServiceImpl_UpdateProfile(t *testing.T) {
//	requestValid := &UpdateProfileRequest{
//		UserID:      uuid.NewV4().String(),
//		Email:       "testservice@test.com",
//		Username:    "testservice",
//		DisplayName: "test service",
//		Biography:   "test service",
//	}
//	requestEmpty := &UpdateProfileRequest{
//		UserID:      "",
//		Email:       "",
//		Username:    "",
//		DisplayName: "",
//		Biography:   "",
//	}
//
//	t.Run("success should not panic", func(t *testing.T) {
//		uRepository, uService := setupServiceTest()
//		uRepository.On("FindById", requestValid.UserID).Return(&User{
//			ID:          requestValid.UserID,
//			Username:    requestValid.Username,
//			DisplayName: requestValid.DisplayName,
//			Email:       requestValid.Email,
//			Biography:   requestValid.Biography,
//		})
//		uRepository.On("FindByUserID", requestValid.Username).Return(nil)
//		uRepository.On("FindByEmail", requestValid.Email).Return(nil)
//		uRepository.On("Update", mock.MatchedBy(func(u *User) bool {
//			assert.Equal(t, requestValid.Email, u.Email)
//			assert.Equal(t, requestValid.Username, u.Username)
//			assert.Equal(t, requestValid.DisplayName, u.DisplayName)
//			assert.Equal(t, requestValid.Biography, u.Biography)
//			return true
//		})).Return(nil)
//
//		assert.NotPanics(t, func() {
//			uService.UpdateProfile(context.Background(), requestValid)
//		})
//	})
//
//	t.Run("empty input should panic", func(t *testing.T) {
//		_, uService := setupServiceTest()
//
//		assert.Panics(t, func() {
//			uService.UpdateProfile(context.Background(), requestEmpty)
//		})
//	})
//}
//
//func TestServiceImpl_UpdatePassword(t *testing.T) {
//
//}
//
//func TestServiceImpl_FindByUsername(t *testing.T) {
//
//}
//
//func TestServiceImpl_SearchLike(t *testing.T) {
//
//}
