package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func NewController(service Service) Controller {
	return &controllerImpl{service: service}
}

func (c *controllerImpl) Create(ctx *gin.Context) {
	var req *CreateRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		panic(err)
	}

	req.PostID = ctx.Param("postID")
	req.UserID = ctx.GetHeader("User_id")
	c.service.Create(ctx, req)
	ctx.IndentedJSON(http.StatusCreated, &model.WebResponse{
		Code:   http.StatusCreated,
		Status: "ok",
	})
}

func (c *controllerImpl) Delete(ctx *gin.Context) {
	postID := ctx.Param("postID")
	userID := ctx.GetHeader("User_id")
	c.service.Delete(ctx, &DeleteRequest{
		PostID: postID,
		UserID: userID,
	})
	ctx.IndentedJSON(http.StatusOK, &model.WebResponse{
		Code:   http.StatusOK,
		Status: "ok",
	})
}
