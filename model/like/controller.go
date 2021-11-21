package like

import (
	"github.com/gin-gonic/gin"
	"go-api/model"
	"net/http"
)

type Controller interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type controllerImpl struct {
	service Service
}

func NewController(service Service) *controllerImpl {
	return &controllerImpl{service: service}
}

func (c *controllerImpl) Create(ctx *gin.Context) {
	postID := ctx.Param("postID")
	userID := ctx.GetHeader("User_id")
	req := &Request{
		PostID: postID,
		UserID: userID,
	}
	c.service.Create(ctx, req)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
	})
}

func (c *controllerImpl) Delete(ctx *gin.Context) {
	postID := ctx.Param("postID")
	userID := ctx.GetHeader("User_id")
	req := &Request{
		PostID: postID,
		UserID: userID,
	}
	c.service.Delete(ctx, req)
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
	})
}
