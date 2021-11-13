package middleware

import (
	"github.com/gin-gonic/gin"
	"go-api/exception"
	"go-api/helper"
	"net/http"
	"strings"
)

func JWTValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.FullPath(), "register") || strings.Contains(c.FullPath(), "login") {
			c.Next()
			return
		}

		key, err := c.Cookie(helper.JWTCookieName)
		if err != nil {
			PanicHandler(c, exception.InvalidSignatureError{Message: "No Signature Provided"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		payload, err := helper.ValidateJWT(key)
		if err != nil {
			PanicHandler(c, exception.InvalidSignatureError{Message: err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Request.Header.Set("Uid", payload.Id)
		c.Next()
	}
}
