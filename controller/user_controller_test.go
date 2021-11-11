package controller

import (
	"context"
	"encoding/json"
	"go-api/app"
	"go-api/domain"
	"go-api/model"
	"go-api/repository"
	"go-api/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var r *gin.Engine
var userRepository repository.UserRepository
var userService service.UserService
var userController UserController

func TestMain(m *testing.M) {
	r = gin.Default()

	db := app.NewDatabase("test")
	userRepository = repository.NewUserRepositoryImpl(db)
	userService = service.NewUserService(validator.New(), userRepository)
	userController = NewUserController(userService)

	r.GET("/api/user/:key", userController.Find)
	r.POST("/api/user", userController.Register)
	r.POST("/api/user/login", userController.Login)
	r.PUT("/api/user/profile", userController.UpdateProfile)
	r.PUT("/api/user/password", userController.UpdatePassword)
	r.DELETE("/api/user", userController.Delete)

	m.Run()
}

func TestUserControllerImpl_RegisterSuccess(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	registerRequestJSON := "{" +
		`"full_name" : "test controller",` +
		`"username": "testctrl",` +
		`"email" : "testctrl@test.com",` +
		`"password" : "testctrl"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(registerRequestJSON))
	w := httptest.NewRecorder()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.NotEmpty(t, resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_RegisterBadFormat(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	registerRequestJSON := "{" +
		`"full_name" : "",` +
		`"username": "",` +
		`"email" : "testctrltest.com",` +
		`"password" : "tes"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(registerRequestJSON))
	w := httptest.NewRecorder()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_RegisterDuplicate(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	_, err := userService.Register(context.Background(), &model.UserRegisterRequest{
		FullName: "test controller",
		Username: "testctrl",
		Email:    "testctrl@test.com",
		Password: "testctrl",
	})
	assert.NoError(t, err)

	registerRequestJSON := "{" +
		`"full_name" : "test controller",` +
		`"username" : "testctrl",` +
		`"email" : "testctrl@test.com",` +
		`"password" : "testctrl"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(registerRequestJSON))
	w := httptest.NewRecorder()

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotEmpty(t, string(resBody))

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_LoginSuccess(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	_, err := userService.Register(context.Background(), &model.UserRegisterRequest{
		FullName: "test controller",
		Username: "testctrl",
		Email:    "testctrl@test.com",
		Password: "testctrl",
	})
	assert.NoError(t, err)

	loginRequestJSON := "{" +
		`"email" : "testctrl@test.com",` +
		`"password" : "testctrl"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(loginRequestJSON))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.NotEmpty(t, resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_LoginNotFound(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())

	loginRequestJSON := "{" +
		`"email" : "testctrl@test.com",` +
		`"password" : "testctrl"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(loginRequestJSON))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	assert.Equal(t, "Record not found", resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_LoginWrongPassword(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	_, err := userService.Register(context.Background(), &model.UserRegisterRequest{
		FullName: "test controller",
		Username: "testctrl",
		Email:    "testctrl@test.com",
		Password: "testctrl",
	})
	assert.NoError(t, err)

	loginRequestJSON := "{" +
		`"email" : "testctrl@test.com",` +
		`"password" : "testctrl2"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(loginRequestJSON))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_LoginBadFormat(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	_, err := userService.Register(context.Background(), &model.UserRegisterRequest{
		FullName: "test controller",
		Username: "testctrl",
		Email:    "testctrl@test.com",
		Password: "testctrl",
	})
	assert.NoError(t, err)

	loginRequestJSON := "{" +
		`"email" : "testctrltest.com",` +
		`"password" : "rl2"` +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(loginRequestJSON))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_FindSuccess(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	_, err := userService.Register(context.Background(), &model.UserRegisterRequest{
		FullName: "test controller",
		Username: "testctrl",
		Email:    "testctrl@test.com",
		Password: "testctrl",
	})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/user/testctrl", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.NotEmpty(t, resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_FindNoResult(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/api/user/testctrl", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Nil(t, resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_UpdateProfileSuccess(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id": "` + uid.String() + `",` +
		`"full_name" : "test controller updt",` +
		`"username" : "testctrlupdt",` +
		`"email" : "testctrlupdt@test.com"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/profile", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.NotEmpty(t, resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_UpdateProfileNotFound(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id": "` + uid.String() + `",` +
		`"full_name" : "test controller updt",` +
		`"username" : "testctrlupdt",` +
		`"email" : "testctrlupdt@test.com"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/profile", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_UpdateProfileBadFormat(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id": "` + uid.String() + `",` +
		`"full_name" : "",` +
		`"username" : "pdt",` +
		`"email" : "testctrlupdttest.com"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/profile", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	t.Log(string(resBody))
}

func TestUserControllerImpl_UpdatePasswordSuccess(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updatePasswordRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "testctrl",` +
		`"email"        : "testctrl@test.com",` +
		`"old_password" : "testctrl",` +
		`"new_password" : "testctrlupdate"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/password", strings.NewReader(updatePasswordRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.NotEmpty(t, resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
}

func TestUserControllerImpl_UpdatePassword(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	updatePasswordRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "testctrl",` +
		`"email"        : "testctrl@test.com",` +
		`"old_password" : "testctrl",` +
		`"new_password" : "testctrlupdate"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/password", strings.NewReader(updatePasswordRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	assert.Equal(t, "Record not found", resBodyModel.Error)
}

func TestUserControllerImpl_UpdatePasswordBadFormat(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updatePasswordRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "",` +
		`"email"        : "testctrltest.com",` +
		`"old_password" : "testctrl",` +
		`"new_password" : "testctrlupdate"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/password", strings.NewReader(updatePasswordRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
}

func TestUserControllerImpl_UpdatePasswordWrongOldPassword(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "testctrl",` +
		`"email"        : "testctrl@test.com",` +
		`"old_password" : "tasdasdctrl",` +
		`"new_password" : "testctrlupdate"` +
		"}"

	req := httptest.NewRequest(http.MethodPut, "/api/user/password", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	assert.Equal(t, "Old password didn't match", resBodyModel.Error)
}

func TestUserControllerImpl_DeleteSuccess(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "testctrl",` +
		`"email"        : "testctrl@test.com",` +
		`"password"     : "testctrl"` +
		"}"

	req := httptest.NewRequest(http.MethodDelete, "/api/user", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.NotEmpty(t, resBodyModel.Data)
	assert.Equal(t, "User Deleted Successfully", resBodyModel.Data)
	assert.Empty(t, resBodyModel.Error)
}

func TestUserControllerImpl_DeleteNotFound(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "testctrl",` +
		`"email"        : "testctrl@test.com",` +
		`"password"     : "testctrl"` +
		"}"

	req := httptest.NewRequest(http.MethodDelete, "/api/user", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
	assert.Equal(t, "Record not found", resBodyModel.Error)
}

func TestUserControllerImpl_DeleteWrongPassword(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "testctrl",` +
		`"email"        : "testctrl@test.com",` +
		`"password"     : "wrongpasws"` +
		"}"

	req := httptest.NewRequest(http.MethodDelete, "/api/user", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
}

func TestUserControllerImpl_DeleteBadFormat(t *testing.T) {
	_ = userRepository.DeleteAll(context.Background())
	uid, err := uuid.NewUUID()
	assert.NoError(t, err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte("testctrl"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = userRepository.Create(context.Background(), &domain.User{
		Id:        uid,
		FullName:  "test controller",
		Username:  "testctrl",
		Email:     "testctrl@test.com",
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	updateProfileRequest := "{" +
		`"id"           : "` + uid.String() + `",` +
		`"username"     : "tl",` +
		`"email"        : "testctrltest.com",` +
		`"password"     : "wrong"` +
		"}"

	req := httptest.NewRequest(http.MethodDelete, "/api/user", strings.NewReader(updateProfileRequest))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))

	assert.Nil(t, err)
	assert.NotNil(t, resBody)

	resBodyModel := &model.WebResponse{}
	err = json.Unmarshal(resBody, resBodyModel)
	assert.NoError(t, err)

	assert.Empty(t, resBodyModel.Data)
	assert.NotEmpty(t, resBodyModel.Error)
}
