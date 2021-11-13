package controller

import (
	"bytes"
	"context"
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

var router *gin.Engine
var userController UserController
var userService service.UserService
var userRepository repository.UserRepository

type userTest struct {
	Id          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	OldPassword string    `json:"old_password"`
	NewPassword string    `json:"new_password"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	pathRegister       = "/api/register"
	pathLogin          = "/api/login"
	pathFind           = "/api/user/:key"
	pathUpdateProfile  = "/api/user/profile"
	pathUpdatePassword = "/api/user/password"
	pathDelete         = "/api/user"

	userTestValid = &userTest{
		FullName:    "test controller",
		Username:    "testctrl",
		Email:       "testctrl@test.com",
		Password:    "testctrl",
		OldPassword: "testctrl",
		NewPassword: "testctrl",
		CreatedAt:   time.Now(),
	}
	userTestUpdate = &userTest{
		FullName:    "test controller update",
		Username:    "testctrlupd",
		Email:       "testctrlupd@test.com",
		Password:    "testctrlupd",
		OldPassword: "testctrl",
		NewPassword: "testctrlupdt",
		CreatedAt:   time.Now(),
	}
	userTestWrongPassword = &userTest{
		FullName:    "test controller",
		Username:    "testctrl",
		Email:       "testctrl@test.com",
		Password:    "wrongpswd",
		OldPassword: "wrongpswd",
		NewPassword: "testctrl",
		CreatedAt:   time.Now(),
	}
	userTestInvalid = &userTest{
		Id:          uuid.UUID{},
		FullName:    "",
		Username:    "",
		Email:       "",
		Password:    "",
		OldPassword: "",
		NewPassword: "",
		CreatedAt:   time.Time{},
	}
)

func registerDummyUser() *model.UserResponse {
	response := userService.Register(context.Background(), &model.UserRegisterRequest{
		FullName: "test controller",
		Username: "testctrl",
		Email:    "testctrl@test.com",
		Password: "testctrl",
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
	router = gin.Default()
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	db := app.NewDatabase("test")
	userRepository = repository.NewUserRepository(db)
	userService = service.NewUserService(validator.New(), userRepository)
	userController = NewUserController(userService)

	router.GET(pathFind, userController.Find)
	router.POST(pathRegister, userController.Register)
	router.POST(pathLogin, userController.Login)
	router.PUT(pathUpdateProfile, userController.UpdateProfile)
	router.PUT(pathUpdatePassword, userController.UpdatePassword)
	router.DELETE(pathDelete, userController.Delete)

	m.Run()
}

func TestUserControllerImpl_Register(t *testing.T) {
	t.Run("register success should return user data and cookie", func(t *testing.T) {
		userRepository.DeleteAll(context.Background())
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, pathRegister, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPost, pathRegister, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		registerDummyUser()
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, pathRegister, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		registerDummyUser()
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		requestJSON := generateUserJSON(userTestValid)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		registerDummyUser()

		requestJSON := generateUserJSON(userTestWrongPassword)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
	path := strings.Replace(pathFind, ":key", "", 1)

	t.Run("find user should get ok, data array of users", func(t *testing.T) {
		userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, path+dummyResponse.Email, nil)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, path+"notregistered", nil)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON := generateUserJSON(userTestUpdate)

		req := httptest.NewRequest(http.MethodPut, pathUpdateProfile, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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

		assert.NotEqual(t, dummyResponse.FullName, webResponse.Data.(map[string]interface{})["full_name"])
		assert.NotEqual(t, dummyResponse.Username, webResponse.Data.(map[string]interface{})["username"])
		assert.NotEqual(t, dummyResponse.Email, webResponse.Data.(map[string]interface{})["email"])

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("update user with no signature (jwt) should get unauthorized and empty data", func(t *testing.T) {
		userRepository.DeleteAll(context.Background())
		req := httptest.NewRequest(http.MethodPut, pathUpdateProfile, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPut, pathUpdateProfile, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())

		requestJSON := generateUserJSON(userTestUpdate)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, pathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		req := httptest.NewRequest(http.MethodPut, pathUpdatePassword, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON := generateUserJSON(userTestInvalid)

		req := httptest.NewRequest(http.MethodPut, pathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())

		requestJSON := generateUserJSON(userTestWrongPassword)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, pathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		requestJSON := generateUserJSON(userTestValid)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, pathDelete, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("delete using no signature (jwt) should get unauthorized and error", func(t *testing.T) {
		userRepository.DeleteAll(context.Background())
		req := httptest.NewRequest(http.MethodDelete, pathDelete, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())
		requestJSON := generateUserJSON(userTestInvalid)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, pathDelete, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
		userRepository.DeleteAll(context.Background())

		requestJSON := generateUserJSON(userTestUpdate)
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, pathDelete, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

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
