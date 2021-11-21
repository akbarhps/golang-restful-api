package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-api/model"
	"net/http"
)

type Controller interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	Search(ctx *gin.Context)
	FindByUsername(ctx *gin.Context)
}

type controllerImpl struct {
	service Service
}

func NewController(service Service) Controller {
	return &controllerImpl{service: service}
}

func (c *controllerImpl) Register(ctx *gin.Context) {
	var req *RegisterRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	res := c.service.Register(context.Background(), req)
	ctx.SetCookie("token", res.Token, 3600, "/", "", false, false)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
		Data:   res,
	})
}

func (c *controllerImpl) Login(ctx *gin.Context) {
	var req *LoginRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	res := c.service.Login(context.Background(), req)
	ctx.SetCookie("token", res.Token, 3600, "/", "", false, false)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   res,
	})
}

func (c *controllerImpl) UpdateProfile(ctx *gin.Context) {
	var req *UpdateProfileRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	req.UserID = ctx.Request.Header.Get("User_id")
	c.service.UpdateProfile(context.Background(), req)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
	})
}

func (c *controllerImpl) UpdatePassword(ctx *gin.Context) {
	var req *UpdatePasswordRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	req.UserID = ctx.Request.Header.Get("User_id")
	c.service.UpdatePassword(context.Background(), req)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
	})
}

func (c *controllerImpl) Search(ctx *gin.Context) {
	keyword := ctx.Query("handler")
	users := c.service.SearchLike(context.Background(), keyword)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   users,
	})
}

func (c *controllerImpl) FindByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	user := c.service.FindByUsername(context.Background(), username)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
	})
}
