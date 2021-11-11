package main

import (
	"go-api/app"
	"go-api/controller"
	"go-api/repository"
	"go-api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func main() {
	db := app.NewDatabase("prod")
	userRepository := repository.NewUserRepositoryImpl(db)
	userService := service.NewUserService(validator.New(), userRepository)
	userController := controller.NewUserController(userService)

	r := gin.Default()

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
