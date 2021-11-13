package controller

import (
	"context"
	"github.com/google/uuid"
	"go-api/helper"
	"go-api/model"
	"go-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userControllerImpl struct {
	Service service.UserService
}

var (
	UserPathRegister       = "/api/user/register"
	UserPathLogin          = "/api/user/login"
	UserPathFind           = "/api/user/:key"
	UserPathUpdateProfile  = "/api/user"
	UserPathUpdatePassword = "/api/user/password"
	UserPathDelete         = "/api/user"
)

func NewUserController(service service.UserService) UserController {
	return &userControllerImpl{Service: service}
}

func (c *userControllerImpl) SetRoutes(r *gin.Engine) {
	r.GET(UserPathFind, c.Find)
	r.PUT(UserPathUpdateProfile, c.UpdateProfile)
	r.PUT(UserPathUpdatePassword, c.UpdatePassword)
	r.POST(UserPathLogin, c.Login)
	r.POST(UserPathRegister, c.Register)
	r.DELETE(UserPathDelete, c.Delete)
}

func (c *userControllerImpl) Register(ctx *gin.Context) {
	requestModel := &model.UserRegister{}
	err := ctx.BindJSON(requestModel)
	helper.PanicIfError(err)

	userResponse := c.Service.Register(context.Background(), requestModel)

	jwtString, err := helper.GenerateJWT(userResponse)
	helper.PanicIfError(err)

	ctx.SetCookie(helper.JWTCookieName, jwtString, 60*60*24 /*24 hours*/, "/", "", false, true)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "User Created Successfully",
		Data:   userResponse,
	})
}

func (c *userControllerImpl) Login(ctx *gin.Context) {
	requestModel := &model.UserLogin{}
	err := ctx.BindJSON(requestModel)
	helper.PanicIfError(err)

	userResponse := c.Service.Login(context.Background(), requestModel)

	jwtString, err := helper.GenerateJWT(userResponse)
	helper.PanicIfError(err)

	ctx.SetCookie(helper.JWTCookieName, jwtString, 60*60*24 /*24 hours*/, "/", "", false, true)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	})
}

func (c *userControllerImpl) Find(ctx *gin.Context) {
	requestModel := &model.UserFind{
		Username: ctx.Params.ByName("key"),
		Email:    ctx.Params.ByName("key"),
		FullName: ctx.Params.ByName("key"),
	}

	response := c.Service.Find(context.Background(), requestModel)

	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (c *userControllerImpl) UpdateProfile(ctx *gin.Context) {
	requestModel := &model.UserUpdateProfile{}
	err := ctx.BindJSON(requestModel)
	helper.PanicIfError(err)

	requestModel.Id = uuid.MustParse(ctx.Request.Header.Get("Uid"))
	response := c.Service.UpdateProfile(context.Background(), requestModel)

	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (c *userControllerImpl) UpdatePassword(ctx *gin.Context) {
	requestModel := &model.UserUpdatePassword{}
	err := ctx.BindJSON(requestModel)
	helper.PanicIfError(err)

	requestModel.Id = uuid.MustParse(ctx.Request.Header.Get("Uid"))
	c.Service.UpdatePassword(context.Background(), requestModel)

	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	})
}

func (c *userControllerImpl) Delete(ctx *gin.Context) {
	requestModel := &model.UserDelete{}
	err := ctx.BindJSON(requestModel)
	helper.PanicIfError(err)

	requestModel.Id = uuid.MustParse(ctx.Request.Header.Get("Uid"))
	c.Service.Delete(context.Background(), requestModel)

	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	})
}
