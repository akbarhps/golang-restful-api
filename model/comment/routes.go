package comment

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.RouterGroup, controller Controller) {
	userGroup := router.Group("/comment")
	userGroup.POST("/:postID", controller.Create)
	userGroup.DELETE("/:postID", controller.Delete)
}
