package user_test

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"go-api/app"
	"go-api/helper"
	"go-api/middleware"
	"go-api/model"
	"go-api/model/user"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	router     *gin.Engine
	repository user.Repository
	service    user.Service
	controller user.Controller

	registerValid = &user.RegisterRequest{
		Email:       "testcontroller@test.com",
		Username:    "testcontroller",
		DisplayName: "test controller",
		Password:    "testcontroller",
	}

	registerEmpty = &user.RegisterRequest{
		Email:       "",
		Username:    "",
		DisplayName: "",
		Password:    "",
	}

	loginValid = &user.LoginRequest{
		Handler:  "testcontroller",
		Password: "testcontroller",
	}

	loginWrongPassword = &user.LoginRequest{
		Handler:  "testcontroller",
		Password: "testcontroller1",
	}

	loginNotRegistered = &user.LoginRequest{
		Handler:  "notregistered",
		Password: "notregistered",
	}

	loginEmpty = &user.LoginRequest{
		Handler:  "",
		Password: "",
	}

	updateProfileValid = &user.UpdateProfileRequest{
		UserID:      "",
		Email:       "controllerupdate@test.com",
		Username:    "controllerupdate",
		DisplayName: "controller update",
		Biography:   "controller update",
	}

	updateProfileEmpty = &user.UpdateProfileRequest{
		UserID:      "",
		Email:       "",
		Username:    "",
		DisplayName: "",
		Biography:   "",
	}

	updatePasswordValid = &user.UpdatePasswordRequest{
		UserID:      "",
		OldPassword: "testcontroller",
		NewPassword: "controllerupdate1",
	}

	updatePasswordWrongOldPassword = &user.UpdatePasswordRequest{
		UserID:      "",
		OldPassword: "controllerupdate1",
		NewPassword: "controllerupdate",
	}

	updatePasswordEmpty = &user.UpdatePasswordRequest{
		UserID:      "",
		OldPassword: "",
		NewPassword: "",
	}
)

func setupControllerTest() {
	app.TestDBInit()
	repository = user.NewRepository()
	service = user.NewService(validator.New(), repository)
	controller = user.NewController(service)

	router = gin.Default()
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	router.POST("/user/login", controller.Login)
	router.POST("/user/register", controller.Register)

	router.GET("/user", controller.Search)
	router.GET("/user/:username", controller.FindByUsername)

	router.PUT("/user/edit", controller.UpdateProfile)
	router.PUT("/user/password", controller.UpdatePassword)

	app.GetDB().Exec("DELETE FROM users")
}

func TestControllerImpl_Register(t *testing.T) {
	t.Run("success should return uid, jwt, and set jwt cookie", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"POST",
			"/user/register",
			helper.StructToJSONReader(registerValid),
		)
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("empty input should return bad request", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"POST",
			"/user/register",
			helper.StructToJSONReader(registerEmpty),
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("taken email or username should return bad request", func(t *testing.T) {
		setupControllerTest()
		service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"POST",
			"/user/register",
			helper.StructToJSONReader(registerValid),
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestControllerImpl_Login(t *testing.T) {
	t.Run("success should return user_id, token and set cookie", func(t *testing.T) {
		setupControllerTest()
		service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"POST",
			"/user/login",
			helper.StructToJSONReader(loginValid),
		)
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("not registered username or password should return not found", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"POST",
			"/user/login",
			helper.StructToJSONReader(loginNotRegistered),
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("wrong password should return bad request", func(t *testing.T) {
		setupControllerTest()
		service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"POST",
			"/user/login",
			helper.StructToJSONReader(loginWrongPassword),
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("empty input should return bad request", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"POST",
			"/user/login",
			helper.StructToJSONReader(loginEmpty),
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestControllerImpl_UpdateProfile(t *testing.T) {
	t.Run("success should return ok", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"PUT",
			"/user/edit",
			helper.StructToJSONReader(updateProfileValid),
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("no token (jwt) should return unauthorized", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"PUT",
			"/user/edit",
			helper.StructToJSONReader(updateProfileValid),
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("empty input should return bad request", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"PUT",
			"/user/edit",
			helper.StructToJSONReader(updateProfileEmpty),
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestControllerImpl_UpdatePassword(t *testing.T) {
	t.Run("success should return ok", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"PUT",
			"/user/password",
			helper.StructToJSONReader(updatePasswordValid),
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("no token (jwt) should return unauthorized", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"PUT",
			"/user/password",
			nil,
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("empty input should return bad request", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"PUT",
			"/user/password",
			helper.StructToJSONReader(updatePasswordEmpty),
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("wrong old password should return bad request", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"PUT",
			"/user/password",
			helper.StructToJSONReader(updatePasswordWrongOldPassword),
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestControllerImpl_FindByUsername(t *testing.T) {
	t.Run("success should return data user", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"GET",
			"/user/testcontroller",
			nil,
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("no token (jwt) should return unauthorized", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"GET",
			"/user/testcontroller",
			nil,
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("not found username should return not found", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"GET",
			"/user/notfoundusername",
			nil,
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

func TestControllerImpl_Search(t *testing.T) {
	t.Run("success should return array data of users", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"GET",
			"/user?handler=test",
			nil,
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("no token (jwt) should return unauthorized", func(t *testing.T) {
		setupControllerTest()
		req := httptest.NewRequest(
			"GET",
			"/user?handler=test",
			nil,
		)
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
		assert.NotEmpty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})

	t.Run("not found should return ok with empty array data", func(t *testing.T) {
		setupControllerTest()
		cred := service.Register(context.Background(), registerValid)
		req := httptest.NewRequest(
			"GET",
			"/user?handler=notfound",
			nil,
		)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: cred.Token,
		})
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
		assert.Empty(t, webResponse.Errors)

		t.Log(res.Cookies())
		t.Log(webResponse)
	})
}

