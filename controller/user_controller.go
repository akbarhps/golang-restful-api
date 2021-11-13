package controller

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	SetRoutes(r *gin.Engine)

	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Find(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
