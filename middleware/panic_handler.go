package middleware

import (
	"github.com/gin-gonic/gin"
	"go-api/exception"
	"go-api/model"
	"net/http"

	"github.com/go-playground/validator"
)

func badRequest(res *model.WebResponse, err string) {
	res.Code = http.StatusBadRequest
	res.Status = "Bad Request"
	res.Error = err
}

func notFound(res *model.WebResponse, err string) {
	res.Code = http.StatusNotFound
	res.Status = "Record Not Found"
	res.Error = err
}

func internalServerError(res *model.WebResponse, err string) {
	res.Code = http.StatusInternalServerError
	res.Status = "Internal Server Error"
	res.Error = err
}

func unauthorizedError(res *model.WebResponse, err string) {
	res.Code = http.StatusUnauthorized
	res.Status = "Unauthorized Access"
	res.Error = err
}

func PanicHandler(c *gin.Context, err interface{}) {
	res := &model.WebResponse{}

	switch err.(type) {
	case exception.RecordDuplicateError:
		badRequest(res, err.(exception.RecordDuplicateError).Error())
	case exception.RecordNotFoundError:
		notFound(res, err.(exception.RecordNotFoundError).Error())
	case validator.ValidationErrors:
		badRequest(res, err.(validator.ValidationErrors).Error())
	case exception.InvalidCredentialError:
		badRequest(res, err.(exception.InvalidCredentialError).Error())
	case exception.InvalidSignatureError:
		unauthorizedError(res, err.(exception.InvalidSignatureError).Error())
	case error:
		internalServerError(res, err.(error).Error())
	}

	c.AbortWithStatusJSON(res.Code, res)
}