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

func NewUserController(service service.UserService) UserController {
	return &userControllerImpl{Service: service}
}

func (controller *userControllerImpl) Register(c *gin.Context) {
	requestModel := &model.UserRegisterRequest{}
	err := c.BindJSON(requestModel)
	helper.PanicIfError(err)

	userResponse := controller.Service.Register(context.Background(), requestModel)

	jwtString, err := helper.GenerateJWT(userResponse)
	helper.PanicIfError(err)

	c.SetCookie(helper.JWTCookieName, jwtString, 60*60*24 /*24 hours*/, "/", "", false, true)
	c.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "User Created Successfully",
		Data:   userResponse,
	})
}

func (controller *userControllerImpl) Login(c *gin.Context) {
	requestModel := &model.UserLoginRequest{}
	err := c.BindJSON(requestModel)
	helper.PanicIfError(err)

	userResponse := controller.Service.Login(context.Background(), requestModel)

	jwtString, err := helper.GenerateJWT(userResponse)
	helper.PanicIfError(err)

	c.SetCookie(helper.JWTCookieName, jwtString, 60*60*24 /*24 hours*/, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	})
}

func (controller *userControllerImpl) Find(c *gin.Context) {
	requestModel := &model.UserFindRequest{
		Username: c.Params.ByName("key"),
		Email:    c.Params.ByName("key"),
		FullName: c.Params.ByName("key"),
	}

	response := controller.Service.Find(context.Background(), requestModel)

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (controller *userControllerImpl) UpdateProfile(c *gin.Context) {
	requestModel := &model.UserUpdateProfileRequest{}
	err := c.BindJSON(requestModel)
	helper.PanicIfError(err)

	requestModel.Id = uuid.MustParse(c.Request.Header.Get("Uid"))
	response := controller.Service.UpdateProfile(context.Background(), requestModel)

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (controller *userControllerImpl) UpdatePassword(c *gin.Context) {
	requestModel := &model.UserUpdatePasswordRequest{}
	err := c.BindJSON(requestModel)
	helper.PanicIfError(err)

	requestModel.Id = uuid.MustParse(c.Request.Header.Get("Uid"))
	controller.Service.UpdatePassword(context.Background(), requestModel)

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		//Data:   response,
	})
}

func (controller *userControllerImpl) Delete(c *gin.Context) {
	requestModel := &model.UserDeleteRequest{}
	err := c.BindJSON(requestModel)
	helper.PanicIfError(err)

	requestModel.Id = uuid.MustParse(c.Request.Header.Get("Uid"))
	controller.Service.Delete(context.Background(), requestModel)

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "User Deleted Successfully",
	})
}
