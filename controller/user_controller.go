package controller

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	SetRoutes(r *gin.Engine)
	Register(c *gin.Context)
	Login(c *gin.Context)
	Find(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdatePassword(c *gin.Context)
	Delete(c *gin.Context)
}
