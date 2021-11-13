package main

import (
	"go-api/app"
	"go-api/controller"
	"go-api/exception"
	"go-api/middleware"
	"go-api/repository"
	"go-api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func main() {
	db := app.NewDatabase("prod")
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(validator.New(), userRepository)
	userController := controller.NewUserController(userService)

	r := gin.Default()
	r.Use(middleware.JWTValidator())
	r.Use(gin.CustomRecovery(exception.PanicHandler))

	//TODO: Refactor routes
	r.GET("/api/user/:key", userController.Find)
	r.POST("/api/user", userController.Register)
	r.POST("/api/user/login", userController.Login)
	r.PUT("/api/user/profile", userController.UpdateProfile)
	r.PUT("/api/user/password", userController.UpdatePassword)
	r.DELETE("/api/user", userController.Delete)

	err := r.Run("localhost:3000")
	if err != nil {
		panic(err)
	}
}
