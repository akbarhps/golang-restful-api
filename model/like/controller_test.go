package like

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go-api/app"
	"go-api/middleware"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestControllerImpl_Create(t *testing.T) {
	app.TestDBInit()
	repository := NewRepository()
	service := NewService(validator.New(), repository)
	controller := NewController(service)

	router := gin.Default()
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	router.POST("/like/:postID", controller.Create)

	req := httptest.NewRequest("POST", "/like/12a30928-e553-4ab9-b245-984e6f238cb0", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc1NTY2MzAsImp0aSI6IjljYjBiMzJiLTkyMTQtNGI3Yy1iYzhlLWM0M2FiM2ZkMjE4OCIsImlhdCI6MTYzNzQ3MDIzMCwiaXNzIjoiaW5zdGFwb3VuZHMiLCJuYmYiOjE2Mzc0NzAyMzAsInN1YiI6InRlc3Rjb250cm9sbGVyIn0.pTwohCrzvYJNGef1iUXK6VuPJKwCGiQmQasIfq1zU6Q",
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	t.Log(string(body))
}

func TestControllerImpl_Delete(t *testing.T) {
	app.TestDBInit()
	repository := NewRepository()
	service := NewService(validator.New(), repository)
	controller := NewController(service)

	router := gin.Default()
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	router.DELETE("/like/:postID", controller.Delete)

	req := httptest.NewRequest("DELETE", "/like/12a30928-e553-4ab9-b245-984e6f238cb0", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc1NTY2MzAsImp0aSI6IjljYjBiMzJiLTkyMTQtNGI3Yy1iYzhlLWM0M2FiM2ZkMjE4OCIsImlhdCI6MTYzNzQ3MDIzMCwiaXNzIjoiaW5zdGFwb3VuZHMiLCJuYmYiOjE2Mzc0NzAyMzAsInN1YiI6InRlc3Rjb250cm9sbGVyIn0.pTwohCrzvYJNGef1iUXK6VuPJKwCGiQmQasIfq1zU6Q",
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	t.Log(string(body))
}
