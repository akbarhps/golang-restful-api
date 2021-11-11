package controller

import (
	"context"
	"go-api/exception"
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
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	response, err := controller.Service.Register(context.Background(), requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "User Created Successfully",
		Data:   response,
	})
}

func (controller *userControllerImpl) Login(c *gin.Context) {
	requestModel := &model.UserLoginRequest{}
	err := c.BindJSON(requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	response, err := controller.Service.Login(context.Background(), requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (controller *userControllerImpl) Find(c *gin.Context) {
	requestModel := &model.UserFindRequest{
		Username: c.Params.ByName("key"),
		Email:    c.Params.ByName("key"),
		FullName: c.Params.ByName("key"),
	}

	response, err := controller.Service.Find(context.Background(), requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (controller *userControllerImpl) UpdateProfile(c *gin.Context) {
	requestModel := &model.UserUpdateProfileRequest{}
	err := c.BindJSON(requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	response, err := controller.Service.UpdateProfile(context.Background(), requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (controller *userControllerImpl) UpdatePassword(c *gin.Context) {
	requestModel := &model.UserUpdatePasswordRequest{}
	err := c.BindJSON(requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	response, err := controller.Service.UpdatePassword(context.Background(), requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (controller *userControllerImpl) Delete(c *gin.Context) {
	requestModel := &model.UserDeleteRequest{}
	err := c.BindJSON(requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	err = controller.Service.Delete(context.Background(), requestModel)
	if err != nil {
		exception.PanicHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "User Deleted Successfully",
	})
}
