package post_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"go-api/app"
	"go-api/middleware"
	"go-api/model/comment"
	"go-api/model/like"
	"go-api/model/post"
	"go-api/model/resource"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	app.TestDBInit()
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	postRepo := post.NewRepository()
	resourceRepo := resource.NewRepository()
	likeRepo := like.NewRepository()
	commentRepo := comment.NewRepository()

	postService := post.NewService(validator.New(), postRepo, resourceRepo, likeRepo, commentRepo)
	postController := post.NewController(postService)

	router.POST("/post", postController.Create)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("caption", "test kepseng 3")

	dat, err := os.ReadDir("../../res")
	assert.NoError(t, err)
	for _, f := range dat {
		if f.IsDir() {
			continue
		}

		file, err := os.Open("../../res/" + f.Name())
		assert.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("images[]", f.Name())
		assert.NoError(t, err)
		part.Write([]byte(file.Name()))
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/post", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc1ODQ5MDksImp0aSI6IjNiNTNiZjA2LWQzYTktNDdlNy1hOTAyLWFiODA1OTZhNzU3OCIsImlhdCI6MTYzNzQ5ODUwOSwiaXNzIjoiaW5zdGFwb3VuZHMiLCJuYmYiOjE2Mzc0OTg1MDksInN1YiI6InRlc3Rjb250cm9sbGVyIn0.YEhuW3v8XBP1RTrkM-Fbb4uYSs-AvkErKlz8BzSDz_A",
	})
	w := httptest.NewRecorder()

	t.Log(req.PostForm.Get("caption"))

	router.ServeHTTP(w, req)

	res := w.Result()
	resBody, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))
}

func TestGetByUserID(t *testing.T) {
	app.TestDBInit()
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	postRepo := post.NewRepository()
	resourceRepo := resource.NewRepository()
	likeRepo := like.NewRepository()
	commentRepo := comment.NewRepository()

	postService := post.NewService(validator.New(), postRepo, resourceRepo, likeRepo, commentRepo)
	postController := post.NewController(postService)

	router.GET("/post", postController.FindByUserID)

	req := httptest.NewRequest("GET", "/post?user_id=3b53bf06-d3a9-47e7-a902-ab80596a7578", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc1ODQ5MDksImp0aSI6IjNiNTNiZjA2LWQzYTktNDdlNy1hOTAyLWFiODA1OTZhNzU3OCIsImlhdCI6MTYzNzQ5ODUwOSwiaXNzIjoiaW5zdGFwb3VuZHMiLCJuYmYiOjE2Mzc0OTg1MDksInN1YiI6InRlc3Rjb250cm9sbGVyIn0.YEhuW3v8XBP1RTrkM-Fbb4uYSs-AvkErKlz8BzSDz_A",
	})
	w := httptest.NewRecorder()

	t.Log(req.PostForm.Get("caption"))

	router.ServeHTTP(w, req)

	res := w.Result()
	resBody, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))
}

func TestGetByPostID(t *testing.T) {
	app.TestDBInit()
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Use(middleware.JWTValidator())
	router.Use(gin.CustomRecovery(middleware.PanicHandler))

	postRepo := post.NewRepository()
	resourceRepo := resource.NewRepository()
	likeRepo := like.NewRepository()
	commentRepo := comment.NewRepository()

	postService := post.NewService(validator.New(), postRepo, resourceRepo, likeRepo, commentRepo)
	postController := post.NewController(postService)

	router.GET("/post/:postID", postController.FindByPostID)

	req := httptest.NewRequest("GET", "/post/0c6fb644-fd44-455b-b6d2-219b05b951ca", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc1NTY2MzAsImp0aSI6IjljYjBiMzJiLTkyMTQtNGI3Yy1iYzhlLWM0M2FiM2ZkMjE4OCIsImlhdCI6MTYzNzQ3MDIzMCwiaXNzIjoiaW5zdGFwb3VuZHMiLCJuYmYiOjE2Mzc0NzAyMzAsInN1YiI6InRlc3Rjb250cm9sbGVyIn0.pTwohCrzvYJNGef1iUXK6VuPJKwCGiQmQasIfq1zU6Q",
	})
	w := httptest.NewRecorder()

	t.Log(req.PostForm.Get("caption"))

	router.ServeHTTP(w, req)

	res := w.Result()
	resBody, _ := ioutil.ReadAll(res.Body)
	t.Log(string(resBody))
}
