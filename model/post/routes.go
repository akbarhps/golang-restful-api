package post

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine, controller Controller) {
	userGroup := router.Group("/post")
	userGroup.GET("/", controller.FindByUserID)
	userGroup.GET("/:postID", controller.FindByPostID)
	userGroup.POST("/", controller.Create)
	userGroup.PUT("/:postID", controller.Update)
	userGroup.DELETE("/:postID", controller.Delete)
}
