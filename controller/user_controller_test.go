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

var (
	pathRegister       = "/api/register"
	pathLogin          = "/api/login"
	pathFind           = "/api/user/:key"
	pathUpdateProfile  = "/api/user/profile"
	pathUpdatePassword = "/api/user/password"
	pathDelete         = "/api/user"

	modeUpdate           = "update"
	modeBadFormat        = "bad_format"
	modeGoodFormat       = "good_format"
	modeWrongOldPassword = "wrong_old_password"
)

func registerDummyUser() *model.UserResponse {
	response, _ := userService.Register(context.Background(), &model.UserRegisterRequest{
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

func generateUserJSON(mode string) *bytes.Reader {
	type UserTest struct {
		Id          uuid.UUID `json:"id"`
		FullName    string    `json:"full_name"`
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		Password    string    `json:"password"`
		OldPassword string    `json:"old_password"`
		NewPassword string    `json:"new_password"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var user UserTest
	if mode == modeBadFormat {
		user = UserTest{
			FullName:  "troller",
			Username:  "trl",
			Email:     "trltest.com",
			Password:  "trl",
			CreatedAt: time.Now(),
		}
	} else if mode == modeWrongOldPassword {
		user = UserTest{
			FullName:    "test controller update",
			Username:    "testctrlupd",
			Email:       "testctrlupd@test.com",
			Password:    "testctrlupd",
			OldPassword: "wrongpasswrd",
			NewPassword: "testctrlupdt",
			CreatedAt:   time.Now(),
		}
	} else if mode == modeUpdate {
		user = UserTest{
			FullName:    "test controller update",
			Username:    "testctrlupd",
			Email:       "testctrlupd@test.com",
			Password:    "testctrlupd",
			OldPassword: "testctrl",
			NewPassword: "testctrlupdt",
			CreatedAt:   time.Now(),
		}
	} else {
		user = UserTest{
			FullName:  "test controller",
			Username:  "testctrl",
			Email:     "testctrl@test.com",
			Password:  "testctrl",
			CreatedAt: time.Now(),
		}
	}

	userJSON, _ := json.Marshal(user)
	return bytes.NewReader(userJSON)
}

func TestMain(m *testing.M) {
	router = gin.Default()
	router.Use(middleware.JWTValidator())

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
	requestJSON := generateUserJSON(modeGoodFormat)
	webResponse := &model.WebResponse{}

	t.Run("Register Success", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeGoodFormat)

		req := httptest.NewRequest(http.MethodPost, pathRegister, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotEmpty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Register Bad Format Input", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeBadFormat)

		req := httptest.NewRequest(http.MethodPost, pathRegister, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Duplicate Email Or Password", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		registerDummyUser()
		requestJSON = generateUserJSON(modeGoodFormat)

		req := httptest.NewRequest(http.MethodPost, pathRegister, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestUserControllerImpl_Login(t *testing.T) {
	requestJSON := generateUserJSON(modeGoodFormat)
	webResponse := &model.WebResponse{}

	t.Run("Login Success", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		registerDummyUser()
		requestJSON = generateUserJSON(modeGoodFormat)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.NotEmpty(t, res.Cookies())

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

	t.Run("Login User Not Registered", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeGoodFormat)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)

		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)
		assert.Equal(t, "Record not found", webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Login Wrong Password", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		registerDummyUser()
		requestJSON = generateUserJSON(modeBadFormat)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Login Bad Format Input", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeBadFormat)

		req := httptest.NewRequest(http.MethodPost, pathLogin, requestJSON)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Empty(t, res.Cookies())

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
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

	t.Run("Find Success", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
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

	t.Run("Find With No Signature (JWT)", func(t *testing.T) {
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

	t.Run("Find Not Found", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
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
	requestJSON := generateUserJSON(modeGoodFormat)
	webResponse := &model.WebResponse{}

	t.Run("Update Profile Success", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON = generateUserJSON(modeUpdate)

		req := httptest.NewRequest(http.MethodPut, pathUpdateProfile, requestJSON)
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

		assert.NotEqual(t, dummyResponse.FullName, webResponse.Data.(map[string]interface{})["full_name"])
		assert.NotEqual(t, dummyResponse.Username, webResponse.Data.(map[string]interface{})["username"])
		assert.NotEqual(t, dummyResponse.Email, webResponse.Data.(map[string]interface{})["email"])

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Update Profile No Signature Key (JWT)", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		req := httptest.NewRequest(http.MethodPut, pathUpdateProfile, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Update Profile Bad Format", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON = generateUserJSON(modeBadFormat)

		req := httptest.NewRequest(http.MethodPut, pathUpdateProfile, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		assert.NotNil(t, resBody)

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Empty(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

}

func TestUserControllerImpl_UpdatePassword(t *testing.T) {
	requestJSON := generateUserJSON(modeUpdate)
	webResponse := &model.WebResponse{}

	t.Run("Update Password Success", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
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

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Update Password No Signature (JWT)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, pathUpdatePassword, nil)
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

	t.Run("Update Password Bad Format Input", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		dummyResponse := registerDummyUser()
		token, err := helper.GenerateJWT(dummyResponse)
		assert.NoError(t, err)

		requestJSON = generateUserJSON(modeBadFormat)

		req := httptest.NewRequest(http.MethodPut, pathUpdatePassword, requestJSON)
		req.AddCookie(generateJWTCookie(token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

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

	t.Run("Update Password Wrong Old Password", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeWrongOldPassword)
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

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestUserControllerImpl_Delete(t *testing.T) {
	requestJSON := generateUserJSON(modeUpdate)
	webResponse := &model.WebResponse{}

	t.Run("Delete Success", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeGoodFormat)
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

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, webResponse.Data)
		assert.Empty(t, webResponse.Error)

		assert.Equal(t, "User Deleted Successfully", webResponse.Data)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Delete No Signature (JWT)", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		req := httptest.NewRequest(http.MethodDelete, pathDelete, nil)
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

	t.Run("Delete Bad Format Input", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeBadFormat)
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

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("Delete Wrong Password", func(t *testing.T) {
		_ = userRepository.DeleteAll(context.Background())
		requestJSON = generateUserJSON(modeWrongOldPassword)
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

		err = json.Unmarshal(resBody, webResponse)
		assert.NoError(t, err)
		assert.Nil(t, webResponse.Data)
		assert.NotEmpty(t, webResponse.Error)

		assert.Equal(t, "Incorrect password", webResponse.Error)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}
