package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go-api/app"
	"go-api/controller"
	"go-api/helper"
	"go-api/middleware"
	"go-api/repository"
	"go-api/service"
)

func main() {
	r := gin.Default()
	db := app.NewDatabase("prod")
	validate := validator.New()

	// middleware
	r.Use(middleware.JWTValidator())
	r.Use(gin.CustomRecovery(middleware.PanicHandler))

	// repository
	userRepository := repository.NewUserRepository(db)

	// service
	userService := service.NewUserService(validate, userRepository)

	// controller
	userController := controller.NewUserController(userService)

	// routes
	userController.SetRoutes(r)

	err := r.Run("localhost:3000")
	helper.PanicIfError(err)
}
