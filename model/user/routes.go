package user

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine, controller Controller) {
	userGroup := router.Group("/user")
	userGroup.POST("/", controller.Register)
	userGroup.PUT("/edit/", controller.UpdateProfile)
	userGroup.PUT("/password/", controller.UpdatePassword)

	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
}
