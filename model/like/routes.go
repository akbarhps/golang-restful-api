package like

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine, controller Controller) {
	userGroup := router.Group("/like")
	userGroup.POST("/:postID", controller.Create)
	userGroup.DELETE("/:postID", controller.Delete)
}
