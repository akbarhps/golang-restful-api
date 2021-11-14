package controller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"go-api/app"
	"go-api/helper"
	"go-api/middleware"
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
)

var r *gin.Engine
var userController UserController
var userService service.UserService
var userRepository repository.UserRepository

type userTest struct {
	Id          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	OldPassword string    `json:"old_password"`
	NewPassword string    `json:"new_password"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	db = app.NewDatabase("test")

	userTestValid = &userTest{
		DisplayName: "test controller",
		Username:    "testctrl",
		Email:       "testctrl@test.com",
		Password:    "testctrl",
		OldPassword: "testctrl",
		NewPassword: "testctrl",
		CreatedAt:   time.Now(),
	}
	userTestUpdate = &userTest{
		DisplayName: "test controller update",
		Username:    "testctrlupd",
		Email:       "testctrlupd@test.com",
		Password:    "testctrlupd",
		OldPassword: "testctrl",
		NewPassword: "testctrlupdt",
		CreatedAt:   time.Now(),
	}
	userTestWrongPassword = &userTest{
		DisplayName: "test controller",
		Username:    "testctrl",
		Email:       "testctrl@test.com",
		Password:    "wrongpswd",
		OldPassword: "wrongpswd",
		NewPassword: "testctrl",
		CreatedAt:   time.Now(),
	}
	userTestInvalid = &userTest{
		Id:          uuid.UUID{}.String(),
		DisplayName: "",
		Username:    "",
		Email:       "",
		Password:    "",
		OldPassword: "",
		NewPassword: "",
		CreatedAt:   time.Time{},
	}
)

func clearRecord() *sql.Tx {
	tx, _ := db.Begin()
	userRepository.DeleteAll(context.Background(), tx)
	return tx
}

func registerDummyUser() *model.UserResponse {
	response := userService.Register(context.Background(), &model.UserRegister{
		DisplayName: "test controller",
		Username:    "testctrl",
		Email:       "testctrl@test.com",
		Password:    "testctrl",
	})
	return response
}

func generateJWTCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:    helper.JWTCookieName,
		Value:   token,
		Path:    "/",
		Domain:  "",
		Expires: jwt.TimeFunc().Add(time.Hour * 24),
	}
}

func generateUserJSON(userTest *userTest) *bytes.Reader {
	userJSON, _ := json.Marshal(userTest)
	return bytes.NewReader(userJSON)
}

func TestMain(m *testing.M) {
	r = gin.Default()
	r.Use(middleware.JWTValidator())
	r.Use(gin.CustomRecovery(middleware.PanicHandler))

	userRepository = repository.NewUserRepository()
	userService = service.NewUserService(db, validator.New(), userRepository)
	userController = NewUserController(userService)
	userController.SetRoutes(r)

	m.Run()
}

func TestUserControllerImpl_Register(t *testing.T) {
	t.Run("register success should return user data and cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, UserPathRegister, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotEmpty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("register using bad input should get bad request, empty data, and no cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPost, UserPathRegister, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("register using registered email or username should get bad request, empty data, and no cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		registerDummyUser()
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, UserPathRegister, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestUserControllerImpl_Login(t *testing.T) {
	t.Run("login using valid credential should get ok, user data and cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		registerDummyUser()
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, UserPathLogin, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.NotEmpty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("login using not registered user should get not found, empty data, and no cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, UserPathLogin, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)

		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("login using wrong password should get bad request, no data, and no cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		registerDummyUser()

		requestJSON := generateUserJSON(userTestWrongPassword)

		req := httptest.NewRequest(http.MethodPost, UserPathLogin, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("login using bad input should get bad request, no data and no cookie", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPost, UserPathLogin, requestJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestUserControllerImpl_Find(t *testing.T) {
	webResponse := &model.WebResponse{}
	path := strings.Replace(UserPathFind, ":key", "", 1)

	t.Run("find user should get ok, data array of users", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, path+dummyResponse.Email, nil)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("find user with no signature cookie (jwt) should get unauthorized and empty data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, path+"notregistered", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("find non-exist user should get ok and empty array data", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, path+"notregistered", nil)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestUserControllerImpl_UpdateProfile(t *testing.T) {
	t.Run("update profile with valid input should get ok and data of updated user", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON := generateUserJSON(userTestUpdate)

		req := httptest.NewRequest(http.MethodPut, UserPathUpdateProfile, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		assert.NotEqual(t, dummyResponse.DisplayName, webResponse.Data.(map[string]interface{})["full_name"])
		assert.NotEqual(t, dummyResponse.Username, webResponse.Data.(map[string]interface{})["username"])
		assert.NotEqual(t, dummyResponse.Email, webResponse.Data.(map[string]interface{})["email"])

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("update user with no signature (jwt) should get unauthorized and empty data", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		req := httptest.NewRequest(http.MethodPut, UserPathUpdateProfile, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("update profile with bad format should get bad request and empty data", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPut, UserPathUpdateProfile, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

}

func TestUserControllerImpl_UpdatePassword(t *testing.T) {
	t.Run("update password success should get ok and no error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)

		requestJSON := generateUserJSON(userTestUpdate)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, UserPathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("update password with no signature (jwt) should get unauthorized and error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, UserPathUpdatePassword, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("update password with bad format input should get bad request and error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPut, UserPathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("update password using wrong old password should get bad request and error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)

		requestJSON := generateUserJSON(userTestWrongPassword)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, UserPathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestUserControllerImpl_Delete(t *testing.T) {
	t.Run("delete using valid input should get ok and no error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		requestJSON := generateUserJSON(userTestValid)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, UserPathDelete, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("delete using no signature (jwt) should get unauthorized and error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		req := httptest.NewRequest(http.MethodDelete, UserPathDelete, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("delete using invalid input should get bad request and error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)
		requestJSON := generateUserJSON(userTestInvalid)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, UserPathDelete, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("delete using wrong old password should get bad request and error", func(t *testing.T) {
		tx := clearRecord()
		helper.TXCommitOrRollback(tx)

		requestJSON := generateUserJSON(userTestUpdate)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, UserPathDelete, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		var webResponse model.WebResponse
		err = json.Unmarshal(resBody, &webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}
