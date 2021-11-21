package post

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-api/model"
	"go-api/model/resource"
	"net/http"
)

type Controller interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByUserID(ctx *gin.Context)
	FindByPostID(ctx *gin.Context)
}

type controllerImpl struct {
	service Service
}

func NewController(service Service) Controller {
	return &controllerImpl{service: service}
}

func (c *controllerImpl) Create(ctx *gin.Context) {
	var req *CreateRequest
	err := ctx.ShouldBindWith(&req, binding.Form)
	if err != nil {
		panic(err)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		panic(err)
	}

	files := form.File["images[]"]
	if len(files) == 0 {
		panic("can't post with empty image")
	}

	for i, file := range files {
		filePath := "../../res/posts/" + file.Filename
		err = ctx.SaveUploadedFile(file, filePath)
		if err != nil {
			panic(err)
		}

		req.Resources = append(req.Resources, resource.Resource{
			Path:        filePath,
			ShareURL:    "localhost:3000/" + filePath,
			IndexInPost: i,
		})
	}

	req.UserID = ctx.GetHeader("User_id")
	res := c.service.Create(context.Background(), req)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
		Data:   res,
	})
}

func (c *controllerImpl) Update(ctx *gin.Context) {
	var req *UpdateRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	req.UserID = ctx.GetHeader("User_id")
	c.service.Update(context.Background(), req)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
	})
}

func (c *controllerImpl) Delete(ctx *gin.Context) {
	var req *DeleteRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	req.UserID = ctx.GetHeader("User_id")
	c.service.Delete(context.Background(), req)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
	})
}

func (c *controllerImpl) FindByUserID(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	res := c.service.FindByUserID(context.Background(), userID)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
		Data:   res,
	})
}

func (c *controllerImpl) FindByPostID(ctx *gin.Context) {
	postID := ctx.Param("postID")
	viewerID := ctx.GetHeader("User_id")
	res := c.service.FindByPostID(context.Background(), postID, viewerID)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
		Data:   res,
	})
}
